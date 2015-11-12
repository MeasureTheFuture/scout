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

//TODO: Shift the MOG2Bindings linkage to someplace else, after I have finished deborking them.

// #cgo darwin CFLAGS: -I/usr/local/opt/opencv3/include -I/usr/local/opt/opencv3/include/opencv -I/Users/cfreeman/Projects/measure-the-future/code/MOG2Bindings
// #cgo linux CFLAGS: -I/usr/local/include -I/usr/local/include/opencv
// #cgo CFLAGS: -Wno-error
// #cgo darwin LDFLAGS: -L/usr/local/opt/opencv3/lib -L/Users/cfreeman/Projects/measure-the-future/code/MOG2Bindings
// #cgo linux LDFLAGS: -L/usr/local/lib
// #cgo darwin LDFLAGS: -lstdc++ -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lopencv_highgui -lopencv_core -lopencv_video -lopencv_hal -lMOG2Bindings
// #cgo linux LDFLAGS: -lm -lstdc++ -lz -ldl -lpthread -lippicv -lopencv_imgcodecs -lopencv_imgproc -lopencv_videoio -lIlmImf -llibpng -llibjasper -llibjpeg -llibwebp -llibtiff -lopencv_highgui -lopencv_core -lopencv_video -lopencv_hal -ltbb
// #include "cv.h"
// #include "highgui.h"
// #include "BackgroundSubtractorMOG2.h"
import "C"

import (
	"log"
	"unsafe"
)

func main() {
	log.Printf("INFO: Starting sensor.\n")

	// Webcam source.
	//camera := C.cvCaptureFromCAM(-1)

	videoFile := C.CString("sample.mp4")
	camera := C.cvCaptureFromFile(videoFile)

	if camera == nil {
		log.Printf("WARNING: No camera detected. Shutting down sensor.\n")
		return
	}

	// Webcam source.
	//C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_WIDTH, 1280)
	//C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_HEIGHT, 720)

	refFrame := C.cvQueryFrame(camera)
	file := C.CString("frame.png")
	C.cvSaveImage(file, unsafe.Pointer(refFrame), nil)
	C.free(unsafe.Pointer(file))

	mask := C.cvCloneImage(refFrame)
	//nexFrame := C.cvQueryFrame(camera)

	mog2 := C.createMOG2(30, 0.5, 1)
	C.applyMOG2(mog2, unsafe.Pointer(C.cvQueryFrame(camera)), unsafe.Pointer(mask), 0.1)

	file = C.CString("mask.png")
	C.cvSaveImage(file, unsafe.Pointer(mask), nil)

	C.cvReleaseImage(&refFrame)
}
