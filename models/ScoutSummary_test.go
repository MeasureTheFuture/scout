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
)

func TestScoutSummary(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scout Summary Suite")
}

var _ = Describe("Scout Summary Model", func() {
	AfterEach(cleaner)

	Context("Insert", func() {
		It("Scout insert should create matching scout summary", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			ss, err := GetScoutSummaryByUUID(db, s.UUID)
			Ω(err).Should(BeNil())
			Ω(ss).Should(Equal(&ScoutSummary{s.UUID, 0, Buckets{}, IntBuckets{}}))
		})

		It("Should be able to update existing scout summary.", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			ss, err := GetScoutSummaryByUUID(db, s.UUID)
			Ω(err).Should(BeNil())
			ss.VisitorCount = 2
			ss.VisitTimeBuckets[1][5] = 0.1
			ss.VisitorBuckets[2][5] = 2
			err = ss.Update(db)
			Ω(err).Should(BeNil())

			ss2, err := GetScoutSummaryByUUID(db, s.UUID)
			Ω(err).Should(BeNil())
			Ω(ss2).Should(Equal(ss))
		})
	})

	Context("Get", func() {
		It("should be able to get scout healths as json", func() {
			ss := ScoutSummary{}
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ss,
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			ss.ScoutUUID = s.UUID
			Ω(err).Should(BeNil())

			jsonF, err := ScoutSummariesAsJSON(db)
			Ω(err).Should(BeNil())

			jsonB, err := ioutil.ReadFile(jsonF)
			Ω(err).Should(BeNil())

			var result []ScoutSummary
			err = json.Unmarshal(jsonB, &result)
			Ω(err).Should(BeNil())
			Ω(result).Should(Equal([]ScoutSummary{ss}))
		})
	})

	Context("Clear", func() {
		It("Should be able to clear an existing scout summary", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			ss, err := GetScoutSummaryByUUID(db, s.UUID)
			Ω(err).Should(BeNil())
			ss.VisitorCount = 2
			ss.VisitTimeBuckets[1][5] = 0.1
			ss.VisitorBuckets[2][5] = 2
			err = ss.Update(db)
			Ω(err).Should(BeNil())

			err = ss.Clear(db)
			Ω(err).Should(BeNil())

			ss2, err := GetScoutSummaryByUUID(db, s.UUID)
			Ω(err).Should(BeNil())
			Ω(ss2).Should(Equal(ss))
			Ω(ss2.VisitorCount).Should(Equal(int64(0)))
			Ω(ss2.VisitTimeBuckets[1][5]).Should(Equal(float32(0)))
			Ω(ss2.VisitorBuckets[2][5]).Should(Equal(int(0)))
		})
	})
})
