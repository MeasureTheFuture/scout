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
	"database/sql"
	"flag"
	"github.com/MeasureTheFuture/scout/configuration"
	"github.com/MeasureTheFuture/scout/controllers"
	"github.com/MeasureTheFuture/scout/models"
	"github.com/MeasureTheFuture/scout/processes"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

func main() {
	var configFile string
	var videoFile string
	var logFile string
	var debug bool

	flag.StringVar(&configFile, "configFile", "scout.json", "The path to the configuration file")
	flag.StringVar(&videoFile, "videoFile", "", "The path to a video file to detect motion from instead of a webcam")
	flag.StringVar(&logFile, "logFile", "scout.log", "The output path for log files.")
	flag.BoolVar(&debug, "debug", false, "Should we run scout in debug mode, and render frames of detected materials")
	flag.Parse()

	// Copy the old log file to a temporary location for transmission to the mothership
	// and start a new log for this instance of scout.
	tmpLog := "scout_tmp.log"
	os.Link("scout.log", tmpLog)
	os.Remove("scout.log")

	f, err := os.OpenFile("scout.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Unable to open log file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("INFO: Starting scout.\n")

	// Parse the configuration file.
	config, err := configuration.Parse(configFile)
	if err != nil {
		log.Fatalf("ERROR: Can't parse configuration - %s", err)
	}

	// Open a connection to the database.
	connection := "user=" + config.DBUserName + " dbname=" + config.DBName
	if config.DBPassword != "" {
		connection = connection + " password=" + config.DBPassword
	}

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatalf("ERROR: Can't open database - %s", err)
	}
	defer db.Close()

	// If no scout exists in the DB, bootstrap the DB by creating one.
	c, err := models.NumScouts(db)
	if err != nil {
		log.Fatalf("ERROR: Unable to cound scouts in DB - %s", err)
	}
	if c == 0 {
		ns := models.Scout{"", "0.0.0.0", 8080, false, "Location " + strconv.FormatInt(c+1, 10), "idle", &models.ScoutSummary{},
			2.0, 2, 2, 2, 2, 2.0, 0, 2.0, 0.2, 0.3, 1}
		err = ns.Insert(db)
		if err != nil {
			log.Fatalf("ERROR: Unable to add initial scout to DB.")
		}
	}

	// Start the background processes.
	go processes.SaveLogToDB(tmpLog, db)
	go processes.HealthHeartbeat(db)
	go processes.Summarise(db, config)

	deltaC := make(chan models.Command)
	// Test to see if the scout is still in measurement mode on boot and resume if necessary.
	go func() {
		if _, err := os.Stat(".mtf-measure"); err == nil {
			log.Printf("INFO: Resuming.")
			deltaC <- models.START_MEASURE
		}
	}()
	go processes.Monitor(db, deltaC, videoFile, debug)

	// Start the user interface.
	e := echo.New()
	e.Static("/", config.StaticAssets)
	e.Static("/css", config.StaticAssets+"/css")
	e.Static("/fonts", config.StaticAssets+"/fonts")
	e.Static("/img", config.StaticAssets+"/img")

	// Front-end API for displaying results from the scouts.
	e.GET("/scouts", func(c echo.Context) error {
		return controllers.GetScouts(db, c)
	})

	e.GET("/scouts/:uuid/frame.jpg", func(c echo.Context) error {
		return controllers.GetScoutFrame(db, c)
	})

	e.GET("/scouts/:uuid", func(c echo.Context) error {
		return controllers.GetScout(db, c)
	})

	e.PUT("/scouts/:uuid", func(c echo.Context) error {
		return controllers.UpdateScout(db, c, deltaC)
	})

	e.GET("/download.zip", func(c echo.Context) error {
		return controllers.DownloadData(db, c)
	})

	// Start scout user-interface.
	if err := e.Start(config.Address); err != nil {
		e.Logger.Fatal(err)
	}
}
