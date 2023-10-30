package file

import (
	"encoding/json"
	"github.com/JuliyaMS/service-metrics-alerting/internal/storage"
	"os"
)

type StorageFileEncode struct {
	file    *os.File
	encoder *json.Encoder
}

func NewStorageFileEncode(filename string) (*StorageFileEncode, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return &StorageFileEncode{file: f, encoder: json.NewEncoder(f)}, nil
}

func (f *StorageFileEncode) WriteToFile() error {
	return f.encoder.Encode(storage.Storage)
}

func (f *StorageFileEncode) Close() {
	f.file.Close()
}

type StorageFileDecode struct {
	file    *os.File
	encoder *json.Decoder
}

func NewStorageFileDecode(filename string) (*StorageFileDecode, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &StorageFileDecode{file: f, encoder: json.NewDecoder(f)}, nil
}

func (f *StorageFileDecode) ReadFromFile() error {
	return f.encoder.Decode(&storage.Storage)
}

func (f *StorageFileDecode) Close() {
	f.file.Close()
}
