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

import "context"

// API is wechat pay api v3.
type API interface {
	Pay(ctx context.Context, r *PayRequest) (*PayResponse, error)
	Query(ctx context.Context, r *QueryRequest) (*QueryResponse, error)
	Cert(ctx context.Context, r *CertificatesRequest) (*CertificatesResponse, error)
	Close(ctx context.Context, r *CloseRequest) error
	Refund(ctx context.Context, r *RefundRequest) (*RefundResponse, error)
	QueryRefund(ctx context.Context, r *RefundQueryRequest) (*RefundQueryResponse, error)
	DownloadTradeBill(ctx context.Context, r *TradeBillRequest) (*TradeBillResponse, error)
	DownloadOriginalTradeBill(ctx context.Context, r *TradeBillRequest) ([]byte, error)
	DownloadFundFlowBill(ctx context.Context, r *FundFlowBillRequest) (*FundFlowBillResponse, error)
	DownloadFundOriginalFlowBill(ctx context.Context, r *FundFlowBillRequest) ([]byte, error)
	CombinePay(ctx context.Context, r *CombinePayRequest) (*CombinePayResponse, error)
	CombineQuery(ctx context.Context, r *CombineQueryRequest) (*CombineQueryResponse, error)
	CombineClose(ctx context.Context, r *CombineCloseRequest) error
}

// Pay send a transaction and invoke wechat payment.
func (c *client) Pay(ctx context.Context, r *PayRequest) (*PayResponse, error) {
	return r.Do(ctx, c)
}

// Query send the request of query transaction.
func (c *client) Query(ctx context.Context, r *QueryRequest) (*QueryResponse, error) {
	return r.Do(ctx, c)
}

// Cert get certificates from wechat pay.
func (c *client) Cert(ctx context.Context, r *CertificatesRequest) (*CertificatesResponse, error) {
	return r.Do(ctx, c)
}

// Close send the request of close transaction.
func (c *client) Close(ctx context.Context, r *CloseRequest) error {
	return r.Do(ctx, c)
}

// Refund send the refund request and return refund response.
func (c *client) Refund(ctx context.Context, r *RefundRequest) (*RefundResponse, error) {
	return r.Do(ctx, c)
}

// QueryRefund send the refund query result.
func (c *client) QueryRefund(ctx context.Context, r *RefundQueryRequest) (*RefundQueryResponse, error) {
	return r.Do(ctx, c)
}

// DownloadTradeBill download and unmarshal the data of trade bill.
func (c *client) DownloadTradeBill(ctx context.Context, r *TradeBillRequest) (*TradeBillResponse, error) {
	return r.UnmarshalDownload(ctx, c)
}

// DownloadOriginalTradeBill download plain text of trade bill.
func (c *client) DownloadOriginalTradeBill(ctx context.Context, r *TradeBillRequest) ([]byte, error) {
	return r.Download(ctx, c)
}

// DownloadFundFlowBill download and unmarshal the data of fundflow bill.
func (c *client) DownloadFundFlowBill(ctx context.Context, r *FundFlowBillRequest) (*FundFlowBillResponse, error) {
	return r.UnmarshalDownload(ctx, c)
}

// DownloadFundOriginalFlowBill download plain text of fundflow bill.
func (c *client) DownloadFundOriginalFlowBill(ctx context.Context, r *FundFlowBillRequest) ([]byte, error) {
	return r.Download(ctx, c)
}

// CombinePay send a transaction and invoke wechat payment.
func (c *client) CombinePay(ctx context.Context, r *CombinePayRequest) (*CombinePayResponse, error) {
	return r.Do(ctx, c)
}

// CombineQuery send the request of query transaction.
func (c *client) CombineQuery(ctx context.Context, r *CombineQueryRequest) (*CombineQueryResponse, error) {
	return r.Do(ctx, c)
}

// CombineClose send the request of combine close transaction.
func (c *client) CombineClose(ctx context.Context, r *CombineCloseRequest) error {
	return r.Do(ctx, c)
}
