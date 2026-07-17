package hw

import (
	"fmt"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"

	"github.com/sskver/koi"
)

type GPIO struct {
	pins map[int]gpio.PinIO
}

func NewGPIO(pins map[int]gpio.PinIO) *GPIO {
	return &GPIO{
		pins: pins,
	}
}

// NewGPIOFromPins resolves and configures the pins described by a koi.Pins
// wiring config (RST/PWR/DC as outputs, BUSY as input) and returns a ready
// to use GPIO.
func NewGPIOFromPins(pins koi.Pins) (*GPIO, error) {
	gpioPins := make(map[int]gpio.PinIO)

	outputs := []int{pins.RST, pins.PWR, pins.DC}
	for _, num := range outputs {
		p, err := resolvePin(num)
		if err != nil {
			return nil, err
		}

		out, ok := p.(gpio.PinOut)
		if !ok {
			return nil, fmt.Errorf("gpio pin %d is not output capable", num)
		}

		if err := out.Out(gpio.Low); err != nil {
			return nil, fmt.Errorf("gpio pin %d: %w", num, err)
		}

		gpioPins[num] = p
	}

	in, err := resolvePin(pins.BUSY)
	if err != nil {
		return nil, err
	}

	inPin, ok := in.(gpio.PinIn)
	if !ok {
		return nil, fmt.Errorf("gpio pin %d is not input capable", pins.BUSY)
	}

	if err := inPin.In(gpio.PullUp, gpio.NoEdge); err != nil {
		return nil, fmt.Errorf("gpio pin %d: %w", pins.BUSY, err)
	}

	gpioPins[pins.BUSY] = in

	return NewGPIO(gpioPins), nil
}

func resolvePin(num int) (gpio.PinIO, error) {
	p := gpioreg.ByName(fmt.Sprintf("GPIO%d", num))
	if p == nil {
		return nil, fmt.Errorf("gpio pin %d not found", num)
	}

	return p, nil
}

func (g *GPIO) Set(pin int, high bool) error {
	p, ok := g.pins[pin]
	if !ok || p == nil {
		return fmt.Errorf("gpio pin %d not found", pin)
	}

	out, ok := p.(gpio.PinOut)
	if !ok {
		return fmt.Errorf("gpio pin %d is not output capable", pin)
	}

	if high {
		return out.Out(gpio.High)
	}

	return out.Out(gpio.Low)
}

func (g *GPIO) Read(pin int) (bool, error) {
	p, ok := g.pins[pin]
	if !ok || p == nil {
		return false, fmt.Errorf("gpio pin %d not found", pin)
	}

	in, ok := p.(gpio.PinIn)
	if !ok {
		return false, fmt.Errorf("gpio pin %d is not input capable", pin)
	}

	v := in.Read()

	return v == gpio.High, nil
}
