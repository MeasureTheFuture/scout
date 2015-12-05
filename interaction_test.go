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
			Ω(s).Should(Equal(Scene{}))
		})
	})

	Context("WayPoint", func() {
		It("should be able to calculate the distance squared between two way points", func() {
			a := Waypoint{5, 5, 10, 10, 0.0}
			b := Waypoint{3, 3, 10, 10, 0.0}

			Ω(a.distanceSq(b)).Should(Equal(8))
		})
	})

	Context("addInteraction", func() {
		wpA := Waypoint{100, 100, 20, 20, 0.0}
		wpAA := Waypoint{102, 100, 20, 20, 0.0}
		wpB := Waypoint{50, 50, 20, 20, 0.0}
		wpBA := Waypoint{55, 53, 20, 20, 0.0}
		wpC := Waypoint{150, 150, 20, 20, 0.0}

		It("should be able to add an interaction to an empty scene", func() {
			s := initScene()
			addInteraction(&s, []Waypoint{wpA})

			Ω(len(s.Interactions)).Should(Equal(1))
			Ω(s.Interactions[0].Path).Should(Equal([]Waypoint{wpA}))
		})

		It("should be able to add multiple interactions to an empty scene,", func() {
			s := initScene()

			addInteraction(&s, []Waypoint{wpA, wpB})

			Ω(len(s.Interactions)).Should(Equal(2))
			Ω(s.Interactions[0].Path).Should(Equal([]Waypoint{wpA}))
			Ω(s.Interactions[1].Path).Should(Equal([]Waypoint{wpB}))
		})

		It("should list the interaction start time truncated to 30 mins", func() {
			s := initScene()
			addInteraction(&s, []Waypoint{wpA})

			Ω(s.Interactions[0].Entered).Should(Equal(time.Now().Truncate(30 * time.Minute)))
		})

		It("should be able to add an interaction to a scene with stuff already going on", func() {
			s := initScene()
			addInteraction(&s, []Waypoint{wpA})
			addInteraction(&s, []Waypoint{wpAA, wpB})

			Ω(len(s.Interactions)).Should(Equal(2))
			Ω(s.Interactions[0].Path).Should(Equal([]Waypoint{wpA, wpAA}))
			Ω(s.Interactions[1].Path).Should(Equal([]Waypoint{wpB}))
		})

		It("should be able to add multiple interactions to a scene with stuff already going on", func() {
			s := initScene()
			addInteraction(&s, []Waypoint{wpA})
			addInteraction(&s, []Waypoint{wpAA, wpB})
			addInteraction(&s, []Waypoint{wpAA, wpBA, wpC})

			Ω(len(s.Interactions)).Should(Equal(3))
			Ω(s.Interactions[0].Path).Should(Equal([]Waypoint{wpA, wpAA, wpAA}))
			Ω(s.Interactions[1].Path).Should(Equal([]Waypoint{wpB, wpBA}))
			Ω(s.Interactions[2].Path).Should(Equal([]Waypoint{wpC}))
		})
	})
})
