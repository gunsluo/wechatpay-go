package wechatpay

import (
	"context"
	"reflect"
	"testing"
)

func TestRefundQueryRequest_Do(t *testing.T) {

	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req       *RefundQueryRequest
		resp      *RefundQueryResponse
		transport *mockTransport
		pass      bool
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
		{
			req: &RefundQueryRequest{},
			resp: nil,
			transport: nil,
			pass: false,
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
