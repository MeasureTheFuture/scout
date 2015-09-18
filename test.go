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

// #cgo darwin CFLAGS: -I/opt/local/include -I/opt/local/include/opencv
// #cgo linux CFLAGS: -I/usr/include -I/usr/include/opencv
// #cgo CFLAGS: -Wno-error
// #cgo darwin LDFLAGS: -L/opt/local/lib
// #cgo linux LDFLAGS: -L/usr/lib
// #cgo LDFLAGS: -lopencv_highgui -lopencv_core -lopencv_video
// #include "cv.h"
// #include "highgui.h"
import "C"

import (
	"log"
	"unsafe"
)

func main() {
	log.Printf("INFO: Starting sensor.\n")
	camera := C.cvCaptureFromCAM(-1)

	if camera == nil {
		log.Printf("WARNING: No camera detected. Shutting down sensor.\n")
		return
	}

	C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_WIDTH, 1280)
	C.cvSetCaptureProperty(camera, C.CV_CAP_PROP_FRAME_HEIGHT, 720)

	frame := C.cvQueryFrame(camera)
	file := C.CString("frame.png")
	C.cvSaveImage(file, unsafe.Pointer(frame), nil)
	C.free(unsafe.Pointer(file))
	C.cvReleaseImage(&frame)
}
