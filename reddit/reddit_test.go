package reddit

import (
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

// testCase is an expectation for a resulting request from a single method call
// on a Bot interface.
type testCase struct {
	name    string
	err     error
	f       func(Bot) error
	correct http.Request
}

func TestAccount(t *testing.T) {
	testRequests(
		[]testCase{
			testCase{
				name: "Reply",
				f: func(b Bot) error {
					return b.Reply("name", "text")
				},
				correct: http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme:   "https",
						Host:     "reddit.com",
						Path:     "/api/comment",
						RawQuery: "text=text&thing_id=name",
					},
					Host:   "reddit.com",
					Header: formEncoding,
				},
			},
			testCase{
				name: "GetReply",
				f: func(b Bot) error {
					_, err := b.GetReply("name", "text")
					return err
				},
				correct: http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme:   "https",
						Host:     "reddit.com",
						Path:     "/api/comment",
						RawQuery: "api_type=json&text=text&thing_id=name",
					},
					Host:   "reddit.com",
					Header: formEncoding,
				},
			},
			testCase{
				name: "SendMessage",
				f: func(b Bot) error {
					return b.SendMessage("user", "subject", "text")
				},
				correct: http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme:   "https",
						Host:     "reddit.com",
						Path:     "/api/compose",
						RawQuery: "subject=subject&text=text&to=user",
					},
					Host:   "reddit.com",
					Header: formEncoding,
				},
			},
			testCase{
				name: "PostSelf",
				f: func(b Bot) error {
					return b.PostSelf("self", "title", "text")
				},
				correct: http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme:   "https",
						Host:     "reddit.com",
						Path:     "/api/submit",
						RawQuery: "kind=self&sr=self&text=text&title=title",
					},
					Host:   "reddit.com",
					Header: formEncoding,
				},
			},
			testCase{
				name: "GetPostSelf",
				f: func(b Bot) error {
					_, err := b.GetPostSelf("self", "title", "text")
					return err

				},
				correct: http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme:   "https",
						Host:     "reddit.com",
						Path:     "/api/submit",
						RawQuery: "api_type=json&kind=self&sr=self&text=text&title=title",
					},
					Host:   "reddit.com",
					Header: formEncoding,
				},
			},
			testCase{
				name: "PostLink",
				f: func(b Bot) error {
					return b.PostLink("link", "title", "url")
				},
				correct: http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme:   "https",
						Host:     "reddit.com",
						Path:     "/api/submit",
						RawQuery: "kind=link&sr=link&title=title&url=url",
					},
					Host:   "reddit.com",
					Header: formEncoding,
				},
			},
			testCase{
				name: "GetPostLink",
				f: func(b Bot) error {
					_, err := b.GetPostLink("link", "title", "url")
					return err
				},
				correct: http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme:   "https",
						Host:     "reddit.com",
						Path:     "/api/submit",
						RawQuery: "api_type=json&kind=link&sr=link&title=title&url=url",
					},
					Host:   "reddit.com",
					Header: formEncoding,
				},
			},
		}, t,
	)
}

func TestScanner(t *testing.T) {
	testRequests(
		[]testCase{
			testCase{
				name: "Listing",
				f: func(b Bot) error {
					_, err := b.Listing("/r/all", "ref")
					return err
				},
				correct: http.Request{
					Method: "GET",
					URL: &url.URL{
						Scheme:   "https",
						Host:     "reddit.com",
						Path:     "/r/all.json",
						RawQuery: "before=ref&limit=100&raw_json=1",
					},
					Host: "reddit.com",
				},
			},
		}, t,
	)
}

func TestLurker(t *testing.T) {
	testRequests(
		[]testCase{
			testCase{
				name: "Thread",
				err:  ThreadDoesNotExistErr,
				f: func(b Bot) error {
					_, err := b.Thread("/permalink")
					return err
				},
				correct: http.Request{
					Method: "GET",
					URL: &url.URL{
						Scheme:   "https",
						Host:     "reddit.com",
						Path:     "/permalink.json",
						RawQuery: "raw_json=1",
					},
					Host: "reddit.com",
				},
			},
		}, t,
	)
}

func testRequests(cases []testCase, t *testing.T) {
	c := &mockClient{}
	r := &reaperImpl{
		cli:        c,
		parser:     &mockParser{},
		hostname:   "reddit.com",
		reapSuffix: ".json",
		scheme:     "https",
		mu:         &sync.Mutex{},
	}
	b := &bot{
		Account: newAccount(r),
		Lurker:  newLurker(r),
		Scanner: newScanner(r),
		Reaper:  r,
	}
	for _, test := range cases {
		if err := test.f(b); err != test.err {
			t.Errorf("[%s] unexpected error: %v", test.name, err)
		}

		if diff := pretty.Compare(c.request, test.correct); diff != "" {
			t.Errorf(
				"[%s] request incorrect; diff: %s",
				test.name,
				diff,
			)
		}
	}
}
