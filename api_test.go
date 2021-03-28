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

func TestClientPay(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req  *PayRequest
		resp *PayResponse
		pass bool
	}{
		{
			&PayRequest{
				Description: "for testing",
				OutTradeNo:  "forxxxxxxxxx",
				TimeExpire:  time.Now().Add(10 * time.Minute),
				Attach:      "cipher code",
				NotifyUrl:   "https://luoji.live/notify",
				Amount: PayAmount{
					Total:    1,
					Currency: "CNY",
				},
				TradeType: Native,
			},
			&PayResponse{
				CodeUrl: "weixin://wxpay/bizpayurl/up?pr=NwY5Mz9&groupid=00",
			},
			true,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		resp, err := client.Pay(ctx, c.req)
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

func TestClientQuery(t *testing.T) {
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
		req  *QueryRequest
		resp *QueryResponse
		pass bool
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
			true,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		resp, err := client.Query(ctx, c.req)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, err)
		}

		if err != nil {
			continue
		}

		if !c.resp.IsSuccess() {
			t.Fatal("invalid resp")
		}

		if !reflect.DeepEqual(c.resp, resp) {
			t.Fatalf("expect %v, got %v", c.resp, resp)
		}
	}
}

func TestClientCert(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req  *CertificatesRequest
		resp *CertificatesResponse
		pass bool
	}{
		{
			req: &CertificatesRequest{},
			resp: &CertificatesResponse{
				Certificates: []Certificate{
					{
						SerialNo:      mockSerialNo,
						EffectiveTime: "2020-09-17T14:26:23+08:00",
						ExpireTime:    "2025-09-16T14:26:23+08:00",
						Encrypt: EncryptCertificate{
							Algorithm:  "AEAD_AES_256_GCM",
							Nonce:      "eabb3e044577",
							Associated: "certificate",
							CipherText: "/M2eAJyVx/0y8JOErsNEWbYpikwKMS0hDahBYrR9Tnqvaxw/WLMHyLq7G3GUoWx3NSwYZlSZ+1JxAMTd4yge1B8bxY7OLrDkXm+BBDVypy5jCi/gcTQduTJpR4nRcBRYtEIxLGLrVaUXlDjDa4nM0mUPk6XA7AAUUAl3z5lYISapsFYUuHO9splBrmUESHxzRhSfsTyW68ll8o+ND7xA5R94slxzZIVdVg2Tz/3uXi5X1Qu5oi9Dn7pFdHD7++msMB7rgSJUTIFMwZ2GhAX3f/vVWemSMCymPPxzYxdiGFJJ8oBaIn+17pwulmz6NodFS0ilJr9wBs/05gqxe5L6S64ApwXNTfq3YJFVIU6munBaHomRZqsMg3MQlji9yNLBdKO2hk2rq/jCaBLsqcrCHEMEEULA5/1ImeYEkKcX2vIiVtKX8WxxP4M/Gq7btAQZVGzvczopb3wZNu1QLnzC13ov0pB5BPMhrx0tE4rLuZ5d+uzGOwuI8CvqOa+8TQ0DNGNaEA/IPrMJCVvmLrDi/aMQB+P4mO9BhUlfGHwQL7Q0anHzZaGHGkYyEGoTPmqQcY1mRbVcXDpIGn7rfHgiXnQTurB886T//ddhcv1/LQmcohSveZJAltcaDlmeqMgc+bXsOlAy6JNIIVPJ04ysI+V7nc0O4k4A32ZYA1hK52CU1YWz3vMoaaHVr/t6AF3dVWE1CphhNIwGbaz9M1sgEsWwT8LKLG5csgVwG20LO8wmLkxNUQ4fSkMdC+2Qv+rSFd8rlT1j+sYEbPVq6E6URkYPUKMqI1mEEudU1Rx0bE/pjj7+++0gX1H7sHp4+02KLdWS27gptHVXdDjNFPyCEshfVL2B8aEhq8PxSDG5zTqWHrKBAl04WU3kjlSsKZPrpKyhpIrKbEZHcrip3wOGeMf+4XDoZ8Iq8KoM8R6m8wkWi0GAW4G743O44PxHFvljKDIkIQm8gWV37jC3+qb/ZwUDxHONw3tHMH8XWsCVq1KAtKeE/iE9CCmE+ht7K4B+w0DeqKEicm0dkdjuFc9IgFa1W+q0HqGFI2Snd6ZX6crUy1I1vkRTQRj1mqjaP7dFOFV0JMpK/4CKMruZfUilNfOnSoKqHA2jPQ3f4ro0H22bF/PNhOWXp6Tzl5ZVbIFBIMdD9+ocq1lDH7vcBfKVwUltKl7jgI9HlpCDPZp++Mt3C4lPDzP/XrqorJnFBKw8eMBHS7N+jDhzhqJnI3ldwlGxUsqS/hj+jUUPpYINe/UtVwlOBi/tfuEfv47H5YgbP+Y3dz78a6KJUcA7caPSSqX+8LBcwEEZELXR8gU/AxwoDAsHM1pb7wc9fslct+awivfRi47AJtFeeZMGF6bb14VnbzvIZdpZRBIzHlvUqP+t8ZKEUvEJ+lVk7vv0/ySWBZbt0oA5XQ2RVwgzKGOgfMzZafsWAqrq1PGYjJqBbm/hudPtqsBridW/QjoE2Bp+Qnp8mWhdlSP8dgdeefLEeZGUSJx0Tzu2hBveEz7jMNQSOyg8HEE=",
						},
					},
				},
			},
			pass: true,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		resp, err := client.Cert(ctx, c.req)
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

func TestClientClose(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req  *CloseRequest
		pass bool
	}{
		{
			&CloseRequest{
				MchId:      client.config.MchId,
				OutTradeNo: "fortest",
			},
			true,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		err := client.Close(ctx, c.req)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, err)
		}

		if err != nil {
			continue
		}
	}
}

