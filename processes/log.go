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

package processes

import (
	"database/sql"
	"github.com/MeasureTheFuture/scout/models"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func SaveLogToDB(tmpLog string, db *sql.DB) {
	// Store the old log in the DB.
	l, err := createLogFromFile(tmpLog, db)
	if err != nil {
		log.Fatalf("ERROR: Unable to create log from file - %s", err)
	}
	err = l.Insert(db)
	if err != nil {
		log.Fatalf("ERROR: Unable to save log to database - %s", err)
	}
	err = os.Remove(tmpLog)
	if err != nil {
		log.Fatalf("ERROR: Unable to delete temporary log - %s", err)
	}
}

func createLogFromFile(tmpLog string, db *sql.DB) (*models.ScoutLog, error) {
	b, err := ioutil.ReadFile(tmpLog)
	if err != nil {
		return nil, err
	}

	sl := models.ScoutLog{models.GetScoutUUID(db), b, time.Now().UTC()}
	return &sl, nil
}
