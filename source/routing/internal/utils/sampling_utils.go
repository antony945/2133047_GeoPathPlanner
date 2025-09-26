package utils

import (
	"geopathplanner/routing/internal/models"
	"math/rand/v2"

	"github.com/paulmach/orb"
)

// Sampler defines the contract for any sampling strategy.
type Sampler interface {
	SampleXY(minX, maxX, minY, maxY float64) (float64, float64)
	SampleXYZ(minX, maxX, minY, maxY, minZ, maxZ float64) (float64, float64, float64)
	SampleZ(minZ, maxZ float64) float64
}

// ------------------------------------------------------

// TODO: Test uniform sampler
type UniformSampler struct {}

func NewUniformSampler() *UniformSampler {
	return &UniformSampler{}
}

func (s *UniformSampler) SampleXY(minX, maxX, minY, maxY float64) (float64, float64) {
	x := s.Sample(minX, maxX)
	y := s.Sample(minY, maxY)
	return x, y
}

func (s *UniformSampler) SampleXYZ(minX, maxX, minY, maxY, minZ, maxZ float64) (float64, float64, float64) {
	x := s.Sample(minX, maxX)
	y := s.Sample(minY, maxY)
	z := s.SampleZ(minZ, maxZ)
	return x, y, z
}

func (s *UniformSampler) SampleZ(minZ, maxZ float64) float64 {
	return s.Sample(minZ, maxZ)
}

func (s *UniformSampler) Sample(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// --------------------------------------------------------

// TODO: Test halton sampler
type HaltonSampler struct {
	idx int
}

func NewHaltonSampler() *HaltonSampler {
	return &HaltonSampler{
		idx: 1,
	}
}

func (s *HaltonSampler) SampleXY(minX, maxX, minY, maxY float64) (float64, float64) {
	x := minX + s._halton(s.idx, 2)*(maxX-minX)
	y := minY + s._halton(s.idx, 3)*(maxY-minY)
	s.idx++
	return x, y
}

func (s *HaltonSampler) SampleXYZ(minX, maxX, minY, maxY, minZ, maxZ float64) (float64, float64, float64) {
	x := minX + s._halton(s.idx, 2)*(maxX-minX)
	y := minY + s._halton(s.idx, 3)*(maxY-minY)
	z := minZ + s._halton(s.idx, 5)*(maxZ-minZ)
	s.idx++
	return x, y, z
}

func (s *HaltonSampler) SampleZ(minZ, maxZ float64) float64 {
	// Don't sample again, assume that a SampleXY was already called before so decrement idx before sampling
	z := minZ + s._halton(s.idx-1, 5)*(maxZ-minZ)
	return z
}

func (s *HaltonSampler) _halton(idx, base int) float64 {
	result := 0.0
	f := 1.0
	i := idx
	for i > 0 {
		f /= float64(base)
		result += f * float64(i%base)
		i /= base
	}
	return result
}

// -------------------------------------------------------------------------------------------

// TODO: Test goal bias sampler
type GoalBiasSampler struct {
	InternalSampler Sampler
	Goal *models.Waypoint
	Bias float64
	last_chosen_goal bool
}

func NewGoalBiasSampler(sampler Sampler, goal *models.Waypoint, bias float64) *GoalBiasSampler {
	return &GoalBiasSampler{
		InternalSampler: sampler,
		Goal: goal,
		Bias: bias,
		last_chosen_goal: false,
	}
}

func (s *GoalBiasSampler) UseGoal() bool {
	s.last_chosen_goal = rand.Float64() < s.Bias
	return s.last_chosen_goal
}

func (s *GoalBiasSampler) SampleXY(minX, maxX, minY, maxY float64) (float64, float64) {
	if s.UseGoal() {
		return s.Goal.Lon, s.Goal.Lat
	}
	return s.InternalSampler.SampleXY(minX, maxX, minY, maxY)
}

func (s *GoalBiasSampler) SampleXYZ(minX, maxX, minY, maxY, minZ, maxZ float64) (float64, float64, float64) {
	if s.UseGoal() {
		return s.Goal.Lon, s.Goal.Lat, s.Goal.Alt.Value
	}
	return s.InternalSampler.SampleXYZ(minX, maxX, minY, maxY, minZ, maxZ)
}

func (s *GoalBiasSampler) SampleZ(minZ, maxZ float64) (float64) {
	// Here do not check again if we have to use goal or no, we assume the choice was already been done
	if s.last_chosen_goal {
		return s.Goal.Alt.Value
	}
	return s.InternalSampler.SampleZ(minZ, maxZ)
}

// -------------------------------------------------------------------------------------------

func Sample2D(geometry orb.Geometry, sampler Sampler) orb.Point {
	// 1. Retrieve bounding box to sample there
	bound := geometry.Bound()
	minLon, minLat := bound.Min.Lon(), bound.Min.Lat()
	maxLon, maxLat := bound.Max.Lon(), bound.Max.Lat()

	// Track if the sampled point is actually valid (if it's inside the geometry or not)
	var sampled orb.Point
	isValid := false
	for !isValid {
		// 2. Sample point
		randLon, randLat := sampler.SampleXY(minLon, maxLon, minLat, maxLat)
				
		// 3. Check if sampled point is inside the geometry (because maybe it's inside the bbox but not the geometry)
		sampled = orb.Point{randLon, randLat}
		isValid = PointInPolygon2D(sampled, geometry.(orb.Polygon))
	}

	return sampled
}

func SampleWithAltitude2D(geometry orb.Geometry, alt models.Altitude, sampler Sampler) *models.Waypoint {
	sampled := Sample2D(geometry, sampler)
	wp, _ := models.NewWaypoint(sampled.Lat(), sampled.Lon(), alt)
	return wp
}

func Sample3D(geometry *models.Feature3D, sampler Sampler) *models.Waypoint {
	sampled := Sample2D(geometry.Geometry, sampler)
	minAlt, maxAlt := geometry.MinAltitude.Normalize().Value, geometry.MaxAltitude.Normalize().Value
	
	// Sample point
	randAlt, _ := models.NewAltitude(sampler.SampleZ(minAlt, maxAlt), models.MT)		
	sampled3D, _ := models.NewWaypoint(sampled.Lat(), sampled.Lon(), randAlt)
	return sampled3D
}