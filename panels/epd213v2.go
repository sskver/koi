package panels

import (
	"time"

	"github.com/sskver/koi"
)

const (
	Width  = 122
	Height = 250
)

type EPD213V2 struct {
	controller *koi.Controller
	fullUpdate bool
}

func NewEPD213V2(controller *koi.Controller, fullUpdate bool) *EPD213V2 {
	return &EPD213V2{
		controller: controller,
		fullUpdate: fullUpdate,
	}
}

var lutFullUpdate = []byte{
	0x80, 0x60, 0x40, 0x00, 0x00, 0x00, 0x00,
	0x10, 0x60, 0x20, 0x00, 0x00, 0x00, 0x00,
	0x80, 0x60, 0x40, 0x00, 0x00, 0x00, 0x00,
	0x10, 0x60, 0x20, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

	0x03, 0x03, 0x00, 0x00, 0x02,
	0x09, 0x09, 0x00, 0x00, 0x02,
	0x03, 0x03, 0x00, 0x00, 0x02,
	0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00,

	0x15, 0x41, 0xA8, 0x32, 0x30, 0x0A,
}

var lutPartialUpdate = []byte{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

	0x0A, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00,

	0x15, 0x41, 0xA8, 0x32, 0x30, 0x0A,
}

func (p *EPD213V2) Init() error {
	if err := p.controller.Reset(); err != nil {
		return err
	}

	if p.fullUpdate {
		if err := p.controller.WaitBusy(); err != nil {
			return err
		}

		// Software reset
		if err := p.controller.Command(0x12); err != nil {
			return err
		}

		if err := p.controller.WaitBusy(); err != nil {
			return err
		}

		// Analog block control
		if err := p.controller.Command(0x74); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{0x54}); err != nil {
			return err
		}

		// Digital block control
		if err := p.controller.Command(0x7E); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{0x3B}); err != nil {
			return err
		}

		// Driver output control
		if err := p.controller.Command(0x01); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{
			0xF9,
			0x00,
			0x00,
		}); err != nil {
			return err
		}

		// Data entry mode
		if err := p.controller.Command(0x11); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{0x01}); err != nil {
			return err
		}

		// RAM X address range
		if err := p.controller.Command(0x44); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{
			0x00,
			0x0F,
		}); err != nil {
			return err
		}

		// RAM Y address range
		if err := p.controller.Command(0x45); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{
			0xF9,
			0x00,
			0x00,
			0x00,
		}); err != nil {
			return err
		}

		// Border waveform
		if err := p.controller.Command(0x3C); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{0x03}); err != nil {
			return err
		}

		// VCOM voltage
		if err := p.controller.Command(0x2C); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{0x55}); err != nil {
			return err
		}

		// Gate voltage
		if err := p.controller.Command(0x03); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{lutFullUpdate[70]}); err != nil {
			return err
		}

		// Source voltage
		if err := p.controller.Command(0x04); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{
			lutFullUpdate[71],
			lutFullUpdate[72],
			lutFullUpdate[73],
		}); err != nil {
			return err
		}

		// Dummy line period
		if err := p.controller.Command(0x3A); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{lutFullUpdate[74]}); err != nil {
			return err
		}

		// Gate time
		if err := p.controller.Command(0x3B); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{lutFullUpdate[75]}); err != nil {
			return err
		}

		// Load LUT
		if err := p.controller.Command(0x32); err != nil {
			return err
		}

		if err := p.controller.Data(lutFullUpdate[:70]); err != nil {
			return err
		}

		// Set RAM address counter X
		if err := p.controller.Command(0x4E); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{0x00}); err != nil {
			return err
		}

		// Set RAM address counter Y
		if err := p.controller.Command(0x4F); err != nil {
			return err
		}

		if err := p.controller.Data([]byte{
			0xF9,
			0x00,
		}); err != nil {
			return err
		}

		return p.controller.WaitBusy()
	}

	// partial update initialization
	p.controller.Command(0x2C) // VCOM voltage
	p.controller.Data([]byte{0x26})

	p.controller.WaitBusy()

	p.controller.Command(0x32) // Load LUT
	p.controller.Data(lutPartialUpdate[:70])

	p.controller.Command(0x37)
	p.controller.Data([]byte{0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00})

	p.controller.Command(0x22)
	p.controller.Data([]byte{0xC0})
	p.controller.Command(0x20)

	p.controller.WaitBusy()

	p.controller.Command(0x3C) // Border waveform

	return p.controller.Data([]byte{0x01})
}

func (p *EPD213V2) Display(fb *koi.Framebuffer) error {
	if p.fullUpdate {
		return p.DisplayBuffer(fb.Data())
	}
	return p.DisplayPartial(fb.Data())
}

func (p *EPD213V2) DisplayBuffer(buf []byte) error {
	if err := p.controller.Command(0x24); err != nil {
		return err
	}

	if err := p.controller.Data(buf); err != nil {
		return err
	}

	// Display update control
	if err := p.controller.Command(0x22); err != nil {
		return err
	}

	if err := p.controller.Data([]byte{0xC7}); err != nil {
		return err
	}

	// Master activation
	if err := p.controller.Command(0x20); err != nil {
		return err
	}

	return p.controller.WaitBusy()
}

func (p *EPD213V2) DisplayPartial(buf []byte) error {
	// New data RAM
	if err := p.controller.Command(0x24); err != nil {
		return err
	}

	if err := p.controller.Data(buf); err != nil {
		return err
	}

	// Old data RAM: inverted so every pixel registers as changed against
	// the partial LUT, forcing a full-looking refresh without flashing.
	inverted := make([]byte, len(buf))
	for i, b := range buf {
		inverted[i] = ^b
	}

	if err := p.controller.Command(0x26); err != nil {
		return err
	}

	if err := p.controller.Data(inverted); err != nil {
		return err
	}

	// Display update control (partial update, no flash)
	if err := p.controller.Command(0x22); err != nil {
		return err
	}

	if err := p.controller.Data([]byte{0x0C}); err != nil {
		return err
	}

	// Master activation
	if err := p.controller.Command(0x20); err != nil {
		return err
	}

	return p.controller.WaitBusy()
}

func (p *EPD213V2) Clear(color byte) error {
	buf := make([]byte, ((Width+7)/8)*Height)

	for i := range buf {
		buf[i] = color
	}

	return p.DisplayBuffer(buf)
}

func (p *EPD213V2) Sleep() error {
	if err := p.controller.Command(0x10); err != nil {
		return err
	}

	if err := p.controller.Data([]byte{0x03}); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	return nil
}
