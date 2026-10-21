# SPEC — mobile-push

| Campo | Valor |
|---|---|
| Module id | `mobile-push` |
| Origen | Alcance nuevo mobile (D-26, D-27, D-28) |
| Estado | `ready` |
| Depende de | `push` (backend: tokens + Expo Push API) |
| Migración goose | Ninguna (la tabla es del backend) |
| Dependencia nueva | `expo-notifications` + `expo-device` + `expo-constants` — **pasa por `.agents/mobile/dependencies.md`** |

## Objective

Que el usuario reciba en el teléfono el aviso diario de repaso y pueda abrirlo.

**Historias**:
- Al iniciar sesión, la app pide permiso de notificaciones y registra el dispositivo (una sola vez).
- Cada mañana, si tengo repasos pendientes, me llega "Tienes N preguntas para repasar".
- Toco la notificación y la app se abre en Inicio, donde está el botón "Repasar hoy".
- Si deniego el permiso o el registro falla, la app funciona igual: simplemente no hay push.

**Éxito**: el token queda registrado en el backend y el aviso llega; sin permiso, nada se rompe (D-28: el aviso in-app sigue existiendo).

## Prerrequisitos (bloqueantes, verificar antes de empezar)

1. **Development build**: desde Expo SDK 53 las notificaciones remotas **no funcionan en Expo Go**. Hace falta un dev build (`expo-dev-client` + `eas build --profile development`) o probar en un build de producción. Documentado en N-09.
2. **EAS projectId**: `getExpoPushTokenAsync()` lo necesita. Hoy `app.json` **no** tiene `expo.extra.eas.projectId`. Hay que ejecutar `eas init` (o añadirlo a mano) y commitear el resultado.
3. Sin esos dos, la app debe **seguir funcionando** (registro omitido con un `console.warn`).

## Commands

```bash
pnpm --filter mobile typecheck
pnpm --filter mobile test:run
pnpm --filter mobile lint
pnpm --filter mobile start       # en dev build, no Expo Go
```

## Project structure

```
apps/mobile/
├── app.json                                  # + plugin expo-notifications + extra.eas.projectId
├── src/core/push/
│   ├── notifications.ts                      # NUEVO: permisos, token, canal, handlers
│   ├── notifications.test.ts                 # NUEVO (mocks)
│   └── register.ts                           # NUEVO: orquesta registro/baja contra el backend
├── src/core/api/services/push.service.ts     # NUEVO: registerToken / unregisterToken
├── src/core/api/services/push.service.test.ts# NUEVO
├── src/core/auth/useAuthBootstrap.ts         # (o donde viva el bootstrap) + registrar al entrar
└── src/features/profile/...                  # logout: dar de baja el token
```

## Contrato

### `app.json`

```jsonc
{
  "expo": {
    "plugins": [
      "expo-router",
      ["expo-notifications", { "color": "#22C55E" }]
    ],
    "extra": { "eas": { "projectId": "<uuid del proyecto EAS>" } }
  }
}
```

- Android: crear canal `default` con `Notifications.setNotificationChannelAsync('default', { importance: AndroidImportance.DEFAULT })` **antes** de pedir el token.
- El `color` debe coincidir con `colors.primary` (hoy `#22C55E`).

### `notifications.ts`

```ts
export async function ensurePermission(): Promise<boolean>
export async function getDevicePushToken(): Promise<string | null>   // null si no hay projectId/permiso/emulador
export async function configureNotificationHandling(onOpen: () => void): () => void
```

- `ensurePermission`: `getPermissionsAsync` → si no está concedido, `requestPermissionsAsync`. Devuelve el resultado.
- `getDevicePushToken`: comprueba `Device.isDevice` (en emulador no hay token remoto), obtiene el `projectId` de `Constants.expoConfig?.extra?.eas?.projectId`; si falta → `console.warn('⚠️ Push deshabilitado: falta eas.projectId')` y `null`. Si todo está, `getExpoPushTokenAsync({ projectId })`.
- `configureNotificationHandling`: 
  - `setNotificationHandler` → `{ shouldShowBanner: true, shouldShowList: true, shouldPlaySound: false, shouldSetBadge: false }` (API de SDK 53+).
  - listener de `addNotificationResponseReceivedListener` → `onOpen()`.
  - en arranque en frío, `getLastNotificationResponseAsync()` → si existe, `onOpen()`.
  - devuelve la función de limpieza.

### `register.ts`

```ts
export async function registerForPush(): Promise<void>
export async function unregisterFromPush(): Promise<void>
```

- `registerForPush()`: permiso → token → `POST /push/tokens { token, platform: Platform.OS }`. Guarda el token en memoria para poder darlo de baja. **Nunca lanza**: cualquier fallo se loguea y se sigue.
- `unregisterFromPush()`: si hay token, `POST /push/tokens/unregister { token }`; limpia la memoria. **Nunca lanza**.
- Se llama en el **bootstrap autenticado** (después del login/restauración de sesión), no en cada render.
- `unregisterFromPush()` se llama en el **logout** (junto a `removeToken`).

