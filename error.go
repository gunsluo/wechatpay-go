package wechatpay

import "strconv"

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
