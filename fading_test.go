package dioder

import (
	"image/color"
	"testing"
	"time"
	"fmt"
)

func TestDioderfade(t *testing.T) {
	d := New(Pins{"1", "2", "3"}, tmpPiBlasterLocation)

	d.ColorConfiguration = color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0,
	}

	duration, _ := time.ParseDuration("2s")
	d.fade(color.RGBA{255, 255, 255, 100}, duration)

	fmt.Println(d.ColorConfiguration)
}
