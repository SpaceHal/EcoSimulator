package world

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Vor.: -
// Erg.: ein neuer Welt
// New (width, height float32, img *ebiten.Image) *data // *data erfüllt das Interface World

type World interface {
	// Vor.: -
	// Eff.: Ändert, ob das Gitter gezeigt wird oder nicht.
	// Erg.: -
	ToggleGrid()

	// Vor.: -
	// Eff.: Schaltet den Debug-Modus an und aus.
	// Erg.: -
	ToggleDebug()

	// Vor.: -
	// Eff.: -
	// Erg.: Liefert den aktuellen Status des Debug-Modus.
	GetDebug() bool

	// Vor.: -
	// Eff.: -
	// Erg.: Liefert die Breite des Weltes in Pixel.
	Width() float32

	// Vor.: -
	// Eff.: -
	// Erg.: Liefert die Hoehe des Weltes in Pixel.
	Height() float32

	// Vor.: -
	// Eff.: -
	// Erg.: Liefert die Entfernung der Küste auf den Kacheln zur Kachelwand ohne Skalierung.
	Margin() float32

	// Vor.: -
	// Eff.: Das geklickte Kaestchen wir durch Klicken zwischen Boden und Wasser umgechaltet.
	// Erg.: -
	ToggleGround(mx, my int)

	// Vor.: -
	// Eff.: Der Welt wird gezeichnet.
	// Erg.: -
	Draw(dst *ebiten.Image, c int)

	// Vor.: -
	// Eff.: Die skailierte Kachelgröße ist geliefert
	// Erg.: -
	GetTileSizeScaled() int

	// Vor.: -
	// Eff.: Gibt für die Pixel (x,y) die entsprechende Kachelkoordinaten (tileX, tileY)
	// Erg.: -
	GetXYTile(x, y int) (tileX, tileY int)

	// Vor.: -
	// Eff.: Gibt für die Kachel mit den Pixelkoordinaten (x,y) an, ob in den Himmelsrichtungen
	// n, no, o, so, s, sw, w, nw eine Wasserkachel liegt.
	// Erg.: -
	GetTileBorders(x, y int) (n, no, o, so, s, sw, w, nw bool)
}
