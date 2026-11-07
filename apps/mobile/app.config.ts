import type { ConfigContext } from 'expo/config'

/**
 * Config dinámica de Expo.
 *
 * `app.json` sigue siendo la fuente de `version` y `android.versionCode`
 * (`scripts/bump-mobile-version.sh` y `scripts/preflight-apk.sh` los leen de
 * ahí). Aquí solo se sobreescribe el nombre visible para poder compilar
 * variantes por entorno: "DevTest", "DevTest Pruebas", etc.
 *
 * El nombre es build-time: cambiarlo exige recompilar el APK. En local,
 * `make mobile-apk-local APP_NAME="DevTest Pruebas"` ejecuta el prebuild
 * necesario para actualizar `android/app/src/main/res/values/strings.xml`.
 */
const DEFAULT_APP_NAME = 'DevTest'

export default ({ config }: ConfigContext) => ({
  ...config,
  name: process.env.EXPO_PUBLIC_APP_NAME?.trim() || config.name || DEFAULT_APP_NAME,
})
