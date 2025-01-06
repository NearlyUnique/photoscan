package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"io"
	"os"
)

func processFile(file *FileOpener, indent bool) error {
	var err error
	var imgFile *os.File
	var metaData *exif.Exif
	var jsonByte []byte
	var jsonString string

	imgFile, err = file.Open()
	if err != nil {
		return fmt.Errorf("open_file:(%s) %w", file.Filename, err)
	}

	metaData, err = exif.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("exif_decode:%w", err)
	}

	walk := NewDecoderWalker()
	err = metaData.Walk(walk)
	if err != nil {
		return fmt.Errorf("walk:%w", err)
	}
	if len(walk.errors) > 0 {
		walk.data["_errors"] = walk.errors
	}
	walk.data["FilePath"] = file.Filename
	var h string
	h, err = hash(file)
	if err != nil {
		return fmt.Errorf("file_hash:%w", err)
	}
	walk.data["FileHash"] = h
	if indent {
		jsonByte, err = json.MarshalIndent(walk.data, "", "  ")
	} else {
		jsonByte, err = json.Marshal(walk.data)
	}
	if err != nil {
		return fmt.Errorf("json_marshal:%w", err)
	}

	jsonString = string(jsonByte)
	fmt.Println(jsonString)

	return nil
}
func hash(file *FileOpener) (string, error) {

	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("file_open:%w", err)
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("calc_sha256:%w", err)
	}
	buf := h.Sum(nil)
	hashStr := hex.EncodeToString(buf)
	return hashStr, nil
}
