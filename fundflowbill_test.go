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
	"strconv"
	"strings"
	"testing"

	"github.com/gunsluo/wechatpay-go/v3/sign"
)

func TestUnmarshalFundFlowBillSummary(t *testing.T) {
	cases := []struct {
		v      []string
		pass   bool
		expect *FundFlowBillSummary
	}{
		{
			[]string{"`3", "`1", "`0.01", "`2", "`0.02"},
			true,
			&FundFlowBillSummary{3, 1, 0.01, 2, 0.02},
		},
		{
			[]string{},
			false,
			&FundFlowBillSummary{},
		},
		{
			[]string{"`a3", "`1", "`0.01", "`2", "`0.02"},
			false,
			&FundFlowBillSummary{},
		},
		{
			[]string{"`3", "`a1", "`0.01", "`2", "`0.02"},
			false,
			&FundFlowBillSummary{},
		},
		{
			[]string{"`3", "`1", "`a0.01", "`2", "`0.02"},
			false,
			&FundFlowBillSummary{},
		},
		{
			[]string{"`3", "`1", "`0.01", "`a2", "`0.02"},
			false,
			&FundFlowBillSummary{},
		},
		{
			[]string{"`3", "`1", "`0.01", "`2", "`a0.02"},
			false,
			&FundFlowBillSummary{},
		},
	}

	for _, c := range cases {
		resp, err := UnmarshalFundFlowBillSummary(c.v)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, err)
		}

		if err != nil {
			continue
		}

		if !reflect.DeepEqual(c.expect, resp) {
			t.Fatalf("expect %v, got %v", c.expect, resp)
		}
	}
}

func TestUnmarshalFundFlowBill(t *testing.T) {
	cases := []struct {
		v      []string
		pass   bool
		expect *FundFlowBill
	}{
		{
			[]string{"`2021-02-01 13:54:01", "`50300806962021020105978994968", "`4200000920202101197964319284", "`退款", "`退款", "`支出", "`0.01", "`0.22", "`1601959334API", "`退款总金额0.01元;含手续费0.00元", "`S20210201135356381941"},
			true,
			&FundFlowBill{"2021-02-01 13:54:01", "50300806962021020105978994968", "4200000920202101197964319284", "退款", "退款", "支出", 0.01, 0.22, "1601959334API", "退款总金额0.01元;含手续费0.00元", "S20210201135356381941"},
		},
		{
			[]string{},
			false,
			&FundFlowBill{},
		},
		{
			[]string{"`2021-02-01 13:54:01", "`50300806962021020105978994968", "`4200000920202101197964319284", "`退款", "`退款", "`支出", "`a0.01", "`0.22", "`1601959334API", "`退款总金额0.01元;含手续费0.00元", "`S20210201135356381941"},
			false,
			&FundFlowBill{},
		},
		{
			[]string{"`2021-02-01 13:54:01", "`50300806962021020105978994968", "`4200000920202101197964319284", "`退款", "`退款", "`支出", "`0.01", "`a0.22", "`1601959334API", "`退款总金额0.01元;含手续费0.00元", "`S20210201135356381941"},
			false,
			&FundFlowBill{},
		},
	}

	for _, c := range cases {
		resp, err := UnmarshalFundFlowBill(c.v)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, err)
		}

		if err != nil {
			continue
		}

		if !reflect.DeepEqual(c.expect, resp) {
			t.Fatalf("expect %v, got %v", c.expect, resp)
		}
	}
}

