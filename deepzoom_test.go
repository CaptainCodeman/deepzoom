package deepzoom

import (
	"fmt"
	"log"
	"math"
	"testing"
)

func TestDeepZoomLevels(t *testing.T) {
	tests := []struct {
		width  int
		height int
		size   int
		min    int
		max    int
	}{
		{800, 600, 32, 5, 10},
		{800, 600, 64, 6, 10},
		{800, 600, 128, 7, 10},
		{800, 600, 256, 8, 10},

		{8000, 6000, 32, 5, 13},
		{8000, 6000, 64, 6, 13},
		{8000, 6000, 128, 7, 13},
		{8000, 6000, 256, 8, 13},

		{10000, 10000, 32, 5, 14},
		{10000, 10000, 64, 6, 14},
		{10000, 10000, 128, 7, 14},
		{10000, 10000, 256, 8, 14},
	}
	for i, test := range tests {
		dz := New(test.width, test.height, test.size, 0)
		if min := dz.MinLevel(); min != test.min {
			t.Errorf("%d expected min %d levels, got %d", i, test.min, min)
		}
		if max := dz.MaxLevel(); max != test.max {
			t.Errorf("%d expected max %d levels, got %d", i, test.max, max)
		}
	}
}

func TestLayers(t *testing.T) {
	type layerTest struct {
		level  int
		width  int
		height int
		cols   int
		rows   int
		scale  float64
	}
	tests := []struct {
		width  int
		height int
		size   int
		layers []layerTest
	}{
		{
			8000, 6000, 256, []layerTest{
				{13, 8000, 6000, 32, 24, 1.0},
				{12, 4000, 3000, 16, 12, 0.5},
				{11, 2000, 1500, 8, 6, 0.25},
				{10, 1000, 750, 4, 3, 0.125},
				{9, 500, 375, 2, 2, 0.0625},
				{8, 250, 188, 1, 1, 0.03125},
			},
		},
		{
			4200, 2800, 128, []layerTest{
				{13, 4200, 2800, 33, 22, 1.0},
				{12, 2100, 1400, 17, 11, 0.5},
				{11, 1050, 700, 9, 6, 0.25},
				{10, 525, 350, 5, 3, 0.125},
				{9, 263, 175, 3, 2, 0.0625},
				{8, 132, 88, 2, 1, 0.03125},
				{7, 66, 44, 1, 1, 0.015625},
			},
		},
	}
	for i, test := range tests {
		dz := New(test.width, test.height, test.size, 1)

		for j, testLayer := range test.layers {
			layer, _ := dz.Layer(testLayer.level)
			r := layer.Bounds()
			if r.Dx() != testLayer.width || r.Dy() != testLayer.height {
				t.Errorf("%d.%d expected %d x %d got %dx%d", i, j, testLayer.width, testLayer.height, r.Dx(), r.Dy())
			}
			cols, rows := layer.Dimensions()
			if cols != testLayer.cols || rows != testLayer.rows {
				t.Errorf("%d.%d expected %d cols %d rows got %d cols %d rows", i, j, testLayer.cols, testLayer.rows, cols, rows)
			}
			if layer.Scale != testLayer.scale {
				t.Errorf("%d expected scale %f got %f", j, testLayer.scale, layer.Scale)
			}
		}
	}
}

