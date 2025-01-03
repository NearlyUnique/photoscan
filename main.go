package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	filePath := flag.String("file", "", "file o read")
	indent := flag.Bool("pretty", true, "pretty json output")
	debug := flag.Bool("debug", false, "print extra info")
	flag.Parse()
	if *filePath == "" {
		fmt.Println("missing -file")
		return
	}

	err := recursiveFileProcessor(ProcessorConfig{
		RootPath:   *filePath,
		IndentJSON: *indent,
		Debug:      *debug,
	}, logger)

	if err != nil {
		logger.Error("main", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
