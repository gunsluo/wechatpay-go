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
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gunsluo/wechatpay-go/v3/sign"
)

func TestNewClient(t *testing.T) {
	cases := []struct {
		appId          string
		mchId          string
		apiv3Secret    string
		serialNo       string
		privateKeyTxt  string
		privateKeyPath string
		expect         bool
	}{
		{
			"wxd678efh567hg6787",
			"1230000109",
			"AES256Key-32Characters1234567890",
			mockSerialNo,
			``,
			"./test_fixtures/mock_private_key_pkcs8.pem",
			true,
		},
		{
			"wxd678efh567hg6787",
			"1230000109",
			"AES256Key-32Characters1234567890",
			mockSerialNo,
			`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCprsmcXPHqLtnP
oPDGUoMULK2WOo5FW8c72Svnqn/4aXPaJhlOtPxtX2frqIhTjwcOs6hNm3XFTGBL
MrdB94YQvj+Q7P12GNmxXG+9Ms+uUyJToYjlYDAG6UFKE10Jkm9cDGuLSkekU1Ao
rKE1G1wndH37w4AzVXoGBQ3NIiyW8jIm8Zi3/WNCVpHUoXYUuyhFEZ23fXytnps4
hARgg6NvPncIKtWvlUh85ZVOSsqc1T8dFaeDRXaj7r3jdJJ74tsGRMvZyUipJXyE
3uR2QkrGyia+0phDpC6zeMMpP+MQO9ohh+xQWBCeyvQjjnOPAlGThl+ThfXImU30
HL17oHdBAgMBAAECggEADm6FSz1Efgx6DgS8NcHy0BZ0tSBJ1XBW46o2579Cnxgo
+FbhNCaEibDhn9N3tNOnYAK7v84HGD7EueCYYY3x4x6rPWJKtG6spT8dadQWgdck
RkSo5glmTFAuc2RuN1AzFHsh8njg2wMTAEKee2vWTKzFwlIAZ11PwY9Qey/65uOT
Bi8q1Rssu6xofNadO5MbqMJ1Tl8DDIaLGnzTzbHrk9thBUo1FwFjJWTVI7nz2En4
Yc/G1/LQJfiQ31F+lkL3j6ABRJqtsgb07r9H/hT6+fd1hGDt2qKuS+E1mLDp9fHw
n6UyS4HyB7DA/XtFZ9z0VtAlmcGoUkyJLtXjEmwsGQKBgQDeCtE3spULpC7VPqk1
xv034C6zybZ7y8kSKwRvyYwkzdgSRgVaKsTVb8RNYor8hoGrVgdXFqQUI8O/v1cN
9wFoGYJT0LHre/YzOg31TkQkBfHHCFH/L50uOJcIQueftctz5Bwj6bJO/ih5iIAK
yjrHse4PdIiEJfz2D9hc4wnxrwKBgQDDogrWlUCTj2fvmZfkWR3Hbs0kIHd7zjIk
bJJONGtD8gE4i562tajC1mKoQEwt4YSwWsBkGAw1LhvMROQFT6AOaIIhHNex1Z3t
c2gAdEeWOMmzZnnhwWzTiYJomixrFkmEwT3EJK89GO3E0FH5S+G1P1tNXq38Vpty
1YVqOgMSDwKBgFrzuWGEQDMljJ2C7lL98KlbpiW1AY/SGMndXxLfTw2gV9qcXgLi
NABtqM4+CEqKWkExmw4cUxeA0uUPXnx06lmW4WCtwsN/4oh3RlJuPdE3siLiEJxk
B5FwUsVqinBMSktta+12A7kBuNiXhkNlNRCpnKcuB+GBog20zd62jVM3AoGBALcA
zFazQ7dFfRq7eUUYwCyhT7Et1dewqWM9VRdnHbhvmAjHQu7zvCyW069Ehn6c6bz3
B+YaQME2orZQ82SsebNAvAoxquwmQhevz2gtXhH+iWASyo0Onbi8d4tWPZrnPFq9
UgQ7tNnYigOEREqKW1drLwOPP/4/Hicr6iPWpKytAoGAEQ6J/RB/olEAC46ACoFo
FBgA+GUbDB0xBcA2inEt3q//208YMkjnKM871n89HpAgms5xrK32T69lduebk7Ar
9wWvkJVUwI9VDXomCFQqtiGzHlTl1Xq31BfeIDyq1ayQmTkRpRqIagbDZVtM+ha/
0I2SEzTObt07wcYcYG2Chvg=
-----END PRIVATE KEY-----`,
			"",
			true,
		},
		{
			"wxd678efh567hg6787",
			"1230000109",
			"AES256Key-32Characters1234567890",
			mockSerialNo,
			``,
			"./test_fixtures/mock_private_key.pem",
			false,
		},
		{
			"wxd678efh567hg6787",
			"1230000109",
			"AES256Key-32Characters1234567890",
			mockSerialNo,
			`-----BEGIN PRIVATE KEY----------END PRIVATE KEY-----`,
			"",
			false,
		},
		{
			"",
			"1230000109",
			"AES256Key-32Characters1234567890",
			mockSerialNo,
			"./test_fixtures/mock_private_key_pkcs8.pem",
			"",
			false,
		},
		{
			"wxd678efh567hg6787",
			"",
			"AES256Key-32Characters1234567890",
			mockSerialNo,
			"",
			"./test_fixtures/mock_private_key_pkcs8.pem",
			false,
		},

		{
			"wxd678efh567hg6787",
			"1230000109",
			"",
			mockSerialNo,
			"./test_fixtures/mock_private_key_pkcs8.pem",
			"",
			false,
		},
		{
			"wxd678efh567hg6787",
			"1230000109",
			"AES256Key-32Characters1234567890",
			"",
			"./test_fixtures/mock_private_key_pkcs8.pem",
			"",
			false,
		},
		{
			"wxd678efh567hg6787",
			"1230000109",
			"AES256Key-32Characters1234567890",
			mockSerialNo,
			"",
			"",
			false,
		},
	}

	for _, c := range cases {
		_, err := NewClient(
			Config{
				AppId:       c.appId,
				MchId:       c.mchId,
				Apiv3Secret: c.apiv3Secret,
				Cert: CertSuite{
					SerialNo:       c.serialNo,
					PrivateKeyPath: c.privateKeyPath,
					PrivateKeyTxt:  c.privateKeyTxt,
				},
			})
		expect := err == nil
		if expect != c.expect {
			t.Fatalf("expect %v, got %v, err: %v", c.expect, expect, err)
		}
	}
}

