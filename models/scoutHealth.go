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
	"github.com/MeasureTheFuture/scout/configuration"
	_ "github.com/lib/pq"
	"time"
)

type ScoutHealth struct {
	ScoutUUID   string
	CPU         float32
	Memory      float32
	TotalMemory float32
	Storage     float32
	CreatedAt   time.Time
}

func GetScoutHealthByUUID(db *sql.DB, scoutUUID string, time time.Time) (*ScoutHealth, error) {
	const query = `SELECT cpu, memory, total_memory, storage FROM scout_healths WHERE scout_uuid = $1 AND created_at = $2`

	var result ScoutHealth
	err := db.QueryRow(query, scoutUUID, time).Scan(&result.CPU, &result.Memory, &result.TotalMemory, &result.Storage)
	result.ScoutUUID = scoutUUID
	result.CreatedAt = time

	return &result, err
}

func GetLastScoutHealth(db *sql.DB, scoutUUID string) (*ScoutHealth, error) {
	const query = `SELECT cpu, memory, total_memory, storage, created_at FROM scout_healths WHERE scout_uuid = $1 ORDER BY created_at DESC LIMIT 1`

	var result ScoutHealth
	err := db.QueryRow(query, scoutUUID).Scan(&result.CPU, &result.Memory, &result.TotalMemory, &result.Storage, &result.CreatedAt)
	result.ScoutUUID = scoutUUID

	return &result, err
}

func DeleteScoutHealths(db *sql.DB, scoutUUID string) error {
	const query = `DELETE FROM scout_healths WHERE scout_uuid = $1`
	_, err := db.Exec(query, scoutUUID)
	return err
}

func NumScoutHealths(db *sql.DB) (int64, error) {
	const query = `SELECT COUNT(*) FROM scout_healths`
	var result int64
	err := db.QueryRow(query).Scan(&result)

	return result, err
}

func (s *ScoutHealth) Insert(db *sql.DB) error {
	const query = `INSERT INTO scout_healths (scout_uuid, cpu, memory, total_memory, storage,
		created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, s.ScoutUUID, s.CPU, s.Memory, s.TotalMemory, s.Storage, s.CreatedAt)
	return err
}

func ScoutHealthsAsJSON(db *sql.DB) (string, error) {
	file := configuration.GetDataDir() + "/scout_healths.json"
	const query = `SELECT * FROM scout_healths`
	rows, err := db.Query(query)
	if err == sql.ErrNoRows {
		return file, nil
	} else if err != nil {
		return file, err
	}
	defer rows.Close()

	var result []ScoutHealth
	for rows.Next() {
		var sh ScoutHealth
		err = rows.Scan(&sh.ScoutUUID, &sh.CPU, &sh.Memory, &sh.TotalMemory, &sh.Storage, &sh.CreatedAt)
		if err != nil {
			return file, err
		}

		result = append(result, sh)
	}

	return file, configuration.SaveAsJSON(result, file)
}
