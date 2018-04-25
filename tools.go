package booktools

import (
	"fmt"
	"sort"

	"github.com/jdkato/prose/chunk"
	"github.com/jdkato/prose/tag"
	"github.com/jdkato/prose/tokenize"
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

func GetNamesViaProse(s string) []string {
	words := tokenize.TextToWords(s)
	regex := chunk.TreebankNamedEntities

	tagger := tag.NewPerceptronTagger()
	prosenames := make(map[string]int)
	for _, entity := range chunk.Chunk(tagger.Tag(words), regex) {
		prosenames[entity] = 1
	}
	names := make([]string, 0, len(prosenames))
	for k, _ := range prosenames {
		names = append(names, k)
	}
	return names
}

func GetNameFrequencyViaProse(s string) map[string]int {
	words := tokenize.TextToWords(s)
	regex := chunk.TreebankNamedEntities

	tagger := tag.NewPerceptronTagger()
	prosenames := make(map[string]int)
	for _, entity := range chunk.Chunk(tagger.Tag(words), regex) {
		prosenames[entity] = prosenames[entity] + 1
	}
	return prosenames
}

func NameFrequencyMapMerge(a map[string]int, b map[string]int) {
	for k, v := range b {
		a[k] = a[k] + v
	}
}
