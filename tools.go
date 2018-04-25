// Copyright © 2018 Howard C. Shaw III <howardcshaw@gmail.com>
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

package booktools

import (
	"fmt"
	"sort"
)

func PrintLines(ss []string) {
	for _, s := range ss {
		fmt.Println(s)
	}
}

func SortAndPrintLines(ss []string) {
	sort.Slice(ss, func(i, j int) bool {
		if ss[i] < ss[j] {
			return true
		}
		return false
	})
	for _, s := range ss {
		fmt.Println(s)
	}
}

func PrintNameFrequency(nf map[string]int) {
	names := make([]string, 0, len(nf))
	for k, _ := range nf {
		names = append(names, k)
	}
	sort.Slice(names, func(i, j int) bool {
		if names[i] < names[j] {
			return true
		}
		return false
	})
	for _, k := range names {
		fmt.Printf("%v: %-10d\n", k, nf[k])
	}
}
