# TASKS — Engagement & Retención

> Backend: Go + Gin + GORM (Postgres). Frontend: Vue 3 + Vite + Vuetify + TanStack Query + Pinia.
> Decisiones: notificaciones **solo in-app** · gamificación **solo individual** · 4 fases.
> Convención: modelos → AutoMigrate en `internal/database/database.go` · `make swagger` tras cada endpoint ·
> `make check` (backend) · `make fe-check` + `make fe-lint` (frontend).

## Estado general
- [ ] Fase 1 — Gamificación individual
- [ ] Fase 2 — Repaso inteligente (temas débiles)
- [ ] Fase 3 — Notificaciones in-app
- [ ] Fase 4 — Reportes semanal/mensual

## Pendiente anterior (Vitest)
- [ ] **5.5** (Opcional) Playwright E2E — flujo `status → setup → login → dashboard → sessions`

---

### Fase 1 — Gamificación individual
#### Backend
- [ ] **1.1** Modelo `internal/models/user_stats.go` → `UserStats` (PK `UserID`, `Xp`, `Level`, `CurrentStreak`, `BestStreak`, `LastActivityDate *time.Time`, `TotalAnswers`, `TotalCorrect`, `DailyGoal int` default 10)
- [ ] **1.2** Modelo `internal/models/user_badge.go` → `UserBadge` (ID, `UserID`, `BadgeKey`, `EarnedAt`, unique `(user_id,badge_key)`)
- [ ] **1.3** Registrar ambos modelos en `database.go` AutoMigrate
- [ ] **1.4** Módulo `internal/modules/gamification/` (store/service/handler/routes, plantilla: `progress`)
- [ ] **1.5** `catalog.go` — catálogo de logros como `const` (no tabla): `first_session`, `streak_3`, `streak_7`, `streak_30`, `perfect_session`, `answers_100`, `topics_mastered_5`, `fast_solver`
- [ ] **1.6** Reglas XP/racha/nivel: acierto +10 XP (`hard` ×1.5, `easy` ×0.5); sesión completada +25; `level = 1 + floor(sqrt(xp/100))`; racha = días consecutivos con ≥1 sesión completada
- [ ] **1.7** Endpoints: `GET /gamification/summary`, `GET /gamification/badges`, `PUT /gamification/goal`
- [ ] **1.8** Hooks: `sessions/service.go` `Answer()` → `RecordAnswer()` y `Finish()` → `RecordSessionFinished()` (fire-and-forget); wiring en `server.go`
- [ ] **1.9** Tests `service_test.go` (store mockeado) + `handler_test.go`

#### Frontend
- [ ] **1.10** `types/gamification.types.ts` + `api/services/gamification.service.ts` + `queries/gamification.queries.ts`
- [ ] **1.11** `features/gamification/pages/AchievementsPage.vue` (grid logros ganados/bloqueados) + ruta
- [ ] **1.12** Nav: item "Logros" (`mdi-trophy`, `userOnly`) en sección **Estudio** de `DefaultLayout.vue`
- [ ] **1.13** Widget racha/XP/nivel + barra meta diaria en `DashboardPage.vue`; chip XP/racha en el `v-app-bar`
- [ ] **1.14** Tests (service/queries/componentes) + `make fe-check` + `make fe-lint`

---

### Fase 2 — Repaso inteligente (temas débiles)
#### Backend
- [ ] **2.1** `progress/store.go` → `FindTopicMastery(userID)` (agregación `SessionAnswer` JOIN `question_topics` GROUP BY tema) + `CountDue(userID)`
- [ ] **2.2** `progress/service.go` → `TopicMastery()` con `isWeak = accuracy < 60%` y `attempts >= 5`
- [ ] **2.3** `response.go` → `TopicMasteryResponse` / `TopicMasteryItem`
- [ ] **2.4** Endpoint `GET /progress/topics-mastery` + anotaciones Swagger
- [ ] **2.5** Tests service/store

#### Frontend
- [ ] **2.6** `types/progress.types.ts` ampliado + service/query
- [ ] **2.7** `ProgressPage.vue` — nueva pestaña "Temas" con barras de dominio + badge "Débil"
- [ ] **2.8** Tests + `make fe-check` + `make fe-lint`

---

### Fase 3 — Notificaciones in-app
#### Backend
- [ ] **3.1** Modelo `internal/models/notification.go` → `Notification` (ID, `UserID`, `Type`, `Title`, `Body`, `IsRead`, `CreatedAt`); registrar en AutoMigrate
- [ ] **3.2** Módulo `internal/modules/notifications/` (store/service/handler/routes)
- [ ] **3.3** Endpoints: `GET /notifications`, `POST /notifications/:id/read`, `POST /notifications/read-all`, `GET /notifications/unread-count`
- [ ] **3.4** Generador: job diario crea "Hoy tienes N preguntas para repasar" (desde `next_review_at <= now`)
- [ ] **3.5** Scheduler en `server.go` (goroutine + `github.com/robfig/cron/v3` — validar con `dependencies.md` antes de `go get`)
- [ ] **3.6** Tests

#### Frontend
- [ ] **3.7** `types/notification.types.ts` + service + queries
- [ ] **3.8** Campana en `v-app-bar` con badge no-leídos + `features/notifications/pages/NotificationsPage.vue` + ruta
- [ ] **3.9** Tests + `make fe-check` + `make fe-lint`

---

### Fase 4 — Reportes semanal/mensual
#### Backend
- [ ] **4.1** Módulo `internal/modules/reports/` (store/service/handler/routes)
- [ ] **4.2** Endpoint `GET /reports/summary?period=week|month` → `{ sessions, questionsAnswered, accuracy, timeSpentMs, topTopics[], weakTopics[], xpGained, streakDays, dailyActivity[] }`
- [ ] **4.3** Agregaciones desde `Session` + `SessionAnswer` + `question_topics` + stats de Fase 1
- [ ] **4.4** Tests

#### Frontend
- [ ] **4.5** `types/report.types.ts` + service + queries
- [ ] **4.6** `features/reports/pages/ReportsPage.vue` (selector semana/mes, KPIs, barras por tema) + ruta
- [ ] **4.7** Nav: item "Informes" (`mdi-file-chart`, `userOnly`) en sección **Estudio**
- [ ] **4.8** Tests + `make fe-check` + `make fe-lint`

---

## Flujo de trabajo
1. Marcar tarea como `in_progress` antes de empezar.
2. Implementar + verificar (`make check` / `make fe-check` + `make fe-lint`; tests).
3. Marcar completada y seguir con la siguiente.
4. Al terminar, commit manual (no automático).

## Notas
- Logros/gamificación: catálogo en código (const), no tabla. Criterios reevaluados en cada `RecordAnswer`/`RecordSessionFinished`.
- Notificaciones in-app: sin secretos ni env vars nuevas; scheduler sin estado compartido.
- Reportes: dependen de datos de Fase 1 (XP) y Fase 2 (temas débiles), por eso van al final.
- Dependencias nuevas (Fase 3: `robfig/cron/v3`) → pasar por `.agents/backend/dependencies.md` antes de instalar.
