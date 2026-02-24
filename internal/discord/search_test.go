package discord

import (
	"net/url"
	"testing"
)

func TestParseSearchResponse(t *testing.T) {
	// Discord returns messages as nested arrays: [[msg1, msg2], [msg3, msg4]]
	// Each inner array is a "group" with context messages; the hit message has "hit": true
	body := []byte(`{
		"total_results": 2,
		"messages": [
			[
				{"id": "100", "channel_id": "10", "content": "before context", "timestamp": "2026-02-24T10:00:00+00:00", "hit": false, "author": {"id": "1", "username": "alice", "bot": false}},
				{"id": "101", "channel_id": "10", "content": "the headless option is great", "timestamp": "2026-02-24T10:01:00+00:00", "hit": true, "author": {"id": "2", "username": "bob", "bot": false}},
				{"id": "102", "channel_id": "10", "content": "after context", "timestamp": "2026-02-24T10:02:00+00:00", "hit": false, "author": {"id": "1", "username": "alice", "bot": false}}
			],
			[
				{"id": "200", "channel_id": "20", "content": "headless mode works well", "timestamp": "2026-02-24T11:00:00+00:00", "hit": true, "author": {"id": "3", "username": "carol", "bot": false}}
			]
		]
	}`)

	result, err := parseSearchResponse(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TotalResults != 2 {
		t.Errorf("expected TotalResults=2, got %d", result.TotalResults)
	}

	if len(result.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(result.Messages))
	}

	// First hit
	msg := result.Messages[0]
	if msg.ID != "101" {
		t.Errorf("expected first hit ID=101, got %s", msg.ID)
	}
	if msg.Content != "the headless option is great" {
		t.Errorf("unexpected content: %s", msg.Content)
	}
	if msg.Author.Username != "bob" {
		t.Errorf("expected author bob, got %s", msg.Author.Username)
	}
	if msg.ChannelID != "10" {
		t.Errorf("expected channel_id=10, got %s", msg.ChannelID)
	}

	// Second hit
	msg = result.Messages[1]
	if msg.ID != "200" {
		t.Errorf("expected second hit ID=200, got %s", msg.ID)
	}
	if msg.Author.Username != "carol" {
		t.Errorf("expected author carol, got %s", msg.Author.Username)
	}
}

func TestParseSearchResponseEmpty(t *testing.T) {
	body := []byte(`{"total_results": 0, "messages": []}`)

	result, err := parseSearchResponse(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TotalResults != 0 {
		t.Errorf("expected TotalResults=0, got %d", result.TotalResults)
	}

	if len(result.Messages) != 0 {
		t.Errorf("expected 0 messages, got %d", len(result.Messages))
	}
}

func TestParseSearchResponseMalformed(t *testing.T) {
	body := []byte(`{not valid json}`)

	_, err := parseSearchResponse(body)
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestBuildSearchParams(t *testing.T) {
	tests := []struct {
		name     string
		opts     SearchOptions
		expected map[string]string
		absent   []string
	}{
		{
			name: "content only",
			opts: SearchOptions{Content: "headless"},
			expected: map[string]string{
				"content": "headless",
			},
			absent: []string{"author_id", "channel_id", "has", "offset", "sort_by", "sort_order"},
		},
		{
			name: "all options",
			opts: SearchOptions{
				Content:   "test query",
				AuthorID:  "123",
				ChannelID: "456",
				Has:       "image",
				Offset:    25,
				SortBy:    "timestamp",
				SortOrder: "desc",
			},
			expected: map[string]string{
				"content":    "test query",
				"author_id":  "123",
				"channel_id": "456",
				"has":        "image",
				"offset":     "25",
				"sort_by":    "timestamp",
				"sort_order": "desc",
			},
		},
		{
			name: "zero offset omitted",
			opts: SearchOptions{Content: "test", Offset: 0},
			expected: map[string]string{
				"content": "test",
			},
			absent: []string{"offset"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := buildSearchParams(tt.opts)

			for key, want := range tt.expected {
				got := params.Get(key)
				if got != want {
					t.Errorf("param %q: expected %q, got %q", key, want, got)
				}
			}

			for _, key := range tt.absent {
				if params.Has(key) {
					t.Errorf("param %q should be absent, but has value %q", key, params.Get(key))
				}
			}
		})
	}
}

// buildSearchParams is a helper tested here, implemented in client.go
// parseSearchResponse is a standalone function tested here, implemented in client.go
// Both need to be package-level functions (not methods) for testability without a live Discord session.
var _ = func() url.Values { return buildSearchParams(SearchOptions{}) }
