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
	"strconv"
)

// Error is more detail error
type Error struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implement Error function for err
func (e *Error) Error() string {
	if e == nil {
		return "{}"
	}

	return `{"status":` + strconv.Itoa(e.Status) + `,"code":"` + e.Code + `","message":"` + e.Message + `"}`
}

const (
	UserPaying           = "USERPAYING"
	TradeError           = "TRADE_ERROR"
	SystemError          = "SYSTEMERROR"
	SignError            = "SIGN_ERROR"
	RuleLimit            = "RULELIMIT"
	ParamError           = "PARAM_ERROR"
	OutTradeNoUsed       = "OUT_TRADE_NO_USED"
	OrderNotExist        = "ORDERNOTEXIST"
	OrderClosed          = "ORDER_CLOSED"
	OpenidMismatch       = "OPENID_MISMATCH"
	NotEnough            = "NOTENOUGH"
	NoAuth               = "NOAUTH"
	MchNotExists         = "MCH_NOT_EXISTS"
	InvalidTransactionid = "INVALID_TRANSACTIONID"
	InvalidRequest       = "INVALID_REQUEST"
	FrequencyLimited     = "FREQUENCY_LIMITED"
	BankError            = "BANKERROR"
	AppidMchidNotMatch   = "APPID_MCHID_NOT_MATCH"
	AccountError         = "ACCOUNTERROR"
)