func TestSignatureForClient(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req    *sign.RequestSignature
		pass   bool
		expect string
	}{
		{
			&sign.RequestSignature{
				Method:    "POST",
				Url:       "https://api.mch.weixin.qq.com/v3/pay/transactions/native",
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Body:      []byte(`{"appid":"wx81be3101902f7cb2","mchid":"1601959334","description":"for testing","out_trade_no":"S20210124144305172434","time_expire":"2021-01-24T14:53:05+08:00","attach":"cipher code","notify_url":"https://luoji.live/notify","amount":{"total":1,"currency":"CNY"},"detail":{},"scene_info":{"payer_client_ip":"","store_info":{"id":""}}}`),
			},
			true,
			`WECHATPAY2-SHA256-RSA2048 mchid="1230000109",nonce_str="AF1404CC2980FB414C99C0B98883BD42",signature="ItuRCG6nAf6ZUi5C5LPa0beCGrG7+G4NdaCHLTmym+UzuZHFgFeqRZ4zKQ0n93qehchFWfQ7s00pgABYvXcOMsV1ld7AUjDTZBPucJK6yhFKz9jd20wtRdDG4LRCZcaTowD2f7LtlixFm8F3/YQaBavxiOe54tc3RX/22flYRzy4YFOpBt+bmjSPZIdSFi53323u7cohwvdHwX+avQCtLZKAUNFJIob66u05BbDEITzYuHjakjpb5btvWemjoZBPxkiETzmd4Oa1y2U+rfFCPZyWT4EV7UxHeEizBL8DkubEBD3KXeArqRX6yoMAU4ywmdFeWDbv1EF0Ndy9hiddZQ==",timestamp="1611368330",serial_no="477ED0046A54F0360A72A63A8F2816312AAEAB53"`,
		},
		{
			&sign.RequestSignature{
				Method:    "POST",
				Url:       "https:\n//api.mch.weixin.qq.com/v3/pay/transactions/native",
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Body:      []byte(`{"appid":"wx81be3101902f7cb2","mchid":"1601959334","description":"for testing","out_trade_no":"S20210124144305172434","time_expire":"2021-01-24T14:53:05+08:00","attach":"cipher code","notify_url":"https://luoji.live/notify","amount":{"total":1,"currency":"CNY"},"detail":{},"scene_info":{"payer_client_ip":"","store_info":{"id":""}}}`),
			},
			false,
			``,
		},
	}

	for _, c := range cases {
		signature, err := client.Signature(c.req)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err:%v", c.pass, pass, err)
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
		t.Fatal("client is nil")
	}

	cases := []struct {
		req    interface{}
		method string
		url    string
		pass   bool
	}{
		{
			&PayRequest{},
			http.MethodPost,
			"https://api.mch.weixin.qq.com/v3/pay/transactions/id/4200000914202101195554393855",
			true,
		},
		{
			nil,
			http.MethodGet,
			"https://api.mch.weixin.qq.com/v3/certificates",
			true,
		},
		{
			(*PayRequest)(nil),
			http.MethodGet,
			"https://api.mch.weixin.qq.com/v3/certificates",
			true,
		},
		{
			&CertificatesRequest{},
			http.MethodGet,
			"https://api.mch.weixin.qq.com/v3/certificates",
			true,
		},
		{
			&CertificatesRequest{},
			http.MethodGet,
			"https:\n//api.mch.weixin.qq.com/v3/certificates",
			false,
		},
		{
			&CertificatesRequest{},
			http.MethodGet,
			"https://api.mch.weixin.qq.com/v3/nocert",
			false,
		},
		{
			&CertificatesRequest{},
			http.MethodGet,
			"https://api.mch.weixin.qq.com/v3/invalidresp",
			false,
		},
		{
			&CertificatesRequest{},
			http.MethodGet,
			"https://api.mch.weixin.qq.com/v3/invalidrespdata",
			false,
		},
		{
			&CertificatesRequest{},
			http.MethodGet,
			"https://api.mch.weixin.qq.com/v3/invalidheader",
			false,
		},
		{
			&CertificatesRequest{},
			http.MethodGet,
			"https://api.mch.weixin.qq.com/v3/nodataresp",
			false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		result := client.Do(ctx, c.method, c.url, c.req)
		pass := result.Err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, result.Err)
		}
	}
}

