# APK manual con Expo (EAS cloud)

GuÃ­a para generar un APK instalable **sin publicar en la Play Store** y sin
Android SDK local. El build se dispara a mano; no hay CI de APK.

> Contexto: decisiÃ³n D-46. Docker/entornos: `docker/README.md`.
> Mobile corre fuera de Docker (Metro en el host).

## Modelo mental (leer antes de compilar)

```
APK  â”€â”€HTTPSâ”€â”€â–¶  https://app.tu-dominio.com  â”€â”€â–¶  nginx  â”€â”€â–¶  api  â”€â”€â–¶  Postgres
```

- El APK **solo** habla con el API. **No** se conecta a Postgres.
- `EXPO_PUBLIC_API_BASE_URL` se **compila dentro del bundle**. No es una
  variable de runtime.
- EAS cloud **no lee** `apps/mobile/.env` ni las variables exportadas en tu
  shell. Usa `apps/mobile/eas.json` â†’ `build.preview.env.EXPO_PUBLIC_API_BASE_URL`.
- Fuera de desarrollo la URL debe ser **HTTPS** (D-41).
- Si el backend cambia de dominio, hay que **volver a compilar** el APK.

## Requisitos (una vez)

```bash
node -v                 # >= 20
cd apps/mobile
npx --yes eas-cli@latest login
npx --yes eas-cli@latest init  # agrega extra.eas.projectId a app.json
```

`eas init` requiere una cuenta Expo. El keystore lo gestiona EAS la primera vez
(te ofrece generarlo). GuardÃ¡ las credenciales de la cuenta: sin ellas no podÃ©s
volver a firmar actualizaciones.

## ConfiguraciÃ³n versionada

`apps/mobile/app.json`:

```json
"android": {
  "package": "com.felipe.devtest",
  "versionCode": 1
}
```

`apps/mobile/eas.json`:

```json
"preview": {
  "distribution": "internal",
  "android": { "buildType": "apk" },
  "env": { "EXPO_PUBLIC_API_BASE_URL": "https://app.tu-dominio.com" }
}
```

Cambiar `app.example.com` por el dominio real **antes** de compilar. El
preflight falla si sigue el placeholder.

### Versionado

Antes de **cada** APK nuevo:

1. Subir `expo.version` (visible al usuario, ej. `1.0.1`).
2. Subir `expo.android.versionCode` (entero, siempre creciente; Android no
   instala un APK con `versionCode` menor o igual al instalado).
3. `npx --yes eas-cli@latest build -p android --profile preview`.

## Preflight (obligatorio)

```bash
make mobile-apk
# o, sin compilar todavÃ­a:
bash scripts/preflight-apk.sh
```

El script verifica:

- La URL de `eas.json` es HTTPS y no es `example.com`.
- `app.json` tiene `package` y `versionCode >= 1`.
- `GET <URL>/health` responde `{"status":"ok"}`.
- `GET <URL>/api/v1/auth/status` responde (proxy de nginx OK).
- `EXPO_PUBLIC_API_BASE_URL` del shell, si existe, coincide con `eas.json`
  (EAS cloud ignora el shell).

Si el backend todavÃ­a no estÃ¡ desplegado, el preflight falla: es intencional.

## Build e instalaciÃ³n

```bash
make mobile-apk
```

Equivale a:

```bash
cd apps/mobile
npx --yes eas-cli@latest build -p android --profile preview
```

Al terminar, EAS da una URL del `.apk`. InstalaciÃ³n por USB:

```bash
adb install -r ruta/al/app.apk
```

O descargÃ¡ el APK en el telÃ©fono y habilita "instalar apps desconocidas" para el
navegador/gestor de archivos.

## VerificaciÃ³n en el dispositivo

1. Abrir la app â†’ login (o setup si la DB estÃ¡ vacÃ­a).
2. Hacer una peticiÃ³n real (listar temas, iniciar sesiÃ³n diaria).
3. Confirmar en los logs del API (Dokploy) que llega el trÃ¡fico.
4. Probar modo aviÃ³n + reintento (mobile-offline).

Si la app no conecta:

- VerificÃ¡ que la URL compilada sea la correcta (no `localhost`).
- `curl https://tu-dominio/health` desde la PC.
- HTTPS vÃ¡lido (certificado de Let's Encrypt, no autofirmado).
- El APK de release **no** permite HTTP en claro; no insistas con `http://LAN`.

## Pruebas locales (sin APK)

Para desarrollo usÃ¡ Expo Go contra el stack dev:

```bash
make docker-dev            # API en :8082
# .env local:
EXPO_PUBLIC_API_BASE_URL=http://localhost:8082    # simulador iOS
EXPO_PUBLIC_API_BASE_URL=http://10.0.2.2:8082     # emulador Android
EXPO_PUBLIC_API_BASE_URL=http://192.168.1.50:8082 # dispositivo en la LAN
pnpm --filter mobile start
```

> Si usÃ¡s el backend nativo (`make dev`), el puerto es `:8080` en vez de `:8082`.

Expo Go permite HTTP en desarrollo. Eso **no** aplica al APK de release.

## Â¿Y la base de datos de la nube?

El APK no la necesita. Si vos necesitÃ¡s inspeccionarla desde la PC (DBeaver,
psql, restore), seguÃ­ `docker/README.md` â†’ "Conectarse a la base de datos":
terminal del contenedor, tÃºnel SSH o External Port temporal. **Nunca** publiques
5432 de forma permanente ni metas credenciales de DB en `EXPO_PUBLIC_*`.

## Alternativas (no son el camino principal)

- **EAS local** (`eas build --local`): usa el builder Docker de EAS desde WSL2.
  Reproducible, pero pesado para 7 GB de RAM.
- **Gradle nativo** (`expo prebuild` + `./gradlew assembleRelease`): control
  total; requiere JDK 17 + Android SDK (~8 GB) y keystore propio.
- **CI** (`expo-cicd-workflows`): automatizar el build mÃ¡s adelante; hoy se hace
  a mano para controlar la URL y la versiÃ³n.

## Checklist antes de entregar un APK

- [ ] Backend desplegado y `https://â€¦/health` OK.
- [ ] Admin creado (o setup disponible si es una DB nueva).
- [ ] `eas.json` apunta al dominio correcto (HTTPS, sin `/` final).
- [ ] `expo.version` y `android.versionCode` incrementados.
- [ ] `bash scripts/preflight-apk.sh` en verde.
- [ ] APK instalado y login verificado en un Android real.
- [ ] Si cambiÃ³ el backend: APK recompilado (no reutilizar el anterior).
