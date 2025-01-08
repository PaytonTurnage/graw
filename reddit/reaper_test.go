package reddit

import (
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
)

func TestNew(t *testing.T) {
	cli := &mockClient{}
	par := &mockParser{}
	cfg := reaperConfig{
		client:   cli,
		parser:   par,
		hostname: "com",
		tls:      true,
	}
	expected := &reaperImpl{
		cli:      cli,
		parser:   par,
		hostname: "com",
		scheme:   "https",
		mu:       &sync.Mutex{},
	}

	if diff := pretty.Compare(newReaper(cfg), expected); diff != "" {
		t.Errorf("reaper construction incorrect; diff: %s", diff)
	}
}

func TestReap(t *testing.T) {
	for i, test := range []struct {
		path    string
		values  map[string]string
		correct http.Request
	}{
		{"", nil, http.Request{
			Method: "GET",
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "",
				RawQuery: "",
			},
		}},
		{"", map[string]string{"key": "value"}, http.Request{
			Method: "GET",
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "",
				RawQuery: "key=value",
			},
		}},
		{"path", nil, http.Request{
			Method: "GET",
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "path",
				RawQuery: "",
			},
		}},
	} {
		expected := Harvest{
			Comments: []*Comment{
				&Comment{
					Body: "comment",
				},
			},
			Posts: []*Post{
				&Post{
					SelfText: "post",
				},
			},
			Messages: []*Message{
				&Message{
					Body: "message",
				},
			},
		}
		c := &mockClient{}
		r := &reaperImpl{
			cli:      c,
			parser:   parserWhich(expected),
			hostname: "com",
			scheme:   "http",
			mu:       &sync.Mutex{},
		}

		Harvest, err := r.Reap(test.path, test.values)
		if err != nil {
			t.Errorf("Error reaping input %d: %v", i, err)
		}

		if diff := pretty.Compare(Harvest, expected); diff != "" {
			t.Errorf("Harvest incorrect; diff: %s", diff)
		}

		if diff := pretty.Compare(c.request, test.correct); diff != "" {
			t.Errorf("request incorrect; diff: %s", diff)
		}
	}
}

func TestSow(t *testing.T) {
	for i, test := range []struct {
		path    string
		values  map[string]string
		correct http.Request
	}{
		{"", nil, http.Request{
			Method: "POST",
			Header: formEncoding,
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "",
				RawQuery: "",
			},
		}},
		{"", map[string]string{"key": "value"}, http.Request{
			Method: "POST",
			Header: formEncoding,
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "",
				RawQuery: "key=value",
			},
		}},
		{"path", nil, http.Request{
			Method: "POST",
			Header: formEncoding,
			Host:   "com",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "com",
				Path:     "path",
				RawQuery: "",
			},
		}},
	} {
		c := &mockClient{}
		r := &reaperImpl{
			cli:      c,
			parser:   &mockParser{},
			hostname: "com",
			scheme:   "http",
			mu:       &sync.Mutex{},
		}

		if err := r.Sow(test.path, test.values); err != nil {
			t.Errorf("Error reaping input %d: %v", i, err)
		}

		if diff := pretty.Compare(c.request, test.correct); diff != "" {
			t.Errorf("request incorrect; diff: %s", diff)
		}
	}
}

func TestRateBlockReap(t *testing.T) {
	testRateBlock(func(r Reaper) { r.Reap("", nil) }, t)
}

func TestRateBlockSow(t *testing.T) {
	testRateBlock(func(r Reaper) { r.Sow("", nil) }, t)
}

func testRateBlock(f func(Reaper), t *testing.T) {
	start := time.Now()
	r := &reaperImpl{
		cli:    &mockClient{},
		parser: &mockParser{},
		rate:   10 * time.Millisecond,
		last:   start,
		mu:     &sync.Mutex{},
	}

	f(r)
	end := time.Now()

	if block := end.Sub(start); block < r.rate {
		t.Errorf("wanted block for %v; blocked for %v", r.rate, block)
	} else if r.last == start {
		t.Errorf("wanted updated timestamp; found same timestamp")
	}
}
