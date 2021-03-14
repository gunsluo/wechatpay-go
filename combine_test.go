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

func TestCombineCloseRequestDo(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req       *CombineCloseRequest
		transport *mockTransport
		pass      bool
	}{
		{
			&CombineCloseRequest{
				AppId:      client.config.AppId,
				OutTradeNo: "fortest",
				Orders: []CloseSubOrder{
					{
						MchId:      client.config.MchId,
						OutTradeNo: "fortest1",
					},
				},
			},
			nil,
			true,
		},
		{
			&CombineCloseRequest{
				OutTradeNo: "fortest",
				Orders: []CloseSubOrder{
					{
						MchId:      client.config.MchId,
						OutTradeNo: "fortest1",
					},
				},
			},
			nil,
			true,
		},

		{
			&CombineCloseRequest{
				OutTradeNo: "fortest",
				Orders:     []CloseSubOrder{},
			},
			nil,
			false,
		},
		{
			&CombineCloseRequest{
				AppId:      client.config.AppId,
				OutTradeNo: "fortest",
				Orders: []CloseSubOrder{
					{
						MchId:      client.config.MchId,
						OutTradeNo: "fortest1",
					},
				},
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
			client.secrets.clear()
		}

		err := c.req.Do(ctx, client)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, err)
		}

		if err != nil {
			continue
		}
	}
}

func TestDoForCombineQuery(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	tm, err := time.Parse(time.RFC3339, "2021-01-19T15:43:01+08:00")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		req       *CombineQueryRequest
		resp      *CombineQueryResponse
		transport *mockTransport
		pass      bool
	}{
		{
			&CombineQueryRequest{
				OutTradeNo: "S20210119074247105778399200",
			},
			&CombineQueryResponse{
				AppId:      "wxd678efh567hg6787",
				MchId:      "1230000109",
				OutTradeNo: "S20210119074247105778399200",
				Orders: []QuerySubOrder{
					{

						MchId:         "1230000109",
						OutTradeNo:    "S20210119074247105778399201",
						TransactionId: "4200000914202101195554393855",
						TradeType:     Native,
						TradeState:    "SUCCESS",
						BankType:      "OTHERS",
						Attach:        "",
						SuccessTime:   tm,
						Amount: CombineSubOrderAmount{
							Total:         1,
							PayerTotal:    1,
							Currency:      "CNY",
							PayerCurrency: "CNY",
						},
					},
				},
				Payer: &Payer{OpenId: "ofyak5qYxYJVnhTlrkk_ACWIVrHI"},
			},
			nil,
			true,
		},
		{
			&CombineQueryRequest{},
			&CombineQueryResponse{},
			nil,
			false,
		},
		{
			&CombineQueryRequest{
				OutTradeNo: "S20210119NOTFOUND",
			},
			&CombineQueryResponse{},
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
