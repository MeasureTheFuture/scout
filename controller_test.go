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
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func createTestServer() (dC chan Command, dCFG chan Configuration,
	cFile string, c Configuration, s *httptest.Server) {
	dC = make(chan Command, 1)
	dCFG = make(chan Configuration, 1)
	cFile = "testdata/web.json"
	os.Remove(cFile) // Make sure we have no configuration file on disk,
	// so a new default is generated each time.

	c, _ = parseConfiguration(cFile)

	s = httptest.NewServer(bindHandlers(dC, dCFG, cFile, c))

	return
}

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = Describe("Controller", func() {
	Context("Inbound commands", func() {
		It("be able to start measuring", func() {
			dC, _, _, _, s := createTestServer()

			req, err := http.NewRequest("GET", fmt.Sprintf("%s/measure/start", s.URL), nil)
			res, err := http.DefaultClient.Do(req)

			Ω(err).Should(BeNil())
			Ω(res.StatusCode).Should(Equal(200))
			c := <-dC
			Ω(c).Should(Equal(START_MEASURE))

			s.Close()
		})

		It("be able to stop measuring", func() {
			dC, _, _, _, s := createTestServer()

			req, err := http.NewRequest("GET", fmt.Sprintf("%s/measure/stop", s.URL), nil)
			res, err := http.DefaultClient.Do(req)

			Ω(err).Should(BeNil())
			Ω(res.StatusCode).Should(Equal(200))
			c := <-dC
			Ω(c).Should(Equal(STOP_MEASURE))

			s.Close()
		})

		It("be able to calibrate", func() {
			dC, dCFG, cFile, _, s := createTestServer()

			req, err := http.NewRequest("GET", fmt.Sprintf("%s/calibrate?MinArea=2.0&DilationIterations=2", s.URL), nil)
			res, err := http.DefaultClient.Do(req)

			Ω(err).Should(BeNil())
			Ω(res.StatusCode).Should(Equal(200))
			com := <-dC
			Ω(com).Should(Equal(CALIBRATE))

			cfg := <-dCFG
			Ω(cfg.MinArea).Should(BeNumerically("==", 2))
			Ω(cfg.DilationIterations).Should(Equal(2))

			newCFG, err := parseConfiguration(cFile)
			Ω(err).Should(BeNil())

			Ω(newCFG.MinArea).Should(BeNumerically("==", 2))
			Ω(newCFG.DilationIterations).Should(Equal(2))

			s.Close()
		})
	})
})
