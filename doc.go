/*
Package deepzoom provides calculations for generating Deep Zoom images

Deep Zoom is a technology developed by Microsoft for efficiently
transmitting and viewing images. It allows users to pan around and
zoom in a large, high resolution image or a large collection of images
with the device only downloading the parts of the image visible in the
viewport at the resolution necessary. It allows you to view multi-Mb
high-res images even on a mobile device with limited bandwidth while
still being able to zoom in to 1:1 resolution.

See: https://msdn.microsoft.com/en-us/library/cc645077(v=vs.95).aspx

For a demonstration of Deep Zoom and a web viewer that you can use in
your own app, visit https://openseadragon.github.io/ This package provides
the image slicing calculations necessary to provide the deepzoom tiles
source described at https://openseadragon.github.io/examples/tilesource-dzi/
*/

package deepzoom
