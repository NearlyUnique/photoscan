package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type DecoderWalker struct {
	data   map[string]any
	errors []string
}

func (d *DecoderWalker) process(name string) {
	switch name {
	case "DateTime",
		"DateTimeDigitized",
		"DateTimeOriginal",
		"GPSDateStamp":
		v, ok := d.data[name].(string)
		if ok && len(v) > 6 {
			buf := []byte(v)
			buf[4] = '-'
			buf[7] = '-'
			v = string(buf)
			d.data[name] = v
		}
		break
	case "GPSTimeStamp":
		v, ok := d.data[name].([]float64)
		if ok && len(v) == 3 {
			d.data[name] =
				fmt.Sprintf("%02d:%02d:%02d", int(v[0]), int(v[1]), int(v[2]))
		}
	case "GPSLatitude",
		"GPSLongitude":
		v, ok := d.data[name].([]float64)
		if ok && len(v) == 3 {
			d.data[name] =
				fmt.Sprintf("%.0f\u00b0%.0f'%.4f\"", v[0], v[1], v[2])
		}
	}
}
func (d *DecoderWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	var err error

	switch tag.Type {
	case tiff.DTAscii:
		d.data[string(name)], err = tag.StringVal()
	case tiff.DTShort:
		d.data[string(name)], err = tag.Int(0)
	case tiff.DTByte:
		var v int
		v, err = tag.Int(0)
		d.data[string(name)] = byte(v)
	case tiff.DTSByte:
		var v int
		v, err = tag.Int(0)
		d.data[string(name)] = int8(v)
	case tiff.DTLong:
		var v1, v2 int64
		v1, err = tag.Int64(0)
		v2, err = tag.Int64(1)
		d.data[string(name)] = (v1 << 32) | v2
	case tiff.DTSRational:
		d.data[string(name)] = "?"
		d.errors = append(d.errors, fmt.Sprintf("%s (unknown tag=%v)", name, tag.Type))
	case tiff.DTRational:
		err = d.extractRational(name, tag)
	case tiff.DTUndefined:
		break
	default:
		d.errors = append(d.errors, fmt.Sprintf("%s (unknown tag=%v)", name, tag.Type))
		//err = errors.New("not implemented")
	}
	d.process(string(name))
	return err
}

func (d *DecoderWalker) extractRational(name exif.FieldName, tag *tiff.Tag) error {
	var num, den int64
	var err error

	if tag.Count == 1 {
		num, den, err = tag.Rat2(0)
		if den == 0 {
			d.errors = append(d.errors, fmt.Sprintf("%s (Div0,num=%d)", name, num))
		}
		result := float64(num) / float64(den)
		d.data[string(name)] = result
	} else {
		imax := int(tag.Count)
		var result []float64
		for i := 0; i < imax; i++ {
			num, den, err = tag.Rat2(i)
			if den == 0 {
				d.errors = append(d.errors, fmt.Sprintf("%s[%d] (Div0,num=%d)", name, i, num))
			}
			result = append(result, float64(num)/float64(den))
		}
		d.data[string(name)] = result
	}
	return err
}
func NewDecoderWalker() *DecoderWalker {
	d := DecoderWalker{
		data: make(map[string]any),
	}
	return &d
}
func (d *DecoderWalker) decode(data *exif.Exif) (map[string]any, error) {
	d.data = make(map[string]any)
	err := data.Walk(d)
	return d.data, err
}
