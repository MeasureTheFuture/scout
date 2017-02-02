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

package models

import (
 	"database/sql"
	"log"
	"time"
)

type Interaction struct {
	UUID     string     // The UUID for the scout that detected the interaction.
	Version  string     // The Version of the protocol used for transmitting data to the mothership
	Entered  time.Time  // The time the interaction started (rounded to nearest half hour)
	started  time.Time  // The actual time the interaction started. Private. Not to be transmitted for privacy concerns
	Duration float32    // The total duration of the interaction.
	Path     []Waypoint // The pathway of the interaction through the scene.
	SceneID  int
	dScout *Scout
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

func NewInteraction(w Waypoint, sId int, s *Scout) Interaction {
	start := time.Now().UTC()

	// The start time broadcasted for the interaction is truncated to the nearest 30 minutes.
	apparentStart := start.Round(15 * time.Minute)

	i := Interaction{s.UUID, "0.1", apparentStart, start, 0.0, []Waypoint{}, sId, s}
	i.addWaypoint(w)
	return i
}

// addWaypoint inserts a new waypoint to the end of the interaction.
func (i *Interaction) addWaypoint(w Waypoint) {
	newW := w
	newW.T = float32(time.Now().UTC().Sub(i.started).Seconds())

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
func (i *Interaction) LastWaypoint() Waypoint {
	return i.Path[len(i.Path)-1]
}

func (i *Interaction) simplify() {
	i.Path = douglasPeucker(i.Path, i.dScout.SimplifyEpsilon)
}

func (i *Interaction) saveToDB(db *sql.DB) {
	i.simplify() // Remove unnecessary segments from the pathway before storing in the DB.

	si := CreateScoutInteraction(i)
	si.ScoutUUID = i.dScout.UUID
	err := si.Insert(db)
	if err != nil {
		log.Printf("ERROR: Unable to save Interaction to DB.")
		log.Print(err)
	}
}