func TestFailedDoForClient(t *testing.T) {
	cases := []struct {
		req       interface{}
		method    string
		url       string
		newClient func() (*client, error)
	}{
		{
			&CertificatesRequest{},
			http.MethodGet,
			"https://api.mch.weixin.qq.com/v3/validsign",
			func() (*client, error) {
				client, err := mockNewClient()
				if err != nil {
					return nil, err
				}

				client.privateKey = &rsa.PrivateKey{
					PublicKey: rsa.PublicKey{
						N: fromBase10("935393046677"),
						E: 65537,
					},
					D: fromBase10("72663984313281163"),
					Primes: []*big.Int{
						fromBase10("9892036654808464"),
						fromBase10("9456020830884701"),
					},
				}

				return client, nil
			},
		},
		{
			&CertificatesRequest{},
			http.MethodGet,
			"https://api.mch.weixin.qq.com/v3/certificates",
			func() (*client, error) {
				client, err := mockNewClient(&mockTransport{
					RoundTripFn: func(req *http.Request) (*http.Response, error) {
						var resp = &http.Response{
							StatusCode: http.StatusOK,
						}

						resp.Header = http.Header{}
						resp.Body = ioutil.NopCloser(strings.NewReader("{}"))
						return resp, nil
					},
				})
				if err != nil {
					return nil, err
				}

				return client, nil
			},
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		client, err := c.newClient()
		if err != nil {
			t.Fatal(err)
		}

		result := client.Do(ctx, c.method, c.url, c.req)
		if result.Err == nil {
			t.Fatal("should be an error")
		}
	}
}

func TestDoExtraWorkflow(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req    *sign.RequestSignature
		result *Result
		pass   bool
	}{
		{
			&sign.RequestSignature{
				Method:    http.MethodGet,
				Url:       "https://api.mch.weixin.qq.com/v3/certificates",
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Body:      []byte(""),
			},
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "",
				SerialNo:  mockSerialNo,
				Body: []byte(`{"data":[{"effective_time":"2020-09-17T14:26:23+08:00","encrypt_certificate":{"algorithm":"AEAD_AES_256_GCM","associated_data":"certificate","ciphertext":"/M2eAJyVx/0y8JOErsNEWbYpikwKMS0hDahBYrR9Tnqvaxw/WLMHyLq7G3GUoWx3NSwYZlSZ+1JxAMTd4yge1B8bxY7OLrDkXm+BBDVypy5jCi/gcTQduTJpR4nRcBRYtEIxLGLrVaUXlDjDa4nM0mUPk6XA7AAUUAl3z5lYISapsFYUuHO9splBrmUESHxzRhSfsTyW68ll8o+ND7xA5R94slxzZIVdVg2Tz/3uXi5X1Qu5oi9Dn7pFdHD7++msMB7rgSJUTIFMwZ2GhAX3f/vVWemSMCymPPxzYxdiGFJJ8oBaIn+17pwulmz6NodFS0ilJr9wBs/05gqxe5L6S64ApwXNTfq3YJFVIU6munBaHomRZqsMg3MQlji9yNLBdKO2hk2rq/jCaBLsqcrCHEMEEULA5/1ImeYEkKcX2vIiVtKX8WxxP4M/Gq7btAQZVGzvczopb3wZNu1QLnzC13ov0pB5BPMhrx0tE4rLuZ5d+uzGOwuI8CvqOa+8TQ0DNGNaEA/IPrMJCVvmLrDi/aMQB+P4mO9BhUlfGHwQL7Q0anHzZaGHGkYyEGoTPmqQcY1mRbVcXDpIGn7rfHgiXnQTurB886T//ddhcv1/LQmcohSveZJAltcaDlmeqMgc+bXsOlAy6JNIIVPJ04ysI+V7nc0O4k4A32ZYA1hK52CU1YWz3vMoaaHVr/t6AF3dVWE1CphhNIwGbaz9M1sgEsWwT8LKLG5csgVwG20LO8wmLkxNUQ4fSkMdC+2Qv+rSFd8rlT1j+sYEbPVq6E6URkYPUKMqI1mEEudU1Rx0bE/pjj7+++0gX1H7sHp4+02KLdWS27gptHVXdDjNFPyCEshfVL2B8aEhq8PxSDG5zTqWHrKBAl04WU3kjlSsKZPrpKyhpIrKbEZHcrip3wOGeMf+4XDoZ8Iq8KoM8R6m8wkWi0GAW4G743O44PxHFvljKDIkIQm8gWV37jC3+qb/ZwUDxHONw3tHMH8XWsCVq1KAtKeE/iE9CCmE+ht7K4B+w0DeqKEicm0dkdjuFc9IgFa1W+q0HqGFI2Snd6ZX6crUy1I1vkRTQRj1mqjaP7dFOFV0JMpK/4CKMruZfUilNfOnSoKqHA2jPQ3f4ro0H22bF/PNhOWXp6Tzl5ZVbIFBIMdD9+ocq1lDH7vcBfKVwUltKl7jgI9HlpCDPZp++Mt3C4lPDzP/XrqorJnFBKw8eMBHS7N+jDhzhqJnI3ldwlGxUsqS/hj+jUUPpYINe/UtVwlOBi/tfuEfv47H5YgbP+Y3dz78a6KJUcA7caPSSqX+8LBcwEEZELXR8gU/AxwoDAsHM1pb7wc9fslct+awivfRi47AJtFeeZMGF6bb14VnbzvIZdpZRBIzHlvUqP+t8ZKEUvEJ+lVk7vv0/ySWBZbt0oA5XQ2RVwgzKGOgfMzZafsWAqrq1PGYjJqBbm/hudPtqsBridW/QjoE2Bp+Qnp8mWhdlSP8dgdeefLEeZGUSJx0Tzu2hBveEz7jMNQSOyg8HEE=","nonce":"eabb3e044577"},"expire_time":"2025-09-16T14:26:23+08:00","serial_no":"477ED0046A54F0360A72A63A8F2816312AAEAB53"}]}
`),
			},
			true,
		},
		{
			&sign.RequestSignature{
				Method:    http.MethodGet,
				Url:       "https://api.mch.weixin.qq.com/v3/certificates",
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Body:      []byte(""),
			},
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{`),
			},
			false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		err := client.doExtraWorkflow(ctx, c.req, c.result)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, err)
		}
	}
}

func TestUpgradeCertWorkflow(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req    *sign.RequestSignature
		result *Result
		pass   bool
	}{
		{
			&sign.RequestSignature{
				Method:    http.MethodGet,
				Url:       "https://api.mch.weixin.qq.com/v3/certificates",
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Body:      []byte(""),
			},
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "",
				SerialNo:  mockSerialNo,
				Body: []byte(`{"data":[{"effective_time":"2020-09-17T14:26:23+08:00","encrypt_certificate":{"algorithm":"AEAD_AES_256_GCM","associated_data":"certificate","ciphertext":"/M2eAJyVx/0y8JOErsNEWbYpikwKMS0hDahBYrR9Tnqvaxw/WLMHyLq7G3GUoWx3NSwYZlSZ+1JxAMTd4yge1B8bxY7OLrDkXm+BBDVypy5jCi/gcTQduTJpR4nRcBRYtEIxLGLrVaUXlDjDa4nM0mUPk6XA7AAUUAl3z5lYISapsFYUuHO9splBrmUESHxzRhSfsTyW68ll8o+ND7xA5R94slxzZIVdVg2Tz/3uXi5X1Qu5oi9Dn7pFdHD7++msMB7rgSJUTIFMwZ2GhAX3f/vVWemSMCymPPxzYxdiGFJJ8oBaIn+17pwulmz6NodFS0ilJr9wBs/05gqxe5L6S64ApwXNTfq3YJFVIU6munBaHomRZqsMg3MQlji9yNLBdKO2hk2rq/jCaBLsqcrCHEMEEULA5/1ImeYEkKcX2vIiVtKX8WxxP4M/Gq7btAQZVGzvczopb3wZNu1QLnzC13ov0pB5BPMhrx0tE4rLuZ5d+uzGOwuI8CvqOa+8TQ0DNGNaEA/IPrMJCVvmLrDi/aMQB+P4mO9BhUlfGHwQL7Q0anHzZaGHGkYyEGoTPmqQcY1mRbVcXDpIGn7rfHgiXnQTurB886T//ddhcv1/LQmcohSveZJAltcaDlmeqMgc+bXsOlAy6JNIIVPJ04ysI+V7nc0O4k4A32ZYA1hK52CU1YWz3vMoaaHVr/t6AF3dVWE1CphhNIwGbaz9M1sgEsWwT8LKLG5csgVwG20LO8wmLkxNUQ4fSkMdC+2Qv+rSFd8rlT1j+sYEbPVq6E6URkYPUKMqI1mEEudU1Rx0bE/pjj7+++0gX1H7sHp4+02KLdWS27gptHVXdDjNFPyCEshfVL2B8aEhq8PxSDG5zTqWHrKBAl04WU3kjlSsKZPrpKyhpIrKbEZHcrip3wOGeMf+4XDoZ8Iq8KoM8R6m8wkWi0GAW4G743O44PxHFvljKDIkIQm8gWV37jC3+qb/ZwUDxHONw3tHMH8XWsCVq1KAtKeE/iE9CCmE+ht7K4B+w0DeqKEicm0dkdjuFc9IgFa1W+q0HqGFI2Snd6ZX6crUy1I1vkRTQRj1mqjaP7dFOFV0JMpK/4CKMruZfUilNfOnSoKqHA2jPQ3f4ro0H22bF/PNhOWXp6Tzl5ZVbIFBIMdD9+ocq1lDH7vcBfKVwUltKl7jgI9HlpCDPZp++Mt3C4lPDzP/XrqorJnFBKw8eMBHS7N+jDhzhqJnI3ldwlGxUsqS/hj+jUUPpYINe/UtVwlOBi/tfuEfv47H5YgbP+Y3dz78a6KJUcA7caPSSqX+8LBcwEEZELXR8gU/AxwoDAsHM1pb7wc9fslct+awivfRi47AJtFeeZMGF6bb14VnbzvIZdpZRBIzHlvUqP+t8ZKEUvEJ+lVk7vv0/ySWBZbt0oA5XQ2RVwgzKGOgfMzZafsWAqrq1PGYjJqBbm/hudPtqsBridW/QjoE2Bp+Qnp8mWhdlSP8dgdeefLEeZGUSJx0Tzu2hBveEz7jMNQSOyg8HEE=","nonce":"eabb3e044577"},"expire_time":"2025-09-16T14:26:23+08:00","serial_no":"477ED0046A54F0360A72A63A8F2816312AAEAB53"}]}
`),
			},
			true,
		},
		{
			&sign.RequestSignature{
				Method:    http.MethodGet,
				Url:       "https://api.mch.weixin.qq.com/v3/certificates",
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Body:      []byte(""),
			},
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{`),
			},
			false,
		},
		{
			&sign.RequestSignature{
				Method:    http.MethodGet,
				Url:       "https://api.mch.weixin.qq.com/v3/certificates",
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Body:      []byte(""),
			},
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "",
				SerialNo:  mockSerialNo,
				Body: []byte(`{"data":[{"effective_time":"2020-09-17T14:26:23+08:00","encrypt_certificate":{"algorithm":"AEAD_AES_256_GCM","associated_data":"c","ciphertext":"/M2eAJyVx/0y8JOErsNEWbYpikwKMS0hDahBYrR9Tnqvaxw/WLMHyLq7G3GUoWx3NSwYZlSZ+1JxAMTd4yge1B8bxY7OLrDkXm+BBDVypy5jCi/gcTQduTJpR4nRcBRYtEIxLGLrVaUXlDjDa4nM0mUPk6XA7AAUUAl3z5lYISapsFYUuHO9splBrmUESHxzRhSfsTyW68ll8o+ND7xA5R94slxzZIVdVg2Tz/3uXi5X1Qu5oi9Dn7pFdHD7++msMB7rgSJUTIFMwZ2GhAX3f/vVWemSMCymPPxzYxdiGFJJ8oBaIn+17pwulmz6NodFS0ilJr9wBs/05gqxe5L6S64ApwXNTfq3YJFVIU6munBaHomRZqsMg3MQlji9yNLBdKO2hk2rq/jCaBLsqcrCHEMEEULA5/1ImeYEkKcX2vIiVtKX8WxxP4M/Gq7btAQZVGzvczopb3wZNu1QLnzC13ov0pB5BPMhrx0tE4rLuZ5d+uzGOwuI8CvqOa+8TQ0DNGNaEA/IPrMJCVvmLrDi/aMQB+P4mO9BhUlfGHwQL7Q0anHzZaGHGkYyEGoTPmqQcY1mRbVcXDpIGn7rfHgiXnQTurB886T//ddhcv1/LQmcohSveZJAltcaDlmeqMgc+bXsOlAy6JNIIVPJ04ysI+V7nc0O4k4A32ZYA1hK52CU1YWz3vMoaaHVr/t6AF3dVWE1CphhNIwGbaz9M1sgEsWwT8LKLG5csgVwG20LO8wmLkxNUQ4fSkMdC+2Qv+rSFd8rlT1j+sYEbPVq6E6URkYPUKMqI1mEEudU1Rx0bE/pjj7+++0gX1H7sHp4+02KLdWS27gptHVXdDjNFPyCEshfVL2B8aEhq8PxSDG5zTqWHrKBAl04WU3kjlSsKZPrpKyhpIrKbEZHcrip3wOGeMf+4XDoZ8Iq8KoM8R6m8wkWi0GAW4G743O44PxHFvljKDIkIQm8gWV37jC3+qb/ZwUDxHONw3tHMH8XWsCVq1KAtKeE/iE9CCmE+ht7K4B+w0DeqKEicm0dkdjuFc9IgFa1W+q0HqGFI2Snd6ZX6crUy1I1vkRTQRj1mqjaP7dFOFV0JMpK/4CKMruZfUilNfOnSoKqHA2jPQ3f4ro0H22bF/PNhOWXp6Tzl5ZVbIFBIMdD9+ocq1lDH7vcBfKVwUltKl7jgI9HlpCDPZp++Mt3C4lPDzP/XrqorJnFBKw8eMBHS7N+jDhzhqJnI3ldwlGxUsqS/hj+jUUPpYINe/UtVwlOBi/tfuEfv47H5YgbP+Y3dz78a6KJUcA7caPSSqX+8LBcwEEZELXR8gU/AxwoDAsHM1pb7wc9fslct+awivfRi47AJtFeeZMGF6bb14VnbzvIZdpZRBIzHlvUqP+t8ZKEUvEJ+lVk7vv0/ySWBZbt0oA5XQ2RVwgzKGOgfMzZafsWAqrq1PGYjJqBbm/hudPtqsBridW/QjoE2Bp+Qnp8mWhdlSP8dgdeefLEeZGUSJx0Tzu2hBveEz7jMNQSOyg8HEE=","nonce":"eabb3e044577"},"expire_time":"2025-09-16T14:26:23+08:00","serial_no":"477ED0046A54F0360A72A63A8F2816312AAEAB53"}]}
`),
			},
			false,
		},
		{
			&sign.RequestSignature{
				Method:    http.MethodGet,
				Url:       "https://api.mch.weixin.qq.com/v3/certificates",
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Body:      []byte(""),
			},
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "",
				SerialNo:  mockSerialNo,
				Body: []byte(`{"data":[{"effective_time":"2020-09-17T14:26:23+08:00","encrypt_certificate":{"algorithm":"AEAD_AES_256_GCM","associated_data":"certificate","ciphertext":"tJjSQMG758oX39qpn/RoZPZ3qh8LRIIwcnQeFhU/alQ=","nonce":"eabb3e044577"},"expire_time":"2025-09-16T14:26:23+08:00","serial_no":"477ED0046A54F0360A72A63A8F2816312AAEAB53"}]}
`),
			},
			false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		err := upgradeCertWorkflow(ctx, client, c.req, c.result)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, err)
		}
	}
}

