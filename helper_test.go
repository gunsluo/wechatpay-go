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
	"bytes"
	"crypto/rsa"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gunsluo/wechatpay-go/v3/sign"
)

const (
	mockAppId          = "wxd678efh567hg6787"
	mockMchId          = "1230000109"
	mockApiv3Secret    = "AES256Key-32Characters1234567890"
	mockSerialNo       = "477ED0046A54F0360A72A63A8F2816312AAEAB53"
	mockPrivateKeyPath = "./test_fixtures/mock_private_key_pkcs8.pem"

	mockTimestamp int64 = 1611368330
	mockNonce           = "AF1404CC2980FB414C99C0B98883BD42"
)

type mockTransport struct {
	RoundTripFn func(req *http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripFn(req)
}

func mockGenRequestSignature(method, url string, body []byte) *sign.RequestSignature {
	return &sign.RequestSignature{
		Method:    method,
		Timestamp: mockTimestamp,
		Url:       url,
		Nonce:     mockNonce,
		Body:      body,
	}
}

func mockNewClient(transports ...*mockTransport) (*client, error) {
	var (
		appId          = mockAppId
		mchId          = mockMchId
		apiv3Secret    = mockApiv3Secret
		serialNo       = mockSerialNo
		privateKeyPath = mockPrivateKeyPath
	)

	var transport *mockTransport
	if len(transports) > 0 {
		transport = transports[0]
	} else {
		/*
			privateKey, err := sign.LoadRSAPrivateKeyFromFile(privateKeyPath)
			if err != nil {
				return nil, err
			}

			transport = &mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					return defaultMockData(req, privateKey)
				},
			}
		*/
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
		},
		Transport(transport),
		Timeout(time.Minute),
		CertRefreshTime(10*time.Minute),
	)
	if err != nil {
		return nil, err
	}

	if client.config.opts.transport == nil {
		client.config.opts.transport = &mockTransport{
			RoundTripFn: func(req *http.Request) (*http.Response, error) {
				return defaultMockData(req, client.privateKey)
			},
		}
	}

	// mock request signature
	client.genRequestSignature = mockGenRequestSignature
	return client, nil
}

var defaultMockDataMapping = map[string]func(*http.Request, *http.Response, *rsa.PrivateKey) error{
	"/v3/certificates":            mockDataWithCert,
	"/v3/pay/transactions/native": mockDataWithPay,
	"/v3/pay/transactions/app":    mockDataWithPay,
	"/v3/pay/transactions/h5":     mockDataWithPay,
	"/v3/pay/transactions/jsapi":  mockDataWithPay,

	"/v3/combine-transactions/native": mockDataWithCombinPay,
	"/v3/combine-transactions/app":    mockDataWithCombinPay,
	"/v3/combine-transactions/h5":     mockDataWithCombinPay,
	"/v3/combine-transactions/jsapi":  mockDataWithCombinPay,

	"/v3/pay/transactions/id/4200000914202101195554393855":          mockDataWithQueryPay,
	"/v3/pay/transactions/out-trade-no/S20210119074247105778399200": mockDataWithQueryPay,
	"/v3/pay/transactions/out-trade-no/S20210119NOTFOUND":           mockDataWithNotFoundQueryPay,
	"/v3/refund/domestic/refunds":                                   mockDataWithRefund,
	"/v3/pay/transactions/out-trade-no/fortest/close":               mockDataWithClose,
	"/v3/refund/domestic/refunds/1217752501201407033233368018":      mockDataWithQueryRefund,
	"/v3/billdownload/file":                                         mockDataWithDownloadFile,
	"/v3/bill/tradebill":                                            mockDataWithTradeBill,
	"/v3/bill/fundflowbill":                                         mockDataWithFundflowBill,
	"/v3/invalidresp":                                               mockDataWithInvalidResp,
	"/v3/invalidrespdata":                                           mockDataWithInvalidRespData,
	"/v3/invalidheader":                                             mockDataWithInvalidHeader,

	"/v3/combine-transactions/out-trade-no/fortest/close": mockDataWithClose,
}

