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
	if !slices.Contains(n.Links, targetID) {
		return errors.New("link already present")
	}
	n.Links = append(n.Links, targetID)
	return nil
}

func (n *Note) RemoveLink(targetID string) error {
	// If note id not found in n.Links list return error else remove
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
	// If note tag not found in n.Tags list return error else remove
	tags := strings.Split(strings.ToLower(tagList), ",")
	indexes := []int{}
	for _, tag := range tags {
		indexes = append(indexes, slices.Index(n.Tags, strings.ToLower(tag)))
	}

	if slices.Contains(indexes, -1) {
		return errors.New(tags[slices.Index(indexes, -1)] + "tag not found")
	}
	for _, index := range indexes {
		n.Tags = slices.Delete(n.Tags, index, index)
	}
	return nil
}
