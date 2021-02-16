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

func TestNotificationAnswer(t *testing.T) {
	answer := &NotificationAnswer{}
	buffer := answer.Bytes()
	actual := string(buffer)
	expect := `{"code":"","message":""}`
	if expect != actual {
		t.Fatalf("expect %s, got %s", expect, actual)
	}
}

func TestParseHttpRequestForRefundNotification(t *testing.T) {
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

				signature := "fjUjJEN2e60uTygpCUNF+gNbm3dY8DJxyZyL1SSHkbRc4lqmYYyLzlot1PWUvzvYdbGF99SMvAGuigkXuxmgJjaFsQ05uZiK5vUhCpv6hEKEIxQxgYPp5n7wBOZX2VOaKWBQp7F3B/4R8ZZEQ9nPKbrDjUFFu+LFRqA1akmukO5MY1sxXFpxBTrqfesnNQSo7pqjigBufJnh3yjmpUNdXbmDhMIUuAWr0dDNmETeK94B4tFjqF7hGza8/WUzwXj4JTy4ZBz8irgyKX9VimILiEVPB6I3afXptSMvTMlUwBG4gpwLCSgscn+vGKQv2mNCHkPDecl9c3dMbiyTkYtR3A=="
				mockBody := `{"id":"9971e868-2144-58dd-99ae-df4c76ccde42","create_time":"2021-02-01T15:13:13+08:00","resource_type":"encrypt-resource","event_type":"REFUND.SUCCESS","summary":"退款成功","resource":{"original_type":"refund","algorithm":"AEAD_AES_256_GCM","ciphertext":"i6LL5pzT9gNoZTx3EtUmdiLPz7cQRJXa6mO+2kZBsn6aU5Rjd/m0+3YLnMXFYT2AKUUr0Iel3iQQw2rN834d305VcR7BQHXaY7qJ2fc0lZSqkx3aszF6yNRQ3rlvHBsqlOjQwPYFydbA0Fu0TO3F3aqbvJfYJ541O/O/EZeYZX31e6nbvY3ZjGGhXnW/CqWYkUSu8v9K3Q/KGP94VVTw/dvqyoOKN9NhGH4YV/62My6HUA6Khf2BQsYhsSqPJ3RzeEZiB6cxwWposXNqtjrwUI+Y7IrlJjwRjg8i0SPyUaDkTEybtdBTFSNzVSt7F5W32qYksgHozjlEIqQ4G/OMZEP+XUelVWeoGXkPgEC6WYHDZixyPlfCLNRxkxkaBhgs2+PKYAVBYtag8Je1/88oQ7Ms+qcUjHTXTRJJgRsBXZLPT20dFDySOOI7iIOkd0V+B8s5/NvUGqiCr3q4kAZDEI9H83qqA0QZvfH5zDr/l5VCh0ko7L8DTibF1w5WcILrnxuJcA==","associated_data":"refund","nonce":"QOXEHLl2XppO"}}`

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

				mockBody := `{"id":"9971e868-2144-58dd-99ae-df4c76ccde42","create_time":"2021-02-01T15:13:13+08:00","resource_type":"encrypt-resource","event_type":"REFUND.SUCCESS","summary":"退款成功","resource":{"original_type":"refund","algorithm":"AEAD_AES_256_GCM","ciphertext":"i6LL5pzT9gNoZTx3EtUmdiLPz7cQRJXa6mO+2kZBsn6aU5Rjd/m0+3YLnMXFYT2AKUUr0Iel3iQQw2rN834d305VcR7BQHXaY7qJ2fc0lZSqkx3aszF6yNRQ3rlvHBsqlOjQwPYFydbA0Fu0TO3F3aqbvJfYJ541O/O/EZeYZX31e6nbvY3ZjGGhXnW/CqWYkUSu8v9K3Q/KGP94VVTw/dvqyoOKN9NhGH4YV/62My6HUA6Khf2BQsYhsSqPJ3RzeEZiB6cxwWposXNqtjrwUI+Y7IrlJjwRjg8i0SPyUaDkTEybtdBTFSNzVSt7F5W32qYksgHozjlEIqQ4G/OMZEP+XUelVWeoGXkPgEC6WYHDZixyPlfCLNRxkxkaBhgs2+PKYAVBYtag8Je1/88oQ7Ms+qcUjHTXTRJJgRsBXZLPT20dFDySOOI7iIOkd0V+B8s5/NvUGqiCr3q4kAZDEI9H83qqA0QZvfH5zDr/l5VCh0ko7L8DTibF1w5WcILrnxuJcA==","associated_data":"refund","nonce":"QOXEHLl2XppO"}}`

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
		n := RefundNotification{}
		req := c.newReq()
		_, err := n.ParseHttpRequest(client, req)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err %v", c.pass, pass, err)
		}
	}
}

