# SPEC — client-hardening

| Campo | Valor |
|---|---|
| Module id | `client-hardening` |
| Origen | Revisión de seguridad 2026-10-21 (M4 web, L1, L7) |
| Estado | `ready` |
| Depende de | `http-hardening` (la API ya emite sus cabeceras) |
| Consumido por | — |
| Migración goose | Ninguna |
| Dependencia nueva | Ninguna |
| Cambio de contrato | No |

## Objective

La API queda endurecida en `http-hardening`, pero los dos clientes siguen con
huecos propios:

- Mobile acepta `http://` como base URL también en builds de producción
  (`apps/mobile/src/core/config.ts:1-9`), así que el JWT puede viajar en claro.
- La SPA guarda el JWT en `localStorage` (`apps/web/src/utils/storage.ts`) y no
  define ninguna política de seguridad en el despliegue: sin CSP, un XSS futuro
  exfiltra el token sin fricción.
- El `index.html` no fija `Referrer-Policy` ni hay plantilla de cabeceras para el
  hosting estático.

**Historias**

- Como usuario de la app mobile, quiero que mi token nunca salga por HTTP.
- Como operador de la web, quiero aplicar CSP/HSTS en el borde sin adivinar los
  valores correctos para Vue + Vuetify.
- Como revisor, quiero que el riesgo de `localStorage` esté documentado y
  mitigado, no ignorado.

**Éxito**

- Un build de producción de mobile con `EXPO_PUBLIC_API_BASE_URL=http://…` falla
  al arrancar con un mensaje claro (no envía ni una petición).
- En desarrollo (`__DEV__`) sigue funcionando `http://localhost:8080`.
- Existen plantillas de cabeceras para nginx y para hostings tipo Netlify/CF
  Pages, con la CSP mínima que funciona con Vuetify.
- `apps/web/README.md` documenta la política de token y la TTL recomendada.

## Tech stack

Expo SDK 57 + TypeScript (mobile), Vue 3 + Vite + Vuetify (web). Solo cambios de
configuración, plantillas y documentación.

## Commands

```bash
pnpm --filter mobile typecheck
pnpm --filter mobile test:run
make fe-format && make fe-check && make fe-lint && make fe-test
make fe-build

# verificación de cabeceras en el entorno desplegado
curl -sI https://<tu-dominio>/ | grep -Ei 'content-security-policy|strict-transport|nosniff|frame|referrer'
```

## Project structure

```
apps/mobile/src/core/config.ts                 # guard HTTPS en !__DEV__
apps/mobile/src/core/config.test.ts            # NUEVO
apps/mobile/.env.example                       # nota de HTTPS obligatorio
apps/web/index.html                            # + meta referrer
apps/web/deploy/nginx-security.conf            # NUEVO: plantilla de cabeceras
apps/web/deploy/_headers                       # NUEVO: Netlify / Cloudflare Pages
apps/web/deploy/README.md                      # NUEVO: cómo aplicar y verificar
apps/web/README.md                             # + sección de seguridad
.env.example                                   # comentario de JWT_EXPIRY_HOURS
.agents/web/architecture.md                    # + sección de seguridad del cliente
```

## Contrato

### 1. Mobile: HTTPS obligatorio fuera de desarrollo

```ts
// apps/mobile/src/core/config.ts
const DEFAULT_API_BASE_URL = 'http://localhost:8080'

function readApiBaseUrl(): string {
  const raw = process.env.EXPO_PUBLIC_API_BASE_URL?.trim() || DEFAULT_API_BASE_URL
  const url = raw.replace(/\/+$/, '')

  if (!__DEV__ && !url.startsWith('https://')) {
    throw new Error(
      'EXPO_PUBLIC_API_BASE_URL debe usar HTTPS en builds de producción ' +
        '(recibido: ' + url.replace(/\/\/.*@/, '//<credenciales>@') + ')',
    )
  }
  return url
}

export const API_BASE_URL = readApiBaseUrl()
```

Reglas:

- `__DEV__` lo define React Native/Expo; el build de producción lo deja en
  `false`. Es la única excepción permitida.
- El mensaje **no** debe imprimir la URL completa si lleva credenciales
  (`user:pass@host`); se sanea como arriba.
- No se añade `expo-build-properties` ni `usesCleartextTraffic`: es una
  dependencia nueva y un cambio nativo. Se documenta como follow-up si se hace
  prebuild.
- `apps/mobile/.env.example` añade:

```dotenv
# Producción: HTTPS obligatorio. http:// solo se acepta en desarrollo (__DEV__).
EXPO_PUBLIC_API_BASE_URL=https://api.tu-dominio.com
```

### 2. Web: política de token (decisión, no código)

- Se **mantiene** `localStorage` (migrar a cookie `httpOnly` exige CSRF y cambios
  de contrato: fuera de alcance).
- Mitigaciones obligatorias: CSP estricta en el borde (§3), `removeToken()` en
  cualquier 401 (ya existe en `apps/web/src/api/client.ts:41`), TTL de JWT corta
  recomendada (8–12 h) y `Referrer-Policy: no-referrer`.
- Se documenta en `apps/web/README.md` y `.agents/web/architecture.md` con este
  texto de riesgo: *"un XSS en la SPA puede leer `auth_token`; la defensa es la
  CSP y la revisión de cualquier `v-html`"*.

### 3. Cabeceras de despliegue (plantillas)

`apps/web/deploy/nginx-security.conf`:

