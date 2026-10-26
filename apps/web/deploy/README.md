# Plantillas de despliegue web

El hosting canónico del proyecto es **nginx dentro del contenedor** (`docker/web/nginx.conf`),
que incluye `nginx-security.conf`. Estas plantillas existen para no improvisar las
cabeceras si algún día se sirve la SPA en otro hosting.

## nginx propio (contenedor)

- `nginx-security.conf` se copia a `/etc/nginx/snippets/security-headers.conf` y se
  incluye en el `server{}` que sirve `dist/`.
- `connect-src 'self'` es correcto porque la SPA y el API comparten dominio
  (nginx proxya `/api` a `api:8080`). Si separás el API, agregá su origen.
- Verificación:

```bash
curl -sI https://app.tu-dominio.com/ | grep -Ei \
  'content-security-policy|strict-transport|nosniff|frame|referrer|permissions'
```

## Netlify / Cloudflare Pages

Usar `_headers` (si se agrega en el futuro) con las mismas cabeceras y sustituir
`connect-src 'self'` por el origen real del API. **No** usar `connect-src *`.

## Reglas

- `style-src 'unsafe-inline'` es obligatorio con Vuetify; no endurecerlo sin migrar
  a estilos no inline.
- Un script inline nuevo (analytics) rompe `script-src 'self'`: en ese momento hay
  que pasar a hashes/nonces.
- Tras cada cambio de dominio, repetir `curl -sI` y revisar la consola del navegador.
