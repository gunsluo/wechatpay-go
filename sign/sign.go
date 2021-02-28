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

// Package sign implements signature and verify for wechat pay. It
// includes all encryption and decryption related implementations.
package sign

import (
	"bytes"
	"crypto/rsa"
	"net/url"
	"strconv"
	"time"
)

// RequestSignature is request signature information.
// The format as shown below:
// HTTP Method\nURL\nTimestamp\nNonce string\nHTTP Body\n
type RequestSignature struct {
	Method    string
	Url       string
	Timestamp int64
	Nonce     string
	Body      []byte
}

// Marshal returns the array byte about the request signature.
func (r *RequestSignature) Marshal() ([]byte, error) {
	u, err := url.Parse(r.Url)
	if err != nil {
		return nil, err
	}
	uri := u.Path
	if u.RawQuery != "" {
		uri += "?" + u.RawQuery
	}

	var b bytes.Buffer
	b.WriteString(r.Method)
	b.WriteString("\n")
	b.WriteString(uri)
	b.WriteString("\n")
	b.WriteString(strconv.FormatInt(r.Timestamp, 10))
	b.WriteString("\n")
	b.WriteString(r.Nonce)
	b.WriteString("\n")
	if len(r.Body) > 0 {
		b.Write(r.Body)
	}
	b.WriteString("\n")

	return b.Bytes(), nil
}

// NewRequestSignature return a request signature
func NewRequestSignature(method, url string, body []byte) *RequestSignature {
	return &RequestSignature{
		Method:    method,
		Timestamp: time.Now().Unix(),
		Url:       url,
		Nonce:     randomHex(32),
		Body:      body,
	}
}

// ResponseSignature is response signature information
// from the response of wechat pay.
// The format as shown below:
// Timestamp\nNonce string\nHTTP Body\n
type ResponseSignature struct {
	Body      []byte
	Timestamp int64
	Nonce     string
}

// Marshal returns the array byte about the response signature.
func (r *ResponseSignature) Marshal() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(strconv.FormatInt(r.Timestamp, 10))
	b.WriteString("\n")
	b.WriteString(r.Nonce)
	b.WriteString("\n")
	if len(r.Body) > 0 {
		b.Write(r.Body)
	}
	b.WriteString("\n")

	return b.Bytes(), nil
}

// GenerateSignature generate a signature string,
// privateKey is an RSA key.
func GenerateSignature(privateKey *rsa.PrivateKey, reqSign *RequestSignature, mchId, serialNo string) (string, error) {
	reqSignature, err := reqSign.Marshal()
	if err != nil {
		return "", err
	}

	signature, err := SignatureSHA256WithRSA(privateKey, reqSignature)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	b.WriteString(`mchid="`)
	b.WriteString(mchId)
	b.WriteString(`",nonce_str="`)
	b.WriteString(reqSign.Nonce)
	b.WriteString(`",signature="`)
	b.WriteString(signature)
	b.WriteString(`",timestamp="`)
	b.WriteString(strconv.FormatInt(reqSign.Timestamp, 10))
	b.WriteString(`",serial_no="`)
	b.WriteString(serialNo)
	b.WriteString(`"`)
	return b.String(), nil
}

// VerifySignature verify that the signature is passed.
// privateKey is an RSA key.
func VerifySignature(publicKey *rsa.PublicKey, respSign *ResponseSignature, signature string) error {
	respSignature, err := respSign.Marshal()
	if err != nil {
		return err
	}

	return VerifySHA256WithRSA(publicKey, signature, respSignature)
}
