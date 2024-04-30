package animal

import (
	//"fmt"

	//termC "github.com/fatih/color"
	"github.com/hajimehoshi/ebiten/v2"

	ve "github.com/quartercastle/vector"
)

type vec = ve.Vector //  Vektoren

// Vor.: -
// Erg.: ein neues Tier
// New (w world.World, x, y float64) *data // *data erfüllt das Interface Animal

type Animal interface {
	// Die neue Position e.pos aus e.vel und e.acc bestimmen.
	Update(others *[]Animal)

	// Vor.: -
	// Eff.: -
	// Erg.: True ist geliefert, wenn das Tier noch Lebensenergie besitzt.
	IsAlive() bool

	SetEnergyLoss(e float64)
	SetEnergy(e float64)

	// Vor.: -
	// Eff.: Das Bild des Tieres wird durch die angegebene Datei ersetzt.
	// 		 Falls es ein Problem mit der angegebenen Datei gibt, wird nichts geaendert.
	// Erg.: -
	SetImageFromFile(file string)

	// Vor.:
	// Eff.: Setzt die RGB Farbe für das Tier und zeichnet es neu.
	// Erg.:
	SetColorRGB(r, g, b uint8)

	// Vor.:
	// Eff.: Das Attribute der Beutetiere ist gesetzt
	// Erg.:
	SetPreys(preys *[]Animal)

	SetMaxVel(v float64)

	SetViewMag(mag float64)

	// Vor.: ?
	// Eff.: ?
	// Erg.: Splice mit Objekten (seen) und deren Abstandsvektoren (direction),
	// die im Sichtfeld des Objekts liegen
	SeeOthers(others *[]Animal) (*[]Animal, *[]vec)

	// Vor.: -
	// Eff.: Das Tier ist gezeichnet.
	// Erg.: -
	Draw(screen *ebiten.Image)

	// Vor.: -
	// Eff.: -
	// Erg.: Liefert die aktuelle Position des Tieres.
	GetPosition() vec

	// Vor.: -
	// Eff.: -
	// Erg.: Prueft, ob das Tier die gleichen Eigenschaften hat wie das
	// 		 uebergebenen Tier.
	IsSame(b *data) bool

	// Vor.:
	// Eff.: Das Objekt beschleunigt von anderen Objekten, die im Sichtfeld liegen weg
	// Erg.:
	//avoidCollisionWithSeenObjects(others []Animal)
}
