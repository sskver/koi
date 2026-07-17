package koi

import (
	"errors"
	"log/slog"
	"time"
)

var ErrNotInitialized = errors.New("controller not initialized")

type ControllerConfig struct {
	Pins Pins
	GPIO GPIO
	SPI  SPI
}

type Pins struct {
	DC   int
	RST  int
	BUSY int
	PWR  int
}

type Controller struct {
	cfg ControllerConfig

	initialized bool
}

type GPIO interface {
	Set(pin int, high bool) error
	Read(pin int) (bool, error)
}

type SPI interface {
	Write(buf []byte) error
}

func NewController(cfg ControllerConfig) *Controller {
	return &Controller{
		cfg: cfg,
	}
}

func (c *Controller) Init() error {
	if c.initialized {
		return nil
	}

	if err := c.cfg.GPIO.Set(c.cfg.Pins.PWR, true); err != nil {
		return err
	}

	time.Sleep(10 * time.Millisecond)

	if err := c.cfg.GPIO.Set(c.cfg.Pins.DC, true); err != nil {
		return err
	}

	if err := c.cfg.GPIO.Set(c.cfg.Pins.RST, true); err != nil {
		return err
	}

	c.initialized = true

	return nil
}

func (c *Controller) Reset() error {
	if !c.initialized {
		return ErrNotInitialized
	}

	if err := c.cfg.GPIO.Set(c.cfg.Pins.RST, true); err != nil {
		return err
	}

	time.Sleep(200 * time.Millisecond)

	if err := c.cfg.GPIO.Set(c.cfg.Pins.RST, false); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)

	if err := c.cfg.GPIO.Set(c.cfg.Pins.RST, true); err != nil {
		return err
	}

	time.Sleep(200 * time.Millisecond)

	return nil
}

func (c *Controller) Command(cmd byte) error {
	if !c.initialized {
		return ErrNotInitialized
	}

	if err := c.cfg.GPIO.Set(c.cfg.Pins.DC, false); err != nil {
		return err
	}

	return c.cfg.SPI.Write([]byte{cmd})
}

func (c *Controller) Data(data []byte) error {
	if !c.initialized {
		return ErrNotInitialized
	}

	if err := c.cfg.GPIO.Set(c.cfg.Pins.DC, true); err != nil {
		return err
	}

	return c.cfg.SPI.Write(data)
}

func (c *Controller) WaitBusy() error {
	if !c.initialized {
		return ErrNotInitialized
	}

	startTime := time.Now()

	for {
		busy, err := c.cfg.GPIO.Read(c.cfg.Pins.BUSY)

		if err != nil {
			slog.Error("error reading busy pin", "error", err)
			return err
		}

		if !busy {
			return nil
		}

		if time.Since(startTime) > 10*time.Second {
			slog.Error("timeout waiting for busy pin to go low")
			return errors.New("timeout waiting for busy pin to go low")
		}

		time.Sleep(10 * time.Millisecond)
	}
}
