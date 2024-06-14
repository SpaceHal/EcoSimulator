package cats

import (
	"ecosim/entity"
	"ecosim/grass"
	"ecosim/rabbits"
)

// New(w *world.World) *data

type Cat interface {
	entity.Animal

	// Vor.: -
	// Eff.: Die neue Postion, das aktuelle Alter und die Gesundheit ist bestimmt
	// Erg.: Ein Nachkommen ist geliefert. Gibt es kein Nachkommen, ist nil geliefert.
	Update(others *[]Cat, preys1 *[]rabbits.Rabbit, preys2 *[]grass.Grass) (offSpring *data)

	// Vor.: -
	// Eff.: -
	// Erg.: Ein Nachkommen ist geliefert. Gibt es kein Nachkommen, ist nil geliefert.
	GetOffspring() *data
}
