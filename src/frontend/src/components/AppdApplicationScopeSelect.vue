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
  <q-select
    label="AppDynamics Application Scope"
    outlined
    dense
    :model="model"
    :options="options"
    option-value="id"
    option-label="name"
    use-input
    input-debounce="0"
    @filter="filterFn"
    @filter-abort="abortFilterFn"
    @input="(val) => $emit('input', val)"
  >
    <template v-slot:no-option>
      <q-item>
        <q-item-section class="text-grey">
          No results
        </q-item-section>
      </q-item>
    </template>
  </q-select>
</template>

<script>
import { appdMixins } from '../appdMixins'

export default {
  name: 'AppdApplicationScopeSelect',
  mixins: [appdMixins],
  props: [
    'value',
    'appId'
  ],
  data () {
    return {
      options: null,
      allOptions: null,
      appLoaded: ""
    }
  },
  computed: {
    model: function () {
      return this.value
    }
  },
  methods: {
    async filterFn (val, update, abort) {
      console.log(this.options, this.model)
      if (this.options !== null && this.appId === this.appLoaded) {
        console.log('Val:', val)
        if (val !== '') {
          update(() => {
            const needle = val.toLowerCase()
            this.options = this.allOptions.filter(v => v.summary.name.toLowerCase().indexOf(needle) > -1)
          })
        } else {
          this.options = this.allOptions.map(v => v)
          update()
        }
        return
      }
      var options
      options = await this.appdGetAppScopes(this.appId, (e) => this.$q.notify('Cannot load AppDynamics applications: ' + e))
      this.allOptions = options.map(v => ({id: v.summary.id, "name": v.summary.name, ...v} ))
      this.options = this.allOptions.map(v => v)
      this.appLoaded = this.appId
      update()
    },

    abortFilterFn () {
      // console.log('delayed filter aborted')
    }
  }
}
</script>
