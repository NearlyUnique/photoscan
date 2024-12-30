package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

	err := processFile(*filePath, *indent)
	if err != nil {
		mainLog.Error("processFile", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
func processFile(filePath string, indent bool) error {
	var err error
	var imgFile *os.File
	var metaData *exif.Exif
	var jsonByte []byte
	var jsonString string

	imgFile, err = os.Open(filePath)
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
