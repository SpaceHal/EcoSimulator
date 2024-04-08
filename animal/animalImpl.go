package eater

import (
	"ecosim/world"

	//"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	//termC "github.com/fatih/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	ve "github.com/quartercastle/vector"
)

type vec = ve.Vector //  Vektoren

type Animal struct {
	//dir                  float64 // Bewegungsrichtung
	pos, vel, acc, z, z1 vec     // Position, Geschwindigkeit, Beschleunigung (temp Werte)
	accBorder            float64 // Beschleunigung weg vom Rand
	maxVel               float64 // Betrag der Maximalgeschwindigkeit,
	absVel               float64 // Betrag der aktuellen Geschwindigkeit
	ahead                float64 // Abstand des "Ziehpunkts" (die Deichsel), an dem die Beschleunigung ansetzt, zum Objekt.
	maxAccPhi            float64 // maximale Winkeländerung für die Beschleunigung auf den Ziehpunkt
	accPhi               float64
	eps                  float64 // Elastizität (Impulserhaltung)

	viewAngle float64 // Öffnungswinkel des Sichtfelds
	viewMag   float64 // Sichtweite
	inView    bool    // wenn etwas im Sichtfeld ist
	atWater   bool    // wenn das Objekt am sich am Wasser befindet

	img, imgDebug *ebiten.Image // das zu zeigende Bild

	w *world.World // Die Simulationswelt

	r, g, b, a          uint8   // Farbe rot, grün, blau, alpha des Objekts
	imgWidth, imgHeight float32 // Größe des Tierbilds
	line, aa            bool
	debug               bool
}

