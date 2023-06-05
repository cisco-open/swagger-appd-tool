/*
	SPDX-License-Identifier: Apache-2.0

	Copyright (c) 2023 Cisco Systems, Inc.

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

package main

import "fmt"

func testCompare() {
	sr := []*SwaggerRule{
		{
			Method:   "GET",
			UriRegex: "/a/b/.*/d",
		},
		{
			Method:   "POST",
			UriRegex: "/a/b/.*/d",
		},
		{
			Method:   "GET",
			UriRegex: "/a/e/.*/d",
		},
		{
			Method:   "POST",
			UriRegex: "/a/e/.*/d",
		},
		{
			Method:   "POST",
			UriRegex: "/a/g/.*/d",
		}}

	ar := []AppDRule{
		{
			Method:   "GET",
			UriRegex: "/a/a/.*/d",
		},
		{
			Method:   "POST",
			UriRegex: "/a/b/.*/d",
		},
		{
			Method:   "GET",
			UriRegex: "/a/e/.*/d",
		},
		{
			Method:   "POST",
			UriRegex: "/a/f/.*/d",
		},
	}

	sra, ard := diffRules(sr, &ar)

	fmt.Printf("Diff:\nAdd: %v\nDelete:%v\n", sra, ard)
}
