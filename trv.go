package trv

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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
	OptName              = "name"
	OptObjtype           = "objecttype"
	OptOrderBy           = "orderby"
	OptRadius            = "radius"
	OptSchemaversion     = "schemaversion"
	OptShape             = "shape"
	OptLimit             = "limit"
	OptValue             = "value"
	OptAuthenticationKey = "authenticationkey"
)

// NewRequest createsa new TRV request
func NewRequest(apiKey string, query *Tag) *Tag {
	return Request().Add(
		Login().Opts(Opts{
			OptAuthenticationKey: apiKey,
		}),
		query,
	)
}

/* // TagInterface ...
type TagInterface interface {
	Add(...*Tag) *Tag
	Opts(Opts) *Tag
	Value(string) *Tag
} */

// Opts holds our tag options <TAG key="value">
type Opts map[string]string

// Set a opt value
func (o Opts) Set(key, value string) {
	o[key] = value
}

// Tag ...
type Tag struct {
	tag        string // The tag
	level      int    // How deep is the tag, for pretty printing
	nochildren bool
	children   []*Tag // Children tags
	opts       Opts   // Tag options
	value      string // The Value inside the <tag></tag>
	short      bool   // true closes the tag with <TAG/> false closes it with a <TAG></TAG>
	inline     bool   // Should the opening & closing tags be inline ( on the same line )
}

func (t *Tag) String() string {
	return fmt.Sprintf("<%s> <opts: %q> <value: %q> <short: %t> <inline: %t>", t.tag, t.opts, t.value, t.short, t.inline)
}

// Add new child(ren) to the tag
func (t *Tag) Add(tags ...*Tag) *Tag {
	t.children = append(t.children, tags...)
	return t
}

// Opts sets options on the tag
func (t *Tag) Opts(opts Opts) *Tag {
	t.opts = opts
	return t
}

// Value sets options on the tag
func (t *Tag) Value(value string) *Tag {
	t.value = value
	return t
}

// Do sends the request
func (t *Tag) Do() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := t.Build(buf); err != nil {
		return nil, err
	}
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return body, fmt.Errorf("return code %d", resp.StatusCode)
	}

	return body, nil
}

// Request tag
func Request() *Tag {
	return &Tag{
		tag: "REQUEST",
	}
}

// Login tag
func Login() *Tag {
	return &Tag{
		tag:        "LOGIN",
		short:      true,
		nochildren: true,
	}
}

// Query tag
func Query() *Tag {
	return &Tag{
		tag: "QUERY",
	}
}

// Filter tag
func Filter() *Tag {
	return &Tag{
		tag: "FILTER",
	}
}

// And tag
func And() *Tag {
	return &Tag{
		tag: "AND",
	}
}

// Or tag
func Or() *Tag {
	return &Tag{
		tag: "OR",
	}
}

// Exists Exists
func Exists() *Tag {
	return &Tag{
		tag:        "EXISTS",
		short:      true,
		nochildren: true,
	}
}

// Eq equal tag
func Eq() *Tag {
	return &Tag{
		tag:        "EQ",
		short:      true,
		nochildren: true,
	}
}

// Gt Greater Than
func Gt() *Tag {
	return &Tag{
		tag:        "GT",
		short:      true,
		nochildren: true,
	}
}

// Gte Greater Than or Equal
func Gte() *Tag {
	return &Tag{
		tag:        "GTE",
		short:      true,
		nochildren: true,
	}
}

// Lt Less Than
func Lt() *Tag {
	return &Tag{
		tag:        "LT",
		short:      true,
		nochildren: true,
	}
}

// Lte Less Than or Equals
func Lte() *Tag {
	return &Tag{
		tag:        "LTE",
		short:      true,
		nochildren: true,
	}
}

// Ne Not Equal
func Ne() *Tag {
	return &Tag{
		tag:        "NE",
		short:      true,
		nochildren: true,
	}
}

// Like Not Equal
func Like() *Tag {
	return &Tag{
		tag:        "LIKE",
		short:      true,
		nochildren: true,
	}
}

// NotLike Not Equal
func NotLike() *Tag {
	return &Tag{
		tag:        "NOTLIKE",
		short:      true,
		nochildren: true,
	}
}

// In ...
func In() *Tag {
	return &Tag{
		tag:        "IN",
		short:      true,
		nochildren: true,
	}
}

// NotIn ...
func NotIn() *Tag {
	return &Tag{
		tag:        "NOTIN",
		short:      true,
		nochildren: true,
	}
}

// Within ...
func Within() *Tag {
	return &Tag{
		tag:        "WITHIN",
		short:      true,
		nochildren: true,
	}
}

// Intersects ...
func Intersects() *Tag {
	return &Tag{
		tag:        "INTERSECTS",
		short:      true,
		nochildren: true,
	}
}

// Near ...
func Near() *Tag {
	return &Tag{
		tag:        "NEAR",
		short:      true,
		nochildren: true,
	}
}

// Include .tag
func Include() *Tag {
	return &Tag{
		tag:        "INCLUDE",
		inline:     true,
		nochildren: true,
	}
}

// Build a TRV query
func (t *Tag) Build(w io.Writer) error {
	if err := t.start(w); err != nil {
		return err
	}
	if !t.nochildren {
		for _, c := range t.children {
			c.level = t.level + 1
			if err := c.Build(w); err != nil {
				return err
			}
		}
	}
	if err := t.end(w); err != nil {
		return err
	}
	return nil
}

// opens the tag
func (t *Tag) start(w io.Writer) error {
	if t.nochildren && len(t.children) > 0 {
		return fmt.Errorf("tag %s should not have any children", t.tag)
	}

	if PrettyPrint {
		w.Write([]byte(strings.Repeat(IndentChar, t.level)))
	}

	// Write opening for tag
	w.Write([]byte("<" + t.tag))

	// If there are any options, write them
	if t.opts != nil {
		w.Write([]byte(t.optsString()))
	}

	// If it's a short tag, close it now
	if t.short {
		w.Write([]byte("/>\n"))
		return nil
	}

	// Write end of tag opener
	w.Write([]byte(">"))

	// If the tag has a value write it now
	if t.value != "" {
		w.Write([]byte(t.value))
	}

	if !t.inline {
		w.Write([]byte("\n"))
	}

	return nil
}

// closes the tag
func (t *Tag) end(w io.Writer) error {
	// If it's a short tag the start() function renders it in whole
	if t.short {
		return nil
	}
	if PrettyPrint && !t.inline {
		w.Write([]byte(strings.Repeat(IndentChar, t.level)))
	}
	w.Write([]byte("</" + t.tag + ">\n"))
	return nil
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

// TagFunc ...
type TagFunc func() *Tag

var verbMap = map[string]TagFunc{
	"QUERY":      Query,
	"INCLUDE":    Include,
	"FILTER":     Filter,
	"AND":        And,
	"OR":         Or,
	"EXISTS":     Exists,
	"EQ":         Eq,
	"GT":         Gt,
	"GTE":        Gte,
	"LT":         Lt,
	"LTE":        Lte,
	"NE":         Ne,
	"LIKE":       Like,
	"NOTLIKE":    NotLike,
	"IN":         In,
	"NOTIN":      NotIn,
	"WITHIN":     Within,
	"INTERSECTS": Intersects,
	"NEAR":       Near,
}

// VerbToFunc ...
func VerbToFunc(v string) (TagFunc, bool) {
	f, found := verbMap[strings.ToUpper(v)]
	return f, found
}
