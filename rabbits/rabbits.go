package rabbits

import (
	"ecosim/entity"
	"ecosim/grass"
)

type Rabbit interface {
	entity.Animal

	// Vor.: -
	// Eff.: Die neue Postion, das aktuelle Alter und die Gesundheit ist bestimmt
	// Erg.: Ein Nachkommen ist geliefert. Gibt es kein Nachkommen, ist nil geliefert.
	Update(others *[]Rabbit, food *[]grass.Grass) (offSpring *data)

	// Vor.: -
	// Eff.: -
	// Erg.: Ein Nachkommen ist geliefert. Gibt es kein Nachkommen, ist nil geliefert.
	GetOffspring() *data
}
