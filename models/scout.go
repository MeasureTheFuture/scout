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
	"fmt"
	"github.com/MeasureTheFuture/mothership/configuration"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
)

type ScoutState string

const (
	IDLE        ScoutState = "idle"
	CALIBRATING ScoutState = "calibrating"
	CALIBRATED  ScoutState = "calibrated"
	MEASURING   ScoutState = "measuring"
)

func (s *ScoutState) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Unable to deserialise ScoutState")
	}

	*s = ScoutState(string(asBytes))
	return nil
}

func (s ScoutState) Value() (driver.Value, error) {
	return string(s), nil
}

type Scout struct {
	UUID       string        `json:"uuid"`
	IpAddress  string        `json:"ip_address"`
	Port       int64         `json:"port"`
	Authorised bool          `json:"authorised"`
	Name       string        `json:"name"`
	State      ScoutState    `json:"state"`
	Summary    *ScoutSummary `json:"summary"`

	MinArea            float64
	DilationIterations int64
	ForegroundThresh   int64
	GaussianSmooth     int64
	MogHistoryLength   int64
	MogThreshold       float64
	MogDetectShadows   int64
	SimplifyEpsilon    float64
	MinDuration        float32
	IdleDuration       float32
	ResumeSqDistance   int64
	MaxArea            float64
}

func GetScoutByUUID(db *sql.DB, uuid string) (*Scout, error) {
	const query = `SELECT ip_address, port, authorised, name, state, min_area,
				   dilation_iterations, foreground_thresh, guassian_smooth,
				   mog_history_length, mog_threshold, mog_detect_shadows,
				   simplify_epsilon, min_duration, idle_duration, resume_sq_distance, max_area
				   FROM scouts WHERE uuid = $1`
	var result Scout
	err := db.QueryRow(query, uuid).Scan(&result.IpAddress, &result.Port, &result.Authorised,
		&result.Name, &result.State, &result.MinArea, &result.DilationIterations,
		&result.ForegroundThresh, &result.GaussianSmooth, &result.MogHistoryLength,
		&result.MogThreshold, &result.MogDetectShadows, &result.SimplifyEpsilon,
		&result.MinDuration, &result.IdleDuration, &result.ResumeSqDistance, &result.MaxArea)
	result.UUID = uuid
	result.Summary, err = GetScoutSummaryByUUID(db, result.UUID)
	if err != nil {
		return &result, err
	}

	return &result, err
}

func GetScout(db *sql.DB) *Scout {
	const query = `SELECT uuid, ip_address, port, authorised, name, state, min_area,
				   dilation_iterations, foreground_thresh, guassian_smooth,
				   mog_history_length, mog_threshold, mog_detect_shadows,
				   simplify_epsilon, min_duration, idle_duration, resume_sq_distance, max_area
				   FROM scouts LIMIT 1`
	var result Scout
	err := db.QueryRow(query).Scan(&result.UUID, &result.IpAddress, &result.Port, &result.Authorised,
		&result.Name, &result.State, &result.MinArea, &result.DilationIterations,
		&result.ForegroundThresh, &result.GaussianSmooth, &result.MogHistoryLength,
		&result.MogThreshold, &result.MogDetectShadows, &result.SimplifyEpsilon,
		&result.MinDuration, &result.IdleDuration, &result.ResumeSqDistance, &result.MaxArea)
	if err != nil {
		log.Fatalf("Unable to get scout %v", err)
	}

	return &result
}

func GetScoutUUID(db *sql.DB) string {
	const query = `SELECT uuid FROM scouts LIMIT 1`
	var result string
	err := db.QueryRow(query).Scan(&result)
	if err != nil {
		log.Fatalf("Unable get scout UUID: %v", err)
	}

	return result
}

func GetAllScouts(db *sql.DB) ([]*Scout, error) {
	const query = `SELECT uuid, ip_address, port, authorised, name, state, min_area,
				   dilation_iterations, foreground_thresh, guassian_smooth,
				   mog_history_length, mog_threshold, mog_detect_shadows,
				   simplify_epsilon, min_duration, idle_duration, resume_sq_distance,
				   max_area FROM scouts`

	var result []*Scout
	rows, err := db.Query(query)
	if err == sql.ErrNoRows {
		return result, nil
	} else if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var s Scout
		err = rows.Scan(&s.UUID, &s.IpAddress, &s.Port, &s.Authorised, &s.Name, &s.State,
			&s.MinArea, &s.DilationIterations, &s.ForegroundThresh,
			&s.GaussianSmooth, &s.MogHistoryLength, &s.MogThreshold,
			&s.MogDetectShadows, &s.SimplifyEpsilon, &s.MinDuration,
			&s.IdleDuration, &s.ResumeSqDistance, &s.MaxArea)
		if err != nil {
			return result, err
		}
		s.Summary, err = GetScoutSummaryByUUID(db, s.UUID)
		if err != nil {
			return result, err
		}

		result = append(result, &s)
	}

	return result, nil
}

