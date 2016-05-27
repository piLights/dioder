// Some parts of this file are taken from Lakrizz' example on GitHub
// https://github.com/lakrizz/go-pidioder/blob/master/main.go

package dioder

import (
	"bufio"
	"errors"
	"image/color"
	"os"
	"strconv"
	"time"
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
	PiBlaster          string
	FadeDuration       time.Duration
}

//New creates a new instance
func New(pinConfiguration Pins, piBlasterFile string, fadeDuration time.Duration) Dioder {
	if piBlasterFile == "" {
		piBlasterFile = "/dev/pi-blaster"
	}

	d := Dioder{}

	d.SetPins(pinConfiguration)
	d.PiBlaster = piBlasterFile
	d.FadeDuration = fadeDuration

	return d
}

// GetCurrentColor returns the current color
func (d *Dioder) GetCurrentColor() color.RGBA {
	return d.ColorConfiguration
}

// SetAll sets the given values for the channels
func (d *Dioder) SetAll(colorSet color.RGBA) {
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
	d.PinConfiguration = pinConfiguration
}

//TurnOff turns off the dioder-strips and saves the current configuration
func (d *Dioder) TurnOff() {
	//Temporary save the configuration
	configuration := d.ColorConfiguration

	fadeChan := make(chan color.RGBA)
	fade(fadeChan, d.ColorConfiguration, color.RGBA{}, d.FadeDuration)

	for fadeConfiguration := range fadeChan {
		d.SetAll(fadeConfiguration)
	}

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

	fadeChan := make(chan color.RGBA)
	fade(fadeChan, colorSet, d.ColorConfiguration, d.FadeDuration)

	for fadeConfiguration := range fadeChan {
		d.SetAll(fadeConfiguration)
	}
}

func floatToString(floatValue float64) string {
	return strconv.FormatFloat(floatValue, 'f', 6, 64)
}

//SetColor Sets a color on the given channel
func (d *Dioder) SetColor(channel string, value float64) error {
	piBlasterCommand := channel + "=" + floatToString(value) + "\n"

	file, error := os.OpenFile(d.PiBlaster, os.O_RDWR, os.ModeNamedPipe)

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

func fade(fadeChan chan color.RGBA, oldColorSet, newColorSet color.RGBA, duration time.Duration) {

	//Calculate the diference between all those values
	//Red
	redDifference := getDifference(oldColorSet.R, newColorSet.R)
	greenDifference := getDifference(oldColorSet.G, newColorSet.G)
	blueDifference := getDifference(oldColorSet.B, newColorSet.B)

	average := (redDifference + greenDifference + blueDifference) / 3

	sleepTimePerStep := (duration.Nanoseconds() / int64(time.Millisecond)) / int64(average)

	for {
		var partialFadedColorConfiguration color.RGBA

		//Transition to the partial neew value

		fadeChan <- partialFadedColorConfiguration

		time.Sleep(time.Duration(sleepTimePerStep))
	}
	close(fadeChan)
}

func getDifference(a, b uint8) uint8 {
	if a > b {
		return a - b
	} else {
		return b - a
	}
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
