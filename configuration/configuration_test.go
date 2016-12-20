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
			c, err := Parse("../mothership.json_example")
			Ω(err).Should(BeNil())

			Ω(c.DBUserName).Should(Equal("mtf"))
			Ω(c.DBName).Should(Equal("mothership"))
			Ω(c.DBTestName).Should(Equal("mothership_test"))
			Ω(c.Address).Should(Equal(":1323"))
			Ω(c.StaticAssets).Should(Equal("public"))
		})
	})
})
