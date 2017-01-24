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
	"github.com/MeasureTheFuture/scout/processes"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var mainfunc = make(chan func())

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

	//TODO: Shift much of the config metadata into the DB.
	//TODO: Remove the scout ID and just have UUID.
	//TODO: Init UUID for scout on first boot - this used to be done in parse above.
	//TODO: Merge interaction and ScoutInteraction.

	l, err := models.CreateLogFromFile(tmpLog, db, config)
	if err != nil {
		log.FatalF("ERROR: Unable to create log from file - %s", err)
	}
	err = l.Insert(db)
	if err != nil {
		lgo.FatalF("ERROR: Unable to save log to database - %s", err)
	}

	go processes.HealthHeartbeat(db, config)
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

	// Start scout user-interface.
	if err := e.Start(config.Address); err != nil {
		e.Logger.Fatal(err)
	}

	// deltaC := make(chan Command)
	// deltaCFG := make(chan Configuration, 1)

	//go controller(deltaC, deltaCFG, configFile, config)
	// Test to see if the scout is still in measurement mode on boot and resume if necssary.
	// go func() {
	// 	if _, err := os.Stat(".mtf-measure"); err == nil {
	// 		log.Printf("INFO: Resuming.")
	// 		deltaC <- START_MEASURE
	// 	}
	// }()
	// monitor(deltaC, deltaCFG, videoFile, debug, config)
}
