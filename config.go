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

// Config is config for mechat pay, all fileds is required.
type Config struct {
	AppId string    `json:"app_id"`
	MchId string    `json:"mch_id"`
	Cert  CertSuite `json:"cert"`

	Apiv3Secret string `json:"api_secret"`
}

// CertSuite is the suite for api cert
type CertSuite struct {
	SerialNo       string
	PrivateKeyTxt  string
	PrivateKeyPath string
}

// Option is optional configuration for mechat pay.
type Option func(o *options)

type options struct {
	domain  string
	schema  string
	certUrl string
}

func defaultOptions() options {
	return options{
		schema:  defaultSchema,
		domain:  defaultDomain,
		certUrl: defaultDomain + "/v3/certificates",
	}
}

const defaultSchema = "WECHATPAY2-SHA256-RSA2048"
const defaultDomain = "https://api.mch.weixin.qq.com"
