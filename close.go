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
)

// CloseRequest is the request for close transaction
type CloseRequest struct {
	MchId      string `json:"mchid"`
	OutTradeNo string `json:"-"`
}

// Do send the request of close transaction
func (r *CloseRequest) Do(ctx context.Context, c Client) error {
	url := r.url(c.Config().Options().Domain)

	if err := c.Do(ctx, http.MethodPost, url, r).Error(); err != nil {
		return err
	}

	return nil
}
func (r *CloseRequest) url(domain string) string {
	return domain + "/v3/pay/transactions/out-trade-no/" + r.OutTradeNo + "/close"
}
