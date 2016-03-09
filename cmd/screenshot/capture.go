/*
Copyright (C) 2016 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const bytesPerPixel = 4

type frameCapture struct {
	taps  [][]byte
	image *image.RGBA

	id, numTaps,
	stride int
}

func newFrameCapture(id, numTaps, stride int) *frameCapture {
	taps := make([][]byte, numTaps)
	return &frameCapture{id: id, taps: taps, numTaps: numTaps, stride: stride}
}

func (fc *frameCapture) isComplete() bool {
	return len(fc.taps) == fc.numTaps
}

func (fc *frameCapture) scale() int {
	return int(math.Sqrt(float64(fc.numTaps)))
}

func (fc *frameCapture) geometry() (int, int) {
	scale := fc.scale()
	return scale * fc.stride, scale * len(fc.taps[0]) / (bytesPerPixel * fc.stride)
}

func (fc *frameCapture) addTap(tap int, data []byte) error {
	if fc.taps[tap] != nil {
		return fmt.Errorf("tap already written: %v", tap)
	}

	fc.taps[tap] = data
	return nil
}

func (fc *frameCapture) writeLine(line int, data []byte) {
	var c color.RGBA
	reader := bytes.NewReader(data)
	width := len(data) / 4

	for x := 0; x < width; x++ {
		if err := binary.Read(reader, binary.BigEndian, &c); err != nil {
			panic(err)
		}
		fc.image.SetRGBA(x, line, c)
	}
}

func (fc *frameCapture) compose() {
	width, height := fc.geometry()
	rect := image.Rect(0, 0, width, height)
	fc.image = image.NewRGBA(rect)

	scale := fc.scale()
	lineData := make([]byte, width*bytesPerPixel)

	for line := 0; line < height; line++ {
		vtap := line % scale
		y := line / scale

		for htap := 0; htap < scale; htap++ {
			tap := fc.taps[vtap*scale+htap]

			for x := 0; x < fc.stride; x++ {
				destOffset := x*scale + htap
				srcOffset := y*fc.stride + x

				for i := 0; i < bytesPerPixel; i++ {
					lineData[destOffset+i] = tap[srcOffset+i]
				}
			}
		}

		fc.writeLine(line, lineData)
	}
}

func (fc *frameCapture) save(path string) error {
	if !fc.isComplete() {
		return errors.New("capture is incompleat")
	}

	fc.compose()
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	return png.Encode(fp, fc.image)
}
