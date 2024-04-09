package world

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	tilesImage *ebiten.Image
	waterImage *ebiten.Image
)

type data struct {
	width, height float32
	tileSize      int     // Ursprüngliche Kachelgröße (z.B. 16x16), wird mit `scale` skaliert
	scale         float32 // Skaliert das Hintergrundbild
	numTileX      int     // Anzahle der Kacheln pro Zeile
	numTileY      int     // Anzahle der Kacheln pro Spalte
	coastMg       float32 // Entfernung der Küste auf den Kacheln zur Kachelwand ohne Skalierung
	margin        float32 // Entfernung der Küste auf den Kacheln zur Kachelwand ohne Skalierung
	r, g, b, a    uint8

	layers [][]int

	grid  bool // aktiviert die Gitterlinien der Kacheln (Debuggen)
	debug bool // Debugmodus Tiers
}

func init() {
	img, _, err := ebitenutil.NewImageFromFile("./resources/grass.png")
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = img

	img, _, err = ebitenutil.NewImageFromFile("./resources/water.png")
	if err != nil {
		log.Fatal(err)
	}
	waterImage = img
}

func New(width float32, height float32, img *ebiten.Image) *data {
	wo := &data{
		width:    width,
		height:   height,
		coastMg:  3,
		margin:   (16 + 6) * 3,
		r:        0xeb,
		g:        0xeb,
		b:        0xeb,
		a:        0xff,
		tileSize: 16,

		//tilesImage: img,
		scale: 3,

		// Karte mit zwei Ebenen, wo welche Tiles abgebildet werden.
		// Die Tiles (16x16) werden von oben links der Reihe nach durchgezählt.
		// Das letzte Tile ist unten rechts.
		layers: [][]int{
			{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			{
				-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
				-1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, -1,
				-1, 11, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, 11, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, 11, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, 11, 12, 12, 70, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, 11, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 58, 12, 13, -1,
				-1, 11, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, 11, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, 11, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 70, 12, 12, 12, 12, 13, -1,
				-1, 11, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, 22, 23, 23, 17, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, -1, -1, -1, 11, 12, 12, 12, 12, 70, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, -1, -1, -1, 11, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, 75,
				-1, 0, 1, 1, 28, 12, 12, 58, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, 11, 12, 12, 58, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, 11, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, 11, 12, 12, 12, 12, 12, 12, 12, 12, 70, 12, 12, 12, 12, 12, 12, 12, 13, -1,
				-1, 22, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 24, -1,
				-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
			},
		},
	}

	wo.numTileX = int(wo.width) / (int(wo.tileSize) * int(wo.scale))
	wo.numTileY = int(wo.height) / (int(wo.tileSize) * int(wo.scale))

	return wo
}

// Vor.: keine
// Eff.: Ändert den Zustand von data.grid
// Erg.: keins
func (wo *data) ToggleGrid() {
	wo.grid = !wo.grid
	fmt.Println("Debug world:", wo.debug)
}

func (wo *data) Debug() {
	wo.debug = !wo.debug
	fmt.Println("Debug Animals:", wo.debug)
}

func (wo *data) GetDebug() bool {
	return wo.debug
}

/*
// Vor.:
// Eff.: Gibt für die Kachel mit den Pixelkoordinaten (x,y) an, ob in der Himmelsrichtung
// N,S,O,W eine Wasserkachel liegt.
// Erg.:
func (wo *data) GetTileBorders(x, y int) (bool, bool, bool, bool) {
	n, _, o, _, s, _, w, _ := wo.areNeighborsGround(x, y, wo.layers[1])
	n, s, o, w = !n, !s, !o, !w

	return n, s, o, w
}
*/

/*
Vor.: keine
Eff.: kein
Erg.: Die Nummer der Kachel (für den Array mit den Kacheln) und true ist geliefert.
Wenn die Kachel nicht existiert, wird -1 und false zurückgegeben.
*/
func (wo *data) getTileNumber(tileX, tileY int) (int, bool) {
	tileCount := tileX + (tileY * wo.numTileX)
	if tileX >= wo.numTileX || tileY >= wo.numTileY || tileX < 0 || tileY < 0 {
		return -1, false
	}
	if tileCount >= wo.numTileX*wo.numTileY || tileCount < 0 {
		return -1, false
	}
	return tileCount, true
}

// Überprüft, ob die Nachbarfelder (N,NO,O,SO,S,SW,W,NW) Land oder Wasser sind
func (wo *data) areNeighborsGround(tileX, tileY int, layer []int) (n, no, o, so, s, sw, w, nw bool) {

	if tileY < 0 {
		n = false
	} else {
		if t, ok := wo.getTileNumber(tileX, tileY-1); ok {
			n = layer[t] != -1
		}
	}

	// Norden-Osten
	if tileY < 0 || tileX >= wo.numTileX {
		no = false
	} else {
		if t, ok := wo.getTileNumber(tileX+1, tileY-1); ok {
			no = layer[t] != -1
		}
	}

	// Osten
	if tileX >= wo.numTileX {
		o = false
	} else {
		if t, ok := wo.getTileNumber(tileX+1, tileY); ok {
			o = layer[t] != -1
		}
	}

	// Süd-Osten
	if tileX >= wo.numTileX || tileY >= wo.numTileY {
		so = false
	} else {
		if t, ok := wo.getTileNumber(tileX+1, tileY+1); ok {
			so = layer[t] != -1
		}
	}

	// Süden
	if tileY >= wo.numTileY {
		s = false
	} else {
		if t, ok := wo.getTileNumber(tileX, tileY+1); ok {
			s = layer[t] != -1
		}
	}

	// Süd-Westen
	if tileY >= wo.numTileY || tileX < 0 {
		sw = false
	} else {
		if t, ok := wo.getTileNumber(tileX-1, tileY+1); ok {
			sw = layer[t] != -1
		}
	}

	// Westen
	if tileX < 0 {
		w = false
	} else {
		if t, ok := wo.getTileNumber(tileX-1, tileY); ok {
			w = layer[t] != -1
		}
	}

	// Mord-Westen
	if tileX < 0 || tileY < 0 {
		nw = false
	} else {
		if t, ok := wo.getTileNumber(tileX-1, tileY-1); ok {
			nw = layer[t] != -1
		}
	}

	return
}

/*
			  a
			+----+
	 	  d |    | b
			+----+
			  c     -> abcd
*/
func getState(a, b, c, d bool) int {
	return boolToInt(a)*8 + boolToInt(b)*4 + boolToInt(c)*2 + boolToInt(d)*1
}

func (wo *data) Width() float32 {
	return wo.width
}

func (wo *data) Height() float32 {
	return wo.height
}

func (wo *data) Margin() float32 {
	return wo.margin
}

func boolToInt(bit bool) int {
	var bitSetVar int
	if bit {
		bitSetVar = 1
	}
	return bitSetVar
}

func (wo *data) setLayer(x, y int, l []int, value int) {
	if t, ok := wo.getTileNumber(x, y); ok {
		l[t] = value
	}
}

func (wo *data) getLayer(x, y int, l []int) int {
	if t, ok := wo.getTileNumber(x, y); ok {
		return l[t]
	} else {
		return -1
	}
}

func (wo *data) ToggleGround(mx, my int) {
	tileX := mx / (int(wo.tileSize) * int(wo.scale))
	tileY := my / (int(wo.tileSize) * int(wo.scale))

	wo.toggle(tileX, tileY)

	// Die Nachbarfelder aktualisieren
	wo.toggle(tileX-1, tileY+1)
	wo.toggle(tileX-1, tileY+1)

	wo.toggle(tileX, tileY+1)
	wo.toggle(tileX, tileY+1)

	wo.toggle(tileX+1, tileY+1)
	wo.toggle(tileX+1, tileY+1)

	wo.toggle(tileX-1, tileY)
	wo.toggle(tileX-1, tileY)

	wo.toggle(tileX+1, tileY)
	wo.toggle(tileX+1, tileY)

	wo.toggle(tileX-1, tileY-1)
	wo.toggle(tileX-1, tileY-1)

	wo.toggle(tileX, tileY-1)
	wo.toggle(tileX, tileY-1)

	wo.toggle(tileX+1, tileY-1)
	wo.toggle(tileX+1, tileY-1)
}

func (wo *data) toggle(tileX, tileY int) {
	if tileX >= wo.numTileX || tileY >= wo.numTileY || tileX < 0 || tileY < 0 {
		return
	}

	n, no, o, so, s, sw, w, nw := wo.areNeighborsGround(tileX, tileY, wo.layers[1])
	stateOrth := getState(n, o, s, w)
	stateDiag := getState(nw, no, so, sw)

	// Wenn die gewählte Kachel Wasser ist, dann die korrekte Landkachel wählen
	tileType := -1 // keine Kachel
	if wo.getLayer(tileX, tileY, wo.layers[1]) == -1 {
		switch stateOrth {
		case 0:
			tileType = 46
		case 1:
			tileType = 35
		case 2:
			tileType = 3
		case 3:
			if sw {
				tileType = 2
			} else {
				tileType = 7
			}
		case 4:
			tileType = 33
		case 5:
			tileType = 34
		case 6:
			if so {
				tileType = 0
			} else {
				tileType = 4
			}
		case 7:
			if sw && so {
				tileType = 1
			} else if sw && !so {
				tileType = 5
			} else if !sw && so {
				tileType = 6
			} else if !sw && !so {
				tileType = 8
			}
		case 8:
			tileType = 25
		case 9:
			if nw {
				tileType = 24
			} else {
				tileType = 40
			}
		case 10:
			tileType = 14
		case 11:
			if nw && sw {
				tileType = 13
			} else if nw && !sw {
				tileType = 18
			} else if !nw && sw {
				tileType = 29
			} else if !nw && !sw {
				tileType = 51
			}
		case 12:
			if no {
				tileType = 22
			} else {
				tileType = 37
			}
		case 13:
			if nw && no {
				tileType = 23
			} else if nw && !no {
				tileType = 38
			} else if !nw && no {
				tileType = 39
			} else if !nw && !no {
				tileType = 41
			}
		case 14:
			if no && so {
				tileType = 11
			} else if no && !so {
				tileType = 15
			} else if !no && so {
				tileType = 26
			} else if !no && !so {
				tileType = 48
			}
		case 15: //
			switch stateDiag {
			case 0:
				tileType = 52
			case 1:
				tileType = 32
			case 2:
				tileType = 31
			case 3:
				tileType = 30
			case 4:
				tileType = 42
			case 5:
				tileType = 9
			case 6:
				tileType = 50
			case 7:
				tileType = 28
			case 8:
				tileType = 43
			case 9:
				tileType = 49
			case 10:
				tileType = 20
			case 11:
				tileType = 27
			case 12:
				tileType = 19
			case 13:
				tileType = 16
			case 14:
				tileType = 17
			case 15:
				tileType = 12
			}
		}
	}
	wo.setLayer(tileX, tileY, wo.layers[1], tileType)
	//fmt.Printf("tx,ty: %d,%d (%d) stateOrth: %d, stateDiag: %d \n", tileX, tileY, tileX+(tileY*numTileX), stateOrth, stateDiag)

}

func (wo *data) Draw(dst *ebiten.Image, c int) {
	//dst.Fill(color.NRGBA{wo.r, wo.g, wo.b, wo.a})
	//vector.StrokeRect(dst, wo.Margin, wo.Margin, wo.Width-2*wo.Margin, wo.Height-2*wo.Margin, 2, color.Gray{200}, true)

	nW := waterImage.Bounds().Dx() // Gibt die Breite des Tilesheets
	tileXCount := nW / wo.tileSize // Anzahl der Tiles in x Richtung im Tile Sheet (25)

	xCount := int(wo.width/wo.scale) / wo.tileSize // Anzahl der Tiles in xRichtung im SCREEN

	// ========  Layer 0 - Wasser =========
	for i, t := range wo.layers[0] {

		t = (c / 10) % 4 // Animationseffekt des Wassers
		op := &ebiten.DrawImageOptions{}

		// Ort, an den das aktuelle Tile hin geschoben werden soll
		op.GeoM.Translate(float64((i%xCount)*wo.tileSize), float64((i/xCount)*wo.tileSize))
		op.GeoM.Scale(float64(wo.scale), float64(wo.scale)) // Skaliert die 16x16 Tiles (auch die Position)

		sx := (t % tileXCount) * wo.tileSize //
		sy := (t / tileXCount) * wo.tileSize //
		dst.DrawImage(waterImage.SubImage(image.Rect(sx, sy, sx+wo.tileSize, sy+wo.tileSize)).(*ebiten.Image), op)
	}

	// ========  Layer 1 - Land =========
	nW = tilesImage.Bounds().Dx() // Gibt die Breite des Tilesheets
	tileXCount = nW / wo.tileSize // Anzahl der Tiles in x Richtung im Tile Sheet (25)

	xCount = int(wo.width/wo.scale) / wo.tileSize // Anzahl der Tiles in x Richtung im SCREEN
	for i, t := range wo.layers[1] {
		op := &ebiten.DrawImageOptions{}

		// Ort, an den das aktuelle Tile hin geschoben werden soll (ohne Skalierung)
		x := float64((i % xCount) * wo.tileSize)
		y := float64((i / xCount) * wo.tileSize)
		op.GeoM.Translate(x, y)
		op.GeoM.Scale(float64(wo.scale), float64(wo.scale)) // Skaliert die 16x16 Tiles (auch die Position)

		sx := (t % tileXCount) * wo.tileSize // x Koordinate der Kachel `t` im Style Sheet
		sy := (t / tileXCount) * wo.tileSize // y Koordinate der Kachel `t` im Style Sheet
		dst.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+wo.tileSize, sy+wo.tileSize)).(*ebiten.Image), op)
		if wo.grid {
			vector.StrokeRect(dst, float32(x)*wo.scale, float32(y)*wo.scale, float32(wo.tileSize)*wo.scale, float32(wo.tileSize)*wo.scale, 1, color.Gray{120}, false)

			vector.StrokeRect(dst, (float32(x)+wo.coastMg)*wo.scale, (float32(y)+wo.coastMg)*wo.scale, (float32(wo.tileSize)-2*wo.coastMg)*wo.scale, (float32(wo.tileSize)-2*wo.coastMg)*wo.scale, 1, color.Gray{150}, false)
		}
	}

}
