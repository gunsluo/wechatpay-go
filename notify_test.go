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
	"strconv"
	"strings"
	"testing"
)

func TestParseHttpRequestForPayNotification(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		newReq func() *http.Request
		pass   bool
	}{
		{
			func() *http.Request {
				req := &http.Request{
					Header: http.Header{},
				}

				signature := "Jook1G0Ex2xkvw5isZNY8Pvxj30X6HOCLNwMBh0wpRCU0LMTD+wQqHCENpYcsaMM/6vFMsRXtZnKldRk1dFmzpLOT8Rh1SwfMp/61oz7Eyh9+y1p2QkC2EW9dEnZk3gl7j5WcSsncy8ccM4ohfZVwQLslZwzKKaLxg5F5MTeiP/0ykYdFHOqIKdp9QMlly0Yb9aUXiVe19u3PEIOUkAawr9vD7EL5VHtnuer90ADrO9b+p4MAFxL1QfqshNhb4KeDjyVAzOqHjkThqAeuY1wv8KjoeVpZOxxrdSAoYcek2c2A8ywKWNMZi/k0Wwpu05UN498a39tKdHPZrqb6Qt4ZA=="
				mockBody := `{"id":"b62e271c-3389-58a0-8146-4a704966e8f1","create_time":"2021-01-28T17:07:11+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"original_type":"transaction","algorithm":"AEAD_AES_256_GCM","ciphertext":"yuKJXXxnqVMulBUy5NoriSab/S9aen3wXNYLqGdvBfxsWmN9JAFAMXO3LgDFPqNeZMrkSmQyFa981IVxLvWHzwrzlBtJk+hOwnxTgDxc8SsGt39QkRBbfGR8rutMr3Goiq03ygWjMA6I+n6qhqQ/zS0/bMIB1dQoFZBSCKiLp8VHbGDLirh9MqYRa7MKJEYziPF2DmdtRHvXie4AWSxcV6hq8Ufao9FQooLOA2gD/9JA+L6BqquOPOnStExxH26cK7QgFFAf22GP7JKXnMH0LF3lJrK6ZMQ7iTXvVxv/q6j3SwUbyWVKmXdMJTqnXtU4H90DjRC6It4cOavr3Gz6xeVyv4S3i1qdAD8rAqgjjF1QWnUQtIm4/TdOw3ro0L73VI07H8c9O6VX/U0TcGMJJrAKMJ/yBZlD6owliffy/pzceEG/MV27euHDS5VW/m23tokNy2G1XJu1T3sUzEUsNil7vngBLYHGEGNw6brOYxwxXEUI2n0tSJOG8upiSGmN0fOnWbPoN9YqtuIhvY4xKOJpKwQrNJSm+ybNrugAwbLf/HMATxK6dGk9RQK8Nn9PHSRSPmTU5sci6zzFGAEHKQ==","associated_data":"transaction","nonce":"fG1l57vn9BCX"}}`

				req.Header.Set("Wechatpay-Nonce", mockNonce)
				req.Header.Set("Wechatpay-Signature", signature)
				req.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(mockTimestamp, 10))
				req.Header.Set("Wechatpay-Serial", mockSerialNo)
				req.Body = ioutil.NopCloser(strings.NewReader(mockBody))

				return req
			},
			true,
		},
		{
			func() *http.Request {
				req := &http.Request{
					Header: http.Header{},
				}

				mockBody := `{"id":"b62e271c-3389-58a0-8146-4a704966e8f1","create_time":"2021-01-28T17:07:11+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"original_type":"transaction","algorithm":"AEAD_AES_256_GCM","ciphertext":"yuKJXXxnqVMulBUy5NoriSab/S9aen3wXNYLqGdvBfxsWmN9JAFAMXO3LgDFPqNeZMrkSmQyFa981IVxLvWHzwrzlBtJk+hOwnxTgDxc8SsGt39QkRBbfGR8rutMr3Goiq03ygWjMA6I+n6qhqQ/zS0/bMIB1dQoFZBSCKiLp8VHbGDLirh9MqYRa7MKJEYziPF2DmdtRHvXie4AWSxcV6hq8Ufao9FQooLOA2gD/9JA+L6BqquOPOnStExxH26cK7QgFFAf22GP7JKXnMH0LF3lJrK6ZMQ7iTXvVxv/q6j3SwUbyWVKmXdMJTqnXtU4H90DjRC6It4cOavr3Gz6xeVyv4S3i1qdAD8rAqgjjF1QWnUQtIm4/TdOw3ro0L73VI07H8c9O6VX/U0TcGMJJrAKMJ/yBZlD6owliffy/pzceEG/MV27euHDS5VW/m23tokNy2G1XJu1T3sUzEUsNil7vngBLYHGEGNw6brOYxwxXEUI2n0tSJOG8upiSGmN0fOnWbPoN9YqtuIhvY4xKOJpKwQrNJSm+ybNrugAwbLf/HMATxK6dGk9RQK8Nn9PHSRSPmTU5sci6zzFGAEHKQ==","associated_data":"transaction","nonce":"fG1l57vn9BCX"}}`

				req.Header.Set("Wechatpay-Nonce", mockNonce)
				req.Header.Set("Wechatpay-Timestamp", "xxx")
				req.Header.Set("Wechatpay-Serial", mockSerialNo)
				req.Body = ioutil.NopCloser(strings.NewReader(mockBody))

				return req
			},
			false,
		},
	}

	for _, c := range cases {
		n := PayNotification{}
		req := c.newReq()
		_, err := n.ParseHttpRequest(client, req)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err %v", c.pass, pass, err)
		}
	}
}

