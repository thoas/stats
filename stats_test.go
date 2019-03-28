package stats

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/websocket"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, map[string]int{"200": 1}, s.ResponseCounts)
}

func TestGetStats(t *testing.T) {
	s := New()

	var stats = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		stats := s.Data()

		b, _ := json.Marshal(stats)

		w.Write(b)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
	})

	res := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	s.Handler(testHandler).ServeHTTP(res, req)

	res = httptest.NewRecorder()

	s.Handler(stats).ServeHTTP(res, req)

	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var data map[string]interface{}

	err := json.Unmarshal(res.Body.Bytes(), &data)

	assert.Nil(t, err)

	assert.Equal(t, float64(1), data["total_count"].(float64))
}

func TestRace(t *testing.T) {
	s := New()

	ch1 := make(chan bool)
	ch2 := make(chan bool)

	go func() {
		now := time.Now()
		for true {
			select {
			case _ = <-ch1:
				return
			default:
				s.End(now, WithStatusCode(200))

			}
		}

	}()

	go func() {
		dt := s.Data()
		for true {
			select {
			case _ = <-ch2:
				return
			default:
				_ = dt.TotalStatusCodeCount["200"]
			}
		}
	}()

	time.Sleep(time.Second)

	ch1 <- true
	ch2 <- true
}

func TestWebsocketIgnore(t *testing.T) {
	s := New()
	handler := websocket.Handler(func(conn *websocket.Conn) {
		conn.Write([]byte("TEST"))
	})

	srv := httptest.NewServer(s.Handler(handler))

	url := strings.Replace(srv.URL, "http", "ws", 1)
	ws, err := websocket.Dial(url, "", srv.URL)
	require.NoError(t, err)

	bytes := make([]byte, 255)
	n, err := ws.Read(bytes)
	require.NoError(t, err)

	assert.Equal(t, "TEST", string(bytes[:n]))
	assert.Equal(t, map[string]int{}, s.ResponseCounts)

}

func TestWebsocketErrorNotIgnore(t *testing.T) {
	s := New()

	srv := httptest.NewServer(s.Handler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	})))

	url := strings.Replace(srv.URL, "http", "ws", 1)
	_, err := websocket.Dial(url, "", srv.URL)
	require.Error(t, err)

	assert.Equal(t, map[string]int{"500": 1}, s.ResponseCounts)

}

func TestIgnoreHijackedConnection(t *testing.T) {
	s := New()

	res := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	s.Handler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.(http.Hijacker).Hijack()
	})).ServeHTTP(res, req)

	assert.Equal(t, map[string]int{}, s.ResponseCounts)
}
