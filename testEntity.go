package main

import (
	"ecosim/world"
	"ecosim/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"fmt"
)

func main() {
	var image *ebiten.Image
	var w world.World
	w = world.New(200,200,2,image)
	e:=entity.New(&w)
	fmt.Println("IsAlive:",e.IsAlive())
	
	fmt.Println("Aktuelle Gesundheit:",e.GetHealth())
	e.SetHealth(350)
	fmt.Println("Gesundheit auf 350 gesetzt.")
	fmt.Println("Neue Gesundheit:",e.GetHealth())
	
	fmt.Println("Aktueller Nähwert:",e.GetHealthWhenEaten())
	e.SetHealthWhenEaten(50)
	fmt.Println("Nähwert auf 50 gesetzt.")
	fmt.Println("Neuer Nähwert:",e.GetHealthWhenEaten())
	
	e.SetHealthLoss(20)
	fmt.Println("Health loss auf 20 gesetzt.")
	for i:=0;i<=5;i++ {
		e.IncAge()
		fmt.Println("Alter:",e.GetAge(),",Gesundheit:",e.GetHealth())
	}
	// CHECK Reduktion immer 17,5??
	
	e.SetMatureAge(5)
	fmt.Println("Die Geschlechtsreife wurde auf 5 gesetzt.")
	fmt.Println("Geschlechtsreife:",e.GetMatureAge())
	
	f:=entity.New(&w)
	fmt.Println("Sind e und e die selben Entitäten?:",e.IsSame(e))
	fmt.Println("Sind e und f die selben Entitäten?:",e.IsSame(f))
	
	fmt.Println("Gesundheit:",e.GetHealth())
	fmt.Println("IsAlive:",e.IsAlive())
	e.SetHealth(0)
	fmt.Println("Gesundheit auf 0 gesetzt.")
	fmt.Println("IsAlive:",e.IsAlive())

	fmt.Println("Aktuelle Position:",f.GetPosition())
	
	f.SetLifeSpan(8)
	fmt.Println("Höchstalter auf 8 gesetzt.")
	for f.IsAlive() {
		f.IncAge()
		fmt.Println("Alter:",f.GetAge(),"Gesundheit:,",f.GetHealth(),"IsAlive:",f.IsAlive())
	}
}
