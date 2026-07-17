package hw

import (
	"fmt"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
)

type SPI struct {
	dev spi.Conn
}

func NewSPI(dev spi.Conn) *SPI {
	return &SPI{
		dev: dev,
	}
}

// NewSPIFromPort opens and configures a periph.io SPI port by name and
// returns a ready to use SPI.
func NewSPIFromPort(name string, speed physic.Frequency, mode spi.Mode, bits int) (*SPI, error) {
	port, err := spireg.Open(name)
	if err != nil {
		return nil, fmt.Errorf("open spi port %s: %w", name, err)
	}

	conn, err := port.Connect(speed, mode, bits)
	if err != nil {
		return nil, fmt.Errorf("connect spi port %s: %w", name, err)
	}

	return NewSPI(conn), nil
}

func (s *SPI) Write(buf []byte) error {
	err := s.dev.Tx(buf, nil)

	return err
}
