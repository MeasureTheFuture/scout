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

// Min returns the minimum value out of a and b.
func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum value out of a and b.
func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// MaxF returns the maximum float value out of a and b.
func MinF(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
