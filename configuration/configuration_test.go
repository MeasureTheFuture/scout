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

package configuration

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
			_, err := Parse("foo")
			Ω(err).ShouldNot(BeNil())
		})

		It("should be able to parse a valid config file", func() {
			c, err := Parse("../scout.json_example")
			Ω(err).Should(BeNil())

			Ω(c.DBUserName).Should(Equal("mtf"))
			Ω(c.DBName).Should(Equal("mothership"))
			Ω(c.DBTestName).Should(Equal("mothership_test"))
			Ω(c.Address).Should(Equal(":80"))
			Ω(c.StaticAssets).Should(Equal("public"))

			Ω(c.MinArea).Should(Equal(2.0))
			Ω(c.DilationIterations).Should(Equal(2))
			Ω(c.ForegroundThresh).Should(Equal(2))
			Ω(c.GaussianSmooth).Should(Equal(2))
			Ω(c.MogHistoryLength).Should(Equal(2))
			Ω(c.MogThreshold).Should(Equal(2.0))
			Ω(c.MogDetectShadows).Should(Equal(0))

			Ω(c.ScoutAddress).Should(Equal(":9090"))
			Ω(c.MothershipAddress).Should(Equal("127.0.0.1:9091"))
			Ω(c.SimplifyEpsilon).Should(Equal(2.0))
			Ω(c.MinDuration).Should(Equal(float32(0.2)))
			Ω(c.IdleDuration).Should(Equal(float32(0.3)))
			Ω(c.ResumeSqDistance).Should(Equal(1))
		})
	})

	Context("Saving", func() {
		It("should be able to save a config file", func() {
			c := Configuration{"mtf", "", "mothership", "mothership_test", ":80", "public", 1000,
				2.0, 2, 2, 2, 2, 2.0, 0, ":9090", "127.0.0.1:9091",
				"0938c583-4140-458c-b267-a8d816d96f4b", 2.0, 0.2, 0.3, 1}
			SaveAsJSON(c, "../testdata/foo.json")

			a, err := Parse("../scout.json_example")
			b, err := Parse("../testdata/foo.json")

			Ω(err).Should(BeNil())
			Ω(a).Should(Equal(b))
		})
	})
})