package dioder

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	d := New(Pins{}, "")

	if reflect.TypeOf(d).String() != "dioder.Dioder" {
		t.Errorf("Got wrong dioder.Dioder object: %s", reflect.TypeOf(d))
	}
}

func TestDioderGetCurrentColor(t *testing.T) {
	dioder := New(Pins{"18", "17", "4"}, "/dev/pi-blaster")

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

	dioder := New(pinConfiguration, "/dev/pi-blaster")

	if dioder.PinConfiguration != pinConfiguration {
		t.Errorf("Pins are not correctly configured. Gave %s, got %s", pinConfiguration, dioder.PinConfiguration)
	}

	if dioder.PiBlaster != "/dev/pi-blaster" {
		t.Errorf("Pi-Blaster file wrong configured. Gave /dev/pi-blaster, got %s", dioder.PiBlaster)
	}
}
