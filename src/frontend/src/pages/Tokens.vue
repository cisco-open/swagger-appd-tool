<!--
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
-->

<template>
  <div>
    <q-input dense v-model="masterToken" label="Master Authentication Token"></q-input>
    <q-btn color="primary" class="q-mr-md col-auto" @click="getTokens">Load Tokens</q-btn>
    <div class="text-h5 q-ma-sm">Token List</div>
    <div v-for="app in applTokens" v-bind:key="app.appId">{{app.appName}} - {{app.token}}</div>
  </div>
</template>

<script>

export default {
  name: 'Tokens',
  data () {
    return {
      masterToken: "",
      applTokens: [],
    }
  },
  methods: {
    async getTokens () {
      try {
        const response = await this.$axios.get('/api/sec/tokens/' + this.masterToken, {
        })
        if (response.status === 200) {
          let tokens = response.data
          this.applTokens = tokens.sort(function(a, b){return a.appName < b.appName ? -1 : a.appName === b.appName ? 0 : 1});
          // this.applTokens = tokens
        }
      } catch (err) {
        this.$q.notify('Cannot load tokens - ' + err)
      }
    }
  }
}
</script>
