package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	filePath := flag.String("file", "", "file o read")
	indent := flag.Bool("pretty", true, "pretty json output")
	flag.Parse()
	if *filePath == "" {
		fmt.Println("missing -file")
		return
	}
	mainLog := logger.With("file", *filePath)
	file := &FileOpener{Filename: *filePath}
	defer func() {
		if err := file.Close(); err != nil {
			mainLog.Error("close_file", slog.String("error", err.Error()))
		}
	}()

	err := processFile(file, *indent)
	if err != nil {
		mainLog.Error("processFile", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
func processFile(file *FileOpener, indent bool) error {
	var err error
	var imgFile *os.File
	var metaData *exif.Exif
	var jsonByte []byte
	var jsonString string

	imgFile, err = file.Open()
	if err != nil {
		return fmt.Errorf("open_file:%w", err)
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
	var h string
	h, err = hash(file)
	if err != nil {
		return fmt.Errorf("file_hash:%w", err)
	}
	walk.data["file_hash"] = string(h)
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
