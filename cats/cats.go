package cats

import (
	"ecosim/entity"
	"ecosim/grass"
	"ecosim/rabbits"
)

type Cat interface {
	entity.Entity
	Update(others *[]Cat, preys1 *[]rabbits.Rabbit, preys2 *[]grass.Grass) (offSpring *data)
}
