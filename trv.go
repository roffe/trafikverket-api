package trv

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Settings
var (
	Debug       = false // Be verbose
	IndentChar  = "  "  // The PrettyPrint indentation character(s)
	PrettyPrint = true  // Should we print with indentation using IndentChar
)

const (
	apiURL = "https://api.trafikinfo.trafikverket.se/v2/data.json"
)

// Named Opt variables
const (
	OptName          = "name"
	OptValue         = "value"
	OptObjtype       = "objecttype"
	OptOrderBy       = "orderby"
	OptSchemaversion = "schemaversion"
)

// Tag ...
type Tag struct {
	tag      string // The tag
	level    int    // How deep is the tag, for pretty printing
	children []*Tag // Children tags
	opts     Opts   // Tag options
	value    string // The Value inside the <tag></tag>
	short    bool   // true closes the tag with <TAG/> false closes it with a <TAG></TAG>
	inline   bool   // Should the opening & closing tags be inline ( on the same line )
}

func (t *Tag) String() string {
	return fmt.Sprintf("<%s> <opts: %q> <value: %q> <short: %t> <inline: %t>", t.tag, t.opts, t.value, t.short, t.inline)
}

// Opts holds our tag options <TAG key="value">
type Opts map[string]string

// NewRequest createsa new TRV request
func NewRequest(apiKey string, query *Tag) *Tag {
	t := &Tag{
		tag:      "REQUEST",
		children: []*Tag{Login(apiKey), query},
	}
	return t
}

// Build a TRV query
func (t *Tag) Build(w io.Writer) {
	t.start(w)
	for _, c := range t.children {
		c.level = t.level + 1
		c.Build(w)
	}
	t.end(w)
}

// Do sends the request
func (t *Tag) Do() (*http.Response, error) {
	body := bytes.NewBuffer([]byte{})
	t.Build(body)
	if Debug {
		fmt.Println(body.String())
	}

	resp, err := http.Post(apiURL, "text/xml", body)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != 200 {
		return resp, fmt.Errorf("return code %d", resp.StatusCode)
	}
	return resp, nil
}

// Login tag
func Login(key string) *Tag {
	return &Tag{
		tag: "LOGIN",
		opts: Opts{
			"authenticationkey": key,
		},
		short: true,
	}
}

// Query tag
func Query(opts Opts, tags ...*Tag) *Tag {
	return &Tag{
		tag:      "QUERY",
		opts:     opts,
		children: tags,
	}
}

// Filter tag
func Filter(tags ...*Tag) *Tag {
	return &Tag{
		tag:      "FILTER",
		children: tags,
	}
}

// And tag
func And(tags ...*Tag) *Tag {
	return &Tag{
		tag:      "AND",
		children: tags,
	}
}

// Or tag
func Or(tags ...*Tag) *Tag {
	return &Tag{
		tag:      "OR",
		children: tags,
	}
}

// Eq equal tag
func Eq(opts Opts, tags ...*Tag) *Tag {
	return &Tag{
		tag:      "EQ",
		opts:     opts,
		short:    true,
		children: tags,
	}
}

// Gt greate then tag
func Gt(opts Opts) *Tag {
	return &Tag{
		tag:   "GT",
		opts:  opts,
		short: true,
	}
}

// Lt Lesser then tag
func Lt(opts Opts) *Tag {
	return &Tag{
		tag:   "LT",
		opts:  opts,
		short: true,
	}
}

// Include .tag
func Include(value string) *Tag {
	return &Tag{
		tag:    "INCLUDE",
		value:  value,
		inline: true,
	}
}

// opens the tag
func (t *Tag) start(w io.Writer) {
	if PrettyPrint {
		w.Write([]byte(strings.Repeat(IndentChar, t.level)))
	}

	w.Write([]byte("<" + t.tag))

	if t.opts != nil {
		w.Write([]byte(t.optsString()))
	}

	if t.short {
		w.Write([]byte("/>"))
	} else {
		w.Write([]byte(">"))
	}

	if t.value != "" {
		w.Write([]byte(t.value))
	}

	if t.inline {
		w.Write([]byte("</" + t.tag + ">\n"))
	} else {
		w.Write([]byte("\n"))
	}
}

// closes the tag
func (t *Tag) end(w io.Writer) {
	// If it's a short or inline tag the start() function renders it in whole
	if t.short || t.inline {
		return
	}
	if PrettyPrint {
		w.Write([]byte(strings.Repeat(IndentChar, t.level)))
	}
	w.Write([]byte("</" + t.tag + ">\n"))
}

func (t *Tag) optsString() []byte {
	if t.opts != nil {
		buf := bytes.NewBuffer([]byte{})
		for k, v := range t.opts {
			buf.WriteString(fmt.Sprintf(" %s=\"%s\"", k, v))
		}
		return buf.Bytes()
	}
	return nil
}
