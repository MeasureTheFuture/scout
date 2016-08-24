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

/*
#cgo darwin CFLAGS: -I/usr/local/opt/opencv3/include -I/usr/local/opt/opencv3/include/opencv
#cgo linux CFLAGS: -I/usr/local/include -I/usr/local/include/opencv
#cgo CFLAGS: -Wno-error
#cgo darwin LDFLAGS: -L/usr/local/opt/opencv3/lib
#cgo linux LDFLAGS: -L/usr/local/lib -L/usr/lib
#cgo darwin LDFLAGS: -lstdc++ -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lopencv_highgui -lopencv_core -lopencv_features2d -lopencv_video -lopencv_hal -lCVBindings
#cgo linux LDFLAGS: -lm -lstdc++ -lz -ldl -lpthread -lv4l1 -lv4l2 -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lopencv_highgui -lCVBindings -lopencv_video -lopencv_core
#include "cv.h"
#include "highgui.h"
#include "CVBindings.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"unsafe"
)

func getVideoSource(videoFile string) (camera *C.CvCapture, err error) {
	if videoFile != "" {
		file := C.CString(videoFile)
		camera = C.cvCaptureFromFile(file)
		C.free(unsafe.Pointer(file))
		if camera == nil {
			return camera, errors.New("Unable to open a video file. Shutting down scout.")
		}

		return camera, nil
	} else {
		camera = C.cvCreateCameraCapture(0)
		if camera == nil {
			return camera, errors.New("Unable to open webcam. Shutting down scout.")
		}

		// Make sure the webcam is set to 720p.
		C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_WIDTH, 1280)
		C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_HEIGHT, 720)
		C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_BUFFERSIZE, 1)

		return camera, nil
	}
}

func monitor(deltaC chan Command, deltaCFG chan Configuration,
	videoFile string, debug bool, config Configuration) {

	runtime.LockOSThread() // All OpenCV operations must run on the OS thread to access the webcam.

	for {
		c := <-deltaC

		// See if the configuration has been changed by calibration
		select {
		case config = <-deltaCFG:
		default:
		}

		switch {
		case c == CALIBRATE:
			log.Printf("INFO: Calibrating scout.")
			calibrate(videoFile, config)

		case c == START_MEASURE:
			log.Printf("INFO: Starting measure")
			measure(deltaC, videoFile, debug, config)

		case c == STOP_MEASURE:
			log.Printf("INFO: Stopping measure")
			// Nothing to do at the moment.
		}
	}

	runtime.UnlockOSThread()
}

func calibrate(videoFile string, config Configuration) {
	camera, err := getVideoSource(videoFile)
	if err != nil {
		// No valid webcam detected either. Shutdown the scout.
		log.Printf("ERROR: %s\n", err)
		return
	}
	defer C.cvReleaseCapture(&camera)

	// Build the calibration image from the first frame that comes off the camera.
	calibrationFrame := C.cvQueryFrame(camera)
	fileName := "calibrationFrame.jpg"
	file := C.CString(fileName)
	C.cvSaveImage(file, unsafe.Pointer(calibrationFrame), nil)
	C.free(unsafe.Pointer(file))

	// Broadcast calibration results to the mothership.
	f, err := os.Open(fileName)
	if err != nil {
		log.Printf("ERROR: Unable to open calibration frame to broadcast")
		return
	}
	defer f.Close()
	post(fileName, config.MothershipAddress+"/scout_api/calibrated", config.UUID, f)
}

func measure(deltaC chan Command, videoFile string, debug bool, config Configuration) {
	camera, err := getVideoSource(videoFile)
	if err != nil {
		// No valid video source. Abort measuring.
		log.Printf("ERROR: %s\n", err)
		return
	}
	defer C.cvReleaseCapture(&camera)

	scene := initScene()

	// Build the calibration frame from disk.
	var calibrationFrame *C.IplImage
	if _, err := os.Stat("calibrationFrame.jpg"); err == nil {
		file := C.CString("calibrationFrame.jpg")

		calibrationFrame = C.cvLoadImage(file, C.CV_LOAD_IMAGE_COLOR)
		defer C.cvReleaseImage(&calibrationFrame)

		C.free(unsafe.Pointer(file))
	} else {
		log.Printf("ERROR: Unable to measure, missing calibration frame")
		return
	}

	// Create a frame to hold the foreground mask results.
	mask := C.cvCreateImage(C.cvSize(calibrationFrame.width, calibrationFrame.height), C.IPL_DEPTH_8U, 1)
	defer C.cvReleaseImage(&mask)

	// Push the initial calibration frame into the MOG2 image subtractor.
	C.initMOG2(C.int(config.MogHistoryLength), C.double(config.MogThreshold), C.int(config.MogDetectShadows))
	C.applyMOG2(unsafe.Pointer(calibrationFrame), unsafe.Pointer(mask))

	// Current frame counter.
	frame := int64(0)
	measuring := true

	// Start monitoring from the camera.
	for measuring && C.cvGrabFrame(camera) != 0 {
		// See if there are any new commands on the deltaC channel.
		select {
		case c := <-deltaC:
			switch {
			case c == STOP_MEASURE:
				log.Printf("INFO: Stopping measure")
				measuring = false
			}

		default:
			// Procceed with measuring.
		}

		// Subtract the calibration frame from the current frame.
		nextFrame := C.cvRetrieveFrame(camera, 0)
		C.applyMOG2(unsafe.Pointer(nextFrame), unsafe.Pointer(mask))

		// Filter the foreground mask to clean up any noise or holes (morphological-closing).
		C.cvSmooth(unsafe.Pointer(mask), unsafe.Pointer(mask), C.CV_GAUSSIAN, C.int(config.GaussianSmooth), 0, 0.0, 0.0)
		C.cvThreshold(unsafe.Pointer(mask), unsafe.Pointer(mask), C.double(config.ForegroundThresh), 255, 0)
		C.cvDilate(unsafe.Pointer(mask), unsafe.Pointer(mask), nil, C.int(config.DilationIterations))

		// Detect contours in filtered foreground mask
		storage := C.cvCreateMemStorage(0)
		contours := C.cvCreateSeq(0, C.size_t(unsafe.Sizeof(C.CvSeq{})), C.size_t(unsafe.Sizeof(C.CvPoint{})), storage)
		offset := C.cvPoint(C.int(0), C.int(0))
		num := int(C.cvFindContours(unsafe.Pointer(mask), storage, &contours, C.int(unsafe.Sizeof(C.CvContour{})),
			C.CV_RETR_LIST, C.CV_CHAIN_APPROX_SIMPLE, offset))

		var detectedObjects []Waypoint

		// Track each of the detected contours.
		for contours != nil {
			area := float64(C.cvContourArea(unsafe.Pointer(contours), C.cvSlice(0, 0x3fffffff), 0))

			// Only track large objects.
			if area > config.MinArea {
				boundingBox := C.cvBoundingRect(unsafe.Pointer(contours), 0)
				w := int(boundingBox.width / 2)
				h := int(boundingBox.height / 2)
				x := int(boundingBox.x) + w
				y := int(boundingBox.y) + h

				detectedObjects = append(detectedObjects, Waypoint{x, y, w, h, 0.0})

				if debug {
					// DEBUG -- Render contours and bounding boxes around detected objects.
					pt1 := C.cvPoint(boundingBox.x, boundingBox.y)
					pt2 := C.cvPoint(boundingBox.x+boundingBox.width, boundingBox.y+boundingBox.height)
					C.cvDrawContours(unsafe.Pointer(nextFrame), contours, C.cvScalar(12.0, 212.0, 250.0, 255), C.cvScalar(0, 0, 0, 0), 2, 1, 8, offset)
					C.cvRectangle(unsafe.Pointer(nextFrame), pt1, pt2, C.cvScalar(16.0, 8.0, 186.0, 255), C.int(5), C.int(8), C.int(0))
				}
			} else {
				num--
			}

			contours = contours.h_next
		}

		scene.update(detectedObjects, config)

		if debug {
			// DEBUG -- render current interaction path for detected objects.
			for _, i := range scene.Interactions {
				for _, w := range i.Path {
					pt1 := C.cvPoint(C.int(w.XPixels), C.int(w.YPixels))
					C.cvCircle(unsafe.Pointer(nextFrame), pt1, C.int(10), C.cvScalar(109.0, 46.0, 0.0, 255), C.int(2), C.int(8), C.int(0))
				}
			}

			file := C.CString("f" + fmt.Sprintf("%03d", frame) + "-detected.png")
			C.cvSaveImage(file, unsafe.Pointer(nextFrame), nil)
			C.free(unsafe.Pointer(file))
			//scene.save(string("f" + strconv.FormatInt(frame, 10) + "-metadata.json"))
			frame++
		}
	}
}
