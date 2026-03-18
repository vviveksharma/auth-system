# API Endpoints Roadmap

> A phased rollout plan covering **50 endpoints** across 10 weeks.

---

## Table of Contents

- [Phase 1 — Foundation](#phase-1--foundation-week-1-2)
- [Phase 2 — Core Tracking](#phase-2--core-tracking-week-3-4)
- [Phase 3 — Analytics](#phase-3--analytics-week-5-6)
- [Phase 4 — Alerts & Monitoring](#phase-4--alerts--monitoring-week-7-8)
- [Phase 5 — Polish & Extras](#phase-5--polish--extras-week-9-10)

---

## Phase 1 — Foundation `Week 1-2`

> **15 Endpoints** — Core auth, organizations, projects, and team management.

### Authentication

| Method | Endpoint                |
| ------ | ----------------------- |
| `POST` | `/api/v1/auth/register` |
| `POST` | `/api/v1/auth/login`    |
| `POST` | `/api/v1/auth/logout`   |

### Organizations

| Method   | Endpoint                        |
| -------- | ------------------------------- |
| `POST`   | `/api/v1/organizations`         |
| `GET`    | `/api/v1/organizations`         |
| `GET`    | `/api/v1/organizations/{orgId}` |
| `PUT`    | `/api/v1/organizations/{orgId}` |
| `DELETE` | `/api/v1/organizations/{orgId}` |

### Projects

| Method   | Endpoint                                             |
| -------- | ---------------------------------------------------- |
| `POST`   | `/api/v1/organizations/{orgId}/projects`             |
| `GET`    | `/api/v1/organizations/{orgId}/projects`             |
| `GET`    | `/api/v1/organizations/{orgId}/projects/{projectId}` |
| `PUT`    | `/api/v1/organizations/{orgId}/projects/{projectId}` |
| `DELETE` | `/api/v1/organizations/{orgId}/projects/{projectId}` |

### Team Management

| Method | Endpoint                                |
| ------ | --------------------------------------- |
| `POST` | `/api/v1/organizations/{orgId}/members` |
| `GET`  | `/api/v1/organizations/{orgId}/members` |

---

After phase one make the delete a shared service as it would connected with the users and complete flow 

## Phase 2 — Core Tracking `Week 3-4`

> **10 Endpoints** ⭐ **CRITICAL** — The heart of the system.

### API Keys

| Method   | Endpoint                                               |
| -------- | ------------------------------------------------------ |
| `POST`   | `/api/v1/organizations/{orgId}/api-keys`               |
| `GET`    | `/api/v1/organizations/{orgId}/api-keys`               |
| `GET`    | `/api/v1/organizations/{orgId}/api-keys/{keyId}`       |
| `DELETE` | `/api/v1/organizations/{orgId}/api-keys/{keyId}`       |
| `GET`    | `/api/v1/organizations/{orgId}/api-keys/{keyId}/usage` |

### Usage Tracking ⭐⭐⭐ MOST IMPORTANT

> These are **THE** core endpoints — everything else supports these!

| Method | Endpoint              |
| ------ | --------------------- |
| `POST` | `/api/v1/track`       |
| `POST` | `/api/v1/track/batch` |

### Real-Time Stats

| Method | Endpoint                                       |
| ------ | ---------------------------------------------- |
| `GET`  | `/api/v1/organizations/{orgId}/stats/realtime` |
| `GET`  | `/api/v1/organizations/{orgId}/stats/today`    |
| `GET`  | `/api/v1/projects/{projectId}/stats/realtime`  |

---

## Phase 3 — Analytics `Week 5-6`

> **12 Endpoints** — Daily/monthly stats, charts, trends, and top lists.

### Daily Statistics

| Method | Endpoint                                            |
| ------ | --------------------------------------------------- |
| `GET`  | `/api/v1/organizations/{orgId}/analytics/daily`     |
| `GET`  | `/api/v1/projects/{projectId}/analytics/daily`      |
| `GET`  | `/api/v1/organizations/{orgId}/analytics/providers` |
| `GET`  | `/api/v1/projects/{projectId}/analytics/providers`  |

### Monthly Statistics

| Method | Endpoint                                          |
| ------ | ------------------------------------------------- |
| `GET`  | `/api/v1/organizations/{orgId}/analytics/monthly` |
| `GET`  | `/api/v1/projects/{projectId}/analytics/monthly`  |

### Charts & Trends

| Method | Endpoint                                                |
| ------ | ------------------------------------------------------- |
| `GET`  | `/api/v1/organizations/{orgId}/analytics/trends`        |
| `GET`  | `/api/v1/organizations/{orgId}/analytics/breakdown`     |
| `GET`  | `/api/v1/projects/{projectId}/analytics/trends`         |
| `GET`  | `/api/v1/projects/{projectId}/analytics/cost-over-time` |

### Top Lists

| Method | Endpoint                                                |
| ------ | ------------------------------------------------------- |
| `GET`  | `/api/v1/organizations/{orgId}/analytics/top-projects`  |
| `GET`  | `/api/v1/organizations/{orgId}/analytics/top-providers` |

---

## Phase 4 — Alerts & Monitoring `Week 7-8`

> **8 Endpoints** — Configurable alerts and budget tracking.

### Alerts

| Method   | Endpoint                                         |
| -------- | ------------------------------------------------ |
| `POST`   | `/api/v1/organizations/{orgId}/alerts`           |
| `GET`    | `/api/v1/organizations/{orgId}/alerts`           |
| `GET`    | `/api/v1/organizations/{orgId}/alerts/{alertId}` |
| `PUT`    | `/api/v1/organizations/{orgId}/alerts/{alertId}` |
| `DELETE` | `/api/v1/organizations/{orgId}/alerts/{alertId}` |

### Budgets

| Method | Endpoint                                        |
| ------ | ----------------------------------------------- |
| `POST` | `/api/v1/organizations/{orgId}/budgets`         |
| `GET`  | `/api/v1/organizations/{orgId}/budgets`         |
| `GET`  | `/api/v1/organizations/{orgId}/budgets/current` |

---

## Phase 5 — Polish & Extras `Week 9-10`

> **5 Endpoints** — Exports, reports, settings, and health check.

### Export & Reports

| Method | Endpoint                                        |
| ------ | ----------------------------------------------- |
| `GET`  | `/api/v1/organizations/{orgId}/export/csv`      |
| `GET`  | `/api/v1/organizations/{orgId}/reports/monthly` |

### Settings

| Method | Endpoint                                 |
| ------ | ---------------------------------------- |
| `GET`  | `/api/v1/organizations/{orgId}/settings` |
| `PUT`  | `/api/v1/organizations/{orgId}/settings` |

### Health Check

| Method | Endpoint  |
| ------ | --------- |
| `GET`  | `/health` |
