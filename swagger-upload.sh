#! /bin/bash

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

SWAGGER_URL="http://localhost:8686/"

APP_NAME=$1
SCOPE_NAME=$2
LANG_ID=$3
SWAGGER_FILE=$4
TOKEN=$5 #_AUaA4E78FzSnkxW

APP_ID=$(curl -s ${SWAGGER_URL}api/appd/apps | jq -r '.[] | select(.name == "'"${APP_NAME}"'") | .id')

SCOPE_ID=$(curl -s ${SWAGGER_URL}api/appd/app/${APP_ID}/scopes | jq -r '.scopes[] | select(.summary.name == "'"${SCOPE_NAME}"'") | .summary.id')

echo ${SWAGGER_URL}api/swagger/upload/${APP_ID}/${SCOPE_ID}/${LANG_ID}/${TOKEN}

UPLOAD=$(curl -s -F 'data=@${SWAGGER_FILE}' ${SWAGGER_URL}api/swagger/upload/${APP_ID}/${SCOPE_ID}/${LANG_ID}/${TOKEN})

echo $UPLOAD
