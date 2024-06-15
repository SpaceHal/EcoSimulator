package foxes

import (
	. "ecosim/config"
	"ecosim/entity"
	"ecosim/rabbits"
	"ecosim/world"
	"math/rand"
)

type data struct {
	entity.Entity
}

// Type Conversion
func ToAnimals(foxes *[]Fox) *[]entity.Entity {
	var animals []entity.Entity
	for _, f := range *foxes {
		animals = append(animals, f)
	}
	return &animals
}

func New(w *world.World) *data {
	var f *data
	f = new(data)
	(*f).Entity = entity.New(w)

	f.SetImageFromFile("foxes/fox.png", 0, 0, 0)
	//f.SetColorRGB(204, 20, 204)

	//f.SetEnergy(900)
	f.SetHealthLoss((rand.Float64() * FoxHealthLoss / 4) + FoxHealthLoss)
	f.SetMaxVel(FoxMaxVelocitiy)
	f.SetViewMag(FoxViewMagnitude)
	f.SetViewAngle(FoxViewAngle)
	f.SetMatureAge(FoxMatureAge)

	return f
}

func (a *data) Update(others *[]Fox, preys *[]rabbits.Rabbit) (offSpring *data) {
	a.IncAge()
	a.ApplyMove(ToAnimals(others), rabbits.ToAnimals(preys))
	offSpring = a.GetOffspring()
	return offSpring
}

func (a *data) GetOffspring() *data {
	if a.GetAge() > a.GetDateOfLastBirth() && a.GetHealth() > 100 {
		a.SetDateOfLastBirth(a.GetAge() + rand.Intn(a.GetMatureAge()) + a.GetMatureAge())
		nFox := New(a.GetWorld())
		return nFox
	}
	return nil
}
