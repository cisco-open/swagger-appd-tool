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
  <q-item v-if="children === null" :to="link">
    <q-item-section
      v-if="icon"
      avatar
    >
      <q-icon :name="icon" />
    </q-item-section>
    <q-item-section>
      <q-item-label>{{ title }}</q-item-label>
      <q-item-label caption>
        {{ caption }}
      </q-item-label>
    </q-item-section>
  </q-item>

  <q-expansion-item
    v-else
    expand-separator
    :icon="icon"
    :label="title"
    :caption="caption"
    group="default"
    :content-inset-level="0.3"
  >
    <q-item
      v-for="sub in children"
      v-bind:key="sub.title"
      :to="sub.link">
      <q-item-section
        v-if="sub.icon"
        avatar
      >
        <q-icon :name="sub.icon" />
      </q-item-section>
      <q-item-section>
        <q-item-label>{{ sub.title }}</q-item-label>
        <q-item-label caption>
          {{ sub.caption }}
        </q-item-label>
      </q-item-section>
    </q-item>
  </q-expansion-item>
</template>

<script>
export default {
  name: 'DrawerLinks',
  props: {
    title: {
      type: String,
      required: true
    },

    caption: {
      type: String,
      default: ''
    },

    link: {
      type: String,
      default: '#'
    },

    icon: {
      type: String,
      default: ''
    },

    children: {
      default: null
    }
  }
}
</script>
