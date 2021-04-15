/**
 * t: The Trivially Translatable Thing
 */
package t

import (
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"
	"net/http"
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

func Init(d fs.FS) tdata {
	ret := tdata{
		lastLocale: "C",
	}
	ret.tbl = make(map[string]trans)
	err := fs.WalkDir(d, ".", func(path string, info fs.DirEntry, _ error) error {
		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		f, e := d.Open(path)
		if e != nil {
			return e
		}
		defer f.Close()

		sl, _ := io.ReadAll(f)
		var tr trans
		e = json.Unmarshal(sl, &tr)
		if e != nil {
			return e
		}

		path = string(regexp.MustCompile("^.*:").ReplaceAll([]byte(path), []byte("")))
		path = strings.ReplaceAll(path, "i18n/", "")
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

func (t *tdata) localeFromRequest( r *http.Request ) string {
	if len(r.Header) < 1 {
		return "C"
	}

	if len(r.Header.Get("Accept-Language")) < 1 {
		return "C"
	}

	for _, s := range strings.Split(r.Header.Get("Accept-Language"),";") {
		// TODO: Grossly inefficient and potentially wasteful and not RFC 7231 compliant
		if t.fallback(strings.TrimSpace(s)) != "C" {
			return t.fallback(strings.TrimSpace(s))
		}
	}

	return "C"
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

func R(r *http.Request, in string) string {
	if glob != nil {
		return glob.R(r, in)
	}

	return in
}

func (t *tdata) R(r *http.Request, in string) string {
	l := t.localeFromRequest(r)
	if _, ok := t.tbl[l][in]; !ok {
		return in
	}
	return t.tbl[l][in]
}
