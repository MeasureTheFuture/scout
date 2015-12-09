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
)

func TestConfiguration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configuration Suite")
}

var _ = Describe("Configuration", func() {
	Context("Parsing", func() {
		It("should throw an error for an invalid config file", func() {
			c, err := parseConfiguration("foo")
			Ω(err).ShouldNot(BeNil())

			Ω(c.MinArea).Should(Equal(14000.0))
			Ω(c.DilationIterations).Should(Equal(10))
			Ω(c.ForegroundThresh).Should(Equal(128))
			Ω(c.GaussianSmooth).Should(Equal(5))
			Ω(c.MogHistoryLength).Should(Equal(500))
			Ω(c.MogThreshold).Should(Equal(30))
			Ω(c.MogDetectShadows).Should(Equal(1))
		})

		It("should be able to parse a valid config file", func() {
			c, err := parseConfiguration("testdata/test-config.json")
			Ω(err).Should(BeNil())

			Ω(c.MinArea).Should(Equal(2.0))
			Ω(c.DilationIterations).Should(Equal(2))
			Ω(c.ForegroundThresh).Should(Equal(2))
			Ω(c.GaussianSmooth).Should(Equal(2))
			Ω(c.MogHistoryLength).Should(Equal(2))
			Ω(c.MogThreshold).Should(Equal(2))
			Ω(c.MogDetectShadows).Should(Equal(0))
		})
	})
})
