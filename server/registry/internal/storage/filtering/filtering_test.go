// Copyright 2020 Google LLC. All Rights Reserved.
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

package filtering

import (
	"testing"
	"time"
)

func TestFilter_Matches(t *testing.T) {
	tests := []struct {
		desc     string
		filter   string
		fields   []Field
		positive map[string]interface{}
		negative map[string]interface{}
	}{
		{
			desc:   "equal to String",
			filter: `k == "match"`,
			fields: []Field{
				{
					Name: "k",
					Type: String,
				},
			},
			positive: map[string]interface{}{
				"k": "match",
			},
			negative: map[string]interface{}{
				"k": "mismatch",
			},
		},
		{
			desc:   "equal to Int",
			filter: `k == 123`,
			fields: []Field{
				{
					Name: "k",
					Type: Int,
				},
			},
			positive: map[string]interface{}{
				"k": 123,
			},
			negative: map[string]interface{}{
				"k": 321,
			},
		},
		{
			desc:   "less than Timestamp",
			filter: `k < timestamp("2021-01-01T00:00:00Z")`,
			fields: []Field{
				{
					Name: "k",
					Type: Timestamp,
				},
			},
			positive: map[string]interface{}{
				"k": time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			negative: map[string]interface{}{
				"k": time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			desc:   "greater than Timestamp",
			filter: `k > timestamp("2021-01-01T00:00:00Z")`,
			fields: []Field{
				{
					Name: "k",
					Type: Timestamp,
				},
			},
			positive: map[string]interface{}{
				"k": time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			negative: map[string]interface{}{
				"k": time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			desc:   "has StringMap key",
			filter: `has(labels.match)`,
			fields: []Field{
				{
					Name: "labels",
					Type: StringMap,
				},
			},
			positive: map[string]interface{}{
				"labels": map[string]string{
					"match": "v",
				},
			},
			negative: map[string]interface{}{
				"labels": map[string]string{
					"mismatch": "v",
				},
			},
		},
		{
			desc:   "in StringMap keys",
			filter: `"match" in labels`,
			fields: []Field{
				{
					Name: "labels",
					Type: StringMap,
				},
			},
			positive: map[string]interface{}{
				"labels": map[string]string{
					"match": "v",
				},
			},
			negative: map[string]interface{}{
				"labels": map[string]string{
					"mismatch": "v",
				},
			},
		},
		{
			desc:   "equal to StringMap value",
			filter: `labels["k"] == "match"`,
			fields: []Field{
				{
					Name: "labels",
					Type: StringMap,
				},
			},
			positive: map[string]interface{}{
				"labels": map[string]string{
					"k": "match",
				},
			},
			negative: map[string]interface{}{
				"labels": map[string]string{
					"k": "mismatch",
				},
			},
		},
		{
			desc:   "substring of StringMap value",
			filter: `labels.k.contains("substring")`,
			fields: []Field{
				{
					Name: "labels",
					Type: StringMap,
				},
			},
			positive: map[string]interface{}{
				"labels": map[string]string{
					"k": "substring_match",
				},
			},
			negative: map[string]interface{}{
				"labels": map[string]string{
					"k": "substr_mismatch",
				},
			},
		},
		{
			desc:   "in StringMap value split",
			filter: `"match" in labels.k.split("_")`,
			fields: []Field{
				{
					Name: "labels",
					Type: StringMap,
				},
			},
			positive: map[string]interface{}{
				"labels": map[string]string{
					"k": "split_match",
				},
			},
			negative: map[string]interface{}{
				"labels": map[string]string{
					"k": "split_mismatch",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			f, err := NewFilter(test.filter, test.fields)
			if err != nil {
				t.Fatalf("NewFilter(%q) returned error: %s", test.filter, err)
			}

			if match, err := f.Matches(test.positive); err != nil {
				t.Fatalf("NewFilter(%q).Matches(%v) returned error: %s", test.filter, test.positive, err)
			} else if !match {
				t.Errorf("NewFilter(%q).Matches(%v) returned unexpected mismatch", test.filter, test.positive)
			}

			if match, err := f.Matches(test.negative); err != nil {
				t.Fatalf("NewFilter(%q).Matches(%v) returned error: %s", test.filter, test.negative, err)
			} else if match {
				t.Errorf("NewFilter(%q).Matches(%v) returned unexpected match", test.filter, test.negative)
			}
		})
	}
}
