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
	r := regexp.MustCompile("t\\.T\\(\"([^\"]*?)\"\\)")
	out := make(map[string]string)
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

	m, _ := json.MarshalIndent(out, "", "  ")
	log.Println(string(m))
}
