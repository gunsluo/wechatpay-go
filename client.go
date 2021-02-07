// Copyright The Wechat Pay Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wechatpay

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gunsluo/wechatpay-go/v3/sign"
)

// client is wechat pay client for api v3
type Client interface {
	Config() *Config
	Do(context.Context, string, string, ...interface{}) *Result
	VerifySignature(context.Context, *Result) error
}

type client struct {
	config     Config
	secrets    secrets
	privateKey *rsa.PrivateKey

	genRequestSignature func(string, string, []byte) *sign.RequestSignature
}

// NewClient creates a new client with configuration from cfg.
func NewClient(cfg Config, opts ...Option) (Client, error) {
	return newClient(cfg, opts...)
}

func newClient(cfg Config, opts ...Option) (*client, error) {
	c := &client{
		config: cfg,
	}
	c.config.opts = defaultOptions()
	for _, opt := range opts {
		opt(&c.config.opts)
	}

	c.secrets.clear()

	if c.config.AppId == "" {
		return nil, errors.New("AppId is required")
	}

	if c.config.MchId == "" {
		return nil, errors.New("MchId is required")
	}

	if c.config.Apiv3Secret == "" {
		return nil, errors.New("Apiv3 Secret is required")
	}

	if c.config.Cert.SerialNo == "" {
		return nil, errors.New("SerialNo is required")
	}

	if c.config.Cert.PrivateKeyTxt == "" &&
		c.config.Cert.PrivateKeyPath == "" {
		return nil, errors.New("private key txt and path have at least one of them")
	}

	// load api private cert
	if c.config.Cert.PrivateKeyTxt != "" {
		privateKey, err := sign.LoadRSAPrivateKeyFromTxt(c.config.Cert.PrivateKeyTxt)
		if err != nil {
			return nil, err
		}
		c.privateKey = privateKey
	} else {
		privateKey, err := sign.LoadRSAPrivateKeyFromFile(c.config.Cert.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
		c.privateKey = privateKey
	}

	c.genRequestSignature = genRequestSignature
	return c, nil
}

// Config return client config
func (c *client) Config() *Config {
	return &c.config
}

// Signature signature a request and return signature string
func (c *client) Signature(reqSign *sign.RequestSignature) (string, error) {
	signature, err := sign.GenerateSignature(c.privateKey,
		reqSign, c.config.MchId, c.config.Cert.SerialNo)
	if err != nil {
		return "", err
	}

	return c.config.opts.Schema + " " + signature, nil
}

// Do sends a request and returns a result.
func (c *client) Do(ctx context.Context, method, url string, req ...interface{}) *Result {
	// 1. serialize the request
	var reqBuffer []byte
	if len(req) > 0 && method != http.MethodGet && req[0] != nil {
		buffer, err := json.Marshal(req[0])
		if err != nil {
			return &Result{Err: err}
		}
		reqBuffer = buffer
	}
	reqSign := c.genRequestSignature(method, url, reqBuffer)

	// 2-5. get data from wechatpay side
	result := c.do(ctx, reqSign)
	if result.Err != nil {
		return result
	}

	// 6. do extra workflow
	if err := c.doExtraWorkflow(ctx, reqSign, result); err != nil {
		result.Err = err
		return result
	}

	// 7. verify the response
	if err := c.VerifySignature(ctx, result); err != nil {
		result.Err = err
	}

	return result
}

func (c *client) do(ctx context.Context, reqSign *sign.RequestSignature) *Result {
	var reader io.Reader
	if len(reqSign.Body) > 0 {
		reader = bytes.NewBuffer(reqSign.Body)
	}

	// 2. create a http request
	httpReq, err := http.NewRequest(reqSign.Method, reqSign.Url, reader)
	if err != nil {
		return &Result{Err: err}
	}

	// 3. signature the request
	authSign, err := c.Signature(reqSign)
	if err != nil {
		return &Result{Err: err}
	}

	httpReq.Header.Set("Authorization", authSign)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// 4. send the request
	client := &http.Client{
		Transport: c.config.opts.transport,
		Timeout:   c.config.opts.timeout,
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return &Result{Err: err}
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode >= http.StatusMultipleChoices {
		message, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return &Result{Err: err}
		}

		e := &Error{Status: httpResp.StatusCode}
		if err := json.Unmarshal(message, e); err != nil {
			return &Result{Err: err}
		}

		return &Result{Err: e}
	}

	// 5. read the response
	nonce := httpResp.Header.Get("Wechatpay-Nonce")
	signature := httpResp.Header.Get("Wechatpay-Signature")
	ts := httpResp.Header.Get("Wechatpay-Timestamp")
	serialNo := httpResp.Header.Get("Wechatpay-Serial")

	var timestamp int64
	if ts != "" {
		i, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return &Result{Err: err}
		}
		timestamp = i
	}

	var body []byte
	if httpResp.StatusCode != http.StatusNoContent {
		body, err = ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return &Result{Err: err}
		}
	}

	result := &Result{
		Body:      body,
		Timestamp: timestamp,
		Nonce:     nonce,
		Signature: signature,
		SerialNo:  serialNo,
	}

	return result
}

