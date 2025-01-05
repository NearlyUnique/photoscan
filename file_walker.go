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
	MaxError   int
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
	case "@eaDir", // NAS meta folders
		"node_modules": // node development
		return filepath.SkipDir
	}
	return nil
}
func recursiveFileProcessor(config ProcessorConfig, logger *slog.Logger) error {
	errorCount := 0
	err := filepath.WalkDir(config.RootPath, func(path string, d fs.DirEntry, inErr error) error {
		if inErr != nil {
			return inErr
		}
		walkLog := logger.With(slog.String("full_path", path))
		if d.IsDir() {
			skip := shouldSkipDir(d.Name())
			if config.Debug && skip != nil {
				walkLog.
					With(slog.String("directory_name", d.Name())).
					Info("dir_skipped")
			}
			return skip
		}
		if shouldIgnoreFile(d.Name()) {
			if config.Debug {
				walkLog.
					With(slog.String("file_name", d.Name())).
					Info("file_skipped")
			}
			return nil
		}
		file := &FileOpener{Filename: d.Name()}
		defer func() {
			if err := file.Close(); err != nil {
				walkLog.
					With(slog.String("file_name", d.Name())).
					Error("close_file_error", slog.String("error", err.Error()))
			}
		}()

		err := processFile(file, config.IndentJSON)
		if err != nil {
			errorCount++
			if config.Debug {
				walkLog.
					With(
						slog.String("error", d.Name()),
						slog.String("file_name", d.Name()),
					).
					Error("process_file_error")
			}
			if errorCount < config.MaxError || config.MaxError == -1 {
				err = nil
			}
		}
		return err
	})
	return err
}