### Navegación al tocar

`onOpen` → `router.push('/')`. El aviso diario siempre lleva a Inicio, donde está el CTA de la sesión diaria; no se inventa un deep link a una sesión que puede no existir.

### Servicio HTTP

```ts
export async function registerToken(token: string, platform: 'ios' | 'android'): Promise<void>
export async function unregisterToken(token: string): Promise<void>
```

Coinciden con `SPEC-push.md` (`POST /push/tokens`, `POST /push/tokens/unregister`).

## Code style

- Todo lo nativo vive en `src/core/push/`; las pantallas no importan `expo-notifications` directamente.
- Sin emojis como icono (la notificación usa el icono de la app).
- Errores de red del registro: silenciosos para el usuario (no hay acción que tomar), con `console.warn`.

## Testing strategy

Los módulos nativos se mockean (`expo-notifications`, `expo-device`, `expo-constants`).

| Archivo | Casos |
|---|---|
| `push.service.test.ts` | URL, método y body de register/unregister |
| `register.test.ts` | registra con permiso concedido; **no** registra sin permiso; **no** registra sin `projectId`; no lanza si el backend falla; unregister borra el token en memoria y llama al backend |
| `notifications.test.ts` | `ensurePermission` pide permiso solo si no está concedido; `getDevicePushToken` devuelve `null` en emulador/sin projectId; `configureNotificationHandling` registra y limpia listeners |

## Boundaries

- **Always**: degradar sin push; registrar una sola vez por sesión; dar de baja al hacer logout; canal Android creado antes del token; nunca bloquear el arranque por push.
- **Ask first**: añadir `expo-dev-client` (cambia el flujo de desarrollo), configurar el proyecto EAS, pedir permiso en el primer arranque en vez de tras login, programar notificaciones locales.
- **Never**: pedir permiso al abrir la app sin sesión; asumir que `getExpoPushTokenAsync` funciona en Expo Go (N-09); guardar el token en AsyncStorage sin necesidad; mostrar un error al usuario si el registro falla; tocar web o backend.

## Success Criteria

1. Con permiso concedido y `projectId`, tras el login se llama `POST /push/tokens` **una vez** con el token y la plataforma.
2. Sin `projectId`, no se llama al backend y queda un `console.warn`; la app funciona.
3. En emulador, no se registra (no hay token remoto) y no se rompe nada.
4. Al cerrar sesión se llama `POST /push/tokens/unregister` y se limpia el token en memoria.
5. Tocar una notificación abre Inicio; abrir la app desde una notificación en frío también.
6. `typecheck`, `test:run` y `lint` en verde.

## Notas de implementación (2026-10-20)

- **Dependencias**: `expo-notifications@~57.0.23`, `expo-device@~57.0.2` (oficiales, MIT, compatibles con SDK 57); `expo-constants` ya estaba. Gate de `.agents/mobile/dependencies.md` OK.
- **`eas.projectId` pendiente**: no se inventa un UUID; hay que ejecutar `eas init` con una cuenta EAS. Mientras falte, `getDevicePushToken()` devuelve `null` + `console.warn` y la app funciona sin push (degradación probada en tests).
- **Canal Android** `default` con `AndroidImportance.DEFAULT` creado antes de `getExpoPushTokenAsync`.
- **Wiring**: `app/_layout.tsx` dispara `registerForPush()` cuando hay token (guard de sesión en el módulo, no en cada render) y configura el handler/listener; `ProfileScreen` llama `unregisterFromPush()` **antes** de borrar el token. El registro nunca lanza.
- **Navegación**: tocar una notificación (o abrir la app desde una en frío) hace `router.push('/')`, sin deep links inventados.
- Tests: **18 nuevos** → suite mobile **129 tests / 28 archivos**; `typecheck` + `test:run` + `lint` (0 errores) en verde.
- ⚠️ **Pendiente manual**: `eas init` + dev build (Expo Go no soporta push remoto desde SDK 53) + envío real desde el backend. Es el único bloqueo para cerrar MPUSH-1/MPUSH-3.

## Open Questions

| # | Pregunta | Propuesta por defecto |
|---|---|---|
| OQ-1 | ¿Pedir permiso en el primer arranque o tras login? | **Tras login**: sin sesión no tiene sentido y molesta menos. |
| OQ-2 | ¿Notificaciones locales (recordatorio sin backend)? | No: duplicaría el job diario del servidor y se desincronizaría. |
| OQ-3 | ¿Badge en el icono de la app? | No por ahora (`shouldSetBadge: false`); requeriría sincronizar el contador. |
| OQ-4 | ¿Qué pasa si el usuario cambia de cuenta en el mismo dispositivo? | El backend reasigna el token al nuevo usuario (upsert por token). Cubierto por `SPEC-push.md`. |
