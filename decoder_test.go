package main

import (
	"bytes"
	"encoding/binary"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"testing"

	"github.com/stretchr/testify/assert"
)

type input struct {
	tagID      string
	dataType   string
	valueCount string
	offset     string
	val        string
}

const (
	InputDTByte      = "0001"
	InputDTAscii     = "0002"
	InputDTShort     = "0003"
	InputDTLong      = "0004"
	InputDTRational  = "0005"
	InputDTSInt32    = "0008"
	InputDTSByte     = "0006"
	InputDTUndefined = "0007"
	InputDTSShort    = "0008"
	InputDTSLong     = "0009"
	InputDTSRational = "000A"
	InputDTFloat     = "000B"
	InputDTDouble    = "000C"
)
const (
	InputValuesCount2 = "00000002"
)

func Test_input_values(t *testing.T) {
	w := DecoderWalker{}
	w.data = make(map[string]any)
	testData := map[string]struct {
		expected any
		data     []byte
	}{
		"a byte":         {byte(255), buildInput(InputDTByte, 255, "")},
		"a signed byte":  {int8(-8), buildInput(InputDTSByte, 248, "")},
		"32 bit integer": {10, buildInput(InputDTShort, 10, "")},
		"64 bit integer": {int64(0x01_23_45_67_8a_bc_de_f0), buildInput(InputDTLong, 2,
			I4(0x01234567, 0x8abcdef0))},
		"rational":        {0.25, buildInput(InputDTRational, 1, I4(25, 100))},
		"signed rational": {0.25, buildInput(InputDTSRational, 1, I4(25, 100))},
		"rational array":  {[]float64{0.25, 0.05, 7.5}, buildInput(InputDTRational, 3, I4(25, 100)+I4(50, 1000)+I4(75, 10))},
		"string":          {"abcdef", buildInput(InputDTAscii, 6, "abcdef")},
		// The following keys are "special"
		"DateTime": {"2019-10-11 19:05:41", buildInput(InputDTAscii, 0, "2019:10:11 19:05:41")},
		"GPSLatitude": {"52°12'30.3951\"", buildInput(InputDTRational, 3, I4(
			52, 1,
			12, 1,
			30000, 987))},
		"GPSTimeStamp": {"18:05:33", buildInput(InputDTRational, 3, I4(
			18, 1,
			5, 1,
			33, 1))},
	}
	for name, td := range testData {
		t.Run(name, func(t *testing.T) {
			raw := bytes.NewReader(td.data)
			tag, err := tiff.DecodeTag(raw, binary.BigEndian)
			assert.NoError(t, err)
			err = w.Walk(exif.FieldName(name), tag)
			assert.NoError(t, err)
			assert.Equal(t, td.expected, w.data[name])
			// 81985529234382576
		})
	}
}
