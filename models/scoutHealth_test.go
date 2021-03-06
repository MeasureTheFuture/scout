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
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"testing"
	"time"
)

func TestScoutHealth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scout Health Suite")
}

var _ = Describe("Scout Health Model", func() {
	AfterEach(cleaner)

	Context("Insert", func() {
		It("should insert a valid scouthealth into the DB.", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1, 4.0}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			t := time.Now()
			sh := ScoutHealth{s.UUID, 0.1, 0.2, 0.3, 0.4, t}
			err = sh.Insert(db)
			Ω(err).Should(BeNil())

			sh2, err := GetScoutHealthByUUID(db, s.UUID, t)
			Ω(err).Should(BeNil())
			Ω(&sh).Should(Equal(sh2))

		})

		It("should return an error when an invalid scout health is inserted into the DB.", func() {
			sh := ScoutHealth{"", 0.1, 0.2, 0.3, 0.4, time.Now().UTC()}
			err := sh.Insert(db)
			Ω(err).ShouldNot(BeNil())
		})
	})

	Context("Delete", func() {
		It("should be able to delete healths for a specified scout", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1, 4.0}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			sh := ScoutHealth{s.UUID, 0.1, 0.2, 0.3, 0.4, time.Now().UTC()}
			err = sh.Insert(db)

			sh2 := ScoutHealth{s.UUID, 0.1, 0.2, 0.3, 0.4, time.Now().UTC()}
			err = sh2.Insert(db)

			err = DeleteScoutHealths(db, s.UUID)
			Ω(err).Should(BeNil())

			n, err := NumScoutHealths(db)
			Ω(err).Should(BeNil())
			Ω(n).Should(Equal(int64(0)))
		})
	})

	Context("Get", func() {
		It("should be able to get scout healths as json", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "idle", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1, 4.0}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			t := time.Now().UTC().Round(time.Second)
			sh := ScoutHealth{s.UUID, 0.1, 0.2, 0.3, 0.4, t}
			err = sh.Insert(db)

			jsonF, err := ScoutHealthsAsJSON(db)
			Ω(err).Should(BeNil())

			jsonB, err := ioutil.ReadFile(jsonF)
			Ω(err).Should(BeNil())

			var result []ScoutHealth
			err = json.Unmarshal(jsonB, &result)
			Ω(err).Should(BeNil())
			Ω(result).Should(Equal([]ScoutHealth{sh}))
		})
	})
})