func TestUnmarshalFundFlowBillResponse(t *testing.T) {
	cases := []struct {
		t      AccountType
		v      []byte
		pass   bool
		expect *FundFlowBillResponse
	}{
		{
			BasicAccount,
			[]byte("记账时间,微信支付业务单号,资金流水单号,业务名称,业务类型,收支类型,收支金额(元),账户结余(元),资金变更提交申请人,备注,业务凭证号\n" +
				"`2021-02-01 13:54:01,`50300806962021020105978994968,`4200000920202101197964319284,`退款,`退款,`支出,`0.01,`0.22,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201135356381941\n" +
				"`2021-02-01 14:00:45,`50300907032021020105978998710,`4200000846202101197461830397,`退款,`退款,`支出,`0.01,`0.21,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201140044552846\n" +
				"资金流水总笔数,收入笔数,收入金额,支出笔数,支出金额\n" +
				"`3,`1,`0.01,`2,`0.02\n"),
			true,
			&FundFlowBillResponse{
				Summary: FundFlowBillSummary{3, 1, 0.01, 2, 0.02},
				Bill: []*FundFlowBill{
					{"2021-02-01 13:54:01", "50300806962021020105978994968", "4200000920202101197964319284", "退款", "退款", "支出", 0.01, 0.22, "1601959334API", "退款总金额0.01元;含手续费0.00元", "S20210201135356381941"},
					{"2021-02-01 14:00:45", "50300907032021020105978998710", "4200000846202101197461830397", "退款", "退款", "支出", 0.01, 0.21, "1601959334API", "退款总金额0.01元;含手续费0.00元", "S20210201140044552846"},
				},
			},
		},
		{
			BasicAccount,
			[]byte{},
			false,
			&FundFlowBillResponse{},
		},
		{
			BasicAccount,
			[]byte("记账时间,微信支付业务单号,资金流水单号,业务名称,业务类型,收支类型,收支金额(元),账户结余(元),资金变更提交申请人,备注,业务凭证号\n" +
				"`2021-02-01 13:54:01,`50300806962021020105978994968,`4200000920202101197964319284,`退款,`退款,`支出,`0.01,`0.22,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201135356381941\n" +
				"`2021-02-01 14:00:45,`50300907032021020105978998710,`4200000846202101197461830397,`退款,`退款,`支出,`0.01,`0.21,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201140044552846\n" +
				"资金流水总笔数,收入笔数,收入金额,支出笔数,支出金额\n" +
				"`3,`1,`0.01,`2,`a0.02\n"),
			false,
			&FundFlowBillResponse{},
		},
		{
			BasicAccount,
			[]byte("记账时间,微信支付业务单号,资金流水单号,业务名称,业务类型,收支类型,收支金额(元),账户结余(元),资金变更提交申请人,备注,业务凭证号\n" +
				"`2021-02-01 13:54:01,`50300806962021020105978994968,`4200000920202101197964319284,`退款,`退款,`支出,`0.01,`0.22,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201135356381941\n" +
				"`2021-02-01 14:00:45,`50300907032021020105978998710,`4200000846202101197461830397,`退款,`退款,`支出,`a0.01,`0.21,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201140044552846\n" +
				"资金流水总笔数,收入笔数,收入金额,支出笔数,支出金额\n" +
				"`3,`1,`0.01,`2,`0.02\n"),
			false,
			&FundFlowBillResponse{},
		},
	}

	for _, c := range cases {
		resp, err := UnmarshalFundFlowBillResponse(c.t, c.v)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, err)
		}

		if err != nil {
			continue
		}

		if !reflect.DeepEqual(c.expect, resp) {
			t.Fatalf("expect %v, got %v", c.expect, resp)
		}
	}
}

func TestDownloadForFundFlowBill(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req       *FundFlowBillRequest
		transport *mockTransport
		pass      bool
		expect    string
	}{
		{
			req: &FundFlowBillRequest{
				BillDate:    "2021-01-01",
				AccountType: BasicAccount,
				TarType:     DataStream,
			},
			pass: true,
			expect: "记账时间,微信支付业务单号,资金流水单号,业务名称,业务类型,收支类型,收支金额(元),账户结余(元),资金变更提交申请人,备注,业务凭证号\n" +
				"`2021-02-01 13:54:01,`50300806962021020105978994968,`4200000920202101197964319284,`退款,`退款,`支出,`0.01,`0.22,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201135356381941\n" +
				"`2021-02-01 14:00:45,`50300907032021020105978998710,`4200000846202101197461830397,`退款,`退款,`支出,`0.01,`0.21,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201140044552846\n" +
				"资金流水总笔数,收入笔数,收入金额,支出笔数,支出金额\n" +
				"`3,`1,`0.01,`2,`0.02\n",
		},
		{
			req: &FundFlowBillRequest{
				BillDate:    "2021-01-01",
				AccountType: BasicAccount,
				TarType:     GZIP,
			},
			pass: true,
			expect: "记账时间,微信支付业务单号,资金流水单号,业务名称,业务类型,收支类型,收支金额(元),账户结余(元),资金变更提交申请人,备注,业务凭证号\n" +
				"`2021-02-01 13:54:01,`50300806962021020105978994968,`4200000920202101197964319284,`退款,`退款,`支出,`0.01,`0.22,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201135356381941\n" +
				"`2021-02-01 14:00:45,`50300907032021020105978998710,`4200000846202101197461830397,`退款,`退款,`支出,`0.01,`0.21,`1601959334API,`退款总金额0.01元;含手续费0.00元,`S20210201140044552846\n" +
				"资金流水总笔数,收入笔数,收入金额,支出笔数,支出金额\n" +
				"`3,`1,`0.01,`2,`0.02\n",
		},
		{
			req: &FundFlowBillRequest{
				BillDate:    "",
				AccountType: BasicAccount,
			},
			pass:   false,
			expect: "",
		},
		{
			req: &FundFlowBillRequest{
				BillDate:    "20210101",
				AccountType: BasicAccount,
			},
			pass:   false,
			expect: "",
		},
		{
			req: &FundFlowBillRequest{
				BillDate:    "2021-01-01",
				AccountType: BasicAccount,
				TarType:     "ABCD\nZ",
			},
			pass:   false,
			expect: "",
		},
		{
			req: &FundFlowBillRequest{
				BillDate:    "2021-01-01",
				AccountType: BasicAccount,
				TarType:     GZIP,
			},
			pass: false,
			transport: &mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					return mockDownloadFundflow(client.privateKey, req)
				},
			},
			expect: "",
		},
		{
			req: &FundFlowBillRequest{
				BillDate:    "2021-01-01",
				AccountType: BasicAccount,
				TarType:     GZIP,
			},
			pass: false,
			transport: &mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					return mockDownloadFundflow2(client.privateKey, req)
				},
			},
			expect: "",
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		if c.transport != nil {
			client.config.opts.transport = c.transport
			client.secrets.clear()
		}

		data, err := c.req.Download(ctx, client)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, err)
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