var (
	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

// Wird genau einmal aufgerufen, wenn das Packet verwendet wird
func init() {
	whiteImage.Fill(color.White)
}

func New(world *world.World, x, y float64) *Animal {

	e := &Animal{
		imgWidth:  7,
		imgHeight: 20,
		r:         0xa0,
		g:         0x00,
		b:         0x50,
		a:         0x60,
		accBorder: 0.5,
		maxVel:    0.5,
		absVel:    1,
		ahead:     1,
		pos:       vec{x, y},
		vel:       vec{1, rand.Float64() * 100},
		accPhi:    rand.Float64() * math.Pi * 2,
		maxAccPhi: math.Pi / 8,
		eps:       0.1,
		viewAngle: math.Pi / 6,
		viewMag:   80,
		aa:        true,
		w:         world,
		debug:     true,
	}
	e.acc = vec{1, 1}.Unit().Scale(e.ahead / 8)

	e.makeAnimal() // erzeugt das Bild vom Tier

	if e.debug {
		size := math.Max(float64(e.imgHeight), float64(e.imgHeight))
		e.imgDebug = ebiten.NewImage(int(20*size), int(20*size))
	}

	return e
}

func (a *Animal) isOutside() bool {

	return float32(a.pos.X()) >= a.w.Width-(2*a.imgWidth)-a.w.Margin ||
		float32(a.pos.X()) <= 0+a.w.Margin ||
		float32(a.pos.Y()) >= a.w.Height-(2*a.imgHeight)-a.w.Margin ||
		float32(a.pos.Y()) <= 0+a.w.Margin
}

// Die neue Position e.pos aus e.vel und e.acc bestimmen.
func (a *Animal) Update(others []*Animal) {
	a.randomStep()
	a.repelAtBorder()
	a.avoidCollisionWithSeenObjekts(others)
	a.applyMove(others)

}

// Vor.:
// Eff.: Die neue Position animal.pos ist bestimmt. Das Überlappen von
// Objekten wird vermieden
// Erg.:
func (a *Animal) applyMove(others []*Animal) {
	//a.vel.Unit().Scale(a.maxVel)
	//a.pos = a.pos.Add(a.vel)
	newPos := a.pos.Add(a.vel)

	collission := false
	sumDiff := vec{0, 0}
	var counts float64
	for _, other := range others {
		dist := newPos.Sub(other.pos)
		if a != other && dist.Magnitude() <= float64(a.imgHeight*1.1) {

			collission = collission || true
			// Einfaches Separieren
			sumDiff = sumDiff.Add(dist.Unit())
			counts++
		}
	}
	if counts > 0 {
		sumDiff = sumDiff.Scale(1 / counts)
		sumDiff = sumDiff.Scale(0.5)
		//fmt.Println(sumDiff)
		a.vel = a.vel.Add(sumDiff)
	}
	if !collission {
		a.pos = newPos
		//a.b = 0x50
	} else {
		a.pos = a.pos.Add(a.vel)
		//a.b = 0xff
	}
	//a.makeAnimal()

	// Grenzen der Welt beachten
	if float32(a.pos[0]) >= a.w.Width-a.w.Margin {
		a.pos[0] = float64(a.w.Width - a.w.Margin)
	} else if float32(a.pos.X()) <= a.w.Margin {
		a.pos[0] = float64(a.w.Margin)
	}

	if float32(a.pos.Y()) >= a.w.Height-a.w.Margin {
		a.pos[1] = float64(a.w.Height - a.w.Margin)
	} else if float32(a.pos.Y()) <= a.w.Margin {
		a.pos[1] = float64(a.w.Margin)
	}

}

// Vor.: ?
// Eff.: Bestimmt den Geschwindigkeitsvektor Animal.vel
// Erg.:
func (a *Animal) randomStep() {

	// Ein Vektor, kurz vor das Objekt zeigend, ist der Punkt, an dem gezogen wird.
	z := a.vel.Unit().Scale(a.ahead)

	// zufällige Richtungsänderung der Beschleunigung
	// TODO: Perlin - Noise
	a.accPhi = (rand.Float64()*2 - 1) * a.maxAccPhi // zufällliger Winkel 0 ... ??
	a.acc = a.acc.Rotate(a.accPhi).Unit().Scale(a.ahead / 8)

	// TODO: Die Länge von z (Ziehpunkt) variiert mit absVel
	z = z.Add(a.acc)
	// Der Betrag der Geschwindigkeit soll sich kontinuierlich aber nicht sprunghaft ändern.
	/*
		a.absVel += (rand.Float64() - 0.5) / 4
		if a.absVel > a.maxVel {
			a.absVel = a.maxVel
			} else if a.absVel < 0 {
				a.absVel = 0
			}
	*/

	a.vel = a.vel.Unit().Scale(a.maxVel).Add(z)
}

func (a *Animal) isAtWater() bool {
	return a.atWater
}

// Vor.:
// Eff.: Stößt das Objekt von der Wand ab
// Erg.:
func (a *Animal) repelAtBorder() {
	// Wenn das Objekt in die Nähe des Bildschirmrandes kommt,
	// wird es senkrecht dazu beschleunigt (dreht also um)
	// TODO: Die Beschleunigung vom Rand weg sollte Proportional zur Entfernung zum Rand sein.

	//n, s, o, w := a.w.GetTileBorders(int(a.pos[0]), int(a.pos[1]))
	//fmt.Println("Grenzen:", n, s, o, w)
	repel := vec{0, 0}
	const d = 6
	a.atWater = false
	if float32(a.pos.X()) >= a.w.Width-(a.w.Margin+d) {
		repel[0] = -a.accBorder
		a.atWater = true
	} else if float32(a.pos.X()) <= a.w.Margin+d {
		repel[0] = a.accBorder
		a.atWater = true
	}

	if float32(a.pos.Y()) >= a.w.Height-(a.w.Margin+d) {
		repel[1] = -a.accBorder
		a.atWater = true
	} else if float32(a.pos.Y()) <= a.w.Margin+d {
		repel[1] = a.accBorder
		a.atWater = true
	}
	a.vel = a.vel.Add(repel)
}

// Vor.:
// Eff.: Das Objekt beschleunigt von anderen Objekten, die im Sichtfeld liegen weg
// Erg.:
func (a *Animal) avoidCollisionWithSeenObjekts(others []*Animal) {
	avg := vec{0, 0}

	_, dirs := a.SeeOthers(others)
	for _, dir := range dirs {
		dir = dir.Unit()
		dir = dir.Scale(a.viewMag - dir.Magnitude())
		avg = avg.Add(dir)
	}
	if len(dirs) > 0 {
		avg.Scale(1 / float64(len(dirs)))
	}

	z := a.vel.Unit().Scale(a.ahead)
	z = z.Add(avg.Unit().Scale(-1))
	a.vel = a.vel.Unit().Scale(a.maxVel).Add(z)

}

// Vor.: ?
// Eff.: ?
// Erg.: Splice mit Objekten (seen) und deren Abstandsvektoren (direction),
// die im Sichtfeld des Objekts liegen
func (a *Animal) SeeOthers(others []*Animal) (seen []*Animal, direction []vec) {
	inView := false
	for _, other := range others {
		delta := other.pos.Sub(a.pos)
		if a != other && delta.Magnitude() < a.viewMag {

			if math.Abs(a.vel.Angle(delta)) < a.viewAngle {
				inView = inView || true
				seen = append(seen, a)
				direction = append(direction, delta.Clone())
			}
		} else {
			inView = inView || false
		}
	}
	a.inView = inView
	return seen, direction
}

func (a *Animal) Draw(screen *ebiten.Image) {
	a.drawAnimal(screen)
}

// Vor.: ?
// Eff.: Zeichnet ein Tier als Vektorgrafik mit Sichtfeld und Geschwindigkeitsvektor
// Erg.: ?
func (a *Animal) drawAnimal(screen *ebiten.Image) {
	//halfImg := e.imgHeight / 2
	if a.w.GetDebug() {
		w := float32(a.imgDebug.Bounds().Dx())
		h := float32(a.imgDebug.Bounds().Dy())
		a.imgDebug.Clear()

		// --- Rahmen ---
		/*
			vector.StrokeRect(e.imgDebug, 1, 1, w-1, h-1, 1, color.Gray{100}, false)
			opRect := &ebiten.DrawImageOptions{}
			opRect.GeoM.Translate(e.pos[0]+float64(-w/2), e.pos[1]-float64(h/2)) // ... und zum Schluss ein die gewünschte Stelle bewegen
			screen.DrawImage(e.imgDebug, opRect)
			e.imgDebug.Clear()
		*/

		// --- Geschwindigkeit ---
		length := float32(a.vel.Magnitude()) * 20
		vector.StrokeLine(a.imgDebug, w/2, h/2, w/2, h/2-length, 1.5, color.Gray{0}, false)
		opD := &ebiten.DrawImageOptions{}
		dirAngle := a.vel.Angle(vec{0, -1}) // Winkel zur Y-Achse bestimmen

		// --- Sichtfeld ---
		viewColor := color.NRGBA{120, 180, 100, 80}
		if a.inView {
			viewColor = color.NRGBA{150, 100, 180, 80}
		}
		a.makeArc(a.imgDebug, float32(a.viewMag), float32(-math.Pi/2-a.viewAngle), float32(-math.Pi/2+a.viewAngle), viewColor, false)

		// --- Transformation ---
		opD.GeoM.Translate(float64(-w/2), -float64(h/2)) // Koordinaten zuerst in die Mitte des Bilder bewegen ...
		opD.GeoM.Rotate(-dirAngle)                       // ... dann drehen ...
		opD.GeoM.Translate(a.pos[0], a.pos[1])           // ... und zum Schluss ein die gewünschte Stelle bewegen
		screen.DrawImage(a.imgDebug, opD)

	}
	// Das zu zeichnende Bild ist in e.img gespeichert.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(a.imgHeight)/2, -float64(a.imgHeight)/2) // Koordinaten zuerst in die Mitte des Bilder bewegen ...
	op.GeoM.Translate(a.pos[0], a.pos[1])                               // ... und zum Schluss ein die gewünschte Stelle bewegen
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(a.img, op)

}

// Vor.: ?
// Eff.: Ein Kreisbogen wird erzeugt, dessen Mittelpunkt mittig in `img` gespeichert wird.
// Erg.:
func (a *Animal) makeArc(img *ebiten.Image, radius float32, startAngle, endAngle float32, c color.NRGBA, line bool) {
	w := float32(img.Bounds().Dx())
	h := float32(img.Bounds().Dy())

	var path vector.Path

	path.MoveTo(w/2, h/2)
	path.Arc(w/2, h/2, radius, startAngle, endAngle, vector.Clockwise)
	path.Close()

	opArc := &vector.StrokeOptions{}
	opArc.Width = 1
	opArc.LineJoin = vector.LineJoinRound

	var vs []ebiten.Vertex
	var is []uint16
	if line {
		op := &vector.StrokeOptions{}
		op.Width = 1
		op.LineJoin = vector.LineJoinRound
		vs, is = path.AppendVerticesAndIndicesForStroke(nil, nil, opArc)
	} else {
		vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	}

	for i := range vs {
		//vs[i].DstX = (vs[i].DstX + float32(getX(e.pos)))
		//vs[i].DstY = (vs[i].DstY + float32(getY(e.pos)))
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		// ColorR/ColorG/ColorB/ColorA represents color scaling values.
		// 1 means the original source image color is used.
		// 0 means a transparent color is used.
		vs[i].ColorR = float32(c.R) / 0xff
		vs[i].ColorG = float32(c.G) / 0xff
		vs[i].ColorB = float32(c.B) / 0xff
		vs[i].ColorA = float32(c.A) / 0xff

	}
	// --- Kantenglättung ---
	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true
	if !line {
		op.FillRule = ebiten.EvenOdd
	}

	img.DrawTriangles(vs, is, whiteSubImage, op)
}

// Vor.: ?
// Eff.: Erstellt ein Bild für ein Tier. Das Bild wird in animal.img gespeichert und
// später mit Animal.DrawShape() jedes mal neu gezeichnet.
// Erg.:
func (a *Animal) makeAnimal() {
	a.img = ebiten.NewImage(int(a.imgHeight), int(a.imgHeight))
	vector.DrawFilledCircle(a.img, a.imgHeight/2, a.imgHeight/2, a.imgHeight/2, color.NRGBA{a.r, a.g, a.b, a.a}, true)
}

// Ein Dreieck als image
func (a *Animal) createTriangle() {

	var path vector.Path

	path.MoveTo(a.imgWidth, 0)
	path.LineTo(2*a.imgWidth, 2*a.imgHeight)
	path.LineTo(0*a.imgWidth, 2*a.imgHeight)
	path.LineTo(a.imgWidth, 0)
	path.Close()

	var vs []ebiten.Vertex
	var is []uint16
	if a.line {
		op := &vector.StrokeOptions{}
		op.Width = 5
		op.LineJoin = vector.LineJoinRound
		vs, is = path.AppendVerticesAndIndicesForStroke(nil, nil, op)
	} else {
		vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	}

	for i := range vs {
		//vs[i].DstX = (vs[i].DstX + float32(getX(e.pos)))
		//vs[i].DstY = (vs[i].DstY + float32(getY(e.pos)))
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(a.r) / 0xff
		vs[i].ColorG = float32(a.g) / 0xff
		vs[i].ColorB = float32(a.b) / 0xff
		vs[i].ColorA = float32(a.a) / 0xff
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = a.aa
	if !a.line {
		//op.FillRule = ebiten.EvenOdd // Für Komplexe konkave Formen
		op.FillRule = ebiten.FillAll
	}

	a.img = ebiten.NewImage(int(2*a.imgWidth), int(2*a.imgHeight))
	//vector.DrawFilledRect(e.img, 0, 0, 2*e.imgWidth, 2*e.imgHeight, color.NRGBA{0xff, 0x00, 0x00, 0xff}, false)

	a.img.DrawTriangles(vs, is, whiteSubImage, op)

}
