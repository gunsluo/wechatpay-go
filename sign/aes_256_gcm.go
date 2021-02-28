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
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// DecryptByAes256Gcm uses algorithm aes-256-gcm to decrypt text.
// The key argument should be the AES key, either 16, 24, or
// 32 bytes to select AES-128, AES-192, or AES-256.
func DecryptByAes256Gcm(key, nonce, additionalData []byte, cipherText string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	cipherBuffer, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	plainText, err := aesGcm.Open(nil, nonce, cipherBuffer, additionalData)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// EncryptByAes256Gcm uses algorithm aes-256-gcm to encrypt text
// and return a base64 string. The key argument should be the AES key,
// either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
func EncryptByAes256Gcm(key, nonce, additionalData []byte, plainText string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	cipherText := aesGcm.Seal(nil, nonce, []byte(plainText), additionalData)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
