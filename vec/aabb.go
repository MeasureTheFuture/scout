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

package vec

import (
	"github.com/MeasureTheFuture/scout/models"
)

// Axis Aligned Bounding Box
type AABB struct {
	Min Vec // Min is the minimum extents of the bounding box.
	Max Vec // Max is the maximum extents of the bounding box.
}

// AABBFromIndex creates an AABB from the supplied indexes, and the grid size.
func AABBFromIndex(i int, j int, iWidth int, jWidth int) AABB {
	min := Vec{i * iWidth, j * jWidth}
	return AABB{min, Vec{min[0] + iWidth, min[1] + jWidth}}
}

// AABBFromWaypoints creates an AABB from the two supplied waypoints and the
// maximum dimensions of the frame.
func AABBFromWaypoints(a models.Waypoint, b models.Waypoint, maxW int, maxH int) AABB {
	minX := Max(0, Min((a.XPixels-a.HalfWidthPixels), (b.XPixels-b.HalfWidthPixels)))
	minY := Max(0, Min((a.YPixels-a.HalfHeightPixels), (b.YPixels-b.HalfHeightPixels)))

	maxX := Min(maxW, Max((a.XPixels+a.HalfWidthPixels), (b.XPixels+b.HalfWidthPixels)))
	maxY := Min(maxH, Max((a.YPixels+a.HalfHeightPixels), (b.YPixels+b.HalfHeightPixels)))

	return AABB{Vec{minX, minY}, Vec{maxX, maxY}}
}

// AABBFromWaypoints creates an AABB from the supplied waypoing and the maximum
// dimensions of the frame.
func AABBFromWaypoint(a models.Waypoint, maxW int, maxH int) AABB {
	minX := Max(0, (a.XPixels - a.HalfWidthPixels))
	minY := Max(0, (a.YPixels - a.HalfHeightPixels))

	maxX := Min(maxW, (a.XPixels + a.HalfWidthPixels))
	maxY := Min(maxH, (a.YPixels + a.HalfHeightPixels))

	return AABB{Vec{minX, minY}, Vec{maxX, maxY}}
}

// Intersects returns true if a and b overlap (intersect), false otherwise.
func (b *AABB) Intersects(a *AABB) bool {
	return a.Max[0] >= b.Min[0] && a.Min[0] <= b.Max[0] &&
		a.Max[1] >= b.Min[1] && a.Min[1] <= b.Max[1]
}
