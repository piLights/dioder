package dioder

import "testing"

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
