<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

interface NavItem {
  title: string
  to: string
  icon: string
  adminOnly?: boolean
  userOnly?: boolean
}

interface NavSection {
  title?: string
  items: Array<NavItem>
}

const navSections: Array<NavSection> = [
  { items: [{ title: 'Dashboard', to: '/', icon: 'mdi-view-dashboard' }] },
  {
    title: 'Preguntas',
    items: [
      { title: 'Preguntas', to: '/questions', icon: 'mdi-help-circle', userOnly: true },
      { title: 'Importar', to: '/questions/import', icon: 'mdi-file-upload', userOnly: true },
      { title: 'Temas', to: '/topics', icon: 'mdi-tag' },
      { title: 'Estadísticas', to: '/stats', icon: 'mdi-chart-bar', userOnly: true },
    ],
  },
  {
    title: 'Estudio',
    items: [
      { title: 'Sesiones', to: '/sessions', icon: 'mdi-play-circle', userOnly: true },
      { title: 'Progreso', to: '/progress', icon: 'mdi-chart-line', userOnly: true },
    ],
  },
  {
    title: 'Administración',
    items: [{ title: 'Usuarios', to: '/users', icon: 'mdi-account-group', adminOnly: true }],
  },
]

const filteredSections = computed(() =>
  navSections
    .map((section) => ({
      ...section,
      items: section.items.filter(
        (item) =>
          (item.adminOnly ? authStore.isAdmin : true) &&
          (item.userOnly ? !authStore.isAdmin : true),
      ),
    }))
    .filter((section) => section.items.length > 0),
)

function logout() {
  authStore.clearSession()
  router.push('/login')
}
</script>

<template>
  <v-app :theme="appStore.theme">
    <v-navigation-drawer v-model="appStore.sidebarOpen" :permanent="true" color="surface">
      <v-list-item
        class="pa-4"
        prepend-icon="mdi-code-braces"
        title="DevTest"
        subtitle="Estudio inteligente"
      />

      <v-divider />

      <v-list density="compact" nav>
        <template v-for="section in filteredSections" :key="section.title ?? '__top'">
          <v-list-subheader v-if="section.title">{{ section.title }}</v-list-subheader>
          <v-list-item
            v-for="item in section.items"
            :key="item.to"
            :to="item.to"
            :prepend-icon="item.icon"
            :title="item.title"
            exact
          />
        </template>
      </v-list>

      <template #append>
        <v-divider />
        <v-list density="compact" nav class="pa-2">
          <v-list-item prepend-icon="mdi-logout" title="Cerrar sesión" @click="logout" />
        </v-list>
      </template>
    </v-navigation-drawer>

    <v-app-bar color="primary" density="compact">
      <v-app-bar-nav-icon @click="appStore.toggleSidebar" />
      <v-toolbar-title>DevTest</v-toolbar-title>

      <v-spacer />

      <v-btn
        :icon="appStore.theme === 'light' ? 'mdi-weather-night' : 'mdi-weather-sunny'"
        variant="text"
        @click="appStore.toggleTheme"
      />

      <v-menu>
        <template #activator="{ props }">
          <v-btn v-bind="props" variant="text" icon="mdi-account-circle" />
        </template>
        <v-list density="compact">
          <v-list-item>
            <v-list-item-title>{{ authStore.user?.email }}</v-list-item-title>
            <v-list-item-subtitle>
              {{ authStore.isAdmin ? 'Administrador' : 'Usuario' }}
            </v-list-item-subtitle>
          </v-list-item>
          <v-divider />
          <v-list-item prepend-icon="mdi-logout" title="Cerrar sesión" @click="logout" />
        </v-list>
      </v-menu>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>

    <v-snackbar
      v-model="appStore.snackbar.show"
      :color="appStore.snackbar.color"
      :timeout="4000"
      location="bottom right"
    >
      {{ appStore.snackbar.message }}
      <template #actions>
        <v-btn variant="text" icon="mdi-close" @click="appStore.hideSnackbar" />
      </template>
    </v-snackbar>
  </v-app>
</template>
