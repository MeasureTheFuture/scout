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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			i := Interaction{"abc", "0.1", t, 0.1, wp}

			si := CreateScoutInteraction(&i)
			Ω(si.ScoutId).Should(Equal(int64(-1)))
			Ω(si.Duration).Should(Equal(i.Duration))
			Ω(si.EnteredAt).Should(Equal(t))
			Ω(si.Waypoints).Should(Equal(Path{[2]int{1, 2}}))
			Ω(si.WaypointWidths).Should(Equal(Path{[2]int{3, 4}}))
			Ω(si.WaypointTimes).Should(Equal(RealArray{0.1}))
		})
	})

	Context("Insert", func() {
		It("Should be able to insert a scout interaction", func() {
			s := Scout{-1, "800fd548-2d2b-4185-885d-6323ccbe88a0", "192.168.0.1", 8080, true, "foo",
				"idle", &ScoutSummary{}}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			et := time.Now().UTC().Round(15 * time.Minute)
			si := ScoutInteraction{-1, s.Id, 0.2, Path{[2]int{1, 2}, [2]int{5, 6}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, et}
			err = si.Insert(db)
			Ω(err).Should(BeNil())

			si2, err := GetScoutInteractionById(db, si.Id)
			Ω(err).Should(BeNil())
			Ω(si2).Should(Equal(&si))
		})
	})

	Context("Delete", func() {
		It("Should be able to delete interactions for a specified scout", func() {
			s := Scout{-1, "800fd548-2d2b-4185-885d-6323ccbe88a0", "192.168.0.1", 8080, true, "foo",
				"idle", &ScoutSummary{}}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			si := ScoutInteraction{-1, s.Id, 0.2, Path{[2]int{1, 2}, [2]int{5, 6}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, time.Now()}
			err = si.Insert(db)
			Ω(err).Should(BeNil())

			si2 := ScoutInteraction{-1, s.Id, 0.2, Path{[2]int{1, 2}, [2]int{5, 6}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, time.Now()}
			err = si2.Insert(db)
			Ω(err).Should(BeNil())

			err = DeleteScoutInteractions(db, s.Id)
			Ω(err).Should(BeNil())

			n, err := NumScoutInteractions(db)
			Ω(err).Should(BeNil())
			Ω(n).Should(Equal(int64(0)))
		})
	})

	Context("Unprocessed", func() {
		It("Should be able to get unproccessed interactions", func() {
			s := Scout{-1, "800fd548-2d2b-4185-885d-6323ccbe88a0", "192.168.0.1", 8080, true, "foo",
				"idle", &ScoutSummary{}}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			et := time.Now().UTC().Round(15 * time.Minute)
			si1 := ScoutInteraction{-1, s.Id, 0.2, Path{[2]int{1, 2}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, et}
			err = si1.Insert(db)
			Ω(err).Should(BeNil())

			si2 := ScoutInteraction{-1, s.Id, 0.3, Path{[2]int{1, 2}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, et}
			err = si2.Insert(db)
			Ω(err).Should(BeNil())

			si3 := ScoutInteraction{-1, s.Id, 0.4, Path{[2]int{1, 2}}, Path{[2]int{3, 4}}, RealArray{0.1}, true, et}
			err = si3.Insert(db)
			Ω(err).Should(BeNil())

			up, err := GetUnprocessed(db)
			Ω(err).Should(BeNil())
			Ω(up).Should(Equal([]*ScoutInteraction{&si1, &si2}))
		})
	})

	Context("MarkProcessed", func() {
		It("Should be able to mark interactions as processed", func() {
			s := Scout{-1, "800fd548-2d2b-4185-885d-6323ccbe88a0", "192.168.0.1", 8080, true, "foo",
				"idle", &ScoutSummary{}}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			et := time.Now().UTC().Round(15 * time.Minute)
			si1 := ScoutInteraction{-1, s.Id, 0.2, Path{[2]int{1, 2}}, Path{[2]int{3, 4}}, RealArray{0.1}, false, et}
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
})
