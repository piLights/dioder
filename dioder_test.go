package dioder

import (
	"image/color"
	"reflect"
	"testing"
)

var piBlasterFile = "/dev/pi-blaster"

func TestNew(t *testing.T) {
	d := New(Pins{}, "")

	if reflect.TypeOf(d).String() != "dioder.Dioder" {
		t.Errorf("Got wrong dioder.Dioder object: %s", reflect.TypeOf(d))
	}
}

func TestDioderGetCurrentColor(t *testing.T) {
	dioder := New(Pins{"18", "17", "4"}, piBlasterFile)

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
	pinConfiguration := Pins{"18", "17", "4"}

	dioder := New(pinConfiguration, piBlasterFile)

	if dioder.PinConfiguration != pinConfiguration {
		t.Errorf("Pins are not correctly configured. Gave %s, got %s", pinConfiguration, dioder.PinConfiguration)
	}

	if dioder.PiBlaster != "/dev/pi-blaster" {
		t.Errorf("Pi-Blaster file wrong configured. Gave /dev/pi-blaster, got %s", dioder.PiBlaster)
	}
}

func TestDioderSetAll(t *testing.T) {
	d := New(Pins{}, "/dev/pi-blaster")

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
}

func TestDioderSetPins(t *testing.T) {
	d := New(Pins{"1", "2", "3"}, piBlasterFile)

	if d.PinConfiguration.Blue != "3" {
		t.Error("Blue pin is not correct")
	}

	if d.PinConfiguration.Green != "2" {
		t.Error("Green pin is not correct")
	}

	if d.PinConfiguration.Red != "1" {
		t.Error("Red pin is not correct")
	}
}

func TestDidoerTurnOff(t *testing.T) {
	d := New(Pins{"1", "2", "3"}, piBlasterFile)

	configuration := d.ColorConfiguration

	d.TurnOff()

	if d.ColorConfiguration != configuration {
		t.Errorf("Didn't saved the current settings, had %d - got %d", configuration, d.ColorConfiguration)
	}
}

func TestDioderTurnOn(t *testing.T) {
	d := New(Pins{"1", "2", "3"}, piBlasterFile)

	d.TurnOff()
	d.TurnOn()

	if d.ColorConfiguration.A != 255 {
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
