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
	"flag"
	"log"
)

func main() {
	log.Printf("INFO: Starting sensor.\n")
	var configFile string
	var videoFile string

	flag.StringVar(&configFile, "configFile", "scout.json", "The path to the configuration file")
	flag.StringVar(&videoFile, "videoFile", "", "The path to a video file to detect motion from instead of a webcam")
	flag.Parse()

	config, err := parseConfiguration(configFile)
	if err != nil {
		log.Printf("INFO: Unable to open '%s', using defaults.", configFile)
		log.Printf("INFO: %s", err)
	}

	monitor(config, videoFile)
}
