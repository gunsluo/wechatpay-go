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
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRefundValidate(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("client is nil")
	}

	type fields struct {
		TransactionId string
		OutTradeNo    string
		OutRefundNo   string
		Reason        string
		NotifyUrl     string
		FundsAccount  string
		Amount        RefundAmount
		GoodsDetail   []RefundGoodDetail
	}
	tests := []struct {
		name            string
		fields          fields
		want            *RefundResponse
		wantErr         bool
		wantErrContains string
	}{
		{
			name: "validate",
			fields: fields{
				TransactionId: "",
			},
			want:            nil,
			wantErr:         true,
			wantErrContains: "transaction_id can't be empty",
		},
		{
			name: "validate",
			fields: fields{
				TransactionId: "1234578945678",
			},
			want:            nil,
			wantErr:         true,
			wantErrContains: "out_refund_no can't be empty",
		},
		{
			name: "validate",
			fields: fields{
				TransactionId: "1234578945678",
				OutRefundNo:   "1234557677",
			},
			want:            nil,
			wantErr:         true,
			wantErrContains: "out_trade_no can't be empty",
		},
		{
			name: "validate",
			fields: fields{
				TransactionId: "1234578945678",
				OutTradeNo:    "123456789",
				OutRefundNo:   "123456789",
				Amount:        RefundAmount{},
			},
			want:            nil,
			wantErr:         true,
			wantErrContains: "refund can't less than 0",
		},
		{
			name: "validate",
			fields: fields{
				TransactionId: "1234578945678",
				OutTradeNo:    "123456789",
				OutRefundNo:   "123456789",
				Amount: RefundAmount{
					Refund: 1,
				},
			},
			want:            nil,
			wantErr:         true,
			wantErrContains: "total can't less than 0",
		},
		{
			name: "validate",
			fields: fields{
				TransactionId: "1234578945678",
				OutTradeNo:    "123456789",
				OutRefundNo:   "123456789",
				Amount: RefundAmount{
					Refund: 1,
					Total:  1,
				},
			},
			want:            nil,
			wantErr:         true,
			wantErrContains: "currency can't be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RefundRequest{
				TransactionId: tt.fields.TransactionId,
				OutTradeNo:    tt.fields.OutTradeNo,
				OutRefundNo:   tt.fields.OutRefundNo,
				Reason:        tt.fields.Reason,
				NotifyUrl:     tt.fields.NotifyUrl,
				FundsAccount:  tt.fields.FundsAccount,
				Amount:        tt.fields.Amount,
				GoodsDetail:   tt.fields.GoodsDetail,
			}
			got, err := r.Do(context.Background(), client)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr {
				if tt.wantErrContains != "" {
					if !strings.Contains(err.Error(), tt.wantErrContains) {
						t.Errorf("Do() error = %v, wantErr %v, don't contains %v", err, tt.wantErr, tt.wantErrContains)
						return
					}
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Do() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefundDo(t *testing.T) {

	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		req       *RefundRequest
		resp      *RefundResponse
		transport *mockTransport
		pass      bool
	}{
		{
			&RefundRequest{
				TransactionId: "for test",
				OutTradeNo:    "for test",
				OutRefundNo:   "for test",
				Reason:        "for test",
				NotifyUrl:     "http://domain.com/notify",
				FundsAccount:  "",
				Amount: RefundAmount{
					Refund:   1,
					Total:    1,
					Currency: "CNY",
				},
				GoodsDetail: nil,
			},
			&RefundResponse{
				RefundId:            "50300807092021020105990201735",
				OutRefundNo:         "S20210201151309277501",
				TransactionId:       "4200000925202101284997714292",
				OutTradeNo:          "S20210128170702357723",
				Channel:             "ORIGINAL",
				UserReceivedAccount: "支付用户零钱",
				SuccessTime:         time.Time{},
				CreateTime:          dateFromString("2021-02-01T15:13:10+08:00"),
				Status:              "PROCESSING",
				FundsAccount:        "UNAVAILABLE",
				Amount: RefundAmountDetail{
					Total:            1,
					Refund:           1,
					PayerTotal:       1,
					PayerRefund:      1,
					SettlementTotal:  1,
					SettlementRefund: 1,
					DiscountRefund:   0,
					Currency:         "CNY",
				},
				Promotion: nil,
			},
			nil,
			true,
		},
		{
			&RefundRequest{
				TransactionId: "for test",
				OutTradeNo:    "for test",
				OutRefundNo:   "for test",
				Reason:        "for test",
				NotifyUrl:     "http://domain.com/notify",
				FundsAccount:  "",
				Amount: RefundAmount{
					Refund:   1,
					Total:    1,
					Currency: "CNY",
				},
				GoodsDetail: nil,
			},
			&RefundResponse{},
			&mockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					var resp = &http.Response{
						StatusCode: http.StatusOK,
					}

					resp.Header = http.Header{}
					resp.Body = ioutil.NopCloser(strings.NewReader("{}"))
					return resp, nil
				},
			},
			false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		if c.transport != nil {
			client.config.opts.transport = c.transport
			client.secrets.clear()
		}

		resp, err := c.req.Do(ctx, client)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err: %v", c.pass, pass, err)
		}

		if err != nil {
			continue
		}

		if !reflect.DeepEqual(c.resp, resp) {
			t.Fatalf("expect %v, got %v", c.resp, resp)
		}
	}
}

func dateFromString(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
