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

package models

import (
	"database/sql"
	"encoding/json"
	"github.com/MeasureTheFuture/scout/configuration"
	"io/ioutil"
	"math"
	"time"
)

type Scene struct {
	Interactions     []Interaction // The current interactions occuring within the scene.
	IdleInteractions []Interaction // The current interactions that are idle (resumable).
	sId              int
}

// initScene creates an empty scene that can be used for monitoring interactions.
func InitScene() *Scene {
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
			d := detected[i].distanceSq(s.Interactions[j].LastWaypoint())
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
func (s *Scene) addInteraction(detected []Waypoint, db *sql.DB, config configuration.Configuration) {
	if len(s.Interactions) == 0 {
		// Empty scene: just add a new interaction for each new waypoint.
		for i := 0; i < len(detected); i++ {
			s.Interactions = append(s.Interactions, NewInteraction(detected[i], s.sId, db))
			s.sId++
		}

	} else {
		// Existing scene:
		// for each of the detected waypoints
		// 	 work which of the existing interactions is the closest.
		//   for interactions that have more than one close detected waypoints
		//		create a new interaction from the furthest detected waypoint
		// 		the nearest waypoint is used to update the existing interaction.
		distances := s.buildDistanceMap(detected)

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
				// This detected element doesn't appear to belong to an existing interaction within the scene.
				// Before we create a new interaction, we check and if we can use the detected element to
				// resume an existing interaction.
				t := time.Now().UTC()
				resumed := false

				for k := len(s.IdleInteractions) - 1; k >= 0; k-- {
					wp := s.IdleInteractions[k].LastWaypoint()
					wpt := s.IdleInteractions[k].started.Add(time.Duration(wp.T) * time.Second)
					dt := float32(t.Sub(wpt).Seconds())

					if detected[i].distanceSq(wp) < config.ResumeSqDistance && dt < config.IdleDuration {
						// Resume idle interaction.
						s.IdleInteractions[k].addWaypoint(detected[i])
						s.Interactions = append(s.Interactions, s.IdleInteractions[k])
						s.IdleInteractions = append(s.IdleInteractions[:k], s.IdleInteractions[k+1:]...)
						resumed = true
					}
				}

				// We haven't resumed an idle interaction, so the detected element must be a new interaction.
				if !resumed {
					s.Interactions = append(s.Interactions, NewInteraction(detected[i], s.sId, db))
					s.sId++
				}
			}
		}
	}
}

func (s *Scene) removeInteraction(detected []Waypoint, config configuration.Configuration) {
	distances := s.buildDistanceMap(detected)
	matched := map[int]int{}

	for i := 0; i < len(distances); i++ {
		matched[distances[i][1]] = i
	}

	for i := len(s.Interactions) - 1; i >= 0; i-- {
		if v, ok := matched[i]; ok {
			s.Interactions[i].addWaypoint(detected[v])
		} else {
			// Interactions are not removed (and broadcasted to the mothership) immediately,
			// they are marked as idle first and can be subsequently resumed by waypoints
			// at a later time. Idle interactions eventual 'expire' and are removed completely
			// from the scene and broadcasted to the mothership.
			s.IdleInteractions = append(s.IdleInteractions, s.Interactions[i])
			s.Interactions = append(s.Interactions[:i], s.Interactions[i+1:]...)
		}
	}
}

func (s *Scene) Update(db *sql.DB, detected []Waypoint, config configuration.Configuration) {
	if len(detected) >= len(s.Interactions) {
		s.addInteraction(detected, db, config)
	} else {
		s.removeInteraction(detected, config)
	}

	// broadcast idle interactions that have expired and are no longer resumable.
	t := time.Now().UTC()
	for i := len(s.IdleInteractions) - 1; i >= 0; i-- {
		wp := s.IdleInteractions[i].LastWaypoint()
		wpt := s.IdleInteractions[i].started.Add(time.Duration(wp.T) * time.Second)
		dt := float32(t.Sub(wpt).Seconds())

		if dt > config.IdleDuration {
			// Only transmit the interaction to the mothership if it is longer than the
			// specified minimum duration. This is to filter out any detected noise.
			if s.IdleInteractions[i].Duration > config.MinDuration {
				s.IdleInteractions[i].saveToDB(db, config)
			}

			s.IdleInteractions = append(s.IdleInteractions[:i], s.IdleInteractions[i+1:]...)
		}
	}
}

func (s *Scene) Close(db *sql.DB, config configuration.Configuration) {
	// Broadcast all results to the mothership.
	for _, i := range s.Interactions {
		i.saveToDB(db, config)
	}

	for _, i := range s.IdleInteractions {
		i.saveToDB(db, config)
	}
}

func (s *Scene) save(filename string) {
	b, _ := json.Marshal(s)
	ioutil.WriteFile(filename, b, 0611)
}
