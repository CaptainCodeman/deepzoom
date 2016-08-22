# DeepZoom calculations for Go

[Deep Zoom](https://msdn.microsoft.com/en-us/library/cc645077(v=vs.95).aspx)
is a technology developed by Microsoft for efficiently transmitting and viewing
images. It allows users to pan around and zoom in a large, high resolution image 
with the device only downloading the parts of the image visible in the viewport
at the resolution necessary. It allows you to view multi-Mb high-res images even
on a mobile device with limited bandwidth while still being able to zoom in to 
the full 1:1 resolution.

[OpenSeadragon](https://openseadragon.github.io/) is an open-source viewer
that you can use in your own app and demonstrates how it works.

This package provides the image slicing calculations necessary to generate a
[deepzoom tiles source](https://openseadragon.github.io/examples/tilesource-dzi/)

It does not generate the images themselves (yet).