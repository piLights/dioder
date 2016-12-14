// Some parts of this file are taken from Lakrizz' example on GitHub
// https://github.com/lakrizz/go-pidioder/blob/master/main.go

package dioder

import (
	"bufio"
	"image/color"
	"os"
	"strconv"
	"sync"
	"fmt"
)

const piBlasterLocation = "/dev/pi-blaster"

//Pins the numbers of the RGB-pins
type Pins struct {
	Red   int
	Green int
	Blue  int
}

//Dioder the main structure
type Dioder struct {
	sync.Mutex
	PinConfiguration   Pins
	ColorConfiguration color.RGBA
	PiBlaster          string
}

//New creates a new instance
func New(pinConfiguration Pins, piBlasterFile string) Dioder {
	if piBlasterFile == "" {
		piBlasterFile = piBlasterLocation
	}

	d := Dioder{}

	d.SetPins(pinConfiguration)
	d.PiBlaster = piBlasterFile

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
	//Ugliy hack, to turn the lights back on
	colorSet := d.ColorConfiguration
	d.ColorConfiguration = color.RGBA{}

	d.SetAll(colorSet)
}

func floatToString(floatValue float64) string {
	return strconv.FormatFloat(floatValue, 'f', 6, 64)
}

//SetColor Sets a color on the given channel
func (d *Dioder) SetColor(channel int, value float64) error {
	piBlasterCommand := fmt.Sprintf("%d=%s\n", channel, floatToString(value))

	file, err := os.OpenFile(d.PiBlaster, os.O_RDWR, os.ModeNamedPipe)

	if err != nil {
		return err
	}

	defer file.Close()

	stream := bufio.NewWriter(file)

	_, err = stream.WriteString(piBlasterCommand)

	if err != nil {
		return err
	}

	stream.Flush()

	return nil
}

//SetChannelInteger check if the value is in the correct range and convert it to float64
func (d *Dioder) SetChannelInteger(value uint8, channel int) error {
	floatval := float64(value) / 255.0

	err := d.SetColor(channel, float64(floatval))
	return err
}

//Release releases all used pins, so that they can be used in other applications
func (d *Dioder) Release() error {
	d.Lock()
	defer d.Unlock()

	piBlasterCommand := fmt.Sprintf("release %d\n", d.PinConfiguration.Red)
	piBlasterCommand += fmt.Sprintf("release %d\n", d.PinConfiguration.Green)
	piBlasterCommand += fmt.Sprintf("release %d\n", d.PinConfiguration.Blue)

	file, err := os.OpenFile(d.PiBlaster, os.O_RDWR, os.ModeNamedPipe)

	if err != nil {
		return err
	}

	defer file.Close()

	stream := bufio.NewWriter(file)

	_, err = stream.WriteString(piBlasterCommand)

	if err != nil {
		return err
	}

	stream.Flush()

	return nil
}

/*func (d *Dioder) fade(currentValue uint8, targetValue uint8, fadeTime time.Duration, channel string) {
	//Fade:
	//Einteilung der Werteminderung pro Zeiteinheit anhand der Werte zwischen currentValue und targetValue
	var neededSteps uint8

	if currentValue < targetValue {
		neededSteps = targetValue - currentValue
	} else {
		neededSteps = currentValue - targetValue
	}

	for neededSteps {
		//set value
		//sleep duration / neededSteps
	}

}*/

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
