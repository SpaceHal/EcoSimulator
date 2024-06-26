package rabbits

import (
	. "ecosim/config"
	"ecosim/entity"
	"ecosim/grass"
	"ecosim/world"
	"math/rand"
)

type data struct {
	entity.Entity
}

// Type Conversion
func ToAnimals(rabbits *[]Rabbit) *[]entity.Entity {
	var animals []entity.Entity
	for _, r := range *rabbits {
		animals = append(animals, r)
	}
	return &animals
}

func New(w *world.World) *data {
	var r *data
	r = new(data)
	(*r).Entity = entity.New(w)

	r.SetImageFromFile("rabbits/rabbit.png", 0, 0, 0)
	//f.SetColorRGB(21, 123, 220)
	r.SetHealthLoss((rand.Float64() * BunnyHealthLoss / 4) + BunnyHealthLoss)
	r.SetMaxVel(BunnyMaxVelocitiy)
	r.SetViewMag(BunnyViewMagnitude)

	r.SetHealthWhenEaten(100)
	r.SetViewAngle(BunnyViewAngle)
	r.SetMatureAge(BunnyMatureAge)

	return r
}

func (a *data) Update(others *[]Rabbit, food *[]grass.Grass) (offSpring *data) {
	a.IncAge()

	a.ApplyMove(ToAnimals(others), grass.ToAnimals(food))
	offSpring = a.GetOffspring()
	return offSpring

}

func (a *data) GetOffspring() *data {
	if a.GetAge() > a.GetDateOfLastBirth() {
		a.SetDateOfLastBirth(a.GetAge() + rand.Intn(a.GetMatureAge()))
		newR := New(a.GetWorld())
		return newR
	}
	return nil
}
