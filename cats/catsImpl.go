package cats

import (
	. "ecosim/config"
	"ecosim/entity"
	"ecosim/grass"
	"ecosim/rabbits"
	"ecosim/world"
	"math/rand"
)

type data struct {
	entity.Animal
}

func ToAnimals(cats *[]Cat) *[]entity.Animal {
	var animals []entity.Animal
	for _, f := range *cats {
		animals = append(animals, f)
	}
	return &animals
}

func New(w *world.World) *data {
	var c *data
	c = new(data)
	(*c).Animal = entity.New(w)

	c.SetImageFromFile("cats/Cat_map.png", 16, 0, 0)
	//f.SetColorRGB(204, 20, 204)

	//f.SetEnergy(900)
	c.SetHealthLoss((rand.Float64() * CatHealthLoss / 4) + CatHealthLoss)
	c.SetMaxVel(CatMaxVelocitiy)
	c.SetViewMag(CatViewMagnitude)
	c.SetViewAngle(CatViewAngle)
	c.SetMatureAge(CatMatureAge)
	c.SetLifeSpan(CatMatureAge * 2.1)

	return c
}

func (a *data) Update(others *[]Cat, preys1 *[]rabbits.Rabbit, preys2 *[]grass.Grass) (offSpring *data) {
	a.IncAge()

	preys := *(rabbits.ToAnimals(preys1))
	pr2 := *(grass.ToAnimals(preys2))
	preys = append(preys, pr2...)
	a.ApplyMove(ToAnimals(others), &preys)
	offSpring = a.GetOffspring()
	return offSpring

}

func (a *data) GetOffspring() *data {
	if a.GetAge() > a.GetDateOfLastBirth() && a.GetHealth() > 100 {
		a.SetDateOfLastBirth(a.GetAge() + rand.Intn(a.GetMatureAge()) + a.GetMatureAge())

		nCat := New(a.GetWorld())
		return nCat
	}
	return nil
}
