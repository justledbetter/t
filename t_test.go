package t

import (
	"testing"
)

func TestT(test *testing.T) {

	// Confirm the t.T function acts as a no-op if the translation system is not yet loaded.
	if T("not loaded yet") != "not loaded yet" {
		test.Fatalf("String mismatch: not loaded yet")
	}

	// Initialize the translation system into a local object named "t"
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

	s = t.SetLocale("en_UK")
	if s != "en_UK" {
		test.Fatalf("Failed loading translation: en_UK got %v", s)
	}

	if t.T("com.testing.error35") != "Something improper has occurred" {
		test.Fatalf("String mismatch: error35 es")
	}

	s = t.SetLocale("en_AU")
	if s != "en" {
		test.Fatalf("Failed loading translation: en got %v", s)
	}

	if t.T("com.testing.error35") != "An error has occurred" {
		test.Fatalf("String mismatch: error35 es")
	}
}
