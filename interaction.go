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
	"strconv"
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

// addInteraction
func addInteraction(s *Scene, detected []Waypoint) {
	// The start time to use for the new interaction -- truncated to the nearest 30 minutes.
	start := time.Now().Truncate(30 * time.Minute)

	// Empty scene - just add a new interaction for each new waypoint.
	if len(s.Interactions) == 0 {
		for i := 0; i < len(detected); i++ {
			s.Interactions = append(s.Interactions, Interaction{start, 0.0, []Waypoint{detected[i]}})
		}

	} else {
		var closest map[int][][]int = make(map[int][][]int)

		for i := 0; i < len(s.Interactions); i++ {
			dist := math.MaxInt32
			minW := -1

			for j := 0; j < len(detected); j++ {
				d := s.Interactions[i].lastWaypoint().distanceSq(detected[j])
				if d < dist {
					dist = d
					minW = j
				}
			}

			closest[i] = append(closest[i], []int{minW, dist})
			//closest[i] = []int{minW, dist}
		}

		// TODO: detected is larger.
		// for each path in interactions
		// work out the closest waypoint in detected.
		// One or more interactions will have two or more close waypoints in detected.
		// The nearest waypoint in detected is used to update the pathway in interactions.
		// While the other waypoints are interactions new to the scene.
	}
}

func updateInteractions(s *Scene, detected []Waypoint) {

}

func monitorScene(s *Scene, detected []Waypoint) {
	for i := 0; i < len(detected); i++ {
		log.Printf("\t D: [" + strconv.Itoa(detected[i].XPixels) + "," + strconv.Itoa(detected[i].YPixels) + "]")
	}

	if len(detected) > len(s.Interactions) {
		// Someone new has entered the frame.
		log.Printf("\t New person")
		addInteraction(s, detected)

	} else if len(detected) == len(s.Interactions) {
		// Update the positions of everyone within the frame.
		log.Printf("\t Updating")

	} else {
		// Someone has left the frame.
		log.Printf("\t Person left")

	}

	log.Printf("")
}

func sendInteraction(i Interaction) {
	// TODO: Broadcast the interaction to the mothership.
}
