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

/*
#cgo darwin CFLAGS: -I/usr/local/opt/opencv@3/include -I/usr/local/opt/opencv@3/include/opencv
#cgo linux CFLAGS: -I/usr/local/include -I/usr/local/include/opencv
#cgo CFLAGS: -Wno-error
#cgo darwin LDFLAGS: -L/usr/local/opt/opencv@3/
#cgo linux LDFLAGS: -L/usr/local/lib -L/usr/lib
#cgo darwin LDFLAGS: -lstdc++ -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lopencv_highgui -lopencv_core -lopencv_features2d -lopencv_video -lopencv_core -lCVBindings
#cgo linux LDFLAGS: -lm -lstdc++ -lz -ldl -lpthread -lv4l1 -lv4l2 -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lopencv_highgui -lopencv_video -lopencv_core -lCVBindings
#include "stdlib.h"
#include "CVBindings.h"
*/
import "C"

import (
	"database/sql"
	"github.com/MeasureTheFuture/scout/configuration"
	"github.com/MeasureTheFuture/scout/models"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"unsafe"
)

func Monitor(db *sql.DB, deltaC chan models.Command, videoFile string, debug bool) {

	// All OpenCV operations must run on the OS thread to access the webcam.
	runtime.LockOSThread()

	for {
		c := <-deltaC

		switch {
		case c == models.CALIBRATE:
			log.Printf("INFO: Calibrating scout.")
			calibrate(db, videoFile)

		case c == models.START_MEASURE:
			log.Printf("INFO: Starting measure")

			// Create a hidden file to store the measuring state across reboots.
			f, err := os.Create(".mtf-measure")
			f.Close()
			if err != nil && os.IsNotExist(err) {
				log.Printf("ERROR: Unable to create .mtf-measure file")
				log.Print(err)
			}

			measure(db, deltaC, videoFile, debug)

		case c == models.STOP_MEASURE:
			log.Printf("INFO: Stopping measure")

			// Delete the hidden file to indicate that measuring has stopped across reboots.
			err := os.Remove(".mtf-measure")
			if err != nil && os.IsNotExist(err) {
				log.Printf("ERROR: Unable to remove .mtf-measure file")
				log.Print(err)
			}
		}
	}

	runtime.UnlockOSThread()
}

func calibrate(db *sql.DB, videoFile string) {
	srcFile := C.CString(videoFile)
	dstFile := C.CString("calibrationFrame.jpg")

	success := C.calibrate(srcFile, dstFile, C.int(configuration.FrameW), C.int(configuration.FrameH))

	C.free(unsafe.Pointer(srcFile))
	C.free(unsafe.Pointer(dstFile))

	if success != true {
		log.Printf("ERROR: Unable to Calibrate")
		return
	}

	// Update the DB with the latest calibration details.
	s, err := models.GetScoutByUUID(db, models.GetScoutUUID(db))
	if err != nil {
		log.Printf("ERROR: Unable to calibrate, can't fetch scout from DB")
		log.Print(err)
		return
	}

	frame, err := ioutil.ReadFile("calibrationFrame.jpg")
	if err != nil {
		log.Printf("ERROR: Unable to calibrate, can't fetch frame from disk.")
		log.Print(err)
		return
	}
	s.UpdateCalibrationFrame(db, frame)

	s.State = models.CALIBRATED
	err = s.Update(db)
	if err != nil {
		log.Printf("ERROR: Unable to calibrate. can't update scout DB")
		log.Print(err)
		return
	}
}

