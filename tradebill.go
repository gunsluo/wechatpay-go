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
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// TradeBillRequest is the request for trade bill
type TradeBillRequest struct {
	BillDate string   `json:"-"`
	BillType BillType `json:"-"`
	TarType  TarType  `json:"-"`
}

// TradeBillResponse is the response for trade bill
type TradeBillResponse struct {
	Summary TradeBillSummary
	Refund  []*RefundTradeBill
	All     []*AllTradeBill
	Success []*SuccessTradeBill
}

// Do send the request and get download url
func (r *TradeBillRequest) Do(ctx context.Context, c Client) (*FileUrl, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	url := r.url(c.Config().Options().Domain)

	fileUrl := &FileUrl{}
	if err := c.Do(ctx, http.MethodGet, url).Scan(fileUrl); err != nil {
		return nil, err
	}

	return fileUrl, nil
}

// Download download original the data of trade bill
func (r *TradeBillRequest) Download(ctx context.Context, c Client) ([]byte, error) {
	fileUrl, err := r.Do(ctx, c)
	if err != nil {
		return nil, err
	}

	data, err := c.Download(ctx, fileUrl)
	if err != nil {
		return nil, err
	}

	if r.TarType == GZIP {
		zr, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}

		var uncompressed bytes.Buffer
		if _, err := io.Copy(&uncompressed, zr); err != nil {
			return nil, err
		}

		if err := zr.Close(); err != nil {
			return nil, err
		}

		data = uncompressed.Bytes()
	}

	return data, nil
}

