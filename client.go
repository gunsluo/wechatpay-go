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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gunsluo/wechatpay-go/v3/sign"
)

// Client is wechat pay client for api v3
type Client struct {
	config Config
	opts   options

	privateKey *rsa.PrivateKey
	publicKeys map[string]*rsa.PublicKey
}

// NewClient creates a new client with configuration from cfg.
func NewClient(cfg Config, opts ...Option) (*Client, error) {
	c := &Client{
		config:     cfg,
		opts:       defaultOptions(),
		publicKeys: make(map[string]*rsa.PublicKey),
	}
	for _, opt := range opts {
		opt(&c.opts)
	}

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

	return c, nil
}

// Signature signature a request and return signature string
func (c *Client) Signature(reqSign *sign.RequestSignature) (string, error) {
	signature, err := sign.GenerateSignature(c.privateKey,
		reqSign, c.config.MchId, c.config.Cert.SerialNo)
	if err != nil {
		return "", err
	}

	return c.opts.schema + " " + signature, nil
}

// Do sends a request and returns a result.
func (c *Client) Do(ctx context.Context, method, url string, req ...interface{}) *Result {
	isCertRequest := c.isCertificateRequest(method, url)
	if !isCertRequest {
		// check and load certificates
		if err := c.lazyLoadCertificates(ctx); err != nil {
			return &Result{Err: err}
		}
	}

	// 1. serialie the request
	var reqBuffer []byte
	if len(req) > 0 {
		buffer, err := json.Marshal(req[0])
		if err != nil {
			return &Result{Err: err}
		}
		reqBuffer = buffer
	}
	reqSign := sign.NewRequestSignature(method, url, reqBuffer)

	result := c.do(ctx, reqSign)
	if result.Err != nil {
		return result
	}

	// 6. verify the response
	if isCertRequest {
		// upgrade certs and then verify signature.
		if err := c.upgradeCertificate(result.Body); err != nil {
			return &Result{Err: err}
		}
	}

	if err := c.VerifySignature(result); err != nil {
		result.Err = err
	}

	return result
}

func (c *Client) do(ctx context.Context, reqSign *sign.RequestSignature) *Result {
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
	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return &Result{Err: err}
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		message, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return &Result{Err: err}
		}
		return &Result{Err: fmt.Errorf("error:%s, statusCode=%v", string(message), httpResp.StatusCode)}
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

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return &Result{Err: err}
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

func (c *Client) isCertificateRequest(method, url string) bool {
	if method == http.MethodGet && url == c.opts.certUrl {
		return true
	}
	return false
}

func (c *Client) lazyLoadCertificates(ctx context.Context) error {
	// TODO: maybe set a expried time for this
	if len(c.publicKeys) > 0 {
		return nil
	}

	rs := c.Do(ctx, http.MethodGet, c.opts.certUrl)
	if rs.Err != nil {
		return rs.Err
	}

	if len(c.publicKeys) == 0 {
		return errors.New("no certificates are available")
	}

	return nil
}

func (c *Client) upgradeCertificate(data []byte) error {
	resp := &CertificatesRespone{}
	if err := json.Unmarshal(data, resp); err != nil {
		return err
	}

	apiv3Secret := []byte(c.config.Apiv3Secret)
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
		c.publicKeys[cert.SerialNo] = publicKey
	}

	return nil
}

// VerifySignature verify the signature from wechat pay's responses
func (c *Client) VerifySignature(result *Result) error {
	publicKey, ok := c.publicKeys[result.SerialNo]
	if !ok {
		return errors.New("no cert")
	}

	respSign := &sign.ResponseSignature{
		Body:      result.Body,
		Timestamp: result.Timestamp,
		Nonce:     result.Nonce,
	}

	return sign.VerifySignature(publicKey, respSign, result.Signature)
}
