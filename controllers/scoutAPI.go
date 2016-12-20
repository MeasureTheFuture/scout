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

package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/MeasureTheFuture/scout/models"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func isScoutAuthorised(db *sql.DB, c echo.Context) (*models.Scout, error) {
	uuid := c.Request().Header.Get("Mothership-Authorization")
	s, err := models.GetScoutByUUID(db, uuid)

	if err != nil {
		c, err := models.NumScouts(db)
		if err != nil {
			return nil, err
		}

		// Scout doesn't exist, create it and mark it as un-authorized.
		ns := models.Scout{-1, uuid, "0.0.0.0", 8080, false, "Location " + strconv.FormatInt(c+1, 10), "idle", &models.ScoutSummary{}}
		err = ns.Insert(db)
		return nil, err
	}

	if !s.Authorised {
		return nil, nil
	}

	return s, nil
}

func ScoutCalibrated(db *sql.DB, c echo.Context) error {
	s, err := isScoutAuthorised(db, c)
	if err != nil {
		return err
	}
	if s == nil {
		return c.HTML(http.StatusNotFound, "")
	}

	img, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := img.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Store the calibration frame.
	var buff bytes.Buffer
	_, err = buff.ReadFrom(src)

	s.UpdateCalibrationFrame(db, buff.Bytes())
	s.State = models.CALIBRATED
	err = s.Update(db)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "Calibration received")
}

func ScoutInteraction(db *sql.DB, c echo.Context) error {
	s, err := isScoutAuthorised(db, c)
	if err != nil {
		return err
	}
	if s == nil {
		return c.HTML(http.StatusNotFound, "")
	}

	data, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := data.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	var i models.Interaction
	err = json.NewDecoder(src).Decode(&i)
	if err != nil {
		return err
	}

	si := models.CreateScoutInteraction(&i)
	si.ScoutId = s.Id
	err = si.Insert(db)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "Interaction received")
}

func ScoutLog(db *sql.DB, c echo.Context) error {
	s, err := isScoutAuthorised(db, c)
	if err != nil {
		return err
	}
	if s == nil {
		return c.HTML(http.StatusNotFound, "")
	}

	data, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := data.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Store scout log.
	var buff bytes.Buffer
	_, err = buff.ReadFrom(src)
	sl := models.ScoutLog{s.Id, buff.Bytes(), time.Now().UTC()}
	err = sl.Insert(db)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "Log received")
}

type Heartbeat struct {
	UUID    string     // The UUID for the scout.
	Version string     // The version of the protocol used used for transmitting data to the mothership.
	Health  HealthData // The current health status of the scout.

}

type HealthData struct {
	IpAddress   string  // The current IP address of the scout.
	CPU         float32 // The amount of CPU load currently being consumed on the scout. 0.0 - no load, 1.0 - full load.
	Memory      float32 // The amount of memory consumed on the scout. 0.0 - no memory used, 1.0 no memory available.
	TotalMemory float32 // The total number of gigabytes of virtual memory currently available.
	Storage     float32 // The amount of storage consumed on the scout. 0.0 - disk unused, 1.0 disk full.
}

func ScoutHeartbeat(db *sql.DB, c echo.Context) error {
	s, err := isScoutAuthorised(db, c)
	if err != nil {
		return err
	}
	if s == nil {
		return c.HTML(http.StatusNotFound, "")
	}

	// Deserialise the heartbeat from JSON.
	data, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := data.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	var hb Heartbeat
	err = json.NewDecoder(src).Decode(&hb)
	if err != nil {
		return err
	}

	// Create new record for the scout health.
	sh := models.ScoutHealth{s.Id, hb.Health.CPU, hb.Health.Memory, hb.Health.TotalMemory, hb.Health.Storage, time.Now().UTC()}
	err = sh.Insert(db)
	if err != nil {
		return err
	}

	add := strings.Split(hb.Health.IpAddress, ":")
	if len(add) != 2 {
		return errors.New("Invalid Ip Address.")
	}

	// Update IP address of scout.
	s.IpAddress = "0.0.0.0"
	if add[0] != "" {
		s.IpAddress = add[0]
	}

	s.Port = 8080
	if add[1] != "" {
		i, err := strconv.Atoi(add[1])
		if err != nil {
			return err
		}
		s.Port = int64(i)
	}

	err = s.Update(db)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "Heartbeat received")
}