func TestClientRefund(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req  *RefundRequest
		resp *RefundResponse
		pass bool
	}{
		{
			&RefundRequest{
				TransactionId: "for test",
				OutTradeNo:    "for test",
				OutRefundNo:   "for test",
				Reason:        "for test",
				NotifyUrl:     "http://domain.com/notify",
				FundsAccount:  "",
				Amount: RefundAmount{
					Refund:   1,
					Total:    1,
					Currency: "CNY",
				},
				GoodsDetail: nil,
			},
			&RefundResponse{
				RefundId:            "50300807092021020105990201735",
				OutRefundNo:         "S20210201151309277501",
				TransactionId:       "4200000925202101284997714292",
				OutTradeNo:          "S20210128170702357723",
				Channel:             "ORIGINAL",
				UserReceivedAccount: "支付用户零钱",
				SuccessTime:         time.Time{},
				CreateTime:          dateFromString("2021-02-01T15:13:10+08:00"),
				Status:              "PROCESSING",
				FundsAccount:        "UNAVAILABLE",
				Amount: RefundAmountInQueryResp{
					Total:            1,
					Refund:           1,
					PayerTotal:       1,
					PayerRefund:      1,
					SettlementTotal:  1,
					SettlementRefund: 1,
					DiscountRefund:   0,
					Currency:         "CNY",
				},
				Promotion: nil,
			},
			true,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		resp, err := client.Refund(ctx, c.req)
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

func TestClientQueryRefund(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req  *RefundQueryRequest
		resp *RefundQueryResponse
		pass bool
	}{
		{
			req: &RefundQueryRequest{
				OutRefundNo: "1217752501201407033233368018",
			},
			resp: &RefundQueryResponse{
				RefundID:            "50000000382019052709732678859",
				OutRefundNo:         "1217752501201407033233368018",
				TransactionID:       "1217752501201407033233368018",
				OutTradeNo:          "1217752501201407033233368018",
				Channel:             "ORIGINAL",
				UserReceivedAccount: "招商银行信用卡0403",
				SuccessTime:         dateFromString("2020-12-01T16:18:12+08:00"),
				CreateTime:          dateFromString("2020-12-01T16:18:12+08:00"),
				Status:              "SUCCESS",
				FundsAccount:        "UNSETTLED",
				Amount: &RefundQueryAmount{
					Total:            100,
					Refund:           100,
					PayerTotal:       90,
					PayerRefund:      90,
					SettlementRefund: 100,
					SettlementTotal:  100,
					DiscountRefund:   10,
					Currency:         "CNY",
				},
				PromotionDetail: []RefundQueryPromotionDetail{
					{
						PromotionID:  "109519",
						Scope:        "SINGLE",
						Type:         "DISCOUNT",
						Amount:       5,
						RefundAmount: 100,
						GoodsDetail: []GoodsDetail{
							{
								MerchantGoodsID:  "1217752501201407033233368018",
								WechatpayGoodsID: "1001",
								GoodsName:        "iPhone6s 16G",
								UnitPrice:        528800,
								RefundAmount:     528800,
								RefundQuantity:   1,
							},
						},
					},
				},
			},
			pass: true,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		resp, err := client.QueryRefund(ctx, c.req)
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

func TestClientDownloadTradeBill(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req  *TradeBillRequest
		pass bool
		resp *TradeBillResponse
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
	}

	ctx := context.Background()
	for _, c := range cases {
		resp, err := client.DownloadTradeBill(ctx, c.req)
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

func TestClientDownloadOriginalTradeBill(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req    *TradeBillRequest
		pass   bool
		expect string
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
	}

	ctx := context.Background()
	for _, c := range cases {
		data, err := client.DownloadOriginalTradeBill(ctx, c.req)
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

func TestClientDownloadFundFlowBill(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req  *FundFlowBillRequest
		pass bool
		resp *FundFlowBillResponse
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
	}

	ctx := context.Background()
	for _, c := range cases {
		resp, err := client.DownloadFundFlowBill(ctx, c.req)
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

func TestClientDownloadFundOriginalFlowBill(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req    *FundFlowBillRequest
		pass   bool
		expect string
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
	}

	ctx := context.Background()
	for _, c := range cases {
		data, err := client.DownloadFundOriginalFlowBill(ctx, c.req)
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

func TestClientCombinePay(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req  *CombinePayRequest
		resp *CombinePayResponse
		pass bool
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
			true,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		resp, err := client.CombinePay(ctx, c.req)
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

func TestClientCombineQuery(t *testing.T) {
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
		req  *CombineQueryRequest
		resp *CombineQueryResponse
		pass bool
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
			true,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		resp, err := client.CombineQuery(ctx, c.req)
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

func TestClientCombineClose(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req  *CombineCloseRequest
		pass bool
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
			true,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
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
