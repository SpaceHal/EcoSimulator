package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Vor.: -
// Erg.: neue Graphiken
// New (x, y int) *data // *data erfüllt das Interface Graphics

type Graphics interface {
	// Vor.: -
	// Eff.: Die Grafiken sind gezeichnet.
	// Erg.: -
	Draw(dst *ebiten.Image)
	
	Update(nG,nR,nC,nF int)
}
