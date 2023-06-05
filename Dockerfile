#----------------------------------------------------------------------------
#  SPDX-License-Identifier: Apache-2.0
#
#  Copyright (c) 2023 Cisco Systems, Inc.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#----------------------------------------------------------------------------


############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/app/appd-swagger/
COPY src/rulexlator/. .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/appd-swagger

############################
# STEP 2 build gui
############################
FROM node:alpine AS gui-builder
WORKDIR /app/swagger-gui
COPY ./src/frontend .
RUN npm install
RUN npm install -g @quasar/cli
RUN quasar build

############################
# STEP 3 build a small image
############################
FROM alpine:latest
# Copy our static executable.
COPY --from=builder /go/bin/appd-swagger /go/bin/appd-swagger
COPY --from=gui-builder /app/swagger-gui/dist/spa /gui

# Expose port 8686
EXPOSE 8686 
# Run the hello binary.
CMD ["/go/bin/appd-swagger"]
