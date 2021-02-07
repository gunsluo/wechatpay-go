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

func TestDownloadForTradeBill(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req       *TradeBillRequest
		transport *mockTransport
		pass      bool
		expect    string
	}{
		{
			req: &TradeBillRequest{
				BillDate: "2021-01-01",
				BillType: AllBill,
				TarType:  DataStream,
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
			req: &TradeBillRequest{
				BillDate: "2021-01-01",
				BillType: AllBill,
				TarType:  GZIP,
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
			req: &TradeBillRequest{
				BillDate: "",
				BillType: AllBill,
			},
			pass:   false,
			expect: "",
		},
		{
			req: &TradeBillRequest{
				BillDate: "20210101",
				BillType: AllBill,
			},
			pass:   false,
			expect: "",
		},
		{
			req: &TradeBillRequest{
				BillDate: "2021-01-01",
				BillType: AllBill,
				TarType:  "ABCD\nZ",
			},
			pass:   false,
			expect: "",
		},
		{
			req: &TradeBillRequest{
				BillDate: "2021-01-01",
				BillType: AllBill,
				TarType:  GZIP,
			},
			pass: false,
			transport: &mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					return mockDownload(client.privateKey, req)
				},
			},
			expect: "",
		},
		{
			req: &TradeBillRequest{
				BillDate: "2021-01-01",
				BillType: AllBill,
				TarType:  GZIP,
			},
			pass: false,
			transport: &mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					return mockDownload2(client.privateKey, req)
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

func TestUnmarshalDownloadForTradeBill(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req       *TradeBillRequest
		transport *mockTransport
		pass      bool
		resp      *TradeBillResponse
	}{
		{
			req: &TradeBillRequest{
				BillDate: "2021-01-01",
				BillType: AllBill,
				TarType:  DataStream,
			},
			pass: true,
			resp: &TradeBillResponse{
				Summary: TradeBillSummary{3, 0.03, 0.00, 0.00, 0.00000, 0.03, 0.00},
				All: []*AllTradeBill{
					{"2021-01-28 17:07:11", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000925202101284997714292", "S20210128170702357723", "ofyak5qR_1wYsC99CsWA6R9MJazA", "NATIVE", "SUCCESS", "OTHERS", "CNY", 0.01, 0.00, "0", "0", 0.00, 0.00, "", "", "for testing", "cipher code", 0.00000, "1.00%", 0.01, 0.00, ""},
					{`2021-01-28 15:35:18`, `wx81be3101902f7cb2`, `1601959334`, "0", "", `4200000910202101282955148400`, `S20210128153505214586`, `ofyak5qR_1wYsC99CsWA6R9MJazA`, `NATIVE`, `SUCCESS`, `OTHERS`, `CNY`, 0.01, 0.00, "0", "0", 0.00, 0.00, ``, ``, `for testing`, `cipher code`, 0.00000, `1.00%`, 0.01, 0.00, ``},
					{`2021-01-28 16:59:46`, `wx81be3101902f7cb2`, `1601959334`, `0`, ``, `4200000926202101281412639609`, `S20210128165824499930`, `ofyak5qR_1wYsC99CsWA6R9MJazA`, `NATIVE`, `SUCCESS`, `OTHERS`, `CNY`, 0.01, 0.00, `0`, `0`, 0.00, 0.00, ``, ``, `for testing`, `cipher code`, 0.00000, `1.00%`, 0.01, 0.00, ""},
				},
			},
		},
		{
			req: &TradeBillRequest{
				BillDate: "2021-01-01",
				BillType: AllBill,
				TarType:  DataStream,
			},
			transport: &mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					return mockDownload(client.privateKey, req)
				},
			},
			pass: false,
		},
		{
			req: &TradeBillRequest{
				BillDate: "2021-01-01",
				BillType: AllBill,
				TarType:  DataStream,
			},
			transport: &mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					return mockDownload2(client.privateKey, req)
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

func TestUnmarshalTradeBillSummary(t *testing.T) {
	cases := []struct {
		v      []string
		pass   bool
		expect *TradeBillSummary
	}{
		{
			[]string{"`3", "`0.03", "`0.00", "`0.00", "`0.00000", "`0.03", "`0.00"},
			true,
			&TradeBillSummary{3, 0.03, 0.00, 0.00, 0.00000, 0.03, 0.00},
		},
		{
			[]string{},
			false,
			&TradeBillSummary{},
		},
		{
			[]string{"`a3", "`0.03", "`0.00", "`0.00", "`0.00000", "`0.03", "`0.00"},
			false,
			&TradeBillSummary{},
		},
		{
			[]string{"`3", "`a0.03", "`0.00", "`0.00", "`0.00000", "`0.03", "`0.00"},
			false,
			&TradeBillSummary{},
		},
		{
			[]string{"`3", "`0.03", "`a0.00", "`0.00", "`0.00000", "`0.03", "`0.00"},
			false,
			&TradeBillSummary{},
		},
		{
			[]string{"`3", "`0.03", "`0.00", "`a0.00", "`0.00000", "`0.03", "`0.00"},
			false,
			&TradeBillSummary{},
		},
		{
			[]string{"`3", "`0.03", "`0.00", "`0.00", "`a0.00000", "`0.03", "`0.00"},
			false,
			&TradeBillSummary{},
		},
		{
			[]string{"`3", "`0.03", "`0.00", "`0.00", "`0.00000", "`a0.03", "`0.00"},
			false,
			&TradeBillSummary{},
		},
		{
			[]string{"`3", "`0.03", "`0.00", "`0.00", "`0.00000", "`0.03", "`a0.00"},
			false,
			&TradeBillSummary{},
		},
	}

	for _, c := range cases {
		resp, err := UnmarshalTradeBillSummary(c.v)
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

func TestUnmarshalAllTradeBill(t *testing.T) {
	cases := []struct {
		v      []string
		pass   bool
		expect *AllTradeBill
	}{
		{
			[]string{"`2021-01-28 17:07:11", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000925202101284997714292", "`S20210128170702357723", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`0.01", "`0.00", "`0", "`0", "`0.00", "`0.00", "`", "`", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.01", "`0.00", "`"},
			true,
			&AllTradeBill{"2021-01-28 17:07:11", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000925202101284997714292", "S20210128170702357723", "ofyak5qR_1wYsC99CsWA6R9MJazA", "NATIVE", "SUCCESS", "OTHERS", "CNY", 0.01, 0.00, "0", "0", 0.00, 0.00, "", "", "for testing", "cipher code", 0.00000, "1.00%", 0.01, 0.00, ""},
		},
		{
			[]string{},
			false,
			&AllTradeBill{},
		},
		{
			[]string{"`2021-01-28 17:07:11", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000925202101284997714292", "`S20210128170702357723", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`a0.01", "`0.00", "`0", "`0", "`0.00", "`0.00", "`", "`", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.01", "`0.00", "`"},
			false,
			&AllTradeBill{},
		},
		{
			[]string{"`2021-01-28 17:07:11", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000925202101284997714292", "`S20210128170702357723", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`0.01", "`a0.00", "`0", "`0", "`0.00", "`0.00", "`", "`", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.01", "`0.00", "`"},
			false,
			&AllTradeBill{},
		},
		{
			[]string{"`2021-01-28 17:07:11", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000925202101284997714292", "`S20210128170702357723", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`0.01", "`0.00", "`0", "`0", "`a0.00", "`0.00", "`", "`", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.01", "`0.00", "`"},
			false,
			&AllTradeBill{},
		},
		{
			[]string{"`2021-01-28 17:07:11", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000925202101284997714292", "`S20210128170702357723", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`0.01", "`0.00", "`0", "`0", "`0.00", "`a0.00", "`", "`", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.01", "`0.00", "`"},
			false,
			&AllTradeBill{},
		},
		{
			[]string{"`2021-01-28 17:07:11", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000925202101284997714292", "`S20210128170702357723", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`0.01", "`0.00", "`0", "`0", "`0.00", "`0.00", "`", "`", "`for testing", "`cipher code", "`a0.00000", "`1.00%", "`0.01", "`0.00", "`"},
			false,
			&AllTradeBill{},
		},
		{
			[]string{"`2021-01-28 17:07:11", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000925202101284997714292", "`S20210128170702357723", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`0.01", "`0.00", "`0", "`0", "`0.00", "`0.00", "`", "`", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`a0.01", "`0.00", "`"},
			false,
			&AllTradeBill{},
		},
		{
			[]string{"`2021-01-28 17:07:11", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000925202101284997714292", "`S20210128170702357723", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`0.01", "`0.00", "`0", "`0", "`0.00", "`0.00", "`", "`", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.01", "`a0.00", "`"},
			false,
			&AllTradeBill{},
		},
	}

	for _, c := range cases {
		resp, err := UnmarshalAllTradeBill(c.v)
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

func TestUnmarshalRefundTradeBill(t *testing.T) {
	cases := []struct {
		v      []string
		pass   bool
		expect *RefundTradeBill
	}{
		{
			[]string{"`2021-01-24 16:16:25", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000844202101245866928772", "`S20210124161554311546", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`REFUND", "`OTHERS", "`CNY", "`0.00", "`0.00", "`2021-02-01 14:33:21", "`2021-02-01 14:33:24", "`50300807172021020106006664916", "`S20210201143320649393", "`0.01", "`0.00", "`ORIGINAL", "`SUCCESS", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.00", "`0.01", "`"},
			true,
			&RefundTradeBill{"2021-01-24 16:16:25", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000844202101245866928772", "S20210124161554311546", "ofyak5qR_1wYsC99CsWA6R9MJazA", "NATIVE", "REFUND", "OTHERS", "CNY", 0.00, 0.00, "2021-02-01 14:33:21", "2021-02-01 14:33:24", "50300807172021020106006664916", "S20210201143320649393", 0.01, 0.00, "ORIGINAL", "SUCCESS", "for testing", "cipher code", 0.00000, "1.00%", 0.00, 0.01, ""},
		},
		{
			[]string{},
			false,
			&RefundTradeBill{},
		},
		{
			[]string{"`2021-01-24 16:16:25", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000844202101245866928772", "`S20210124161554311546", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`REFUND", "`OTHERS", "`CNY", "`a0.00", "`0.00", "`2021-02-01 14:33:21", "`2021-02-01 14:33:24", "`50300807172021020106006664916", "`S20210201143320649393", "`0.01", "`0.00", "`ORIGINAL", "`SUCCESS", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.00", "`0.01", "`"},
			false,
			&RefundTradeBill{},
		},
		{
			[]string{"`2021-01-24 16:16:25", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000844202101245866928772", "`S20210124161554311546", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`REFUND", "`OTHERS", "`CNY", "`0.00", "`a0.00", "`2021-02-01 14:33:21", "`2021-02-01 14:33:24", "`50300807172021020106006664916", "`S20210201143320649393", "`0.01", "`0.00", "`ORIGINAL", "`SUCCESS", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.00", "`0.01", "`"},
			false,
			&RefundTradeBill{},
		},
		{
			[]string{"`2021-01-24 16:16:25", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000844202101245866928772", "`S20210124161554311546", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`REFUND", "`OTHERS", "`CNY", "`0.00", "`0.00", "`2021-02-01 14:33:21", "`2021-02-01 14:33:24", "`50300807172021020106006664916", "`S20210201143320649393", "`a0.01", "`0.00", "`ORIGINAL", "`SUCCESS", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.00", "`0.01", "`"},
			false,
			&RefundTradeBill{},
		},
		{
			[]string{"`2021-01-24 16:16:25", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000844202101245866928772", "`S20210124161554311546", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`REFUND", "`OTHERS", "`CNY", "`0.00", "`0.00", "`2021-02-01 14:33:21", "`2021-02-01 14:33:24", "`50300807172021020106006664916", "`S20210201143320649393", "`0.01", "`a0.00", "`ORIGINAL", "`SUCCESS", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.00", "`0.01", "`"},
			false,
			&RefundTradeBill{},
		},
		{
			[]string{"`2021-01-24 16:16:25", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000844202101245866928772", "`S20210124161554311546", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`REFUND", "`OTHERS", "`CNY", "`0.00", "`0.00", "`2021-02-01 14:33:21", "`2021-02-01 14:33:24", "`50300807172021020106006664916", "`S20210201143320649393", "`0.01", "`0.00", "`ORIGINAL", "`SUCCESS", "`for testing", "`cipher code", "`a0.00000", "`1.00%", "`0.00", "`0.01", "`"},
			false,
			&RefundTradeBill{},
		},
		{
			[]string{"`2021-01-24 16:16:25", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000844202101245866928772", "`S20210124161554311546", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`REFUND", "`OTHERS", "`CNY", "`0.00", "`0.00", "`2021-02-01 14:33:21", "`2021-02-01 14:33:24", "`50300807172021020106006664916", "`S20210201143320649393", "`0.01", "`0.00", "`ORIGINAL", "`SUCCESS", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`a0.00", "`0.01", "`"},
			false,
			&RefundTradeBill{},
		},
		{
			[]string{"`2021-01-24 16:16:25", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000844202101245866928772", "`S20210124161554311546", "`ofyak5qR_1wYsC99CsWA6R9MJazA", "`NATIVE", "`REFUND", "`OTHERS", "`CNY", "`0.00", "`0.00", "`2021-02-01 14:33:21", "`2021-02-01 14:33:24", "`50300807172021020106006664916", "`S20210201143320649393", "`0.01", "`0.00", "`ORIGINAL", "`SUCCESS", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.00", "`a0.01", "`"},
			false,
			&RefundTradeBill{},
		},
	}

	for _, c := range cases {
		resp, err := UnmarshalRefundTradeBill(c.v)
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

func TestUnmarshalSuccessTradeBill(t *testing.T) {
	cases := []struct {
		v      []string
		pass   bool
		expect *SuccessTradeBill
	}{
		{
			[]string{"`2021-02-01 14:38:45", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000922202102014836880592", "`S20210201143829466741", "`ofyak5lCyFIsihOYEX0Zx9smR0g0", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`0.01", "`0.00", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.01", "`"},
			true,
			&SuccessTradeBill{"2021-02-01 14:38:45", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000922202102014836880592", "S20210201143829466741", "ofyak5lCyFIsihOYEX0Zx9smR0g0", "NATIVE", "SUCCESS", "OTHERS", "CNY", 0.01, 0.00, "for testing", "cipher code", 0.00000, "1.00%", 0.01, ""},
		},
		{
			[]string{},
			false,
			&SuccessTradeBill{},
		},
		{
			[]string{"`2021-02-01 14:38:45", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000922202102014836880592", "`S20210201143829466741", "`ofyak5lCyFIsihOYEX0Zx9smR0g0", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`a0.01", "`0.00", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.01", "`"},
			false,
			&SuccessTradeBill{},
		},
		{
			[]string{"`2021-02-01 14:38:45", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000922202102014836880592", "`S20210201143829466741", "`ofyak5lCyFIsihOYEX0Zx9smR0g0", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`0.01", "`a0.00", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`0.01", "`"},
			false,
			&SuccessTradeBill{},
		},
		{
			[]string{"`2021-02-01 14:38:45", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000922202102014836880592", "`S20210201143829466741", "`ofyak5lCyFIsihOYEX0Zx9smR0g0", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`0.01", "`0.00", "`for testing", "`cipher code", "`a0.00000", "`1.00%", "`0.01", "`"},
			false,
			&SuccessTradeBill{},
		},
		{
			[]string{"`2021-02-01 14:38:45", "`wx81be3101902f7cb2", "`1601959334", "`0", "`", "`4200000922202102014836880592", "`S20210201143829466741", "`ofyak5lCyFIsihOYEX0Zx9smR0g0", "`NATIVE", "`SUCCESS", "`OTHERS", "`CNY", "`0.01", "`0.00", "`for testing", "`cipher code", "`0.00000", "`1.00%", "`a0.01", "`"},
			false,
			&SuccessTradeBill{},
		},
	}

	for _, c := range cases {
		resp, err := UnmarshalSuccessTradeBill(c.v)
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

func TestUnmarshalTradeBillResponse(t *testing.T) {
	cases := []struct {
		t      BillType
		v      []byte
		pass   bool
		expect *TradeBillResponse
	}{
		{
			BillType(""),
			[]byte("交易时间,公众账号ID,商户号,特约商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,应结订单金额,代金券金额,微信退款单号,商户退款单号,退款金额,充值券退款金额,退款类型,退款状态,商品名称,商户数据包,手续费,费率,订单金额,申请退款金额,费率备注\n" +
				"`2021-01-28 17:07:11,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000925202101284997714292,`S20210128170702357723,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 15:35:18,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000910202101282955148400,`S20210128153505214586,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 16:59:46,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000926202101281412639609,`S20210128165824499930,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"总交易单数,应结订单总金额,退款总金额,充值券退款总金额,手续费总金额,订单总金额,申请退款总金额\n" +
				"`3,`0.03,`0.00,`0.00,`0.00000,`0.03,`0.00\n"),
			true,
			&TradeBillResponse{
				Summary: TradeBillSummary{3, 0.03, 0.00, 0.00, 0.00000, 0.03, 0.00},
				All: []*AllTradeBill{
					{"2021-01-28 17:07:11", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000925202101284997714292", "S20210128170702357723", "ofyak5qR_1wYsC99CsWA6R9MJazA", "NATIVE", "SUCCESS", "OTHERS", "CNY", 0.01, 0.00, "0", "0", 0.00, 0.00, "", "", "for testing", "cipher code", 0.00000, "1.00%", 0.01, 0.00, ""},
					{`2021-01-28 15:35:18`, `wx81be3101902f7cb2`, `1601959334`, "0", "", `4200000910202101282955148400`, `S20210128153505214586`, `ofyak5qR_1wYsC99CsWA6R9MJazA`, `NATIVE`, `SUCCESS`, `OTHERS`, `CNY`, 0.01, 0.00, "0", "0", 0.00, 0.00, ``, ``, `for testing`, `cipher code`, 0.00000, `1.00%`, 0.01, 0.00, ``},
					{`2021-01-28 16:59:46`, `wx81be3101902f7cb2`, `1601959334`, `0`, ``, `4200000926202101281412639609`, `S20210128165824499930`, `ofyak5qR_1wYsC99CsWA6R9MJazA`, `NATIVE`, `SUCCESS`, `OTHERS`, `CNY`, 0.01, 0.00, `0`, `0`, 0.00, 0.00, ``, ``, `for testing`, `cipher code`, 0.00000, `1.00%`, 0.01, 0.00, ""},
				},
			},
		},
		{
			AllBill,
			[]byte("交易时间,公众账号ID,商户号,特约商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,应结订单金额,代金券金额,微信退款单号,商户退款单号,退款金额,充值券退款金额,退款类型,退款状态,商品名称,商户数据包,手续费,费率,订单金额,申请退款金额,费率备注\n" +
				"`2021-01-28 17:07:11,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000925202101284997714292,`S20210128170702357723,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 15:35:18,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000910202101282955148400,`S20210128153505214586,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 16:59:46,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000926202101281412639609,`S20210128165824499930,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"总交易单数,应结订单总金额,退款总金额,充值券退款总金额,手续费总金额,订单总金额,申请退款总金额\n" +
				"`3,`0.03,`0.00,`0.00,`0.00000,`0.03,`0.00\n"),
			true,
			&TradeBillResponse{
				Summary: TradeBillSummary{3, 0.03, 0.00, 0.00, 0.00000, 0.03, 0.00},
				All: []*AllTradeBill{
					{"2021-01-28 17:07:11", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000925202101284997714292", "S20210128170702357723", "ofyak5qR_1wYsC99CsWA6R9MJazA", "NATIVE", "SUCCESS", "OTHERS", "CNY", 0.01, 0.00, "0", "0", 0.00, 0.00, "", "", "for testing", "cipher code", 0.00000, "1.00%", 0.01, 0.00, ""},
					{`2021-01-28 15:35:18`, `wx81be3101902f7cb2`, `1601959334`, "0", "", `4200000910202101282955148400`, `S20210128153505214586`, `ofyak5qR_1wYsC99CsWA6R9MJazA`, `NATIVE`, `SUCCESS`, `OTHERS`, `CNY`, 0.01, 0.00, "0", "0", 0.00, 0.00, ``, ``, `for testing`, `cipher code`, 0.00000, `1.00%`, 0.01, 0.00, ``},
					{`2021-01-28 16:59:46`, `wx81be3101902f7cb2`, `1601959334`, `0`, ``, `4200000926202101281412639609`, `S20210128165824499930`, `ofyak5qR_1wYsC99CsWA6R9MJazA`, `NATIVE`, `SUCCESS`, `OTHERS`, `CNY`, 0.01, 0.00, `0`, `0`, 0.00, 0.00, ``, ``, `for testing`, `cipher code`, 0.00000, `1.00%`, 0.01, 0.00, ""},
				},
			},
		},
		{
			AllBill,
			[]byte{},
			false,
			&TradeBillResponse{},
		},
		{
			BillType(""),
			[]byte("交易时间,公众账号ID,商户号,特约商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,应结订单金额,代金券金额,微信退款单号,商户退款单号,退款金额,充值券退款金额,退款类型,退款状态,商品名称,商户数据包,手续费,费率,订单金额,申请退款金额,费率备注\n" +
				"`2021-01-28 17:07:11,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000925202101284997714292,`S20210128170702357723,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 15:35:18,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000910202101282955148400,`S20210128153505214586,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 16:59:46,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000926202101281412639609,`S20210128165824499930,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"总交易单数,应结订单总金额,退款总金额,充值券退款总金额,手续费总金额,订单总金额,申请退款总金额\n" +
				"`3,`0.03,`0.00,`0.00,`0.00000,`0.03,`a0.00\n"),
			false,
			&TradeBillResponse{},
		},
		{
			AllBill,
			[]byte("交易时间,公众账号ID,商户号,特约商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,应结订单金额,代金券金额,微信退款单号,商户退款单号,退款金额,充值券退款金额,退款类型,退款状态,商品名称,商户数据包,手续费,费率,订单金额,申请退款金额,费率备注\n" +
				"`2021-01-28 17:07:11,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000925202101284997714292,`S20210128170702357723,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 15:35:18,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000910202101282955148400,`S20210128153505214586,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 16:59:46,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000926202101281412639609,`S20210128165824499930,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"总交易单数,应结订单总金额,退款总金额,充值券退款总金额,手续费总金额,订单总金额,申请退款总金额\n" +
				"`a3,`0.03,`0.00,`0.00,`0.00000,`0.03,`0.00\n"),
			false,
			&TradeBillResponse{},
		},
		{
			AllBill,
			[]byte("交易时间,公众账号ID,商户号,特约商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,应结订单金额,代金券金额,微信退款单号,商户退款单号,退款金额,充值券退款金额,退款类型,退款状态,商品名称,商户数据包,手续费,费率,订单金额,申请退款金额,费率备注\n" +
				"`2021-01-28 17:07:11,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000925202101284997714292,`S20210128170702357723,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 15:35:18,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000910202101282955148400,`S20210128153505214586,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`0.00,`\n" +
				"`2021-01-28 16:59:46,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000926202101281412639609,`S20210128165824499930,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`0,`0,`0.00,`0.00,`,`,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`a0.00,`\n" +
				"总交易单数,应结订单总金额,退款总金额,充值券退款总金额,手续费总金额,订单总金额,申请退款总金额\n" +
				"`3,`0.03,`0.00,`0.00,`0.00000,`0.03,`0.00\n"),
			false,
			&TradeBillResponse{},
		},
		{
			RefundBill,
			[]byte("交易时间,公众账号ID,商户号,特约商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,应结订单金额,代金券金额,退款申请时间,退款成功时间,微信退款单号,商户退款单号,退款金额,充值券退款金额,退款类型,退款状态,商品名称,商户数据包,手续费,费率,订单金额,申请退款金额,费率备注\n" +
				"`2021-01-24 16:16:25,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000844202101245866928772,`S20210124161554311546,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`REFUND,`OTHERS,`CNY,`0.00,`0.00,`2021-02-01 14:33:21,`2021-02-01 14:33:24,`50300807172021020106006664916,`S20210201143320649393,`0.01,`0.00,`ORIGINAL,`SUCCESS,`for testing,`cipher code,`0.00000,`1.00%,`0.00,`0.01,`\n" +
				"`2021-01-19 16:31:18,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000846202101197461830397,`S20210119083100844726118382,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`REFUND,`OTHERS,`CNY,`0.00,`0.00,`2021-02-01 14:00:45,`2021-02-01 14:00:50,`50300907032021020105978998710,`S20210201140044552846,`0.01,`0.00,`ORIGINAL,`SUCCESS,`Package Venue,`,`0.00000,`1.00%,`0.00,`0.01,`\n" +
				"总交易单数,应结订单总金额,退款总金额,充值券退款总金额,手续费总金额,订单总金额,申请退款总金额\n" +
				"`2,`0.00,`0.02,`0.00,`0.00000,`0.00,`0.02\n"),
			true,
			&TradeBillResponse{
				Summary: TradeBillSummary{2, 0.00, 0.02, 0.00, 0.00000, 0.00, 0.02},
				Refund: []*RefundTradeBill{
					{"2021-01-24 16:16:25", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000844202101245866928772", "S20210124161554311546", "ofyak5qR_1wYsC99CsWA6R9MJazA", "NATIVE", "REFUND", "OTHERS", "CNY", 0.00, 0.00, "2021-02-01 14:33:21", "2021-02-01 14:33:24", "50300807172021020106006664916", "S20210201143320649393", 0.01, 0.00, "ORIGINAL", "SUCCESS", "for testing", "cipher code", 0.00000, "1.00%", 0.00, 0.01, ""},
					{"2021-01-19 16:31:18", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000846202101197461830397", "S20210119083100844726118382", "ofyak5qR_1wYsC99CsWA6R9MJazA", "NATIVE", "REFUND", "OTHERS", "CNY", 0.00, 0.00, "2021-02-01 14:00:45", "2021-02-01 14:00:50", "50300907032021020105978998710", "S20210201140044552846", 0.01, 0.00, "ORIGINAL", "SUCCESS", "Package Venue", "", 0.00000, "1.00%", 0.00, 0.01, ""},
				},
			},
		},
		{
			RefundBill,
			[]byte("交易时间,公众账号ID,商户号,特约商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,应结订单金额,代金券金额,退款申请时间,退款成功时间,微信退款单号,商户退款单号,退款金额,充值券退款金额,退款类型,退款状态,商品名称,商户数据包,手续费,费率,订单金额,申请退款金额,费率备注\n" +
				"`2021-01-24 16:16:25,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000844202101245866928772,`S20210124161554311546,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`REFUND,`OTHERS,`CNY,`0.00,`0.00,`2021-02-01 14:33:21,`2021-02-01 14:33:24,`50300807172021020106006664916,`S20210201143320649393,`0.01,`0.00,`ORIGINAL,`SUCCESS,`for testing,`cipher code,`0.00000,`1.00%,`0.00,`0.01,`\n" +
				"`2021-01-19 16:31:18,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000846202101197461830397,`S20210119083100844726118382,`ofyak5qR_1wYsC99CsWA6R9MJazA,`NATIVE,`REFUND,`OTHERS,`CNY,`0.00,`0.00,`2021-02-01 14:00:45,`2021-02-01 14:00:50,`50300907032021020105978998710,`S20210201140044552846,`0.01,`0.00,`ORIGINAL,`SUCCESS,`Package Venue,`,`0.00000,`1.00%,`0.00,`a0.01,`\n" +
				"总交易单数,应结订单总金额,退款总金额,充值券退款总金额,手续费总金额,订单总金额,申请退款总金额\n" +
				"`2,`0.00,`0.02,`0.00,`0.00000,`0.00,`0.02\n"),
			false,
			&TradeBillResponse{
				Summary: TradeBillSummary{2, 0.00, 0.02, 0.00, 0.00000, 0.00, 0.02},
				Refund: []*RefundTradeBill{
					{"2021-01-24 16:16:25", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000844202101245866928772", "S20210124161554311546", "ofyak5qR_1wYsC99CsWA6R9MJazA", "NATIVE", "REFUND", "OTHERS", "CNY", 0.00, 0.00, "2021-02-01 14:33:21", "2021-02-01 14:33:24", "50300807172021020106006664916", "S20210201143320649393", 0.01, 0.00, "ORIGINAL", "SUCCESS", "for testing", "cipher code", 0.00000, "1.00%", 0.00, 0.01, ""},
					{"2021-01-19 16:31:18", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000846202101197461830397", "S20210119083100844726118382", "ofyak5qR_1wYsC99CsWA6R9MJazA", "NATIVE", "REFUND", "OTHERS", "CNY", 0.00, 0.00, "2021-02-01 14:00:45", "2021-02-01 14:00:50", "50300907032021020105978998710", "S20210201140044552846", 0.01, 0.00, "ORIGINAL", "SUCCESS", "Package Venue", "", 0.00000, "1.00%", 0.00, 0.01, ""},
				},
			},
		},
		{
			SuccessBill,
			[]byte("交易时间,公众账号ID,商户号,特约商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,应结订单金额,代金券金额,商品名称,商户数据包,手续费,费率,订单金额,费率备注\n" +
				"`2021-02-01 14:38:45,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000922202102014836880592,`S20210201143829466741,`ofyak5lCyFIsihOYEX0Zx9smR0g0,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`for testing,`cipher code,`0.00000,`1.00%,`0.01,`\n" +
				"总交易单数,应结订单总金额,退款总金额,充值券退款总金额,手续费总金额,订单总金额,申请退款总金额\n" +
				"`1,`0.01,`0.00,`0.00,`0.00000,`0.01,`0.00\n"),
			true,
			&TradeBillResponse{
				Summary: TradeBillSummary{1, 0.01, 0.00, 0.00, 0.00000, 0.01, 0.00},
				Success: []*SuccessTradeBill{
					{"2021-02-01 14:38:45", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000922202102014836880592", "S20210201143829466741", "ofyak5lCyFIsihOYEX0Zx9smR0g0", "NATIVE", "SUCCESS", "OTHERS", "CNY", 0.01, 0.00, "for testing", "cipher code", 0.00000, "1.00%", 0.01, ""},
				},
			},
		},
		{
			SuccessBill,
			[]byte("交易时间,公众账号ID,商户号,特约商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,应结订单金额,代金券金额,商品名称,商户数据包,手续费,费率,订单金额,费率备注\n" +
				"`2021-02-01 14:38:45,`wx81be3101902f7cb2,`1601959334,`0,`,`4200000922202102014836880592,`S20210201143829466741,`ofyak5lCyFIsihOYEX0Zx9smR0g0,`NATIVE,`SUCCESS,`OTHERS,`CNY,`0.01,`0.00,`for testing,`cipher code,`0.00000,`1.00%,`a0.01,`\n" +
				"总交易单数,应结订单总金额,退款总金额,充值券退款总金额,手续费总金额,订单总金额,申请退款总金额\n" +
				"`1,`0.01,`0.00,`0.00,`0.00000,`0.01,0.00\n"),
			false,
			&TradeBillResponse{
				Summary: TradeBillSummary{1, 0.01, 0.00, 0.00, 0.00000, 0.01, 0.00},
				Success: []*SuccessTradeBill{
					{"2021-02-01 14:38:45", "wx81be3101902f7cb2", "1601959334", "0", "", "4200000922202102014836880592", "S20210201143829466741", "ofyak5lCyFIsihOYEX0Zx9smR0g0", "NATIVE", "SUCCESS", "OTHERS", "CNY", 0.01, 0.00, "for testing", "cipher code", 0.00000, "1.00%", 0.01, ""},
				},
			},
		},
	}

	for _, c := range cases {
		resp, err := UnmarshalTradeBillResponse(c.t, c.v)
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

func mockDownload(privateKey *rsa.PrivateKey, req *http.Request) (*http.Response, error) {
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
	case "/v3/bill/tradebill":
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

func mockDownload2(privateKey *rsa.PrivateKey, req *http.Request) (*http.Response, error) {
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
	case "/v3/bill/tradebill":
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

const mockRSAPrivateKeyCert = `-----BEGIN PRIVATE KEY-----
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
