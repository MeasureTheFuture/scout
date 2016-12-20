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

package main

import (
	"database/sql"
	"flag"
	"github.com/MeasureTheFuture/mothership/configuration"
	"github.com/MeasureTheFuture/mothership/controllers"
	"github.com/MeasureTheFuture/mothership/processes"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "configFile", "mothership.json", "The path to the configuration file")
	flag.Parse()

	f, err := os.OpenFile("mothership.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Unable to open log file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("INFO: Starting mothership")

	config, err := configuration.Parse(configFile)
	if err != nil {
		log.Fatalf("ERROR: Can't parse configuration - %s", err)
	}

	connection := "user=" + config.DBUserName + " dbname=" + config.DBName
	if config.DBPassword != "" {
		connection = connection + " password=" + config.DBPassword
	}

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatalf("ERROR: Can't open database - %s", err)
	}
	defer db.Close()

	// Start background proccesses.
	go processes.Summarise(db, config)

	e := echo.New()
	e.Static("/", config.StaticAssets)
	e.Static("/css", config.StaticAssets+"/css")
	e.Static("/fonts", config.StaticAssets+"/fonts")
	e.Static("/img", config.StaticAssets+"/img")

	// Front-end API for displaying results from the scouts.
	e.GET("/scouts", func(c echo.Context) error {
		return controllers.GetScouts(db, c)
	})

	e.GET("/scouts/:id/frame.jpg", func(c echo.Context) error {
		return controllers.GetScoutFrame(db, c)
	})

	e.GET("/scouts/:id", func(c echo.Context) error {
		return controllers.GetScout(db, c)
	})

	e.PUT("/scouts/:id", func(c echo.Context) error {
		return controllers.UpdateScout(db, c)
	})

	e.GET("/download.zip", func(c echo.Context) error {
		return controllers.DownloadData(db, c)
	})

	// SCOUT_API for recieving data from the scout hardware.
	e.POST("/scout_api/calibrated", func(c echo.Context) error {
		return controllers.ScoutCalibrated(db, c)
	})
	e.POST("/scout_api/interaction", func(c echo.Context) error {
		return controllers.ScoutInteraction(db, c)
	})
	e.POST("/scout_api/log", func(c echo.Context) error {
		return controllers.ScoutLog(db, c)
	})
	e.POST("/scout_api/heartbeat", func(c echo.Context) error {
		return controllers.ScoutHeartbeat(db, c)
	})

	if err := e.Start(config.Address); err != nil {
		e.Logger.Fatal(err)
	}
}
