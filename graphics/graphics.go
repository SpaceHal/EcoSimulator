package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Vor.: -
// Erg.: neue Liniendiagramme
// New (x, y int) *data // *data erfüllt das Interface Graphics

type Graphics interface {
	// Vor.: -
	// Eff.: Die Liniendiagramme sind gezeichnet.
	// Erg.: -
	Draw(dst *ebiten.Image)

	// Vor.: -
	// Eff.: Neue Werte für die Liniendiagramme sind ergänzt
	// Erg.: -
	Update(nG, nR, nC, nF int)
}
