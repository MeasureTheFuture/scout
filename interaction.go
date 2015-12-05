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
	"log"
	"math"
	"time"
)

type Waypoint struct {
	XPixels          int // x-coordinate of waypoint centroid in pixels
	YPixels          int // y-coordinate of waypoint centroid in pixels
	HalfWidthPixels  int // Half the width of the waypoint in pixels
	HalfHeightPixels int // Half the height of the waypoint in pixels
	T                float32
}

// distanceSq calculates the distance squared between this and the
// supplied waypoint.
func (a Waypoint) distanceSq(b Waypoint) int {
	dx := a.XPixels - b.XPixels
	dy := a.YPixels - b.YPixels

	return (dx * dx) + (dy * dy)
}

type Interaction struct {
	Entered  time.Time  // The time the interaction started (rounded to nearest half hour)
	Duration float32    // The total duration of the interaction.
	Path     []Waypoint // The pathway of the interaction through the scene.
}

// lastWaypoint returns the last waypoint within the interaction.
func (i Interaction) lastWaypoint() Waypoint {
	return i.Path[len(i.Path)-1]
}

type Scene struct {
	Interactions []Interaction // The current interactions occuring within the scene.
}

// initScene creates an empty scene that can be used for monitoring interactions.
func initScene() Scene {
	return Scene{}
}

func buildDistanceMap(s *Scene, detected []Waypoint) map[int][][]int {
	var distances map[int][][]int = make(map[int][][]int)

	// For each of the detected waypoints, work out the
	// closest interaction in the scene.
	for i := 0; i < len(detected); i++ {
		dist := math.MaxInt32
		minW := -1

		for j := 0; j < len(s.Interactions); j++ {
			d := detected[i].distanceSq(s.Interactions[j].lastWaypoint())
			if d < dist {
				dist = d
				minW = j
			}
		}

		distances[minW] = append(distances[minW], []int{dist, i})
	}

	return distances
}

// addInteraction
func addInteraction(s *Scene, detected []Waypoint) {
	// The start time to use for the new interaction -- truncated to the nearest 30 minutes.
	start := time.Now().Truncate(30 * time.Minute)
	//log.Printf("\t detected: " + strconv.Itoa(len(detected)))

	if len(s.Interactions) == 0 {
		// Empty scene: just add a new interaction for each new waypoint.
		for i := 0; i < len(detected); i++ {
			s.Interactions = append(s.Interactions, Interaction{start, 0.0, []Waypoint{detected[i]}})
		}

	} else {
		// Existing scene:
		// for each of the detected waypoints
		// 	 work which of the existing interactions are closest.
		//   for interactions that have more than one close detected waypoints
		//		create a new interaction from the furthest detected waypoint
		// 		the nearest waypoint is used to update the interaction.
		distances := buildDistanceMap(s, detected)

		for i := 0; i < len(distances); i++ {
			if len(distances[i]) > 1 {
				dist := math.MaxInt32
				minW := -1

				for j := 0; j < len(distances[i]); j++ {
					if distances[i][j][0] < dist {
						dist = distances[i][j][0]
						minW = j
					}
				}

				for j := 0; j < len(distances[i]); j++ {
					if j != minW {
						s.Interactions = append(s.Interactions, Interaction{start, 0.0, []Waypoint{detected[distances[i][j][1]]}})
					} else {
						s.Interactions[i].Path = append(s.Interactions[i].Path, detected[distances[i][j][1]])
					}
				}
			} else if len(distances[i]) == 1 {
				s.Interactions[i].Path = append(s.Interactions[i].Path, detected[distances[i][0][1]])
			}
		}
	}
}

func removeInteraction(s *Scene, detected []Waypoint) {
	distances := buildDistanceMap(s, detected)
	matched := map[int]bool{}

	for i := 0; i < len(distances); i++ {
		for j := 0; j < len(distances[i]); j++ {
			matched[distances[i][j][1]] = true
		}
	}

	for i := 0; i < len(s.Interactions); i++ {
		if _, ok := matched[i]; ok {
			s.Interactions[i].Path = append(s.Interactions[i].Path, detected[distances[i][0][1]])
		} else {
			sendInteraction(s.Interactions[i])
			s.Interactions = append(s.Interactions[:i], s.Interactions[i+1:]...)
		}
	}
}

func monitorScene(s *Scene, detected []Waypoint) {
	if len(detected) >= len(s.Interactions) {
		addInteraction(s, detected)
	} else {
		removeInteraction(s, detected)
	}

	log.Printf("")
}

func sendInteraction(i Interaction) {
	// TODO: Broadcast the interaction to the mothership.
}