func TestUnmarshalDownloadForFundFlowBill(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req       *FundFlowBillRequest
		transport *mockTransport
		pass      bool
		resp      *FundFlowBillResponse
	}{
		{
			req: &FundFlowBillRequest{
				BillDate:    "2021-01-01",
				AccountType: BasicAccount,
				TarType:     DataStream,
			},
			pass: true,
			resp: &FundFlowBillResponse{
				Summary: FundFlowBillSummary{3, 1, 0.01, 2, 0.02},
				Bill: []*FundFlowBill{
					{"2021-02-01 13:54:01", "50300806962021020105978994968", "4200000920202101197964319284", "退款", "退款", "支出", 0.01, 0.22, "1601959334API", "退款总金额0.01元;含手续费0.00元", "S20210201135356381941"},
					{"2021-02-01 14:00:45", "50300907032021020105978998710", "4200000846202101197461830397", "退款", "退款", "支出", 0.01, 0.21, "1601959334API", "退款总金额0.01元;含手续费0.00元", "S20210201140044552846"},
				},
			},
		},
		{
			req: &FundFlowBillRequest{
				BillDate:    "2021-01-01",
				AccountType: BasicAccount,
				TarType:     DataStream,
			},
			transport: &mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					return mockDownloadFundflow(client.privateKey, req)
				},
			},
			pass: false,
		},
		{
			req: &FundFlowBillRequest{
				BillDate:    "2021-01-01",
				AccountType: BasicAccount,
				TarType:     DataStream,
			},
			transport: &mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					return mockDownloadFundflow2(client.privateKey, req)
				},
			},
			pass: false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		if c.transport != nil {
			client.config.opts.transport = c.transport
			client.secrets.clear()
		}

		resp, err := c.req.UnmarshalDownload(ctx, client)
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

func mockDownloadFundflow(privateKey *rsa.PrivateKey, req *http.Request) (*http.Response, error) {
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
	case "/v3/bill/fundflowbill":
		mockBody := `{"hash_type":"SHA1","hash_value":"dcd7ceb3d382a1181798368bb15d8437de46c00f","download_url":"https:\n//api.mch.weixin.qq.com/v3/billdownload/file?token=g44bIUH1GyQtE7ZmeTAPQx5b69qABpYuC_oZq6Aalf-gQP-lJ_FHRMLnyj2O8ujG"}`

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
			return nil, err
		}

		signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
		if err != nil {
			return nil, err
		}
		resp.Header.Set("Wechatpay-Nonce", mockNonce)
		resp.Header.Set("Wechatpay-Signature", signature)
		resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
		resp.Header.Set("Wechatpay-Serial", mockSerialNo)
		resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))
	}

	return resp, nil
}

func mockDownloadFundflow2(privateKey *rsa.PrivateKey, req *http.Request) (*http.Response, error) {
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
	case "/v3/bill/fundflowbill":
		mockBody := `{"hash_type":"SHA1","hash_value":"dcd7ceb3d382a1181798368bb15d8437de46c00f","download_url":"https://api.mch.weixin.qq.com/v3/billdownload/file?token=g44bIUH1GyQtE7ZmeTAPQx5b69qABpYuC_oZq6Aalf-gQP-lJ_FHRMLnyj2O8ujG"}`

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
			return nil, err
		}

		signature, err := sign.SignatureSHA256WithRSA(privateKey, plain)
		if err != nil {
			return nil, err
		}
		resp.Header.Set("Wechatpay-Nonce", mockNonce)
		resp.Header.Set("Wechatpay-Signature", signature)
		resp.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
		resp.Header.Set("Wechatpay-Serial", mockSerialNo)
		resp.Body = ioutil.NopCloser(strings.NewReader(mockBody))
	}

	return resp, nil
}
