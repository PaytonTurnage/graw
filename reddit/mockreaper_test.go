package reddit

// mockReaper saves the paths it is sent and returns preconfigured results.
type mockReaper struct {
	// path is the path received by the most recent Reap or Sow call.
	path string

	h   Harvest
	s   Submission
	err error
}

func (m *mockReaper) Reap(path string, _ map[string]string) (Harvest, error) {
	m.path = path
	return m.h, m.err
}

func (m *mockReaper) Sow(path string, _ map[string]string) error {
	m.path = path
	return m.err
}

func (m *mockReaper) GetSow(path string, _ map[string]string) (Submission, error) {
	m.path = path
	return m.s, m.err
}

func (m *mockReaper) RateBlock() {}

func reaperWhich(h Harvest, err error) *mockReaper {
	return &mockReaper{
		h:   h,
		err: err,
	}
}
