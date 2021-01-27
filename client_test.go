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
	"strconv"
	"strings"
	"testing"

	"github.com/gunsluo/wechatpay-go/v3/sign"
)

func mockNewClient() (*client, error) {
	var (
		appId          = "wxd678efh567hg6787"
		mchId          = "1230000109"
		apiv3Secret    = "AES256Key-32Characters1234567890"
		serialNo       = mockSerialNo
		privateKeyPath = "./test_fixtures/mock_private_key_pkcs8.pem"
	)

	// use mock data
	privateKey, err := sign.LoadRSAPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	mocktransport := &mockTransport{
		RoundTripFn: func(req *http.Request) (*http.Response, error) {
			return mockData(req, privateKey)
		},
	}

	client, err := newClient(
		Config{
			AppId:       appId,
			MchId:       mchId,
			Apiv3Secret: apiv3Secret,
			Cert: CertSuite{
				SerialNo:       serialNo,
				PrivateKeyPath: privateKeyPath,
			},
		}, Transport(mocktransport))
	if err != nil {
		return nil, err
	}

	// mock request signature
	client.genRequestSignature = func(method, url string, body []byte) *sign.RequestSignature {
		return &sign.RequestSignature{
			Method:    method,
			Timestamp: mockTimestamp,
			Url:       url,
			Nonce:     mockNonce,
			Body:      body,
		}
	}

	return client, nil
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

	cases := []struct {
		req    *sign.RequestSignature
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
			`WECHATPAY2-SHA256-RSA2048 mchid="1230000109",nonce_str="AF1404CC2980FB414C99C0B98883BD42",signature="ItuRCG6nAf6ZUi5C5LPa0beCGrG7+G4NdaCHLTmym+UzuZHFgFeqRZ4zKQ0n93qehchFWfQ7s00pgABYvXcOMsV1ld7AUjDTZBPucJK6yhFKz9jd20wtRdDG4LRCZcaTowD2f7LtlixFm8F3/YQaBavxiOe54tc3RX/22flYRzy4YFOpBt+bmjSPZIdSFi53323u7cohwvdHwX+avQCtLZKAUNFJIob66u05BbDEITzYuHjakjpb5btvWemjoZBPxkiETzmd4Oa1y2U+rfFCPZyWT4EV7UxHeEizBL8DkubEBD3KXeArqRX6yoMAU4ywmdFeWDbv1EF0Ndy9hiddZQ==",timestamp="1611368330",serial_no="477ED0046A54F0360A72A63A8F2816312AAEAB53"`,
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

	cases := []struct {
		result *Result
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
		},
	}

	for _, c := range cases {
		err := client.upgradeCertificate([]byte(c.result.Body))
		if err != nil {
			t.Fatal(err)
		}

		err = client.VerifySignature(c.result)
		if err != nil {
			t.Fatal(err)
		}
	}
}

var (
	mockTimestamp int64 = 1611368330
	mockNonce           = "AF1404CC2980FB414C99C0B98883BD42"
	mockSerialNo        = "477ED0046A54F0360A72A63A8F2816312AAEAB53"
)

type mockTransport struct {
	RoundTripFn func(req *http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripFn(req)
}

func mockData(req *http.Request, privateKey *rsa.PrivateKey) (*http.Response, error) {
	path := req.URL.Path

	var resp = &http.Response{
		StatusCode: http.StatusOK,
	}
	switch path {
	case "/v3/certificates":
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

		signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
		if err != nil {
			return nil, err
		}

		resp.Header = http.Header{}
		resp.Header.Set("Wechatpay-Nonce", mockNonce)
		resp.Header.Set("Wechatpay-Signature", signature)
		resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
		resp.Header.Set("Wechatpay-Serial", mockSerialNo)
		resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))
	case "/v3/pay/transactions/native":
		mockBody := `{"code_url":"weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00"}`

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

		signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
		if err != nil {
			return nil, err
		}

		resp.Header = http.Header{}
		resp.Header.Set("Wechatpay-Nonce", mockNonce)
		resp.Header.Set("Wechatpay-Signature", signature)
		resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
		resp.Header.Set("Wechatpay-Serial", mockSerialNo)
		resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))
	default:
		resp.Body = ioutil.NopCloser(strings.NewReader(`{}`))
	}

	return resp, nil
}
