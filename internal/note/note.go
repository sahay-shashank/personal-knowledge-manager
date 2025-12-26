package note

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

func NewNote(title string, content string) *Note {
	return &Note{
		Id:        uuid.New().String(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
}

func (n *Note) AddLink(targetID string) error {
	// If note id not found return error
	if slices.Contains(n.Links, targetID) {
		return errors.New("link already present")
	}
	n.Links = append(n.Links, targetID)
	return nil
}

func (n *Note) RemoveLink(targetID string) error {
	// If note id not found in n.Links list return error else remove
	index := slices.Index(n.Links, targetID)
	if index == -1 {
		return errors.New(targetID + "link not found")
	}
	n.Links = slices.Delete(n.Links, index, index+1)
	return nil
}

func (n *Note) AddTag(tagList string) error {
	// if tags has some issues return error
	tags := strings.Split(strings.ToLower(tagList), ",")
	for _, tag := range tags {
		if slices.Contains(n.Tags, strings.ToLower(tag)) {
			return errors.New("tag already present")
		}
	}

	n.Tags = append(n.Tags, tags...)
	return nil
}

func (n *Note) RemoveTag(tagList string) error {
	tags := strings.Split(strings.ToLower(tagList), ",")
	var indexes []int
	for _, tag := range tags {
		tag = strings.TrimSpace(tag) // trim whitespace
		index := slices.Index(n.Tags, tag)
		if index == -1 {
			return errors.New("tag not found")
		}
		indexes = append(indexes, index)
	}

	// Delete in reverse order to avoid index shifting
	for i := len(indexes) - 1; i >= 0; i-- {
		n.Tags = slices.Delete(n.Tags, indexes[i], indexes[i]+1)
	}
	return nil
}
