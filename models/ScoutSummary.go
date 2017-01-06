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
	"github.com/MeasureTheFuture/mothership/configuration"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
)

type Buckets [configuration.HBuckets][configuration.WBuckets]float32
type IntBuckets [configuration.HBuckets][configuration.WBuckets]int

type ScoutSummary struct {
	ScoutId          int64
	VisitorCount     int64
	VisitTimeBuckets Buckets
	VisitorBuckets   IntBuckets
}

func (b IntBuckets) Value() (driver.Value, error) {
	res := "{"
	for i, r := range b {
		res = res + "{"

		for j, v := range r {
			res = res + strconv.Itoa(v)

			if j < len(r)-1 {
				res = res + ","
			}
		}

		res = res + "}"
		if i < len(r)-1 {
			res = res + ","
		}
	}
	res = res + "}"

	return res, nil
}

func (b Buckets) Value() (driver.Value, error) {
	res := "{"
	for i, r := range b {
		res = res + "{"

		for j, v := range r {
			res = res + strconv.FormatFloat(float64(v), 'f', -1, 32)

			if j < len(r)-1 {
				res = res + ","
			}
		}

		res = res + "}"
		if i < len(r)-1 {
			res = res + ","
		}
	}
	res = res + "}"

	return res, nil
}

func (b *IntBuckets) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Unable to deserialise IntBuckets")
	}

	asString := string(asBytes)
	asString = asString[2 : len(asString)-2]
	elements := strings.Split(asString, "},{")

	var res IntBuckets
	for i, r := range elements {
		wp := strings.Split(r, ",")
		for j, v := range wp {
			bv, err := strconv.Atoi(v)
			if err != nil {
				return err
			}

			res[i][j] = bv
		}
	}

	*b = res
	return nil
}

func (b *Buckets) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Unable to deserialise Buckets")
	}

	asString := string(asBytes)
	asString = asString[2 : len(asString)-2]
	elements := strings.Split(asString, "},{")

	var res Buckets
	for i, r := range elements {
		wp := strings.Split(r, ",")
		for j, v := range wp {
			bv, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return err
			}

			res[i][j] = float32(bv)
		}
	}

	*b = res
	return nil
}

func GetScoutSummaryById(db *sql.DB, scoutId int64) (*ScoutSummary, error) {
	const query = `SELECT visitor_count, visit_time_buckets, visitor_buckets FROM scout_summaries WHERE scout_id = $1`

	var result ScoutSummary
	err := db.QueryRow(query, scoutId).Scan(&result.VisitorCount, &result.VisitTimeBuckets, &result.VisitorBuckets)
	result.ScoutId = scoutId

	return &result, err
}

func (si *ScoutSummary) Clear(db *sql.DB) error {
	si.VisitorCount = 0
	si.VisitTimeBuckets = Buckets{}
	si.VisitorBuckets = IntBuckets{}

	return si.Update(db)
}

func (si *ScoutSummary) Insert(db *sql.DB) error {
	const query = `INSERT INTO scout_summaries (scout_id, visitor_count, visit_time_buckets, visitor_buckets) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, si.ScoutId, si.VisitorCount, si.VisitTimeBuckets, si.VisitorBuckets)

	return err
}

func (si *ScoutSummary) Update(db *sql.DB) error {
	const query = `UPDATE scout_summaries SET visitor_count = $1, visit_time_buckets = $2, visitor_buckets = $3 WHERE scout_id = $4`
	_, err := db.Exec(query, si.VisitorCount, si.VisitTimeBuckets, si.VisitorBuckets, si.ScoutId)

	return err
}

func ScoutSummariesAsJSON(db *sql.DB) (string, error) {
	file := configuration.GetDataDir() + "/scout_summaries.json"

	const query = `SELECT * FROM scout_summaries`
	rows, err := db.Query(query)
	if err == sql.ErrNoRows {
		return file, nil
	} else if err != nil {
		return file, err
	}
	defer rows.Close()

	var result []ScoutSummary
	for rows.Next() {
		var ss ScoutSummary
		err = rows.Scan(&ss.ScoutId, &ss.VisitorCount, &ss.VisitTimeBuckets, &ss.VisitorBuckets)
		if err != nil {
			return file, err
		}

		result = append(result, ss)
	}

	return file, configuration.SaveAsJSON(result, file)
}
