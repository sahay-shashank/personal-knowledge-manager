package note

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"slices"
)

func InitStore(storeDirectory string) *Store {
	return &Store{
		StoreLocation: storeDirectory,
	}
}

func (fileStore *Store) Save(note *Note) error {
	indexFileFS, err := os.OpenFile(fileStore.StoreLocation+"/"+".index.pkm", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer indexFileFS.Close()
	indexFile, err := io.ReadAll(indexFileFS)
	if err != nil {
		return err
	}
	var index Index
	index.KeywordIndex = make(map[string][]string)
	index.TagIndex = make(map[string][]string)
	if len(indexFile) > 0 {
		if err := json.Unmarshal(indexFile, &index); err != nil {
			return err
		}
	}

	allString := normalize(note.Title + " " + note.Content)
	words := filterStopwords(allString)

	for _, word := range words {
		if !slices.Contains(index.KeywordIndex[word], note.Id) {
			index.KeywordIndex[word] = append(index.KeywordIndex[word], note.Id)
		}
	}
	for _, tag := range note.Tags {
		if !slices.Contains(index.TagIndex[tag], note.Id) {
			index.TagIndex[tag] = append(index.TagIndex[tag], note.Id)
		}
	}
	indexPayload, err := json.Marshal(index)
	if err != nil {
		return err
	}

	if err := indexFileFS.Truncate(0); err != nil {
		return err
	}
	if _, err := indexFileFS.Seek(0, 0); err != nil {
		return err
	}
	if _, err := indexFileFS.Write(indexPayload); err != nil {
		return err
	}

	noteFileFS, err := os.OpenFile(fileStore.StoreLocation+"/"+note.Id+".pkm", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer noteFileFS.Close()
	jsonBody, err := json.Marshal(note)
	if err != nil {
		return err
	}

	payload := []byte("PKM\n")
	payload = append(payload, jsonBody...)

	if err := noteFileFS.Truncate(0); err != nil {
		return err
	}
	if _, err := noteFileFS.Seek(0, 0); err != nil {
		return err
	}
	if _, err := noteFileFS.Write(payload); err != nil {
		return err
	}

	return nil
}

func (fileStore *Store) Load(noteLocation string) (*Note, error) {
	fileDataFS, err := os.OpenFile(fileStore.StoreLocation+"/"+noteLocation+".pkm", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fileDataFS.Close()
	fileData, err := io.ReadAll(fileDataFS)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(fileData[:4], []byte("PKM\n")) {
		return nil, errors.New("note corrupted")
	}

	note := Note{}
	if err := json.Unmarshal(fileData[4:], &note); err != nil {
		return nil, err
	}
	return &note, nil
}

func (fileStore *Store) Delete(noteLocation string) error {
	return os.Remove(fileStore.StoreLocation + "/" + noteLocation + ".pkm")
}

func (fileStore *Store) Search(searchType string, terms []string) ([]string, error) {
	fileDataFS, err := os.OpenFile(fileStore.StoreLocation+"/"+".index.pkm", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fileDataFS.Close()
	fileData, err := io.ReadAll(fileDataFS)
	if err != nil {
		return nil, err
	}
	var index Index
	if err := json.Unmarshal(fileData, &index); err != nil {
		return nil, err
	}

	var candidates []string
	switch searchType {
	case "tag":
		candidates = index.TagIndex[terms[0]]
		for _, term := range terms {
			candidates = intersect(candidates, index.TagIndex[term])
		}
	case "keyword":
		candidates = index.KeywordIndex[terms[0]]
		for _, term := range terms {
			candidates = intersect(candidates, index.KeywordIndex[term])
		}
	}

	return candidates, nil
}

func intersect[T comparable](a []T, b []T) []T {
	result := make([]T, 0)
	hash := make(map[T]struct{})
	for _, v := range a {
		hash[v] = struct{}{}
	}
	seen := make(map[T]struct{})
	for _, v := range b {
		if _, ok := hash[v]; ok {
			if _, added := seen[v]; !added {
				result = append(result, v)
				seen[v] = struct{}{}
			}
		}
	}
	return result
}
