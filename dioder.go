// Some parts of this file are taken from Lakrizz' example on GitHub
// https://github.com/lakrizz/go-pidioder/blob/master/main.go

package dioder

import (
	"bufio"
	"errors"
	"image/color"
	"os"
	"strconv"
)

//Pins the numbers of the RGB-pins
type Pins struct {
	Red   string
	Green string
	Blue  string
}

//Dioder the main structure
type Dioder struct {
	PinConfiguration   Pins
	ColorConfiguration color.RGBA
}

//New creates a new instance
func New(pinConfiguration Pins) Dioder {
	d := Dioder{}

	d.SetPins(pinConfiguration)

	return d
}

// GetCurrentColor returns the current color
func (d *Dioder) GetCurrentColor() color.RGBA {
	return d.ColorConfiguration
}

// SetAll sets the given values for the channels
func (d *Dioder) SetAll(colorSet color.RGBA) {
	d.SetRed(colorSet.R)
	d.SetGreen(colorSet.G)
	d.SetBlue(colorSet.B)

	d.ColorConfiguration = colorSet
}

// SetBlue sets the given value on the Blue channel
func (d *Dioder) SetBlue(value uint8) error {
	// Do nothing if the new value is the same as the old
	if value == d.ColorConfiguration.B {
		return nil
	}

	d.ColorConfiguration.B = value
	return setChannelInteger(value, d.PinConfiguration.Blue)
}

// SetGreen sets the given value on the Green channel
func (d *Dioder) SetGreen(value uint8) error {
	// Do nothing if the new value is the same as the old
	if value == d.ColorConfiguration.G {
		return nil
	}

	d.ColorConfiguration.G = value
	return setChannelInteger(value, d.PinConfiguration.Green)
}

//SetPins configures the pin-layout
func (d *Dioder) SetPins(pinConfiguration Pins) {
	d.PinConfiguration = pinConfiguration
}

// SetRed sets the given value on the Red channel
func (d *Dioder) SetRed(value uint8) error {
	// Do nothing if the new value is the same as the old
	if value == d.ColorConfiguration.R {
		return nil
	}

	d.ColorConfiguration.R = value
	return setChannelInteger(value, d.PinConfiguration.Red)
}

//TurnOff turns off the dioder-strips and saves the current configuration
func (d *Dioder) TurnOff() {
	//Temporary save the configuration
	configuration := d.ColorConfiguration
	d.SetAll(color.RGBA{})
	d.ColorConfiguration = configuration
}

//TurnOn turns the dioder-strips on and restores the previous configuration
func (d *Dioder) TurnOn() {
	if d.ColorConfiguration.A == 0 && d.ColorConfiguration.B == 0 && d.ColorConfiguration.G == 0 && d.ColorConfiguration.R == 0 {
		d.ColorConfiguration = color.RGBA{255, 255, 255, 100}
	}

	d.SetAll(d.ColorConfiguration)
}

func floatToString(floatValue float64) string {
	return strconv.FormatFloat(floatValue, 'f', 6, 64)
}

func setColor(channel string, value float64) error {
	piBlasterCommand := channel + "=" + floatToString(value) + "\n"

	file, error := os.OpenFile("/dev/pi-blaster", os.O_RDWR, os.ModeNamedPipe)

	if error != nil {
		panic(error)
	}

	defer file.Close()

	stream := bufio.NewWriter(file)

	_, error = stream.WriteString(piBlasterCommand)

	if error != nil {
		panic(error)
	}

	stream.Flush()

	return nil
}

func setChannelInteger(value uint8, channel string) error {
	if value > 255 {
		return errors.New("Value can not be over 255")
	}

	if value < 0 {
		return errors.New("Value can not be under 0")
	}

	floatval := float64(value) / 255.0

	setColor(channel, float64(floatval))

	return nil
}
