package reddit

// deletedAuthor is the author field of deleted posts on Reddit.
const deletedAuthor = "[deleted]"

// Scanner defines a low level interface for fetching reading Reddit listings.
type Scanner interface {
	// Listing returns a harvest from a listing endpoint at Reddit.
	//
	// There are many things to consider when using this. A listing on
	// Reddit is like an infinite, unstable list. It stretches functionally
	// infinitely forward and backward in time, and elements from it vanish
	// over time for many reasons (users delete their posts, posts get
	// caught in a spam filter, mods remove them, etc).
	//
	// The "after" parameter is the name (in the form tx_xxxxxx, found by
	// the Name field of any Reddit struct defined in this package) of an
	// element known to be in the list. If an empty string the latest 100
	// elements are returned.
	//
	// The way to subscribe to a listing is continually poll this, and keep
	// track of your reference point, and replace it if it gets deleted or
	// dropped from the listing for any reason, which is nontrivial. If your
	// reference point becomes invalid, you will get no elements in the
	// harvest and have no way to find your place unless you planned ahead.
	//
	// If you want a stream where all of this is handled for you, see graw
	// or graw/streams.
	Listing(path, after string) (Harvest, error)
	ListingWithParams(path string, params map[string]string) (Harvest, error)
}

type scanner struct {
	r Reaper
}

func newScanner(r Reaper) Scanner {
	return &scanner{r: r}
}

func (s *scanner) Listing(path, after string) (Harvest, error) {
	return s.r.Reap(
		path, map[string]string{
			"raw_json": "1",
			"limit":    "100",
			"before":   after,
		},
	)
}

func (s *scanner) ListingWithParams(path string, params map[string]string) (
	Harvest,
	error,
) {
	reaperParams := map[string]string{
		"raw_json": "1",
		"limit":    "100",
	}
	for key, value := range params {
		reaperParams[key] = value
	}
	return s.r.Reap(path, reaperParams)
}