func defaultMockData(req *http.Request, privateKey *rsa.PrivateKey) (*http.Response, error) {
	path := req.URL.Path

	var resp = &http.Response{
		StatusCode: http.StatusOK,
	}

	rundTripFn, ok := defaultMockDataMapping[path]
	if !ok {
		resp.Body = ioutil.NopCloser(strings.NewReader(`{}`))
		return resp, nil
	}

	err := rundTripFn(req, resp, privateKey)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func mockDataWithCert(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
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
		return err
	}

	signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
	if err != nil {
		return err
	}

	resp.Header = http.Header{}
	resp.Header.Set("Wechatpay-Nonce", mockNonce)
	resp.Header.Set("Wechatpay-Signature", signature)
	resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
	resp.Header.Set("Wechatpay-Serial", mockSerialNo)
	resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))

	return nil
}

func mockDataWithPay(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	mockBody := `{"code_url":"weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00"}`

	// mock certificates signature
	mockResp := &sign.ResponseSignature{
		Body:      []byte(mockBody),
		Timestamp: mockTimestamp,
		Nonce:     mockNonce,
	}
	plain, err := mockResp.Marshal()
	if err != nil {
		return err
	}

	signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
	if err != nil {
		return err
	}

	resp.Header = http.Header{}
	resp.Header.Set("Wechatpay-Nonce", mockNonce)
	resp.Header.Set("Wechatpay-Signature", signature)
	resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
	resp.Header.Set("Wechatpay-Serial", mockSerialNo)
	resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))

	return nil
}

func mockDataWithCombinPay(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	return mockDataWithPay(req, resp, privateKey)
}

func mockDataWithQueryPay(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	mockBody := `{"appid":"wxd678efh567hg6787","mchid":"1230000109","out_trade_no":"S20210119074247105778399200","transaction_id":"4200000914202101195554393855","trade_type":"NATIVE","trade_state":"SUCCESS","trade_state_desc":"支付成功","bank_type":"OTHERS","success_time":"2021-01-19T15:43:01+08:00","payer":{"openid":"ofyak5qYxYJVnhTlrkk_ACWIVrHI"},"amount":{"total":1,"payer_total":1,"currency":"CNY","payer_currency":"CNY"}}`
	// mock certificates signature
	mockResp := &sign.ResponseSignature{
		Body:      []byte(mockBody),
		Timestamp: mockTimestamp,
		Nonce:     mockNonce,
	}
	plain, err := mockResp.Marshal()
	if err != nil {
		return err
	}

	signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
	if err != nil {
		return err
	}

	resp.Header = http.Header{}
	resp.Header.Set("Wechatpay-Nonce", mockNonce)
	resp.Header.Set("Wechatpay-Signature", signature)
	resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
	resp.Header.Set("Wechatpay-Serial", mockSerialNo)
	resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))

	return nil
}

func mockDataWithNotFoundQueryPay(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	mockBody := `{"status":404,"code":"ORDER_NOT_EXIST","message":"订单不存在"}`
	// mock certificates signature
	mockResp := &sign.ResponseSignature{
		Body:      []byte(mockBody),
		Timestamp: mockTimestamp,
		Nonce:     mockNonce,
	}
	plain, err := mockResp.Marshal()
	if err != nil {
		return err
	}

	signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
	if err != nil {
		return err
	}
	resp.Header = http.Header{}
	resp.Header.Set("Wechatpay-Nonce", mockNonce)
	resp.Header.Set("Wechatpay-Signature", signature)
	resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
	resp.Header.Set("Wechatpay-Serial", mockSerialNo)
	resp.StatusCode = http.StatusNotFound
	resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))

	return nil
}