func TestVerifySignatureForClient(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		result        *Result
		mocktransport *mockTransport
		pass          bool
	}{
		{
			&Result{
				Body: []byte(`{"data":[{"effective_time":"2020-09-17T14:26:23+08:00","encrypt_certificate":{"algorithm":"AEAD_AES_256_GCM","associated_data":"certificate","ciphertext":"/M2eAJyVx/0y8JOErsNEWbYpikwKMS0hDahBYrR9Tnqvaxw/WLMHyLq7G3GUoWx3NSwYZlSZ+1JxAMTd4yge1B8bxY7OLrDkXm+BBDVypy5jCi/gcTQduTJpR4nRcBRYtEIxLGLrVaUXlDjDa4nM0mUPk6XA7AAUUAl3z5lYISapsFYUuHO9splBrmUESHxzRhSfsTyW68ll8o+ND7xA5R94slxzZIVdVg2Tz/3uXi5X1Qu5oi9Dn7pFdHD7++msMB7rgSJUTIFMwZ2GhAX3f/vVWemSMCymPPxzYxdiGFJJ8oBaIn+17pwulmz6NodFS0ilJr9wBs/05gqxe5L6S64ApwXNTfq3YJFVIU6munBaHomRZqsMg3MQlji9yNLBdKO2hk2rq/jCaBLsqcrCHEMEEULA5/1ImeYEkKcX2vIiVtKX8WxxP4M/Gq7btAQZVGzvczopb3wZNu1QLnzC13ov0pB5BPMhrx0tE4rLuZ5d+uzGOwuI8CvqOa+8TQ0DNGNaEA/IPrMJCVvmLrDi/aMQB+P4mO9BhUlfGHwQL7Q0anHzZaGHGkYyEGoTPmqQcY1mRbVcXDpIGn7rfHgiXnQTurB886T//ddhcv1/LQmcohSveZJAltcaDlmeqMgc+bXsOlAy6JNIIVPJ04ysI+V7nc0O4k4A32ZYA1hK52CU1YWz3vMoaaHVr/t6AF3dVWE1CphhNIwGbaz9M1sgEsWwT8LKLG5csgVwG20LO8wmLkxNUQ4fSkMdC+2Qv+rSFd8rlT1j+sYEbPVq6E6URkYPUKMqI1mEEudU1Rx0bE/pjj7+++0gX1H7sHp4+02KLdWS27gptHVXdDjNFPyCEshfVL2B8aEhq8PxSDG5zTqWHrKBAl04WU3kjlSsKZPrpKyhpIrKbEZHcrip3wOGeMf+4XDoZ8Iq8KoM8R6m8wkWi0GAW4G743O44PxHFvljKDIkIQm8gWV37jC3+qb/ZwUDxHONw3tHMH8XWsCVq1KAtKeE/iE9CCmE+ht7K4B+w0DeqKEicm0dkdjuFc9IgFa1W+q0HqGFI2Snd6ZX6crUy1I1vkRTQRj1mqjaP7dFOFV0JMpK/4CKMruZfUilNfOnSoKqHA2jPQ3f4ro0H22bF/PNhOWXp6Tzl5ZVbIFBIMdD9+ocq1lDH7vcBfKVwUltKl7jgI9HlpCDPZp++Mt3C4lPDzP/XrqorJnFBKw8eMBHS7N+jDhzhqJnI3ldwlGxUsqS/hj+jUUPpYINe/UtVwlOBi/tfuEfv47H5YgbP+Y3dz78a6KJUcA7caPSSqX+8LBcwEEZELXR8gU/AxwoDAsHM1pb7wc9fslct+awivfRi47AJtFeeZMGF6bb14VnbzvIZdpZRBIzHlvUqP+t8ZKEUvEJ+lVk7vv0/ySWBZbt0oA5XQ2RVwgzKGOgfMzZafsWAqrq1PGYjJqBbm/hudPtqsBridW/QjoE2Bp+Qnp8mWhdlSP8dgdeefLEeZGUSJx0Tzu2hBveEz7jMNQSOyg8HEE=","nonce":"eabb3e044577"},"expire_time":"2025-09-16T14:26:23+08:00","serial_no":"477ED0046A54F0360A72A63A8F2816312AAEAB53"}]}
`),
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				SerialNo:  mockSerialNo,
				Signature: "KDrEP098zDlbX6ioHrS7sKLUNIqxzQcf+JXCkG5W44EKno1/qmI4WBf/sh63fwC++ZKBn/4gfEj7Iv4W3YH5kfgki6fFvfrRrGAxROiLSn/FZhbVu9E8pR4McxOR04UP+opyFhDL3lpPKqFB5AnUsTHhoCcZADzuHmCVHwU20DMGa00/Wr3kEcNYByy5hqz5sn7VbjoMs1KAMzmEKxXiIZIu5nvf4b4gk7zNvNWjMAUzsFHELHLfNqNMetzW/TIc0RL4S9vQL+GR7qRnzgKGkd5bfOn611jPEv1ut7UbWV+qvIYKeyaMe9xfyH83fobzSD9sbfZFwmb0wYMqPIgMtw==",
			},
			nil,
			true,
		},
		{
			&Result{},
			&mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					var resp = &http.Response{
						StatusCode: http.StatusInternalServerError,
					}

					resp.Body = ioutil.NopCloser(strings.NewReader(`{"code":"ERROR_NAME","message":"ERROR_DESCRIPTION"}`))
					return resp, nil
				},
			},
			false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		if c.mocktransport != nil {
			client.config.opts.transport = c.mocktransport
			client.secrets.clear()
		}
		err = client.VerifySignature(ctx, c.result)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err %v", c.pass, pass, err)
		}
	}
}

