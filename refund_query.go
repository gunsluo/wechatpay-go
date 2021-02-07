package wechatpay

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// RefundQueryResponse is the result for refund query
type RefundQueryResponse struct {
	RefundID            string                       `json:"refund_id"`
	OutRefundNo         string                       `json:"out_refund_no"`
	TransactionID       string                       `json:"transaction_id"`
	OutTradeNo          string                       `json:"out_trade_no"`
	Channel             string                       `json:"channel"`
	UserReceivedAccount string                       `json:"user_received_account"`
	SuccessTime         time.Time                    `json:"success_time"`
	CreateTime          time.Time                    `json:"create_time"`
	Status              string                       `json:"status"`
	FundsAccount        string                       `json:"funds_account"`
	Amount              *RefundQueryAmount           `json:"amount"`
	PromotionDetail     []RefundQueryPromotionDetail `json:"promotion_detail"`
}

type RefundQueryAmount struct {
	Total            int    `json:"total"`
	Refund           int    `json:"refund"`
	PayerTotal       int    `json:"payer_total"`
	PayerRefund      int    `json:"payer_refund"`
	SettlementRefund int    `json:"settlement_refund"`
	SettlementTotal  int    `json:"settlement_total"`
	DiscountRefund   int    `json:"discount_refund"`
	Currency         string `json:"currency"`
}

type GoodsDetail struct {
	MerchantGoodsID  string `json:"merchant_goods_id"`
	WechatpayGoodsID string `json:"wechatpay_goods_id"`
	GoodsName        string `json:"goods_name"`
	UnitPrice        int    `json:"unit_price"`
	RefundAmount     int    `json:"refund_amount"`
	RefundQuantity   int    `json:"refund_quantity"`
}

type RefundQueryPromotionDetail struct {
	PromotionID  string        `json:"promotion_id"`
	Scope        string        `json:"scope"`
	Type         string        `json:"type"`
	Amount       int           `json:"amount"`
	RefundAmount int           `json:"refund_amount"`
	GoodsDetail  []GoodsDetail `json:"goods_detail"`
}

// RefundQueryRequest is the request for query transaction
type RefundQueryRequest struct {
	OutRefundNo string `json:"-"`
}

// Do send the refund query result
func (r *RefundQueryRequest) Do(ctx context.Context, c Client) (*RefundQueryResponse, error) {
	url := r.url(c.Config().Options().Domain)

	if err := r.validate(); err != nil {
		return nil, err
	}

	resp := &RefundQueryResponse{}
	if err := c.Do(ctx, http.MethodGet, url, nil).Scan(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *RefundQueryRequest) validate() error {
	if r.OutRefundNo == "" {
		return errors.New("out_refund_no can't be empty")
	}

	return nil
}

func (r *RefundQueryRequest) url(domain string) string {
	return domain + `/v3/refund/domestic/refunds/` + r.OutRefundNo
}
