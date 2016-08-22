package deepzoom

import (
	"fmt"
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
		{800, 600, 32, 6, 10},
		{800, 600, 64, 7, 10},
		{800, 600, 128, 8, 10},
		{800, 600, 256, 9, 10},

		{8000, 6000, 32, 6, 13},
		{8000, 6000, 64, 7, 13},
		{8000, 6000, 128, 8, 13},
		{8000, 6000, 256, 9, 13},

		{10000, 10000, 32, 6, 14},
		{10000, 10000, 64, 7, 14},
		{10000, 10000, 128, 8, 14},
		{10000, 10000, 256, 9, 14},
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
	}
	tests := []struct {
		width  int
		height int
		size   int
		layers []layerTest
	}{
		{8000, 6000, 256, []layerTest{
			{13, 8000, 6000, 32, 24},
			{12, 4000, 3000, 16, 12},
			{11, 2000, 1500, 8, 6},
			{10, 1000, 750, 4, 3},
			{9, 500, 375, 2, 2}},
		},
		{4200, 2800, 128, []layerTest{
			{13, 4200, 2800, 33, 22},
			{12, 2100, 1400, 17, 11},
			{11, 1050, 700, 9, 6},
			{10, 525, 350, 5, 3},
			{9, 263, 175, 3, 2},
			{8, 132, 88, 2, 1}},
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
		}
	}
}

func TestTiles(t *testing.T) {
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
