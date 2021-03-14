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
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestDoForCombinePay(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req       *CombinePayRequest
		resp      *CombinePayResponse
		transport *mockTransport
		pass      bool
	}{
		{
			&CombinePayRequest{
				OutTradeNo: "forxxxxxxxxx",
				TimeStart:  time.Now(),
				TimeExpire: time.Now().Add(10 * time.Minute),
				NotifyUrl:  "https://luoji.live/notify",
				Orders: []SubOrder{
					{
						MchId:  mockMchId,
						Attach: "cipher code",
						Amount: CombinePayAmount{
							Total:    1,
							Currency: "CNY",
						},
						OutTradeNo:  "forxxxxxxxxx1",
						Description: "for testing",
					},
				},
				TradeType: Native,
			},
			&CombinePayResponse{
				CodeUrl: "weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00",
			},
			nil,
			true,
		},
		{
			&CombinePayRequest{
				AppId:      client.config.AppId,
				MchId:      client.config.MchId,
				OutTradeNo: "forxxxxxxxxx",
				TimeStart:  time.Now(),
				TimeExpire: time.Now().Add(10 * time.Minute),
				NotifyUrl:  "https://luoji.live/notify",
				Orders: []SubOrder{
					{
						MchId:  mockMchId,
						Attach: "cipher code",
						Amount: CombinePayAmount{
							Total:    1,
							Currency: "CNY",
						},
						OutTradeNo:  "forxxxxxxxxx1",
						Description: "for testing",
					},
				},
				Payer:     &Payer{"openid"},
				TradeType: JSAPI,
			},
			&CombinePayResponse{
				CodeUrl: "weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00",
			},
			nil,
			true,
		},
		{
			&CombinePayRequest{
				AppId:      client.config.AppId,
				MchId:      client.config.MchId,
				OutTradeNo: "forxxxxxxxxx",
				TimeStart:  time.Now(),
				TimeExpire: time.Now().Add(10 * time.Minute),
				NotifyUrl:  "https://luoji.live/notify",
				Orders: []SubOrder{
					{
						MchId:  mockMchId,
						Attach: "cipher code",
						Amount: CombinePayAmount{
							Total:    1,
							Currency: "CNY",
						},
						OutTradeNo:  "forxxxxxxxxx1",
						Description: "for testing",
					},
				},
			},
			&CombinePayResponse{
				CodeUrl: "weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00",
			},
			nil,
			true,
		},
		{
			&CombinePayRequest{
				AppId:      client.config.AppId,
				MchId:      client.config.MchId,
				OutTradeNo: "forxxxxxxxxx",
				TimeStart:  time.Now(),
				TimeExpire: time.Now().Add(10 * time.Minute),
				NotifyUrl:  "https://luoji.live/notify",
				Orders: []SubOrder{
					{
						MchId:  mockMchId,
						Attach: "cipher code",
						Amount: CombinePayAmount{
							Total:    1,
							Currency: "CNY",
						},
						OutTradeNo:  "forxxxxxxxxx1",
						Description: "for testing",
					},
				},
			},
			&CombinePayResponse{
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
		{
			&CombinePayRequest{
				AppId:      client.config.AppId,
				MchId:      client.config.MchId,
				OutTradeNo: "forxxxxxxxxx",
				TimeStart:  time.Now(),
				TimeExpire: time.Now().Add(10 * time.Minute),
				NotifyUrl:  "https://luoji.live/notify",
				Orders: []SubOrder{
					{
						MchId:  mockMchId,
						Attach: "cipher code",
						Amount: CombinePayAmount{
							Total:    1,
							Currency: "CNY",
						},
						OutTradeNo:  "forxxxxxxxxx1",
						Description: "for testing",
					},
				},
				Payer:     &Payer{},
				TradeType: JSAPI,
			},
			&CombinePayResponse{
				CodeUrl: "weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00",
			},
			nil,
			false,
		},
		{
			&CombinePayRequest{
				AppId:      client.config.AppId,
				MchId:      client.config.MchId,
				OutTradeNo: "forxxxxxxxxx",
				TimeStart:  time.Now(),
				TimeExpire: time.Now().Add(10 * time.Minute),
				NotifyUrl:  "https://luoji.live/notify",
				Orders: []SubOrder{
					{
						MchId:  mockMchId,
						Attach: "cipher code",
						Amount: CombinePayAmount{
							Total:    1,
							Currency: "CNY",
						},
						OutTradeNo:  "forxxxxxxxxx1",
						Description: "for testing",
					},
				},
				TradeType: JSAPI,
			},
			&CombinePayResponse{
				CodeUrl: "weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00",
			},
			nil,
			false,
		},
		{
			&CombinePayRequest{
				AppId:      client.config.AppId,
				MchId:      client.config.MchId,
				OutTradeNo: "forxxxxxxxxx",
				TimeStart:  time.Now(),
				TimeExpire: time.Now().Add(10 * time.Minute),
				NotifyUrl:  "https://luoji.live/notify",
				Orders:     []SubOrder{},
			},
			&CombinePayResponse{
				CodeUrl: "weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00",
			},
			nil,
			false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		if c.transport != nil {
			client.config.opts.transport = c.transport
			client.secrets.clear()
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
