package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var defaultHandlerTmpl = `
	<html>
		<head>
			<meta charset="utf-8" />
			<title>Choose your own adventure</title>
		</head>
		<body>
			<h1>{{.Title}}</h1>
			{{range .Paragraphs}}
				<p>{{.}}</p>
			{{end}}

			<ul>
				{{range .Options}}
					<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
				{{end}}
			</ul>
		</body>
	</html>
`

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tpl.Execute(w, h.s["intro"])
	if err != nil {
		panic(err)
	}

}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

// Story is ...
type Story map[string]Chapter

// Chapter was generated with https://mholt.github.io/json-to-go/
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option is ...
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
