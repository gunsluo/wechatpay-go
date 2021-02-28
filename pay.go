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

// Package wechatpay implements V3 endpoints for wechat pay. It is general
// SDK and provides the featrues, such as pay/query/close/notify transaction,
// refund/download bill.
//
// As a quick start:
//	client, err := NewClient(Config{})
//	// check error
//
// If you want to apply a pay request, use PayRequest
//	// create a pay request
//	req := &.PayRequest{
//		AppId:       appId,
//		MchId:       mchId,
//		Description: "for testing",
//		TradeType: .Native,
//	}
//
//	resp, err := req.Do(r.Context(), client)
//	if err != nil {
//		// do something
//	}
//	codeUrl := resp.CodeUrl
package wechatpay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// PayAmount is total amount paid, have total and currency.
type PayAmount struct {
	Total    int    `json:"total"`
	Currency string `json:"currency,omitempty"`
}

// PayDetail is the promotion information about the transaction.
type PayDetail struct {
	CostPrice   int          `json:"cost_price,omitempty"`
	InvoiceId   string       `json:"invoice_id,omitempty"`
	GoodsDetail []GoodDetail `json:"goods_detail,omitempty"`
}

// GoodDetail is the good information about the transaction.
type GoodDetail struct {
	MerchantGoodsId  string `json:"merchant_goods_id"`
	WechatpayGoodsId string `json:"wechatpay_goods_id,omitempty"`
	GoodsName        string `json:"goods_name,omitempty"`
	Quantity         int    `json:"quantity"`
	UnitPrice        int    `json:"unit_price"`
}

// PaySceneInfo is the scene information about the transaction.
type PaySceneInfo struct {
	PayerClientIp string     `json:"payer_client_ip"`
	DeviceId      string     `json:"device_id,omitempty"`
	StoreInfo     *StoreInfo `json:"store_info,omitempty"`
}

// StoreInfo  the store information about the transaction.
type StoreInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name,omitempty"`
	AreaCode string `json:"area_code,omitempty"`
	Address  string `json:"address,omitempty"`
}

// PayRequest is request when send a payment.
type PayRequest struct {
	AppId       string    `json:"appid"`
	MchId       string    `json:"mchid"`
	Description string    `json:"description"`
	OutTradeNo  string    `json:"out_trade_no"`
	TimeExpire  time.Time `json:"time_expire,omitempty"`
	Attach      string    `json:"attach,omitempty"`
	NotifyUrl   string    `json:"notify_url"`
	GoodsTag    string    `json:"goods_tag,omitempty"`
	Amount      PayAmount `json:"amount"`
	// Only set up Payer for JSAPI
	Payer     *Payer        `json:"payer,omitempty"`
	Detail    *PayDetail    `json:"detail,omitempty"`
	SceneInfo *PaySceneInfo `json:"scene_info,omitempty"`
	TradeType TradeType     `json:"-"`
}

// TradeType is trade type and defined by wechat pay.
type TradeType string

const (
	JSAPI  TradeType = "JSAPI"
	APP    TradeType = "APP"
	H5     TradeType = "H5"
	Native TradeType = "NATIVE"
)

// PayResponse is response when send a payment.
type PayResponse struct {
	// The CodeUrl is returned when the merchant used Native
	CodeUrl string `json:"code_url"`
	// The CodeUrl is returned when the merchant used JSAPI APP
	PrepayId string `json:"prepay_id"`
	// The CodeUrl is returned when the merchant used H5
	H5Url string `json:"h5_url"`
}

// Pay send a transaction and invoke wechat payment.
func (r *PayRequest) Do(ctx context.Context, c Client) (*PayResponse, error) {
	if r.AppId == "" {
		r.AppId = c.Config().AppId
	}

	if r.MchId == "" {
		r.MchId = c.Config().MchId
	}

	if r.TradeType == "" {
		r.TradeType = Native
	}

	switch r.TradeType {
	case JSAPI:
		if r.Payer == nil || r.Payer.OpenId == "" {
			return nil, errors.New("payer is required for JSAPI")
		}
	default:
		if r.Payer != nil {
			return nil, fmt.Errorf("don't set payer is for %v", r.TradeType)
		}
	}

	url := r.url(c.Config().Options().Domain)

	resp := &PayResponse{}
	if err := c.Do(ctx, http.MethodPost, url, r).Scan(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *PayRequest) url(domain string) string {
	return domain + "/v3/pay/transactions/" + strings.ToLower(string(r.TradeType))
}
