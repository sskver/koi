package main

import (
	"log/slog"
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

	panel := panels.NewEPD213V2(controller)

	if err := panel.Init(); err != nil {
		fatal("panel init failed", err)
	}

	fb := koi.NewFramebuffer(
		panels.Width,
		panels.Height,
	)

	fb.Clear(0xFF)

	for y := 50; y < 100; y++ {
		for x := 20; x < 80; x++ {
			fb.ClearPixel(x, y)
		}
	}

	if err := panel.Display(fb); err != nil {
		fatal("panel display failed", err)
	}
}
