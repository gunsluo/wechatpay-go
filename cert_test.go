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
	"reflect"
	"testing"
)

func TestDoForCert(t *testing.T) {
	client, err := mockNewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fail()
	}

	cases := []struct {
		req  *CertificatesRequest
		resp *CertificatesRespone
	}{
		{
			&CertificatesRequest{},
			&CertificatesRespone{
				Certificates: []Certificate{
					{
						SerialNo:      mockSerialNo,
						EffectiveTime: "2020-09-17T14:26:23+08:00",
						ExpireTime:    "2025-09-16T14:26:23+08:00",
						Encrypt: EncryptCertificate{
							Algorithm:  "AEAD_AES_256_GCM",
							Nonce:      "eabb3e044577",
							Associated: "certificate",
							CipherText: "/M2eAJyVx/0y8JOErsNEWbYpikwKMS0hDahBYrR9Tnqvaxw/WLMHyLq7G3GUoWx3NSwYZlSZ+1JxAMTd4yge1B8bxY7OLrDkXm+BBDVypy5jCi/gcTQduTJpR4nRcBRYtEIxLGLrVaUXlDjDa4nM0mUPk6XA7AAUUAl3z5lYISapsFYUuHO9splBrmUESHxzRhSfsTyW68ll8o+ND7xA5R94slxzZIVdVg2Tz/3uXi5X1Qu5oi9Dn7pFdHD7++msMB7rgSJUTIFMwZ2GhAX3f/vVWemSMCymPPxzYxdiGFJJ8oBaIn+17pwulmz6NodFS0ilJr9wBs/05gqxe5L6S64ApwXNTfq3YJFVIU6munBaHomRZqsMg3MQlji9yNLBdKO2hk2rq/jCaBLsqcrCHEMEEULA5/1ImeYEkKcX2vIiVtKX8WxxP4M/Gq7btAQZVGzvczopb3wZNu1QLnzC13ov0pB5BPMhrx0tE4rLuZ5d+uzGOwuI8CvqOa+8TQ0DNGNaEA/IPrMJCVvmLrDi/aMQB+P4mO9BhUlfGHwQL7Q0anHzZaGHGkYyEGoTPmqQcY1mRbVcXDpIGn7rfHgiXnQTurB886T//ddhcv1/LQmcohSveZJAltcaDlmeqMgc+bXsOlAy6JNIIVPJ04ysI+V7nc0O4k4A32ZYA1hK52CU1YWz3vMoaaHVr/t6AF3dVWE1CphhNIwGbaz9M1sgEsWwT8LKLG5csgVwG20LO8wmLkxNUQ4fSkMdC+2Qv+rSFd8rlT1j+sYEbPVq6E6URkYPUKMqI1mEEudU1Rx0bE/pjj7+++0gX1H7sHp4+02KLdWS27gptHVXdDjNFPyCEshfVL2B8aEhq8PxSDG5zTqWHrKBAl04WU3kjlSsKZPrpKyhpIrKbEZHcrip3wOGeMf+4XDoZ8Iq8KoM8R6m8wkWi0GAW4G743O44PxHFvljKDIkIQm8gWV37jC3+qb/ZwUDxHONw3tHMH8XWsCVq1KAtKeE/iE9CCmE+ht7K4B+w0DeqKEicm0dkdjuFc9IgFa1W+q0HqGFI2Snd6ZX6crUy1I1vkRTQRj1mqjaP7dFOFV0JMpK/4CKMruZfUilNfOnSoKqHA2jPQ3f4ro0H22bF/PNhOWXp6Tzl5ZVbIFBIMdD9+ocq1lDH7vcBfKVwUltKl7jgI9HlpCDPZp++Mt3C4lPDzP/XrqorJnFBKw8eMBHS7N+jDhzhqJnI3ldwlGxUsqS/hj+jUUPpYINe/UtVwlOBi/tfuEfv47H5YgbP+Y3dz78a6KJUcA7caPSSqX+8LBcwEEZELXR8gU/AxwoDAsHM1pb7wc9fslct+awivfRi47AJtFeeZMGF6bb14VnbzvIZdpZRBIzHlvUqP+t8ZKEUvEJ+lVk7vv0/ySWBZbt0oA5XQ2RVwgzKGOgfMzZafsWAqrq1PGYjJqBbm/hudPtqsBridW/QjoE2Bp+Qnp8mWhdlSP8dgdeefLEeZGUSJx0Tzu2hBveEz7jMNQSOyg8HEE=",
						},
					},
				},
			},
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		resp, err := c.req.Do(ctx, client)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c.resp, resp) {
			t.Fatalf("expect %v, got %v", c.resp, resp)
		}
	}
}
