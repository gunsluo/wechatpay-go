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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

// SignatureSHA256WithRSA calculates the signature of hashed
// using SHA256 with RSA.
func SignatureSHA256WithRSA(privateKey *rsa.PrivateKey, plain []byte) (string, error) {
	h := sha256.New()
	h.Write(plain)
	d := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, d)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySHA256WithRSA verify that the signature is available
// using SHA256 with RSA.
func VerifySHA256WithRSA(publicKey *rsa.PublicKey, signature string, plain []byte) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256(plain)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], sig)
	if err != nil {
		return err
	}

	return nil
}
