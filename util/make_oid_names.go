// Copyright 2023 The DER ASCII Authors. All Rights Reserved.
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

//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"encoding/asn1"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	oidNamesTxt = "util/oid_names.txt"
	oidNamesGo  = "cmd/der2ascii/oid_names.go"
)

func makeOIDNames() error {
	inp, err := os.Open(oidNamesTxt)
	if err != nil {
		return err
	}
	defer inp.Close()

	var b bytes.Buffer
	b.WriteString(`// Copyright 2016 The DER ASCII Authors. All Rights Reserved.
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

// This file is generated by make_oid_names.go. Do not edit by hand.
// To regenerate, run "go run util/make_oid_names.go" from the top-level directory.

package main

var oidNames = []struct {
	oid  []byte
	name string
}{
`)

	var lineNo int
	scanner := bufio.NewScanner(inp)
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		if idx := strings.IndexByte(line, '#'); idx >= 0 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		colon := strings.IndexByte(line, ':')
		if colon < 0 {
			return fmt.Errorf("could not parse line %d: missing colon separator", lineNo)
		}

		name := strings.TrimSpace(line[:colon])
		oidStr := strings.Split(strings.TrimSpace(line[colon+1:]), ".")

		oid := make(asn1.ObjectIdentifier, len(oidStr))
		for i, s := range oidStr {
			oid[i], err = strconv.Atoi(s)
			if err != nil || oid[i] < 0 {
				return fmt.Errorf("could not parse line %d: invalid OID component %q", lineNo, s)
			}
		}

		// Encode the OID, excluding the outer wrapper.
		der, err := asn1.Marshal(oid)
		if err != nil {
			return fmt.Errorf("could not parse line %d: error encoding OID: %s", lineNo, err)
		}

		// Parse the OID back out to remove the header.
		var raw asn1.RawValue
		rest, err := asn1.Unmarshal(der, &raw)
		if err != nil || len(rest) != 0 {
			// This should be impossible.
			panic("could not reparse asn1.Marshal output")
		}

		fmt.Fprintf(&b, "\t{%#v, %q},\n", raw.Bytes, name)
	}
	b.WriteString("}\n")

	return os.WriteFile(oidNamesGo, b.Bytes(), 0666)
}

func main() {
	if err := makeOIDNames(); err != nil {
		fmt.Fprintf(os.Stderr, "Error making oid_names.go: %s\n", err)
		os.Exit(1)
	}
}