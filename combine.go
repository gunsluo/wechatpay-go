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
