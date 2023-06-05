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
    <div class="text-h5 q-ma-sm">Upload Swagger File</div>
    <q-input class="q-ma-sm" style="max-width: 300px" v-model="appdAppToken" label="Application Access Token"/>
    <appd-application-select class="q-ma-sm" style="max-width: 300px" v-model="appdApplication"/>
    <appd-application-scope-select class="q-ma-sm" style="max-width: 300px" v-model="appdApplicationScope" :appId="appdApplication?.id"/>
    <q-select class="q-ma-sm" style="max-width: 300px" outlined dense label="Select language" v-model="langId" :options="languages" emitValue mapOptions />
    <q-uploader
      :disable="appdApplication === null || appdApplicationScope === null || langId === null || appdAppToken === ''"
      class="q-ma-sm" 
      v-bind:url="'/api/swagger/upload/' + appdApplication?.id + '/' + appdApplicationScope?.name + '/' + langId + '/' + appdAppToken"
      label="Swagger File"
      style="max-width: 300px"
    />
  </div>
</template>

<script>
import AppdApplicationSelect from '../components/AppdApplicationSelect.vue'
import AppdApplicationScopeSelect from '../components/AppdApplicationScopeSelect.vue'

export default {
  name: 'Swagger',
  components: {AppdApplicationSelect, AppdApplicationScopeSelect},
  data () {
    return {
      appdApplication: null,
      appdApplicationScope: null,
      appdAppToken: "",
      langId: null,
      languages: [
        {label: "Java", value: 0},
        {label: ".NET", value: 1},
        {label: "NodeJS", value: 2},
        {label: "Python", value: 3},
        {label: "PHP", value: 4},
        {label: "Apache HTTPD", value: 5},
      ]
    }
  },
  watch: {
    appdApplication: function (val) {
      this.appdApplicationScope = null
    }
  }
}
</script>