func mockDataWithRefund(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	mockBody := `{ "refund_id": "50300807092021020105990201735", "out_refund_no": "S20210201151309277501", "transaction_id": "4200000925202101284997714292", "out_trade_no": "S20210128170702357723", "channel": "ORIGINAL", "user_received_account": "支付用户零钱", "success_time": "0001-01-01T00:00:00Z", "create_time": "2021-02-01T15:13:10+08:00", "status": "PROCESSING", "funds_account": "UNAVAILABLE", "amount": { "total": 1, "refund": 1, "payer_total": 1, "payer_refund": 1, "settlement_total": 1, "settlement_refund": 1, "discount_refund": 0, "currency": "CNY" } }`

	// mock certificates signature
	mockResp := &sign.ResponseSignature{
		Body:      []byte(mockBody),
		Timestamp: mockTimestamp,
		Nonce:     mockNonce,
	}
	plain, err := mockResp.Marshal()
	if err != nil {
		return err
	}

	signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
	if err != nil {
		return err
	}
	resp.Header = http.Header{}
	resp.Header.Set("Wechatpay-Nonce", mockNonce)
	resp.Header.Set("Wechatpay-Signature", signature)
	resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
	resp.Header.Set("Wechatpay-Serial", mockSerialNo)
	resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))

	return nil
}

func mockDataWithClose(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	resp.Header = http.Header{}
	resp.StatusCode = 204
	mockBody := ``
	// mock certificates signature
	mockResp := &sign.ResponseSignature{
		Body:      []byte(mockBody),
		Timestamp: mockTimestamp,
		Nonce:     mockNonce,
	}
	plain, err := mockResp.Marshal()
	if err != nil {
		return err
	}

	signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
	if err != nil {
		return err
	}
	resp.Header.Set("Wechatpay-Nonce", mockNonce)
	resp.Header.Set("Wechatpay-Signature", signature)
	resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
	resp.Header.Set("Wechatpay-Serial", mockSerialNo)
	resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))

	return nil
}

func mockDataWithQueryRefund(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	mockBody := `{"refund_id":"50000000382019052709732678859","out_refund_no":"1217752501201407033233368018","transaction_id":"1217752501201407033233368018","out_trade_no":"1217752501201407033233368018","channel":"ORIGINAL","user_received_account":"招商银行信用卡0403","success_time":"2020-12-01T16:18:12+08:00","create_time":"2020-12-01T16:18:12+08:00","status":"SUCCESS","funds_account":"UNSETTLED","amount":{"total":100,"refund":100,"payer_total":90,"payer_refund":90,"settlement_refund":100,"settlement_total":100,"discount_refund":10,"currency":"CNY"},"promotion_detail":[{"promotion_id":"109519","scope":"SINGLE","type":"DISCOUNT","amount":5,"refund_amount":100,"goods_detail":[{"merchant_goods_id":"1217752501201407033233368018","wechatpay_goods_id":"1001","goods_name":"iPhone6s 16G","unit_price":528800,"refund_amount":528800,"refund_quantity":1}]}]}`

	resp.Header = http.Header{}
	resp.StatusCode = 200
	// mock certificates signature
	mockResp := &sign.ResponseSignature{
		Body:      []byte(mockBody),
		Timestamp: mockTimestamp,
		Nonce:     mockNonce,
	}
	plain, err := mockResp.Marshal()
	if err != nil {
		return err
	}

	signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
	if err != nil {
		return err
	}
	resp.Header.Set("Wechatpay-Nonce", mockNonce)
	resp.Header.Set("Wechatpay-Signature", signature)
	resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
	resp.Header.Set("Wechatpay-Serial", mockSerialNo)
	resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))

	return nil
}

