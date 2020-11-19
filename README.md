t: The Trivially Translatable Thing
===================================

A very simple way to load and use embedded translations from within your Golang software.
The goal is to provide, at the application level, a way to package all translation files
within the application itself, so the application is capable of running in any locale
without requiring the user to unpack or manage any external files. Note that no
consideration has been given whatsoever to module-level translations: This package assumes
all translatable strings will exist in files underneath your current directory, included
and built into a single set of translation files. More modular design may or may not come
in the future.

## Usage:
* Include the package `github.com/justledbetter/t`
* In an initialization function, initialize the package
* Surround all translatable strings with t.T()
* Use the tscan utility to extract strings into JSON file
* Copy JSON file to other locales
* ...?
* Profit!

## Simplest Example:

main.go:
```
package main

import (
   "github.com/justledbetter/t"
   "github.com/markbates/pkger"
)

func main() {
   tr := t.Init(pkger.Dir("/i18n"))
   tr.SetGlob()

   log.Println(t.T("Some string in English"))
}
```

i18n/C.go:
```
{
"Some string in English": "Some string in English"
}
```

i18n/es.go:
```
{
"Some string in English": "Some string in Spanish :)"
}
```

```
$ go generate
$ go run
Some string in English
$ LANG=es go run
Some string in Spanish
```

## TODO
* Initialization is fairly UNIX-centric right now, need to port the locale detection function to Windows.
* Improve the tscan utility to be able to ingest existing JSON files and add newly-discovered strings (defaulting to ./lang/C.json)
