// Copyright 2016 yryz Author. All Rights Reserved.

package ds18b20

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Sensors get all connected sensor IDs as array
func Sensors() ([]string, error) {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/w1_bus_master1/w1_master_slaves")
	if err != nil {
		return nil, fmt.Errorf("error reading sensor list: %w", err)
	}

	sensors := strings.Split(string(data), "\n")
	if len(sensors) > 0 {
		sensors = sensors[:len(sensors)-1]
	}

	return sensors, nil
}

// Temperature get the temperature of a given sensor
func Temperature(sensor string) (float64, error) {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/" + sensor + "/w1_slave")
	if err != nil {
		return 0.0, fmt.Errorf("error reading data from sensor: %w", err)
	}

	raw := string(data)

	if !strings.Contains(raw, " YES") {
		return 0.0, fmt.Errorf("checksum verification failed")
	}

	i := strings.LastIndex(raw, "t=")
	if i == -1 {
		return 0.0, fmt.Errorf("could not find temperature in sensor output")
	}

	c, err := strconv.ParseFloat(raw[i+2:len(raw)-1], 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing float value: %w", err)
	}

	return c / 1000.0, nil
}
