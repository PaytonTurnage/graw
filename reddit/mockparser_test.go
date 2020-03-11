package reddit

import (
	"encoding/json"
)

type mockParser struct {
	comments   []*Comment
	posts      []*Post
	messages   []*Message
	submission Submission
}

func (m *mockParser) parse(
	blob json.RawMessage,
) ([]*Comment, []*Post, []*Message, error) {
	return m.comments, m.posts, m.messages, nil
}

func (m *mockParser) parse_submitted(
	blob json.RawMessage,
) (Submission, error) {
	return m.submission, nil
}

func parserWhich(h Harvest) parser {
	return &mockParser{
		comments: h.Comments,
		posts:    h.Posts,
		messages: h.Messages,
	}
}
