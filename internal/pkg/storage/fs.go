package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
)

type storageFile struct {
	Inner map[string]Value    `json:"inner"`
	List  map[string][]string `json:"list"`
}

func (s *Storage) SaveToFile(path string) error {

	data := storageFile{
		Inner: s.inner,
		List:  make(map[string][]string),
	}

	for k, v := range s.list {
		data.List[k] = v.elem
	}

	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}

	s.logger.Info("Storage saved to file", zap.String("file", path))

	return nil

}

func (s *Storage) LoadFromFile(path string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var data storageFile
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}

	s.inner = data.Inner
	s.list = make(map[string]*List)

	for k, v := range data.List {
		s.list[k] = &List{elem: v}
	}

	fmt.Println(jsonData)
	s.logger.Info("Stotage loaded fron file", zap.String("file", path))
	return nil
}
