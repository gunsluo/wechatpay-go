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
	"context"
	"net/http"
)

// CertificatesRequest is the request for certificates.
type CertificatesRequest struct {
}

// CertificatesResponse is the response for certificates.
type CertificatesResponse struct {
	Certificates []Certificate `json:"data"`
}

// Certificate is certificate information
type Certificate struct {
	SerialNo      string             `json:"serial_no"`
	EffectiveTime string             `json:"effective_time"`
	ExpireTime    string             `json:"expire_time"`
	Encrypt       EncryptCertificate `json:"encrypt_certificate"`
}

// EncryptCertificate is the information of encrypt certificate.
type EncryptCertificate struct {
	Algorithm  string `json:"algorithm"`
	Nonce      string `json:"nonce"`
	Associated string `json:"associated_data"`
	CipherText string `json:"ciphertext"`
}

// Do get certificates from wechat pay.
func (r *CertificatesRequest) Do(ctx context.Context, c Client) (*CertificatesResponse, error) {
	url := c.Config().Options().CertUrl

	resp := &CertificatesResponse{}
	if err := c.Do(ctx, http.MethodGet, url).Scan(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
