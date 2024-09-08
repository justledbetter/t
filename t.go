/**
 * t: The Trivially Translatable Thing
 */
package t

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
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
	ctx        context.Context
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

func (t *tdata) localeFromRequest(r *http.Request) string {
	if len(r.Header) < 1 {
		return "C"
	}

	if len(r.Header.Get("Accept-Language")) < 1 {
		return "C"
	}

	for _, s := range strings.Split(r.Header.Get("Accept-Language"), ";") {
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

func Locale() string {
	if glob != nil {
		return glob.Locale()
	}

	return "C" // TODO this is not correct
}

func (t tdata) Locale() string {
	if len(t.lastLocale) > 0 {
		return t.lastLocale
	} else {
		return "C" // TODO this is not correct
	}
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

// R refers to tdata.R if it can.
func R(r *http.Request, in string) string {
	if glob != nil {
		return glob.R(r, in)
	}

	return in
}

// R infers the user's locale from the HTTP request using the tdata.localeFromRequest helper func
func (t *tdata) R(r *http.Request, in string) string {
	l := t.localeFromRequest(r)
	if _, ok := t.tbl[l][in]; !ok {
		return in
	}
	return t.tbl[l][in]
}

// ContextForLocale creates a new per-user context that contains a tdata struct _just for their requested locale_.
func ContextForLocale(ctx context.Context, in string) context.Context {
	if glob == nil {
		return ctx
	}

	t := *glob
	t.lastLocale = t.fallback(in)

	return context.WithValue(ctx, "t.locale", &t)
}

// C calls tdata.C
func C(ctx context.Context, in string) string {
	td := &tdata{} // WE DO NOT NEED GLOB FOR THIS (it'll be pulled up later if needed)
	return td.C(ctx, in)
}

// C is responsible for translating based on the locale selected by the _user_.  If none is found in the context.Value("t.locale"), it reverts to glob.
func (t *tdata) C(ctx context.Context, in string) string {
	tt := t.tdataFromContext(ctx)
	if tt == nil {
		return in
	}

	return tt.T(in)
}

func (t tdata) tdataFromContext(ctx context.Context) *tdata {
	v := ctx.Value("t.locale")
	if locale, ok := v.(*tdata); ok {
		return locale
	}

	return glob
}
