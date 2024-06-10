package entity

import (
	//"fmt"

	//termC "github.com/fatih/color"
	"ecosim/world"

	"github.com/hajimehoshi/ebiten/v2"

	ve "github.com/quartercastle/vector"
)

type vec = ve.Vector //  Vektoren

// Vor.: -
// Erg.: ein neues Tier
// New (w world.World, x, y float64) *data // *data erfüllt das Interface Entity

type Entity interface {
	// Die neue Position e.pos aus e.vel und e.acc bestimmen.
	//Update(others *[]Entity) (offSpring *data)

	// Vor.: -
	// Eff.: -
	// Erg.: True ist geliefert, wenn das Tier noch Lebensenergie besitzt.
	IsAlive() bool

	SetHealthLoss(e float64)
	SetHealth(e float64)
	GetHealth() float64
	SetHealthWhenEaten(e float64)
	GetHealthWhenEaten() float64

	GetAge() int
	IncAge()
	SetMatureAge(mAge int)
	SetLifeSpan(ls int)
	GetDateOfLastBirth() int
	SetDateOfLastBirth(d int)

	GetWorld() *world.World

	/*
		GetPreys() *[]Entity
		GetNumOfPreys() int
		GetPredators() *[]Entity
		SetPreys(preys *[]Entity)
		SetPredators(preds *[]Entity)
	*/

	GetMatureAge() int

	SetViewAngle(ang float64)

	SetMoveable(m bool)
	ApplyMove(others *[]Entity, preys *[]Entity)

	// Vor.: -
	// Eff.: Das Bild des Tieres wird durch die angegebene Datei ersetzt.
	// 		 Falls es ein Problem mit der angegebenen Datei gibt, wird nichts geaendert.
	// Erg.: -
	SetImageFromFile(file string, size, x, y int)

	// Vor.:
	// Eff.: Setzt die RGB Farbe für das Tier und zeichnet es neu.
	// Erg.:
	SetColorRGB(r, g, b uint8)

	SetMaxVel(v float64)

	SetViewMag(mag float64)

	// Vor.: ?
	// Eff.: ?
	// Erg.: Splice mit Objekten (seen) und deren Abstandsvektoren (direction),
	// die im Sichtfeld des Objekts liegen
	SeeOthers(others *[]Entity) (*[]Entity, *[]vec)

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
	IsSame(b *EntityData) bool

	// Vor.:
	// Eff.: Das Objekt beschleunigt von anderen Objekten, die im Sichtfeld liegen weg
	// Erg.:
	//avoidCollisionWithSeenObjects(others []Entity)
}
