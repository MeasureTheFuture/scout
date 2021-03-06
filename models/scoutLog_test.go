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
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestScoutLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scout Log Suite")
}

var _ = Describe("Scout Health Model", func() {
	AfterEach(cleaner)

	Context("Insert", func() {
		It("should insert a valid scout_log into the DB.", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "calibrated", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1, 4.0}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			t := time.Now()
			sl := ScoutLog{s.UUID, []byte("abc"), t}
			err = sl.Insert(db)
			Ω(err).Should(BeNil())

			sl2, err := GetScoutLogByUUID(db, s.UUID, t)
			Ω(err).Should(BeNil())
			Ω(&sl).Should(Equal(sl2))

		})

		It("should return an error when an invalid scout health is inserted into the DB.", func() {
			sl := ScoutLog{"", []byte("abc"), time.Now()}
			err := sl.Insert(db)
			Ω(err).ShouldNot(BeNil())
		})
	})

	Context("Delete", func() {
		It("should be able to delete logs for a specified scout", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "calibrated", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1, 4.0}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			sl := ScoutLog{s.UUID, []byte("abc"), time.Now()}
			err = sl.Insert(db)
			Ω(err).Should(BeNil())

			sl2 := ScoutLog{s.UUID, []byte("abc"), time.Now()}
			err = sl2.Insert(db)
			Ω(err).Should(BeNil())

			err = DeleteScoutLogs(db, s.UUID)
			Ω(err).Should(BeNil())

			n, err := NumScoutLogs(db)
			Ω(err).Should(BeNil())
			Ω(n).Should(Equal(int64(0)))
		})
	})
})