// UnmarshalDownload download and unmarshal the data of trade bill
func (r *TradeBillRequest) UnmarshalDownload(ctx context.Context, c Client) (*TradeBillResponse, error) {
	data, err := r.Download(ctx, c)
	if err != nil {
		return nil, err
	}

	resp, err := UnmarshalTradeBillResponse(r.BillType, data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *TradeBillRequest) validate() error {
	if r.BillDate == "" {
		return errors.New("bill date is required")
	}

	if _, err := time.Parse("2006-01-02", r.BillDate); err != nil {
		return fmt.Errorf("invalid bill date, the format: YYYY-MM-DD.")
	}

	return nil
}

func (r *TradeBillRequest) url(domain string) string {
	v := url.Values{}
	v.Add("bill_date", r.BillDate)
	if r.BillType != "" {
		v.Add("bill_type", string(r.BillType))
	}
	if r.TarType != "" {
		v.Add("tar_type", string(r.TarType))
	}

	return domain + "/v3/bill/tradebill?" + v.Encode()
}

// UnmarshalTradeBillResponse parses the bill data
// and stores the result in this response.
func UnmarshalTradeBillResponse(billType BillType, data []byte) (*TradeBillResponse, error) {
	if len(data) == 0 {
		return nil, errors.New("invaild data length")
	}

	r := &TradeBillResponse{}
	first := true
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for i := 0; scanner.Scan(); i++ {
		// skip title
		if i == 0 {
			continue
		}
		values := strings.Split(scanner.Text(), ",")

		// last line
		if len(values) == 7 {
			// skip title
			if first {
				first = false
				continue
			}
			summary, err := UnmarshalTradeBillSummary(values)
			if err != nil {
				return nil, err
			}
			r.Summary = *summary
			break
		}

		switch billType {
		case AllBill:
			b, err := UnmarshalAllTradeBill(values)
			if err != nil {
				return nil, err
			}
			r.All = append(r.All, b)
		case RefundBill:
			b, err := UnmarshalRefundTradeBill(values)
			if err != nil {
				return nil, err
			}
			r.Refund = append(r.Refund, b)
		case SuccessBill:
			b, err := UnmarshalSuccessTradeBill(values)
			if err != nil {
				return nil, err
			}
			r.Success = append(r.Success, b)
		default:
			b, err := UnmarshalAllTradeBill(values)
			if err != nil {
				return nil, err
			}
			r.All = append(r.All, b)
		}
	}
	return r, nil
}

// BillType is bill type
type BillType string

const (
	AllBill     BillType = "ALL"
	SuccessBill BillType = "SUCCESS"
	RefundBill  BillType = "REFUND"
)

// TarType is file tar type
type TarType string

const (
	DataStream TarType = ""
	GZIP       TarType = "GZIP"
)

// TradeBillSummary is summary trade bill
type TradeBillSummary struct {
	TotalNumberOfTransactions int
	TotalSettlementFee        float64
	TotalRefundFee            float64
	TotalCouponFee            float64
	TotalCommissionFee        float64
	TotalApplyRefundFee       float64
	TotalAmount               float64
}

// UnmarshalTradeBillSummary parses the bill data
// and stores the result in the bill summary.
func UnmarshalTradeBillSummary(values []string) (*TradeBillSummary, error) {
	if len(values) != 7 {
		return nil, errors.New("values length is invalid")
	}

	summary := &TradeBillSummary{}

	if i, err := atoi(values[0]); err != nil {
		return nil, err
	} else {
		summary.TotalNumberOfTransactions = i
	}

	if i, err := parseFloat(values[1]); err != nil {
		return nil, err
	} else {
		summary.TotalSettlementFee = i
	}

	if i, err := parseFloat(values[2]); err != nil {
		return nil, err
	} else {
		summary.TotalRefundFee = i
	}

	if i, err := parseFloat(values[3]); err != nil {
		return nil, err
	} else {
		summary.TotalCouponFee = i
	}

	if i, err := parseFloat(values[4]); err != nil {
		return nil, err
	} else {
		summary.TotalCommissionFee = i
	}

	if i, err := parseFloat(values[5]); err != nil {
		return nil, err
	} else {
		summary.TotalApplyRefundFee = i
	}

	if i, err := parseFloat(values[6]); err != nil {
		return nil, err
	} else {
		summary.TotalAmount = i
	}

	return summary, nil
}

// RefundTradeBill is data for refund trade bill
type RefundTradeBill struct {
	TradeTime          string
	AppId              string
	MchId              string
	SpecialMechId      string
	DeviceId           string
	TransactionId      string
	OutTradeNo         string
	OpenId             string
	TardeType          string
	TradeState         string
	BankType           string
	Currency           string
	SettlementTotalFee float64
	CouponAmount       float64
	RefundApplyTime    string
	RefundSuccessTime  string
	PayerRefundId      string
	MerchantRefundId   string
	RefundAmount       float64
	CouponRefundAmount float64
	RefundType         string
	RefundStatus       string
	GoodName           string
	Attach             string
	CommissionFee      float64
	Rate               string
	Amount             float64
	RefundApplyAmount  float64
	RateComment        string
}

// UnmarshalRefundTradeBill parses the bill data
// and stores the result in the bill .
func UnmarshalRefundTradeBill(values []string) (*RefundTradeBill, error) {
	if len(values) != 29 {
		return nil, errors.New("values length is invalid")
	}

	b := &RefundTradeBill{
		TradeTime:         removeDot(values[0]),
		AppId:             removeDot(values[1]),
		MchId:             removeDot(values[2]),
		SpecialMechId:     removeDot(values[3]),
		DeviceId:          removeDot(values[4]),
		TransactionId:     removeDot(values[5]),
		OutTradeNo:        removeDot(values[6]),
		OpenId:            removeDot(values[7]),
		TardeType:         removeDot(values[8]),
		TradeState:        removeDot(values[9]),
		BankType:          removeDot(values[10]),
		Currency:          removeDot(values[11]),
		RefundApplyTime:   removeDot(values[14]),
		RefundSuccessTime: removeDot(values[15]),
		PayerRefundId:     removeDot(values[16]),
		MerchantRefundId:  removeDot(values[17]),
		RefundType:        removeDot(values[20]),
		RefundStatus:      removeDot(values[21]),
		GoodName:          removeDot(values[22]),
		Attach:            removeDot(values[23]),
		Rate:              removeDot(values[25]),
		RateComment:       removeDot(values[28]),
	}

	if i, err := parseFloat(values[12]); err != nil {
		return nil, err
	} else {
		b.SettlementTotalFee = i
	}

	if i, err := parseFloat(values[13]); err != nil {
		return nil, err
	} else {
		b.CouponAmount = i
	}

	if i, err := parseFloat(values[18]); err != nil {
		return nil, err
	} else {
		b.RefundAmount = i
	}

	if i, err := parseFloat(values[19]); err != nil {
		return nil, err
	} else {
		b.CouponRefundAmount = i
	}

	if i, err := parseFloat(values[24]); err != nil {
		return nil, err
	} else {
		b.CommissionFee = i
	}

	if i, err := parseFloat(values[26]); err != nil {
		return nil, err
	} else {
		b.Amount = i
	}

	if i, err := parseFloat(values[27]); err != nil {
		return nil, err
	} else {
		b.RefundApplyAmount = i
	}

	return b, nil
}

// AllTradeBill is data for all trade bill
type AllTradeBill struct {
	TradeTime          string
	AppId              string
	MchId              string
	SpecialMechId      string
	DeviceId           string
	TransactionId      string
	OutTradeNo         string
	OpenId             string
	TardeType          string
	TradeState         string
	BankType           string
	Currency           string
	SettlementTotalFee float64
	CouponAmount       float64
	PayerRefundId      string
	MerchantRefundId   string
	RefundAmount       float64
	CouponRefundAmount float64
	RefundType         string
	RefundStatus       string
	GoodName           string
	Attach             string
	CommissionFee      float64
	Rate               string
	Amount             float64
	RefundApplyAmount  float64
	RateComment        string
}

// UnmarshalAllTradeBill parses the bill data
// and stores the result in the bill .
func UnmarshalAllTradeBill(values []string) (*AllTradeBill, error) {
	if len(values) != 27 {
		return nil, errors.New("values length is invalid")
	}

	b := &AllTradeBill{
		TradeTime:        removeDot(values[0]),
		AppId:            removeDot(values[1]),
		MchId:            removeDot(values[2]),
		SpecialMechId:    removeDot(values[3]),
		DeviceId:         removeDot(values[4]),
		TransactionId:    removeDot(values[5]),
		OutTradeNo:       removeDot(values[6]),
		OpenId:           removeDot(values[7]),
		TardeType:        removeDot(values[8]),
		TradeState:       removeDot(values[9]),
		BankType:         removeDot(values[10]),
		Currency:         removeDot(values[11]),
		PayerRefundId:    removeDot(values[14]),
		MerchantRefundId: removeDot(values[15]),
		RefundType:       removeDot(values[18]),
		RefundStatus:     removeDot(values[19]),
		GoodName:         removeDot(values[20]),
		Attach:           removeDot(values[21]),
		Rate:             removeDot(values[23]),
		RateComment:      removeDot(values[26]),
	}

	if i, err := parseFloat(values[12]); err != nil {
		return nil, err
	} else {
		b.SettlementTotalFee = i
	}

	if i, err := parseFloat(values[13]); err != nil {
		return nil, err
	} else {
		b.CouponAmount = i
	}

	if i, err := parseFloat(values[16]); err != nil {
		return nil, err
	} else {
		b.RefundAmount = i
	}

	if i, err := parseFloat(values[17]); err != nil {
		return nil, err
	} else {
		b.CouponRefundAmount = i
	}

	if i, err := parseFloat(values[22]); err != nil {
		return nil, err
	} else {
		b.CommissionFee = i
	}

	if i, err := parseFloat(values[24]); err != nil {
		return nil, err
	} else {
		b.Amount = i
	}

	if i, err := parseFloat(values[25]); err != nil {
		return nil, err
	} else {
		b.RefundApplyAmount = i
	}

	return b, nil
}

// SuccessTradeBill is data for success trade bill
type SuccessTradeBill struct {
	TradeTime          string
	AppId              string
	MchId              string
	SpecialMechId      string
	DeviceId           string
	TransactionId      string
	OutTradeNo         string
	OpenId             string
	TardeType          string
	TradeState         string
	BankType           string
	Currency           string
	SettlementTotalFee float64
	CouponAmount       float64
	GoodName           string
	Attach             string
	CommissionFee      float64
	Rate               string
	Amount             float64
	RateComment        string
}

// UnmarshalSuccessTradeBill parses the bill data
// and stores the result in the bill .
func UnmarshalSuccessTradeBill(values []string) (*SuccessTradeBill, error) {
	if len(values) != 20 {
		return nil, errors.New("values length is invalid")
	}

	b := &SuccessTradeBill{
		TradeTime:     removeDot(values[0]),
		AppId:         removeDot(values[1]),
		MchId:         removeDot(values[2]),
		SpecialMechId: removeDot(values[3]),
		DeviceId:      removeDot(values[4]),
		TransactionId: removeDot(values[5]),
		OutTradeNo:    removeDot(values[6]),
		OpenId:        removeDot(values[7]),
		TardeType:     removeDot(values[8]),
		TradeState:    removeDot(values[9]),
		BankType:      removeDot(values[10]),
		Currency:      removeDot(values[11]),
		GoodName:      removeDot(values[14]),
		Attach:        removeDot(values[15]),
		Rate:          removeDot(values[17]),
		RateComment:   removeDot(values[19]),
	}

	if i, err := parseFloat(values[12]); err != nil {
		return nil, err
	} else {
		b.SettlementTotalFee = i
	}

	if i, err := parseFloat(values[13]); err != nil {
		return nil, err
	} else {
		b.CouponAmount = i
	}

	if i, err := parseFloat(values[16]); err != nil {
		return nil, err
	} else {
		b.CommissionFee = i
	}

	if i, err := parseFloat(values[18]); err != nil {
		return nil, err
	} else {
		b.Amount = i
	}

	return b, nil
}

func removeDot(s string) string {
	if strings.HasPrefix(s, "`") {
		return s[1:]
	}

	return s
}

func atoi(s string) (int, error) {
	s = removeDot(s)
	return strconv.Atoi(s)
}

func parseFloat(s string) (float64, error) {
	s = removeDot(s)
	return strconv.ParseFloat(s, 64)
}
