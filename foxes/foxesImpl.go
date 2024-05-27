package foxes

import (
	. "ecosim/config"
	"ecosim/entity"
	"ecosim/rabbits"
	"ecosim/world"
	"math/rand"
)

type data struct {
	entity.Animal
}

// Type Conversion
func ToAnimals(foxes *[]Fox) *[]entity.Animal {
	var animals []entity.Animal
	for _, f := range *foxes {
		animals = append(animals, f)
	}
	return &animals
}

func New(w *world.World) *data {
	var f *data
	f = new(data)
	(*f).Animal = entity.New(w)

	f.SetImageFromFile("foxes/fox.png")
	//f.SetColorRGB(204, 20, 204)

	//f.SetEnergy(900)
	f.SetHealthLoss((rand.Float64() * FoxHealthLoss / 4) + FoxHealthLoss)
	f.SetMaxVel(FoxMaxVelocitiy)
	f.SetViewMag(FoxViewMagnitude)
	f.SetViewAngle(FoxViewAngle)
	f.SetMatureAge(FoxMatureAge)

	return f
}

// Die neue Position e.pos aus e.vel und e.acc bestimmen und die Lebensenergie aktualisieren
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
		//fmt.Println("Geburt Fuchs: Eltern-Energy", a.GetHealth(), ", E-Alter", a.GetAge(), ", Kind-Energy", nFox.GetHealth(), ", K-Alter", nFox.GetAge())

		return nFox
	}
	return nil
}
