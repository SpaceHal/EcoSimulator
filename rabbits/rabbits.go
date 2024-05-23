package rabbits

import (
	"ecosim/entity"
	"ecosim/grass"
)

type Rabbit interface {
	entity.Animal
	Update(others *[]Rabbit, food *[]grass.Grass) (offSpring *data)
}
