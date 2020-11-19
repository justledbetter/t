/**
 * t: The Trivially Translatable Thing
 */
package t

import (
	"encoding/json"
	"github.com/markbates/pkger"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type trans map[string]string
type alltrans map[string]trans

type t interface {
	T(string) string
}
type tdata struct {
	t
	tbl        alltrans
	lastLocale string
}

var glob *tdata

func Init(d pkger.Dir) tdata {
	ret := tdata{
		lastLocale: "C",
	}
	ret.tbl = make(map[string]trans)
	err := pkger.Walk("/i18n", func(path string, info os.FileInfo, _ error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		f, _ := pkger.Open(path)
		defer f.Close()

		sl, _ := ioutil.ReadAll(f)
		var tr trans
		e := json.Unmarshal(sl, &tr)
		if e != nil {
			return e
		}

		path = string(regexp.MustCompile("^.*:").ReplaceAll([]byte(path), []byte("")))
		path = strings.ReplaceAll(path, "/i18n/", "")
		path = strings.ReplaceAll(path, ".json", "")

		ret.tbl[path] = tr

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return ret
}

func (t *tdata) SetGlobal() {
	glob = t
}

func (t *tdata) fallback(in string) string {
	if t.lastLocale == in {
		return in
	}

	if _, ok := t.tbl[in]; !ok {
		if strings.Contains(in, ".") {
			return t.fallback(strings.Split(in, ".")[0])
		}
		if strings.Contains(in, "_") {
			return t.fallback(strings.Split(in, "_")[0])
		}

		return "C"
	}

	return in
}

func (t *tdata) locale() string {

	// TODO: Is this necessary?  How to do this in Windows?
	if len(t.lastLocale) < 1 {
		return t.fallback(os.Getenv("LANG"))
	}

	return t.lastLocale
}

func SetLocale(in string) string {
	if glob == nil {
		return in
	}

	return glob.SetLocale(in)
}

func (t *tdata) SetLocale(in string) string {

	// Cache this value
	t.lastLocale = t.fallback(in)

	return t.lastLocale
}

func (t *tdata) T(in string) string {
	l := t.locale()
	if _, ok := t.tbl[l][in]; !ok {
		return in
	}
	return t.tbl[l][in]
}

func T(in string) string {
	if glob != nil {
		return glob.T(in)
	}

	return in
}
