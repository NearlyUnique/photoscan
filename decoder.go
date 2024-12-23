package main

import (
	"errors"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type decoderWalker struct {
	data   map[string]any
	errors []string
}

func (d *decoderWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	var err error
	switch tag.Type {
	case tiff.DTShort:
		d.data[string(name)], err = tag.Int(0)
	case tiff.DTRational:
		if tag.Count != 1 {
			d.errors = append(d.errors, fmt.Sprintf("%s (Count=%d)", name, tag.Count))
		}
		var (
			num, den int64
		)
		num, den, err = tag.Rat2(0)
		if den == 0 {
			d.errors = append(d.errors, fmt.Sprintf("%s (Div0,num=%d)", name, num))
		}
		d.data[string(name)] = float64(num) / float64(den)
	default:
		err = errors.New("not implemented")
	}
	return err
}

func (d *decoderWalker) decode(data *exif.Exif) (map[string]any, error) {
	d.data = make(map[string]any)
	err := data.Walk(d)
	return d.data, err
}
