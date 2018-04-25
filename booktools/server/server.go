package server

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	bt "github.com/TheGrum/booktools"
)

type BooktoolsServer struct {
	root *bt.Chunk
}

func (b BooktoolsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request (%v)\n", r.URL.Path)
	p := path.Clean(r.URL.Path)
	elements := strings.Split(strings.TrimPrefix(strings.ToLower(p), "/"), "/")
	if elements == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 Page Not Found")
		return
	}
	action := elements[0]

	log.Printf("Checking possible action: (%v) elements:[%v]", action, elements)
	switch action {
	case "booktools.css":
		SendCSS(w, r)
	case "edit":
		log.Print("edit")
		p = path.Dir(p)
	case "structure":
		log.Print("structure")
		b.SendStructure(w, r)
	case "chaptercharacters":
		log.Print("chaptercharacters")
		b.SendChapterCharacters(w, r)
	case "chaptermatches":
		log.Print("chaptermatches")
		b.SendChapterMatches(w, r)
	case "chapter":
		log.Print("chapter")
		b.SendChapter(w, r)
	default:
		// index.html
		sb := strings.Builder{}
		sb.WriteString(`<head></head><body><h1>`)
		sb.WriteString(b.root.GetFirstSentence())
		sb.WriteString(`</h1></br><a href="structure/">Display Structure</a></p>
			<a href="chaptercharacters/">Display Characters By Chapter</a></p>
			<a href="chaptermatches/add/names/here/">Display Specific Characters By Chapter</a></p>
			<a href="chapter/1/">Chapter 1</a>
			</body>
			`)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write([]byte(sb.String()))
		if err != nil {
			log.Fatalf("Error serving default: %v", err)
		}
	}
}

func (b BooktoolsServer) SendStructure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	sb := strings.Builder{}
	sb.WriteString("<head><link rel=\"stylesheet\" href=\"/booktools.css\"><h1>")
	sb.WriteString(b.root.GetFirstSentence())
	sb.WriteString("</h1></head><body>")

	iter := bt.NewChunkIterator(b.root)

	for iter.NextChunk() != nil {
		if iter.Value().Unit > 0 {
			sb.WriteString(fmt.Sprintf("<p style=\"margin-left: %dpx\">", iter.GetDepth()*50))
			for i := 0; i < iter.GetDepth(); i++ {
				sb.WriteString("    ")
			}
			sb.WriteString("[" + bt.UnitToString(iter.Value().Unit) + "]")
			if iter.Value().Unit == bt.Sentence {
				sb.WriteString(iter.Value().String())
			}
			sb.WriteString("</p>\n")
		}
	}
	sb.WriteString("</body>")
	_, err := w.Write([]byte(sb.String()))
	if err != nil {
		log.Fatalf("Error serving structure: %v", err)
	}
}

func (b BooktoolsServer) SendChapter(w http.ResponseWriter, r *http.Request) {
	p := path.Clean(r.URL.Path)
	target := strings.ToLower(path.Base(p))
	i, err := strconv.Atoi(target)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("<body>Could not parse target %v.</body>", target)))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	sb := strings.Builder{}
	sb.WriteString("<head><link rel=\"stylesheet\" href=\"/booktools.css\"><h1>")
	sb.WriteString(b.root.GetFirstSentence())
	sb.WriteString("</h1></head><body>")
	sb.WriteString(b.root.GetChapterHTML(i))
	sb.WriteString("</body>")
	_, err = w.Write([]byte(sb.String()))
	if err != nil {
		log.Fatalf("Error serving structure: %v", err)
	}
}

