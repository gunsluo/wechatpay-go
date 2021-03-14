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
	"time"
)

const (
	TradeStateSuccess    = "SUCCESS"
	TradeStateRefund     = "REFUND"
	TradeStateNotPay     = "NOTPAY"
	TradeStateClosed     = "CLOSED"
	TradeStateRevoked    = "REVOKED"
	TradeStateUserPaying = "USERPAYING"
	TradeStatePayError   = "PAYERROR"
	TradeStateAccept     = "ACCEPT"
)

// QueryRequest is the request for query transaction.
type QueryRequest struct {
	MchId         string `json:"-"`
	OutTradeNo    string `json:"-"`
	TransactionId string `json:"-"`
}

// QueryResponse is the response for query transaction.
type QueryResponse struct {
	AppId          string    `json:"appid"`
	MchId          string    `json:"mchid"`
	OutTradeNo     string    `json:"out_trade_no"`
	TransactionId  string    `json:"transaction_id,omitempty"`
	TradeType      TradeType `json:"trade_type,omitempty"`
	TradeState     string    `json:"trade_state"`
	TradeStateDesc string    `json:"trade_state_desc"`
	BankType       string    `json:"bank_type,omitempty"`
	Attach         string    `json:"attach,omitempty"`
	SuccessTime    time.Time `json:"success_time,omitempty"`
	Payer          Payer     `json:"payer"`

	Amount    TransactionAmount     `json:"amount,omitempty"`
	SceneInfo *TransactionSceneInfo `json:"scene_info,omitempty"`
	Promotion []*PromotionDetail    `json:"promotion_detail,omitempty"`
}

// IsSuccess check if the transactions pay success.
func (q QueryResponse) IsSuccess() bool {
	return q.TradeState == TradeStateSuccess
}

// Payer is the payer of the transaction.
type Payer struct {
	OpenId string `json:"openid"`
}

// TransactionAmount is tatal amount paid, have total and currency.
type TransactionAmount struct {
	Total         int    `json:"total,omitempty"`
	PayerTotal    int    `json:"payer_total,omitempty"`
	Currency      string `json:"currency,omitempty"`
	PayerCurrency string `json:"payer_currency,omitempty"`
}

// TransactionSceneInfo is the scene information about the transaction.
type TransactionSceneInfo struct {
	DeviceId string `json:"device_id,omitempty"`
}

// PromotionDetail is the promotion information about the transaction.
type PromotionDetail struct {
	CouponId            string `json:"coupon_id"`
	Name                string `json:"name,omitempty"`
	Scope               string `json:"scope,omitempty"`
	Type                string `json:"type,omitempty"`
	Amount              int    `json:"amount"`
	StockId             string `json:"stock_id,omitempty"`
	WechatpayContribute int    `json:"wechatpay_contribute,omitempty"`
	MerchantContribute  int    `json:"merchant_contribute,omitempty"`
	OtherContribute     int    `json:"other_contribute,omitempty"`
	Currency            string `json:"currency,omitempty"`

	GoodsDetail []TransactionGoodDetail `json:"goods_detail,omitempty"`
}

// TransactionGoodDetail is the good information about the transaction.
type TransactionGoodDetail struct {
	GoodsId        string `json:"goods_id"`
	Quantity       int    `json:"quantity"`
	UnitPrice      int    `json:"unit_price"`
	DiscountAmount int    `json:"discount_amount"`
	GoodsRemark    string `json:"goods_remark,omitempty"`
}

// Do send the request of query transaction.
func (r *QueryRequest) Do(ctx context.Context, c Client) (*QueryResponse, error) {
	if r.MchId == "" {
		r.MchId = c.Config().MchId
	}

	url := r.url(c.Config().Options().Domain)

	resp := &QueryResponse{}
	if err := c.Do(ctx, http.MethodGet, url).Scan(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// return the url according to querying parameters.
func (r *QueryRequest) url(domain string) string {
	if r.TransactionId != "" {
		return domain + "/v3/pay/transactions/id/" + r.TransactionId + "?mchid=" + r.MchId
	}

	return domain + "/v3/pay/transactions/out-trade-no/" + r.OutTradeNo + "?mchid=" + r.MchId
}
