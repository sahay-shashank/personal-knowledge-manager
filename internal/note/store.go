package note

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/sahay-shashank/personal-knowledge-manager/internal/crypt"
)

func InitStore(storeDirectory string) *Store {
	return &Store{
		StoreLocation: storeDirectory,
	}
}

func (fileStore *Store) Save(note *Note, username string, kp *crypt.KeyProvider) error {
	userdir := filepath.Join(fileStore.StoreLocation, username)
	if err := os.MkdirAll(userdir, 0755); err != nil {
		return err
	}

	jsonBody, err := json.Marshal(note)
	if err != nil {
		return err
	}

	encryptedBody, err := kp.Encrypt(jsonBody)
	if err != nil {
		return err
	}
	payload := []byte("PKM\n")
	payload = append(payload, encryptedBody...)

	noteFilePath := filepath.Join(userdir, note.Id+".pkm")
	noteFileFS, err := os.OpenFile(noteFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer noteFileFS.Close()

	if err := noteFileFS.Truncate(0); err != nil {
		return err
	}
	if _, err := noteFileFS.Seek(0, 0); err != nil {
		return err
	}
	if _, err := noteFileFS.Write(payload); err != nil {
		return err
	}

	indexFilePath := filepath.Join(userdir, ".index.pkm")
	indexFileFS, err := os.OpenFile(indexFilePath, os.O_CREATE|os.O_RDWR, 0644)
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
	if len(indexFile) > 4 && bytes.Equal(indexFile[:4], []byte("PKM\n")) {
		encryptedIndex := indexFile[4:]
		decryptedIndex, err := kp.Decrypt(encryptedIndex)
		if err == nil {
			if err := json.Unmarshal(decryptedIndex, &index); err != nil {
				return err
			}
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
	indexJson, err := json.Marshal(index)
	if err != nil {
		return err
	}

	encryptedIndex, err := kp.Encrypt(indexJson)
	if err != nil {
		return err
	}
	indexPayload := []byte("PKM\n")
	indexPayload = append(indexPayload, encryptedIndex...)

	if err := indexFileFS.Truncate(0); err != nil {
		return err
	}
	if _, err := indexFileFS.Seek(0, 0); err != nil {
		return err
	}
	if _, err := indexFileFS.Write(indexPayload); err != nil {
		return err
	}

	return nil
}

func (fileStore *Store) Load(noteLocation string, username string, kp *crypt.KeyProvider) (*Note, error) {
	fileDataPath := filepath.Join(fileStore.StoreLocation, username, noteLocation+".pkm")

	fileDataFS, err := os.OpenFile(fileDataPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fileDataFS.Close()
	fileData, err := io.ReadAll(fileDataFS)
	if err != nil {
		return nil, err
	}
	if len(fileData) < 4 || !bytes.Equal(fileData[:4], []byte("PKM\n")) {
		return nil, errors.New("note corrupted")
	}

	encrpyted := fileData[4:]
	jsonData, err := kp.Decrypt(encrpyted)
	if err != nil {
		return nil, err
	}
	note := Note{}
	if err := json.Unmarshal(jsonData, &note); err != nil {
		return nil, err
	}
	return &note, nil
}

func (fileStore *Store) Delete(noteLocation string, username string) error {
	fileLoc := filepath.Join(fileStore.StoreLocation, username, noteLocation+".pkm")
	return os.Remove(fileLoc)
}

func (fileStore *Store) Search(searchType string, terms []string, username string, kp *crypt.KeyProvider) ([]string, error) {
	indexPath := filepath.Join(fileStore.StoreLocation, username, ".index.pkm")
	fileDataFS, err := os.OpenFile(indexPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fileDataFS.Close()
	fileData, err := io.ReadAll(fileDataFS)
	if err != nil {
		return nil, err
	}
	if len(fileData) < 4 || !bytes.Equal(fileData[:4], []byte("PKM\n")) {
		return nil, errors.New("index corrupted")
	}

	encrpyted := fileData[4:]
	jsonData, err := kp.Decrypt(encrpyted)
	if err != nil {
		return nil, err
	}
	var index Index
	if err := json.Unmarshal(jsonData, &index); err != nil {
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

func (fileStore *Store) List(username string, kp *crypt.KeyProvider) ([]NoteSummary, error) {

	fileDirPath := filepath.Join(fileStore.StoreLocation, username)
	fileDirFS, err := os.ReadDir(fileDirPath)
	if err != nil {
		return nil, err
	}

	var noteSummaryList []NoteSummary

	for _, file := range fileDirFS {

		if file.IsDir() {
			continue
		}
		name := file.Name()

		if !strings.HasSuffix(name, ".pkm") || name == ".index.pkm" {
			continue
		}

		fileDataPath := filepath.Join(fileDirPath, name)
		fileDataFS, err := os.OpenFile(fileDataPath, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}

		fileData, err := io.ReadAll(fileDataFS)
		if err != nil {
			return nil, err
		}

		if err := fileDataFS.Close(); err != nil {
			return nil, err
		}

		if len(fileData) < 4 || !bytes.Equal(fileData[:4], []byte("PKM\n")) {
			// Skip corrupted notes
			continue
		}

		encrpyted := fileData[4:]
		jsonData, err := kp.Decrypt(encrpyted)
		if err != nil {
			// Skip notes that fail to decrypt
			// continue
			return nil, err
		}

		var note Note
		if err := json.Unmarshal(jsonData, &note); err != nil {
			return nil, err
		}

		noteSummaryList = append(noteSummaryList, NoteSummary{
			Id:    note.Id,
			Title: note.Title,
			Tags:  note.Tags,
		})
	}

	sort.Slice(noteSummaryList, func(i, j int) bool {
		return noteSummaryList[i].Title < noteSummaryList[j].Title
	})

	return noteSummaryList, nil
}