func (b BooktoolsServer) SendChapterCharacters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	sb := strings.Builder{}
	sb.WriteString("<head><link rel=\"stylesheet\" href=\"/booktools.css\"><h1>")
	sb.WriteString(b.root.GetFirstSentence())
	sb.WriteString("</h1></head><body><table class=\"simpleTable\">\n")
	sb.WriteString("<tr><td>Chapter</td><td>WordCount</td><td>Characters</td><td>First Sentence</td></tr>\n")
	iter := bt.NewChunkIterator(b.root)

	i := 0
	for iter.NextChunk() != nil {
		if iter.Value().Unit == bt.Chapter {
			i = i + 1
			sb.WriteString("<tr>")
			wc := iter.Value().GetWordCount()
			chars := ""
			charFreq := bt.CharacterFrequencies(iter.Value(), 1, 0)
			pl := bt.RankByFrequency(charFreq)
			for i, p := range pl {
				if i < 8 {
					chars = chars + p.Key + ", "
				}
			}
			chars = strings.TrimSuffix(chars, ", ")
			sb.WriteString(fmt.Sprintf("<td><a href=\"/chapter/%d\">Chapter %d</a></td>", i, i))
			sb.WriteString(fmt.Sprintf("<td>%d</td><td>%v</td><td>%v</td>\n", wc, chars, iter.Value().GetFirstSentence()))
			sb.WriteString("</tr>\n")
		}
	}
	sb.WriteString("</table></body>\n")
	_, err := w.Write([]byte(sb.String()))
	if err != nil {
		log.Fatalf("Error serving chapter characters: %v", err)
	}
}

func (b BooktoolsServer) SendChapterMatches(w http.ResponseWriter, r *http.Request) {
	p := path.Clean(r.URL.Path)
	elements := strings.Split(strings.TrimPrefix(p, "/"), "/")
	if elements == nil || len(elements) < 2 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 Page Not Found")
		return
	}
	elements = elements[1:]
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	sb := strings.Builder{}
	sb.WriteString("<head><link rel=\"stylesheet\" href=\"/booktools.css\"><h1>")
	sb.WriteString(b.root.GetFirstSentence())
	sb.WriteString("</h1></head><body><table class=\"simpleTable\">\n")
	sb.WriteString("<tr><td>Chapter</td>")
	for _, element := range elements {
		sb.WriteString("<td>")
		sb.WriteString(element)
		sb.WriteString("</td>")
	}
	sb.WriteString("</tr>\n")
	iter := bt.NewChunkIterator(b.root)

	i := 0
	for iter.NextChunk() != nil {
		if iter.Value().Unit == bt.Chapter {
			i = i + 1
			sb.WriteString("<tr>")
			sb.WriteString(fmt.Sprintf("<td><a href=\"/chapter/%d\">Chapter %d</a></td>", i, i))
			for _, element := range elements {
				wc := iter.Value().GetSpecificWordCount(element)
				if wc > 0 {
					sb.WriteString("<td>")
					sb.WriteString(strconv.Itoa(wc))
					sb.WriteString(" - " + string([]rune(element)[0:10]))
				} else {
					sb.WriteString("<td style=\"background-color: ffffff;\" class=\"empty\">")
				}
				sb.WriteString("</td>")
			}
			sb.WriteString("</tr>\n")
		}
	}
	sb.WriteString("</table></body>\n")
	_, err := w.Write([]byte(sb.String()))
	if err != nil {
		log.Fatalf("Error serving chapter matches: %v", err)
	}
}

func SendCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte(`
empty td {
	border: none;
	background-color: #FFFFFF;
}
table.simpleTable {
  border: 1px solid #1C6EA4;
  background-color: #DCEEDB;
  width: 100%;
  text-align: left;
  empty-cells: hide;
}
table.simpleTable td, table.simpleTable th {
  border: 1px solid #AAAAAA;
  padding: 3px 2px;
}
table.simpleTable tbody td {
  font-size: 13px;
}
table.simpleTable tr:nth-child(even) {
  background: #B8F5C5;
}
table.simpleTable thead {
  background: #14A44D;
  border-bottom: 2px solid #444444;
}
table.simpleTable thead th {
  font-size: 15px;
  font-weight: bold;
  color: #FFFFFF;
  border-left: 2px solid #D0E4F5;
}
table.simpleTable thead th:first-child {
  border-left: none;
}

table.simpleTable tfoot td {
  font-size: 14px;
}
table.simpleTable tfoot .links {
  text-align: right;
}
table.simpleTable tfoot .links a{
  display: inline-block;
  background: #1C6EA4;
  color: #FFFFFF;
  padding: 2px 8px;
  border-radius: 5px;
}`))
	if err != nil {
		log.Fatalf("Error serving CSS: %v", err)
	}

}

func Listen(root *bt.Chunk, listenPort int) {
	listenOn := fmt.Sprintf(":%d", listenPort)
	http.Handle("/", BooktoolsServer{root: root})
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
