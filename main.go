package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	var err error
	var imgFile *os.File
	var metaData *exif.Exif
	var jsonByte []byte
	var jsonString string

	filePath := flag.String("file", "", "file o read")
	flag.Parse()
	if *filePath == "" {
		fmt.Println("missing -file")
		return
	}

	imgFile, err = os.Open(*filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	metaData, err = exif.Decode(imgFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	walk := NewDecoderWalker()
	err = metaData.Walk(walk)
	if err != nil {
		log.Fatal(err.Error())
	}
	jsonByte, err = json.Marshal(walk.data)
	if err != nil {
		log.Fatal(err.Error())
	}

	jsonString = string(jsonByte)
	fmt.Println(jsonString)
}
