package grass

import "ecosim/entity"

// Vor.: -
// Erg.: neues Grass
// New (x, y int) *data // *data erfüllt das Interface Grass

type Grass interface {
	entity.Animal
}
