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
	"strings"
	"time"
)

// FundFlowBillRequest is the request for trade bill.
type FundFlowBillRequest struct {
	BillDate    string      `json:"-"`
	AccountType AccountType `json:"-"`
	TarType     TarType     `json:"-"`
}

// FundFlowBillResponse is the response for trade bill.
type FundFlowBillResponse struct {
	Summary FundFlowBillSummary
	Bill    []*FundFlowBill
}

// FundFlowBill is summary fundflow.
type FundFlowBillSummary struct {
	TotalNumber          int
	TotalNumberOfIncome  int
	IncomeAomunt         float64
	TotalNumberOfOutcome int
	OutcomeAomunt        float64
}

// FundFlowBill is data for fund flow.
type FundFlowBill struct {
	AccountingTime      string
	TransactionId       string
	OrderNo             string
	BusinessName        string
	BusinessType        string
	InOutcomeType       string
	InOutcomeAmount     float64
	AccountBalance      float64
	FundChangeApplicant string
	Remark              string
	BusinessNumber      string
}

// Do send the request of downloading fundflow bill.
func (r *FundFlowBillRequest) Do(ctx context.Context, c Client) (*FileUrl, error) {
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

// Download download original the data of fundflow bill.
func (r *FundFlowBillRequest) Download(ctx context.Context, c Client) ([]byte, error) {
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

// UnmarshalDownload download and unmarshal the data of fundflow bill.
func (r *FundFlowBillRequest) UnmarshalDownload(ctx context.Context, c Client) (*FundFlowBillResponse, error) {
	data, err := r.Download(ctx, c)
	if err != nil {
		return nil, err
	}

	resp, err := UnmarshalFundFlowBillResponse(r.AccountType, data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *FundFlowBillRequest) validate() error {
	if r.BillDate == "" {
		return errors.New("bill date is required")
	}

	if _, err := time.Parse("2006-01-02", r.BillDate); err != nil {
		return fmt.Errorf("invalid bill date, the format: YYYY-MM-DD.")
	}

	return nil
}

func (r *FundFlowBillRequest) url(domain string) string {
	v := url.Values{}
	v.Add("bill_date", r.BillDate)
	if r.AccountType != "" {
		v.Add("account_type", string(r.AccountType))
	}
	if r.TarType != "" {
		v.Add("tar_type", string(r.TarType))
	}

	return domain + "/v3/bill/fundflowbill?" + v.Encode()
}

// UnmarshalFundFlowBillResponse parses the bill data
// and stores the result in this response.
func UnmarshalFundFlowBillResponse(accountType AccountType, data []byte) (*FundFlowBillResponse, error) {
	if len(data) == 0 {
		return nil, errors.New("invaild data length")
	}

	r := &FundFlowBillResponse{}
	first := true
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for i := 0; scanner.Scan(); i++ {
		// skip title
		if i == 0 {
			continue
		}
		values := strings.Split(scanner.Text(), ",")

		// last line
		if len(values) == 5 {
			// skip title
			if first {
				first = false
				continue
			}
			summary, err := UnmarshalFundFlowBillSummary(values)
			if err != nil {
				return nil, err
			}
			r.Summary = *summary
			break
		}

		b, err := UnmarshalFundFlowBill(values)
		if err != nil {
			return nil, err
		}
		r.Bill = append(r.Bill, b)
	}

	return r, nil
}

// UnmarshalFundFlowBillSummary parses the bill data
// and stores the result in the bill summary.
func UnmarshalFundFlowBillSummary(values []string) (*FundFlowBillSummary, error) {
	if len(values) != 5 {
		return nil, errors.New("values length is invalid")
	}

	summary := &FundFlowBillSummary{}
	if i, err := atoi(values[0]); err != nil {
		return nil, err
	} else {
		summary.TotalNumber = i
	}

	if i, err := atoi(values[1]); err != nil {
		return nil, err
	} else {
		summary.TotalNumberOfIncome = i
	}

	if i, err := parseFloat(values[2]); err != nil {
		return nil, err
	} else {
		summary.IncomeAomunt = i
	}

	if i, err := atoi(values[3]); err != nil {
		return nil, err
	} else {
		summary.TotalNumberOfOutcome = i
	}

	if i, err := parseFloat(values[4]); err != nil {
		return nil, err
	} else {
		summary.OutcomeAomunt = i
	}

	return summary, nil
}

// UnmarshalFundFlowBill parses the bill data
// and stores the result in the bill.
func UnmarshalFundFlowBill(values []string) (*FundFlowBill, error) {
	if len(values) != 11 {
		return nil, errors.New("values length is invalid")
	}

	b := &FundFlowBill{
		AccountingTime:      removeDot(values[0]),
		TransactionId:       removeDot(values[1]),
		OrderNo:             removeDot(values[2]),
		BusinessName:        removeDot(values[3]),
		BusinessType:        removeDot(values[4]),
		InOutcomeType:       removeDot(values[5]),
		FundChangeApplicant: removeDot(values[8]),
		Remark:              removeDot(values[9]),
		BusinessNumber:      removeDot(values[10]),
	}

	if i, err := parseFloat(values[6]); err != nil {
		return nil, err
	} else {
		b.InOutcomeAmount = i
	}

	if i, err := parseFloat(values[7]); err != nil {
		return nil, err
	} else {
		b.AccountBalance = i
	}

	return b, nil
}

// AccountType is account type.
type AccountType string

const (
	BasicAccount     AccountType = "BASIC"
	OperationAccount AccountType = "OPERATION"
	FEESAccount      AccountType = "FEES"
)
