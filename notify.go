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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// PayNotification is a paying notification from wechatpay.
type PayNotification struct {
	Notification
}

// PayNotifyTransaction is the transaction after being decrypted.
type PayNotifyTransaction = QueryResponse

// NotificationAnswer is sent to wechat pay after
// processing the notification.
type NotificationAnswer struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// String return a json string.
func (a *NotificationAnswer) String() string {
	return `{"code":"` + a.Code + `","message":"` + a.Message + `"}`
}

// Bytes return a json array bytes.
func (a *NotificationAnswer) Bytes() []byte {
	return []byte(a.String())
}

// ParseHttpRequest pasre the data that read from the http request.
// return a transaction.
func (n *PayNotification) ParseHttpRequest(c Client, req *http.Request) (*PayNotifyTransaction, error) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	nonce := req.Header.Get("Wechatpay-Nonce")
	signature := req.Header.Get("Wechatpay-Signature")
	ts := req.Header.Get("Wechatpay-Timestamp")
	serialNo := req.Header.Get("Wechatpay-Serial")

	var timestamp int64
	if ts != "" {
		i, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return nil, err
		}
		timestamp = i
	}

	result := &Result{
		Body:      data,
		Timestamp: timestamp,
		Nonce:     nonce,
		Signature: signature,
		SerialNo:  serialNo,
	}

	return n.Parse(req.Context(), c, result)
}

// Parse pasre the data from result and return a transaction.
func (n *PayNotification) Parse(ctx context.Context, c Client, result *Result) (*PayNotifyTransaction, error) {
	on, data, err := c.ParseNotification(ctx, result)
	if err != nil {
		return nil, err
	}

	n.Notification = *on

	var trans PayNotifyTransaction
	if err := json.Unmarshal(data, &trans); err != nil {
		return nil, err
	}

	return &trans, nil
}

// RefundNotification is a refund notification from wechatpay.
type RefundNotification struct {
	Notification
}

// RefundNotifyTransaction is the transaction after being decrypted.
type RefundNotifyTransaction struct {
	MchId               string    `json:"mchid"`
	OutTradeNo          string    `json:"out_trade_no"`
	TransactionId       string    `json:"transaction_id"`
	OutRefundNo         string    `json:"out_refund_no"`
	RefundId            string    `json:"refund_id"`
	RefundStatus        string    `json:"refund_status"`
	SuccessTime         time.Time `json:"success_time,omitempty"`
	UserReceivedAccount string    `json:"user_received_account"`

	Amount RefundAmountInNotify `json:"amount"`
}

// RefundAmountInNotify is total amount refund.
type RefundAmountInNotify struct {
	Total       int `json:"total"`
	Refund      int `json:"refund"`
	PayerTotal  int `json:"payer_total"`
	PayerRefund int `json:"payer_refund"`
}

// ParseHttpRequest pasre the data that read from the http request.
// return a refund transaction.
func (n *RefundNotification) ParseHttpRequest(c Client, req *http.Request) (*RefundNotifyTransaction, error) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	nonce := req.Header.Get("Wechatpay-Nonce")
	signature := req.Header.Get("Wechatpay-Signature")
	ts := req.Header.Get("Wechatpay-Timestamp")
	serialNo := req.Header.Get("Wechatpay-Serial")

	var timestamp int64
	if ts != "" {
		i, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return nil, err
		}
		timestamp = i
	}

	result := &Result{
		Body:      data,
		Timestamp: timestamp,
		Nonce:     nonce,
		Signature: signature,
		SerialNo:  serialNo,
	}

	return n.Parse(req.Context(), c, result)
}

// Parse pasre the data from result and return a refund transcation.
func (n *RefundNotification) Parse(ctx context.Context, c Client, result *Result) (*RefundNotifyTransaction, error) {
	on, data, err := c.ParseNotification(ctx, result)
	if err != nil {
		return nil, err
	}
	n.Notification = *on

	var trans RefundNotifyTransaction
	if err := json.Unmarshal(data, &trans); err != nil {
		return nil, err
	}

	return &trans, nil
}
