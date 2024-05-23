package config

import "math"

const (
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