func mockDataWithDownloadFile(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	vs := req.URL.Query()
	billType := vs.Get("bill_type")
	accountType := vs.Get("account_type")
	tarType := vs.Get("tar_type")

	var reader io.Reader
	if accountType == "" {
		switch billType {
		case "REFUND":
		case "SUCCESS":
		case "ALL":
			fallthrough
		default:
			if tarType == "GZIP" {
				mockBody := []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 212, 84, 65, 79, 219, 48, 24, 189, 243, 43, 184, 236, 246, 129, 108, 39, 78, 226, 220, 80, 135, 52, 38, 141, 73, 148, 109, 226, 52, 3, 43, 27, 154, 52, 54, 64, 98, 219, 41, 28, 74, 96, 208, 49, 84, 162, 114, 217, 52, 88, 69, 57, 116, 165, 136, 170, 165, 13, 234, 254, 76, 237, 36, 255, 98, 74, 210, 54, 13, 55, 110, 91, 43, 89, 126, 159, 191, 60, 191, 239, 41, 47, 189, 78, 89, 158, 20, 101, 169, 25, 148, 26, 32, 242, 213, 222, 109, 201, 111, 156, 139, 195, 214, 204, 67, 16, 206, 142, 220, 109, 137, 195, 22, 120, 123, 109, 175, 115, 158, 96, 191, 214, 21, 101, 59, 220, 137, 110, 173, 247, 231, 212, 175, 157, 137, 130, 19, 225, 168, 39, 193, 222, 241, 133, 220, 109, 201, 159, 182, 127, 185, 3, 241, 101, 222, 149, 43, 126, 236, 15, 192, 151, 166, 180, 182, 161, 231, 158, 200, 106, 55, 40, 54, 253, 211, 3, 240, 27, 21, 113, 179, 237, 85, 10, 222, 149, 11, 162, 115, 236, 185, 197, 152, 48, 176, 143, 130, 179, 239, 208, 115, 127, 5, 246, 145, 216, 189, 233, 227, 88, 67, 96, 89, 178, 218, 77, 201, 72, 149, 98, 48, 120, 36, 159, 23, 214, 109, 72, 49, 90, 141, 65, 95, 95, 31, 196, 250, 132, 179, 35, 138, 219, 226, 91, 193, 171, 212, 251, 228, 210, 169, 203, 66, 77, 28, 228, 65, 238, 237, 123, 238, 111, 191, 209, 6, 191, 209, 246, 190, 218, 144, 82, 235, 29, 95, 251, 151, 173, 212, 61, 113, 155, 40, 219, 242, 250, 98, 140, 19, 68, 240, 4, 194, 19, 196, 24, 199, 186, 137, 116, 19, 99, 224, 91, 31, 13, 188, 148, 83, 48, 194, 12, 145, 21, 125, 121, 137, 0, 199, 26, 194, 140, 50, 69, 81, 129, 35, 224, 192, 85, 130, 194, 31, 35, 52, 228, 64, 152, 24, 42, 99, 186, 142, 85, 194, 8, 240, 236, 160, 136, 117, 164, 35, 162, 80, 93, 39, 10, 240, 181, 149, 79, 139, 111, 233, 135, 185, 151, 120, 107, 97, 35, 195, 88, 102, 227, 197, 148, 54, 199, 158, 60, 94, 252, 60, 5, 124, 118, 106, 126, 230, 249, 52, 240, 236, 179, 76, 102, 58, 155, 5, 254, 116, 254, 209, 244, 92, 22, 120, 102, 118, 1, 56, 154, 68, 56, 90, 81, 164, 0, 13, 247, 209, 26, 254, 87, 214, 214, 199, 55, 115, 27, 155, 171, 239, 94, 3, 95, 94, 125, 255, 38, 183, 62, 190, 188, 246, 42, 23, 247, 160, 176, 13, 79, 34, 244, 32, 205, 149, 54, 129, 154, 10, 53, 177, 113, 79, 19, 48, 26, 204, 75, 24, 165, 88, 53, 212, 144, 57, 49, 129, 42, 20, 81, 130, 85, 106, 104, 255, 131, 9, 154, 73, 153, 169, 106, 247, 125, 19, 180, 225, 188, 42, 38, 154, 194, 52, 196, 70, 77, 208, 168, 65, 84, 149, 49, 166, 160, 127, 212, 4, 105, 185, 241, 151, 65, 20, 28, 233, 212, 83, 249, 151, 150, 155, 74, 106, 130, 239, 228, 57, 57, 24, 102, 51, 41, 221, 229, 26, 13, 232, 176, 58, 198, 149, 72, 146, 146, 30, 108, 40, 61, 57, 26, 251, 27, 0, 0, 255, 255, 36, 43, 30, 24, 67, 5, 0, 0}
				reader = bytes.NewReader(mockBody)
			} else {
				mockBody := "交易时间,公众账号ID,商户号,特约商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,应结订单金额,代金券金额,微信退款单号,商户退款单号,退款金额,充值券退款金额,退款类型,退款状态,商品名称,商户数据包,手续费,费率,订单金额,申请退款金额,费率备注\n" +
					"`2021-01-28 17:07:11,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000925202101284997714292,`S20210128170702357723,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
					"`2021-01-28 15:35:18,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000910202101282955148400,`S20210128153505214586,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
					"`2021-01-28 16:59:46,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000926202101281412639609,`S20210128165824499930,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
					"总交易单数,应结订单总金额,退款总金额,充值券退款总金额,手续费总金额,订单总金额,申请退款总金额\n" +
					"`3,`0.03,`0.00,`0.00,`0.00000,`0.03,`0.00\n"
				reader = strings.NewReader(mockBody)
			}
		}
	} else {
		if tarType == "GZIP" {
			mockBody := []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 172, 146, 207, 110, 149, 64, 20, 198, 247, 60, 133, 75, 77, 166, 55, 231, 204, 63, 102, 112, 229, 210, 157, 137, 47, 48, 143, 66, 115, 37, 168, 77, 195, 109, 164, 105, 213, 104, 82, 83, 255, 36, 182, 5, 83, 90, 20, 226, 245, 101, 152, 1, 222, 194, 112, 129, 171, 215, 149, 139, 178, 97, 190, 243, 157, 156, 243, 253, 50, 211, 93, 229, 93, 241, 201, 157, 220, 246, 39, 5, 177, 235, 171, 230, 215, 153, 75, 179, 166, 62, 109, 190, 191, 177, 47, 207, 236, 225, 177, 77, 74, 210, 221, 60, 235, 227, 35, 119, 179, 239, 242, 98, 42, 77, 254, 234, 176, 253, 156, 79, 162, 253, 86, 219, 247, 7, 196, 165, 183, 46, 205, 118, 68, 31, 31, 245, 31, 222, 221, 183, 209, 242, 1, 25, 214, 61, 47, 219, 250, 85, 243, 243, 245, 84, 217, 76, 183, 201, 169, 123, 91, 184, 100, 213, 84, 231, 109, 122, 221, 101, 101, 83, 85, 196, 158, 199, 238, 250, 203, 188, 45, 190, 236, 178, 125, 155, 148, 158, 161, 64, 113, 15, 232, 30, 224, 61, 100, 129, 224, 1, 32, 49, 2, 24, 128, 2, 169, 229, 96, 3, 5, 4, 161, 125, 165, 53, 215, 82, 17, 195, 41, 12, 159, 166, 176, 177, 17, 181, 175, 37, 103, 168, 169, 226, 196, 244, 97, 232, 46, 214, 127, 254, 46, 205, 108, 92, 17, 3, 139, 97, 52, 44, 40, 37, 6, 37, 160, 22, 154, 49, 254, 232, 201, 227, 185, 213, 133, 245, 200, 55, 116, 218, 104, 249, 208, 174, 190, 186, 23, 7, 109, 125, 217, 21, 63, 96, 1, 96, 163, 37, 49, 79, 231, 72, 200, 4, 19, 146, 41, 212, 28, 119, 57, 120, 0, 16, 112, 49, 113, 104, 240, 129, 253, 195, 161, 124, 132, 45, 135, 226, 114, 203, 193, 37, 42, 6, 76, 251, 255, 193, 129, 119, 196, 193, 1, 56, 23, 130, 42, 46, 189, 191, 31, 136, 11, 235, 246, 34, 117, 199, 249, 112, 249, 54, 250, 184, 35, 198, 13, 100, 12, 181, 117, 6, 49, 58, 158, 97, 196, 224, 156, 150, 110, 14, 212, 251, 29, 0, 0, 255, 255, 22, 13, 183, 141, 166, 2, 0, 0}
			reader = bytes.NewReader(mockBody)
		} else {
			mockBody := "记账时间,微信支付业务单号,资金流水单号,业务名称,业务类型,收支类型,收支金额(元),账户结余(元),资金变更提交申请人,备注,业务凭证号\n" +
				"`2021-02-01 13:54:01,`50300806962021020105978994968,`4200000920202101197964319284,`退款,`退款,`支出,`0.01,`0.22,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201135356381941\n" +
				"`2021-02-01 14:00:45,`50300907032021020105978998710,`4200000846202101197461830397,`退款,`退款,`支出,`0.01,`0.21,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201140044552846\n" +
				"资金流水总笔数,收入笔数,收入金额,支出笔数,支出金额\n" +
				"`3,`1,`0.01,`2,`0.02\n"
			reader = strings.NewReader(mockBody)
		}
	}

	resp.Body = ioutil.NopCloser(reader)

	return nil
}

