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

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/netascode/go-secureaccess"
	"github.com/tidwall/gjson"
)

// exceptionApplicationListPageSize is the largest `limit` that the routes of the exception application lists accept (1001 is a 400).
// Without a `limit` a page holds only 100 items (seen on the applications of a list).
const exceptionApplicationListPageSize = 1000

// exceptionApplicationListMaxPages stops a read that never reaches the total, for example because every page holds one item. It is a variable so that
// a test does not have to make a thousand requests.
var exceptionApplicationListMaxPages = 1000

// readExceptionApplicationListPages reads every page of a collection route of the exception application lists API and returns the items of all pages
// as one JSON array. The routes answer {"status": ..., "meta": {"page", "limit", "total"}, "data": [...]}, the pages start at 1 and are selected with
// `page` and `limit`. The client does not handle them: the path prefix /policies/v2/ is registered with a paginator for a bare array with `offset`,
// and the API reads `offset` as the index of a page.
//
// The read fails instead of returning a part of the collection: every page must have a `data` array and a numeric `meta.total`, the total must not change
// while reading, a page must not hold more than `limit` items, be empty before the total is reached or repeat the page before, and the items read must be
// the total.
func readExceptionApplicationListPages(ctx context.Context, client *secureaccess.Client, path string, reqMods ...func(*secureaccess.Req)) (string, error) {
	var items []string
	var previous string
	var total int64
	for page := 1; ; page++ {
		if page > exceptionApplicationListMaxPages {
			return "", fmt.Errorf("read %d items in %d pages and the total of %d was not reached", len(items), page-1, total)
		}
		res, err := client.Get(ctx, fmt.Sprintf("%s?page=%d&limit=%d", path, page, exceptionApplicationListPageSize), reqMods...)
		if err != nil {
			// the answer of the API says what is wrong, the client returns it next to the error
			return "", fmt.Errorf("page %d: %w, %s", page, err, res.String())
		}
		data := res.Get("data")
		if !data.IsArray() {
			return "", fmt.Errorf("page %d: expected a list in `data`, got: %s", page, data.Type)
		}
		metaTotal := res.Get("meta.total")
		if metaTotal.Type != gjson.Number {
			return "", fmt.Errorf("page %d: expected a number in `meta.total`, got: %s", page, metaTotal.Type)
		}
		if page == 1 {
			total = metaTotal.Int()
			// a total that the page limit can never reach would only end at the page limit, after a thousand requests that the rate limit of the API refuses
			if total > int64(exceptionApplicationListMaxPages)*exceptionApplicationListPageSize {
				return "", fmt.Errorf("page 1: `meta.total` is %d, more than %d pages of %d items", total, exceptionApplicationListMaxPages, exceptionApplicationListPageSize)
			}
		} else if metaTotal.Int() != total {
			return "", fmt.Errorf("page %d: `meta.total` changed from %d to %d while the pages were read", page, total, metaTotal.Int())
		}
		var pageItems []string
		data.ForEach(func(_, item gjson.Result) bool {
			pageItems = append(pageItems, item.Raw)
			return true
		})
		n := len(pageItems)
		if n > exceptionApplicationListPageSize {
			return "", fmt.Errorf("page %d: %d items, more than the limit of %d", page, n, exceptionApplicationListPageSize)
		}
		// an API that ignores `page` answers the first page again, and the counts can still add up to the total
		if joined := strings.Join(pageItems, ","); n > 0 && joined == previous {
			return "", fmt.Errorf("page %d is identical to page %d, the API ignores `page`", page, page-1)
		} else {
			previous = joined
		}
		items = append(items, pageItems...)
		if int64(len(items)) >= total {
			break
		}
		if n == 0 {
			return "", fmt.Errorf("page %d is empty, but only %d of %d items were read", page, len(items), total)
		}
	}
	if int64(len(items)) != total {
		return "", fmt.Errorf("read %d items, but `meta.total` is %d", len(items), total)
	}
	return "[" + strings.Join(items, ",") + "]", nil
}
