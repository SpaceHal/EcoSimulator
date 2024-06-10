package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Vor.: -
// Erg.: ein neuer UI
// New () *data // *data erfüllt das Interface UI

type UI interface {
	// Vor.: -
	// Erg.: Der gewuenschte Anzahl der Grassflaechen am Anfang der Simulation ist geliefert. 
	GetNumberOfGrass () int
	
	// Vor.: -
	// Erg.: Der gewuenschte Anzahl der Hasen am Anfang der Simulation ist geliefert. 
	GetNumberOfBunnies () int

	// Vor.: -
	// Erg.: Der gewuenschte Anzahl der Katzen am Anfang der Simulation ist geliefert. 
	GetNumberOfCats () int

	// Vor.: -
	// Erg.: Der gewuenschte Anzahl der Fuechse am Anfang der Simulation ist geliefert. 
	GetNumberOfFoxes () int
	
	// Vor.: -
	// Eff.: Der UI wird gezeichnet.
	// Erg.: -
	Draw(dst *ebiten.Image)
	
	// Vor.: -
	// Eff.: Der UI wird aktualisiert.
	// Erg.: -
	Update()
}
