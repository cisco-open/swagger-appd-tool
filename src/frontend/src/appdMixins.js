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

export const appdMixins = {
  methods: {
    appdGetApplications (onError) {
      return this.$axios.get('/api/appd/apps')
        .then((response) => {
          var apps = response.data
          apps.sort((a, b) => a.name < b.name ? -1 : 1)
          return apps
        })
        .catch(function (error) {
          onError(error)
        })
    },
    appdGetAppTiers (appId, onError) {
      return this.$axios.get('/api/appd/app/' + appId + '/tiers')
        .then((response) => {
          var appTiers = response.data
          appTiers.sort((a, b) => a.name < b.name ? -1 : 1)
          return appTiers
        })
        .catch(function (error) {
          onError(error)
        })
    },
    appdGetAppScopes (appId, onError) {
      return this.$axios.get('/api/appd/app/' + appId + '/scopes')
        .then((response) => {
          var appScopes = response.data.scopes
          appScopes.sort((a, b) => a.summary.name < b.summary.name ? -1 : 1)
          return appScopes
        })
        .catch(function (error) {
          onError(error)
        })
    }
  }
}
