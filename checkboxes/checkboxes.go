package checkboxes

import "github.com/hajimehoshi/ebiten/v2"

// Vor.: -
// Erg.: ein neues Checkbox
// New (x,y float64, text string, checked bool) *data // *data erfüllt das Interface Checkbox

type Checkbox interface {
	// Vor.: -
	// Eff.: Das Checkbox wird gezeichnet.
	// Erg.: -
	Draw(dst *ebiten.Image)
	
	// Vor.: -
	// Eff.: -
	// Erg.: Liefert den aktuellen Status des Checkboxes
	IsChecked() bool
	
	// Vor.: -
	// Eff.: Das Checkbox wird aktualisiert.
	// Erg.: -
	Update()
	
	// Vor.: -
	// Eff.: Die auszurufende Funktion beim Clicken auf das Checkbox wird gespeichert. 
	// Erg.: -
	SetOnClicked(f func ())
}
