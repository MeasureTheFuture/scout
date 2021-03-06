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
	"database/sql"
	"github.com/MeasureTheFuture/scout/configuration"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

func TestScout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scout Suite")
}

var (
	db  *sql.DB
	err error
)

func cleaner() {
	_, err := db.Exec(`DELETE FROM scout_summaries`)
	Ω(err).Should(BeNil())

	_, err = db.Exec(`DELETE FROM scout_interactions`)
	Ω(err).Should(BeNil())

	_, err = db.Exec(`DELETE FROM scout_logs`)
	Ω(err).Should(BeNil())

	_, err = db.Exec(`DELETE FROM scout_healths`)
	Ω(err).Should(BeNil())

	_, err = db.Exec(`DELETE FROM scouts`)
	Ω(err).Should(BeNil())
}

var _ = Describe("Scout Model", func() {

	BeforeSuite(func() {
		config, err := configuration.Parse(os.Getenv("GOPATH") + "/scout.json")
		Ω(err).Should(BeNil())
		db, err = sql.Open("postgres", "user="+config.DBUserName+" dbname="+config.DBTestName)
		Ω(err).Should(BeNil())
	})

	AfterEach(cleaner)
	AfterSuite(cleaner)

	Context("Insert", func() {
		It("should insert a valid scout into the DB.", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "calibrated", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1, 4.0}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			s2, err := GetScoutByUUID(db, s.UUID)
			Ω(err).Should(BeNil())
			Ω(&s).Should(Equal(s2))
		})

		It("should return an error when an invalid scout is inserted into the DB.", func() {
			s := Scout{"aa", "192.168.0.1", 8080, true, "foo", "calibratingas", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1, 4.0}
			err := s.Insert(db)
			Ω(err).ShouldNot(BeNil())
			Ω(s.UUID).Should(Equal("aa"))
		})
	})

	Context("Get", func() {
		It("should be able to get all scouts", func() {
			al, err := GetAllScouts(db)
			Ω(err).Should(BeNil())
			Ω(len(al)).Should(Equal(0))

			s1 := Scout{"", "192.168.0.1", 8080, true, "foo", "calibrated", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1, 4.0}
			err = s1.Insert(db)
			Ω(err).Should(BeNil())

			s2 := Scout{"", "192.168.0.2", 8080, true, "foo", "calibrated", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1, 4.0}
			err = s2.Insert(db)
			Ω(err).Should(BeNil())

			al, err = GetAllScouts(db)
			Ω(err).Should(BeNil())
			Ω(len(al)).Should(Equal(2))
			Ω(al).Should(Equal([]*Scout{&s1, &s2}))
		})
	})

	Context("Update", func() {
		It("should be able to update a scout in the DB", func() {
			s := Scout{"", "192.168.0.1", 8080, true, "foo", "measuring", &ScoutSummary{},
				2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1, 4.0}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			s.IpAddress = "192.168.0.2"
			err = s.Update(db)
			Ω(err).Should(BeNil())
			s2, err := GetScoutByUUID(db, s.UUID)
			Ω(err).Should(BeNil())
			Ω(&s).Should(Equal(s2))
		})
	})
})
