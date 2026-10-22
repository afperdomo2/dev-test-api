# Security Policy

## Versiones soportadas

Solo se da soporte de seguridad a la rama `main`. No hay releases mantenidas en
paralelo: si encuentras un problema, comprueba primero que sigue presente en el
último commit de `main`.

| Versión | Soporte |
|---|---|
| `main` | ✅ |
| cualquier otra rama o commit | ❌ |

## Cómo reportar una vulnerabilidad

**No abras un issue público.** Usa una de estas vías:

1. **GitHub Security Advisories** (preferido): pestaña *Security* → *Advisories*
   → *Report a vulnerability*. El reporte queda privado hasta la divulgación.
2. Si esa opción no está disponible, abre un issue **sin detalles** pidiendo un
   canal privado y espera respuesta antes de compartir nada.

Incluye en el reporte:

- Versión o commit afectado.
- Pasos de reproducción y resultado esperado vs. real.
- Impacto (confidencialidad, integridad, disponibilidad) y requisitos previos
  (¿hace falta sesión? ¿rol admin?).
- PoC **sin datos reales** ni credenciales de terceros.

Nunca envíes secretos, tokens ni datos personales en el reporte.

## Plazos

| Etapa | Compromiso |
|---|---|
| Acuse de recibo | ≤ 72 h |
| Triage inicial (severidad + alcance) | ≤ 7 días |
| Corrección de severidad alta | ≤ 30 días |
| Divulgación coordinada | ≤ 90 días desde el reporte |

Si el problema está siendo explotado, dilo en el reporte: se prioriza.

## Alcance

**Dentro de alcance**

- Autenticación y autorización (JWT, roles, acceso a objetos de otros usuarios).
- Inyección (SQL, plantillas, comandos) y XSS/CSRF.
- Fuga de información sensible (tokens, hashes, datos de otros usuarios).
- Denegación de servicio trivial (payloads sin límite, fuerza bruta).
- Fallos de configuración con impacto de seguridad en el despliegue.

**Fuera de alcance**

- Ingeniería social y phishing.
- Ataques que requieren acceso físico al dispositivo.
- Ausencia de MFA o de rotación automática de tokens (mejora conocida, no
  vulnerabilidad).
- Reportes generados solo por escáneres, sin PoC ni impacto demostrado.
- DoS volumétrico a nivel de red (se mitiga en el borde, no en la aplicación).

## Rotación si se commitea una credencial

1. **Revocar** la credencial en el proveedor (no esperar a limpiar el repo).
2. **Rotar**: generar una nueva y actualizarla en `.env` (nunca en el repo).
3. **Purgar la historia** (`git filter-repo` o BFG) y forzar push coordinado.
4. **Auditar el uso** de la credencial filtrada (logs del proveedor) desde la
   exposición.
5. Abrir un advisory privado si la credencial estuvo pública.

`SETUP_TOKEN` y `JWT_SECRET` se rotan igual. Rotar `JWT_SECRET` invalida todas
las sesiones activas (comportamiento esperado).

## Detección de secretos

El repo es público y escanea secretos con **gitleaks** en cada PR y push a
`main` (`.github/workflows/secrets.yml`, reglas en `.gitleaks.toml`).

- Para ejecutarlo en local (o instalar el binario), ver **README → Seguridad**.
- Recomendado antes de commitear: `gitleaks protect --staged --redact`.
- No hay hook de pre-commit por defecto: exigiría tener gitleaks instalado en
  todas las máquinas y rompería el `make setup` actual.

### Cobertura de GitHub (verificada el 2026-10-21)

| Ajuste | Estado | Cómo se verificó |
|---|---|---|
| Secret scanning | ✅ On | `repos/{owner}/{repo}` → `security_and_analysis` |
| Push protection | ✅ On | idem |
| Dependabot alerts | ✅ On | `GET /vulnerability-alerts` → 204 |
| Dependabot security updates | ✅ On | `GET /automated-security-fixes` → 200 |
| Non-provider patterns | ❌ No disponible | El PATCH responde 200 pero no aplica (requiere GitHub Secret Protection) |
| Validity checks | ❌ No disponible | idem |

**Limitación importante de Push Protection en este repo**

GitHub solo bloquea por defecto los **patrones de proveedor** (tokens de AWS,
Stripe, GitHub…). Los **patrones genéricos** —claves privadas, cadenas de
conexión— necesitan *GitHub Secret Protection* (plan de pago) y no se pueden
activar por API en un repo público de usuario.

Se verificó empíricamente con 4 canarios en ramas desechables (RSA PKCS#1, EC,
RSA PKCS#8 y una API key de Google falsa): **ninguno fue bloqueado**, y no se
creó ninguna alerta (`secret-scanning/alerts` → `[]`). Los canarios se borraron
inmediatamente (rama local y remota).

**Consecuencia práctica**: para claves privadas y secretos genéricos, la única
barrera real de este repo es **gitleaks en CI** (`.github/workflows/secrets.yml`)
y el chequeo manual `gitleaks protect --staged`. Por eso el workflow no es
opcional: cubre justo lo que GitHub no bloquea.

Si un escaneo falla por un falso positivo, se corrige el patrón o se añade la
ruta a la allowlist **solo si no es código de producción** (`internal/`,
`apps/*/src/`), dejando el motivo en el propio `.gitleaks.toml`.

## Riesgos aceptados

| Fecha | Elemento | Motivo | Revisión |
|---|---|---|---|
| 2026-10-21 | `braces@3.0.3` (GHSA-vfj7-8cjw-p6xm) | Sin versión parcheada; alcanzable solo desde `expo/@expo/cli` en build | 2027-01-21 |
| 2026-10-21 | `node-forge@1.4.0` (GHSA-86w9-cpqp-85rv) | Sin versión parcheada; alcanzable solo desde `@expo/cli` en build | 2027-01-21 |
| 2026-10-21 | `sprintf-js@1.0.3` (GHSA-hp3w-g68c-fv3c) | Cadena de `jest` (solo test) | 2027-01-21 |
| 2026-10-21 | JWT en `localStorage` (web) | Migrar a cookie `httpOnly` exige CSRF y cambios de contrato; se mitiga con CSP + TTL | 2027-01-21 |

Cada nueva entrada en `pnpm.auditConfig.ignoreCves` debe añadirse aquí con motivo
y fecha de revisión. **No** se ignoran CVEs de código alcanzable en producción.

## Notas operativas

- **`SETUP_TOKEN`**: si está vacío, `POST /api/v1/auth/setup` responde 403 y la
  instancia no se puede inicializar desde la API. Genéralo con
  `openssl rand -hex 32`, ponlo en `.env` y retíralo después de crear el admin.
- **Variables de entorno**: `.env` nunca se commitea (`.gitignore` cubre
  `.env`, `.env.*`, `*.pem`, `*.key`, `*.p12` y `*.jks`).

## Checklist de configuración (GitHub)

- [x] Secret scanning: **On** — verificado por API el 2026-10-21
- [x] Push protection: **On** — solo patrones de proveedor (ver limitación)
- [x] Dependabot alerts: **On** — verificado por API el 2026-10-21
- [x] Dependabot security updates: **On** — verificado por API el 2026-10-21
- [ ] Canario de patrón genérico rechazado — **no aplica** en este plan: requiere GitHub Secret Protection. La cobertura de claves privadas la da gitleaks en CI.
