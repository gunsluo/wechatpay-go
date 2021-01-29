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
	"net/url"
)

// FundFlowBillRequest is the request for trade bill
type FundFlowBillRequest struct {
	BillDate    string `json:"-"`
	AccountType string `json:"-"`
	TarType     string `json:"-"`
}

// FundFlowBillRespone is the response for trade bill
type FundFlowBillRespone struct {
	HashType    string `json:"hash_type"`
	HashValue   string `json:"hash_value"`
	DownloadUrl string `json:"download_url"`
}

// Do send the request of close transaction
func (r *FundFlowBillRequest) Do(ctx context.Context, c Client) error {
	url := r.url(c.Config().Options().Domain)

	resp := &FundFlowBillRespone{}
	if err := c.Do(ctx, http.MethodGet, url).Scan(resp); err != nil {
		return err
	}

	return nil
}

func (r *FundFlowBillRequest) url(domain string) string {
	v := url.Values{}
	v.Add("bill_date", r.BillDate)
	if r.AccountType != "" {
		v.Add("account_type", r.AccountType)
	}
	if r.TarType != "" {
		v.Add("tar_type", r.TarType)
	}

	return domain + "/v3/bill/fundflowbill?" + v.Encode()
}
