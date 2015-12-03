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

type Scene struct {
	interactions []Interaction // The current interactions occuring within the scene.
}

// initScene creates an empty scene that can be used for monitoring interactions.
func initScene() Scene {
	return Scene{}
}

// addInteraction
func addInteraction(s *Scene, detected []Waypoint) {
	// Work out which of the detected waypoints are new.
	var matched []bool = make([]bool, len(s.interactions))

	for i := 0; i < len(detected); i++ {
		for j := 0; j < len(s.interactions); j++ {

			//dist :=
		}
	}
}

func updateInteractions(s *Scene, detected []Waypoint) {

}

func monitorScene(s *Scene, detected []Waypoint) {
	for i := 0; i < len(detected); i++ {
		log.Printf("\t D: [" + strconv.Itoa(detected[i].XPixels) + "," + strconv.Itoa(detected[i].YPixels) + "]")
	}

	if len(detected) > len(s.interactions) {
		// Someone new has entered the frame.
		log.Printf("\t New person")
		addInteraction(s, detected)

	} else if len(detected) == len(s.interactions) {
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
