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
	"errors"
	"net/http"
	"strings"
	"time"
)

// CombinePayAmount is total amount paid, have total and currency.
type CombinePayAmount struct {
	Total    int    `json:"total_amount"`
	Currency string `json:"currency,omitempty"`
}

// SettleInfo is settle information
type SettleInfo struct {
	ProfitSharing bool `json:"profit_sharing"`
	SubsidyAmount bool `json:"subsidy_amount"`
}

// SubOrder is the order under the combine transcation
type SubOrder struct {
	MchId       string           `json:"mchid"`
	Attach      string           `json:"attach,omitempty"`
	Amount      CombinePayAmount `json:"amount"`
	OutTradeNo  string           `json:"out_trade_no"`
	Description string           `json:"description"`
}

// CombinePayRequest is request when send a combin payment.
type CombinePayRequest struct {
	AppId      string        `json:"combine_appid"`
	MchId      string        `json:"combine_mchid"`
	OutTradeNo string        `json:"combine_out_trade_no"`
	TimeStart  time.Time     `json:"time_start,omitempty"`
	TimeExpire time.Time     `json:"time_expire,omitempty"`
	NotifyUrl  string        `json:"notify_url"`
	SceneInfo  *PaySceneInfo `json:"scene_info,omitempty"`
	Payer      *Payer        `json:"combine_payer_info,omitempty"`
	Orders     []SubOrder    `json:"sub_orders,omitempty"`
	TradeType  TradeType     `json:"-"`
}

// CombinePayResponse is response when send a combine payment.
type CombinePayResponse struct {
	// The CodeUrl is returned when the merchant used Native
	CodeUrl string `json:"code_url"`
	// The CodeUrl is returned when the merchant used JSAPI APP
	PrepayId string `json:"prepay_id"`
	// The CodeUrl is returned when the merchant used H5
	H5Url string `json:"h5_url"`
}

// Do send a transaction and invoke wechat payment.
func (r *CombinePayRequest) Do(ctx context.Context, c Client) (*CombinePayResponse, error) {
	if r.AppId == "" {
		r.AppId = c.Config().AppId
	}

	if r.MchId == "" {
		r.MchId = c.Config().MchId
	}

	if r.TradeType == "" {
		r.TradeType = Native
	}

	if len(r.Orders) == 0 {
		return nil, errors.New("orders is required")
	}

	switch r.TradeType {
	case JSAPI:
		if r.Payer == nil || r.Payer.OpenId == "" {
			return nil, errors.New("payer is required for JSAPI")
		}
	}

	url := r.url(c.Config().Options().Domain)

	resp := &CombinePayResponse{}
	if err := c.Do(ctx, http.MethodPost, url, r).Scan(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *CombinePayRequest) url(domain string) string {
	return domain + "/v3/combine-transactions/" + strings.ToLower(string(r.TradeType))
}

// CloseSubOrder is the order under the combine close transcation
type CloseSubOrder struct {
	MchId      string `json:"mchid"`
	OutTradeNo string `json:"out_trade_no"`
}

// CombineCloseRequest is the request for close transaction.
type CombineCloseRequest struct {
	AppId      string          `json:"combine_appid"`
	OutTradeNo string          `json:"combine_out_trade_no"`
	Orders     []CloseSubOrder `json:"sub_orders,omitempty"`
}

// Do send the request of combine close transaction.
func (r *CombineCloseRequest) Do(ctx context.Context, c Client) error {
	if r.AppId == "" {
		r.AppId = c.Config().AppId
	}

	if len(r.Orders) == 0 {
		return errors.New("orders is required")
	}

	url := r.url(c.Config().Options().Domain)

	if err := c.Do(ctx, http.MethodPost, url, r).Error(); err != nil {
		return err
	}

	return nil
}

// return the url for combine close transcation
func (r *CombineCloseRequest) url(domain string) string {
	return domain + "/v3/combine-transactions/out-trade-no/" + r.OutTradeNo + "/close"
}

// CombineQueryRequest is the request for query transaction.
type CombineQueryRequest struct {
	OutTradeNo string `json:"combine_out_trade_no"`
}

// QuerySubOrder is the order under the combine transcation
type QuerySubOrder struct {
	MchId         string    `json:"mchid"`
	OutTradeNo    string    `json:"out_trade_no"`
	TradeType     TradeType `json:"trade_type,omitempty"`
	TradeState    string    `json:"trade_state"`
	BankType      string    `json:"bank_type,omitempty"`
	Attach        string    `json:"attach,omitempty"`
	SuccessTime   time.Time `json:"success_time,omitempty"`
	TransactionId string    `json:"transaction_id,omitempty"`

	Amount CombineSubOrderAmount `json:"amount,omitempty"`
}

// CombineSubOrderAmount is tatal amount paid, have total and currency.
type CombineSubOrderAmount struct {
	Total         int    `json:"total_amount,omitempty"`
	PayerTotal    int    `json:"payer_total,omitempty"`
	Currency      string `json:"currency,omitempty"`
	PayerCurrency string `json:"payer_currency,omitempty"`
}

// CombineQueryResponse is the response for query transaction.
type CombineQueryResponse struct {
	AppId      string                `json:"combine_appid"`
	MchId      string                `json:"combine_mchid"`
	OutTradeNo string                `json:"combine_out_trade_no"`
	SceneInfo  *TransactionSceneInfo `json:"scene_info,omitempty"`
	Orders     []QuerySubOrder       `json:"sub_orders,omitempty"`
	Payer      *Payer                `json:"combine_payer_info,omitempty"`
}

// Do send the request of query transaction.
func (r *CombineQueryRequest) Do(ctx context.Context, c Client) (*CombineQueryResponse, error) {
	if r.OutTradeNo == "" {
		return nil, errors.New("out trader no is required")
	}

	url := r.url(c.Config().Options().Domain)

	resp := &CombineQueryResponse{}
	if err := c.Do(ctx, http.MethodGet, url).Scan(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// return the url according to querying parameters.
func (r *CombineQueryRequest) url(domain string) string {
	return domain + "/v3/combine-transactions/out-trade-no/" + r.OutTradeNo
}
