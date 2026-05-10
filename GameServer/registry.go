package main

import (
	"sync"
)

type SpatialRegistry struct {
	mu   sync.RWMutex
	X    []float32
	Y    []float32
	Z    []float32
	VX   []float32
	VY   []float32
	VZ   []float32
	AX   []float32
	AY   []float32
	AZ   []float32
	Mass []float32
}

func NewSpatialRegistry() *SpatialRegistry {
	return &SpatialRegistry{
		X:    make([]float32, 0),
		Y:    make([]float32, 0),
		Z:    make([]float32, 0),
		VX:   make([]float32, 0),
		VY:   make([]float32, 0),
		VZ:   make([]float32, 0),
		AX:   make([]float32, 0),
		AY:   make([]float32, 0),
		AZ:   make([]float32, 0),
		Mass: make([]float32, 0),
	}
}

func (sr *SpatialRegistry) AddEntity(x, y, z, mass float32) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.X = append(sr.X, x)
	sr.Y = append(sr.Y, y)
	sr.Z = append(sr.Z, z)
	sr.VX = append(sr.VX, 0)
	sr.VY = append(sr.VY, 0)
	sr.VZ = append(sr.VZ, 0)
	sr.AX = append(sr.AX, 0)
	sr.AY = append(sr.AY, 0)
	sr.AZ = append(sr.AZ, -9.8)
	sr.Mass = append(sr.Mass, mass)
}

type PotentialFieldRegistry struct {
	mu        sync.RWMutex
	X         []float32
	Y         []float32
	Z         []float32
	Amplitude []float32
	Sigma     []float32
}

func NewPotentialFieldRegistry() *PotentialFieldRegistry {
	return &PotentialFieldRegistry{
		X:         make([]float32, 0),
		Y:         make([]float32, 0),
		Z:         make([]float32, 0),
		Amplitude: make([]float32, 0),
		Sigma:     make([]float32, 0),
	}
}

func (pfr *PotentialFieldRegistry) AddEmitters(x, y, z, amplitude, sigma float32) {
	pfr.mu.Lock()
	defer pfr.mu.Unlock()
	pfr.X = append(pfr.X, x)
	pfr.Y = append(pfr.Y, y)
	pfr.Z = append(pfr.Z, z)
	pfr.Amplitude = append(pfr.Amplitude, amplitude)
	pfr.Sigma = append(pfr.Sigma, sigma)
}
