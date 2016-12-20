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
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net"
	"os"
	"syscall"
)

type Heartbeat struct {
	UUID    string     // The UUID for the scout.
	Version string     // The version of the protocol used used for transmitting data to the mothership.
	Health  HealthData // The current health status of the scout.

}

type HealthData struct {
	IpAddress   string  // The current IP address of the scout.
	CPU         float32 // The amount of CPU load currently being consumed on the scout. 0.0 - no load, 1.0 - full load.
	Memory      float32 // The amount of memory consumed on the scout. 0.0 - no memory used, 1.0 no memory available.
	TotalMemory float32 // The total number of gigabytes of virtual memory currently available.
	Storage     float32 // The amount of storage consumed on the scout. 0.0 - disk unused, 1.0 disk full.
}

func NewHeartbeat(config Configuration) *Heartbeat {
	t, u := getMemoryUsage()
	h := Heartbeat{config.UUID, "0.1", HealthData{getIpAddress(), getCPULoad(), u, t, getStorageUsage()}}

	return &h
}

func getIpAddress() string {
	addys, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("ERROR: Unable to get the IP address for the scout.")
		log.Print(err)
		return ""
	}

	for _, address := range addys {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func getCPULoad() float32 {
	c, err := load.Avg()
	if err != nil {
		log.Printf("ERROR: Unable to get CPU load for the scout.")
		log.Print(err)
	}

	return float32(c.Load5)
}

func getMemoryUsage() (float32, float32) {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("ERROR: Unable to get memory usage for the scout.")
		log.Print(err)
	}

	return float32(v.Total), (float32(v.UsedPercent) / 100.0)
}

func getStorageUsage() float32 {
	var stat syscall.Statfs_t
	syscall.Statfs("/", &stat)

	size := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)

	return float32(size-free) / float32(size)
}

func postLog(config Configuration, tmpLog string) {
	f, err := os.Open(tmpLog)
	if err != nil {
		log.Printf("Unable to open temporary log")
		log.Print(err)
		return
	}

	post("scout.log", config.MothershipAddress+"/scout_api/log", config.UUID, bufio.NewReader(f))
	f.Close()
	os.Remove(tmpLog)
}

func (h *Heartbeat) post(config Configuration) {
	body := bytes.Buffer{}
	encoder := json.NewEncoder(&body)

	err := encoder.Encode(h)
	if err != nil {
		log.Printf("ERROR: Unable to encode configuration for transport to mothership")
		log.Print(err)
	}

	post("heartbeat.json", config.MothershipAddress+"/scout_api/heartbeat", config.UUID, &body)
}
