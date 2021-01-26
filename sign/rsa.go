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
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

// LoadRSAPrivateKey load the buffer about rsa private cert and
// return private key.
func LoadRSAPrivateKey(buffer []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(buffer)
	if block == nil {
		return nil, errors.New("invalid private key")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not rsa private key")
	}

	return privateKey, nil
}

// LoadRSAPrivateKeyFromTxt load the string about rsa private key
// and return private key.
func LoadRSAPrivateKeyFromTxt(privateKeyTxt string) (*rsa.PrivateKey, error) {
	return LoadRSAPrivateKey([]byte(privateKeyTxt))
}

// LoadRSAPrivateKeyFromFile load the file about rsa private key and
// return private key.
func LoadRSAPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	privateKeyBuffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return LoadRSAPrivateKey(privateKeyBuffer)
}

// LoadRSAPublicKeyFromCert load the buffer about rsa cert and
// return public key.
func LoadRSAPublicKeyFromCert(buffer []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(buffer)
	if block == nil {
		return nil, errors.New("invalid publicKey key")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not rsa public key")
	}

	return publicKey, nil
}
