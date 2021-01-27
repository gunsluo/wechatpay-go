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

package sign

import (
	"fmt"
	"testing"
)

func TestEncryptByAes256Gcm(t *testing.T) {
	cases := []struct {
		key    []byte
		noce   []byte
		data   []byte
		text   string
		expect bool
	}{
		{
			[]byte("AES256Key-32Characters1234567890"),
			[]byte("eabb3e044577"),
			[]byte("certificate"),
			"exampleplaintext",
			true,
		},
		{
			[]byte("AES256Key-"),
			[]byte("eabb3e044577"),
			[]byte("certificate"),
			"exampleplaintext",
			false,
		},
	}

	for _, c := range cases {
		_, err := EncryptByAes256Gcm(c.key, c.noce, c.data, c.text)
		expect := err == nil
		if c.expect != expect {
			t.Fatalf("expect %v, got %v, %v", c.expect, expect, err)
		}
	}
}

func TestDecryptByAes256Gcm(t *testing.T) {
	cases := []struct {
		key    []byte
		noce   []byte
		data   []byte
		secret string
		expect bool
	}{
		{
			[]byte("AES256Key-32Characters1234567890"),
			[]byte("eabb3e044577"),
			[]byte("certificate"),
			"tJjSQMG758oX39qpn/RoZPZ3qh8LRIIwcnQeFhU/alQ=",
			true,
		},
		{
			[]byte("AES256Key-"),
			[]byte("eabb3e044577"),
			[]byte("certificate"),
			"tJjSQMG758oX39qpn/RoZPZ3qh8LRIIwcnQeFhU/alQ=",
			false,
		},
		{
			[]byte("AES256Key-32Characters1234567890"),
			[]byte("eabb3e044577"),
			[]byte("certificate"),
			"tJjSQMG75/RoZP/alQ=",
			false,
		},
		{
			[]byte("AES256Key-32Characters1234567890"),
			[]byte("eabb3e044577"),
			[]byte("certificate"),
			"exampleplaintext",
			false,
		},
	}

	for _, c := range cases {
		_, err := DecryptByAes256Gcm(c.key, c.noce, c.data, c.secret)
		expect := err == nil
		if c.expect != expect {
			t.Fatalf("expect %v, got %v, %v", c.expect, expect, err)
		}
	}
}

func TestAes256Gcm(t *testing.T) {
	cases := []struct {
		key  []byte
		noce []byte
		data []byte
		text string
	}{
		{
			[]byte("AES256Key-32Characters1234567890"),
			[]byte("eabb3e044577"),
			[]byte("certificate"),
			"exampleplaintext",
		},
	}

	for _, c := range cases {
		secret, err := EncryptByAes256Gcm(c.key, c.noce, c.data, c.text)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("-->", secret)

		plain, err := DecryptByAes256Gcm(c.key, c.noce, c.data, secret)
		if err != nil {
			t.Fatal(err)
		}

		plainTxt := string(plain)
		if plainTxt != c.text {
			t.Fatal("invalid aes-256-gcm")
		}
	}
}
