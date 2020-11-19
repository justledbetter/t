t: The Trivially Translatable Thing
===================================

A very simple way to load and use embedded translations from within your Golang software.

The goal is to provide, at the application level, a way to package all translation files
within the application itself, so the application is capable of running in any locale
without requiring the user to unpack or manage any external files.

Note that no consideration has been given whatsoever to module-level translations: This
package assumes all translatable strings will exist in files underneath your current directory,
included and built into a single set of translation files. More modular design may or may
not come in the future.

## Usage:
* Include the package `github.com/justledbetter/t`
* In an initialization function, initialize the package
* Surround all translatable strings with t.T()
* Use the tscan utility to extract strings into JSON file
* Copy JSON file to other locales
* ...?
* Profit!

By default, the application initializes using the UNIX LANG environment variable (pending update).
The application can call `t.SetLocale()` to manually override this, if necessary. We (will try to)
follow common i18n fallback rules to determine the closest reasonable translation, so calling
`t.SetLocale("es_MX.UTF-8")` will search for translations in the following order:

* es_MX.UTF-8
* es_MX
* es
* C

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
   tr.SetLocale("en_UK")
   log.Println(t.T("Some string in English"))
}
```

Find strings (pending improvement):
```
$ mkdir -p ./i18n
$ go run github.com/justledbetter/t/cmd/tscan
Translations updated in ./i18n/C.json
$ cp ./i18n/C.json ./i18n/es.json
$ cp ./i18n/C.json ./i18n/en_UK.json
$ vi ./i18n/es.json ## (and customize)
$ vi ./i18n/en_UK.json ## (and customize)
```

i18n/C.json:
```
{
"Some string in English": "Some string in English"
}
```

i18n/es.json:
```
{
"Some string in English": "Some string in Spanish :)"
}
```

i18n/en_UK.json:
```
{
"Some string in English": "Some string in colourful English"
}
```

```
$ go generate
$ go run
Some string in English
Some string in colourful English
$ LANG=es go run
Some string in Spanish :)
Some string in colourful English
```

## TODO
* Initialization is fairly UNIX-centric right now, need to port the locale detection function to Windows.
* Locale names need to be harmonized with i18n standards, it's a little bit fast and loose right now.
* Add the ability to have static translations loaded from in-application structs, if necessary.
* Add the ability to import the package with a name other than "t" (and still pick up strings with tscan).
* Allow the user to customize the behavior of the tscan utility.
