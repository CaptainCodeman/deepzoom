package deepzoom

import (
	"fmt"
	"image"
	"math"
)

type (
	// DeepZoom represents a Deep Zoom image hierarchy of layers and tiles
	DeepZoom struct {
		// Width of the image in pixels
		Width int
		// Height of the image in pixels
		Height int
		// Size of each tile in pixels
		Size int
		// Overlap tiles by number of pixels
		Overlap int
	}

	// Layer represents a layer within a Deep Zoom hierarchy
	Layer struct {
		*DeepZoom
		// Level is the depth of the layer
		Level int
		// Scale is the ratio to resize the original image by for this layer
		Scale float64
	}

	// Tile identifies a specific image tile within a Layer
	Tile struct {
		*Layer
		// Col is the column number of this tile (X coordinate)
		Col int
		// Row is the row number of this tile (Y coordinate)
		Row int
	}
)

// New creates a new DeepZoom struct for an image with the given
// height and width. The tile size and tile overlap must match the
// values expected by the viewer.
func New(width, height, size, overlap int) *DeepZoom {
	return &DeepZoom{
		Width:   width,
		Height:  height,
		Size:    size,
		Overlap: overlap,
	}
}

func (dz *DeepZoom) maxDimension() float64 {
	return math.Max(float64(dz.Width), float64(dz.Height))
}

// MinLevel returns the minimum level that is the complete image
// Levels below this are just smaller scale versions of the full image.
func (dz *DeepZoom) MinLevel() int {
	return dz.MaxLevel() - int(math.Ceil(math.Log2(dz.maxDimension()/float64(dz.Size))))
}

// MaxLevel returns the maximum level corresponding to 1:1 resolution.
// Levels beyond would just be scaling up the image which adds nothing.
func (dz *DeepZoom) MaxLevel() int {
	return int(math.Ceil(math.Log2(dz.maxDimension())))
}

// Layer returns the layer for the requested level
func (dz *DeepZoom) Layer(level int) (*Layer, error) {
	min := dz.MinLevel()
	max := dz.MaxLevel()

	if level < min || level > max {
		return nil, fmt.Errorf("invalid level")
	}

	scale := math.Pow(0.5, float64(max-level))

	layer := &Layer{
		DeepZoom: dz,
		Level:    level,
		Scale:    scale,
	}

	return layer, nil
}

// Bounds returns the image bounds for the layer
// This is what the original image needs to be scaled to
func (l *Layer) Bounds() image.Rectangle {
	width := int(math.Ceil(float64(l.Width) * l.Scale))
	height := int(math.Ceil(float64(l.Height) * l.Scale))

	return image.Rect(0, 0, width, height)
}

// Dimensions returns the number of columns and rows for this level
func (l *Layer) Dimensions() (int, int) {
	r := l.Bounds()

	cols := int(math.Ceil(float64(r.Dx()) / float64(l.Size)))
	rows := int(math.Ceil(float64(r.Dy()) / float64(l.Size)))

	return cols, rows
}

// Tile returns the tile pointer for the given col and row within this layer
func (l *Layer) Tile(col, row int) (*Tile, error) {
	cols, rows := l.Dimensions()
	if col < 0 || col > cols {
		return nil, fmt.Errorf("invalid col")
	}
	if row < 0 || row > rows {
		return nil, fmt.Errorf("invalid row")
	}

	tile := &Tile{
		Layer: l,
		Col:   col,
		Row:   row,
	}

	return tile, nil
}

// Bounds returns the image bounding box for the tile relative to the layer image
func (t *Tile) Bounds() image.Rectangle {
	x1 := t.Col * t.Size
	y1 := t.Row * t.Size
	x2 := x1 + t.Size - 1 + t.Overlap
	y2 := y1 + t.Size - 1 + t.Overlap

	if t.Col > 0 {
		x1 -= t.Overlap
	}
	if t.Row > 0 {
		y1 -= t.Overlap
	}

	r := t.Layer.Bounds()
	if x2 >= r.Dx() {
		x2 = r.Dx() - 1
	}

	if y2 >= r.Dy() {
		y2 = r.Dy() - 1
	}

	return image.Rect(x1, y1, x2, y2)
}

// CropScale returns the crop and scale values relative to the source image
// for image decoders that can resize and crop in a single operation (e.g. WebP)
func (t *Tile) CropScale() (image.Rectangle, image.Rectangle) {
	x1 := t.Col * t.Size
	y1 := t.Row * t.Size
	x2 := x1 + t.Size - 1 + t.Overlap
	y2 := y1 + t.Size - 1 + t.Overlap

	if t.Col > 0 {
		x1 -= t.Overlap
	}
	if t.Row > 0 {
		y1 -= t.Overlap
	}

	r := t.Layer.Bounds()
	if x2 >= r.Dx() {
		x2 = r.Dx() - 1
	}

	if y2 >= r.Dy() {
		y2 = r.Dy() - 1
	}

	w := x2 - x1 + 1
	h := y2 - y1 + 1

	x1 = int(math.Ceil(float64(x1) / t.Scale))
	y1 = int(math.Ceil(float64(y1) / t.Scale))
	x2 = int(math.Ceil(float64(x2) / t.Scale))
	y2 = int(math.Ceil(float64(y2) / t.Scale))

	if x2 >= t.DeepZoom.Width {
		x2 = t.DeepZoom.Width - 1
	}

	if y2 >= t.DeepZoom.Height {
		y2 = t.DeepZoom.Height - 1
	}

	return image.Rect(x1, y1, x2, y2), image.Rect(0, 0, w, h)
}
