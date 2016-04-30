package dioder

import (
	"image/color"
	"testing"
)

func TestNew(t *testing.T) {
	t.SkipNow()
}

func TestDioderGetCurrentColor(t *testing.T) {
	dioder := New(Pins{"18", "17", "4"})

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

	dioder := New(pinConfiguration)

	if dioder.PinConfiguration != pinConfiguration {
		t.Errorf("Pins are not correctly configured. Gave %s, got %s", pinConfiguration, dioder.PinConfiguration)
	}
}

func TestDioderSetAll(t *testing.T) {
	d := New(Pins{})

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
	d := New(Pins{"1", "2", "3"})

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
