package foxes

import (
	"ecosim/entity"
	"ecosim/rabbits"
)

type Fox interface {
	entity.Entity
	Update(others *[]Fox, preys *[]rabbits.Rabbit) (offSpring *data)
}
