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

	"github.com/gunsluo/wechatpay-go/v3/sign"
)

// PayNotification is a paying notification from wechatpay
type PayNotification struct {
	Id           string `json:"id"`
	CreateTime   string `json:"create_time"`
	EventType    string `json:"event_type"`
	ResourceType string `json:"resource_type"`
	Summary      string `json:"summary"`

	Resource PayNotifyResource `json:"resource"`
}

// EncryptCertificate is the information of encrypt certificate
type PayNotifyResource struct {
	Algorithm    string `json:"algorithm"`
	CipherText   string `json:"ciphertext"`
	Associated   string `json:"associated_data"`
	OriginalType string `json:"original_type"`
	Nonce        string `json:"nonce"`
}

// PayNotifyTransaction is the transaction after being decrypted
type PayNotifyTransaction = QueryResponse

// PayNotificationAnswer is sent to wechat pay after
// processing the notification.
type PayNotificationAnswer struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// String return a json string
func (a *PayNotificationAnswer) String() string {
	return `{"code":"` + a.Code + `","message":"` + a.Message + `"}`
}

// Bytes return a json array bytes
func (a *PayNotificationAnswer) Bytes() []byte {
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

// ParseHttpRequest pasre the data from result.
func (n *PayNotification) Parse(ctx context.Context, c Client, result *Result) (*PayNotifyTransaction, error) {
	if err := json.Unmarshal(result.Body, n); err != nil {
		return nil, err
	}

	// verify signature
	if err := c.VerifySignature(ctx, result); err != nil {
		return nil, err
	}

	// using apiv3 secret decrypt data
	apiv3Secret := []byte(c.Config().Apiv3Secret)
	data, err := sign.DecryptByAes256Gcm(
		apiv3Secret,
		[]byte(n.Resource.Nonce),
		[]byte(n.Resource.Associated),
		n.Resource.CipherText)
	if err != nil {
		return nil, err
	}

	var trans PayNotifyTransaction
	if err := json.Unmarshal(data, &trans); err != nil {
		return nil, err
	}

	return &trans, nil
}
