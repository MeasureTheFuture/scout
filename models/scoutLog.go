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
	_ "github.com/lib/pq"
	"time"
)

type ScoutLog struct {
	ScoutId   int64
	Log       []byte
	CreatedAt time.Time
}

func GetScoutLogById(db *sql.DB, scoutId int64, time time.Time) (*ScoutLog, error) {
	const query = `SELECT log FROM scout_logs WHERE scout_id = $1 AND created_at = $2`

	var result ScoutLog
	err := db.QueryRow(query, scoutId, time).Scan(&result.Log)
	result.ScoutId = scoutId
	result.CreatedAt = time

	return &result, err
}

func GetLastScoutLog(db *sql.DB, scoutId int64) (*ScoutLog, error) {
	const query = `SELECT log, created_at FROM scout_logs WHERE scout_id = $1 ORDER by created_at DESC LIMIT 1`

	var result ScoutLog
	err := db.QueryRow(query, scoutId).Scan(&result.Log, &result.CreatedAt)
	result.ScoutId = scoutId

	return &result, err
}

func NumScoutLogs(db *sql.DB) (int64, error) {
	const query = `SELECT COUNT(*) FROM scout_logs`
	var result int64
	err := db.QueryRow(query).Scan(&result)

	return result, err
}

func DeleteScoutLogs(db *sql.DB, scoutId int64) error {
	const query = `DELETE FROM scout_logs WHERE scout_id = $1`
	_, err := db.Exec(query, scoutId)

	return err
}

func (s *ScoutLog) Insert(db *sql.DB) error {
	const query = `INSERT INTO scout_logs (scout_id, log, created_at) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, s.ScoutId, s.Log, s.CreatedAt)

	return err
}
