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
	"database/sql/driver"
	"errors"
	"github.com/MeasureTheFuture/scout/configuration"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
	"time"
)

type RealArray []float32
type Path [][2]int

type ScoutInteraction struct {
	Id             int64
	ScoutUUID      string
	Duration       float32
	Waypoints      Path
	WaypointWidths Path
	WaypointTimes  RealArray
	Processed      bool
	EnteredAt      time.Time
}

func (a *RealArray) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Unable to deserialise RealArray")
	}

	asString := string(asBytes)
	asString = asString[1 : len(asString)-1]
	elements := strings.Split(asString, ",")

	res := make([]float32, len(elements))
	for i, v := range elements {
		t, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return err
		}

		res[i] = float32(t)
	}

	*a = res
	return nil
}

func (p *Path) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Unable to deserialise Path")
	}

	asString := string(asBytes)
	asString = asString[2 : len(asString)-2]
	elements := strings.Split(asString, "),(")

	res := make([][2]int, len(elements))
	for i, v := range elements {
		wp := strings.Split(v, ",")
		wpv, err := strconv.ParseInt(wp[0], 10, 32)
		if err != nil {
			return err
		}
		res[i][0] = int(wpv)

		wpv, err = strconv.ParseInt(wp[1], 10, 32)
		if err != nil {
			return err
		}
		res[i][1] = int(wpv)
	}

	*p = res
	return nil
}

func (p Path) Value() (driver.Value, error) {
	res := "["
	for i, v := range p {
		res = res + "(" + strconv.FormatInt(int64(v[0]), 10) + "," + strconv.FormatInt(int64(v[1]), 10) + ")"

		if i < len(p)-1 {
			res = res + ","
		}
	}
	res = res + "]"

	return res, nil
}

func (a RealArray) Value() (driver.Value, error) {
	res := "{"
	for i, v := range a {
		res = res + strconv.FormatFloat(float64(v), 'f', -1, 32)
		if i < len(a)-1 {
			res = res + ","
		}
	}
	res = res + "}"

	return res, nil
}

func CreateScoutInteraction(i *Interaction) ScoutInteraction {
	var result ScoutInteraction

	result.Id = -1
	result.ScoutUUID = ""
	result.Duration = i.Duration

	wpLength := len(i.Path)
	result.Waypoints = make([][2]int, wpLength)
	result.WaypointWidths = make([][2]int, wpLength)
	result.WaypointTimes = make([]float32, wpLength)

	for i, wp := range i.Path {
		result.Waypoints[i] = [2]int{wp.XPixels, wp.YPixels}
		result.WaypointWidths[i] = [2]int{wp.HalfWidthPixels, wp.HalfHeightPixels}
		result.WaypointTimes[i] = wp.T
	}

	result.Processed = false
	result.EnteredAt = i.Entered

	return result
}

func GetScoutInteractionById(db *sql.DB, id int64) (*ScoutInteraction, error) {
	const query = `SELECT id, scout_uuid, duration, waypoints, waypoint_widths, waypoint_times,
	processed, entered_at FROM scout_interactions WHERE id = $1`

	var result ScoutInteraction
	var et time.Time
	err := db.QueryRow(query, id).Scan(&result.Id, &result.ScoutUUID, &result.Duration,
		&result.Waypoints, &result.WaypointWidths, &result.WaypointTimes,
		&result.Processed, &et)
	result.EnteredAt = et.UTC()

	return &result, err
}

func GetLastScoutInteraction(db *sql.DB, scoutUUID string) (*ScoutInteraction, error) {
	const query = `SELECT id, duration, waypoints, waypoint_widths, waypoint_times, processed, entered_at
		FROM scout_interactions WHERE scout_uuid = $1 ORDER BY id DESC LIMIT 1`

	var result ScoutInteraction
	err := db.QueryRow(query, scoutUUID).Scan(&result.Id, &result.Duration, &result.Waypoints, &result.WaypointWidths,
		&result.WaypointTimes, &result.Processed, &result.EnteredAt)
	result.ScoutUUID = scoutUUID

	return &result, err
}

func (si *ScoutInteraction) MarkProcessed(db *sql.DB) error {
	const query = `UPDATE scout_interactions SET processed = true WHERE id = $1`
	_, err := db.Exec(query, si.Id)
	si.Processed = true

	return err
}

func GetUnprocessed(db *sql.DB) ([]*ScoutInteraction, error) {
	const query = `SELECT * FROM scout_interactions WHERE processed = false`
	var result []*ScoutInteraction

	rows, err := db.Query(query)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var si ScoutInteraction
		var et time.Time
		err = rows.Scan(&si.Id, &si.Duration, &si.Waypoints, &si.WaypointWidths,
			&si.WaypointTimes, &si.Processed, &et, &si.ScoutUUID)
		si.EnteredAt = et.UTC()
		if err != nil {
			return result, err
		}
		result = append(result, &si)
	}

	return result, rows.Err()
}

func NumScoutInteractions(db *sql.DB) (int64, error) {
	const query = `SELECT COUNT(*) FROM scout_interactions`
	var result int64
	err := db.QueryRow(query).Scan(&result)

	return result, err
}

func DeleteScoutInteractions(db *sql.DB, scoutUUID string) error {
	const query = `DELETE FROM scout_interactions WHERE scout_uuid = $1`
	_, err := db.Exec(query, scoutUUID)
	return err
}

func (si *ScoutInteraction) Insert(db *sql.DB) error {
	const query = `INSERT INTO scout_interactions (scout_uuid, duration, waypoints,
		waypoint_widths, waypoint_times, processed, entered_at) VALUES
		($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return db.QueryRow(query, si.ScoutUUID, si.Duration, si.Waypoints, si.WaypointWidths,
		si.WaypointTimes, si.Processed, si.EnteredAt).Scan(&si.Id)
}

func ScoutInteractionsAsJSON(db *sql.DB) (string, error) {
	file := configuration.GetDataDir() + "/scout_interactions.json"

	const query = `SELECT * FROM scout_interactions`
	rows, err := db.Query(query)
	if err == sql.ErrNoRows {
		return file, nil
	} else if err != nil {
		return file, err
	}
	defer rows.Close()

	var result []ScoutInteraction
	for rows.Next() {
		var si ScoutInteraction
		err = rows.Scan(&si.Id, &si.Duration, &si.Waypoints, &si.WaypointWidths,
			&si.WaypointTimes, &si.Processed, &si.EnteredAt, &si.ScoutUUID)
		if err != nil {
			return file, err
		}

		result = append(result, si)
	}

	return file, configuration.SaveAsJSON(result, file)
}