func measure(db *sql.DB, deltaC chan models.Command, videoFile string, debug bool) {
	s := models.GetScout(db)

	if _, err := os.Stat("calibrationFrame.jpg"); err != nil {
		log.Printf("ERROR: Unable to measure, missing calibration frame")
		log.Print(err)
		return
	}

	srcFile := C.CString(videoFile)
	calFile := C.CString("calibrationFrame.jpg")

	success := C.startMeasure(srcFile, calFile,
							  C.int(configuration.FrameW), C.int(configuration.FrameH),
							  C.int(s.MogHistoryLength), C.double(s.MogThreshold), C.int(s.MogDetectShadows))

	C.free(unsafe.Pointer(srcFile))
	C.free(unsafe.Pointer(calFile))

	if success != true {
		log.Printf("ERROR: Unable to get video source")
		return
	}
	defer C.stopMeasure()

	// Make sure we release the camera when the operating system crushes us.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Printf("INFO: The OS shut down the scout.")
		C.stopMeasure()
		return
	}()

	scene := models.InitScene(s)

	// Current frame counter.
	//frame := int64(0)
	measuring := true

	// Start monitoring from the camera.
	for measuring {
		// See if there are any new commands on the deltaC channel.
		select {
		case c := <-deltaC:
			switch {
			case c == models.STOP_MEASURE:
				log.Printf("INFO: Stopping measure")
				measuring = false
			}

		default:
			// Procceed with measuring.
		}

		numObjects := C.int(0)
		objects := C.grabFrame(&numObjects,
							   C._Bool(debug),
					           C.double(s.GaussianSmooth),
					           C.double(s.ForegroundThresh),
					           C.int(s.DilationIterations),
					           C.double(s.MinArea),
					           C.double(s.MaxArea))
		o := (*[1<<30]C.int)(unsafe.Pointer(objects))

		var detectedObjects []models.Waypoint
		for i := C.int(0); i < numObjects; i = i + 4 {
			detectedObjects = append(detectedObjects,
									 models.Waypoint{int(o[i]),
									 				 int(o[i + 1]),
									 				 int(o[i + 2]),
									 				 int(o[i + 3]), 0.0})
		}

		C.free(unsafe.Pointer(objects))

		scene.Update(db, detectedObjects)

		/**
		TODO: Need a new method call for debug printing the interaction path.
		if debug {
			var font C.CvFont
			C.cvInitFont(&font, C.CV_FONT_HERSHEY_SIMPLEX, C.double(0.5), C.double(0.5), C.double(1.0), C.int(2), C.CV_AA)
			txt := C.CString("Hello friend.")
			C.cvPutText(unsafe.Pointer(nextFrame), txt, C.cvPoint(2, 2), &font, C.cvScalar(255.0, 255.0, 255.0, 255))
			C.free(unsafe.Pointer(txt))

			// DEBUG -- render current interaction path for detected objects.
			for _, i := range scene.Interactions {
				for _, w := range i.Path {
					pt1 := C.cvPoint(C.int(w.XPixels), C.int(w.YPixels))
					C.cvCircle(unsafe.Pointer(nextFrame), pt1, C.int(10), C.cvScalar(109.0, 46.0, 0.0, 255), C.int(2), C.int(8), C.int(0))
				}

				w := i.LastWaypoint()
				txt := C.CString(fmt.Sprintf("%01d", i.SceneID))
				C.cvPutText(unsafe.Pointer(nextFrame), txt, C.cvPoint(C.int(w.XPixels+10), C.int(w.YPixels+10)), &font, C.cvScalar(255.0, 255.0, 255.0, 255))
				C.free(unsafe.Pointer(txt))
			}

			for _, i := range scene.IdleInteractions {
				w := i.LastWaypoint()
				pt1 := C.cvPoint(C.int(w.XPixels-w.HalfWidthPixels+5), C.int(w.YPixels-w.HalfHeightPixels+5))
				pt2 := C.cvPoint(C.int(w.XPixels+w.HalfWidthPixels-5), C.int(w.YPixels+w.HalfHeightPixels-5))
				C.cvRectangle(unsafe.Pointer(nextFrame), pt1, pt2, C.cvScalar(16.0, 186.0, 8.0, 255), C.int(5), C.int(8), C.int(0))

				txt := C.CString("i:" + fmt.Sprintf("%01d", i.SceneID))
				C.cvPutText(unsafe.Pointer(nextFrame), txt, C.cvPoint(C.int(w.XPixels+10), C.int(w.YPixels+10)), &font, C.cvScalar(255.0, 255.0, 255.0, 255))
				C.free(unsafe.Pointer(txt))
			}

			file := C.CString("f" + fmt.Sprintf("%03d", frame) + "-detected.jpg")
			C.cvSaveImage(file, unsafe.Pointer(nextFrame), nil)
			C.free(unsafe.Pointer(file))
			frame++

		}
		**/
	}

	log.Printf("INFO: Finished measure")
	scene.Close(db)
}
