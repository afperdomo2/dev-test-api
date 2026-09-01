import { mount, type VueWrapper } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory, type RouteRecordRaw } from 'vue-router'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import type { Component } from 'vue'

export interface MountWithProvidersOptions {
  props?: Record<string, unknown>
  slots?: Record<string, string>
  routerRoutes?: Array<RouteRecordRaw>
  initialRoute?: string
  pinia?: ReturnType<typeof createPinia>
  queryClient?: QueryClient
  stubs?: Record<string, boolean>
  global?: Record<string, unknown>
}

export async function mountWithProviders(
  component: Component,
  options: MountWithProvidersOptions = {},
): Promise<VueWrapper<unknown>> {
  const pinia = options.pinia ?? createPinia()
  setActivePinia(pinia)

  const queryClient = options.queryClient ?? new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  })

  const routes: Array<RouteRecordRaw> = options.routerRoutes ?? [{ path: '/', component: { template: '<div />' } }]
  // ensure we have a fallback route for pushes
  if (!routes.some((r) => r.path === '/login')) {
    routes.push({ path: '/login', name: 'Login', component: { template: '<div>login</div>' } })
  }
  if (!routes.some((r) => r.path === '/')) {
    routes.push({ path: '/', name: 'Home', component: { template: '<div>home</div>' } })
  }

  const router = createRouter({
    history: createMemoryHistory(),
    routes,
  })

  if (options.initialRoute) {
    await router.push(options.initialRoute)
    await router.isReady()
  }

  const wrapper = mount(component, {
    props: options.props as never,
    slots: options.slots as never,
    global: {
      plugins: [pinia, [VueQueryPlugin, { queryClient }], router],
      stubs: {
        RouterLink: {
          template: '<a><slot /></a>',
        },
        'router-link': {
          template: '<a><slot /></a>',
        },
        teleport: true,
        ...options.stubs,
      },
      ...(options.global as object),
    },
  })

  // Wait for router ready
  await router.isReady()

  return wrapper as unknown as VueWrapper<unknown>
}

export function createTestQueryClient(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  })
}
