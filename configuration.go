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
	"encoding/json"
	"os"
)

type Configuration struct {
	// Computer vision parameters.
	MinArea            float64 // The minimum area enclosed by a contour to be counted as an interaction.
	DilationIterations int     // The number of iterations to perform while dilating the foreground mask.
	ForegroundThresh   int     // A value between 0 and 255 to use when thresholding the foreground mask.
	GaussianSmooth     int     // The size of the filter to use when gaussian smoothing.
	MogHistoryLength   int     // The length of history to use for the MOG2 subtractor.
	MogThreshold       float64 // Threshold to use with the MOG2 subtractor.
	MogDetectShadows   int     // 1 if you want the MOG2 subtractor to detect shadows, 0 otherwise.

}

func parseConfiguration(configFile string) (c Configuration, err error) {
	c = Configuration{14000.0, 10, 128, 5, 500, 30, 1}

	// Open the configuration file.
	file, err := os.Open(configFile)
	if err != nil {
		return c, err
	}

	// Parse JSON in the configuration file.
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	return c, err
}