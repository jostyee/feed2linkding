package main

type LinkdingURL struct {
	Link        string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Notes       string   `json:"notes"`
	IsArchived  bool     `json:"is_archived"`
	IsUnread    bool     `json:"unread"`
	IsShared    bool     `json:"shared"`
	TagNames    []string `json:"tag_names,omitempty"`
}
