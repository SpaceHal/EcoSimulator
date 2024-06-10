package entity

import (
	"bytes"
	"ecosim/world"
	"fmt"
	"log"

	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

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

func rand_ab(a, b int) float64 {
	return float64(a + rand.Intn(b-a))
}

/////////////////////////////////////////////////////
//   DrawableData
/////////////////////////////////////////////////////

type DrawableData struct {
	w                   *world.World // Die Simulationswelt
	pos                 vec          // Position
	img                 *ebiten.Image
	r, g, b, a          uint8   // Farbe rot, grün, blau, alpha des Objekts
	imgWidth, imgHeight float32 // Größe des Bildes
	imgIsSetFromFile    bool    // true, wenn ein Image von einer Datei geladen wurde
	aa                  bool
}

func NewDrawable(w *world.World) DrawableData {
	ok := false
	var rx, ry float64
	var count = 0
	// 1000 Versuche, um eine Stelle auf dem Land zu finden
	for !ok || count < 1000 {
		rx = rand_ab((*w).GetTileSizeScaled(), int((*w).Width())-(*w).GetTileSizeScaled())
		ry = rand_ab((*w).GetTileSizeScaled(), int((*w).Height())-(*w).GetTileSizeScaled())
		ok = (*w).IsLand(int(rx), int(ry))
		count++
	}

	return DrawableData{
		pos:       vec{rx, ry},
		w:         w,
		r:         0xa0,
		g:         0x00,
		b:         0x50,
		a:         0x60,
		imgWidth:  32,
		imgHeight: 32,
		aa:        true,
	}
}

func (a *DrawableData) SetImageFromFile(file string, size, x, y int) {
	var img *ebiten.Image
	var err error
	img, _, err = ebitenutil.NewImageFromFile(file)
	if err == nil {
		if size == 0 {
			a.img = img
		} else {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			imgN := img.SubImage(image.Rect(size*x, size*y, size*(x+1), size*(y+1))).(*ebiten.Image)
			a.img.Clear()
			a.img.DrawImage(imgN, op)

		}
		a.imgIsSetFromFile = true

	}
}

// Vor.: ?
// Eff.: Erstellt ein Bild für ein Tier. Das Bild wird in entity.img gespeichert und
// später mit Entity.DrawShape() jedes mal neu gezeichnet.
// Erg.:
func (a *DrawableData) makeEntity() {
	a.img = ebiten.NewImage(int(a.imgHeight), int(a.imgHeight))
	vector.DrawFilledCircle(a.img, a.imgHeight/2, a.imgHeight/2, a.imgHeight/2, color.NRGBA{a.r, a.g, a.b, a.a}, true)
}

func (a *DrawableData) SetColorRGB(r, g, b uint8) {
	a.r = r
	a.g = g
	a.b = b
	a.makeEntity()
}

// Liefert die aktuelle Position
func (a *DrawableData) GetPosition() vec {
	return a.pos
}

/////////////////////////////////////////////////////
//   MoveableData
/////////////////////////////////////////////////////

type MoveableData struct {
	DrawableData
	vel, acc  vec     // Position, Geschwindigkeit, Beschleunigung (temp Werte)
	accBorder float64 // Beschleunigung weg vom Rand
	maxVel    float64 // Betrag der Maximalgeschwindigkeit,
	absVel    float64 // Betrag der aktuellen Geschwindigkeit
	ahead     float64 // Abstand des "Ziehpunkts" (die Deichsel), an dem die Beschleunigung ansetzt, zum Objekt.
	maxAccPhi float64 // maximale Winkeländerung für die Beschleunigung auf den Ziehpunkt
	accPhi    float64
	eps       float64 // Elastizität (Impulserhaltung)

	viewAngle float64 // Öffnungswinkel des Sichtfelds
	viewMag   float64 // Sichtweite
	inView    bool    // wenn etwas im Sichtfeld ist

	moveable bool // Ist das Objekt beweglich
}

func NewMoveable(w *world.World) MoveableData {
	mD := MoveableData{
		DrawableData: NewDrawable(w),
		accBorder:    0.8,
		absVel:       1,
		ahead:        1,
		vel:          vec{1, rand.Float64() * 100},
		accPhi:       rand.Float64() * math.Pi * 2,
		maxAccPhi:    math.Pi / 10,
		eps:          0.1,
		viewAngle:    math.Pi / 6,
		moveable:     true,
	}

	mD.acc = vec{1, 1}.Unit().Scale(mD.ahead / 8)
	return mD
}

// Vor.: ?
// Eff.: Bestimmt den Geschwindigkeitsvektor Entity.vel
// Erg.:
func (a *MoveableData) randomStep() {

	// Ein Vektor, kurz vor das Objekt zeigend, ist der Punkt, an dem gezogen wird.
	z := a.vel.Unit().Scale(a.ahead)

	// zufällige Richtungsänderung der Beschleunigung
	// TODO: Perlin - Noise
	a.accPhi = (rand.Float64()*2 - 1) * a.maxAccPhi // zufällliger Winkel 0 ... ??
	a.acc = a.acc.Rotate(a.accPhi).Unit().Scale(a.ahead / 8)

	// TODO: Die Länge von z (Ziehpunkt) variiert mit absVel
	z = z.Add(a.acc)

	a.vel = a.vel.Unit().Scale(a.maxVel).Add(z)
}

/////////////////////////////////////////////////////
//   HealthData
/////////////////////////////////////////////////////

type HealthData struct {
	health          float64 // Lebensenergie zwischen 100 und 0 (0 sterben)
	healthLoss      float64 // Reduziert die Lebensenergie (energy) pro update (ageingFactor <= 1.0)
	healthWhenEaten float64 // Energie, die das Objekt beim Essen liefert
	age             int     // Lebenssekunden
	lifeSpan        int     // Lebenserwartung
	matureAge       int     // ab wann Nachwuchs möglich ist (Frames aka Sekunden)
	birthNotBefore  int
}

func NewHealth() HealthData {
	hd := HealthData{
		health:    100,
		matureAge: 240,
	}
	hd.healthLoss = hd.health / 60
	hd.lifeSpan = hd.matureAge * 3
	hd.birthNotBefore = rand.Intn(hd.matureAge) + hd.matureAge

	return hd
}

func (a *HealthData) IsAlive() bool {
	return a.health > 0 && a.age < a.lifeSpan
}

func (a *HealthData) SetLifeSpan(ls int) {
	a.lifeSpan = ls
}

func (a *HealthData) SetHealthLoss(e float64) {
	a.healthLoss = a.health / e
}
func (a *HealthData) SetHealth(e float64) {
	a.health = e
}

func (a *HealthData) GetHealth() float64 {
	return a.health
}

func (a *HealthData) SetHealthWhenEaten(e float64) {
	a.healthWhenEaten = e
}
func (a *HealthData) GetHealthWhenEaten() float64 {
	return a.healthWhenEaten
}

func (a *HealthData) SetMatureAge(mAge int) {
	a.matureAge = mAge
}

// Vor.:
// Eff.: Das Alter ist im 1 erhöht, die Energie reduziert
// Erg.:
func (a *HealthData) IncAge() {
	a.age++
	a.health -= a.healthLoss
}

// ///////////////////////////////////////////////////
//	EntityData
// ///////////////////////////////////////////////////

type EntityData struct {
	//DrawableData
	MoveableData
	HealthData

	//preys     *[]Entity // die Beute des Tiers (was wird gegessen)
	//predators *[]Entity // die Jäger des Tiers (wovor wird geflüchtet)

	imgDebug *ebiten.Image // das zu zeigende Bild

	debug bool
	font  *text.GoTextFaceSource // Font für Debug-Text
}

func New(w *world.World) *EntityData {

	a := &EntityData{
		MoveableData: NewMoveable(w),
		HealthData:   NewHealth(),
		debug:        true,
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	a.font = s

	if a.debug {
		size := math.Max(float64(a.imgHeight), float64(a.imgHeight))
		a.imgDebug = ebiten.NewImage(int(20*size), int(20*size))
	}

	a.makeEntity()

	return a
}

func (a *EntityData) IsSame(b *EntityData) bool {
	return a == b
}

func (a *EntityData) SetMoveable(m bool) {
	a.moveable = m
}

func (a *EntityData) SetViewAngle(ang float64) {
	a.viewAngle = ang
}

/*
func (a *EntityData) SetPreys(preys *[]Entity) {
	a.preys = preys
}
func (a *EntityData) GetNumOfPreys() int {
	return len(*a.preys)
}

func (a *EntityData) SetPredators(predators *[]Entity) {
	a.predators = predators
}
func (a *EntityData) GetPreys() *[]Entity {
	return a.preys
}
func (a *EntityData) GetPredators() *[]Entity {
	return a.predators
}
*/

func (a *EntityData) SetMaxVel(v float64) {
	a.maxVel = v
}

func (a *EntityData) SetViewMag(mag float64) {
	a.viewMag = mag
}

func (a *EntityData) GetAge() int {
	return a.age
}

func (a *EntityData) GetDateOfLastBirth() int {
	return a.birthNotBefore
}
func (a *EntityData) SetDateOfLastBirth(d int) {
	a.birthNotBefore = d
}
func (a *EntityData) GetMatureAge() int {
	return a.matureAge
}

func (a *EntityData) GetWorld() *world.World {
	return a.w
}

// Die neue Position e.pos aus e.vel und e.acc bestimmen und die Lebensenergie aktualisieren
/*
func (a *data) Update(others *[]Entity) (offSpring *data) {
	a.age++
	if a.moveable {
		a.energy -= a.energyLoss
		a.randomStep()
		if a.preys != nil {
			a.searchFood(a.preys)
			a.eatFood(a.preys)
		}
		a.avoidCollisionWithSeenObjects(others)
		a.repelFromWater()

		a.applyMove(others)  // <== ACHTUNG: hier ist jetzt alles für die Bewegung drin
		offSpring = a.GetOffspring()
	}
	return
}
*/

/*
func (a *data) GetOffspring() *data {
	if a.age > a.dateOfLastBirth {
		a.dateOfLastBirth += rand.Intn(a.matureAge)
		//fmt.Println("Neues Geburtsdatum:", a.dateOfLastBirth)
		return New(a.w)
	}
	return nil

}
*/

// Vor.:
// Eff.: Die neue Position Entity.pos ist bestimmt. Das Überlappen von
// Objekten wird vermieden
// Erg.:
func (a *EntityData) ApplyMove(others *[]Entity, preys *[]Entity) {

	a.randomStep()
	if preys != nil {
		a.searchFood(preys)
		a.eatFood(preys)
	}
	a.avoidCollisionWithSeenObjects(others)
	a.repelFromWater()

	// Überlappen von Ojekten gleicher Art vermeiden
	newPos := a.pos.Add(a.vel)
	collision := false
	sumDiff := vec{0, 0}
	var counts float64
	for _, other := range *others {
		dist := newPos.Sub(other.GetPosition())
		if !other.IsSame(a) && dist.Magnitude() <= float64(a.imgHeight*0.5) {
			collision = collision || true
			// Einfaches Separieren
			sumDiff = sumDiff.Add(dist.Unit())
			counts++
		}
	}
	if counts > 0 {
		sumDiff = sumDiff.Scale(1 / counts)
		sumDiff = sumDiff.Scale(0.5)
		a.vel = a.vel.Add(sumDiff)
	}
	if !collision {
		a.pos = newPos
	} else {
		a.pos = a.pos.Add(a.vel)
	}

}

func (a *EntityData) avoidCollisionWithSeenObjects(others *[]Entity) {
	avg := vec{0, 0}
	_, dirs := a.SeeOthers(others)
	for _, dir := range *dirs {
		dir = dir.Unit()
		dir = dir.Scale(a.viewMag - dir.Magnitude())
		avg = avg.Add(dir)
	}
	if len(*dirs) > 0 {
		avg.Scale(1 / float64(len(*dirs)))
	}
	z := a.vel.Unit().Scale(a.ahead)
	z = z.Add(avg.Unit().Scale(-0.25)) // <== Ändert die Stärke, mit welcher die Objekte voneinander wegstreben
	a.vel = a.vel.Unit().Scale(a.maxVel).Add(z)
}

//func (a *EntityData) flee(others *[]Entity) {}

// Vor.:
// Eff.: Bewegt sich in Richtung des nächsten im Sichtfeld gelegenen Tiers/Essens
// Erg.:

func (a *EntityData) searchFood(others *[]Entity) {
	if others == nil || a.health > 500 {
		return
	}

	seenEntities, directions := a.SeeOthers(others)
	iClosest := 0
	for i := 0; i < len(*seenEntities); i++ {
		if (*directions)[i].Magnitude() < (*directions)[iClosest].Magnitude() {
			iClosest = i
		}
	}

	if len(*seenEntities) > 0 {
		huntDir := (*directions)[iClosest].Unit()
		huntDir = huntDir.Scale(a.maxAccPhi * 10)
		z := a.vel.Unit().Scale(a.ahead)
		z = z.Add(huntDir)
		a.vel = a.vel.Unit().Scale(a.maxVel).Add(z)
	}
}

// BUG: Irgendetwas stimmt hier nicht ...
func (a *EntityData) eatFood(others *[]Entity) {
	if a.health > 500 {
		return
	}

	newPos := a.pos.Add(a.vel)
	for _, other := range *others {
		dist := newPos.Sub(other.GetPosition())
		if !other.IsSame(a) && dist.Magnitude() <= float64(a.imgHeight*1.0) {
			other.SetHealth(0) //wurde gegessen
			//fmt.Println(" >> Einer hat gegessen: ", i)
			a.health += other.GetHealthWhenEaten()
		}
	}
}

// Vor.:
// Eff.: Addiert in der Nähe vom Wasser ein Geschwindigkeitskomponente vom Wasser weg
// auf die aktuelle Geschwindigkeit.
// Erg.:
func (a *EntityData) repelFromWater() {
	// Wenn das Objekt in die Nähe des Bildschirmrandes kommt,
	// wird es senkrecht dazu beschleunigt (dreht also um)
	// TODO: Die Beschleunigung vom Rand weg sollte Proportional zur Entfernung zum Rand sein.

	n, no, o, so, s, sw, w, nw := (*a.w).GetTileBorders(int(a.pos[0]), int(a.pos[1]))
	tileX, tileY := (*a.w).GetXYTile(int(a.pos.X()), int(a.pos.Y()))

	repel := vec{0, 0}
	const d = 15

	if n && a.pos.Y() <= float64((*a.w).GetTileSizeScaled()*tileY+d) {
		repel[1] = a.accBorder
	}

	if no && a.pos.Y() <= float64((*a.w).GetTileSizeScaled()*tileY+d) && float64(a.pos.X()) >= float64((*a.w).GetTileSizeScaled()*(tileX+1)-d) {
		repel[1] = a.accBorder
		repel[0] = -a.accBorder
	}

	if o && a.pos.X() >= float64((*a.w).GetTileSizeScaled()*(tileX+1)-d) {
		repel[0] = -a.accBorder
	}

	if so && a.pos.X() >= float64((*a.w).GetTileSizeScaled()*(tileX+1)-d) && float64(a.pos.Y()) >= float64((*a.w).GetTileSizeScaled()*(tileY+1)-d) {
		repel[0] = -a.accBorder
		repel[1] = -a.accBorder
	}

	if s && a.pos.Y() >= float64((*a.w).GetTileSizeScaled()*(tileY+1)-d) {
		repel[1] = -a.accBorder
	}

	if sw && a.pos.Y() >= float64((*a.w).GetTileSizeScaled()*(tileY+1)-d) && float64(a.pos.X()) <= float64((*a.w).GetTileSizeScaled()*tileX+d) {
		repel[1] = -a.accBorder
		repel[0] = a.accBorder
	}

	if w && a.pos.X() <= float64((*a.w).GetTileSizeScaled()*tileX+d) {
		repel[0] = a.accBorder
	}

	if nw && a.pos.X() <= float64((*a.w).GetTileSizeScaled()*tileX+d) && float64(a.pos.Y()) <= float64((*a.w).GetTileSizeScaled()*tileY+d) {
		repel[0] = a.accBorder
		repel[1] = a.accBorder
	}

	a.vel = a.vel.Add(repel)
}

// Vor.: ?
// Eff.: ?
// Erg.: Splice mit Objekten (seen) und deren Abstandsvektoren (direction),
// die im Sichtfeld des Objekts liegen
func (a *EntityData) SeeOthers(others *[]Entity) (*[]Entity, *[]vec) {
	inView := false
	var seen []Entity
	var direction []vec

	for _, other := range *others {
		delta := other.GetPosition().Sub(a.pos)
		if !other.IsSame(a) && delta.Magnitude() < a.viewMag {

			if math.Abs(a.vel.Angle(delta)) < a.viewAngle {
				inView = inView || true
				seen = append(seen, other)
				direction = append(direction, delta.Clone())
			}
		} else {
			inView = inView || false
		}
	}
	a.inView = inView
	return &seen, &direction
}

func (a *EntityData) Draw(screen *ebiten.Image) {
	a.drawEntity(screen)
}

func (a *EntityData) drawStats(screen *ebiten.Image) {

	if a.moveable {
		optFont := &text.GoTextFace{
			Source: a.font,
			Size:   13,
		}

		opTxt := &text.DrawOptions{}
		livePoints := 100 - float32(a.age)/float32(a.lifeSpan)*100

		opTxt.GeoM.Translate(-float64(a.imgWidth*1.5), -float64(a.imgHeight*1.5)) // Koordinaten zuerst in die Mitte des Bilder bewegen ...
		opTxt.GeoM.Translate(a.pos[0], a.pos[1])
		if livePoints < 10 || a.health < 30 {
			opTxt.ColorScale.ScaleWithColor(color.RGBA{255, 0, 0, 255})
		} else if livePoints > 90 {
			//opTxt.ColorScale.ScaleWithColor(color.RGBA{0, 255, 0, 255})
		} else {
			opTxt.ColorScale.ScaleWithColor(color.Black) // Rot
		}
		text.Draw(screen, fmt.Sprintf("H %2.0f", a.health), optFont, opTxt)

		opTxt.GeoM.Translate(0, 43)
		//opTxt.ColorScale.ScaleWithColor(color.Black) // Rot

		text.Draw(screen, fmt.Sprintf("%2.0f %%", livePoints), optFont, opTxt)
	}
}

func (a *EntityData) drawView() *ebiten.DrawImageOptions {
	w := float32(a.imgDebug.Bounds().Dx())
	h := float32(a.imgDebug.Bounds().Dy())
	a.imgDebug.Clear()

	// --- Rahmen ---
	/*
		vector.StrokeRect(a.imgDebug, 1, 1, w-1, h-1, 1, color.Gray{100}, false)
		opRect := &ebiten.DrawImageOptions{}
		opRect.GeoM.Translate(a.pos[0]+float64(-w/2), a.pos[1]-float64(h/2)) // ... und zum Schluss ein die gewünschte Stelle bewegen
		screen.DrawImage(a.imgDebug, opRect)
		a.imgDebug.Clear()
	*/

	opD := &ebiten.DrawImageOptions{}

	// --- Geschwindigkeit ---
	length := float32(a.vel.Magnitude()) * 20
	vector.StrokeLine(a.imgDebug, w/2, h/2, w/2, h/2-length, 1.5, color.Gray{0}, false)
	dirAngle := a.vel.Angle(vec{0, -1}) // Winkel zur Y-Achse bestimmen

	// --- Sichtfeld ---
	viewColor := color.NRGBA{120, 180, 100, 80}
	if a.inView {
		viewColor = color.NRGBA{150, 100, 180, 80}
	}
	a.makeArc(a.imgDebug, float32(a.viewMag), float32(-math.Pi/2-a.viewAngle), float32(-math.Pi/2+a.viewAngle), viewColor, false)
	// --- Transformation ---
	opD.GeoM.Translate(float64(-w/2), -float64(h/2)) // Koordinaten zuerst in die Mitte des Bilder bewegen ...
	opD.GeoM.Rotate(-dirAngle)
	// ... dann drehen ...
	opD.GeoM.Translate(a.pos[0], a.pos[1]) // ... und zum Schluss ein die gewünschte Stelle bewegen
	return opD
}

/*
func (a *EntityData) drawDebug(screen *ebiten.Image) {
	a.imgDebug.Clear()
	w := float32(a.imgDebug.Bounds().Dx())
	h := float32(a.imgDebug.Bounds().Dy())

	if a.moveable {
	} else {
		// --- Transformation ---
		opD.GeoM.Translate(float64(-w/2), -float64(h/2)) // Koordinaten zuerst in die Mitte des Bilder bewegen ...
		opD.GeoM.Translate(a.pos[0], a.pos[1])           // ... und zum Schluss ein die gewünschte Stelle bewegen

	}

	screen.DrawImage(a.imgDebug, opD)
}
*/

// Vor.: ?
// Eff.: Zeichnet ein Tier als Vektorgrafik mit Sichtfeld und Geschwindigkeitsvektor
// Erg.: ?
func (a *EntityData) drawEntity(screen *ebiten.Image) {
	if a.moveable {
		if (*a.w).GetDebug() {
			opD := a.drawView()
			screen.DrawImage(a.imgDebug, opD)
		}
		if (*a.w).GetShowStats() {
			a.drawStats(screen)
		}
	}
	// Das zu zeichnende Bild ist in e.img gespeichert.
	op := &ebiten.DrawImageOptions{}
	if !a.imgIsSetFromFile {
		op.GeoM.Translate(-float64(a.imgHeight)/2, -float64(a.imgHeight)/2) // Koordinaten zuerst in die Mitte des Bilder bewegen ...
	} else {
		op.GeoM.Translate(-float64(a.imgHeight), -float64(a.imgHeight)) // Koordinaten zuerst in die Mitte des Bilder bewegen ...
	}
	op.GeoM.Translate(a.pos[0], a.pos[1]) // ... und zum Schluss ein die gewünschte Stelle bewegen
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(a.img, op)
}

// Vor.: ?
// Eff.: Ein Kreisbogen wird erzeugt, dessen Mittelpunkt mittig in `img` gespeichert wird.
// Erg.:
func (a *EntityData) makeArc(img *ebiten.Image, radius float32, startAngle, endAngle float32, c color.NRGBA, line bool) {
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
