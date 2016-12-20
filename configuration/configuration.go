/*
 * Copyright (C) 2016 Clinton Freeman
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

package configuration

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

type Configuration struct {
	DBUserName        string // The name of the user with read/write privileges on DBName
	DBPassword        string // The password of the user with read/write privileges on DBName.
	DBName            string // The name of the database that holds the production data.
	DBTestName        string // The name of the database that holds testing data.
	Address           string // The address and port that the mothership is accessible on.
	StaticAssets      string // The path to the static assets rendered by the mothership.
	SummariseInterval int    // The number of milliseconds to wait between updating the interaction summaries.
}

func GetDataDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]) + "/data")
	if err != nil {
		return ""
	}

	err = os.Mkdir(dir, os.ModeDir)
	if os.IsNotExist(err) {
		return ""
	}

	return dir
}

func SaveAsJSON(v interface{}, fileName string) error {
	f, err := os.Create(fileName)
	if os.IsNotExist(err) {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return w.Flush()
}

func Parse(configFile string) (c Configuration, err error) {
	c = Configuration{"mtf", "mothership", "", "mothership_test", ":80", "public", 1000}

	// Open the configuration file.
	file, err := os.Open(configFile)
	if err != nil {
		return c, err
	}
	defer file.Close()

	// Parse JSON in the configuration file.
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	return c, err
}