func TestOnceDownloadCertificates(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		mocktransport *mockTransport
		pass          bool
	}{
		{
			&mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					var resp = &http.Response{
						StatusCode: http.StatusOK,
					}
					mockBody := `{"data":[{"effective_time":"2020-09-17T14:26:23+08:00","encrypt_certificate":{"algorithm":"AEAD_AES_256_GCM","associated_data":"certificate","ciphertext":"/M2eAJyVx/0y8JOErsNEWbYpikwKMS0hDahBYrR9Tnqvaxw/WLMHyLq7G3GUoWx3NSwYZlSZ+1JxAMTd4yge1B8bxY7OLrDkXm+BBDVypy5jCi/gcTQduTJpR4nRcBRYtEIxLGLrVaUXlDjDa4nM0mUPk6XA7AAUUAl3z5lYISapsFYUuHO9splBrmUESHxzRhSfsTyW68ll8o+ND7xA5R94slxzZIVdVg2Tz/3uXi5X1Qu5oi9Dn7pFdHD7++msMB7rgSJUTIFMwZ2GhAX3f/vVWemSMCymPPxzYxdiGFJJ8oBaIn+17pwulmz6NodFS0ilJr9wBs/05gqxe5L6S64ApwXNTfq3YJFVIU6munBaHomRZqsMg3MQlji9yNLBdKO2hk2rq/jCaBLsqcrCHEMEEULA5/1ImeYEkKcX2vIiVtKX8WxxP4M/Gq7btAQZVGzvczopb3wZNu1QLnzC13ov0pB5BPMhrx0tE4rLuZ5d+uzGOwuI8CvqOa+8TQ0DNGNaEA/IPrMJCVvmLrDi/aMQB+P4mO9BhUlfGHwQL7Q0anHzZaGHGkYyEGoTPmqQcY1mRbVcXDpIGn7rfHgiXnQTurB886T//ddhcv1/LQmcohSveZJAltcaDlmeqMgc+bXsOlAy6JNIIVPJ04ysI+V7nc0O4k4A32ZYA1hK52CU1YWz3vMoaaHVr/t6AF3dVWE1CphhNIwGbaz9M1sgEsWwT8LKLG5csgVwG20LO8wmLkxNUQ4fSkMdC+2Qv+rSFd8rlT1j+sYEbPVq6E6URkYPUKMqI1mEEudU1Rx0bE/pjj7+++0gX1H7sHp4+02KLdWS27gptHVXdDjNFPyCEshfVL2B8aEhq8PxSDG5zTqWHrKBAl04WU3kjlSsKZPrpKyhpIrKbEZHcrip3wOGeMf+4XDoZ8Iq8KoM8R6m8wkWi0GAW4G743O44PxHFvljKDIkIQm8gWV37jC3+qb/ZwUDxHONw3tHMH8XWsCVq1KAtKeE/iE9CCmE+ht7K4B+w0DeqKEicm0dkdjuFc9IgFa1W+q0HqGFI2Snd6ZX6crUy1I1vkRTQRj1mqjaP7dFOFV0JMpK/4CKMruZfUilNfOnSoKqHA2jPQ3f4ro0H22bF/PNhOWXp6Tzl5ZVbIFBIMdD9+ocq1lDH7vcBfKVwUltKl7jgI9HlpCDPZp++Mt3C4lPDzP/XrqorJnFBKw8eMBHS7N+jDhzhqJnI3ldwlGxUsqS/hj+jUUPpYINe/UtVwlOBi/tfuEfv47H5YgbP+Y3dz78a6KJUcA7caPSSqX+8LBcwEEZELXR8gU/AxwoDAsHM1pb7wc9fslct+awivfRi47AJtFeeZMGF6bb14VnbzvIZdpZRBIzHlvUqP+t8ZKEUvEJ+lVk7vv0/ySWBZbt0oA5XQ2RVwgzKGOgfMzZafsWAqrq1PGYjJqBbm/hudPtqsBridW/QjoE2Bp+Qnp8mWhdlSP8dgdeefLEeZGUSJx0Tzu2hBveEz7jMNQSOyg8HEE=","nonce":"eabb3e044577"},"expire_time":"2025-09-16T14:26:23+08:00","serial_no":"477ED0046A54F0360A72A63A8F2816312AAEAB53"}]}
`

					// mock certificates signature
					mockResp := &sign.ResponseSignature{
						Body:      []byte(mockBody),
						Timestamp: mockTimestamp,
						Nonce:     mockNonce,
					}
					plain, err := mockResp.Marshal()
					if err != nil {
						return nil, err
					}

					signature, err := sign.SignatureSHA256WithRSA(client.privateKey, plain)
					if err != nil {
						return nil, err
					}

					resp.Header = http.Header{}
					resp.Header.Set("Wechatpay-Nonce", mockNonce)
					resp.Header.Set("Wechatpay-Signature", signature)
					resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
					resp.Header.Set("Wechatpay-Serial", mockSerialNo)
					resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))

					return resp, nil
				},
			},
			true,
		},
		{
			&mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					var resp = &http.Response{
						StatusCode: http.StatusInternalServerError,
					}

					resp.Body = ioutil.NopCloser(strings.NewReader(`{"code":"ERROR_NAME","message":"ERROR_DESCRIPTION"}`))
					return resp, nil
				},
			},
			false,
		},
		{
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
		client.config.opts.transport = c.mocktransport
		client.secrets.clear()
		err := client.onceDownloadCertificates(ctx)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err %v", c.pass, pass, err)
		}
	}
}

