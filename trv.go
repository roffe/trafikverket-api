package trv

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	OptObjtype       = "objecttype"
	OptOrderBy       = "orderby"
	OptRadius        = "radius"
	OptSchemaversion = "schemaversion"
	OptShape         = "shape"
	OptLimit         = "limit"
	OptValue         = "value"
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
func (t *Tag) Do() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	t.Build(buf)
	if Debug {
		fmt.Println(buf.String())
	}

	resp, err := http.Post(apiURL, "text/xml", buf)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("return code %d", resp.StatusCode)
	}

	if err != nil {
		log.Print(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
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

// Exists Exists
func Exists(opts Opts, tags ...*Tag) *Tag {
	return &Tag{
		tag:      "EXISTS",
		opts:     opts,
		short:    true,
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

// Gt Greater Than
func Gt(opts Opts) *Tag {
	return &Tag{
		tag:   "GT",
		opts:  opts,
		short: true,
	}
}

// Gte Greater Than or Equal
func Gte(opts Opts) *Tag {
	return &Tag{
		tag:   "GTE",
		opts:  opts,
		short: true,
	}
}

// Lt Less Than
func Lt(opts Opts) *Tag {
	return &Tag{
		tag:   "LT",
		opts:  opts,
		short: true,
	}
}

// Lte Less Than or Equals
func Lte(opts Opts) *Tag {
	return &Tag{
		tag:   "LTE",
		opts:  opts,
		short: true,
	}
}

// Ne Not Equal
func Ne(opts Opts) *Tag {
	return &Tag{
		tag:   "NE",
		opts:  opts,
		short: true,
	}
}

// Like Not Equal
func Like(opts Opts) *Tag {
	return &Tag{
		tag:   "LIKE",
		opts:  opts,
		short: true,
	}
}

// NotLike Not Equal
func NotLike(opts Opts) *Tag {
	return &Tag{
		tag:   "NOTLIKE",
		opts:  opts,
		short: true,
	}
}

// In ...
func In(opts Opts) *Tag {
	return &Tag{
		tag:   "IN",
		opts:  opts,
		short: true,
	}
}

// NotIn ...
func NotIn(opts Opts) *Tag {
	return &Tag{
		tag:   "NOTIN",
		opts:  opts,
		short: true,
	}
}

// Within ...
func Within(opts Opts) *Tag {
	return &Tag{
		tag:   "WITHIN",
		opts:  opts,
		short: true,
	}
}

// Intersects ...
func Intersects(opts Opts) *Tag {
	return &Tag{
		tag:   "INTERSECTS",
		opts:  opts,
		short: true,
	}
}

// Near ...
func Near(opts Opts) *Tag {
	return &Tag{
		tag:   "NEAR",
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

var countyMap = map[int]string{
	1:  "Stockholms län",
	2:  "DEPRECATED, Användes tidigare för Stockholms län",
	3:  "Uppsala län",
	4:  "Södermanlands län",
	5:  "Östergötlands län",
	6:  "Jönköpings län",
	7:  "Kronobergs län",
	8:  "Kalmar län",
	9:  "Gotlands län",
	10: "Blekinge län",
	12: "Skåne län",
	13: "Hallands län",
	14: "Västra Götalands län",
	17: "Värmlands län",
	18: "Örebro län",
	19: "Västmanlands län",
	20: "Dalarnas län",
	21: "Gävleborgs län",
	22: "Västernorrlands län",
	23: "Jämtlands län",
	24: "Västerbottens län",
	25: "Norrbottens län",
}

// CountyNoToName returns the läns name or "Undefined" if unknown input
func CountyNoToName(n int) string {
	name, found := countyMap[n]
	if !found {
		return "Undefined"
	}
	return name
}
