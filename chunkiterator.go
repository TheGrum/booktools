package booktools

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type ChunkIterator struct {
	root    *Chunk
	stack   []*Chunk
	indices []int
	depth   int
	current *Chunk

	FirstWordInSentence    bool
	FirstSentenceInSection bool
}

func NewChunkIterator(root *Chunk) *ChunkIterator {
	return &ChunkIterator{
		root:    root,
		stack:   make([]*Chunk, 0, 10),
		indices: make([]int, 0, 10),
		depth:   -1,
		current: root,
	}
}

func (c *ChunkIterator) GetDepth() int {
	return c.depth
}

func (c *ChunkIterator) NextChunk() *Chunk {
	//	fmt.Println("Next chunk")
	if c.root == nil {
		return nil
	}
	if c.current == nil {
		c.current = c.root
		//fmt.Println("Root")
		return c.current
	}
	if c.current.Children != nil {
		if len(c.current.Children) > 0 {
			//fmt.Println("Diving into children")
			c.depth = c.depth + 1
			c.stack = append(c.stack, c.current)
			c.indices = append(c.indices, 0)
			//fmt.Printf("%v, %v, %v", c.stack, c.indices, c.current)
			c.current = c.current.Children[0]
			return c.current
		}
	}
	for c.depth > -1 {
		//fmt.Printf("::%v %v %v\n", c.depth, c.indices, c.stack)
		// we are in a leaf node, go to next sibling
		c.indices[c.depth] = c.indices[c.depth] + 1
		//fmt.Printf("%v %v\n", c.indices[c.depth], c.stack[c.depth])
		if c.indices[c.depth] >= len(c.stack[c.depth].Children) {
			//fmt.Println("Pop stack")
			//fmt.Printf(":BEFORE:%v %v %v\n", c.depth, c.indices, c.stack)
			// We are out of siblings, pop the stack and try again
			c.depth = c.depth - 1
			if c.depth == -1 {
				//fmt.Println("Out of stack")
				// we are out of stack, so we are done
				c.current = nil
				return nil
			}
			c.stack = c.stack[0 : c.depth+1]
			c.indices = c.indices[0 : c.depth+1]
			//fmt.Printf(":AFTER:%v %v %v\n", c.depth, c.indices, c.stack)
		} else {
			c.current = c.stack[c.depth].Children[c.indices[c.depth]]
			return c.current
		}
	}
	return nil
}

func (c *ChunkIterator) NextWord() *Chunk {
	//fmt.Println("Next word")
	if c.root == nil {
		return nil
	}
	if c.current == nil {
		c.current = c.root
	}
	for true {
		c.NextChunk()
		//fmt.Println(c.current)
		if c.current == nil {
			return nil
		}
		if c.current.Unit == Word {
			if c.depth > -1 && c.indices[c.depth] == 0 {
				c.FirstWordInSentence = true
			} else {
				c.FirstWordInSentence = false
			}
			if c.depth > 0 && c.indices[c.depth-1] == 0 {
				c.FirstSentenceInSection = true
			} else {
				c.FirstSentenceInSection = false
			}
			return c.current
		}
	}
	return nil
}

func (c *ChunkIterator) Value() *Chunk {
	return c.current
}

func (c *Chunk) String() string {
	if c.Children == nil {
		return c.Word
	}
	iter := NewChunkIterator(c)
	sb := strings.Builder{}
	for iter.NextChunk() != nil {
		switch iter.Value().Unit {
		case Word:
			sb.WriteString(iter.Value().Word + " ")
		case Sentence:
			sb.WriteString(" ")
		case Paragraph:
			sb.WriteString("\n\n")
		case Section:
			sb.WriteString("\n\n----\n\n")
		case Chapter:
			sb.WriteString("\nChapter :\n")
		}
	}
	return sb.String()
}

func (c *Chunk) HTML() string {
	if c.Children == nil {
		return c.Word
	}
	iter := NewChunkIterator(c)
	sb := strings.Builder{}
	sb.WriteString("<p>")
	for iter.NextChunk() != nil {
		switch iter.Value().Unit {
		case Word:
			sb.WriteString(iter.Value().Word + " ")
		case Sentence:
			sb.WriteString(" ")
		case Paragraph:
			sb.WriteString("\n</p><p>\n")
		case Section:
			sb.WriteString("\n</p><p>\n----\n</p><p>\n")
		case Chapter:
			sb.WriteString("\n</br>Chapter :\n</br>")
		}
	}
	sb.WriteString("</p>")
	return sb.String()
}

