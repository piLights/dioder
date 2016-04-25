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

var (
	redPin   = "18"
	bluePin  = "17"
	greenPin = "4"
)

var currentColorConfiguration color.RGBA

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

// SetAll sets the given values for the channels
func SetAll(colorSet color.RGBA) {
	SetRed(colorSet.R)
	SetGreen(colorSet.G)
	SetBlue(colorSet.B)
}

// SetRed sets the given value on the Red channel
func SetRed(value uint8) error {
	// Do nothing if the new value is the same as the old
	if value == currentColorConfiguration.R {
		return nil
	}

	currentColorConfiguration.R = value
	return setChannelInteger(value, redPin)
}

// SetGreen sets the given value on the Green channel
func SetGreen(value uint8) error {
	// Do nothing if the new value is the same as the old
	if value == currentColorConfiguration.G {
		return nil
	}

	currentColorConfiguration.G = value
	return setChannelInteger(value, greenPin)
}

// SetBlue sets the given value on the Blue channel
func SetBlue(value uint8) error {
	// Do nothing if the new value is the same as the old
	if value == currentColorConfiguration.B {
		return nil
	}

	currentColorConfiguration.B = value
	return setChannelInteger(value, bluePin)
}

func SetPins(red, green, blue int) {
	redPin = strconv.Itoa(red)
	greenPin = strconv.Itoa(green)
	bluePin = strconv.Itoa(blue)
}

// GetCurrentColor returns the current color
func GetCurrentColor() color.RGBA {
	return currentColorConfiguration
}
