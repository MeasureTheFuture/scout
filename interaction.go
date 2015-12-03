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
	XPixels      int
	YPixels      int
	WidthPixels  int
	HeightPixels int
	T            float32
}

type Interaction struct {
	Entered  time.Time
	Duration float32

	Path []Waypoint
}

type Scene struct {
	interactions []Interaction
}

func initScene() Scene {
	return Scene{make([]Interaction, 10)}
}

func monitorScene(s Scene, detected []Waypoint) Scene {
	for i := 0; i < len(detected); i++ {
		log.Printf("\t D: [" + strconv.Itoa(detected[i].XPixels) + "," + strconv.Itoa(detected[i].YPixels) + "]")
	}

	log.Printf("")
	return s
}

func sendInteraction(i Interaction) {
	// TODO: Broadcast the interaction to the mothership.
}
