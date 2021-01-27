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
	"reflect"
	"testing"
	"time"
)

func TestDoForPay(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fail()
	}

	cases := []struct {
		req  *PayRequest
		resp *PayRespone
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
			&PayRespone{
				CodeUrl: "weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00",
			},
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		resp, err := c.req.Do(ctx, client)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c.resp, resp) {
			t.Fatalf("expect %v, got %v", c.resp, resp)
		}
	}
}
