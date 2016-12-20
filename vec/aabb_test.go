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

func TestAABB(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AABB Suite")
}

var _ = Describe("AABB", func() {
	Context("CreateAABB", func() {
		It("should create an AABB enclosing two waypoints", func() {
			wpA := models.Waypoint{5, 5, 2, 4, 0.0}
			wpB := models.Waypoint{2, 3, 2, 2, 0.0}

			Ω(AABBFromWaypoints(wpA, wpB, 10, 10)).Should(Equal(AABB{Vec{0, 1}, Vec{7, 9}}))
		})

		It("should create an AABB enclosing one waypoint", func() {
			wp := models.Waypoint{5, 5, 2, 4, 0.0}

			Ω(AABBFromWaypoint(wp, 10, 10)).Should(Equal(AABB{Vec{3, 1}, Vec{7, 9}}))
		})

		It("should create an AABB from i and j bucket indexes", func() {
			Ω(AABBFromIndex(0, 0, 51, 36)).Should(Equal(AABB{Vec{0, 0}, Vec{51, 36}}))
			Ω(AABBFromIndex(2, 3, 51, 36)).Should(Equal(AABB{Vec{102, 108}, Vec{153, 144}}))
		})
	})

	Context("Intersects", func() {
		PIt("should return true when two AABBs intersect", func() {
			a := AABB{Vec{0, 0}, Vec{2, 2}}
			b := AABB{Vec{1, 1}, Vec{2, 2}}
			c := AABB{Vec{1, 1}, Vec{2, 4}}
			d := AABB{Vec{1, 1}, Vec{4, 2}}

			Ω(a.Intersects(&b)).Should(BeTrue())
			Ω(b.Intersects(&a)).Should(BeTrue())
			Ω(a.Intersects(&c)).Should(BeTrue())
			Ω(c.Intersects(&a)).Should(BeTrue())
			Ω(a.Intersects(&d)).Should(BeTrue())
			Ω(d.Intersects(&a)).Should(BeTrue())
		})

		It("should return false when two AABBs don't intersect", func() {
			a := AABB{Vec{0, 0}, Vec{2, 2}}
			b := AABB{Vec{3, 0}, Vec{5, 2}}
			c := AABB{Vec{3, 3}, Vec{5, 5}}
			d := AABB{Vec{0, 3}, Vec{2, 5}}

			Ω(a.Intersects(&b)).Should(BeFalse())
			Ω(b.Intersects(&a)).Should(BeFalse())
			Ω(a.Intersects(&c)).Should(BeFalse())
			Ω(c.Intersects(&a)).Should(BeFalse())
			Ω(a.Intersects(&d)).Should(BeFalse())
			Ω(d.Intersects(&a)).Should(BeFalse())
		})
	})
})
