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
	//"log"
	"encoding/json"
	"io/ioutil"
	"math"
	//"strconv"
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

func (a Waypoint) compare(b Waypoint) bool {
	return a.XPixels == b.XPixels && a.YPixels == b.YPixels && a.HalfHeightPixels == b.HalfHeightPixels && a.HalfWidthPixels == b.HalfWidthPixels && math.Abs(float64(a.T-b.T)) < 0.007
}

type Interaction struct {
	Entered  time.Time  // The time the interaction started (rounded to nearest half hour)
	started  time.Time  // The actual time the interaction started. Private. Not to be transmitted for privacy concerns
	Duration float32    // The total duration of the interaction.
	Path     []Waypoint // The pathway of the interaction through the scene.
}

func (i *Interaction) addWaypoint(w Waypoint) {
	newW := w
	newW.T = float32(time.Now().Sub(i.started).Seconds())

	i.Path = append(i.Path, newW)
}

// lastWaypoint returns the last waypoint within the interaction.
func (i *Interaction) lastWaypoint() Waypoint {
	return i.Path[len(i.Path)-1]
}

type Scene struct {
	Interactions []Interaction // The current interactions occuring within the scene.
}

// initScene creates an empty scene that can be used for monitoring interactions.
func initScene() Scene {
	return Scene{}
}

func buildDistanceMap(s *Scene, detected []Waypoint) map[int][]int {
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
func addInteraction(s *Scene, detected []Waypoint) {
	// The start time to use for the new interaction -- truncated to the nearest 30 minutes.
	start := time.Now().Truncate(30 * time.Minute)

	if len(s.Interactions) == 0 {
		// Empty scene: just add a new interaction for each new waypoint.
		for i := 0; i < len(detected); i++ {
			s.Interactions = append(s.Interactions, Interaction{start, time.Now(), 0.0, []Waypoint{detected[i]}})
		}

	} else {
		// Existing scene:
		// for each of the detected waypoints
		// 	 work which of the existing interactions are closest.
		//   for interactions that have more than one close detected waypoints
		//		create a new interaction from the furthest detected waypoint
		// 		the nearest waypoint is used to update the interaction.
		distances := buildDistanceMap(s, detected)
		// assert(len(distances) == len(detected))

		for i := 0; i < len(distances); i++ {
			//log.Printf("\t len(detected) = " + strconv.Itoa(len(detected)))
			//log.Printf("\t len(distances[" + strconv.Itoa(i) + "]) = " + strconv.Itoa(len(distances[i])))

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
				s.Interactions = append(s.Interactions, Interaction{start, time.Now(), 0.0, []Waypoint{detected[i]}})
			}

			//log.Printf("\t closestI[" + strconv.Itoa(closestI) + "], closestD[" + strconv.Itoa(closestD) + "] = " + strconv.Itoa(dist))
		}
	}
}

func removeInteraction(s *Scene, detected []Waypoint) {
	distances := buildDistanceMap(s, detected)
	matched := map[int]int{}

	for i := 0; i < len(distances); i++ {
		matched[distances[i][1]] = i
	}

	for i := len(s.Interactions) - 1; i >= 0; i-- {
		if v, ok := matched[i]; ok {
			//log.Printf("\t matched and updating: " + strconv.Itoa(i))
			s.Interactions[i].addWaypoint(detected[v])
		} else {
			//log.Printf("\t not matched and removing: " + strconv.Itoa(i))
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
}

func saveScene(filename string, s *Scene) {
	b, _ := json.Marshal(s)
	ioutil.WriteFile(filename, b, 0611)

}

func sendInteraction(i Interaction) {
	// TODO: Broadcast the interaction to the mothership.
}
