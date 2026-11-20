// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestInt64SetElements(t *testing.T) {
	ctx := context.Background()

	if got := Int64SetElements(ctx, types.SetNull(types.Int64Type)); got != nil {
		t.Errorf("null set: got %v, want nil", got)
	}
	if got := Int64SetElements(ctx, types.SetUnknown(types.Int64Type)); got != nil {
		t.Errorf("unknown set: got %v, want nil", got)
	}
	if got := Int64SetElements(ctx, types.SetValueMust(types.Int64Type, []attr.Value{})); len(got) != 0 {
		t.Errorf("empty set: got %v, want no elements", got)
	}

	set := types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(3), types.Int64Value(1), types.Int64Value(2)})
	got := Int64SetElements(ctx, set)
	slices.Sort(got)
	if want := []int64{1, 2, 3}; !slices.Equal(got, want) {
		t.Errorf("set of 3: got %v, want %v", got, want)
	}
}

func TestDifferenceInt64(t *testing.T) {
	tests := []struct {
		name string
		a, b []int64
		want []int64
	}{
		{"both nil", nil, nil, nil},
		{"a nil", nil, []int64{1}, nil},
		{"b nil", []int64{1, 2}, nil, []int64{1, 2}},
		{"disjoint", []int64{1, 2}, []int64{3, 4}, []int64{1, 2}},
		{"overlap", []int64{1, 2, 3}, []int64{2}, []int64{1, 3}},
		{"b contains all of a", []int64{1, 2}, []int64{2, 1, 3}, nil},
		{"keeps the order of a", []int64{5, 1, 4, 2}, []int64{1}, []int64{5, 4, 2}},
		{"only a's elements, not b's", []int64{1}, []int64{2}, []int64{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DifferenceInt64(tt.a, tt.b); !slices.Equal(got, tt.want) {
				t.Errorf("DifferenceInt64(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
