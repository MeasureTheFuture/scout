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
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

type Command int

const (
	CALIBRATE Command = iota
	START_MEASURE
	STOP_MEASURE
)

func calibrateHandler(deltaC chan Command, deltaCFG chan Configuration, configFile string,
	config Configuration, w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	newConfig := config

	f, err := strconv.ParseFloat(q.Get("MinArea"), 64)
	if err == nil {
		newConfig.MinArea = f
	}

	i, err := strconv.ParseInt(q.Get("DilationIterations"), 10, 64)
	if err == nil {
		newConfig.DilationIterations = int(i)
	}

	i, err = strconv.ParseInt(q.Get("ForegroundThresh"), 10, 64)
	if err == nil {
		newConfig.ForegroundThresh = int(i)
	}

	i, err = strconv.ParseInt(q.Get("GaussianSmooth"), 10, 64)
	if err == nil {
		newConfig.GaussianSmooth = int(i)
	}

	i, err = strconv.ParseInt(q.Get("MogHistoryLength"), 10, 64)
	if err == nil {
		newConfig.MogHistoryLength = int(i)
	}

	f, err = strconv.ParseFloat(q.Get("MogThreshold"), 64)
	if err == nil {
		newConfig.MogThreshold = f
	}

	i, err = strconv.ParseInt(q.Get("MogDetectShadows"), 10, 64)
	if err == nil {
		newConfig.MogDetectShadows = int(i)
	}

	saveConfiguration(configFile, newConfig)

	deltaCFG <- newConfig
	deltaC <- CALIBRATE
}

func measureStartHandler(deltaC chan Command, w http.ResponseWriter, r *http.Request) {
	deltaC <- START_MEASURE
}

func measureStopHandler(deltaC chan Command, w http.ResponseWriter, r *http.Request) {
	deltaC <- STOP_MEASURE
}

func bindHandlers(deltaC chan Command, deltaCFG chan Configuration, configFile string, config Configuration) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/calibrate", func(w http.ResponseWriter, r *http.Request) {
		calibrateHandler(deltaC, deltaCFG, configFile, config, w, r)
	})

	mux.HandleFunc("/measure/start", func(w http.ResponseWriter, r *http.Request) {
		measureStartHandler(deltaC, w, r)
	})

	mux.HandleFunc("/measure/stop", func(w http.ResponseWriter, r *http.Request) {
		measureStopHandler(deltaC, w, r)
	})

	return mux
}

func controller(deltaC chan Command, deltaCFG chan Configuration, configFile string, config Configuration) {
	mux := bindHandlers(deltaC, deltaCFG, configFile, config)
	log.Fatal(http.ListenAndServe(config.ScoutAddress, mux))
}

func post(fileName string, url string, uuid string, src io.Reader) {
	body := bytes.Buffer{}
	w := multipart.NewWriter(&body)

	part, err := w.CreateFormFile("file", fileName)
	if err != nil {
		log.Printf("ERROR: Unable to create form element for broadcast")
		w.Close()
	}

	_, err = io.Copy(part, src)
	if err != nil {
		log.Printf("ERROR: unable to copy frame into multipart message")
		w.Close()
	}

	contentType := w.FormDataContentType()
	w.Close()

	req, err := http.NewRequest("POST", url, &body)
	req.Header.Add("Mothership-Authorization", uuid)
	req.Header.Set("Content-Type", contentType)
	if err != nil {
		log.Printf("ERROR: Unable to create multipart message. %+v\n", err)
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Printf("ERROR: Unable to send multipart message. %+v\n", err)
	}
}
