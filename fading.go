package dioder

import (
	"fmt"
	"image/color"
	"reflect"
	"time"
)

func getField(v color.RGBA, field string) uint8 {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return uint8(f.Uint())
}

type Action struct {
	Minus      bool
	Difference uint8
}

type FadingHelper struct {
	R             Action
	G             Action
	B             Action
	A             Action
	MaxDifference uint8
}

func (f *FadingHelper) SetMaxDifference() uint8 {
	maxDifference := f.R.Difference

	if maxDifference < f.G.Difference {
		maxDifference = f.G.Difference
	}

	if maxDifference < f.B.Difference {
		maxDifference = f.B.Difference
	}

	if maxDifference < f.A.Difference {
		maxDifference = f.A.Difference
	}

	f.MaxDifference = maxDifference
	return f.MaxDifference
}

func (d *Dioder) fade(targetColor color.RGBA, fadeTime time.Duration) {
	// Holen der größten Differenz
	// Errechnung der benötigten Slots: Differenz / FadeTime
	//
	// Für jeden Slot:
	// 	Berechnung des Steps für Farbe:
	//		FarbDifferenz / Slots

	var determineAction = func(colorType string) Action {
		action := Action{}

		if getField(d.ColorConfiguration, colorType) > getField(targetColor, colorType) {
			action.Minus = true
			action.Difference = getField(d.ColorConfiguration, colorType) - getField(targetColor, colorType)
		} else {
			action.Minus = false
			action.Difference = getField(targetColor, colorType) - getField(d.ColorConfiguration, colorType)
		}

		return action
	}

	var fadingHelper = FadingHelper{
		R: determineAction("R"),
		G: determineAction("G"),
		B: determineAction("B"),
		A: determineAction("A"),
	}
	fadingHelper.SetMaxDifference()

	if fadingHelper.MaxDifference != 0 {
		numberOfSlots := float64(fadeTime.Nanoseconds() / int64(time.Millisecond)) / float64(fadingHelper.MaxDifference)
		fmt.Println(numberOfSlots)
		fmt.Println(fadingHelper)

		for i := 0; i < int(numberOfSlots); i++ {
			// Red
			if fadingHelper.R.Difference != 0 {
				// Calculate the fraction to minimize
				var fraction float64
				if fadingHelper.R.Minus {
					fraction = (float64(d.ColorConfiguration.R) - (float64(fadingHelper.R.Difference) / numberOfSlots))
					d.ColorConfiguration.R = d.ColorConfiguration.R - uint8(fraction)
				} else {
					fraction = (float64(d.ColorConfiguration.R) + (float64(fadingHelper.R.Difference) / numberOfSlots))
					d.ColorConfiguration.R = d.ColorConfiguration.R + uint8(fraction)
				}


				d.SetColor(d.PinConfiguration.Red, fraction)
				fmt.Println(fraction, d.ColorConfiguration)
			}
			// Green
			if fadingHelper.G.Difference != 0 {
				// Calculate the fraction to minimize
				fraction := (float64(d.ColorConfiguration.G) + (float64(fadingHelper.G.Difference) / numberOfSlots)) / 255.0
				d.SetColor(d.PinConfiguration.Green, fraction)
				d.ColorConfiguration.G += uint8(fraction)
			}
			// Blue
			if fadingHelper.B.Difference != 0 {
				// Calculate the fraction to minimize
				fraction := (float64(d.ColorConfiguration.B) + (float64(fadingHelper.B.Difference) / numberOfSlots)) / 255.0
				d.SetColor(d.PinConfiguration.Blue, fraction)
				d.ColorConfiguration.B += uint8(fraction)
			}

			// @ToDo: Fade ove the alpha channel
			sleepDuration := time.Duration{int64(float64(fadeTime.Nanoseconds() / int64(time.Millisecond)) / numberOfSlots)}
			time.Sleep(sleepDuration)
		}
	}
}
