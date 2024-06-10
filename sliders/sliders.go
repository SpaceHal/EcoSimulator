package sliders

import "github.com/hajimehoshi/ebiten/v2"

// Vor.: -
// Erg.: ein neues Slider
// New (x,y float64, current,max int, text string, active bool) *data // *data erfüllt das Interface Slider

type Slider interface {
	// Vor.: -
	// Eff.: Das Checkbox wird gezeichnet.
	// Erg.: -
	Draw(dst *ebiten.Image)
	
	// Vor.: -
	// Eff.: Das Checkbox wird aktualisiert.
	// Erg.: -
	Update()
		
	// Vor.: -
	// Eff.: -
	// Erg.: Liefert den aktuellen Wert des Sliders
	GetValue() int
	
	// Vor.: -
	// Eff.: Aendert den aktuellen Wert des Sliders
	// Erg.: -
	SetValue(v int) 	
	
	// Vor.: -
	// Eff.: Setzt, ob das Slider aktiviert ist oder nicht
	// Erg.: -
	SetActive(a bool) 

	// Vor.: -
	// Eff.: Die auszurufende Funktion beim Bewegen des Sliders wird gespeichert. 
	// Erg.: -
	SetOnMoved(f func ())
}
