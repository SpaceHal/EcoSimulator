package entity

import (
	"ecosim/world"

	"github.com/hajimehoshi/ebiten/v2"

	ve "github.com/quartercastle/vector"
)

type vec = ve.Vector

type LivingEntity interface {
	IsAlive() bool
	Update()
}

// Vor.: -
// Erg.: ein neues Objekt
// New (w world.World, x, y float64) *data

type Animal interface {

	// Vor.: -
	// Eff.: -
	// Erg.: True ist geliefert, wenn das Tier noch gesund ist und
	// des Maximalalter nicht erreicht ist.
	IsAlive() bool

	// Vor.:-
	// Eff.: Die Reduzierung der Gesundheit bei Zeitschritt ist gesetzt.
	// Erg.:-
	SetHealthLoss(e float64)

	// Vor.:-
	// Eff.: Die Gesundheit ist gesetzt.
	// Erg.:-
	SetHealth(e float64)

	// Vor.:-
	// Eff.:-
	// Erg.: Die Gesundheit ist gegeben.
	GetHealth() float64

	// Vor.:-
	// Eff.: Der Nähwert ist gesetzt.
	// Erg.:-
	SetHealthWhenEaten(e float64)

	// Vor.:-
	// Eff.:-
	// Erg.: Der Nähwert ist gegeben.
	GetHealthWhenEaten() float64

	// Vor.:-
	// Eff.:-
	// Erg.: Das Alter ist gegeben.
	GetAge() int

	// Vor.: -
	// Eff.: Das Alter ist im eins erhöht und die Gesundheit entsprechend `SetHealthLoss()` reduziert.
	// Erg.: -
	IncAge()

	// Vor.:-
	// Eff.: Die Geschlechtsreife ist gesetzt.
	// Erg.:-
	SetMatureAge(mAge int)

	// Vor.: -
	// Eff.: -
	// Erg.: Die Geschlechtsreife ist geliefert.
	GetMatureAge() int

	// Vor.: -
	// Eff.: Das Höchstalter ist gesetzt.
	// Erg.: -
	SetLifeSpan(ls int)

	// Vor.: -
	// Eff.: -
	// Erg.: Alter beim zuletzt erzeugten Nachwuchs ist geliefert
	GetDateOfLastBirth() int

	// Vor.: -
	// Eff.: Alter beim zuletzt erzeugten Nachwuchs ist gesetzt
	// Erg.: -
	SetDateOfLastBirth(d int)

	// Vor.: -
	// Eff.: -
	// Erg.: Zeiger auf die Simulationswelt ist geliefert.
	GetWorld() *world.World

	// Vor.: -
	// Eff.: Der Winkel des Sichfeld ist gesetzt
	// Erg.: -
	SetViewAngle(ang float64)

	// Vor.: -
	// Eff.: Die Sichtweite ist gesetzt.
	// Erg.: -
	SetViewMag(mag float64)

	// Vor.: -
	// Eff.: Der Betrag der Maximalgeschwindigkeit ist gesetzt.
	// Erg.: -
	SetMaxVel(v float64)

	// Vor.: -
	// Eff.: -
	// Erg.: Liefert die aktuelle Position des Objekts.
	GetPosition() vec

	// Vor.: -
	// Eff.: Beweglichkeit eines Objektes ist gesetzt.
	// Erg.: -
	SetMoveable(m bool)

	// Vor.: -
	// Eff.: Das Objekt hat sich auf eine angrenzende Landfläche bewegt.
	// Vermeidet dabei Kollisionen mit gleichen Objekten und vermeidet Wasser.
	// Erg.: -
	ApplyMove(others *[]Animal, preys *[]Animal)

	// Vor.: -
	// Eff.: Das Bild des Tieres wird durch die angegebene Datei ersetzt.
	// 		 Falls es ein Problem mit der angegebenen Datei gibt, wird nichts geaendert.
	// Erg.: -
	SetImageFromFile(file string, size, x, y int)

	// Vor.:
	// Eff.: Setzt die RGB Farbe für das Tier und zeichnet es neu.
	// Erg.:
	SetColorRGB(r, g, b uint8)

	// Vor.: -
	// Eff.: -
	// Erg.: Splice mit Objekten (seen) und deren Abstandsvektoren (direction),
	// die im Sichtfeld des Objekts liegen ist geliefert.
	SeeOthers(others *[]Animal) (*[]Animal, *[]vec)

	// Vor.: -
	// Eff.: Das Tier ist gezeichnet.
	// Erg.: -
	Draw(screen *ebiten.Image)

	// Vor.: -
	// Eff.: -
	// Erg.: Prueft, ob das Tier die gleichen Eigenschaften hat wie das
	// 		 uebergebenen Tier.
	IsSame(b *AnimalData) bool
}
