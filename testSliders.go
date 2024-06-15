package main

import (
	"ecosim/sliders"
	"fmt"
)

func main() {
	s := sliders.New(0,0,10,50,"Test",true)
	fmt.Println("Startwert:",s.GetValue())
	s.SetValue(35)
	fmt.Println("Wert auf 35 gesetzt.")
	fmt.Println("Aktueller Wert:",s.GetValue())
	s.SetValue(-50)
	fmt.Println("Wert auf -50 gesetzt. Das Minimum erlaubt ist 0.")
	fmt.Println("Aktueller Wert:",s.GetValue())
	s.SetValue(100)
	fmt.Println("Wert auf 100 gesetzt. Das Maximum erlaubt ist 50.")
	fmt.Println("Aktueller Wert:",s.GetValue())
	s.SetActive(false)
	fmt.Println("Slider deaktiviert.")
	fmt.Println("Aktueller Wert:",s.GetValue())
}	
