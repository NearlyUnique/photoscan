package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

var Version = "no_version"
var Build = "no_build"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	filePath := flag.String("file", "", "file o read")
	indent := flag.Bool("pretty", true, "pretty json output")
	maxError := flag.Int("max-error", 100, "pretty json output")
	debug := flag.Bool("debug", false, "print extra info")
	version := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *version {
		fmt.Printf("version:%v build:%v\n", Version, Build)
		return
	}

	if *filePath == "" {
		fmt.Println("missing -file")
		return
	}

	err := recursiveFileProcessor(ProcessorConfig{
		RootPath:   *filePath,
		IndentJSON: *indent,
		Debug:      *debug,
		MaxError:   *maxError,
	}, logger)

	if err != nil {
		logger.Error("main", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
