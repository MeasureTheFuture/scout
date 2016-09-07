/*
 * Copyright (C) 2015 Clinton Freeman
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"time"
)

type Waypoint struct {
	XPixels          int     // x-coordinate of waypoint centroid in pixels
	YPixels          int     // y-coordinate of waypoint centroid in pixels
	HalfWidthPixels  int     // Half the width of the waypoint in pixels
	HalfHeightPixels int     // Half the height of the waypoint in pixels
	T                float32 // The number of seconds elapsed since the beginning of the interaction
}

// distanceSq calculates the distance squared between this and the
// supplied waypoint.
func (a Waypoint) distanceSq(b Waypoint) int {
	dx := a.XPixels - b.XPixels
	dy := a.YPixels - b.YPixels

	return (dx * dx) + (dy * dy)
}

// perpendicularDistance calulates the distance from a point (x) to a line
// (defined by a and b).
func (x Waypoint) perpendicularDistance(a Waypoint, b Waypoint) float64 {
	n := float64(((b.YPixels - a.YPixels) * x.XPixels) - ((b.XPixels - a.XPixels) * x.YPixels) + (b.XPixels * a.YPixels) - (b.YPixels * a.XPixels))
	d := float64(((b.YPixels - a.YPixels) * (b.YPixels - a.YPixels)) + ((b.XPixels - a.XPixels) * (b.XPixels - a.XPixels)))

	return (math.Abs(n) / math.Sqrt(d))
}

// compare returns true if two waypoints are the same, false otherwise.
func (a Waypoint) Equal(b Waypoint) bool {
	return a.XPixels == b.XPixels &&
		a.YPixels == b.YPixels &&
		a.HalfHeightPixels == b.HalfHeightPixels &&
		a.HalfWidthPixels == b.HalfWidthPixels &&
		math.Abs(float64(a.T-b.T)) < 0.007
}

type Interaction struct {
	UUID     string     // The UUID for the scout that detected the interaction.
	Version  string     // The Version of the protocol used for transmitting data to the mothership
	Entered  time.Time  // The time the interaction started (rounded to nearest half hour)
	started  time.Time  // The actual time the interaction started. Private. Not to be transmitted for privacy concerns
	Duration float32    // The total duration of the interaction.
	Path     []Waypoint // The pathway of the interaction through the scene.
}

func (i Interaction) Equal(wp []Waypoint) bool {
	if len(i.Path) != len(wp) {
		return false
	}

	for k, v := range i.Path {
		if !v.Equal(wp[k]) {
			return false
		}
	}

	return true
}

func NewInteraction(w Waypoint, config Configuration) Interaction {
	start := time.Now()

	// The start time broadcasted for the interaction is truncated to the nearest 30 minutes.
	apparentStart := start.Round(15 * time.Minute)

	i := Interaction{config.UUID, "0.1", apparentStart, start, 0.0, []Waypoint{}}
	i.addWaypoint(w)
	return i
}

// addWaypoint inserts a new waypoint to the end of the interaction.
func (i *Interaction) addWaypoint(w Waypoint) {
	newW := w
	newW.T = float32(time.Now().Sub(i.started).Seconds())

	i.Duration = newW.T
	i.Path = append(i.Path, newW)
}

func douglasPeucker(path []Waypoint, epsilon float64) []Waypoint {
	if len(path) == 1 {
		return path
	}

	dMax := 0.0
	iMax := 0
	end := len(path) - 1

	for i := 1; i < end; i++ {
		d := path[i].perpendicularDistance(path[0], path[end])
		if d > dMax {
			iMax = i
			dMax = d
		}
	}

	if dMax > epsilon {
		a := douglasPeucker(path[0:iMax+1], epsilon)
		b := douglasPeucker(path[iMax:len(path)], epsilon)

		if len(b) > 1 {
			return append(a, b[1:len(b)]...)
		} else {
			return a
		}
	}

	return []Waypoint{path[0], path[end]}
}

// lastWaypoint returns the last waypoint within the interaction.
func (i *Interaction) lastWaypoint() Waypoint {
	return i.Path[len(i.Path)-1]
}

func (i *Interaction) simplify(config Configuration) {
	i.Path = douglasPeucker(i.Path, config.SimplifyEpsilon)
}

func (i *Interaction) post(debug bool, config Configuration) {
	i.simplify(config) // Remove unessary segments from the pathway before sending.

	body := bytes.Buffer{}
	encoder := json.NewEncoder(&body)

	err := encoder.Encode(i)
	if err != nil {
		log.Printf("ERROR: Unable to encode interaction for transport to mothership")
	}

	if debug {
		b, _ := json.Marshal(i)
		filename := string("f" + fmt.Sprintf("%09d", i.started.Unix()) + "-metadata.json")
		ioutil.WriteFile(filename, b, 0611)
	}

	post("interaction.json", config.MothershipAddress+"/scout_api/interaction", config.UUID, &body)
}

type Scene struct {
	Interactions []Interaction // The current interactions occuring within the scene.
}

// initScene creates an empty scene that can be used for monitoring interactions.
func initScene() *Scene {
	return &Scene{}
}

func (s *Scene) buildDistanceMap(detected []Waypoint) map[int][]int {
	var distances map[int][]int = make(map[int][]int)

	// For each of the detected waypoints, work out the
	// closest interaction in the scene.
	for i := 0; i < len(detected); i++ {
		dist := math.MaxInt32
		closestInteraction := -1

		for j := 0; j < len(s.Interactions); j++ {
			d := detected[i].distanceSq(s.Interactions[j].lastWaypoint())
			if d < dist {
				dist = d
				closestInteraction = j
			}
		}

		distances[i] = []int{dist, closestInteraction}
	}

	return distances
}

// addInteraction
func (s *Scene) addInteraction(detected []Waypoint, config Configuration) {

	if len(s.Interactions) == 0 {
		// Empty scene: just add a new interaction for each new waypoint.
		for i := 0; i < len(detected); i++ {
			s.Interactions = append(s.Interactions, NewInteraction(detected[i], config))
		}

	} else {
		// Existing scene:
		// for each of the detected waypoints
		// 	 work which of the existing interactions are closest.
		//   for interactions that have more than one close detected waypoints
		//		create a new interaction from the furthest detected waypoint
		// 		the nearest waypoint is used to update the existing interaction.
		distances := s.buildDistanceMap(detected)
		// assert(len(distances) == len(detected))

		for i := 0; i < len(distances); i++ {
			dist := math.MaxInt32
			closestI := -1

			// Work out if this detected waypoint is the closest one to an existing interaction.
			for j := 0; j < len(distances); j++ {
				if distances[i][1] == distances[j][1] && distances[j][0] < dist {
					dist = distances[j][0]
					closestI = j
				}
			}

			if i == closestI {
				// If this detected element is the closest to an interaction - update the interaction with the
				// detected waypoint.
				s.Interactions[distances[i][1]].addWaypoint(detected[i])
			} else {
				// Otherwise this must be a new interaction, create it and add it to the scene.
				s.Interactions = append(s.Interactions, NewInteraction(detected[i], config))
			}
		}
	}
}

func (s *Scene) removeInteraction(detected []Waypoint, debug bool, config Configuration) {
	distances := s.buildDistanceMap(detected)
	matched := map[int]int{}

	for i := 0; i < len(distances); i++ {
		matched[distances[i][1]] = i
	}

	for i := len(s.Interactions) - 1; i >= 0; i-- {
		if v, ok := matched[i]; ok {
			s.Interactions[i].addWaypoint(detected[v])
		} else {
			// Only transmit the interaction to the mothership if it is longer than the
			// specified minimum duration. This is to filter out any detected noise.
			if s.Interactions[i].Duration > config.MinDuration {
				s.Interactions[i].post(debug, config)
			}
			s.Interactions = append(s.Interactions[:i], s.Interactions[i+1:]...)
		}
	}
}

func (s *Scene) update(detected []Waypoint, debug bool, config Configuration) {
	if len(detected) >= len(s.Interactions) {
		s.addInteraction(detected, config)
	} else {
		s.removeInteraction(detected, debug, config)
	}
}

func (s *Scene) save(filename string) {
	b, _ := json.Marshal(s)
	ioutil.WriteFile(filename, b, 0611)
}
