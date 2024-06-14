package foxes

import (
	"ecosim/entity"
	"ecosim/rabbits"
)

type Fox interface {
  entity.Entity

	// Vor.: -
	// Eff.: Die neue Postion, das aktuelle Alter und die Gesundheit ist bestimmt
	// Erg.: Ein Nachkommen ist geliefert. Gibt es kein Nachkommen, ist nil geliefert.
	Update(others *[]Fox, preys *[]rabbits.Rabbit) (offSpring *data)

	// Vor.: -
	// Eff.: -
	// Erg.: Ein Nachkommen ist geliefert. Gibt es kein Nachkommen, ist nil geliefert.
	GetOffspring() *data
}
