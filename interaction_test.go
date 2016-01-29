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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestInteraction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Interaction Suite")
}

var _ = Describe("Interaction", func() {
	Context("Init scene", func() {
		It("should be able to init an empty scene", func() {
			s := initScene()
			Ω(*s).Should(Equal(Scene{}))
		})
	})

	Context("WayPoint", func() {
		It("should be able to calculate the distance squared between two way points", func() {
			a := Waypoint{5, 5, 10, 10, 0.0}
			b := Waypoint{3, 3, 10, 10, 0.0}

			Ω(a.distanceSq(b)).Should(Equal(8))
		})

		It("should be able to calculate the perpendicular distance to a line", func() {
			x := Waypoint{2, 0, 0, 0, 0.0}
			a := Waypoint{0, 0, 0, 0, 0.0}
			b := Waypoint{0, 2, 0, 0, 0.0}

			Ω(x.perpendicularDistance(a, b)).Should(BeNumerically("~", 2.0, 0.001))
		})
	})

	Context("douglasPeucker", func() {
		It("should handle small sized paths", func() {
			a := Waypoint{0, 0, 0, 0, 0.0}
			b := Waypoint{0, 2, 0, 0, 0.0}

			Ω(douglasPeucker([]Waypoint{a, b}, 2)).Should(Equal([]Waypoint{a, b}))
			Ω(douglasPeucker([]Waypoint{a}, 2)).Should(Equal([]Waypoint{a}))
		})

		It("should remove waypoints if perpendicular distance is less than epsilon", func() {
			a := Waypoint{0, 0, 0, 0, 0.0}
			b := Waypoint{2, 1, 0, 0, 0.0}
			c := Waypoint{2, 2, 0, 0, 0.0}
			d := Waypoint{0, 4, 0, 0, 0.0}

			Ω(douglasPeucker([]Waypoint{a, b, d}, 3)).Should(Equal([]Waypoint{a, d}))
			Ω(douglasPeucker([]Waypoint{a, b, d}, 1)).Should(Equal([]Waypoint{a, b, d}))
			Ω(douglasPeucker([]Waypoint{a, b, c, d}, 1.9)).Should(Equal([]Waypoint{a, b, d}))
		})
	})

	Context("addInteraction", func() {
		wpA := Waypoint{100, 100, 20, 20, 0.0}
		wpAA := Waypoint{102, 100, 20, 20, 0.0}
		wpAAT := Waypoint{102, 100, 20, 20, 0.10}
		wpB := Waypoint{50, 50, 20, 20, 0.0}
		wpBA := Waypoint{55, 53, 20, 20, 0.0}
		wpC := Waypoint{150, 150, 20, 20, 0.0}
		c := Configuration{2.0, 2, 2, 2, 2, 2.0, 0, ":9090", "127.0.0.1:9091", "abc", 2.0}

		It("should be able to add an interaction to an empty scene", func() {
			s := initScene()
			s.addInteraction([]Waypoint{wpA}, c)

			Ω(len(s.Interactions)).Should(Equal(1))
			Ω(s.Interactions[0].Path).Should(Equal([]Waypoint{wpA}))
		})

		It("should be able to add multiple interactions to an empty scene,", func() {
			s := initScene()
			s.addInteraction([]Waypoint{wpA, wpB}, c)

			Ω(len(s.Interactions)).Should(Equal(2))
			Ω(s.Interactions[0].Path).Should(Equal([]Waypoint{wpA}))
			Ω(s.Interactions[1].Path).Should(Equal([]Waypoint{wpB}))
		})

		It("should list the interaction start time truncated to 30 mins", func() {
			s := initScene()
			s.addInteraction([]Waypoint{wpA}, c)

			Ω(s.Interactions[0].Entered).Should(Equal(time.Now().Round(15 * time.Minute)))
		})

		It("should be able to add an interaction to a scene with stuff already going on", func() {
			s := initScene()
			s.addInteraction([]Waypoint{wpA}, c)

			time.Sleep(100 * time.Millisecond)
			s.addInteraction([]Waypoint{wpAA, wpB}, c)

			Ω(len(s.Interactions)).Should(Equal(2))
			Ω(s.Interactions[0].Path[0].compare(wpA)).Should(BeTrue())
			Ω(s.Interactions[0].Path[1].compare(wpAAT)).Should(BeTrue())
			Ω(s.Interactions[1].Path).Should(Equal([]Waypoint{wpB}))
		})

		It("should be able to add multiple interactions to a scene with stuff already going on", func() {
			s := initScene()
			s.addInteraction([]Waypoint{wpA}, c)
			s.addInteraction([]Waypoint{wpAA, wpB}, c)
			s.addInteraction([]Waypoint{wpAA, wpBA, wpC}, c)

			Ω(len(s.Interactions)).Should(Equal(3))
			Ω(s.Interactions[0].Path[0].compare(wpA)).Should(BeTrue())
			Ω(s.Interactions[0].Path[1].compare(wpAA)).Should(BeTrue())
			Ω(s.Interactions[0].Path[2].compare(wpAA)).Should(BeTrue())

			Ω(s.Interactions[1].Path[0].compare(wpB)).Should(BeTrue())
			Ω(s.Interactions[1].Path[1].compare(wpBA)).Should(BeTrue())

			Ω(s.Interactions[2].Path[0].compare(wpC)).Should(BeTrue())
		})

		It("should be able to remove interactions when a person leaves the scene", func() {
			s := initScene()
			s.addInteraction([]Waypoint{wpA, wpB}, c)
			s.removeInteraction([]Waypoint{wpAA}, c)

			Ω(len(s.Interactions)).Should(Equal(1))
			Ω(len(s.Interactions[0].Path)).Should(Equal(2))

			Ω(s.Interactions[0].Path[0].compare(wpA)).Should(BeTrue())
			Ω(s.Interactions[0].Path[1].compare(wpAA)).Should(BeTrue())
		})

		It("should be able to remove multiple interactions when more than one person leaves the scene", func() {
			s := initScene()
			s.addInteraction([]Waypoint{wpA, wpB, wpC}, c)
			s.removeInteraction([]Waypoint{wpBA}, c)

			Ω(len(s.Interactions)).Should(Equal(1))
			Ω(s.Interactions[0].Path[0].compare(wpB)).Should(BeTrue())
			Ω(s.Interactions[0].Path[1].compare(wpBA)).Should(BeTrue())
		})
	})
})
