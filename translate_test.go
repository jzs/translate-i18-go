package translate_test

import (
	"bytes"
	"testing"

	"github.com/jzs/translate"
)

type log struct {
	t testing.TB
}

func (l log) Errorf(msg string, args ...interface{}) {
	l.t.Errorf(msg, args...)
}

var enfile = `
# English file
title:
  zero: No world
  one: One world
  many: Many worlds
  other: Other world

apple.count:
  zero: No apples
  one: 1 apple
  few: "{{.Count}} apples"
  many: "Many apples"
  other: "{{.Fart}} other apples"
`

var dkstream = bytes.NewReader([]byte(`
# Danish file
title:
  zero: Ingen verden
  one: En verden
  many: Mange verdener
  other: Andre verdener
`))

func TestTranslatePlural(t *testing.T) {
	l := log{t}

	enstream := bytes.NewReader([]byte(enfile))

	en, err := translate.LoadYaml(enstream, "en-us")
	if err != nil {
		t.Fatalf("Expected successful load of yaml file, got: %v", err)
	}
	da, err := translate.LoadYaml(dkstream, "da-dk")
	if err != nil {
		t.Fatalf("Expected successful load of yaml file, got: %v", err)
	}
	ts := translate.New(en, da)
	ts.SetLog(l)
	fun := ts.Tfunc("en-us")

	res := fun("title").String()
	if res != "One world" {
		t.Fatalf("Expected one world, got %v", res)
	}

	res = fun("apple.count").Plural(3, 10).String()
	if res != "3 apples" {
		t.Fatalf("Expected 3 apples, got %v", res)
	}

	res = fun("apple.count").With(map[string]string{"Fart": "fart"}).Other().String()
	if res != "fart other apples" {
		t.Fatalf("Expected fart other apples, got '%v'", res)
	}

	fun = ts.Tfunc("da-dk")

	res = fun("title").String()
	if res != "En verden" {
		t.Fatalf("Expected En verden, got %v", res)
	}
}

func BenchmarkTranslatePlural(b *testing.B) {
	l := log{b}
	enstream := bytes.NewReader([]byte(enfile))
	en, err := translate.LoadYaml(enstream, "en-us")
	if err != nil {
		b.Fatalf("Expected successful load of yaml file, got: %v", err)
	}

	ts := translate.New(en)
	ts.SetLog(l)
	fun := ts.Tfunc("en-us")

	for n := 0; n < b.N; n++ {
		res := fun("apple.count").With(map[string]string{"Fart": "fart"}).Other().String()
		if res != "fart other apples" {
			b.Fatalf("Expected fart other apples, got '%v'", res)
		}
	}
}

func BenchmarkTranslateOne(b *testing.B) {
	l := log{b}
	enstream := bytes.NewReader([]byte(enfile))
	en, err := translate.LoadYaml(enstream, "en-us")
	if err != nil {
		b.Fatalf("Expected successful load of yaml file, got: %v", err)
	}

	ts := translate.New(en)
	ts.SetLog(l)
	fun := ts.Tfunc("en-us")

	for n := 0; n < b.N; n++ {
		res := fun("apple.count").String()
		if res != "1 apple" {
			b.Fatalf("Expected 1 apple, got '%v'", res)
		}
	}
}

func BenchmarkTranslateOneWithTemplate(b *testing.B) {
	l := log{b}
	enstream := bytes.NewReader([]byte(enfile))
	en, err := translate.LoadYaml(enstream, "en-us")
	if err != nil {
		b.Fatalf("Expected successful load of yaml file, got: %v", err)
	}

	ts := translate.New(en)
	ts.SetLog(l)
	fun := ts.Tfunc("en-us")

	for n := 0; n < b.N; n++ {
		res := fun("apple.count").With(map[string]string{"Hello": "World"}).String()
		if res != "1 apple" {
			b.Fatalf("Expected 1 apple, got '%v'", res)
		}
	}
}
