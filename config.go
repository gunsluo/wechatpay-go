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
	"net/http"
	"reflect"
	"time"
)

// Config is config for wechat pay, all fields is required.
type Config struct {
	AppId string
	MchId string
	Cert  CertSuite

	Apiv3Secret string
	opts        options
}

// CertSuite is the suite for api cert.
type CertSuite struct {
	SerialNo       string
	PrivateKeyTxt  string
	PrivateKeyPath string
}

// Option is optional configuration for wechat pay.
type Option func(o *options)

// Transport set transport to http client.
func Transport(transport http.RoundTripper) Option {
	return func(o *options) {
		if transport == nil || reflect.ValueOf(transport).IsNil() {
			return
		}
		o.transport = transport
	}
}

// Timeout set timeout for http client.
func Timeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// CertRefreshTime set max cert refresh time, default
// value is 12h.
func CertRefreshTime(refreshTime time.Duration) Option {
	return func(o *options) {
		o.refreshTime = refreshTime
	}
}

// Options return the options
func (c *Config) Options() *options {
	return &c.opts
}

type options struct {
	Domain  string
	Schema  string
	CertUrl string

	transport   http.RoundTripper
	timeout     time.Duration
	refreshTime time.Duration
}

func defaultOptions() options {
	return options{
		Schema:      defaultSchema,
		Domain:      defaultDomain,
		CertUrl:     defaultDomain + "/v3/certificates",
		refreshTime: 12 * time.Hour,
	}
}

const defaultSchema = "WECHATPAY2-SHA256-RSA2048"
const defaultDomain = "https://api.mch.weixin.qq.com"
