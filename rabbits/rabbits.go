package rabbits

import (
	"ecosim/entity"
	"ecosim/grass"
)

type Rabbit interface {
	entity.Entity
	Update(others *[]Rabbit, food *[]grass.Grass) (offSpring *data)
}
