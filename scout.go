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

// #cgo darwin CFLAGS: -I/usr/local/opt/opencv3/include -I/usr/local/opt/opencv3/include/opencv
// #cgo linux CFLAGS: -I/usr/local/include -I/usr/local/include/opencv
// #cgo CFLAGS: -Wno-error
// #cgo darwin LDFLAGS: -L/usr/local/opt/opencv3/lib
// #cgo linux LDFLAGS: -L/usr/local/lib
// #cgo darwin LDFLAGS: -lstdc++ -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lopencv_highgui -lopencv_core -lopencv_features2d -lopencv_video -lopencv_hal -lCVBindings
// #cgo linux LDFLAGS: -lm -lstdc++ -lz -ldl -lpthread -lippicv -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lIlmImf -llibpng -llibjasper -llibjpeg -llibwebp -llibtiff -lopencv_highgui -lCVBindings -lopencv_video -lopencv_core -lopencv_hal -ltbb
// #include "cv.h"
// #include "highgui.h"
// #include "CVBindings.h"
import "C"

import (
	"log"
	"strconv"
	"unsafe"
)

func monitor(config Configuration, videoFile string) {
	// Try starting monitor with a video file source first.
	camera := C.cvCaptureFromFile(C.CString(videoFile))
	if camera == nil {
		// No valid video file found, fallback to webcam.
		camera = C.cvCaptureFromCAM(-1)

		if camera == nil {
			// No valid webcam detected either. Shutdown monitoring.
			log.Printf("ERROR: Unable to open a video source. Shutting down scout.\n")
			return
		}

		C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_WIDTH, 1280)
		C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_HEIGHT, 720)
	}

	scene := initScene()

	// Build the calibration frame from the first frame from the camera.
	calibrationFrame := C.cvQueryFrame(camera)
	file := C.CString("calibrationFrame.png")
	C.cvSaveImage(file, unsafe.Pointer(calibrationFrame), nil)
	C.free(unsafe.Pointer(file))

	// Create a frame to hold the foreground mask results.
	mask := C.cvCreateImage(C.cvSize(calibrationFrame.width, calibrationFrame.height), C.IPL_DEPTH_8U, 1)

	// Push the initial calibration frame into the MOG2 image subtractor.
	C.initMOG2(C.int(config.MogHistoryLength), C.double(config.MogThreshold), C.int(config.MogDetectShadows))
	C.applyMOG2(unsafe.Pointer(calibrationFrame), unsafe.Pointer(mask))

	// DEBUG - frame counter for use with filenames.
	i := 0

	// Start monitoring from the camera.
	for C.cvGrabFrame(camera) != 0 {
		// Subtract the calibration frame from the current frame.
		nextFrame := C.cvQueryFrame(camera)
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
			//log.Printf("A: " + strconv.FormatFloat(float64(area), 'E', -1, 32))

			// Only track large objects.
			if area > config.MinArea {
				boundingBox := C.cvBoundingRect(unsafe.Pointer(contours), 0)
				w := int(boundingBox.width / 2)
				h := int(boundingBox.height / 2)
				x := int(boundingBox.x) + w
				y := int(boundingBox.y) + h

				detectedObjects = append(detectedObjects, Waypoint{x, y, w, h, 0.0})

				// Debug -- Render contours and bounding boxes.
				pt1 := C.cvPoint(boundingBox.x, boundingBox.y)
				pt2 := C.cvPoint(boundingBox.x+boundingBox.width, boundingBox.y+boundingBox.height)
				C.cvDrawContours(unsafe.Pointer(nextFrame), contours, C.cvScalar(12.0, 212.0, 250.0, 255), C.cvScalar(0, 0, 0, 0), 2, 1, 8, offset)
				C.cvRectangle(unsafe.Pointer(nextFrame), pt1, pt2, C.cvScalar(16.0, 8.0, 186.0, 255), C.int(5), C.int(8), C.int(0))
			} else {
				num--
			}

			contours = contours.h_next
		}

		monitorScene(&scene, detectedObjects)

		// DEBUG - render the interaction path we have detected so far.
		for _, i := range scene.Interactions {
			for _, w := range i.Path {
				pt1 := C.cvPoint(C.int(w.XPixels), C.int(w.YPixels))
				C.cvCircle(unsafe.Pointer(nextFrame), pt1, C.int(10), C.cvScalar(0.0, 46.0, 109.0, 255), C.int(2), C.int(8), C.int(0))
			}
		}

		file = C.CString("f" + strconv.Itoa(i) + "-detected.png")
		C.cvSaveImage(file, unsafe.Pointer(nextFrame), nil)
		saveScene(string("f"+strconv.Itoa(i)+"-metadata.json"), &scene)
		i++
	}

	C.cvReleaseImage(&mask)
	C.cvReleaseImage(&calibrationFrame)
}
