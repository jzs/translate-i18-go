package translate

import (
	"bytes"
	"html/template"
)

// Errlog is an interface that one can embed to catch strings that is not translated. It could be helpful if you want to be notified about missing translations.
type Errlog func(f string, data ...interface{})

// Language represents a language that has been loaded
type Language struct {
	ID   string
	Keys map[string]Value
}

// Value is the translated values for a key
type Value struct {
	Zero  string `json:"zero"`
	One   string `json:"one"`
	Few   string `json:"few"`
	Many  string `json:"many"`
	Other string `json:"other"`
}

type plurality uint8

const (
	pluralzero plurality = iota
	pluralone
	pluralfew
	pluralmany
	pluralother
)

// Translator struct
type Translator struct {
	langs map[string]*Language
	log   Errlog
}

// New returns a new translator with given languages
func New(langs ...*Language) *Translator {
	l := map[string]*Language{}
	for _, lang := range langs {
		l[lang.ID] = lang
	}
	return &Translator{langs: l}
}

// SetLog sets a logger on the translator for catching strings where a key or a plurality is missing
func (t *Translator) SetLog(log Errlog) {
	t.log = log
}

// T is the instance of a particular translation
type T struct {
	key    Value
	plural plurality
	count  uint64
	log    Errlog
	args   interface{}
}

// With assigns an object to T that will be merged with the translated value
func (t T) With(args interface{}) T {
	t.args = args
	return t
}

// Zero returns T with plurality set to zero
func (t T) Zero() T {
	t.plural = pluralzero
	t.count = 0
	return t
}

// Other returns T with plurality set to other
func (t T) Other() T {
	t.plural = pluralother
	return t
}

// Plural returns T with a plurality set to a meaningful value
func (t T) Plural(count, many uint64) T {
	if count == 0 {
		t.plural = pluralzero
		t.count = 0
		return t
	}
	if count == 1 {
		t.plural = pluralone
		t.count = 1
		return t
	}
	if count < many {
		t.plural = pluralfew
		t.count = count
		return t
	}

	t.plural = pluralmany
	return t
}

// renderTranslation merges a translation with given args
func renderTranslation(val string, args interface{}, log Errlog) string {
	tmpl, err := template.New("test").Parse(val)
	if err != nil && log != nil {
		log("Failed parsing translation value: %v, reason: %v", val, err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, args)
	if err != nil && log != nil {
		log("Failed executing template value: %v, reason: %v", val, err)
	}
	return buf.String()
}

// String renders the resulting translation with count and args merged
func (t T) String() string {
	var key string

	switch t.plural {
	case pluralzero:
		key = t.key.Zero
	case pluralone:
		key = t.key.One
	case pluralfew:
		key = renderTranslation(t.key.Few, map[string]uint64{"Count": t.count}, t.log)
	case pluralmany:
		key = t.key.Many
	case pluralother:
		key = t.key.Other
	default:
		key = t.key.One
	}

	if t.args != nil {
		key = renderTranslation(key, t.args, t.log)
	}

	return key
}

// Tfunc returns a translator func with languages in a given order. First language in the array has first priority, if a key is not found there, it will look in the second language etc.
func (t *Translator) Tfunc(languages ...string) func(string) T {
	return func(s string) T {
		for _, lang := range languages {
			if l, ok := t.langs[lang]; ok {
				if k, ok := l.Keys[s]; ok {
					return T{key: k, plural: pluralone, log: t.log} //initialized with default values
				}
				if t.log != nil {
					t.log("No translation match for key: %v in language %v, trying next language", s, l.ID)
					continue
				}
			} else {
				if t.log != nil {
					t.log("Language %v does not exist", l.ID)
					continue
				}
			}
		}

		if t.log != nil {
			t.log("No translation match for key: %v in any of the languages given", s)
		}
		return newT(s)
	}
}

func newT(key string) T {
	return T{
		key: Value{
			Zero:  key,
			One:   key,
			Few:   key,
			Many:  key,
			Other: key,
		},
	}
}
