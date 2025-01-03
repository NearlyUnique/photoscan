package main

import (
	"io/fs"
	"log/slog"
	"path"
	"path/filepath"
	"strings"
)

type ProcessorConfig struct {
	RootPath   string
	IndentJSON bool
	Debug      bool
}

func shouldIgnoreFile(name string) bool {
	ext := strings.ToLower(path.Ext(name))
	switch ext {
	case ".jpg", ".jpeg", ".jpe", ".jif", ".jfif", ".jfi":
		return false
	}
	return true
}
func shouldSkipDir(name string) error {
	if len(name) >= 2 && name[0] == '.' {
		return filepath.SkipDir
	}
	ext := strings.ToLower(path.Ext(name))
	switch ext {
	case "node_modules":
		return filepath.SkipDir
	}
	return nil
}
func recursiveFileProcessor(config ProcessorConfig, logger *slog.Logger) error {
	err := filepath.WalkDir(config.RootPath, func(path string, d fs.DirEntry, inErr error) error {
		if inErr != nil {
			return inErr
		}
		if d.IsDir() {
			skip := shouldSkipDir(d.Name())
			if config.Debug && skip != nil {
				logger.With(slog.String("directory_name", d.Name())).Info("dir_skipped")
			}
			return skip
		}
		if shouldIgnoreFile(d.Name()) {
			if config.Debug {
				logger.With(slog.String("file_name", d.Name())).Info("file_skipped")
			}
			return nil
		}
		filename := d.Name()
		loopLogger := logger.With("file", filename)
		file := &FileOpener{Filename: filename}
		defer func() {
			if err := file.Close(); err != nil {
				loopLogger.Error("close_file", slog.String("error", err.Error()))
			}
		}()

		err := processFile(file, config.IndentJSON)
		if config.Debug && err != nil {
			logger.
				With(
					slog.String("error", d.Name()),
					slog.String("file_name", d.Name()),
				).
				Debug("process_file_error")
		}
		return err
	})
	return err
}
