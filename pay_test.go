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
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestDoForPay(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req       *PayRequest
		resp      *PayResponse
		transport *mockTransport
		pass      bool
	}{
		{
			&PayRequest{
				AppId:       client.config.AppId,
				MchId:       client.config.MchId,
				Description: "for testing",
				OutTradeNo:  "forxxxxxxxxx",
				TimeExpire:  time.Now().Add(10 * time.Minute).Format(time.RFC3339),
				Attach:      "cipher code",
				NotifyUrl:   "https://luoji.live/notify",
				Amount: PayAmount{
					Total:    1,
					Currency: "CNY",
				},
				TradeType: Native,
			},
			&PayResponse{
				CodeUrl: "weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00",
			},
			nil,
			true,
		},
		{
			&PayRequest{
				AppId:       client.config.AppId,
				MchId:       client.config.MchId,
				Description: "for testing",
				OutTradeNo:  "forxxxxxxxxx",
				TimeExpire:  time.Now().Add(10 * time.Minute).Format(time.RFC3339),
				Attach:      "cipher code",
				NotifyUrl:   "https://luoji.live/notify",
				Amount: PayAmount{
					Total:    1,
					Currency: "CNY",
				},
			},
			&PayResponse{
				CodeUrl: "weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00",
			},
			nil,
			true,
		},
		{
			&PayRequest{
				AppId:       client.config.AppId,
				MchId:       client.config.MchId,
				Description: "for testing",
				OutTradeNo:  "forxxxxxxxxx",
				TimeExpire:  time.Now().Add(10 * time.Minute).Format(time.RFC3339),
				Attach:      "cipher code",
				NotifyUrl:   "https://luoji.live/notify",
				Amount: PayAmount{
					Total:    1,
					Currency: "CNY",
				},
			},
			&PayResponse{
				CodeUrl: "weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00",
			},
			&mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					var resp = &http.Response{
						StatusCode: http.StatusOK,
					}

					resp.Header = http.Header{}
					resp.Body = ioutil.NopCloser(strings.NewReader("{}"))
					return resp, nil
				},
			},
			false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		if c.transport != nil {
			client.config.opts.transport = c.transport
			client.publicKeys = make(map[string]*rsa.PublicKey)
		}

		resp, err := c.req.Do(ctx, client)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, err)
		}

		if err != nil {
			continue
		}

		if !reflect.DeepEqual(c.resp, resp) {
			t.Fatalf("expect %v, got %v", c.resp, resp)
		}
	}
}