func TestDownloadForClient(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		f      *FileUrl
		pass   bool
		expect string
	}{
		{
			f: &FileUrl{
				HashType:    "SHA1",
				HashValue:   "dcd7ceb3d382a1181798368bb15d8437de46c00f",
				DownloadUrl: "https://api.mch.weixin.qq.com/v3/billdownload/file?token=g44bIUH1GyQtE7ZmeTAPQx5b69qABpYuC_oZq6Aalf-gQP-lJ_FHRMLnyj2O8ujG",
			},
			pass: true,
			expect: "交易时间,公众账号ID,商户号,特约商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,应结订单金额,代金券金额,微信退款单号,商户退款单号,退款金额,充值券退款金额,退款类型,退款状态,商品名称,商户数据包,手续费,费率,订单金额,申请退款金额,费率备注\n" +
				"`2021-01-28 17:07:11,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000925202101284997714292,`S20210128170702357723,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 15:35:18,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000910202101282955148400,`S20210128153505214586,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 16:59:46,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000926202101281412639609,`S20210128165824499930,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"总交易单数,应结订单总金额,退款总金额,充值券退款总金额,手续费总金额,订单总金额,申请退款总金额\n" +
				"`3,`0.03,`0.00,`0.00,`0.00000,`0.03,`0.00\n",
		},
		{
			f: &FileUrl{
				HashType:    "SHA1",
				HashValue:   "dcd7ceb3d382a1181798368bb15d8437de46c00f",
				DownloadUrl: "https:\n//api.mch.weixin.qq.com/v3/billdownload/file?token=g44bIUH1GyQtE7ZmeTAPQx5b69qABpYuC_oZq6Aalf-gQP-lJ_FHRMLnyj2O8ujG",
			},
			pass: false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		data, err := client.Download(ctx, c.f)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err %v", c.pass, pass, err)
		}

		if err != nil {
			continue
		}

		actual := string(data)
		if c.expect != actual {
			t.Fatalf("expect %v, got %v", c.expect, actual)
		}
	}
}

