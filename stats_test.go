package stats

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("bar"))
})

func TestSimple(t *testing.T) {
	s := New()

	res := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	s.Handler(testHandler).ServeHTTP(res, req)

	assert.Equal(t, s.ResponseCounts, map[string]int{"200": 1})
}
