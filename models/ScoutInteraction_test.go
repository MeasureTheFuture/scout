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
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"testing"
	"time"
)

func TestScoutInteraction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scout Interaction Suite")
}

var _ = Describe("Scout Interaction Model", func() {
	AfterEach(cleaner)

	Context("CreateScoutInteraction", func() {
		It("Should be able to create a scout interaction", func() {
			t := time.Now().UTC()

			wp := []Waypoint{Waypoint{1, 2, 3, 4, 0.1}}
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())
			i := Interaction{"abc", "0.1", t, t, 0.1, wp, 1, &s}

			si := CreateScoutInteraction(&i)
			Ω(si.ScoutUUID).Should(Equal(""))
			Ω(si.Duration).Should(Equal(i.Duration))
			Ω(si.EnteredAt).Should(Equal(t))
			Ω(si.Waypoints).Should(Equal(Path{[2]int{1, 2}}))
			Ω(si.WaypointWidths).Should(Equal(Path{[2]int{3, 4}}))
			Ω(si.WaypointTimes).Should(Equal(RealArray{0.1}))
		})
	})

	Context("Get", func() {
		It("should be able to get scout interactions as json", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			et := time.Now().UTC().Round(15 * time.Minute)
			si := ScoutInteraction{-1, s.UUID, 0.2, Path{[2]int{1, 2}, [2]int{5, 6}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, et}
			err = si.Insert(db)
			Ω(err).Should(BeNil())

			jsonF, err := ScoutInteractionsAsJSON(db)
			Ω(err).Should(BeNil())

			jsonB, err := ioutil.ReadFile(jsonF)
			Ω(err).Should(BeNil())

			var result []ScoutInteraction
			err = json.Unmarshal(jsonB, &result)
			Ω(err).Should(BeNil())
			Ω(result).Should(Equal([]ScoutInteraction{si}))
		})
	})

	Context("Insert", func() {
		It("Should be able to insert a scout interaction", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			et := time.Now().UTC().Round(15 * time.Minute)
			si := ScoutInteraction{-1, s.UUID, 0.2, Path{[2]int{1, 2}, [2]int{5, 6}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, et}
			err = si.Insert(db)
			Ω(err).Should(BeNil())

			si2, err := GetScoutInteractionById(db, si.Id)
			Ω(err).Should(BeNil())
			Ω(si2).Should(Equal(&si))
		})
	})

	Context("Delete", func() {
		It("Should be able to delete interactions for a specified scout", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			si := ScoutInteraction{-1, s.UUID, 0.2, Path{[2]int{1, 2}, [2]int{5, 6}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, time.Now()}
			err = si.Insert(db)
			Ω(err).Should(BeNil())

			si2 := ScoutInteraction{-1, s.UUID, 0.2, Path{[2]int{1, 2}, [2]int{5, 6}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, time.Now()}
			err = si2.Insert(db)
			Ω(err).Should(BeNil())

			err = DeleteScoutInteractions(db, s.UUID)
			Ω(err).Should(BeNil())

			n, err := NumScoutInteractions(db)
			Ω(err).Should(BeNil())
			Ω(n).Should(Equal(int64(0)))
		})
	})

	Context("Unprocessed", func() {
		It("Should be able to get unproccessed interactions", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			et := time.Now().UTC().Round(15 * time.Minute)
			si1 := ScoutInteraction{-1, s.UUID, 0.2, Path{[2]int{1, 2}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, et}
			err = si1.Insert(db)
			Ω(err).Should(BeNil())

			si2 := ScoutInteraction{-1, s.UUID, 0.3, Path{[2]int{1, 2}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, et}
			err = si2.Insert(db)
			Ω(err).Should(BeNil())

			si3 := ScoutInteraction{-1, s.UUID, 0.4, Path{[2]int{1, 2}}, Path{[2]int{3, 4}}, RealArray{0.1}, true, et}
			err = si3.Insert(db)
			Ω(err).Should(BeNil())

			up, err := GetUnprocessed(db)
			Ω(err).Should(BeNil())
			Ω(up).Should(Equal([]*ScoutInteraction{&si1, &si2}))
		})
	})

	Context("MarkProcessed", func() {
		It("Should be able to mark interactions as processed", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			et := time.Now().UTC().Round(15 * time.Minute)
			si1 := ScoutInteraction{-1, s.UUID, 0.2, Path{[2]int{1, 2}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, et}
			err = si1.Insert(db)
			Ω(err).Should(BeNil())

			err = si1.MarkProcessed(db)
			Ω(err).Should(BeNil())

			Ω(si1.Processed).Should(BeTrue())
			up, err := GetUnprocessed(db)
			Ω(err).Should(BeNil())
			Ω(up).Should(BeNil())
		})
	})

	Context("Init scene", func() {
		It("should be able to init an empty scene", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())
			si := InitScene(&s)
			Ω(*si).Should(Equal(Scene{[]Interaction{}, []Interaction{}, 0, &s}))
		})
	})

	Context("WayPoint", func() {
		It("should be able to calculate the distance squared between two way points", func() {
			a := Waypoint{5, 5, 10, 10, 0.0}
			b := Waypoint{3, 3, 10, 10, 0.0}

			Ω(a.distanceSq(b)).Should(Equal(int64(8)))
		})

		It("should be able to calculate the perpendicular distance to a line", func() {
			x := Waypoint{2, 0, 0, 0, 0.0}
			a := Waypoint{0, 0, 0, 0, 0.0}
			b := Waypoint{0, 2, 0, 0, 0.0}

			Ω(x.perpendicularDistance(a, b)).Should(BeNumerically("~", 2.0, 0.001))
		})

		It("should be able to test equality between two waypoints", func() {
			a := Waypoint{5, 4, 10, 10, 0.1}
			b := Waypoint{5, 4, 10, 10, 0.1}
			c := Waypoint{5, 4, 10, 10, 0.2}

			Ω(a.Equal(b)).Should(BeTrue())
			Ω(b.Equal(a)).Should(BeTrue())
			Ω(a.Equal(c)).Should(BeFalse())
			Ω(c.Equal(b)).Should(BeFalse())
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

	Context("NewInteraction", func() {
		It("should create a new interaction", func() {
			a := Waypoint{0, 0, 0, 0, 0.0}
			tr := time.Now().UTC().Round(15 * time.Minute)

			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			i := NewInteraction(a, 0, &s)
			Ω(i.UUID).Should(Equal(s.UUID))
			Ω(i.Version).Should(Equal("0.1"))
			Ω(i.Entered).Should(Equal(tr))
			Ω(i.Duration).Should(BeNumerically("~", float32(0.0), 0.007))
			Ω(i.Equal([]Waypoint{a})).Should(BeTrue())
		})

		It("should add a new waypoint", func() {
			a := Waypoint{0, 0, 0, 0, 0.0}
			b := Waypoint{1, 1, 1, 1, 0.005}

			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			i := NewInteraction(a, 0, &s)
			time.Sleep(50 * time.Millisecond)
			i.addWaypoint(b)
			Ω(i.Duration).Should(BeNumerically("~", float32(0.05), 0.007))
			Ω(len(i.Path)).Should(Equal(2))
		})
	})

	Context("addInteraction", func() {
		wpA := Waypoint{100, 100, 20, 20, 0.0}
		wpAA := Waypoint{102, 100, 20, 20, 0.0}
		wpAAT := Waypoint{102, 100, 20, 20, 0.10}
		wpB := Waypoint{50, 50, 20, 20, 0.0}
		wpBA := Waypoint{55, 53, 20, 20, 0.0}
		wpC := Waypoint{150, 150, 20, 20, 0.0}

		It("should be able to add an interaction to an empty scene", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			si := InitScene(&s)
			si.addInteraction([]Waypoint{wpA})

			Ω(len(si.Interactions)).Should(Equal(1))
			Ω(si.Interactions[0].Equal([]Waypoint{wpA})).Should(BeTrue())
		})

		It("should be able to add multiple interactions to an empty scene,", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			si := InitScene(&s)
			si.addInteraction([]Waypoint{wpA, wpB})

			Ω(len(si.Interactions)).Should(Equal(2))
			Ω(si.Interactions[0].Equal([]Waypoint{wpA})).Should(BeTrue())
			Ω(si.Interactions[1].Equal([]Waypoint{wpB})).Should(BeTrue())
		})

		It("should list the interaction start time truncated to 30 mins", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			si := InitScene(&s)
			si.addInteraction([]Waypoint{wpA})

			Ω(si.Interactions[0].Entered).Should(Equal(time.Now().UTC().Round(15 * time.Minute)))
		})

		It("should be able to add an interaction to a scene with stuff already going on", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			si := InitScene(&s)
			si.addInteraction([]Waypoint{wpA})

			time.Sleep(100 * time.Millisecond)
			si.addInteraction([]Waypoint{wpAA, wpB})

			Ω(len(si.Interactions)).Should(Equal(2))
			Ω(si.Interactions[0].Equal([]Waypoint{wpA, wpAAT})).Should(BeTrue())
			Ω(si.Interactions[1].Equal([]Waypoint{wpB})).Should(BeTrue())
		})

		It("should be able to add multiple interactions to a scene with stuff already going on", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			si := InitScene(&s)
			si.addInteraction([]Waypoint{wpA})
			si.addInteraction([]Waypoint{wpAA, wpB})
			si.addInteraction([]Waypoint{wpAA, wpBA, wpC})

			Ω(len(si.Interactions)).Should(Equal(3))
			Ω(si.Interactions[0].Equal([]Waypoint{wpA, wpAA, wpAA})).Should(BeTrue())
			Ω(si.Interactions[1].Equal([]Waypoint{wpB, wpBA})).Should(BeTrue())
			Ω(si.Interactions[2].Equal([]Waypoint{wpC})).Should(BeTrue())
		})

		It("should be able to remove interactions when a person leaves the scene", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			si := InitScene(&s)
			si.addInteraction([]Waypoint{wpA, wpB})
			si.removeInteraction([]Waypoint{wpAA})

			Ω(len(si.Interactions)).Should(Equal(1))
			Ω(si.Interactions[0].Equal([]Waypoint{wpA, wpAA})).Should(BeTrue())
		})

		It("should be able to remove multiple interactions when more than one person leaves the scene", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			si := InitScene(&s)
			si.addInteraction([]Waypoint{wpA, wpB, wpC})
			si.removeInteraction([]Waypoint{wpBA})

			Ω(len(si.Interactions)).Should(Equal(1))
			Ω(si.Interactions[0].Equal([]Waypoint{wpB, wpBA})).Should(BeTrue())
		})
	})
})
