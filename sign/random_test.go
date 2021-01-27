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
	"testing"
)

func TestRandomHex(t *testing.T) {
	hex := randomHex(10)
	if len(hex) != 10 {
		t.Fail()
	}
}

func TestRandomBytesMod(t *testing.T) {
	defer func() { recover() }()
	randomBytesMod(10, 0)
	t.Errorf("did not panic")
}

func BenchmarkRandomHex(b *testing.B) {
	checkDuplicate := map[string]struct{}{}
	for n := 0; n < b.N; n++ {
		h := randomHex(32)
		if _, ok := checkDuplicate[h]; ok {
			b.Fatal("duplicate random string")
		}
	}
}
