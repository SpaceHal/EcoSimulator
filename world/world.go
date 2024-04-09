package world

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Vor.: -
// Erg.: ein neuer Welt 
// New (width float32, height float32, img *ebiten.Image) *data // *data erfüllt das Interface World 

type World interface {
	// Vor.: keine
	// Eff.: Ändert den Zustand von world.grid
	// Erg.: keins
	Grid() 

	Debug()

	GetDebug() bool 
	
	// Vor.:
	// Eff.: Gibt für die Kachel mit den Pixelkoordinaten (x,y) an, ob in der Himmelsrichtung
	// N,S,O,W eine Wasserkachel liegt.
	// Erg.:
	GetTileBorders(x, y int) (bool, bool, bool, bool) 

	// Vor.: keine
	// Eff.: kein
	// Erg.: Die Nummer der Kachel (für den Array mit den Kacheln) und true ist geliefert.
	// Wenn die Kachel nicht existiert, wird -1 und false zurückgegeben.
	getTileNumber(tileX, tileY int) (int, bool) 

	// Überprüft, ob die Nachbarfelder (N,NO,O,SO,S,SW,W,NW) Land oder Wasser sind
	areNeighborsGround(tileX, tileY int, layer []int) (n, no, o, so, s, sw, w, nw bool) 
	
	Width() float32
	
	Height() float32
	
	Margin() float32
	
	setLayer(x, y int, l []int, value int) 

	getLayer(x, y int, l []int) int 

	ToggleGround(mx, my int) 
	
	toggle(tileX, tileY int) 

	Draw(dst *ebiten.Image, c int) 
}
