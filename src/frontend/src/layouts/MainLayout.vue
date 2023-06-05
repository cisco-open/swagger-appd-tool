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
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="leftDrawerOpen = !leftDrawerOpen"
        />

        <q-toolbar-title>
          Swagger to AppDynamics Translator
        </q-toolbar-title>

        <div>v{{ "0.1" }}</div>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
      content-class="bg-grey-1"
    >
      <q-list>
        <q-item-label
          header
          class="text-grey-8"
        >
        Menu Options
        </q-item-label>
        <DrawerLinks
          v-for="link in drawerLinks"
          :key="link.title"
          v-bind="link"
        />
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script>
import DrawerLinks from 'components/DrawerLinks.vue'

const linksList = [
  {
    title: 'Application Tokens',
    caption: 'Get List of Application Tokens',
    icon: 'school',
    link: '/gui/t0k3nz'
  },
  {
    title: 'Swagger Rules',
    caption: 'Load Swagger Rules to AppDynamics',
    icon: 'school',
    link: '/gui/swagger'
  },
];

import { defineComponent, ref } from 'vue'

export default defineComponent({
  name: 'MainLayout',

  components: {
    DrawerLinks
  },

  setup () {
    const leftDrawerOpen = ref(false)

    return {
      drawerLinks: linksList,
      leftDrawerOpen,
      toggleLeftDrawer () {
        leftDrawerOpen.value = !leftDrawerOpen.value
      }
    }
  }
})
</script>
