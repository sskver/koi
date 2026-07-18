package main

import (
	"log/slog"
	"math/rand/v2"
	"os"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/host/v3"

	"github.com/sskver/koi"
	"github.com/sskver/koi/hw"
	"github.com/sskver/koi/panels"
)

func fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}

func main() {
	if _, err := host.Init(); err != nil {
		fatal("host init failed", err)
	}

	pins := koi.Pins{
		DC:   25,
		RST:  17,
		BUSY: 24,
		PWR:  18,
	}

	gpioClient, err := hw.NewGPIOFromPins(pins)
	if err != nil {
		fatal("gpio init failed", err)
	}

	spiClient, err := hw.NewSPIFromPort("SPI0.0", physic.MegaHertz*4, spi.Mode0, 8)
	if err != nil {
		fatal("spi init failed", err)
	}

	controller := koi.NewController(
		koi.ControllerConfig{
			Pins: pins,
			GPIO: gpioClient,
			SPI:  spiClient,
		},
	)

	if err := controller.Init(); err != nil {
		fatal("controller init failed", err)
	}

	fullPanel := panels.NewEPD213V2(controller, true)

	if err := fullPanel.Init(); err != nil {
		fatal("panel init failed", err)
	}

	if err := fullPanel.Clear(0xFF); err != nil {
		fatal("panel clear failed", err)
	}

	panel := panels.NewEPD213V2(controller, false)

	if err := panel.Init(); err != nil {
		fatal("panel init failed", err)
	}

	fb := koi.NewFramebuffer(
		panels.Width,
		panels.Height,
	)

	x, y := 0, 0
	size := 20

	xNeg, yNeg := false, false

	for {
		fb.Clear(0xFF)

		stepX := size + 1 + rand.IntN(15)
		stepY := size + 1 + rand.IntN(15)

		if xNeg {
			x -= stepX
		} else {
			x += stepX
		}
		if yNeg {
			y -= stepY
		} else {
			y += stepY
		}

		if x >= panels.Width-1-size {
			x = panels.Width - 1 - size
			xNeg = true
		}
		if x <= 0 {
			x = 0
			xNeg = false
		}
		if y >= panels.Height-1-size {
			y = panels.Height - 1 - size
			yNeg = true
		}
		if y <= 0 {
			y = 0
			yNeg = false
		}

		for i := range size {
			fb.ClearPixel(x+i, y)
			fb.ClearPixel(x+i, y+size-1)
			fb.ClearPixel(x, y+i)
			fb.ClearPixel(x+size-1, y+i)
		}

		if err := panel.Display(fb); err != nil {
			fatal("panel display failed", err)
		}
	}

}
