// Some parts of this file are taken from Lakrizz' example on GitHub
// https://github.com/lakrizz/go-pidioder/blob/master/main.go

package dioder

import (
	"bufio"
	"errors"
	"image/color"
	"os"
	"strconv"
	"sync"
)

const piBlasterLocation = "/dev/pi-blaster"

//Pins the numbers of the RGB-pins
type Pins struct {
	Red   string
	Green string
	Blue  string
}

//Dioder the main structure
type Dioder struct {
	sync.Mutex
	PinConfiguration   Pins
	ColorConfiguration color.RGBA
	PiBlaster          string
	File               *os.File
}

//New creates a new instance
func New(pinConfiguration Pins, piBlasterFile string) Dioder {
	if piBlasterFile == "" {
		piBlasterFile = piBlasterLocation
	}

	d := Dioder{}

	d.SetPins(pinConfiguration)
	d.PiBlaster = piBlasterFile

	var err error
	d.File, err = os.OpenFile(d.PiBlaster, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}

	return d
}

// GetCurrentColor returns the current color
func (d *Dioder) GetCurrentColor() color.RGBA {
	d.Lock()
	defer d.Unlock()
	return d.ColorConfiguration
}

// SetAll sets the given values for the channels
func (d *Dioder) SetAll(colorSet color.RGBA) {
	d.Lock()
	defer d.Unlock()

	d.ColorConfiguration = colorSet
	//Red
	colorSet.R = calculateOpacity(colorSet.R, colorSet.A)
	d.SetChannelInteger(colorSet.R, d.PinConfiguration.Red)
	//Green
	colorSet.G = calculateOpacity(colorSet.G, colorSet.A)
	d.SetChannelInteger(colorSet.G, d.PinConfiguration.Green)

	//Blue
	colorSet.B = calculateOpacity(colorSet.B, colorSet.A)
	d.SetChannelInteger(colorSet.B, d.PinConfiguration.Blue)
}

//SetPins configures the pin-layout
func (d *Dioder) SetPins(pinConfiguration Pins) {
	d.Lock()
	defer d.Unlock()

	d.PinConfiguration = pinConfiguration
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

	//@ToDo: Refactor
	//Ugly hack, to turn the lights back on
	colorSet := d.ColorConfiguration
	d.ColorConfiguration = color.RGBA{}

	d.SetAll(colorSet)
}

func floatToString(floatValue float64) string {
	return strconv.FormatFloat(floatValue, 'f', 6, 64)
}

//SetColor Sets a color on the given channel
func (d *Dioder) SetColor(channel string, value float64) error {
	piBlasterCommand := channel + "=" + floatToString(value) + "\n"

	stream := bufio.NewWriter(d.File)

	_, err := stream.WriteString(piBlasterCommand)
	if err != nil {
		panic(err)
	}

	stream.Flush()

	return nil
}

//SetChannelInteger check if the value is in the correct range and convert it to float64
func (d *Dioder) SetChannelInteger(value uint8, channel string) error {
	if value > 255 {
		return errors.New("Value can not be over 255")
	}

	if value < 0 {
		return errors.New("Value can not be under 0")
	}

	floatval := float64(value) / 255.0

	d.SetColor(channel, float64(floatval))

	return nil
}

//Release releases all used pins, so that they can be used in other applications
func (d *Dioder) Release() {
	d.Lock()
	defer d.Unlock()

	piBlasterCommand := "release " + d.PinConfiguration.Red + "\n"
	piBlasterCommand += "release " + d.PinConfiguration.Green + "\n"
	piBlasterCommand += "release " + d.PinConfiguration.Blue + "\n"

	stream := bufio.NewWriter(d.File)

	_, err := stream.WriteString(piBlasterCommand)
	if err != nil {
		panic(err)
	}

	stream.Flush()

	d.File.Close()

	return
}

//calculateOpacity calculates the value of colorValue after applying some opacity
func calculateOpacity(colorValue uint8, opacity uint8) uint8 {
	var calculatedValue float32

	if opacity != 100 {
		calculatedValue = float32(colorValue) / 100 * float32(opacity)
	} else {
		calculatedValue = float32(colorValue)
	}

	return uint8(calculatedValue)
}