func NumScouts(db *sql.DB) (int64, error) {
	const query = `SELECT COUNT(*) FROM scouts`
	var result int64
	err := db.QueryRow(query).Scan(&result)

	return result, err
}

func (s *Scout) ClearCalibrationFrame(db *sql.DB) error {
	const query = `UPDATE scouts SET calibration_frame = NULL WHERE uuid = $1`
	_, err := db.Exec(query, s.UUID)
	return err
}

func (s *Scout) UpdateCalibrationFrame(db *sql.DB, frame []byte) error {
	const query = `UPDATE scouts SET calibration_frame = $1 WHERE uuid = $2`
	_, err := db.Exec(query, frame, s.UUID)
	return err
}

func (s *Scout) GetCalibrationFrame(db *sql.DB) ([]byte, error) {
	const query = `SELECT calibration_frame FROM scouts WHERE uuid = $1`
	var result []byte
	err := db.QueryRow(query, s.UUID).Scan(&result)

	return result, err
}

func (s *Scout) Insert(db *sql.DB) error {
	const query = `INSERT INTO scouts (ip_address, port, authorised, name, state, min_area,
				   dilation_iterations, foreground_thresh, guassian_smooth,
				   mog_history_length, mog_threshold, mog_detect_shadows,
				   simplify_epsilon, min_duration, idle_duration, resume_sq_distance, max_area)
				   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING uuid`
	err := db.QueryRow(query, s.IpAddress, s.Port, s.Authorised, s.Name, s.State,
		s.MinArea, s.DilationIterations, s.ForegroundThresh,
		s.GaussianSmooth, s.MogHistoryLength, s.MogThreshold,
		s.MogDetectShadows, s.SimplifyEpsilon, s.MinDuration,
		s.IdleDuration, s.ResumeSqDistance, s.MaxArea).Scan(&s.UUID)
	if err != nil {
		return err
	}

	// Create matching empty summary.
	s.Summary = &ScoutSummary{s.UUID, 0, Buckets{}, IntBuckets{}}
	return s.Summary.Insert(db)
}

func (s *Scout) Update(db *sql.DB) error {
	const query = `UPDATE scouts SET ip_address = $1, port = $2, authorised = $3, name = $4,
				   state = $5, min_area = $6, dilation_iterations = $7,
				   foreground_thresh = $8, guassian_smooth = $9, mog_history_length = $10,
				   mog_threshold = $11, mog_detect_shadows = $12, simplify_epsilon = $13,
				   min_duration = $14, idle_duration = $15, resume_sq_distance = $16,
				   max_area = $17 WHERE uuid = $18`
	_, err := db.Exec(query, s.IpAddress, s.Port, s.Authorised, s.Name, s.State,
		s.MinArea, s.DilationIterations, s.ForegroundThresh,
		s.GaussianSmooth, s.MogHistoryLength, s.MogThreshold,
		s.MogDetectShadows, s.SimplifyEpsilon, s.MinDuration,
		s.IdleDuration, s.ResumeSqDistance, s.MaxArea, s.UUID)
	return err
}

func ScoutsAsJSON(db *sql.DB) ([]string, error) {
	var files []string
	file := configuration.GetDataDir() + "/scouts.json"

	const query = `SELECT * FROM scouts`
	rows, err := db.Query(query)
	if err == sql.ErrNoRows {
		return files, nil
	} else if err != nil {
		return files, err
	}
	defer rows.Close()

	var result []Scout
	var image []byte

	for rows.Next() {
		var s Scout
		err = rows.Scan(&s.UUID, &s.IpAddress, &s.Authorised, &image, &s.Name, &s.State, &s.Port,
			&s.MinArea, &s.DilationIterations, &s.ForegroundThresh, &s.GaussianSmooth,
			&s.MogHistoryLength, &s.MogThreshold, &s.MogDetectShadows, &s.SimplifyEpsilon,
			&s.MinDuration, &s.IdleDuration, &s.ResumeSqDistance, &s.MaxArea)
		if err != nil {
			return files, err
		}

		// Write image.
		if len(image) > 0 {
			imgF := configuration.GetDataDir() + "/scout-" + fmt.Sprintf("%s", s.UUID) + ".jpg"
			err = ioutil.WriteFile(imgF, image, 0644)
			if err != nil {
				return files, err
			}

			files = append(files, imgF)
		}

		result = append(result, s)
	}

	err = configuration.SaveAsJSON(result, file)
	if err != nil {
		return files, err
	}

	return append(files, file), nil
}