func TestParseForPayNotification(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("client is nil")
	}

	cases := []struct {
		result *Result
		pass   bool
	}{
		{
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "Jook1G0Ex2xkvw5isZNY8Pvxj30X6HOCLNwMBh0wpRCU0LMTD+wQqHCENpYcsaMM/6vFMsRXtZnKldRk1dFmzpLOT8Rh1SwfMp/61oz7Eyh9+y1p2QkC2EW9dEnZk3gl7j5WcSsncy8ccM4ohfZVwQLslZwzKKaLxg5F5MTeiP/0ykYdFHOqIKdp9QMlly0Yb9aUXiVe19u3PEIOUkAawr9vD7EL5VHtnuer90ADrO9b+p4MAFxL1QfqshNhb4KeDjyVAzOqHjkThqAeuY1wv8KjoeVpZOxxrdSAoYcek2c2A8ywKWNMZi/k0Wwpu05UN498a39tKdHPZrqb6Qt4ZA==",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{"id":"b62e271c-3389-58a0-8146-4a704966e8f1","create_time":"2021-01-28T17:07:11+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"original_type":"transaction","algorithm":"AEAD_AES_256_GCM","ciphertext":"yuKJXXxnqVMulBUy5NoriSab/S9aen3wXNYLqGdvBfxsWmN9JAFAMXO3LgDFPqNeZMrkSmQyFa981IVxLvWHzwrzlBtJk+hOwnxTgDxc8SsGt39QkRBbfGR8rutMr3Goiq03ygWjMA6I+n6qhqQ/zS0/bMIB1dQoFZBSCKiLp8VHbGDLirh9MqYRa7MKJEYziPF2DmdtRHvXie4AWSxcV6hq8Ufao9FQooLOA2gD/9JA+L6BqquOPOnStExxH26cK7QgFFAf22GP7JKXnMH0LF3lJrK6ZMQ7iTXvVxv/q6j3SwUbyWVKmXdMJTqnXtU4H90DjRC6It4cOavr3Gz6xeVyv4S3i1qdAD8rAqgjjF1QWnUQtIm4/TdOw3ro0L73VI07H8c9O6VX/U0TcGMJJrAKMJ/yBZlD6owliffy/pzceEG/MV27euHDS5VW/m23tokNy2G1XJu1T3sUzEUsNil7vngBLYHGEGNw6brOYxwxXEUI2n0tSJOG8upiSGmN0fOnWbPoN9YqtuIhvY4xKOJpKwQrNJSm+ybNrugAwbLf/HMATxK6dGk9RQK8Nn9PHSRSPmTU5sci6zzFGAEHKQ==","associated_data":"transaction","nonce":"fG1l57vn9BCX"}}`),
			},
			true,
		},
		{
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "Jook1G0Ex2xkvw5isZNY8Pvxj30X6HOCLNwMBh0wpRCU0LMTD+wQqHCENpYcsaMM/6vFMsRXtZnKldRk1dFmzpLOT8Rh1SwfMp/61oz7Eyh9+y1p2QkC2EW9dEnZk3gl7j5WcSsncy8ccM4ohfZVwQLslZwzKKaLxg5F5MTeiP/0ykYdFHOqIKdp9QMlly0Yb9aUXiVe19u3PEIOUkAawr9vD7EL5VHtnuer90ADrO9b+p4MAFxL1QfqshNhb4KeDjyVAzOqHjkThqAeuY1wv8KjoeVpZOxxrdSAoYcek2c2A8ywKWNMZi/k0Wwpu05UN498a39tKdHPZrqb6Qt4ZA==",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{`),
			},
			false,
		},
		{
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{"id":"b62e271c-3389-58a0-8146-4a704966e8f1","create_time":"2021-01-28T17:07:11+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"original_type":"transaction","algorithm":"AEAD_AES_256_GCM","ciphertext":"yuKJXXxnqVMulBUy5NoriSab/S9aen3wXNYLqGdvBfxsWmN9JAFAMXO3LgDFPqNeZMrkSmQyFa981IVxLvWHzwrzlBtJk+hOwnxTgDxc8SsGt39QkRBbfGR8rutMr3Goiq03ygWjMA6I+n6qhqQ/zS0/bMIB1dQoFZBSCKiLp8VHbGDLirh9MqYRa7MKJEYziPF2DmdtRHvXie4AWSxcV6hq8Ufao9FQooLOA2gD/9JA+L6BqquOPOnStExxH26cK7QgFFAf22GP7JKXnMH0LF3lJrK6ZMQ7iTXvVxv/q6j3SwUbyWVKmXdMJTqnXtU4H90DjRC6It4cOavr3Gz6xeVyv4S3i1qdAD8rAqgjjF1QWnUQtIm4/TdOw3ro0L73VI07H8c9O6VX/U0TcGMJJrAKMJ/yBZlD6owliffy/pzceEG/MV27euHDS5VW/m23tokNy2G1XJu1T3sUzEUsNil7vngBLYHGEGNw6brOYxwxXEUI2n0tSJOG8upiSGmN0fOnWbPoN9YqtuIhvY4xKOJpKwQrNJSm+ybNrugAwbLf/HMATxK6dGk9RQK8Nn9PHSRSPmTU5sci6zzFGAEHKQ==","associated_data":"transaction","nonce":"fG1l57vn9BCX"}}`),
			},
			false,
		},
		{
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "g0A/VGU569/iS8MtR2SRfFg0YOqSzKYipRTJebnm6bLsWgSwWL92KMoAwNrP8qMqf1LUKWWrb2o0XpmLt2DMV7MStrJNmcViV6yVKzVRuS2XE3kUiQNFnbIdvNRiLI0hLDGA9W6dH5YUF/yVPanRo3rBLID8mFxD1tz2XyVpVKsDu7EhUQmwCKpgZ0a+lPILxZfMjnVI7VL6AFuf/iCrb/xaoVzGCJ1hLcPe7QV89MqNp2M4IP1wbiBqJezC7vBF/t/Rtyn+kK+My/S7iB+XDrHHXn/6ldp7RXBcDmjVbnp551oS2s8jyBFN1z/K+BIg+gYmyN9vOGgRFWcV2NGpUQ==",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{"id":"b62e271c-3389-58a0-8146-4a704966e8f1","create_time":"2021-01-28T17:07:11+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"original_type":"transaction","algorithm":"AEAD_AES_256_GCM","ciphertext":"yuKJXXxnqVMulBUy5NoriSab/S9aen3wXNYLqGdvBfxsWmN9JAFAMXO3LgDFPqNeZMrkSmQyFa981IVxLvWHzwrzlBtJk+hOwnxTgDxc8SsGt39QkRBbfGR8rutMr3Goiq03ygWjMA6I+n6qhqQ/zS0/bMIB1dQoFZBSCKiLp8VHbGDLirh9MqYRa7MKJEYziPF2DmdtRHvXie4AWSxcV6hq8Ufao9FQooLOA2gD/9JA+L6BqquOPOnStExxH26cK7QgFFAf22GP7JKXnMH0LF3lJrK6ZMQ7iTXvVxv/q6j3SwUbyWVKmXdMJTqnXtU4H90DjRC6It4cOavr3Gz6xeVyv4S3i1qdAD8rAqgjjF1QWnUQtIm4/TdOw3ro0L73VI07H8c9O6VX/U0TcGMJJrAKMJ/yBZlD6owliffy/pzceEG/MV27euHDS5VW/m23tokNy2G1XJu1T3sUzEUsNil7vngBLYHGEGNw6brOYxwxXEUI2n0tSJOG8upiSGmN0fOnWbPoN9YqtuIhvY4xKOJpKwQrNJSm+ybNrugAwbLf/HMATxK6dGk9RQK8Nn9PkHYuBnwDft8oxSDkqLO7KA==","associated_data":"transaction","nonce":"fG1l57vn9BCX"}}`),
			},
			false,
		},
		{
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "K/aIQbi+bDKKC0nFc9/M/nOe2M/nnBaHbYs3Gf7ZzjD6Cq/thfph+0+h+8gzRay91z+f2/ggSQ18+ZxkI/VBpDFWE4ZoPuwoHk4BDvU6gy39Xcb1dcKs3sfjVu3wTyTrWychJ9CwBE9mVjsHlsFFrrz3+Nk+rBwJZp1h1B94PPImxB//bcXKfcBl/TXzDGFyJcgcpiyMwFe9x6EGIhJh51AXVe17W/qUSZtR8YRSaeqHdMasgTBOPLq93QfAkwz7wmxzUYeIMdmYJMU+Arhbc3EdmAnGfSuR2M6Pl65Y867RwgvAegYX2j8B7Xor7B+F+BdPqjg0mkFuvmHESToZSA==",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{"id":"b62e271c-3389-58a0-8146-4a704966e8f1","create_time":"2021-01-28T17:07:11+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"original_type":"transaction","algorithm":"AEAD_AES_256_GCM","ciphertext":"yriQVTmulDmba+GD/be3QUMtoDg=","associated_data":"transaction","nonce":"fG1l57vn9BCX"}}`),
			},
			false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		n := PayNotification{}
		_, err := n.Parse(ctx, client, c.result)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err %v", c.pass, pass, err)
		}
	}
}