func (c *client) doExtraWorkflow(ctx context.Context, reqSign *sign.RequestSignature, result *Result) error {
	workflows := c.getExtraWorkflows(reqSign)
	for _, workflow := range workflows {
		if err := workflow(ctx, c, reqSign, result); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) getExtraWorkflows(reqSign *sign.RequestSignature) []extraWorkflow {
	var workflows []extraWorkflow

	// cert
	if reqSign.Method == http.MethodGet && reqSign.Url == c.config.opts.CertUrl {
		if workflow, ok := extraWorkflowsMapping["cert"]; ok {
			workflows = append(workflows, workflow)
		}
	}

	return workflows
}

type extraWorkflow func(context.Context, *client, *sign.RequestSignature, *Result) error

var extraWorkflowsMapping = map[string]extraWorkflow{
	"cert": upgradeCertWorkflow,
}

func upgradeCertWorkflow(ctx context.Context, c *client, reqSign *sign.RequestSignature, result *Result) error {
	resp := &CertificatesResponse{}
	if err := json.Unmarshal(result.Body, resp); err != nil {
		return err
	}

	apiv3Secret := []byte(c.Config().Apiv3Secret)
	for _, cert := range resp.Certificates {
		// using apiv3 secret decrypt cert
		certBuffer, err := sign.DecryptByAes256Gcm(
			apiv3Secret,
			[]byte(cert.Encrypt.Nonce),
			[]byte(cert.Encrypt.Associated),
			cert.Encrypt.CipherText)
		if err != nil {
			return err
		}

		publicKey, err := sign.LoadRSAPublicKeyFromCert(certBuffer)
		if err != nil {
			return err
		}

		c.secrets.add(cert.SerialNo, publicKey, c.Config().opts.refreshTime)
	}

	return nil
}

// VerifySignature verify the signature from wechat pay's responses
func (c *client) VerifySignature(ctx context.Context, result *Result) error {
	// check and download certificates
	if err := c.onceDownloadCertificates(ctx); err != nil {
		return err
	}

	publicKey := c.secrets.get(result.SerialNo)
	if publicKey == nil {
		return errors.New("certificate not found")
	}

	respSign := &sign.ResponseSignature{
		Body:      result.Body,
		Timestamp: result.Timestamp,
		Nonce:     result.Nonce,
	}

	return sign.VerifySignature(publicKey, respSign, result.Signature)
}

type ctxOnceDlCert struct{}

var ctxKeyOnceDlCert = ctxOnceDlCert{}

func (c *client) onceDownloadCertificates(ctx context.Context) error {
	// avoid infinite loops
	if v := ctx.Value(ctxKeyOnceDlCert); v != nil {
		return nil
	}
	ctx = context.WithValue(ctx, ctxKeyOnceDlCert, struct{}{})

	if !c.secrets.isUpgrade() {
		return nil
	}

	rs := c.Do(ctx, http.MethodGet, c.config.opts.CertUrl)
	if rs.Err != nil {
		return rs.Err
	}

	//if len(c.publicKeys) == 0 {
	//	return errors.New("no certificates are available")
	//}

	return nil
}

func genRequestSignature(method, url string, body []byte) *sign.RequestSignature {
	return sign.NewRequestSignature(method, url, body)
}

type secrets struct {
	mutex    sync.RWMutex
	deadline time.Time
	all      map[string]*rsa.PublicKey
}

func (s *secrets) isUpgrade() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.deadline.Before(time.Now()) {
		return true
	}

	return len(s.all) == 0
}

func (s *secrets) add(key string, val *rsa.PublicKey, d time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.all[key] = val
	s.deadline = time.Now().Add(d)
}

func (s *secrets) get(key string) *rsa.PublicKey {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	val, _ := s.all[key]
	return val
}

func (s *secrets) clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.all = make(map[string]*rsa.PublicKey)
	s.deadline = time.Now()
}