func TestParseForRefundNotification(t *testing.T) {
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
				Signature: "fjUjJEN2e60uTygpCUNF+gNbm3dY8DJxyZyL1SSHkbRc4lqmYYyLzlot1PWUvzvYdbGF99SMvAGuigkXuxmgJjaFsQ05uZiK5vUhCpv6hEKEIxQxgYPp5n7wBOZX2VOaKWBQp7F3B/4R8ZZEQ9nPKbrDjUFFu+LFRqA1akmukO5MY1sxXFpxBTrqfesnNQSo7pqjigBufJnh3yjmpUNdXbmDhMIUuAWr0dDNmETeK94B4tFjqF7hGza8/WUzwXj4JTy4ZBz8irgyKX9VimILiEVPB6I3afXptSMvTMlUwBG4gpwLCSgscn+vGKQv2mNCHkPDecl9c3dMbiyTkYtR3A==",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{"id":"9971e868-2144-58dd-99ae-df4c76ccde42","create_time":"2021-02-01T15:13:13+08:00","resource_type":"encrypt-resource","event_type":"REFUND.SUCCESS","summary":"退款成功","resource":{"original_type":"refund","algorithm":"AEAD_AES_256_GCM","ciphertext":"i6LL5pzT9gNoZTx3EtUmdiLPz7cQRJXa6mO+2kZBsn6aU5Rjd/m0+3YLnMXFYT2AKUUr0Iel3iQQw2rN834d305VcR7BQHXaY7qJ2fc0lZSqkx3aszF6yNRQ3rlvHBsqlOjQwPYFydbA0Fu0TO3F3aqbvJfYJ541O/O/EZeYZX31e6nbvY3ZjGGhXnW/CqWYkUSu8v9K3Q/KGP94VVTw/dvqyoOKN9NhGH4YV/62My6HUA6Khf2BQsYhsSqPJ3RzeEZiB6cxwWposXNqtjrwUI+Y7IrlJjwRjg8i0SPyUaDkTEybtdBTFSNzVSt7F5W32qYksgHozjlEIqQ4G/OMZEP+XUelVWeoGXkPgEC6WYHDZixyPlfCLNRxkxkaBhgs2+PKYAVBYtag8Je1/88oQ7Ms+qcUjHTXTRJJgRsBXZLPT20dFDySOOI7iIOkd0V+B8s5/NvUGqiCr3q4kAZDEI9H83qqA0QZvfH5zDr/l5VCh0ko7L8DTibF1w5WcILrnxuJcA==","associated_data":"refund","nonce":"QOXEHLl2XppO"}}`),
			},
			true,
		},
		{
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "fjUjJEN2e60uTygpCUNF+gNbm3dY8DJxyZyL1SSHkbRc4lqmYYyLzlot1PWUvzvYdbGF99SMvAGuigkXuxmgJjaFsQ05uZiK5vUhCpv6hEKEIxQxgYPp5n7wBOZX2VOaKWBQp7F3B/4R8ZZEQ9nPKbrDjUFFu+LFRqA1akmukO5MY1sxXFpxBTrqfesnNQSo7pqjigBufJnh3yjmpUNdXbmDhMIUuAWr0dDNmETeK94B4tFjqF7hGza8/WUzwXj4JTy4ZBz8irgyKX9VimILiEVPB6I3afXptSMvTMlUwBG4gpwLCSgscn+vGKQv2mNCHkPDecl9c3dMbiyTkYtR3A==",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{`),
			},
			false,
		},
		{
			&Result{
				Timestamp: mockTimestamp,
				Nonce:     mockNonce,
				Signature: "fjUjJEN2e60uTygpCUNF+gNbm3dY8DJxyZyL1SSHkbRc4lqmYYyLzlot1PWUvzvYdbGF99SMvAGuigkXuxmgJjaFsQ05uZiK5vUhCpv6hEKEIxQxgYPp5n7wBOZX2VOaKWBQp7F3B/4R8ZZEQ9nPKbrDjUFFu+LFRqA1akmukO5MY1sxXFpxBTrqfesnNQSo7pqjigBufJnh3yjmpUNdXbmDhMIUuAWr0dDNmETeK94B4tFjqF7hGza8/WUzwXj4JTy4ZBz8irgyKX9VimILiEVPB6I3afXptSMvTMlUwBG4gpwLCSgscn+vGKQv2mNCHkPDecl9c3dMbiyTkYtR3A==",
				SerialNo:  mockSerialNo,
				Body:      []byte(`{"id":"9971e868-2144-58dd-99ae-df4c76ccde42","create_time":"2021-02-01T15:13:13+08:00","resource_type":"encrypt-resource","event_type":"REFUND.SUCCESS","summary":"退款成功","resource":{"original_type":"refund","algorithm":"AEAD_AES_256_GCM","ciphertext":"i6LL5pzT9gNoZTx3EtUmdiLPz7cQRJXa6mO","associated_data":"refund","nonce":"QOXEHLl2XppO"}}`),
			},
			false,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		n := RefundNotification{}
		_, err := n.Parse(ctx, client, c.result)
		pass := err == nil
		if pass != c.pass {
			t.Fatalf("expect %v, got %v, err %v", c.pass, pass, err)
		}
	}
}
