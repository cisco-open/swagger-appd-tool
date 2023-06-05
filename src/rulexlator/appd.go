/*
	SPDX-License-Identifier: Apache-2.0

	Copyright (c) 2023 Cisco Systems, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	appdrest "gitlab.com/Divis/appdynamics-rest-api"
)

func ensureConnected() {
	if !appdClientInitted {
		var err error

		hostname := os.Getenv("APPDYNAMICS_CONTROLLER_HOST_NAME")
		port := os.Getenv("APPDYNAMICS_CONTROLLER_PORT")
		ssl := os.Getenv("APPDYNAMICS_CONTROLLER_SSL_ENABLED")
		account := os.Getenv("APPDYNAMICS_AGENT_ACCOUNT_NAME")

		username := os.Getenv("APPDYNAMICS_API_KEY_NAME")
		password := os.Getenv("APPDYNAMICS_API_KEY_SECRET")

		appdPort, err := strconv.Atoi(port)
		if err != nil {
			log.Print("Cannot parse AppD controller port " + "443")
		}
		scheme := "http"
		if ssl == "true" {
			scheme = "https"
		}
		appdClient, err = appdrest.NewClient(
			scheme,
			hostname,
			appdPort,
			username,
			password,
			account)

		if err != nil {
			log.Print("Cannot connect to AppDynamics controller")
			appdClientInitted = false
		} else {
			appdClientInitted = true
		}
	}
}

// GetAppDAppList = get list of AppD applications as json string
func GetAppDAppList() string {
	ensureConnected()
	apps, err := appdClient.Application.GetApplications()
	if err != nil {
		log.Print(fmt.Errorf("%s", err))
	}
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(apps)
	// fmt.Println()

	buffer, err := json.Marshal(apps)
	if err != nil {
		log.Printf("error: %v", err)
	}
	return string(buffer)
}

// GetAppDAppScopeList = get list of AppD application Transaction Detection Scopes as []string
func GetAppDAppScopeList(appId int) string {
	ensureConnected()
	scopes, err := appdClient.TxDetectionRule.GetApplicationsScopes(appId)
	if err != nil {
		log.Print(fmt.Errorf("%s", err))
	}
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("Scopes: %v\n", scopes)

	buffer, err := json.Marshal(scopes)
	if err != nil {
		log.Printf("error: %v, %v", err, buffer)
	}
	return string(buffer)
}
