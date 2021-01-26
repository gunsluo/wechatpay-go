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
	"testing"

	"github.com/gunsluo/wechatpay-go/v3/sign"
)

func mockNewClient() (*client, error) {
	var (
		appId          = "wxd678efh567hg6787"
		mchId          = "1230000109"
		apiv3Secret    = "AES256Key-32Characters1234567890"
		serialNo       = "5157F09EFDC096DE15EBE81A47057A7232F1B8E1"
		privateKeyPath = "./test_fixtures/mock_private_key_pkcs8.pem"
	)
	return newClient(
		Config{
			AppId:       appId,
			MchId:       mchId,
			Apiv3Secret: apiv3Secret,
			Cert: CertSuite{
				SerialNo:       serialNo,
				PrivateKeyPath: privateKeyPath,
			},
		})
}

func TestNewClient(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fail()
	}
}

func TestSignatureForClient(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fail()
	}

	var ts int64 = 1611368330
	cases := []struct {
		req    *sign.RequestSignature
		expect string
	}{
		{
			&sign.RequestSignature{
				Method:    "POST",
				Url:       "https://api.mch.weixin.qq.com/v3/pay/transactions/native",
				Timestamp: ts,
				Nonce:     "AF1404CC2980FB414C99C0B98883BD42",
				Body:      []byte(`{"appid":"wx81be3101902f7cb2","mchid":"1601959334","description":"for testing","out_trade_no":"S20210124144305172434","time_expire":"2021-01-24T14:53:05+08:00","attach":"cipher code","notify_url":"https://luoji.live/notify","amount":{"total":1,"currency":"CNY"},"detail":{},"scene_info":{"payer_client_ip":"","store_info":{"id":""}}}`),
			},
			`WECHATPAY2-SHA256-RSA2048 mchid="1230000109",nonce_str="AF1404CC2980FB414C99C0B98883BD42",signature="ItuRCG6nAf6ZUi5C5LPa0beCGrG7+G4NdaCHLTmym+UzuZHFgFeqRZ4zKQ0n93qehchFWfQ7s00pgABYvXcOMsV1ld7AUjDTZBPucJK6yhFKz9jd20wtRdDG4LRCZcaTowD2f7LtlixFm8F3/YQaBavxiOe54tc3RX/22flYRzy4YFOpBt+bmjSPZIdSFi53323u7cohwvdHwX+avQCtLZKAUNFJIob66u05BbDEITzYuHjakjpb5btvWemjoZBPxkiETzmd4Oa1y2U+rfFCPZyWT4EV7UxHeEizBL8DkubEBD3KXeArqRX6yoMAU4ywmdFeWDbv1EF0Ndy9hiddZQ==",timestamp="1611368330",serial_no="5157F09EFDC096DE15EBE81A47057A7232F1B8E1"`,
		},
	}

	for _, c := range cases {
		signature, err := client.Signature(c.req)
		if err != nil {
			t.Fatal(err)
		}

		if signature != c.expect {
			t.Fatalf("expect %s, got %s", c.expect, signature)
		}
	}
}

func TestDoForClient(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fail()
	}

	cases := []struct {
		req    interface{}
		method string
		url    string
		//expect string
	}{
		{
			&CertificatesRequest{},
			http.MethodGet,
			"https://api.mch.weixin.qq.com/v3/certificates",
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		result := client.Do(ctx, c.method, c.url, c.req)
		if result.Err != nil {
			t.Fatal(result.Err)
		}
	}
}

func TestVerifySignatureForClient(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fail()
	}
}
