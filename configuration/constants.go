/*
 * Copyright (C) 2017 Clinton Freeman
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

package configuration

const (
	FrameW   = 1920              // The width of the frames captured by the Webcam in pixels.
	FrameH   = 1080              // The height of the frames captured by the webcam in pixels.
	WBuckets = 20                // The number of horizontal buckets a frame is broken into.
	HBuckets = 20                // The number of vertical buckets a frame is broken into.
	BucketW  = FrameW / WBuckets // The width of a bucket in pixels.
	BucketH  = FrameH / HBuckets // The height of a bucket in pixels
)
