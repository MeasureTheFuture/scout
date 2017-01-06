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

package processes

import (
	"database/sql"
	"github.com/MeasureTheFuture/mothership/configuration"
	"github.com/MeasureTheFuture/mothership/models"
	"github.com/MeasureTheFuture/mothership/vec"
	"log"
	"time"
)

func Summarise(db *sql.DB, c configuration.Configuration) {
	poll := time.NewTicker(time.Millisecond * time.Duration(c.SummariseInterval)).C

	for {
		select {
		case <-poll:
			updateUnprocessed(db)
		}
	}
}

func updateUnprocessed(db *sql.DB) {
	up, err := models.GetUnprocessed(db)
	if err != nil {
		log.Printf("ERROR: Summarise unable to get unprocessed scout interactions.")
		log.Print(err)
		return
	}

	for _, si := range up {
		ss, err := models.GetScoutSummaryById(db, si.ScoutId)
		if err != nil {
			log.Printf("ERROR: Summarise unable to get scout summary")
			log.Print(err)
			return
		}

		ss.VisitorCount += 1
		updateTimeBuckets(db, ss, si)

		err = ss.Update(db)
		if err != nil {
			log.Printf("ERROR: Summarise unable to update scout summary")
			log.Print(err)
		}

		err = si.MarkProcessed(db)
		if err != nil {
			log.Printf("ERROR: Summarise unable to make scout interaction as processed")
			log.Print(err)
			return
		}
	}
}

const (
	FrameW   = 1920
	FrameH   = 1080
	WBuckets = 20
	HBuckets = 20
	BucketW  = FrameW / WBuckets
	BucketH  = FrameH / HBuckets
)

func maxTravelTime(a models.Waypoint, b models.Waypoint) float32 {
	travelD := vec.Vec{(b.XPixels - a.XPixels), (b.YPixels - a.YPixels)}
	travelG := float32(travelD[1]) / float32(travelD[0])

	x := float32(FrameW) / float32(WBuckets)
	y := x * travelG
	bucketD := vec.Vec{int(x), int(y)}

	// Make sure that we don't overallocate time for the bucket, when the
	// travelD is shorter than the bucket itself, the maximum multiplication
	// value for the maxTravel time is 1.0.
	f := vec.MinF(float32(1.0), float32(bucketD.Length()/travelD.Length()))
	return (b.T - a.T) * f
}

func updateTimeBuckets(db *sql.DB, ss *models.ScoutSummary, si *models.ScoutInteraction) {
	var intersected [configuration.HBuckets][configuration.WBuckets]bool

	// For each segment in an interaction.
	for k := 0; k < (len(si.Waypoints) - 1); k++ {
		// Generate a shaft AABB from the two waypoints.
		wpA := models.Waypoint{si.Waypoints[k][0], si.Waypoints[k][1],
			si.WaypointWidths[k][0], si.WaypointWidths[k][1], si.WaypointTimes[k]}
		wpB := models.Waypoint{si.Waypoints[k+1][0], si.Waypoints[k+1][1],
			si.WaypointWidths[k+1][0], si.WaypointWidths[k+1][1], si.WaypointTimes[k+1]}

		s := vec.ShaftFromWaypoints(wpA, wpB, FrameW, FrameH)

		// Work out maximum travel time that can be spent in a bucket.
		mt := maxTravelTime(wpA, wpB)

		// For each of the buckets, see if it intersects the shaft AABB and if it does
		// increment the bucket time by the maximum travel time.
		for i := 0; i < WBuckets; i++ {
			for j := 0; j < HBuckets; j++ {
				bucket := vec.AABBFromIndex(i, j, BucketW, BucketH)
				// TODO: Possibly improve time estimate by working out how much of
				// the bucket overlaps the shaft. Use it as a ratio between 0 and 1
				// to multiply max time.
				//
				// At the moment we allocate mt (the maximum possible travel time)
				// to the bucket once per interaction. Event if more than one segment
				// intersects this bucket.
				if s.Intersects(&bucket) {
					if !intersected[i][j] {
						ss.VisitTimeBuckets[i][j] += mt
						ss.VisitorBuckets[i][j] += 1
						intersected[i][j] = true
					}
				}
			}
		}
	}
}
