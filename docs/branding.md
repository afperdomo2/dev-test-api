# Branding — icono web y mobile

Fuente única: `brand/icon.svg` (tile redondeado) y `brand/mark.svg` (monograma
suelto). Todo lo que consume la web y la app mobile es **generado**: no se edita
a mano.

## Paleta

| Token | Hex | Uso |
|---|---|---|
| Marca | `#0F766E` | fondo del tile, `theme_color`, background del adaptive icon |
| Monograma | `#F0FDFA` | trazo de la D |
| Borde interior | `#5EEAD4` al 28% | hairline del tile (solo visible en tamaños grandes) |

El icono **no** sigue `VITE_APP_COLOR` ni los tokens de mobile: es branding
estático y se comparte entre entornos. El nombre instalado del PWA también queda
fijo (`DevTest`); `VITE_APP_NAME` solo cambia el título dentro de la app.

## Regenerar

```bash
make icons
```

Requiere Google Chrome/Chromium (render de SVG), Python 3 y Pillow (armado del
`.ico` y normalización de PNG). El script es determinista: correrlo dos veces
seguidas no cambia los bytes. Salida:

### Web (`apps/web/public/`)

| Archivo | Tamaño | Uso |
|---|---|---|
| `favicon.svg` | vector | pestaña en navegadores modernos |
| `favicon.ico` | 16/24/32/48 | fallback y marcadores |
| `apple-touch-icon.png` | 180×180 | pantalla de inicio iOS (opaco) |
| `icon-192.png` | 192×192 | PWA `any` |
| `icon-512.png` | 512×512 | PWA `any` |
| `maskable-512.png` | 512×512 | PWA `maskable` (con margen seguro) |
| `site.webmanifest` | — | manifest PWA (estático, no lo genera el script) |

### Mobile (`apps/mobile/assets/`)

| Archivo | Tamaño | Uso |
|---|---|---|
| `icon.png` | 1024×1024 | icono legacy / iOS (opaco, full-bleed) |
| `android-icon-foreground.png` | 1024×1024 | capa foreground del adaptive icon (transparente) |
| `android-icon-monochrome.png` | 432×432 | capa monochrome (transparente) |
| `favicon.png` | 48×48 | Expo web (no se usa en el APK) |
| `splash-icon.png` | 1024×1024 | reservado por si se configura splash |

`android-icon-background.png` se eliminó: `app.json` usa
`android.adaptiveIcon.backgroundColor: "#0F766E"`.

## Zonas seguras

- **Adaptive icon Android**: lienzo de 108dp y zona segura de 66dp. La D se
  renderiza al **77%** del lienzo para que ningún recorte la corte.
- **Maskable PWA**: zona segura del 80% del ancho; se renderiza al **95%**.
- **Apple touch**: opaco y full-bleed; iOS aplica su propia máscara. No debe
  tener transparencia (se rellenaría de negro).

## PWA

El manifest se sirve desde `apps/web/public/site.webmanifest` y se enlaza en
`apps/web/index.html` junto con `theme-color`, `favicon` y `apple-touch-icon`.

nginx 1.29 **no** trae el MIME de `.webmanifest` en `mime.types`. Por eso
`docker/web/nginx.conf` tiene un `location = /site.webmanifest` con
`default_type application/manifest+json`; sin eso se sirve como
`application/octet-stream` y el navegador ignora el manifest. La CSP no cambia:
`manifest-src` cae en `default-src 'self'` y los iconos entran por
`img-src 'self' data:`.

La instalación real se prueba en `localhost` o HTTPS (DevTools → Application →
Manifest). Prod-local en `:8088` es HTTP: el manifest parsea, pero el navegador
no ofrece instalar.

## Rebuild mobile

`android/` es prebuild generado (gitignored) y no lee `app.json` para los
iconos. Cambiar assets **sin** prebuild deja el APK con el icono viejo.

```bash
make mobile-bump VERSION=1.0.2          # icono = PATCH; versionCode sube siempre
make mobile-apk-local APP_NAME="DevTest"  # fuerza prebuild y respalda el keystore
```

Si corrés `expo prebuild` a mano, respaldá y restaurá
`android/app/debug.keystore` (misma firma = se instala encima). Un cambio de
iconos no se cierra entregando dos APKs con el mismo `versionCode`; ver
`docs/mobile-apk.md` → "Versionado".

## Flujo para cambiar el icono

1. Editar `brand/icon.svg` y/o `brand/mark.svg`.
2. `make icons`.
3. Revisar el diff de `apps/web/public/` y `apps/mobile/assets/` (mismo cambio).
4. Web: `make fe-dev` y verificar pestaña + Manifest.
5. Mobile: bump, prebuild y APK (sección anterior).
