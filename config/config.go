package config

import "math"

const (
	CatStartNumber	 = 5
	CatHealthLoss    = 300 // stirbt frühestens nach CatHealthLoss Frames ohne Nahrung
	CatMaxVelocitiy  = 0.5
	CatViewMagnitude = 200
	CatViewAngle     = math.Pi * 0.5
	CatMatureAge     = 130

	FoxStartNumber	 = 5
	FoxHealthLoss    = 200 // stirbt frühestens nach FoxHealthLoss Frames ohne Nahrung
	FoxMaxVelocitiy  = 0.55
	FoxViewMagnitude = 300
	FoxViewAngle     = math.Pi * 0.12
	FoxMatureAge     = 130

	BunnyStartNumber   = 10
	BunnyHealthLoss    = 240 // stirbt frühestens nach BunnyHealthLoss Frames ohne Nahrung
	BunnyMaxVelocitiy  = 0.5
	BunnyViewMagnitude = 80
	BunnyViewAngle     = math.Pi * 0.9
	BunnyMatureAge     = 100
	
	GrassStartNumber = 20
)
