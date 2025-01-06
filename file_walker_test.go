package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_filter_unimportant_folders(t *testing.T) {
	var skip error = filepath.SkipDir
	var noSkip error = nil
	testData := map[string]error{
		".anything":     skip,
		"node_modules":  skip,
		"@eaDir":        skip,
		"normal_folder": noSkip,
		".":             noSkip,
	}
	for filename, expected := range testData {
		t.Run(fmt.Sprintf("skip check on '%s'", filename), func(t *testing.T) {
			assert.Equal(t, expected, shouldSkipDir(filename))
		})
	}
}
