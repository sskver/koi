package koi

import (
	"log/slog"
)

type Framebuffer struct {
	data   []byte
	width  int
	height int
	stride int
}

func NewFramebuffer(width, height int) *Framebuffer {
	stride := (width + 7) / 8

	fb := &Framebuffer{
		data:   make([]byte, stride*height),
		width:  width,
		height: height,
		stride: stride,
	}

	fb.Clear(0xFF)

	return fb
}

func (fb *Framebuffer) Data() []byte {
	return fb.data
}

func (fb *Framebuffer) Width() int {
	return fb.width
}

func (fb *Framebuffer) Height() int {
	return fb.height
}

func (fb *Framebuffer) Stride() int {
	return fb.stride
}

func (fb *Framebuffer) Clear(val byte) {
	for i := range fb.data {
		fb.data[i] = val
	}
}

func (fb *Framebuffer) SetPixel(x, y int) {
	if x < 0 || y < 0 || x >= fb.width || y >= fb.height {
		slog.Warn("setPixel out of bounds", "x", x, "y", y, "width", fb.width, "height", fb.height)
		return
	}

	fb.data[y*fb.stride+x/8] |= 1 << (7 - (x % 8))
}

func (fb *Framebuffer) ClearPixel(x, y int) {
	if x < 0 || y < 0 || x >= fb.width || y >= fb.height {
		slog.Warn("clearPixel out of bounds", "x", x, "y", y, "width", fb.width, "height", fb.height)
		return
	}

	fb.data[y*fb.stride+x/8] &^= 1 << (7 - (x % 8))
}
