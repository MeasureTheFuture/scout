/*
 * Copyright (C) 2016 Clinton Freeman
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
	"encoding/json"
	"io/ioutil"
	"math"
)

type Scene struct {
	Interactions     []Interaction  // The current interactions occuring within the scene.
	idleInteractions []*Interaction // The current interactions that are idle. We wait
	// a user specified time.
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
				// Detect if this belongs to a resumable interaction.
				// check time and distance.

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
			// // Only transmit the interaction to the mothership if it is longer than the
			// // specified minimum duration. This is to filter out any detected noise.
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
