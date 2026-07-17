package koi

type Panel interface {
	Init() error
	Display(*Framebuffer) error
	Sleep() error
}
