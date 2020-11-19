t: The Trivially Translatable Thing
===================================

This package provides a very simple way to load and use embedded translations from within your
Golang software.

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

/i18n/C.go:
```
{
"Some string in English": "Some string in English"
}
```

/i18n/es.go:
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
