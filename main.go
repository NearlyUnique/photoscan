package main

import (
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

	d := decoderWalker{}
	if v, err := d.decode(metaData); err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", v)
	}
	jsonByte, err = metaData.MarshalJSON()
	if err != nil {
		log.Fatal(err.Error())
	}

	jsonString = string(jsonByte)
	fmt.Println(jsonString)
	//
	//fmt.Println("Make: " + gjson.Get(jsonString, "Make").String())
	//fmt.Println("Model: " + gjson.Get(jsonString, "Model").String())
	//fmt.Println("Software: " + gjson.Get(jsonString, "Software").String())
}