func TestParseNotificationForClient(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		result *Result
		pass   bool
	}{
		{
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "Jook1G0Ex2xkvw5isZNY8Pvxj30X6HOCLNwMBh0wpRCU0LMTD+wQqHCENpYcsaMM/6vFMsRXtZnKldRk1dFmzpLOT8Rh1SwfMp/61oz7Eyh9+y1p2QkC2EW9dEnZk3gl7j5WcSsncy8ccM4ohfZVwQLslZwzKKaLxg5F5MTeiP/0ykYdFHOqIKdp9QMlly0Yb9aUXiVe19u3PEIOUkAawr9vD7EL5VHtnuer90ADrO9b+p4MAFxL1QfqshNhb4KeDjyVAzOqHjkThqAeuY1wv8KjoeVpZOxxrdSAoYcek2c2A8ywKWNMZi/k0Wwpu05UN498a39tKdHPZrqb6Qt4ZA==",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{"id":"b62e271c-3389-58a0-8146-4a704966e8f1","create_time":"2021-01-28T17:07:11+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"original_type":"transaction","algorithm":"AEAD_AES_256_GCM","ciphertext":"yuKJXXxnqVMulBUy5NoriSab/S9aen3wXNYLqGdvBfxsWmN9JAFAMXO3LgDFPqNeZMrkSmQyFa981IVxLvWHzwrzlBtJk+hOwnxTgDxc8SsGt39QkRBbfGR8rutMr3Goiq03ygWjMA6I+n6qhqQ/zS0/bMIB1dQoFZBSCKiLp8VHbGDLirh9MqYRa7MKJEYziPF2DmdtRHvXie4AWSxcV6hq8Ufao9FQooLOA2gD/9JA+L6BqquOPOnStExxH26cK7QgFFAf22GP7JKXnMH0LF3lJrK6ZMQ7iTXvVxv/q6j3SwUbyWVKmXdMJTqnXtU4H90DjRC6It4cOavr3Gz6xeVyv4S3i1qdAD8rAqgjjF1QWnUQtIm4/TdOw3ro0L73VI07H8c9O6VX/U0TcGMJJrAKMJ/yBZlD6owliffy/pzceEG/MV27euHDS5VW/m23tokNy2G1XJu1T3sUzEUsNil7vngBLYHGEGNw6brOYxwxXEUI2n0tSJOG8upiSGmN0fOnWbPoN9YqtuIhvY4xKOJpKwQrNJSm+ybNrugAwbLf/HMATxK6dGk9RQK8Nn9PHSRSPmTU5sci6zzFGAEHKQ==","associated_data":"transaction","nonce":"fG1l57vn9BCX"}}`),
			},
			true,
		},
		{
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "Jook1G0Ex2xkvw5isZNY8Pvxj30X6HOCLNwMBh0wpRCU0LMTD+wQqHCENpYcsaMM/6vFMsRXtZnKldRk1dFmzpLOT8Rh1SwfMp/61oz7Eyh9+y1p2QkC2EW9dEnZk3gl7j5WcSsncy8ccM4ohfZVwQLslZwzKKaLxg5F5MTeiP/0ykYdFHOqIKdp9QMlly0Yb9aUXiVe19u3PEIOUkAawr9vD7EL5VHtnuer90ADrO9b+p4MAFxL1QfqshNhb4KeDjyVAzOqHjkThqAeuY1wv8KjoeVpZOxxrdSAoYcek2c2A8ywKWNMZi/k0Wwpu05UN498a39tKdHPZrqb6Qt4ZA==",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{`),
			},
			false,
		},
		{
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{"id":"b62e271c-3389-58a0-8146-4a704966e8f1","create_time":"2021-01-28T17:07:11+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"original_type":"transaction","algorithm":"AEAD_AES_256_GCM","ciphertext":"yuKJXXxnqVMulBUy5NoriSab/S9aen3wXNYLqGdvBfxsWmN9JAFAMXO3LgDFPqNeZMrkSmQyFa981IVxLvWHzwrzlBtJk+hOwnxTgDxc8SsGt39QkRBbfGR8rutMr3Goiq03ygWjMA6I+n6qhqQ/zS0/bMIB1dQoFZBSCKiLp8VHbGDLirh9MqYRa7MKJEYziPF2DmdtRHvXie4AWSxcV6hq8Ufao9FQooLOA2gD/9JA+L6BqquOPOnStExxH26cK7QgFFAf22GP7JKXnMH0LF3lJrK6ZMQ7iTXvVxv/q6j3SwUbyWVKmXdMJTqnXtU4H90DjRC6It4cOavr3Gz6xeVyv4S3i1qdAD8rAqgjjF1QWnUQtIm4/TdOw3ro0L73VI07H8c9O6VX/U0TcGMJJrAKMJ/yBZlD6owliffy/pzceEG/MV27euHDS5VW/m23tokNy2G1XJu1T3sUzEUsNil7vngBLYHGEGNw6brOYxwxXEUI2n0tSJOG8upiSGmN0fOnWbPoN9YqtuIhvY4xKOJpKwQrNJSm+ybNrugAwbLf/HMATxK6dGk9RQK8Nn9PHSRSPmTU5sci6zzFGAEHKQ==","associated_data":"transaction","nonce":"fG1l57vn9BCX"}}`),
			},
			false,
		},
		{
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "g0A/VGU569/iS8MtR2SRfFg0YOqSzKYipRTJebnm6bLsWgSwWL92KMoAwNrP8qMqf1LUKWWrb2o0XpmLt2DMV7MStrJNmcViV6yVKzVRuS2XE3kUiQNFnbIdvNRiLI0hLDGA9W6dH5YUF/yVPanRo3rBLID8mFxD1tz2XyVpVKsDu7EhUQmwCKpgZ0a+lPILxZfMjnVI7VL6AFuf/iCrb/xaoVzGCJ1hLcPe7QV89MqNp2M4IP1wbiBqJezC7vBF/t/Rtyn+kK+My/S7iB+XDrHHXn/6ldp7RXBcDmjVbnp551oS2s8jyBFN1z/K+BIg+gYmyN9vOGgRFWcV2NGpUQ==",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{"id":"b62e271c-3389-58a0-8146-4a704966e8f1","create_time":"2021-01-28T17:07:11+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"original_type":"transaction","algorithm":"AEAD_AES_256_GCM","ciphertext":"yuKJXXxnqVMulBUy5NoriSab/S9aen3wXNYLqGdvBfxsWmN9JAFAMXO3LgDFPqNeZMrkSmQyFa981IVxLvWHzwrzlBtJk+hOwnxTgDxc8SsGt39QkRBbfGR8rutMr3Goiq03ygWjMA6I+n6qhqQ/zS0/bMIB1dQoFZBSCKiLp8VHbGDLirh9MqYRa7MKJEYziPF2DmdtRHvXie4AWSxcV6hq8Ufao9FQooLOA2gD/9JA+L6BqquOPOnStExxH26cK7QgFFAf22GP7JKXnMH0LF3lJrK6ZMQ7iTXvVxv/q6j3SwUbyWVKmXdMJTqnXtU4H90DjRC6It4cOavr3Gz6xeVyv4S3i1qdAD8rAqgjjF1QWnUQtIm4/TdOw3ro0L73VI07H8c9O6VX/U0TcGMJJrAKMJ/yBZlD6owliffy/pzceEG/MV27euHDS5VW/m23tokNy2G1XJu1T3sUzEUsNil7vngBLYHGEGNw6brOYxwxXEUI2n0tSJOG8upiSGmN0fOnWbPoN9YqtuIhvY4xKOJpKwQrNJSm+ybNrugAwbLf/HMATxK6dGk9RQK8Nn9PkHYuBnwDft8oxSDkqLO7KA==","associated_data":"transaction","nonce":"fG1l57vn9BCX"}}`),
			},
			false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		_, _, err := client.ParseNotification(ctx, c.result)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err %v", c.pass, pass, err)
		}
	}
}

