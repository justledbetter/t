package t

import (
	"testing"
)

func TestT(test *testing.T) {
	t := Init("/i18n")

	if t.T("hello") != "hello" {
		test.Fatalf("String mismatch: hello")
	}

	if t.T("com.testing.error35") != "Error 35" {
		test.Fatalf("String mismatch: error35 C")
	}

	s := t.SetLocale("es")
	if s != "es" {
		test.Fatalf("Failed loading translation: es got %v", s)
	}

	if t.T("com.testing.error35") != "Ha ocurrido un error del sistema" {
		test.Fatalf("String mismatch: error35 es")
	}

	if t.T("hello") != "hola" {
		test.Fatalf("String mismatch: hola got %v", t.T("hello"))
	}

	s = t.SetLocale("en-UK")
	if s != "en-UK" {
		test.Fatalf("Failed loading translation: en-UK got %v", s)
	}

	if t.T("com.testing.error35") != "Something improper has occurred" {
		test.Fatalf("String mismatch: error35 es")
	}

	s = t.SetLocale("en-AU")
	if s != "en" {
		test.Fatalf("Failed loading translation: en got %v", s)
	}

	if t.T("com.testing.error35") != "An error has occurred" {
		test.Fatalf("String mismatch: error35 es")
	}
}
