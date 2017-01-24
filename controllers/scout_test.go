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

package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/MeasureTheFuture/scout/configuration"
	"github.com/MeasureTheFuture/scout/models"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
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

func TestScout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scout controller Suite")
}

var _ = Describe("Scout controller", func() {

	BeforeSuite(func() {
		config, err := configuration.Parse(os.Getenv("GOPATH") + "/mothership.json")
		Ω(err).Should(BeNil())
		db, err = sql.Open("postgres", "user="+config.DBUserName+" dbname="+config.DBTestName)
		Ω(err).Should(BeNil())
	})

	AfterEach(cleaner)
	AfterSuite(cleaner)

	Context("GetScouts", func() {
		It("should not generate an error if their are no detected scouts", func() {
			e := echo.New()
			req, err := http.NewRequest(echo.GET, "/scouts", strings.NewReader(""))
			Ω(err).Should(BeNil())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = GetScouts(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(200))

			var sl []models.Scout
			err = json.Unmarshal(rec.Body.Bytes(), &sl)
			Ω(err).Should(BeNil())

			Ω(len(sl)).Should(Equal(0))
		})

		It("should return a list of all the attached scouts", func() {
			s := models.Scout{-1, "59ef7180-f6b2-4129-99bf-970eb4312b4b", "192.168.0.1",
				8080, true, "foo", "calibrating", &models.ScoutSummary{}}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			s2 := models.Scout{-1, "eeef7180-f6b2-4129-99bf-970eb4312b4b", "192.168.0.2",
				8080, true, "foop", "calibrating", &models.ScoutSummary{}}
			err = s2.Insert(db)
			Ω(err).Should(BeNil())

			e := echo.New()
			req, err := http.NewRequest(echo.GET, "/scouts", strings.NewReader(""))
			Ω(err).Should(BeNil())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = GetScouts(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(200))

			var sl []models.Scout
			err = json.Unmarshal(rec.Body.Bytes(), &sl)
			Ω(err).Should(BeNil())

			Ω(len(sl)).Should(Equal(2))
			Ω(sl[0]).Should(Equal(s))
			Ω(sl[1]).Should(Equal(s2))
		})

		It("should return a single scout", func() {
			s := models.Scout{-1, "59ef7180-f6b2-4129-99bf-970eb4312b4b", "192.168.0.1",
				8080, true, "foo", "calibrated", &models.ScoutSummary{}}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			e := echo.New()
			req, err := http.NewRequest(echo.GET, "/scouts/", strings.NewReader(""))
			Ω(err).Should(BeNil())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/scouts/:id")
			c.SetParamNames("id")
			c.SetParamValues(strconv.FormatInt(s.Id, 10))

			err = GetScout(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(200))

			var ns models.Scout
			err = json.Unmarshal(rec.Body.Bytes(), &ns)
			Ω(err).Should(BeNil())
			Ω(ns).Should(Equal(s))
		})

		It("should be able to update a single scout", func() {
			s := models.Scout{-1, "59ef7180-f6b2-4129-99bf-970eb4312b4b", "192.168.0.1",
				8080, true, "foo", "calibrated", &models.ScoutSummary{}}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			s.IpAddress = "192.168.0.5"

			e := echo.New()
			b, err := json.Marshal(s)
			Ω(err).Should(BeNil())
			req, err := http.NewRequest(echo.PUT, "/scouts/", bytes.NewReader(b))
			Ω(err).Should(BeNil())

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/scouts/:id")
			c.SetParamNames("id")
			c.SetParamValues(strconv.FormatInt(s.Id, 10))

			err = UpdateScout(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(200))

			ns, err := models.GetScoutById(db, s.Id)
			Ω(err).Should(BeNil())
			Ω(ns).Should(Equal(&s))
		})
	})
})
