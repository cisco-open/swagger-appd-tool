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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func httpServer() {
	fmt.Printf("Hello Http server\n")

	router := httprouter.New()

	router.POST("/api/login", handleLogin)
	router.GET("/api/appd/apps", handleAppDynamicsApplicationList)
	router.GET("/api/appd/app/:appId/scopes", handleAppDynamicsApplicationScopeList)
	router.GET("/api/sec/verify/:appName/:token", handleAppTokenVerification)
	router.GET("/api/sec/tokens/:authToken", handleGetAppTokens)
	router.GET("/api/sec/token/:appId", handleGetAppToken)

	router.POST("/api/swagger/upload", handleSwaggerUpload)
	router.POST("/api/swagger/upload/:appId", handleSwaggerUpload)
	router.POST("/api/swagger/upload/:appId/:scopeId/:langId/:token", handleSwaggerUpload)

	staticPath := "gui"
	router.NotFound = http.FileServer(http.Dir(staticPath))

	if err := http.ListenAndServe(":8686", router); err != nil {
		log.Print(err)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Default Handler: %s\n", r.URL)
}

func handleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Printf("Default Handler: %s\n", r.URL)
}

func handleAppTokenVerification(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == "GET" {
		appName := params.ByName("appName")
		token := params.ByName("token")

		fmt.Fprintf(w, "{\"verified\":\"%t\"}", verifyToken(appName, token))
	} else {
		http.Error(w, "Method not supported", http.StatusBadRequest)
	}
}

type AppToken struct {
	AppId   int    `json:"appId"`
	AppName string `json:"appName"`
	Token   string `json:"token"`
}

type AppTokens []AppToken

func handleGetAppToken(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == "GET" {
		_, err := strconv.Atoi(params.ByName("appId"))
		if err != nil {
			log.Println(fmt.Errorf("error parsing application ID: %s", err))
			http.Error(w, "Error parsing application ID", http.StatusBadRequest)
		}
		ensureConnected()
		appId := params.ByName("appId")
		app, err := appdClient.Application.GetApplication(appId)
		if err != nil {
			log.Print(fmt.Errorf("%s", err))
		}
		appToken := AppToken{}
		if appId == "0" {
			appToken = AppToken{
				AppId:   0,
				AppName: string(masterPassword),
				Token:   getToken(string(masterPassword)),
			}
		} else {
			appToken = AppToken{
				AppId:   app.ID,
				AppName: app.Name,
				Token:   getToken(app.Name),
			}
		}

		buffer, err := json.Marshal(appToken)
		if err != nil {
			log.Printf("error: %v", err)
		}
		fmt.Fprint(w, string(buffer))
	} else {
		http.Error(w, "Method not supported", http.StatusBadRequest)
	}
}

func handleGetAppTokens(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == "GET" {
		authToken := params.ByName("authToken")
		appTokens := AppTokens{}

		if !verifyToken(string(masterPassword), authToken) {
			fmt.Print("Ivalid credentials\n")
			http.Error(w, "Invalid credentials", http.StatusBadRequest)
			return
		}

		ensureConnected()
		apps, err := appdClient.Application.GetApplications()
		if err != nil {
			log.Print(fmt.Errorf("%s", err))
		}
		/*
			appTokens = append(appTokens, AppToken{
				AppId:   0,
				AppName: "MASTER_TOKEN",
				Token:   getToken("MASTER_TOKEN"),
			})
		*/
		for _, app := range apps {
			appTokens = append(appTokens, AppToken{
				AppId:   app.ID,
				AppName: app.Name,
				Token:   getToken(app.Name),
			})
		}
		buffer, err := json.Marshal(appTokens)
		if err != nil {
			log.Printf("error: %v", err)
		}
		fmt.Fprint(w, string(buffer))
	} else {
		http.Error(w, "Method not supported", http.StatusBadRequest)
	}
}

func handleGenKeyForPassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == "GET" {
		password := params.ByName("password")
		fmt.Fprint(w, getToken(password))
	}

}

func handleAppDynamicsApplicationList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "GET" {
		fmt.Fprint(w, GetAppDAppList())
	} else {
		http.Error(w, "Method not supported", http.StatusBadRequest)
	}
}

func handleAppDynamicsApplicationScopeList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == "GET" {
		appID, err := strconv.Atoi(params.ByName("appId"))
		if err != nil {
			log.Println(fmt.Errorf("error parsing application ID: %s", err))
			http.Error(w, "Error parsing application ID", http.StatusBadRequest)
		}
		fmt.Fprint(w, GetAppDAppScopeList(appID))
	} else {
		http.Error(w, "Method not supported", http.StatusBadRequest)
	}
}

func handleSwaggerUpload(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method == "POST" {
		appID := params.ByName("appId")
		scopeId := params.ByName("scopeId")
		token := params.ByName("token")
		langIdStr := params.ByName("langId")
		langId, err := strconv.Atoi(langIdStr)

		if err != nil {
			http.Error(w, "langId must be int", http.StatusBadRequest)
			return
		}

		if appID == "" || scopeId == "" || langIdStr == "" {
			http.Error(w, "Missing one of mandatory parameters in URL - must include appId, scopeId, langId, and token", http.StatusBadRequest)
		}

		ensureConnected()
		app, err := appdClient.Application.GetApplication(appID)
		if err != nil {
			fmt.Printf("Error getting app by id %s - %v\n", appID, err)
			http.Error(w, "Error looking for app", http.StatusBadRequest)
			return
		}

		/*
			envMasterPassword := os.Getenv("MASTER_ENC_PASSWORD")
			if envMasterPassword != "" {
				masterPassword = []byte(envMasterPassword)
			}
		*/

		if !verifyToken(string(masterPassword), token) {
			if !verifyToken(app.Name, token) {
				fmt.Print("Ivalid credentials\n")
				http.Error(w, "Invalid credentials", http.StatusBadRequest)
				return
			}
		}

		// Maximum upload of 10 MB files
		r.ParseMultipartForm(10 << 20)
		// fmt.Printf("Files: %v\n", r.MultipartForm.File)

		for filename := range r.MultipartForm.File {
			// Get handler for filename, size and headers
			file, handler, err := r.FormFile(filename)
			if err != nil {
				fmt.Println("Error Retrieving the File")
				fmt.Println(err)
				return
			}

			defer file.Close()
			fmt.Printf("Uploaded File: %+v\n", handler.Filename)
			fmt.Printf("File Size: %+v\n", handler.Size)
			fmt.Printf("MIME Header: %+v\n", handler.Header)

			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			processSwaggerFile(buf.String(), appID, scopeId, langId)
		}
		/*
			// Create file
			dst, err := os.Create(handler.Filename)
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Copy the uploaded file to the created file on the filesystem
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		*/
		fmt.Fprintf(w, "Successfully Uploaded File\n")
		fmt.Fprint(w, "")
	} else {
		http.Error(w, "Method not supported", http.StatusBadRequest)
	}
}
