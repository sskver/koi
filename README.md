# koi

waveshare e-paper display library in go, built on [periph.io](https://periph.io). since it's built on periph.io's `conn`/`host` interfaces instead of a specific board's SDK, koi isn't tied to one piece of hardware - it runs anywhere periph.io has a host driver (raspberry pi and other linux SBCs with spi/gpio).

## features

- native go driver, no python or C bindings
- thin `controller` that handles reset timing, command/data framing over spi, and busy-pin polling
- for now, only a simple 1-bit packed `framebuffer` for building the image to send to the panel
- `panel` interface so new waveshare models can be added without touching the rest of the library
- depends only on `periph.io/x/conn/v3` and `periph.io/x/host/v3`

## install

```bash
go get github.com/sskver/koi
```

requires go 1.26+.

## how it fits together

- **`Controller`** - wraps gpio (dc/rst/busy/pwr pins) and spi into the low-level commands a waveshare panel expects: `Command()`, `Data()`, `Reset()`, `WaitBusy()`
- **`Framebuffer`** - a packed 1-bit-per-pixel buffer (`NewFramebuffer(width, height)`, `SetPixel`, `ClearPixel`, `Clear`)
- **`Panel`** - the interface a specific waveshare model implements: `Init()`, `Display(*Framebuffer) error`, `Sleep()`
- **`hw/`** - periph.io-backed implementations of the `GPIO`/`SPI` interfaces the controller needs
- **`panels/`** - per-model panel implementations built on top of `Controller`
- **`cmd/`** - a working example wiring everything together on real hardware

## quick start

```go
package main

import (
	"log"

	"periph.io/x/host/v3"

	"github.com/sskver/koi"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	ctrl := koi.NewController(koi.ControllerConfig{
		Pins: koi.Pins{
			DC:   25,
			RST:  17,
			BUSY: 24,
			PWR:  18,
		},
		GPIO: myGPIO, // your koi.GPIO implementation, e.g. from hw/
		SPI:  mySPI,  // your koi.SPI implementation, e.g. from hw/
	})

	if err := ctrl.Init(); err != nil {
		log.Fatal(err)
	}
  
    // something something, you would need to create the panel, e.g. your specific model
	fb := koi.NewFramebuffer(250, 122) // use the panel specs or define your own
	fb.SetPixel(10, 10)

	// then call panel.Display(fb)
}
```

see `cmd/` for the full, runnable version of this on actual hardware.

## supported hardware

any waveshare e-paper panel with an implementation under `panels/`, on any board periph.io's `host` package supports (raspberry pi, etc).

## license

mit - see [LICENSE](LICENSE)
