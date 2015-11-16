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

// TODO: Shift the MOG2Bindings linkage to someplace else, after I have finished deborking them.

// #cgo darwin CFLAGS: -I/usr/local/opt/opencv3/include -I/usr/local/opt/opencv3/include/opencv -I/Users/cfreeman/Projects/measure-the-future/code/MOG2Bindings
// #cgo linux CFLAGS: -I/usr/local/include -I/usr/local/include/opencv
// #cgo CFLAGS: -Wno-error
// #cgo darwin LDFLAGS: -L/usr/local/opt/opencv3/lib -L/Users/cfreeman/Projects/measure-the-future/code/MOG2Bindings
// #cgo linux LDFLAGS: -L/usr/local/lib
// #cgo darwin LDFLAGS: -lstdc++ -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lopencv_highgui -lopencv_core -lopencv_features2d -lopencv_video -lopencv_hal -lCVBindings
// #cgo linux LDFLAGS: -lm -lstdc++ -lz -ldl -lpthread -lippicv -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lIlmImf -llibpng -llibjasper -llibjpeg -llibwebp -llibtiff -lopencv_highgui -lopencv_core -lopencv_video -lopencv_hal -ltbb
// #include "cv.h"
// #include "highgui.h"
// #include "CVBindings.h"
import "C"

import (
	"log"
	"unsafe"
	"strconv"
)

func monitor() { 	
	// Webcam source.
	//camera := C.cvCaptureFromCAM(-1)

	videoFile := C.CString("sample2.mov")
	camera := C.cvCaptureFromFile(videoFile)

	if camera == nil {
		log.Printf("WARNING: No camera detected. Shutting down sensor.\n")
		return
	}

	// Webcam source.
	//C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_WIDTH, 1280)
	//C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_HEIGHT, 720)

	// Build the calibration frame from the first frame from the camera.
	calibrationFrame := C.cvQueryFrame(camera)
	file := C.CString("calibrationFrame.png")
	C.cvSaveImage(file, unsafe.Pointer(calibrationFrame), nil)
	C.free(unsafe.Pointer(file))

	// Create a frame to hold the foreground mask results.
	mask := C.cvCreateImage(C.cvSize(calibrationFrame.width, calibrationFrame.height), C.IPL_DEPTH_8U, 1)
	maskI := C.cvCreateImage(C.cvSize(calibrationFrame.width, calibrationFrame.height), C.IPL_DEPTH_8U, 1)

	C.initBlob(15, 1000, 0, 255)

	// Push the initial calibration frame into the MOG2 image subtractor.
	C.initMOG2(500, 30, 0)
	C.applyMOG2(unsafe.Pointer(calibrationFrame), unsafe.Pointer(mask))	

	// Start monitoring from the camera.
	for i := 0; i < 100; i++ {
		// Subtract the calibration frame from the current frame.
		C.cvGrabFrame(camera)
		nextFrame := C.cvQueryFrame(camera)
		C.applyMOG2(unsafe.Pointer(nextFrame), unsafe.Pointer(mask))
		
		// Detect blobs in foreground mask
		var blobs [600]int32

		C.cvNot(unsafe.Pointer(mask), unsafe.Pointer(maskI))

		nBlobs := int(C.detectBlobs(unsafe.Pointer(maskI), (*C.int)(unsafe.Pointer(&blobs[0]))))

		log.Printf("Frame" + strconv.Itoa(i) + ": " + strconv.Itoa(nBlobs))
		for j := 0; j < nBlobs; j++ {
		 	log.Printf("=[" + strconv.Itoa(int(blobs[j*3])) + ", " + strconv.Itoa(int(blobs[(j*3) + 1])) + ", " + strconv.Itoa(int(blobs[(j*3) + 2])) + "]")
		}

		// DEBUG - save what we have so far.
		
		// if nBlobs > 0 {
		// 	point := C.cvPoint(C.int(blobs[0]), C.int(blobs[1]))
		// 	C.cvCircle(unsafe.Pointer(maskI), point, C.int(blobs[2]), C.cvScalar(255.0, 0.0, 255.0, 0.0), C.int(5), C.int(8), 0)
		// }

		file = C.CString("mask" + strconv.Itoa(i) + ".png")
		C.cvSaveImage(file, unsafe.Pointer(maskI), nil)
	}
	
	C.cvReleaseImage(&mask)
	C.cvReleaseImage(&calibrationFrame)
}