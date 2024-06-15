package grass

import (
	"ecosim/entity"
	"ecosim/world"
)

type data struct {
	entity.Entity
}

// Type Conversion
func ToAnimals(gras *[]Grass) *[]entity.Entity {
	var animals []entity.Entity
	for _, g := range *gras {
		animals = append(animals, g)
	}
	return &animals
}

func New(w *world.World) *data {
	var gr *data = new(data)
	(*gr).Entity = entity.New(w)
	gr.SetColorRGB(10, 150, 10)
	gr.SetHealthLoss(0)
	gr.SetMoveable(false)
	gr.SetHealthWhenEaten(100)

	return gr
}
