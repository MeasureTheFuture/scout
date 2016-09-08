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
	"math"
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
