/*
  The purpose of this code is to scan wifi based on OS (mac, linux and windows).
  First trying in mac first and with airport mac command.
  Next step :
  1. Covering for linux
  2. Covering for windows
  3. Connecting to SSID
  4. Getting data from another machine through WLan
  5. Sending data to another machine through WLan
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// WlanInfo structure model
type WlanInfo struct {
	SSID  string `json:"SSID"`
	BSSID string `json:"mac"`
	Power int    `json:"powerx,omitempty"`
}

// AccessPoint structure model
type AccessPoint struct {
	Wlan []WlanInfo `json:"wlan"`
}

func main() {
	var accessPoint AccessPoint

	// executing wlan scan using mac airport command.
	airportCmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-s")
	airportCmdOutput, err := airportCmd.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lines := strings.Split(string(airportCmdOutput), "\n")

	var wlanList []WlanInfo
	for _, line := range lines {
		columns := strings.Fields(line)
		if len(columns) > 0 {
			match, _ := regexp.MatchString("^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$", columns[1])
			if match {
				power, _ := strconv.Atoi(columns[2])
				wlan := WlanInfo{
					SSID:  columns[0],
					BSSID: columns[1],
					Power: power,
				}

				wlanList = append(wlanList, wlan)
			}
		}
	}

	accessPoint = AccessPoint{Wlan: wlanList}
	fmt.Println(accessPoint)
}
