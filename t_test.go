package t

import (
	"testing"
	"net/http"
	"embed"
)

//go:embed i18n
var i18n embed.FS

func TestT(test *testing.T) {

	// Confirm the t.T function acts as a no-op if the translation system is not yet loaded.
	if T("not loaded yet") != "not loaded yet" {
		test.Fatalf("String mismatch: not loaded yet")
	}

	// Initialize the translation system into a local object named "t" (largely for compatibility purpose)
	t := Init(i18n)

	// Set this object as the global translator.
	t.SetGlobal()

	// Confirm the basic C translation is present
	if t.T("hello") != "hello" {
		test.Fatalf("String mismatch: hello")
	}

	if t.T("com.testing.error35") != "Error 35" {
		test.Fatalf("String mismatch: error35 C")
	}

	// Confirm we can switch to a Spanish translation
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

	// Confirm we can switch to an English translation
	s = t.SetLocale("en_UK")
	if s != "en_UK" {
		test.Fatalf("Failed loading translation: en_UK got %v", s)
	}

	if t.T("com.testing.error35") != "Something improper has occurred" {
		test.Fatalf("String mismatch: error35 en_UK")
	}

	// Confirm we can fallback to "standard" English after failing to load Australian English
	s = t.SetLocale("en_AU")
	if s != "en" {
		test.Fatalf("Failed loading translation: en got %v", s)
	}

	if t.T("com.testing.error35") != "An error has occurred" {
		test.Fatalf("String mismatch: error35 en")
	}

	r, _ := http.NewRequest("GET","http://localhost",nil)
	r.Header.Add("Accept-Language", "en-us; en")
	if t.R(r, "com.testing.error35") != "An error has occurred" {
		test.Fatalf("String mismatch: error35 http en %v", t.R(r, "com.testing.error35"))
	}
}
