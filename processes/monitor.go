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
#cgo darwin CFLAGS: -I/usr/local/opt/opencv3/include -I/usr/local/opt/opencv3/include/opencv
#cgo linux CFLAGS: -I/usr/local/include -I/usr/local/include/opencv
#cgo CFLAGS: -Wno-error
#cgo darwin LDFLAGS: -L/usr/local/opt/opencv3/lib
#cgo linux LDFLAGS: -L/usr/local/lib -L/usr/lib
#cgo darwin LDFLAGS: -lstdc++ -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lopencv_highgui -lopencv_core -lopencv_features2d -lopencv_video -lCVBindings -lopencv_core
#cgo linux LDFLAGS: -lm -lstdc++ -lz -ldl -lpthread -lv4l1 -lv4l2 -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lopencv_highgui -lCVBindings -lopencv_video -lopencv_core
#include "opencv2/videoio/videoio_c.h"
#include "opencv2/imgcodecs/imgcodecs_c.h"
#include "cv.h"
#include "highgui.h"
#include "CVBindings.h"
*/
import "C"

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MeasureTheFuture/scout/models"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"
	"unsafe"
)

func getVideoSource(videoFile string) (camera *C.CvCapture, err error) {
	log.Printf("getting video source %v", videoFile)
	if videoFile != "" {
		log.Printf("INFO: Getting video file.")
		file := C.CString(videoFile)
		camera = C.cvCaptureFromFile(file)
		C.free(unsafe.Pointer(file))
		if camera == nil {
			return camera, errors.New("Unable to open a video file. Shutting down scout.")
		}

		return camera, nil
	} else {
		log.Printf("INFO: Opening webcam.")
		camera = C.cvCreateCameraCapture(0)
		if camera == nil {
			return camera, errors.New("Unable to open webcam. Shutting down scout.")
		}

		// Make sure the webcam is set to 1080p.
		C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_WIDTH, 1920)
		C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_HEIGHT, 1080)
		C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_BUFFERSIZE, 1)

		return camera, nil
	}
}

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
	// It takes a little while for the white balance to stabalise on the logitech.
	// So grab a frame, wait to stabalise for white balance to stabalise, then grab again
	// for and save as the calibration frame.
	camera, err := getVideoSource(videoFile)
	if err != nil {
		// No valid video source, abort.
		log.Printf("ERROR: Unable to get video source")
		log.Print(err)
		return
	}
	calibrationFrame := C.cvQueryFrame(camera)
	time.Sleep(1250 * time.Millisecond)
	C.cvReleaseCapture(&camera)

	camera, err = getVideoSource(videoFile)
	if err != nil {
		// No valid video source, abort
		log.Printf("ERROR: Unable to get video source")
		log.Print(err)
		return
	}
	defer C.cvReleaseCapture(&camera)

	// Build the calibration image from the first frame that comes off the camera.
	calibrationFrame = C.cvQueryFrame(camera)
	fileName := "calibrationFrame.jpg"
	file := C.CString(fileName)
	C.cvSaveImage(file, unsafe.Pointer(calibrationFrame), nil)
	C.free(unsafe.Pointer(file))

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

	camera, err := getVideoSource(videoFile)
	if err != nil {
		// No valid video source. Abort measuring.
		log.Printf("ERROR: Unable to get video source")
		log.Print(err)
		return
	}
	defer C.cvReleaseCapture(&camera)

	// Make sure we release the camera when the operating system crushes us.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Printf("INFO: The OS shut down the scout.")
		C.cvReleaseCapture(&camera)
		return
	}()

	scene := models.InitScene(s)

	// Build the calibration frame from disk.
	var calibrationFrame *C.IplImage
	if _, err := os.Stat("calibrationFrame.jpg"); err == nil {
		file := C.CString("calibrationFrame.jpg")

		calibrationFrame = C.cvLoadImage(file, C.CV_LOAD_IMAGE_COLOR)
		defer C.cvReleaseImage(&calibrationFrame)

		C.free(unsafe.Pointer(file))
	} else {
		log.Printf("ERROR: Unable to measure, missing calibration frame")
		log.Print(err)
		return
	}

	// Create a frame to hold the foreground mask results.
	mask := C.cvCreateImage(C.cvSize(calibrationFrame.width, calibrationFrame.height), C.IPL_DEPTH_8U, 1)
	defer C.cvReleaseImage(&mask)

	// Push the initial calibration frame into the MOG2 image subtractor.
	C.initMOG2(C.int(s.MogHistoryLength), C.double(s.MogThreshold), C.int(s.MogDetectShadows))
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
			case c == models.STOP_MEASURE:
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
		C.cvSmooth(unsafe.Pointer(mask), unsafe.Pointer(mask), C.CV_GAUSSIAN, 3, 0, 0.0, C.double(s.GaussianSmooth))
		C.cvThreshold(unsafe.Pointer(mask), unsafe.Pointer(mask), C.double(s.ForegroundThresh), 255, 0)
		C.cvDilate(unsafe.Pointer(mask), unsafe.Pointer(mask), nil, C.int(s.DilationIterations))

		// Detect contours in filtered foreground mask
		storage := C.cvCreateMemStorage(0)
		contours := C.cvCreateSeq(0, C.size_t(unsafe.Sizeof(C.CvSeq{})), C.size_t(unsafe.Sizeof(C.CvPoint{})), storage)
		offset := C.cvPoint(C.int(0), C.int(0))
		num := int(C.cvFindContours(unsafe.Pointer(mask), storage, &contours, C.int(unsafe.Sizeof(C.CvContour{})),
			C.CV_RETR_LIST, C.CV_CHAIN_APPROX_SIMPLE, offset))

		var detectedObjects []models.Waypoint

		// Track each of the detected contours.
		for contours != nil {
			area := float64(C.cvContourArea(unsafe.Pointer(contours), C.cvSlice(0, 0x3fffffff), 0))

			// Only track large objects.
			if area > s.MinArea {
				boundingBox := C.cvBoundingRect(unsafe.Pointer(contours), 0)
				w := int(boundingBox.width / 2)
				h := int(boundingBox.height / 2)
				x := int(boundingBox.x) + w
				y := int(boundingBox.y) + h

				detectedObjects = append(detectedObjects, models.Waypoint{x, y, w, h, 0.0})

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

		scene.Update(db, detectedObjects)
		C.cvClearMemStorage(storage)
		C.cvReleaseMemStorage(&storage)

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
	}

	log.Printf("INFO: Finished measure")
	scene.Close(db)
}
