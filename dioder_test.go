package dioder

import (
	"image/color"
	"reflect"
	"testing"
)

var (
	piBlasterFile = "/tmp/pi-blaster"
	pinConfiguration = Pins{
		Red: 1,
		Green: 2,
		Blue: 3,
	}
)

func TestNew(t *testing.T) {
	d := New(Pins{}, "")

	if reflect.TypeOf(d).String() != "dioder.Dioder" {
		t.Errorf("Got wrong dioder.Dioder object: %s", reflect.TypeOf(d))
	}
}

func TestDioderGetCurrentColor(t *testing.T) {
	dioder := New(pinConfiguration, piBlasterFile)

	colorSet := dioder.GetCurrentColor()

	if colorSet.A != 0 {
		t.Error("Opacity is not 0")
	}

	if colorSet.R != 0 {
		t.Error("Red is not 0")
	}

	if colorSet.G != 0 {
		t.Error("Green is not 0")
	}

	if colorSet.B != 0 {
		t.Error("Blue is not 0")
	}
}

func TestDioderPinConfiguration(t *testing.T) {
	dioder := New(pinConfiguration, piBlasterFile)

	if dioder.PinConfiguration != pinConfiguration {
		t.Errorf("Pins are not correctly configured. Gave %s, got %s", pinConfiguration, dioder.PinConfiguration)
	}

	if dioder.PiBlaster != piBlasterFile {
		t.Errorf("Pi-Blaster file wrong configured. Gave %s, got %s", piBlasterFile, dioder.PiBlaster)
	}
}

func TestDioderSetAll(t *testing.T) {
	d := New(Pins{}, piBlasterFile)

	//All values at zero
	d.SetAll(color.RGBA{})

	if d.ColorConfiguration.A != 0 {
		t.Error("Opacity is not correct")
	}

	if d.ColorConfiguration.R != 0 {
		t.Error("Red is not correct")
	}

	if d.ColorConfiguration.G != 0 {
		t.Error("Green is not correct")
	}

	if d.ColorConfiguration.B != 0 {
		t.Error("Blue is not correct")
	}

	d.SetAll(color.RGBA{255, 255, 255, 0})
	if d.ColorConfiguration.R != 255 {
		t.Errorf("Red is not correct after applying opacity of zero. Got: %d\n", d.ColorConfiguration.R)
	}

	if d.ColorConfiguration.G != 255 {
		t.Errorf("Green is not correct after applying opacity of zero. Got: %d\n", d.ColorConfiguration.G)
	}

	if d.ColorConfiguration.B != 255 {
		t.Errorf("Blue is not correct after applying opacity of zero. Got: %d\n", d.ColorConfiguration.B)
	}
}

func TestDioderSetPins(t *testing.T) {
	d := New(pinConfiguration, piBlasterFile)

	if d.PinConfiguration.Blue != 3 {
		t.Error("Blue pin is not correct")
	}

	if d.PinConfiguration.Green != 2 {
		t.Error("Green pin is not correct")
	}

	if d.PinConfiguration.Red != 1 {
		t.Error("Red pin is not correct")
	}
}

func TestDioderTurnOff(t *testing.T) {
	d := New(pinConfiguration, piBlasterFile)

	configuration := d.ColorConfiguration

	d.TurnOff()

	if d.ColorConfiguration != configuration {
		t.Errorf("Didn't saved the current settings, had %d - got %d", configuration, d.ColorConfiguration)
	}
}

func TestDioderTurnOn(t *testing.T) {
	d := New(pinConfiguration, piBlasterFile)

	d.TurnOff()
	d.TurnOn()

	if d.ColorConfiguration.A != 100 {
		t.Errorf("Value for opacity is wrong. Expected 0, got %d", d.ColorConfiguration.A)
	}
	if d.ColorConfiguration.R != 255 {
		t.Errorf("Value for red is wrong. Expected 0, got %d", d.ColorConfiguration.R)
	}
	if d.ColorConfiguration.G != 255 {
		t.Errorf("Value for green is wrong. Expected 0, got %d", d.ColorConfiguration.G)
	}
	if d.ColorConfiguration.B != 255 {
		t.Errorf("Value for blue is wrong. Expected 0, got %d", d.ColorConfiguration.B)
	}
}

func TestDioder_SetChannelInteger(t *testing.T) {
	d := New(pinConfiguration, piBlasterFile)

	err := d.SetChannelInteger(255, 1)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestDioder_SetColor(t *testing.T) {
	d := New(pinConfiguration, "nonExistentFile")

	err := d.SetColor(0, 0)
	if err == nil {
		t.Error("Called with an non existent filename. SetColor() should return an error")
	}

	d.PiBlaster = piBlasterFile

	err = d.SetColor(0, 0)
	if err != nil {
		t.Errorf("Returned an error: %s", err)
	}
}
