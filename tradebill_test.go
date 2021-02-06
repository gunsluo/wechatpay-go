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
	"fmt"
	"testing"

	"github.com/gunsluo/wechatpay-go/v3/sign"
)

func TestA(t *testing.T) {
	method := "GET"
	url := "https://api.mch.weixin.qq.com/v3/billdownload/file?token=xxx"

	reqSign := &sign.RequestSignature{
		Method:    method,
		Url:       url,
		Timestamp: 1554208460,
		Nonce:     "593BEC0C930BF1AFEB40B4A08C8FB242",
		//Body      []byte
	}

	privateKey, err := sign.LoadRSAPrivateKeyFromTxt(mockRSAPrivateKeyCert)
	if err != nil {
		panic(err)
	}

	signature, err := sign.GenerateSignature(privateKey,
		reqSign, "1900009191", "1DDE55AD98ED71D6EDD4A4A16996DE7B47773A8C")
	if err != nil {
		panic(err)
	}

	fmt.Println("-->", signature)
}

const mockRSAPrivateKeyCert = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCprsmcXPHqLtnP
oPDGUoMULK2WOo5FW8c72Svnqn/4aXPaJhlOtPxtX2frqIhTjwcOs6hNm3XFTGBL
MrdB94YQvj+Q7P12GNmxXG+9Ms+uUyJToYjlYDAG6UFKE10Jkm9cDGuLSkekU1Ao
rKE1G1wndH37w4AzVXoGBQ3NIiyW8jIm8Zi3/WNCVpHUoXYUuyhFEZ23fXytnps4
hARgg6NvPncIKtWvlUh85ZVOSsqc1T8dFaeDRXaj7r3jdJJ74tsGRMvZyUipJXyE
3uR2QkrGyia+0phDpC6zeMMpP+MQO9ohh+xQWBCeyvQjjnOPAlGThl+ThfXImU30
HL17oHdBAgMBAAECggEADm6FSz1Efgx6DgS8NcHy0BZ0tSBJ1XBW46o2579Cnxgo
+FbhNCaEibDhn9N3tNOnYAK7v84HGD7EueCYYY3x4x6rPWJKtG6spT8dadQWgdck
RkSo5glmTFAuc2RuN1AzFHsh8njg2wMTAEKee2vWTKzFwlIAZ11PwY9Qey/65uOT
Bi8q1Rssu6xofNadO5MbqMJ1Tl8DDIaLGnzTzbHrk9thBUo1FwFjJWTVI7nz2En4
Yc/G1/LQJfiQ31F+lkL3j6ABRJqtsgb07r9H/hT6+fd1hGDt2qKuS+E1mLDp9fHw
n6UyS4HyB7DA/XtFZ9z0VtAlmcGoUkyJLtXjEmwsGQKBgQDeCtE3spULpC7VPqk1
xv034C6zybZ7y8kSKwRvyYwkzdgSRgVaKsTVb8RNYor8hoGrVgdXFqQUI8O/v1cN
9wFoGYJT0LHre/YzOg31TkQkBfHHCFH/L50uOJcIQueftctz5Bwj6bJO/ih5iIAK
yjrHse4PdIiEJfz2D9hc4wnxrwKBgQDDogrWlUCTj2fvmZfkWR3Hbs0kIHd7zjIk
bJJONGtD8gE4i562tajC1mKoQEwt4YSwWsBkGAw1LhvMROQFT6AOaIIhHNex1Z3t
c2gAdEeWOMmzZnnhwWzTiYJomixrFkmEwT3EJK89GO3E0FH5S+G1P1tNXq38Vpty
1YVqOgMSDwKBgFrzuWGEQDMljJ2C7lL98KlbpiW1AY/SGMndXxLfTw2gV9qcXgLi
NABtqM4+CEqKWkExmw4cUxeA0uUPXnx06lmW4WCtwsN/4oh3RlJuPdE3siLiEJxk
B5FwUsVqinBMSktta+12A7kBuNiXhkNlNRCpnKcuB+GBog20zd62jVM3AoGBALcA
zFazQ7dFfRq7eUUYwCyhT7Et1dewqWM9VRdnHbhvmAjHQu7zvCyW069Ehn6c6bz3
B+YaQME2orZQ82SsebNAvAoxquwmQhevz2gtXhH+iWASyo0Onbi8d4tWPZrnPFq9
UgQ7tNnYigOEREqKW1drLwOPP/4/Hicr6iPWpKytAoGAEQ6J/RB/olEAC46ACoFo
FBgA+GUbDB0xBcA2inEt3q//208YMkjnKM871n89HpAgms5xrK32T69lduebk7Ar
9wWvkJVUwI9VDXomCFQqtiGzHlTl1Xq31BfeIDyq1ayQmTkRpRqIagbDZVtM+ha/
0I2SEzTObt07wcYcYG2Chvg=
-----END PRIVATE KEY-----`
