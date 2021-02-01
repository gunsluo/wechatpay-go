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
	"reflect"
	"testing"
	"time"
)

func TestDoForQuery(t *testing.T) {
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
		req       *QueryRequest
		resp      *QueryResponse
		transport *mockTransport
		pass      bool
	}{
		{
			&QueryRequest{
				OutTradeNo: "S20210119074247105778399200",
			},
			&QueryResponse{
				AppId:          "wxd678efh567hg6787",
				MchId:          "1230000109",
				OutTradeNo:     "S20210119074247105778399200",
				TransactionId:  "4200000914202101195554393855",
				TradeType:      Native,
				TradeState:     "SUCCESS",
				TradeStateDesc: "支付成功",
				BankType:       "OTHERS",
				Attach:         "",
				SuccessTime:    tm,
				Payer:          Payer{OpenId: "ofyak5qYxYJVnhTlrkk_ACWIVrHI"},
				Amount: TransactionAmount{
					Total:         1,
					PayerTotal:    1,
					Currency:      "CNY",
					PayerCurrency: "CNY",
				},
			},
			nil,
			true,
		},
		{
			&QueryRequest{
				MchId:         client.Config().MchId,
				TransactionId: "4200000914202101195554393855",
			},
			&QueryResponse{
				AppId:          "wxd678efh567hg6787",
				MchId:          "1230000109",
				OutTradeNo:     "S20210119074247105778399200",
				TransactionId:  "4200000914202101195554393855",
				TradeType:      Native,
				TradeState:     "SUCCESS",
				TradeStateDesc: "支付成功",
				BankType:       "OTHERS",
				Attach:         "",
				SuccessTime:    tm,
				Payer:          Payer{OpenId: "ofyak5qYxYJVnhTlrkk_ACWIVrHI"},
				Amount: TransactionAmount{
					Total:         1,
					PayerTotal:    1,
					Currency:      "CNY",
					PayerCurrency: "CNY",
				},
			},
			nil,
			true,
		},
		{
			&QueryRequest{
				OutTradeNo: "S20210119NOTFOUND",
			},
			&QueryResponse{},
			nil,
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