func mockDataWithTradeBill(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	vs := req.URL.Query()
	fileUrl := "https://api.mch.weixin.qq.com/v3/billdownload/file?token=g44bIUH1GyQtE7ZmeTAPQx5b69qABpYuC_oZq6Aalf-gQP-lJ_FHRMLnyj2O8ujG"

	fileUrl += "&bill_type=" + vs.Get("bill_type")
	fileUrl += "&tar_type=" + vs.Get("tar_type")

	mockBody := `{"hash_type":"SHA1","hash_value":"dcd7ceb3d382a1181798368bb15d8437de46c00f","download_url":"` + fileUrl + `"}`

	resp.Header = http.Header{}
	resp.StatusCode = 200
	// mock certificates signature
	mockResp := &sign.ResponseSignature{
		Body:      []byte(mockBody),
		Timestamp: mockTimestamp,
		Nonce:     mockNonce,
	}
	plain, err := mockResp.Marshal()
	if err != nil {
		return err
	}

	signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
	if err != nil {
		return err
	}
	resp.Header.Set("Wechatpay-Nonce", mockNonce)
	resp.Header.Set("Wechatpay-Signature", signature)
	resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
	resp.Header.Set("Wechatpay-Serial", mockSerialNo)
	resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))

	return nil
}

func mockDataWithFundflowBill(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	vs := req.URL.Query()
	accountType := vs.Get("account_type")
	if accountType == "" {
		accountType = "BASIC"
	}

	fileUrl := "https://api.mch.weixin.qq.com/v3/billdownload/file?token=g44bIUH1GyQtE7ZmeTAPQx5b69qABpYuC_oZq6Aalf-gQP-lJ_FHRMLnyj2O8ujG"
	fileUrl += "&account_type=" + accountType
	fileUrl += "&tar_type=" + vs.Get("tar_type")

	mockBody := `{"hash_type":"SHA1","hash_value":"dcd7ceb3d382a1181798368bb15d8437de46c00f","download_url":"` + fileUrl + `"}`

	resp.Header = http.Header{}
	resp.StatusCode = 200
	// mock certificates signature
	mockResp := &sign.ResponseSignature{
		Body:      []byte(mockBody),
		Timestamp: mockTimestamp,
		Nonce:     mockNonce,
	}
	plain, err := mockResp.Marshal()
	if err != nil {
		return err
	}

	signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
	if err != nil {
		return err
	}
	resp.Header.Set("Wechatpay-Nonce", mockNonce)
	resp.Header.Set("Wechatpay-Signature", signature)
	resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
	resp.Header.Set("Wechatpay-Serial", mockSerialNo)
	resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))

	return nil
}

func mockDataWithInvalidResp(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	resp.StatusCode = http.StatusInternalServerError
	resp.Body = ioutil.NopCloser(strings.NewReader(`{"code":"ERROR_NAME","message":"ERROR_DESCRIPTION"}`))
	return nil
}

func mockDataWithInvalidRespData(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	resp.StatusCode = http.StatusInternalServerError
	resp.Body = ioutil.NopCloser(strings.NewReader(`{xxxxx}`))
	return nil
}

func mockDataWithInvalidHeader(req *http.Request, resp *http.Response, privateKey *rsa.PrivateKey) error {
	resp.Header = http.Header{}
	resp.Header.Set("Wechatpay-Nonce", mockNonce)
	resp.Header.Set("Wechatpay-Timestamp", "timestamp")
	resp.Header.Set("Wechatpay-Serial", mockSerialNo)
	resp.Body = ioutil.NopCloser(strings.NewReader(`{}`))
	return nil
}

func fromBase10(base10 string) *big.Int {
	i, ok := new(big.Int).SetString(base10, 10)
	if !ok {
		panic("bad number: " + base10)
	}
	return i
}