func TestGenRequestSignature(t *testing.T) {
	cases := []struct {
		method string
		url    string
		body   []byte
	}{
		{
			"POST",
			"https://api.mch.weixin.qq.com/v3/pay/transactions/native",
			[]byte(`{"appid":"wx81be3101902f7cb2","mchid":"1601959334","description":"for testing","out_trade_no":"S20210124144305172434","time_expire":"2021-01-24T14:53:05+08:00","attach":"cipher code","notify_url":"https://luoji.live/notify","amount":{"total":1,"currency":"CNY"},"detail":{},"scene_info":{"payer_client_ip":"","store_info":{"id":""}}}`),
		},
	}

	for _, c := range cases {
		req := genRequestSignature(c.method, c.url, c.body)
		if req == nil {
			t.Fatal("req is nil")
		}
	}
}

func TestSecrets(t *testing.T) {
	cases := []struct {
		secrets *secrets
		expect  bool
	}{
		{
			&secrets{},
			true,
		},
		{
			&secrets{
				all: map[string]*rsa.PublicKey{
					"m": {},
				},
			},
			true,
		},
		{
			&secrets{
				deadline: time.Now().Add(time.Minute),
				all:      map[string]*rsa.PublicKey{},
			},
			true,
		},
		{
			&secrets{
				deadline: time.Now().Add(time.Minute),
				all: map[string]*rsa.PublicKey{
					"m": {},
				},
			},
			false,
		},
	}

	for _, c := range cases {
		// c.secrets.clear()
		actual := c.secrets.isUpgrade()
		if actual != c.expect {
			t.Fatalf("expect %v, got %v", c.expect, actual)
		}
	}
}

func TestSecretsWithGoroutine(t *testing.T) {
	var secrets secrets
	secrets.clear()

	cases := []struct {
		expect bool
	}{
		{false},
		{false},
	}

	actual := []bool{false, false}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		secrets.add("m", &rsa.PublicKey{}, time.Minute)
		secrets.add("m1", &rsa.PublicKey{}, time.Minute)
		wg.Done()
	}()

	go func() {
		secrets.add("m", &rsa.PublicKey{}, time.Minute)
		secrets.add("m2", &rsa.PublicKey{}, time.Minute)
		wg.Done()
	}()

	wg.Wait()
	wg.Add(2)
	go func() {
		isUpgrade := secrets.isUpgrade()
		actual[0] = isUpgrade
		wg.Done()
	}()

	go func() {
		isUpgrade := secrets.isUpgrade()
		actual[1] = isUpgrade
		wg.Done()
	}()

	wg.Wait()
	for i, c := range cases {
		if actual[i] != c.expect {
			t.Fatalf("expect %v, got %v", c.expect, actual[i])
		}
	}
}

func ExampleNewClient() {
	appId := "wxd678efh567hg6787"
	mchId := "1230000109"
	apiv3Secret := "AES256Key-32Characters1234567890"
	serialNo := "477ED0046A54F0360A72A63A8F2816312AAEAB53"
	privateKeyTxt := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCprsmcXPHqLtnP
oPDGUoMULK2WOo5FW8c72Svnqn/4aXPaJhlOtPxtX2frqIhTjwcOs6hNm3XFTGBL
MrdB94YQvj+Q7P12GNmxXG+9Ms+uUyJToYjlYDAG6UFKE10Jkm9cDGuLSkekU1Ao
rKE1G1wndH37w4AzVXoGBQ3NIiyW8jIm8Zi3/WNCVpHUoXYUuyhFEZ23fXytnps4
hARgg6NvPncIKtWvlUh85ZVOSsqc1T8dFaeDRXaj7r3jdJJ74tsGRMvZyUipJXyE
3uR2QkrGyia+0phDpC6zeMMpP+MQO9ohh+xQWBCeyvQjjnOPAlGThl+ThfXImU30
HL17oHdBAgMBAAECggEADm6FSz1Efgx6DgS8NcHy0BZ0tSBJ1XBW46o2579Cnxgo
+FbhNCaEibDhn9N3tNOnYAK7v84HGD7EueCYYY3x4x6rPWJKtG6spT8dadQWgdck
RkSo5glmTFAuc2RuN1AzFHsh8njg2wMTAEKee2vWTKzFwlIAZ11PwY9Qey/65uOT
Bi8q1Rssu6xofNadO5MbqMJ1Tl8DDIaLGnzTzbHrk9thBUo1FwFjJWTVI7nz2En4
Yc/G1/LQJfiQ31F+lkL3j6ABRJqtsgb07r9H/hT6+fd1hGDt2qKuS+E1mLDp9fHw
n6UyS4HyB7DA/XtFZ9z0VtAlmcGoUkyJLtXjEmwsGQKBgQDeCtE3spULpC7VPqk1
xv034C6zybZ7y8kSKwRvyYwkzdgSRgVaKsTVb8RNYor8hoGrVgdXFqQUI8O/v1cN
9wFoGYJT0LHre/YzOg31TkQkBfHHCFH/L50uOJcIQueftctz5Bwj6bJO/ih5iIAK
yjrHse4PdIiEJfz2D9hc4wnxrwKBgQDDogrWlUCTj2fvmZfkWR3Hbs0kIHd7zjIk
bJJONGtD8gE4i562tajC1mKoQEwt4YSwWsBkGAw1LhvMROQFT6AOaIIhHNex1Z3t
c2gAdEeWOMmzZnnhwWzTiYJomixrFkmEwT3EJK89GO3E0FH5S+G1P1tNXq38Vpty
1YVqOgMSDwKBgFrzuWGEQDMljJ2C7lL98KlbpiW1AY/SGMndXxLfTw2gV9qcXgLi
NABtqM4+CEqKWkExmw4cUxeA0uUPXnx06lmW4WCtwsN/4oh3RlJuPdE3siLiEJxk
B5FwUsVqinBMSktta+12A7kBuNiXhkNlNRCpnKcuB+GBog20zd62jVM3AoGBALcA
zFazQ7dFfRq7eUUYwCyhT7Et1dewqWM9VRdnHbhvmAjHQu7zvCyW069Ehn6c6bz3
B+YaQME2orZQ82SsebNAvAoxquwmQhevz2gtXhH+iWASyo0Onbi8d4tWPZrnPFq9
UgQ7tNnYigOEREqKW1drLwOPP/4/Hicr6iPWpKytAoGAEQ6J/RB/olEAC46ACoFo
FBgA+GUbDB0xBcA2inEt3q//208YMkjnKM871n89HpAgms5xrK32T69lduebk7Ar
9wWvkJVUwI9VDXomCFQqtiGzHlTl1Xq31BfeIDyq1ayQmTkRpRqIagbDZVtM+ha/
0I2SEzTObt07wcYcYG2Chvg=
-----END PRIVATE KEY-----`
	client, _ := NewClient(
		Config{
			AppId:       appId,
			MchId:       mchId,
			Apiv3Secret: apiv3Secret,
			Cert: CertSuite{
				SerialNo:      serialNo,
				PrivateKeyTxt: privateKeyTxt,
			},
		})

	fmt.Println(client != nil)
	// Output:
	// true
}
