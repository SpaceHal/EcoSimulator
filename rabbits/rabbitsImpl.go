package rabbits

import (
	. "ecosim/config"
	"ecosim/entity"
	"ecosim/grass"
	"ecosim/world"
	"math/rand"
)

type data struct {
	entity.Animal
}

// Type Conversion
func ToAnimals(rabbits *[]Rabbit) *[]entity.Animal {
	var animals []entity.Animal
	for _, r := range *rabbits {
		animals = append(animals, r)
	}
	return &animals
}

func New(w *world.World) *data {
	var r *data
	r = new(data)
	(*r).Animal = entity.New(w)

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

// Die neue Position e.pos aus e.vel und e.acc bestimmen und die Lebensenergie aktualisieren
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