```nginx
# Aplicar dentro del server{} que sirve el build de Vite.
# <API_ORIGIN> se sustituye en el despliegue (p. ej. https://api.midominio.com).

add_header Content-Security-Policy "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; connect-src 'self' <API_ORIGIN>; object-src 'none'; base-uri 'self'; frame-ancestors 'none'; form-action 'self'" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-Frame-Options "DENY" always;
add_header Referrer-Policy "no-referrer" always;
add_header Permissions-Policy "geolocation=(), microphone=(), camera=()" always;
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

location / {
    try_files $uri $uri/ /index.html;
}
```

`apps/web/deploy/_headers` (Netlify / Cloudflare Pages):

```
/*
  Content-Security-Policy: default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; connect-src 'self' <API_ORIGIN>; object-src 'none'; base-uri 'self'; frame-ancestors 'none'; form-action 'self'
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY
  Referrer-Policy: no-referrer
  Permissions-Policy: geolocation=(), microphone=(), camera=()
  Strict-Transport-Security: max-age=31536000; includeSubDomains
```

Notas:

- `style-src 'unsafe-inline'` es **necesario**: Vuetify inyecta estilos. No se
  puede endurecer sin migrar a estilos no inline.
- `script-src 'self'` es viable: el build de Vite emite `<script type="module"
  src="...">`, sin scripts inline. Si se añade un script inline (analytics), hay
  que pasar a hashes/nonces en ese momento.
- `connect-src` debe incluir el origen real de la API y **no** usar `*`.
- En desarrollo (`pnpm dev`) estas cabeceras no se aplican (las sirve Vite); no
  se añade CSP por meta tag para no romper HMR.

`apps/web/index.html`:

```html
<meta name="referrer" content="no-referrer" />
```

`apps/web/deploy/README.md` documenta: elegir plantilla según hosting, sustituir
`<API_ORIGIN>`, verificar con `curl -sI`, y repetir tras cada cambio de dominio.

### 4. TTL del JWT (documentación)

- `JWT_EXPIRY_HOURS=24` se mantiene como default de código (no rompe sesiones
  existentes).
- `.env.example` pasa a documentar el rango recomendado:

```dotenv
# Sesión: 8–12 h recomendado en producción. 24 h es el default del código.
JWT_EXPIRY_HOURS=12
```

- `apps/web/README.md` y `.agents/backend/responses.md` (sección JWT) mencionan
  que no hay refresh token ni revocación: un token robado es válido hasta `exp`.

## Testing strategy

| # | Caso | Archivo | Esperado |
|---|---|---|---|
| T1 | `__DEV__ = true` + `http://localhost:8080` | `config.test.ts` | `API_BASE_URL` = `http://localhost:8080` |
| T2 | `__DEV__ = false` + `http://…` | idem | `require('./config')` lanza con mensaje de HTTPS |
| T3 | `__DEV__ = false` + `https://…` | idem | carga y normaliza (sin `/` final) |
| T4 | `__DEV__ = false` + env vacío | idem | lanza (el default es http) |
| T5 | URL con credenciales en el error | idem | el mensaje no contiene `user:pass@` |
| T6 | `index.html` tiene el meta referrer | revisión manual / snapshot | presente |
| T7 | Plantillas de cabeceras | `curl -sI` en staging | 5 cabeceras + CSP con el origen correcto |
| T8 | App web sigue arrancando con CSP aplicada | manual en staging | sin errores de consola por CSP |

> `config.test.ts` usa `jest.isolateModules` + `global.__DEV__` para resetear el
> módulo entre casos (el valor se lee a nivel de import).

## Boundaries

- **Always**: `pnpm --filter mobile typecheck` + `test:run` y los 4 checks de web
  tras cualquier cambio; documentar el riesgo de `localStorage`; sustituir
  `<API_ORIGIN>` antes de desplegar.
- **Ask first**: añadir `expo-build-properties`/`usesCleartextTraffic`; cambiar a
  cookie `httpOnly`; añadir scripts inline (rompe `script-src 'self'`); bajar la
  TTL a < 8 h (impacto UX).
- **Never**: permitir `http://` en producción, poner secretos en `VITE_*` o
  `EXPO_PUBLIC_*` (se compilan al bundle), usar `connect-src *`, desactivar la CSP
  "temporalmente" en producción.

## Success criteria

1. Un build/arranque de mobile en modo producción con URL `http://` falla con un
   error explícito y no realiza peticiones.
2. En desarrollo la app sigue apuntando a `http://localhost:8080`.
3. `apps/web/deploy/` contiene las dos plantillas y un README con el
   procedimiento de sustitución y verificación.
4. En staging, `curl -sI` muestra CSP, HSTS, `nosniff`, `X-Frame-Options` y
   `Referrer-Policy`, y la SPA carga sin errores de CSP en consola.
5. `apps/web/README.md` documenta la política de token y la TTL recomendada.
6. Los 4 checks de web y los 2 de mobile están verdes.

## Open questions

- ¿Qué hosting sirve la SPA (nginx propio, Netlify, Cloudflare Pages, Vercel)?
  Se entregan dos plantillas; hay que quedarse con una y borrar la otra o dejarlas
  como referencia (propuesto: dejar ambas, el README indica cuál aplicar).
- ¿Se hace prebuild nativo en el corto plazo? Si sí, añadir
  `expo-build-properties` con `usesCleartextTraffic: false` en una spec aparte.
- ¿Se migra el token a cookie `httpOnly`? Requiere CSRF + cambios de contrato
  (nueva spec, fuera de alcance).
