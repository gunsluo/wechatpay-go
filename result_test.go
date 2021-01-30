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
	"testing"
)

func TestResult(t *testing.T) {
	cases := []struct {
		result *Result
		expect *PayResponse
		pass   bool
	}{
		{
			&Result{
				Body: []byte(`{"code_url":"https://xxx.com"}`),
			},
			&PayResponse{"https://xxx.com"},
			true,
		},
		{
			&Result{
				Err: &Error{},
			},
			&PayResponse{"https://xxx.com"},
			false,
		},
		{
			&Result{
				Body: []byte(``),
			},
			&PayResponse{"https://xxx.com"},
			true,
		},
		{
			&Result{},
			&PayResponse{"https://xxx.com"},
			true,
		},
		{
			&Result{
				Body: []byte(`{"`),
			},
			&PayResponse{"https://xxx.com"},
			false,
		},
	}

	for _, c := range cases {
		dest := &PayResponse{}
		err := c.result.Scan(dest)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err %v", c.pass, pass, err)
		}
	}
}

func TestError(t *testing.T) {
	cases := []struct {
		err    *Error
		expect string
	}{
		{
			&Error{400, "code", "message"},
			`{"status":400,"code":"code","message":"message"}`,
		},
		{
			nil,
			"{}",
		},
	}

	for _, c := range cases {
		actual := c.err.Error()
		if actual != c.expect {
			t.Fatalf("expect %s, got %s", c.expect, actual)
		}
	}
}
