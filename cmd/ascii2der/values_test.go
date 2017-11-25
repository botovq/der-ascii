// Copyright 2015 The DER ASCII Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/google/der-ascii/lib"
)

var decodeTagStringTests = []struct {
	input string
	tag   lib.Tag
	ok    bool
}{
	{"SEQUENCE", lib.Tag{lib.ClassUniversal, 16, true, 0}, true},
	{"long-form:5 SEQUENCE", lib.Tag{lib.ClassUniversal, 16, true, 5}, true},
	{"SEQUENCE CONSTRUCTED", lib.Tag{lib.ClassUniversal, 16, true, 0}, true},
	{"SEQUENCE PRIMITIVE", lib.Tag{lib.ClassUniversal, 16, false, 0}, true},
	{"INTEGER", lib.Tag{lib.ClassUniversal, 2, false, 0}, true},
	{"INTEGER CONSTRUCTED", lib.Tag{lib.ClassUniversal, 2, true, 0}, true},
	{"INTEGER PRIMITIVE", lib.Tag{lib.ClassUniversal, 2, false, 0}, true},
	{"long-form:5 2", lib.Tag{lib.ClassContextSpecific, 2, true, 5}, true},
	{"2 PRIMITIVE", lib.Tag{lib.ClassContextSpecific, 2, false, 0}, true},
	{"APPLICATION 2", lib.Tag{lib.ClassApplication, 2, true, 0}, true},
	{"PRIVATE 2", lib.Tag{lib.ClassPrivate, 2, true, 0}, true},
	{"long-form:5 PRIVATE 2", lib.Tag{lib.ClassPrivate, 2, true, 5}, true},
	{"UNIVERSAL 2", lib.Tag{lib.ClassUniversal, 2, true, 0}, true},
	{"UNIVERSAL 2", lib.Tag{lib.ClassUniversal, 2, true, 0}, true},
	{"UNIVERSAL 2 CONSTRUCTED", lib.Tag{lib.ClassUniversal, 2, true, 0}, true},
	{"UNIVERSAL 2 PRIMITIVE", lib.Tag{lib.ClassUniversal, 2, false, 0}, true},
	{"UNIVERSAL 2 CONSTRUCTED EXTRA", lib.Tag{}, false},
	{"UNIVERSAL 2 EXTRA", lib.Tag{}, false},
	{"UNIVERSAL NOT_A_NUMBER", lib.Tag{}, false},
	{"UNIVERSAL SEQUENCE", lib.Tag{}, false},
	{"UNIVERSAL", lib.Tag{}, false},
	{"SEQUENCE 2", lib.Tag{}, false},
	{"", lib.Tag{}, false},
	{" SEQUENCE", lib.Tag{}, false},
	{"SEQUENCE ", lib.Tag{}, false},
	{"SEQUENCE  CONSTRUCTED", lib.Tag{}, false},
	{"long-form:2", lib.Tag{}, false},
	{"long-form:0 SEQUENCE", lib.Tag{}, false},
	{"long-form:-1 SEQUENCE", lib.Tag{}, false},
	{"long-form:garbage SEQUENCE", lib.Tag{}, false},
}

func TestDecodeTagString(t *testing.T) {
	for i, tt := range decodeTagStringTests {
		tag, err := decodeTagString(tt.input)
		if tag != tt.tag || (err == nil) != tt.ok {
			t.Errorf("%d. decodeTagString(%v) = %v, err=%s, wanted %v, success=%v", i, tt.input, tag, err, tt.tag, tt.ok)
		}
	}
}
