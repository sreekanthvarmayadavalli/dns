/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	// valid
	for _, testCase := range []Config{
		{Federations: map[string]string{}},
		{Federations: map[string]string{"abc": "d.e.f"}},
		{StubDomains: map[string][]string{}},
		{StubDomains: map[string][]string{"foo.com": []string{"1.2.3.4"}}},
		{StubDomains: map[string][]string{"foo.com": []string{"1.2.3.4:32564"}}},
		{StubDomains: map[string][]string{"foo.com": []string{"ns.foo.com"}}},
		{StubDomains: map[string][]string{
			"foo.com": []string{"ns.foo.com"},
			"bar.com": []string{"1.2.3.4"},
		}},
		{StubDomains: map[string][]string{
			"acme.local":   {"2001:db8:1:1:1::1"},
			"google.local": {"google-public-dns-a.google.com"},
			"widget.local": {"[2001:db8:2:2:2::2]:10053", "2001:db8:3:3:3::3"},
		}},
		{UpstreamNameservers: []string{}},
		{UpstreamNameservers: []string{"1.2.3.4"}},
		{UpstreamNameservers: []string{"1.2.3.4", "8.8.4.4", "8.8.8.8"}},
		{UpstreamNameservers: []string{"1.2.3.4:53"}},
		{UpstreamNameservers: []string{"[2001:db8:2:2:2::2]:10053", "2001:db8:3:3:3::3"}},
	} {
		err := testCase.Validate()
		assert.Nil(t, err, "should be valid: %+v", testCase)
	}

	// invalid
	for _, testCase := range []Config{
		{Federations: map[string]string{"a.b": "cdef"}},
		{StubDomains: map[string][]string{"": []string{"1.2.3.4"}}},
		{StubDomains: map[string][]string{"$$$$": []string{"1.2.3.4"}}},
		{StubDomains: map[string][]string{"foo": []string{"$$$$"}}},
		{StubDomains: map[string][]string{"foo.com": []string{"1.2.3.4:65564"}}},
		{UpstreamNameservers: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3", "4.4.4.4"}},
		{UpstreamNameservers: []string{"1.1.1.1:abc", "1.1.1.1:", "1.1.1.1:123456789"}},
	} {
		err := testCase.Validate()
		assert.NotNil(t, err, "should not be valid: %+v", testCase)
	}
}

func TestValidateNodeLocalCacheConfig(t *testing.T) {
	for _, tc := range []struct {
		name      string
		config    Config
		wantError bool
	}{
		{
			name:   "empty config",
			config: Config{},
		},
		{
			name: "valid config",
			config: Config{
				StubDomains:         map[string][]string{"": {"1.1.1.1"}},
				UpstreamNameservers: []string{"1.1.1.1", "2.2.2.2"},
			},
		},
		{
			name: "invalid config: FQDN stub domain",
			config: Config{
				StubDomains:         map[string][]string{"": {"cluster.local"}},
				UpstreamNameservers: []string{"1.1.1.1", "2.2.2.2"},
			},
			wantError: true,
		},
		{
			name: "invalid config: FQDN upstream name server",
			config: Config{
				StubDomains:         map[string][]string{"": {"1.1.1.1"}},
				UpstreamNameservers: []string{"1.1.1.1", "cluster.local"},
			},
			wantError: true,
		},
	} {
		err := tc.config.ValidateNodeLocalCacheConfig()
		gotError := err != nil
		if gotError != tc.wantError {
			t.Fatalf("ValidateNodeLocalCacheConfig(%v) got error = %v, want = %v", tc.config, gotError, tc.wantError)
		}
	}
}
