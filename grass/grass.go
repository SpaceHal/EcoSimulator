package grass

import (
	"ecosim/entity"
	"ecosim/world"
)

type data struct {
	entity.Animal
}

// Type Conversion
func ToAnimals(gras *[]Grass) *[]entity.Animal {
	var animals []entity.Animal
	for _, g := range *gras {
		animals = append(animals, g)
	}
	return &animals
}

func New(w *world.World) *data {
	var gr *data = new(data)
	(*gr).Animal = entity.New(w)

	//f.SetImageFromFile("rabbits/rabbit.png")
	gr.SetColorRGB(10, 150, 10)
	gr.SetHealthLoss(0)
	gr.SetMoveable(false)
	gr.SetHealthWhenEaten(100)

	return gr
}

func (gr *data) Update() {

}
