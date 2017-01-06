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

import (
	"database/sql"
	"github.com/MeasureTheFuture/scout/configuration"
	"github.com/MeasureTheFuture/scout/models"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net"
	"syscall"
	"time"
)

func HealthHeartbeat(db *sql.DB, config configuration.Configuration) {
	s, err := models.GetScoutByUUID(db, config.UUID)
	if err != nil {
		log.Printf("ERROR: Unable to start health heartbeat. Scout UUID missing")
		log.Print(err)
		return
	}

	// Send initial health heartbeat on startup.
	err = SaveHeartbeat(db, s)
	if err != nil {
		log.Printf("ERROR: Unable to save health heartbeat. ")
		log.Print(err)
		return
	}

	// Send periodic health heartbeats to the mothership.
	poll := time.NewTicker(time.Minute * 15).C
	for {
		select {
		case <-poll:
			err = SaveHeartbeat(db, s)
			if err != nil {
				log.Printf("ERROR: Unable to save health heartbeat. ")
				log.Print(err)
				return
			}
		}
	}
}

func SaveHeartbeat(db *sql.DB, s *models.Scout) error {
	t, u := getMemoryUsage()
	sh := models.ScoutHealth{s.Id, getCPULoad(), u, t, getStorageUsage(), time.Now().UTC()}

	return sh.Insert(db)
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
