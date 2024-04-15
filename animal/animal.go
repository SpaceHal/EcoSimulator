package animal

import (
	//"fmt"

	//termC "github.com/fatih/color"
	"github.com/hajimehoshi/ebiten/v2"

	ve "github.com/quartercastle/vector"
)

type vec = ve.Vector //  Vektoren

// Vor.: -
// Erg.: ein neuer Tier 
// New (w world.World, x, y float64) *data // *data erfüllt das Interface World 

type Animal interface {
	// Die neue Position e.pos aus e.vel und e.acc bestimmen.
	Update(others []Animal)
	
	// Vor.: ?
	// Eff.: ?
	// Erg.: Splice mit Objekten (seen) und deren Abstandsvektoren (direction),
	// die im Sichtfeld des Objekts liegen
	SeeOthers(others []Animal) (seen []Animal, direction []vec)
	
	Draw(screen *ebiten.Image)
	
	GetPosition() vec
	
	IsSame(b *data) bool
}