func (c *Chunk) PrintStructure(includeSentences bool, maxDepth int) string {
	if c.Children == nil {
		return c.Word
	}
	iter := NewChunkIterator(c)
	sb := strings.Builder{}
	for iter.NextChunk() != nil {
		if iter.Value().Unit > 0 && iter.GetDepth() < maxDepth {
			for i := 0; i < iter.GetDepth(); i++ {
				sb.WriteString("    ")
			}
			sb.WriteString("[" + UnitToString(iter.Value().Unit) + "]")
			if includeSentences && iter.Value().Unit == Sentence {
				sb.WriteString(iter.Value().String())
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func IdentifyCharacters(root *Chunk, minAppearance int, minNonFirst int) []string {
	var incidence = make(map[string]int)
	var nonfirst = make(map[string]int)
	var confirmed = make([]string, 0)

	name := ""

	//	fmt.Println("New iterator")
	iter := NewChunkIterator(root)
	for iter.NextWord() != nil {
		k := iter.Value()
		w := k.Word
		w = strings.Replace(w, "\"", "", -1)
		w = strings.TrimSuffix(w, "'ve")
		w = strings.TrimSuffix(w, "'re")
		w = strings.TrimSuffix(w, "'d")
		w = strings.TrimSuffix(w, "'ll")
		w = strings.TrimSuffix(w, "'s")
		w = strings.Replace(w, "'", "", -1)
		w = strings.Replace(w, ",", "", -1)
		w = strings.Replace(w, "?", "", -1)
		w = strings.Replace(w, "!", "", -1)
		w = strings.Replace(w, ".", "", -1)
		c, _ := utf8.DecodeRuneInString(w)
		if unicode.IsUpper(c) {
			//fmt.Printf("%v\n", k)
			incidence[w] = incidence[w] + 1
			if !iter.FirstWordInSentence {
				nonfirst[w] = nonfirst[w] + 1
			}
			if name == "" {
				name = w
			} else {
				name = name + " " + w
			}
			if name != k.Word {
				incidence[name] = incidence[name] + 1
			}
		} else {
			name = ""
		}
		//fmt.Printf("%v [%v]\n", w, name)
	}

	for k, v := range incidence {
		if v > minAppearance && nonfirst[k] > minNonFirst {
			// If we have not seen a name at least X times
			// It is probably not a significant character
			// If it is only ever capitalized as the first
			// word in a sentence, probably not a significant
			// character
			confirmed = append(confirmed, k)
		}
	}

	return confirmed
}

func CharacterFrequencies(root *Chunk, minAppearance int, minNonFirst int) map[string]int {
	var incidence = make(map[string]int)
	var nonfirst = make(map[string]int)
	var confirmed = make(map[string]int)

	name := ""

	//	fmt.Println("New iterator")
	iter := NewChunkIterator(root)
	for iter.NextWord() != nil {
		k := iter.Value()
		w := k.Word
		w = strings.Replace(w, "\"", "", -1)
		w = strings.TrimSuffix(w, "'ve")
		w = strings.TrimSuffix(w, "'re")
		w = strings.TrimSuffix(w, "'d")
		w = strings.TrimSuffix(w, "'ll")
		w = strings.TrimSuffix(w, "'s")
		w = strings.Replace(w, "'", "", -1)
		w = strings.Replace(w, ",", "", -1)
		w = strings.Replace(w, "?", "", -1)
		w = strings.Replace(w, "!", "", -1)
		w = strings.Replace(w, ".", "", -1)
		c, _ := utf8.DecodeRuneInString(w)
		if unicode.IsUpper(c) {
			//fmt.Printf("%v\n", k)
			incidence[w] = incidence[w] + 1
			if !iter.FirstWordInSentence {
				nonfirst[w] = nonfirst[w] + 1
			}
			if name == "" {
				name = w
			} else {
				name = name + " " + w
			}
			if name != k.Word {
				incidence[name] = incidence[name] + 1
			}
		} else {
			name = ""
		}
		//fmt.Printf("%v [%v]\n", w, name)
	}

	for k, v := range incidence {
		if v > minAppearance && nonfirst[k] > minNonFirst {
			// If we have not seen a name at least X times
			// It is probably not a significant character
			// If it is only ever capitalized as the first
			// word in a sentence, probably not a significant
			// character
			confirmed[k] = v
			//			confirmed = append(confirmed, k)
		}
	}

	return confirmed
}

func (c *Chunk) GetWordCount() int {
	iter := NewChunkIterator(c)
	wc := 0
	for iter.NextChunk() != nil {
		if iter.Value().Unit == Word {
			wc = wc + 1
		}
	}
	return wc
}

func (c *Chunk) GetSpecificWordCount(w string) int {
	iter := NewChunkIterator(c)
	wc := 0
	if strings.Contains(w, " ") {
		for iter.NextChunk() != nil {
			if iter.Value().Unit == Paragraph {
				s := iter.Value().String()
				wc = wc + strings.Count(s, w)
			}
		}
		return wc
	}
	for iter.NextChunk() != nil {
		if iter.Value().Unit == Word && strings.Trim(iter.Value().Word, " ") == strings.Trim(w, " ") {
			wc = wc + 1
		}
	}
	return wc
}

func (c *Chunk) GetFirstSentence() string {
	iter := NewChunkIterator(c)
	for iter.NextChunk() != nil {
		if iter.Value().Unit == Sentence {
			return iter.Value().String()
		}
	}
	return ""
}

func (c *Chunk) GetNthSentence(n int) string {
	iter := NewChunkIterator(c)
	i := 0
	for iter.NextChunk() != nil {
		if iter.Value().Unit == Sentence {
			i = i + 1
			if i == n {
				return iter.Value().String()
			}
		}
	}
	return ""
}

func (c *Chunk) GetChapter(chapter int) string {
	if c.Children == nil {
		return ""
	}
	iter := NewChunkIterator(c)
	i := 0
	for iter.NextChunk() != nil {
		if iter.Value().Unit == Chapter {
			i = i + 1
			if i == chapter {
				return iter.Value().String()
			}
		}
	}
	return ""
}

func (c *Chunk) GetChapterHTML(chapter int) string {
	if c.Children == nil {
		return ""
	}
	iter := NewChunkIterator(c)
	i := 0
	for iter.NextChunk() != nil {
		if iter.Value().Unit == Chapter {
			i = i + 1
			if i == chapter {
				return iter.Value().HTML()
			}
		}
	}
	return ""
}

func (c *Chunk) PrintTopXCharactersPerChapter(includeSentences bool, includeXthSentence int, topX int, tabDelimit bool, wordCount bool) {
	if c.Children == nil {
		return
	}
	iter := NewChunkIterator(c)
	i := 0
	var sent string
	var wc int
	for iter.NextChunk() != nil {
		if iter.Value().Unit == Chapter {
			i = i + 1
			sent = ""
			wc = 0
			if includeSentences {
				sent = "[" + iter.Value().GetFirstSentence() + "] "
			}
			if includeXthSentence > -1 {
				sent = sent + " " + "[" + iter.Value().GetNthSentence(includeXthSentence) + "] "
			}
			if wordCount {
				wc = iter.Value().GetWordCount()
			}
			chars := ""
			charFreq := CharacterFrequencies(iter.Value(), 1, 0)
			pl := RankByFrequency(charFreq)
			for i, p := range pl {
				if i < topX {
					chars = chars + p.Key + ","
				}
			}
			chars = strings.TrimSuffix(chars, ",")
			if wordCount {
				if tabDelimit {
					fmt.Printf("%03d\t%03d\t%v\t%v\n", i, wc, sent, chars)
				} else {
					fmt.Printf("%03d: [%03d] %v%v\n", i, wc, sent, chars)
				}
			} else {
				if tabDelimit {
					fmt.Printf("%03d\t%v\t%v\n", i, sent, chars)
				} else {
					fmt.Printf("%03d: %v%v\n", i, sent, chars)
				}

			}

		}
	}
}

func (c *Chunk) PrintChapter(chapter int) {
	fmt.Println(c.GetChapter(chapter))
}

func DigestChunks(c chan *Chunk, out chan *Chunk) {
	root := &Chunk{Position: 0, Length: 0, Unit: Work, Word: "", Children: make([]*Chunk, 0)}
	var words = make([]*Chunk, 0)
	var sentences = make([]*Chunk, 0)
	var paragraphs = make([]*Chunk, 0)
	var sections = make([]*Chunk, 0)
	var chapters = make([]*Chunk, 0)
	for ch := range c {
		switch ch.Unit {
		case Word:
			words = append(words, ch)
		case Sentence:
			if len(words) > 0 {
				ch.Children = words
				words = make([]*Chunk, 0)
				sentences = append(sentences, ch)
			}
		case Paragraph:
			if len(sentences) > 0 {
				ch.Children = sentences
				sentences = make([]*Chunk, 0)
				paragraphs = append(paragraphs, ch)
			}
		case Section:
			if len(paragraphs) > 0 {
				ch.Children = paragraphs
				paragraphs = make([]*Chunk, 0)
				sections = append(sections, ch)
			}
		case Chapter:
			if len(sections) > 0 {
				ch.Children = sections
				sections = make([]*Chunk, 0)
				chapters = append(chapters, ch)
			}
		}
	}

	if len(words) > 0 {
		ch := &Chunk{Position: -1, Length: -1, Unit: Sentence, Word: "", Children: make([]*Chunk, 0)}
		ch.Children = words
		words = make([]*Chunk, 0)
		sentences = append(sentences, ch)
	}
	if len(sentences) > 0 {
		ch := &Chunk{Position: -1, Length: -1, Unit: Paragraph, Word: "", Children: make([]*Chunk, 0)}
		ch.Children = sentences
		sentences = make([]*Chunk, 0)
		paragraphs = append(paragraphs, ch)
	}
	if len(paragraphs) > 0 {
		ch := &Chunk{Position: -1, Length: -1, Unit: Section, Word: "", Children: make([]*Chunk, 0)}
		ch.Children = paragraphs
		paragraphs = make([]*Chunk, 0)
		sections = append(sections, ch)
	}
	if len(sections) > 0 {
		ch := &Chunk{Position: -1, Length: -1, Unit: Chapter, Word: "", Children: make([]*Chunk, 0)}
		ch.Children = sections
		sections = make([]*Chunk, 0)
		chapters = append(chapters, ch)
	}

	root.Children = chapters
	out <- root
}

func RankByFrequency(freq map[string]int) PairList {
	pl := make(PairList, len(freq))
	i := 0
	for k, v := range freq {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
