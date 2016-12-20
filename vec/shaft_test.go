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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestShaft(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shaft Suite")
}

var _ = Describe("Shaft", func() {
	Context("ShaftFromWaypoints", func() {
		It("Should create an enclosing shaft from two waypoints", func() {
			wpA := models.Waypoint{5, 5, 2, 4, 0.0}
			wpB := models.Waypoint{2, 3, 2, 2, 0.0}

			wpC := models.Waypoint{5, 1, 2, 1, 0.0}
			wpD := models.Waypoint{2, 2, 2, 1, 0.0}

			s := ShaftFromWaypoints(wpA, wpB, 10, 10)
			Ω(s).Should(Equal(Shaft{AABB{Vec{0, 1}, Vec{7, 9}},
				[2]Vec{Vec{0, 5}, Vec{3, 9}},
				[2]Vec{Vec{4, 1}, Vec{7, 1}}}))

			s2 := ShaftFromWaypoints(wpB, wpA, 10, 10)
			Ω(s2).Should(Equal(Shaft{AABB{Vec{0, 1}, Vec{7, 9}},
				[2]Vec{Vec{0, 5}, Vec{3, 9}},
				[2]Vec{Vec{4, 1}, Vec{7, 1}}}))

			s3 := ShaftFromWaypoints(wpC, wpD, 10, 10)
			Ω(s3).Should(Equal(Shaft{AABB{Vec{0, 0}, Vec{7, 3}},
				[2]Vec{Vec{3, 0}, Vec{0, 1}},
				[2]Vec{Vec{7, 2}, Vec{4, 3}}}))

			s4 := ShaftFromWaypoints(wpD, wpC, 10, 10)
			Ω(s4).Should(Equal(Shaft{AABB{Vec{0, 0}, Vec{7, 3}},
				[2]Vec{Vec{3, 0}, Vec{0, 1}},
				[2]Vec{Vec{7, 2}, Vec{4, 3}}}))
		})
	})

	Context("Intersects", func() {
		It("Should return true when an AABB intersects a shaft", func() {
			wpA := models.Waypoint{5, 5, 2, 4, 0.0}
			wpB := models.Waypoint{2, 3, 2, 2, 0.0}
			s := ShaftFromWaypoints(wpA, wpB, 10, 10)
			s2 := ShaftFromWaypoints(wpB, wpA, 10, 10)

			b1 := &AABB{Vec{6, 8}, Vec{9, 11}}
			Ω(s.Intersects(b1)).Should(BeTrue())
			Ω(s2.Intersects(b1)).Should(BeTrue())

			b2 := &AABB{Vec{0, 6}, Vec{2, 10}}
			Ω(s.Intersects(b2)).Should(BeTrue())
			Ω(s2.Intersects(b2)).Should(BeTrue())

			// TODO: More intersects tests.
		})

		It("Should return false when an AABB doesn't intersect a shaft", func() {
			wpA := models.Waypoint{5, 5, 2, 4, 0.0}
			wpB := models.Waypoint{2, 3, 2, 2, 0.0}
			s := ShaftFromWaypoints(wpA, wpB, 10, 10)
			s2 := ShaftFromWaypoints(wpB, wpA, 10, 10)

			b1 := &AABB{Vec{0, 10}, Vec{4, 14}}
			Ω(s.Intersects(b1)).Should(BeFalse())
			Ω(s2.Intersects(b1)).Should(BeFalse())

			b2 := &AABB{Vec{0, 8}, Vec{1, 10}}
			Ω(s.Intersects(b2)).Should(BeFalse())
			Ω(s2.Intersects(b2)).Should(BeFalse())

			// TODO: More intersects tests.
		})
	})

	Context("isLeft", func() {
		It("Should return true if a point is left of a line", func() {
			l := &[2]Vec{Vec{7, 1}, Vec{4, 3}}
			Ω(isLeft(l, Vec{5, 2})).Should(BeTrue())
		})

		It("Should return false if a point is right of a line", func() {
			l := &[2]Vec{Vec{4, 3}, Vec{7, 1}}
			Ω(isLeft(l, Vec{5, 2})).Should(BeFalse())
		})
	})
})
