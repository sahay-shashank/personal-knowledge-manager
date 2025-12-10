package note

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
)

func InitStore(storeDirectory string) *Store {
	return &Store{
		StoreLocation: storeDirectory,
	}
}

func (fileStore *Store) Save(note *Note) error {
	noteFile, err := os.Create(fileStore.StoreLocation + "/" + note.Id + ".pkm")
	if err != nil {
		return err
	}
	jsonBody, err := json.Marshal(note)
	if err != nil {
		return err
	}

	payload := []byte("PKM\n")
	payload = append(payload, jsonBody...)
	noteFile.Write(payload)

	noteFile.Close()
	return nil
}

func (fileStore *Store) Load(noteLocation string) (*Note, error) {
	fileData, err := os.ReadFile(fileStore.StoreLocation + "/" + noteLocation + ".pkm")
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
