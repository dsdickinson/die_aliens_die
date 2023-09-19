package main

import (
	"image"
	"os"
	_ "image/png"
	"github.com/faiface/pixel"
)

var (
	as actionSquare

	pos, dir, plane pixel.Vec
)

func minimap() *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, 24, 26))

	for x, row := range world {
		for y, _ := range row {
			c := getColor(x, y)
			if c.A == 255 {
				c.A = 96
			}
			m.Set(x, y, c)
		}
	}

	m.Set(int(pos.X), int(pos.Y), color.RGBA{255, 0, 0, 255})

	if as.active {
		m.Set(as.X, as.Y, color.RGBA{255, 255, 255, 255})
	} else {
		m.Set(as.X, as.Y, color.RGBA{64, 64, 64, 255})
	}

	return m
}
