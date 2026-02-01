package urlshort

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMapHandler(t *testing.T) {
	pathsToUrls := map[string]string{
		"/go":     "https://golang.org",
		"/github": "https://github.com",
	}

	fallbackCalled := false
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fallbackCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := MapHandler(pathsToUrls, fallback)

	t.Run("redirects when path exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/go", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusFound {
			t.Errorf("expected status %d, got %d", http.StatusFound, rr.Code)
		}

		location := rr.Header().Get("Location")
		if location != "https://golang.org" {
			t.Errorf("expected redirect to %q, got %q", "https://golang.org", location)
		}
	})

	t.Run("redirects to correct url for different path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/github", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusFound {
			t.Errorf("expected status %d, got %d", http.StatusFound, rr.Code)
		}

		location := rr.Header().Get("Location")
		if location != "https://github.com" {
			t.Errorf("expected redirect to %q, got %q", "https://github.com", location)
		}
	})

	t.Run("calls fallback when path not found", func(t *testing.T) {
		fallbackCalled = false
		req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if !fallbackCalled {
			t.Error("expected fallback to be called")
		}
	})
}

func TestYAMLHandler(t *testing.T) {
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("fallback"))
	})

	t.Run("handles valid yaml", func(t *testing.T) {
		yaml := []byte(`
- path: /go
  url: https://golang.org
- path: /github
  url: https://github.com
`)
		handler, err := YAMLHandler(yaml, fallback)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/go", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusFound {
			t.Errorf("expected status %d, got %d", http.StatusFound, rr.Code)
		}

		location := rr.Header().Get("Location")
		if location != "https://golang.org" {
			t.Errorf("expected redirect to %q, got %q", "https://golang.org", location)
		}
	})

	t.Run("returns error for invalid yaml", func(t *testing.T) {
		invalidYaml := []byte(`not: valid: yaml: [`)

		_, err := YAMLHandler(invalidYaml, fallback)
		if err == nil {
			t.Error("expected error for invalid yaml")
		}
	})

	t.Run("uses fallback for unknown path", func(t *testing.T) {
		yaml := []byte(`
- path: /go
  url: https://golang.org
`)
		handler, err := YAMLHandler(yaml, fallback)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		body := rr.Body.String()
		if body != "fallback" {
			t.Errorf("expected body %q, got %q", "fallback", body)
		}
	})
}

func TestParseYaml(t *testing.T) {
	t.Run("parses valid yaml", func(t *testing.T) {
		yaml := []byte(`
- path: /go
  url: https://golang.org
- path: /github
  url: https://github.com
`)
		result, err := parseYaml(yaml)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 entries, got %d", len(result))
		}

		if result["/go"] != "https://golang.org" {
			t.Errorf("expected /go to map to https://golang.org, got %q", result["/go"])
		}

		if result["/github"] != "https://github.com" {
			t.Errorf("expected /github to map to https://github.com, got %q", result["/github"])
		}
	})

	t.Run("returns error for invalid yaml", func(t *testing.T) {
		invalidYaml := []byte(`{invalid`)

		_, err := parseYaml(invalidYaml)
		if err == nil {
			t.Error("expected error for invalid yaml")
		}
	})

	t.Run("handles empty yaml", func(t *testing.T) {
		yaml := []byte(``)

		result, err := parseYaml(yaml)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected empty map, got %d entries", len(result))
		}
	})
}
