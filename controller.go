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
	"log"
	"net/http"
)

type Command int

const (
	CALIBRATE Command = iota
	START_MEASURE
	STOP_MEASURE
)

func controller(deltaC chan Command, config Configuration) {
	http.HandleFunc("/calibrate", func(w http.ResponseWriter, r *http.Request) {
		deltaC <- CALIBRATE
	})

	http.HandleFunc("/measure/start", func(w http.ResponseWriter, r *http.Request) {
		deltaC <- START_MEASURE
	})

	http.HandleFunc("/measure/stop", func(w http.ResponseWriter, r *http.Request) {
		deltaC <- STOP_MEASURE
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
