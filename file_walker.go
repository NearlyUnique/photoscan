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
	ext := strings.ToLower(name)
	switch ext {
	case "@eadir", // NAS meta folders
		"node_modules": // node development
		return filepath.SkipDir
	}
	return nil
}
func recursiveFileProcessor(config ProcessorConfig, logger *slog.Logger) error {
	return filepath.WalkDir(config.RootPath, (&walkWithState{
		logger: logger,
		config: config,
	}).walk)
}

type walkWithState struct {
	logger     *slog.Logger
	errorCount int
	config     ProcessorConfig
}

func (w *walkWithState) walk(path string, d fs.DirEntry, inErr error) error {
	if inErr != nil {
		return inErr
	}
	walkLog := w.logger.With(slog.String("full_path", path))

	if err, skip := w.shouldSkip(d, walkLog); skip {
		return err
	}

	file := &FileOpener{Filename: d.Name()}
	defer func() {
		if err := file.Close(); err != nil {
			walkLog.
				With(slog.String("file_name", d.Name())).
				Error("close_file_error", slog.String("error", err.Error()))
		}
	}()

	err := processFile(file, w.config.IndentJSON)
	if err != nil {
		w.errorCount++
		if w.config.Debug {
			walkLog.
				With(
					slog.String("error", err.Error()),
					slog.String("file_name", d.Name()),
					slog.String("full_path", file.Filename),
				).
				Error("process_file_error")
		}
		if w.suppressError() {
			err = nil
		}
	}
	return err
}

func (w *walkWithState) shouldSkip(d fs.DirEntry, walkLog *slog.Logger) (error, bool) {
	if d.IsDir() {
		skipErr := shouldSkipDir(d.Name())
		if w.config.Debug && skipErr != nil {
			walkLog.
				With(slog.String("directory_name", d.Name())).
				Info("dir_skipped")
		}
		return skipErr, true
	}
	if shouldIgnoreFile(d.Name()) {
		if w.config.Debug {
			walkLog.
				With(slog.String("file_name", d.Name())).
				Info("file_skipped")
		}
		return nil, true
	}
	return nil, false
}

func (w *walkWithState) suppressError() bool {
	return w.errorCount < w.config.MaxError || w.config.MaxError == -1
}
