package note

import (
	"time"
)

type Note struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Links     []string  `json:"links"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	StoreLocation string
}

type Index struct {
	TagIndex     map[string][]string `json:"tags"`
	KeywordIndex map[string][]string `json:"keywords"`
}

type NoteSummary struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}
