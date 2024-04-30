package foxes

import (
	"ecosim/animal"
	"ecosim/world"
	"math/rand"
)

type data struct {
	animal.Animal
}

func New(w *world.World, x, y float64, preys *[]animal.Animal) *data {
	var f *data
	f = new(data)
	(*f).Animal = animal.New(w, x, y)
	
	f.SetImageFromFile("foxes/fox.png")
	f.SetEnergyLoss(((rand.Float64()*30 + 60) * 60))
	f.SetPreys(preys)
	f.SetMaxVel(0.7)
	f.SetViewMag(200)

	return f
}
