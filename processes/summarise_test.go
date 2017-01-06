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

package processes

import (
	"database/sql"
	"github.com/MeasureTheFuture/mothership/configuration"
	"github.com/MeasureTheFuture/mothership/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
	"time"
)

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

func TestSummarise(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "summarise process Suite")
}

var _ = Describe("Summarise Process", func() {
	BeforeSuite(func() {
		config, err := configuration.Parse(os.Getenv("GOPATH") + "/mothership.json")
		Ω(err).Should(BeNil())
		db, err = sql.Open("postgres", "user="+config.DBUserName+" dbname="+config.DBTestName)
		Ω(err).Should(BeNil())
	})

	AfterEach(cleaner)
	AfterSuite(cleaner)

	Context("updateUnprocessed", func() {
		It("should ignore proccessed interactions", func() {
			s := models.Scout{-1, "59ef7180-f6b2-4129-99bf-970eb4312b4b",
				"192.168.0.1", 8080, true, "foo", "calibrating", &models.ScoutSummary{}}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			et := time.Now().UTC().Round(15 * time.Minute)
			si := models.ScoutInteraction{-1, s.Id, 0.2, models.Path{[2]int{1, 2}},
				models.Path{[2]int{3, 4}}, models.RealArray{0.1}, true, et}
			err = si.Insert(db)
			Ω(err).Should(BeNil())

			updateUnprocessed(db)
			ss, err := models.GetScoutSummaryById(db, s.Id)
			Ω(err).Should(BeNil())
			Ω(ss.VisitorCount).Should(Equal(int64(0)))
		})

		It("should increment the visitor count", func() {
			s := models.Scout{-1, "59ef7180-f6b2-4129-99bf-970eb4312b4b",
				"192.168.0.1", 8080, true, "foo", "calibrating", &models.ScoutSummary{}}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			et := time.Now().UTC().Round(15 * time.Minute)
			si := &models.ScoutInteraction{-1, s.Id, 0.2, models.Path{[2]int{1, 2}},
				models.Path{[2]int{3, 4}}, models.RealArray{0.1}, false, et}
			err = si.Insert(db)
			Ω(err).Should(BeNil())

			si2 := &models.ScoutInteraction{-1, s.Id, 0.2, models.Path{[2]int{1, 2}},
				models.Path{[2]int{3, 4}}, models.RealArray{0.1}, false, et}
			err = si2.Insert(db)
			Ω(err).Should(BeNil())

			updateUnprocessed(db)
			ss, err := models.GetScoutSummaryById(db, s.Id)
			Ω(err).Should(BeNil())
			Ω(ss.VisitorCount).Should(Equal(int64(2)))

			si, err = models.GetScoutInteractionById(db, si.Id)
			Ω(err).Should(BeNil())
			Ω(si.Processed).Should(BeTrue())

			si2, err = models.GetScoutInteractionById(db, si2.Id)
			Ω(err).Should(BeNil())
			Ω(si2.Processed).Should(BeTrue())
		})
	})

	Context("maxTravelTime", func() {
		It("should return the max travel time for a bucket", func() {
			wpA := models.Waypoint{0, 0, 10, 10, 0.0}
			wpB := models.Waypoint{0, 192, 10, 10, 1.0}
			wpC := models.Waypoint{0, 25, 10, 10, 1.0}

			Ω(maxTravelTime(wpA, wpB)).Should(Equal(float32(0.5)))
			Ω(maxTravelTime(wpA, wpC)).Should(Equal(float32(1.0)))
		})
	})

	Context("updateTimeBuckets", func() {
		PIt("it should update the travel times for the buckets in a scout summary", func() {
			ss := &models.ScoutSummary{}
			s := models.Scout{-1, "59ef7180-f6b2-4129-99bf-970eb4312b4b",
				"192.168.0.1", 8080, true, "foo", "calibrating", ss}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			et := time.Now().UTC().Round(15 * time.Minute)
			si := &models.ScoutInteraction{-1, s.Id, 0.2, models.Path{[2]int{1, 2}},
				models.Path{[2]int{0, 0}, [2]int{0, 25}}, models.RealArray{0.0, 1.0}, false, et}
			err = si.Insert(db)
			Ω(err).Should(BeNil())

			updateTimeBuckets(db, ss, si)

			ssA, err := models.GetScoutSummaryById(db, s.Id)
			Ω(err).Should(BeNil())
			Ω(ssA.VisitorCount).Should(Equal(int64(1)))

			tBuckets := models.Buckets{}
			tBuckets[0][0] = 1.0
			vBuckets := models.IntBuckets{}
			vBuckets[0][0] = 1

			Ω(ssA.VisitTimeBuckets).Should(Equal(tBuckets))
			Ω(ssA.VisitorBuckets).Should(Equal(vBuckets))
		})
	})
})
