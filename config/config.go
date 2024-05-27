package config

import "math"

const (
	CatHealthLoss    = 300 // stirbt frühestens nach CatHealthLoss Frames ohne Nahrung
	CatMaxVelocitiy  = 0.5
	CatViewMagnitude = 200
	CatViewAngle     = math.Pi * 0.5
	CatMatureAge     = 130

	FoxHealthLoss    = 200 // stirbt frühestens nach FoxHealthLoss Frames ohne Nahrung
	FoxMaxVelocitiy  = 0.55
	FoxViewMagnitude = 300
	FoxViewAngle     = math.Pi * 0.12
	FoxMatureAge     = 130

	BunnyHealthLoss    = 240 // stirbt frühestens nach FoxHealthLoss Frames ohne Nahrung
	BunnyMaxVelocitiy  = 0.5
	BunnyViewMagnitude = 80
	BunnyViewAngle     = math.Pi * 0.9
	BunnyMatureAge     = 100
)
