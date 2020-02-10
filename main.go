/*
  The purpose of this code is to scan wifi based on OS (mac, linux and windows).
  First trying in mac first and with airport mac command.
  Next step :
  1. Connecting to SSID
  2. Getting data from another machine through WLan
  3. Sending data to another machine through WLan
  4. Covering for linux OS
  5. Covering for windows OS
*/

package main

import (
	"errors"
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

func listAvailableWifi() (AccessPoint, error) {

	var accessPoint AccessPoint

	// executing wlan scan using mac airport command.
	airportCmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-s")
	airportCmdOutput, err := airportCmd.Output()
	if err != nil {
		return accessPoint, err
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
	return accessPoint, nil
}

func connect(netInterface, ssid, password string) error {

	// executing wlan scan using mac airport command.
	networkCmd := exec.Command("networksetup", "-setairportnetwork", netInterface, ssid, password)
	networkCmdOutput, err := networkCmd.Output()
	if err != nil {
		return err
	}

	if strings.ContainsAny(string(networkCmdOutput), "Failed") {
		return errors.New("failed connect to access point")
	}

	return nil
}

func main() {

	// get list of available access point
	accessPointList, err := listAvailableWifi()
	if err != nil {
		fmt.Println("failed get list of available access point, ", err)
		os.Exit(1)
	}

	for _, accessPoint := range accessPointList.Wlan {
		if accessPoint.SSID == "sample" {
			err = connect("en0", accessPoint.SSID, "sample")
			if err != nil {
				fmt.Println(fmt.Sprintf("failed connect to %s using connection interface %s, %v", accessPoint.SSID, "en0", err))
				os.Exit(1)
			}

			fmt.Println(fmt.Sprintf("successfully connected to %s using interface %s", accessPoint.SSID, "en0"))
			break
		}
	}
}