func TestTileBounds(t *testing.T) {
	tests := []struct {
		col int
		row int
		x1  int
		y1  int
		x2  int
		y2  int
		w   int
		h   int
	}{
		{0, 0, 0, 0, 256, 256, 257, 257},
		{0, 1, 0, 255, 256, 512, 257, 258},
		{0, 2, 0, 511, 256, 768, 257, 258},
		{0, 3, 0, 767, 256, 791, 257, 25},
		{1, 0, 255, 0, 512, 256, 258, 257},
		{1, 1, 255, 255, 512, 512, 258, 258},
		{1, 2, 255, 511, 512, 768, 258, 258},
		{1, 3, 255, 767, 512, 791, 258, 25},
		{2, 0, 511, 0, 768, 256, 258, 257},
		{2, 1, 511, 255, 768, 512, 258, 258},
		{2, 2, 511, 511, 768, 768, 258, 258},
		{2, 3, 511, 767, 768, 791, 258, 25},
		{3, 0, 767, 0, 1024, 256, 258, 257},
		{3, 1, 767, 255, 1024, 512, 258, 258},
		{3, 2, 767, 511, 1024, 768, 258, 258},
		{3, 3, 767, 767, 1024, 791, 258, 25},
		{4, 0, 1023, 0, 1055, 256, 33, 257},
		{4, 1, 1023, 255, 1055, 512, 33, 258},
		{4, 2, 1023, 511, 1055, 768, 33, 258},
		{4, 3, 1023, 767, 1055, 791, 33, 25},
	}

	dz := New(4224, 3168, 256, 1)
	level := 11
	for _, test := range tests {
		layer, _ := dz.Layer(level)
		tile, _ := layer.Tile(test.col, test.row)
		r := tile.Bounds()
		if test.x1 != r.Min.X || test.y1 != r.Min.Y || test.x2 != r.Max.X || test.y2 != r.Max.Y {
			fmt.Printf("col %d row %d expected %d,%d - %d,%d got %d,%d - %d,%d\n", test.col, test.row, test.x1, test.y1, test.x2, test.y2, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
		}
	}
}

func TestTileCropScale(t *testing.T) {
	t.Skip()
	tests := []struct {
		col int
		row int
		x1  int
		y1  int
		x2  int
		y2  int
		w   int
		h   int
	}{
		{0, 0, 0, 0, 256, 256, 257, 257},
		{0, 1, 0, 255, 256, 512, 257, 258},
		{0, 2, 0, 511, 256, 768, 257, 258},
		{0, 3, 0, 767, 256, 791, 257, 25},
		{1, 0, 255, 0, 512, 256, 258, 257},
		{1, 1, 255, 255, 512, 512, 258, 258},
		{1, 2, 255, 511, 512, 768, 258, 258},
		{1, 3, 255, 767, 512, 791, 258, 25},
		{2, 0, 511, 0, 768, 256, 258, 257},
		{2, 1, 511, 255, 768, 512, 258, 258},
		{2, 2, 511, 511, 768, 768, 258, 258},
		{2, 3, 511, 767, 768, 791, 258, 25},
		{3, 0, 767, 0, 1024, 256, 258, 257},
		{3, 1, 767, 255, 1024, 512, 258, 258},
		{3, 2, 767, 511, 1024, 768, 258, 258},
		{3, 3, 767, 767, 1024, 791, 258, 25},
		{4, 0, 1023, 0, 1055, 256, 33, 257},
		{4, 1, 1023, 255, 1055, 512, 33, 258},
		{4, 2, 1023, 511, 1055, 768, 33, 258},
		{4, 3, 1023, 767, 1055, 791, 33, 25},
	}

	dz := New(4224, 3168, 256, 1)
	level := 11
	for _, test := range tests {
		layer, _ := dz.Layer(level)
		tile, _ := layer.Tile(test.col, test.row)
		c, s := tile.CropScale()
		fmt.Printf("col %d row %d %d,%d - %d,%d (%d x %d)\n", test.col, test.row, c.Min.X, c.Min.Y, c.Max.X, c.Max.Y, s.Dx(), s.Dy())
		/*
			if test.x1 != c.Min.X || test.y1 != c.Min.Y || test.x2 != c.Max.X || test.y2 != c.Max.Y {
				fmt.Printf("col %d row %d expected %d,%d - %d,%d got %d,%d - %d,%d\n", test.col, test.row, test.x1, test.y1, test.x2, test.y2, c.Min.X, c.Min.Y, c.Max.X, c.Max.Y)
			}
		*/
	}
}

func TestX(t *testing.T) {
	dz := New(4256, 2832, 256, 1)
	for level := dz.MinLevel(); level < dz.MaxLevel(); level++ {
		layer, _ := dz.Layer(level)
		cols, rows := layer.Dimensions()
		fmt.Printf("level %d, %dx%d tiles\n", level, cols, rows)

		for col := 0; col < cols; col++ {
			for row := 0; row < rows; row++ {
				tile, _ := layer.Tile(col, row)
				r := tile.Bounds()
				fmt.Printf("level %d, %d:%d %v (%dx%d)\n", level, col, row, r, r.Dx(), r.Dy())
			}
		}
	}
}

func TestY(t *testing.T) {
	var count int
	dz := New(10240, 6400, 254, 1)
	fmt.Printf("levels min %d max %d\n", dz.MinLevel(), dz.MaxLevel())
	for level := dz.MinLevel(); level <= dz.MaxLevel(); level++ {
		layer, _ := dz.Layer(level)
		r := layer.Bounds()
		cols, rows := layer.Dimensions()
		fmt.Printf("level %d, %dx%d tiles, size %dx%d, scale %f\n", level, cols, rows, r.Dx(), r.Dy(), layer.Scale)

		for col := 0; col < cols; col++ {
			for row := 0; row < rows; row++ {
				tile, _ := layer.Tile(col, row)
				r := tile.Bounds()

				crop, scale := tile.CropScale()
				fmt.Printf("level %d, %d:%d %v (%dx%d) %v %v\n", level, col, row, r, r.Dx()+1, r.Dy()+1, crop, scale)

				count++
			}
		}
	}

	fmt.Printf("total %d tiles\n", count)
}

func TestCalc(t *testing.T) {
	tests := []struct {
		width  int
		height int
		size   int
	}{
		{1024, 1024, 254},
		{608, 798, 254},
		{2995, 4493, 254},
		{2400, 3600, 254},
		{10000, 10000, 254},
		{254, 254, 254},
		{6000, 4000, 128},
		{6000, 4000, 16},
	}

	for _, test := range tests {
		w := test.width
		h := test.height

		maxDimension := math.Max(float64(test.width), float64(test.height))

		maxLevel := int(math.Ceil(math.Log2(maxDimension)))
		minLevel := maxLevel - int(math.Ceil(math.Log2(maxDimension/float64(test.size))))

		log.Printf("%d x %d, %d, levels %d - %d", test.width, test.height, test.size, minLevel, maxLevel)

		for level := maxLevel; level >= minLevel; level-- {
			scale := math.Pow(0.5, float64(maxLevel-level))
			lw := int(math.Ceil(float64(test.width) * scale))
			lh := int(math.Ceil(float64(test.height) * scale))
			cols := int(math.Ceil(float64(w) / float64(test.size)))
			rows := int(math.Ceil(float64(h) / float64(test.size)))
			log.Printf("level %d / %d: %d x %d, %f %d x %d, %d x %d", level, level-minLevel, w, h, scale, lw, lh, cols, rows)
			w = (w + 1) >> 1
			h = (h + 1) >> 1
		}
	}
}
