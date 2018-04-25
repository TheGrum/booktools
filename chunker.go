package booktools

import (
	"bytes"
	"io"
)

const (
	Word      = iota
	Sentence  = iota
	Paragraph = iota
	Section   = iota
	Chapter   = iota
	Work      = iota
)

func UnitToString(unit int) string {
	switch unit {
	case Word:
		return "Word"
	case Sentence:
		return "Sentence"
	case Paragraph:
		return "Paragraph"
	case Section:
		return "Section"
	case Chapter:
		return "Chapter"
	case Work:
		return "Work"
	}
	return "Unknown"
}

type Chunk struct {
	Position int64
	Length   int64
	Unit     int
	Word     string
	Children []*Chunk
}

type Chunker struct {
	r        io.Reader
	position int64

	curWord       string
	curSentence   string
	lastWord      int64
	lastSentence  int64
	lastParagraph int64
	lastSection   int64
	lastChapter   int64
	lastRune      rune

	out chan *Chunk
	b   bytes.Buffer

	OnSentence       func(c *Chunker, s string)
	OnBeforeSentence func(c *Chunker, s string)
}

// NewChunker creates a Chunker to chunk the contents of io.Reader
// into channel out
func NewChunker(r io.Reader, out chan *Chunk) *Chunker {
	return &Chunker{
		r:             r,
		position:      0,
		lastWord:      0,
		lastSentence:  0,
		lastParagraph: 0,
		lastChapter:   0,
		lastRune:      '|',

		out: out,
	}
}

func (c *Chunker) Read(p []byte) (n int, err error) {
	n, err = c.r.Read(p)
	if !(err == nil) {
		close(c.out)
		//fmt.Println(err)
		return 0, err
	}
	c.b.Write(p)
	c.process()
	return n, nil
}

func (c *Chunker) Word() {
	if c.lastWord == c.position || c.curWord == "" {
		return
	}
	c.out <- &Chunk{Position: c.lastWord, Length: c.position - c.lastWord, Unit: Word, Word: c.curWord}
	c.lastWord = c.position
}

func (c *Chunker) Sentence() {
	if c.lastSentence == c.position {
		return
	}
	if c.OnBeforeSentence != nil {
		c.OnBeforeSentence(c, c.curSentence)
	}
	c.Word()
	c.out <- &Chunk{Position: c.lastSentence, Length: c.position - c.lastSentence, Unit: Sentence}
	c.lastSentence = c.position
	if c.OnSentence != nil {
		c.OnSentence(c, c.curSentence)
	}
	c.curSentence = ""
}

func (c *Chunker) Paragraph() {
	if c.lastParagraph == c.position {
		return
	}
	c.Sentence()
	c.out <- &Chunk{Position: c.lastParagraph, Length: c.position - c.lastParagraph, Unit: Paragraph}
	c.lastParagraph = c.position
}

func (c *Chunker) Section() {
	if c.lastSection == c.position {
		return
	}
	//c.Paragraph()
	c.out <- &Chunk{Position: c.lastSection, Length: c.position - c.lastSection, Unit: Section}
	c.lastSection = c.position
}

func (c *Chunker) Chapter() {
	if c.lastChapter == c.position {
		return
	}
	if c.lastSection != c.position {
		c.Section()
	}
	c.out <- &Chunk{Position: c.lastChapter, Length: c.position - c.lastChapter, Unit: Chapter}
	c.lastChapter = c.position
}

func (c *Chunker) process() {
	var err error
	for err == nil {
		r, size, err := c.b.ReadRune()
		if !(err == nil) {
			//fmt.Printf("Error, size: %v, %v\n", err, size)
			c.position += int64(size)
			//c.chapter()
			return
		}
		c.position += int64(size)
		switch r {
		case ' ', '\r', '\n':
//			if c.lastRune != '#' {
				switch c.curWord {
				case "---":
					// Section marker
					c.Section()
				case "# Part":
					// Chapter marker
					c.Chapter()
				default:
					switch c.lastRune {
					case '.', '!', '?':
						//fmt.Printf("Sentence: [%v] %v\n", c.curWord, r)
						c.Sentence()
					case '\r', '\n':
						//fmt.Println("Para")
						c.Paragraph()
					default:
						if c.curWord != "" {
							c.curSentence = c.curSentence + " " + c.curWord
							c.Word()
						}
					}
					//fmt.Printf("|%v| [%v][%v]\n", c.curWord, c.lastRune, r)
				}
				c.curWord = ""
			//}
		default:
			c.curWord = c.curWord + string(r)
		}
		c.lastRune = r
	}
}
