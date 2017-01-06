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
