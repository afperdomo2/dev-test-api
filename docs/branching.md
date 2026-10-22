# Ramas y releases

Modelo para un solo desarrollador con despliegue automático desde `main`.

```
develop  ──(PR semanal)──►  main  ──►  Dokploy (producción)
   ▲
   └── trabajo diario + backup
```

- **`develop`**: rama de trabajo. Aquí viven los commits diarios y el trabajo a
  medias. **No despliega.**
- **`main`**: producción. Solo recibe cambios por PR desde `develop`. Cada merge
  dispara el deploy.
- **Tags**: marcan qué versión está en producción (`v0.4.0`).

## Reglas

| Rama | Reglas |
|---|---|
| `main` | PR obligatorio, checks obligatorios (`backend`, `shared`, `web`, `gitleaks`), rama actualizada, squash/lineal, sin force-push, sin borrado, **sin bypass** |
| `develop` | Sin PR obligatorio; force-push y borrado bloqueados |

Sin bypass significa que **nadie** (ni el dueño) puede saltarse el PR o los
checks. Si GitHub Actions cae, la salida es desactivar la regla temporalmente,
no un atajo permanente.

## Flujo diario

```bash
git switch develop
git pull

# ... commits del día, aunque el trabajo quede a medias ...

git push            # backup; NO despliega
```

Cada push a `develop` ejecuta CI (`Backend`, `Frontend`, `Secrets`) según los
archivos tocados.

## Release (semanal)

La PR se abre **el mismo día** que se va a mergear, no se deja abierta.

```bash
git switch develop
git pull
gh pr create --base main --head develop --fill

# cuando los checks estén verdes
gh pr diff --stat          # tamaño real del bloque
gh pr merge --squash       # dispara el deploy

# tras verificar producción
gh release create v0.4.0 --target main --generate-notes
```

Checklist antes de mergear:

- [ ] Los 4 checks en verde.
- [ ] `gh pr diff --stat` con un tamaño manejable.
- [ ] Si hay migraciones: backup de la DB y migración compatible hacia atrás.
- [ ] Sabes cómo revertir (ver abajo).

## Hotfix

Producción rota: se corrige desde `main`, no desde `develop` (para no publicar el
trabajo pendiente).

```bash
git switch main && git pull
git switch -c fix/login-500
# fix + commit
git push -u origin fix/login-500
gh pr create --fill
gh pr merge --squash --delete-branch

# volver a sincronizar develop
git switch develop
git merge main
git push
```

## Rollback

| Qué | Cómo |
|---|---|
| Código | `git revert -m 1 <merge>` en `main` + push (despliega la reversión) |
| Imagen | Redeploy del deployment/tag anterior en Dokploy |
| DB | **No** se revierte con un tag. Restaurar backup |

Regla: un tag no deshace migraciones. Para cambios de esquema, migraciones
*expand/contract* (añadir → usar → borrar en releases distintas).

## Versiones

SemVer: `vMAJOR.MINOR.PATCH`.

| Cambio | Versión |
|---|---|
| Funcionalidad nueva | `MINOR` |
| Corrección | `PATCH` |
| Cambio incompatible (API/DB) | `MAJOR` |

## CI y caché

- Push a `develop` → 1 ejecución por push (sin PR abierta).
- PR `develop` → `main` → 1 ejecución adicional el día del release.
- La caché de `develop` **no** se reutiliza en la PR (las cachés están aisladas
  por rama; la PR restaura la de `main`). GitHub borra cachés sin usar a los 7
  días.
- Los checks obligatorios corren siempre en PRs (sin filtro `paths`); si un PR
  solo toca docs, igualmente se validan.

## Dokploy

- App apuntando a la rama `main`.
- Auto-deploy activado en cada push a `main`.
- Healthcheck: `GET /health`.
- Las migraciones goose corren al arrancar el contenedor.
