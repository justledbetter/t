package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {

	log.SetFlags(0)

	// Default package name is t.T.
	r := regexp.MustCompile("t\\.T\\(\"([^\"]*?)\"\\)")

	out := make(map[string]string)

	// TODO: allow the user to override package name

	// Load existing translations.
	if f, e := os.Open("./i18n/C.json"); e == nil {
		defer f.Close()
		sl, _ := ioutil.ReadAll(f)
		json.Unmarshal(sl, &out)
	} else {
		log.Fatal("Unable to load existing translations from ./i18n/C.json.  Please create this file before continuing.")
	}

	// Look for strings in all .go files under the current directory.
	e := filepath.Walk("./", func(file string, info os.FileInfo, e error) error {
		if !strings.HasSuffix(file, ".go") {
			return nil
		}
		o, e := os.Open(file)
		defer o.Close()
		if e != nil {
			return e
		}
		sl, e := ioutil.ReadAll(o)
		if e != nil {
			return e
		}
		// Change to reader
		bi := r.FindAllStringSubmatch(string(sl), -1)
		for _, i := range bi {
			if _, ok := out[i[1]]; !ok {
				out[i[1]] = i[1]
			}
		}
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}

	if f, e := os.OpenFile("./i18n/C.json",os.O_WRONLY|os.O_TRUNC|os.O_CREATE,0644); e == nil {
		defer f.Close()
		m, _ := json.MarshalIndent(out, "", "\t")
		f.Write(m)
		f.Write([]byte("\n"))
		log.Println("Translations updated in ./i18n/C.json")
	} else {
		m, _ := json.MarshalIndent(out, "", "\t")
		log.Println(string(m))
	}
}
