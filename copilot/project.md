## **Overview**

This document contains all API endpoints required for the Projects feature, including the project list page and project detail modal.

---

## **API Endpoints**

---

## **1. GET /api/v1/organizations/{orgId}/projects**

**Purpose:** List all projects for an organization

**Query Parameters:**

| Parameter     | Type    | Required | Default    | Description           |
| ------------- | ------- | -------- | ---------- | --------------------- |
| `page`        | INTEGER | No       | 1          | Page number           |
| `limit`       | INTEGER | No       | 20         | Results per page      |
| `search`      | STRING  | No       | null       | Search query          |
| `environment` | STRING  | No       | null       | Filter by environment |
| `sort_by`     | STRING  | No       | created_at | Sort field            |
| `sort_order`  | STRING  | No       | desc       | Sort order: asc, desc |

**Request Example:**

```http
GET /api/v1/organizations/org_abc123/projects?page=1&limit=20
Authorization: Bearer {jwt_token}
```

**Response (200):**

```json
{
  "organization_id": "org_abc123",
  "projects": [
    {
      "id": "proj_abc123",
      "name": "Production API",
      "description": "Main production API for customer-facing applications",
      "environment": "production",
      "icon": "rocket",
      "today_stats": {
        "requests": 9800,
        "cost_usd": 98.2,
        "tokens": 3500000,
        "errors": 78,
        "error_rate": 0.8
      },
      "month_stats": {
        "requests": 320000,
        "cost_usd": 7200.0,
        "tokens": 115000000
      },
      "api_keys": {
        "active_count": 5,
        "total_count": 7
      },
      "created_at": "2025-12-01T00:00:00Z",
      "updated_at": "2026-03-07T10:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 3,
    "per_page": 20
  }
}
```

---

## **2. GET /api/v1/projects/{projectId}/details**

**Purpose:** Get detailed project information (for the detail modal)

**Query Parameters:**

| Parameter           | Type    | Required | Default | Description                         |
| ------------------- | ------- | -------- | ------- | ----------------------------------- |
| `date`              | DATE    | No       | today   | Date for "today" stats (YYYY-MM-DD) |
| `include_providers` | BOOLEAN | No       | true    | Include provider breakdown          |
| `include_alerts`    | BOOLEAN | No       | true    | Include alert configuration         |

**Request Example:**

```http
GET /api/v1/projects/proj_ghi789/details?date=2026-03-08&include_providers=true&include_alerts=true
Authorization: Bearer {jwt_token}
```

**Response (200):**

```json
{
  "project": {
    "id": "proj_ghi789",
    "organization_id": "org_abc123",
    "name": "Development",
    "description": "Development environment for testing new features",
    "environment": "development",
    "icon": "code",
    "created_at": "2026-01-05T00:00:00Z",
    "updated_at": "2026-03-08T00:00:00Z"
  },
  "today_stats": {
    "date": "2026-03-08",
    "requests": {
      "total": 150,
      "successful": 148,
      "failed": 2,
      "error_rate": 1.3
    },
    "cost_usd": 1.5,
    "tokens": {
      "total": 45000,
      "input": 20000,
      "output": 25000
    },
    "performance": {
      "avg_duration_ms": 312,
      "p50_duration_ms": 280,
      "p95_duration_ms": 520,
      "p99_duration_ms": 890
    }
  },
  "month_stats": {
    "year": 2026,
    "month": 3,
    "requests": {
      "total": 4500,
      "successful": 4455,
      "failed": 45,
      "error_rate": 1.0
    },
    "cost_usd": 45.0,
    "tokens": {
      "total": 1400000,
      "avg_per_request": 311
    },
    "performance": {
      "avg_duration_ms": 298,
      "p95_duration_ms": 480
    }
  },
  "providers_usage": [
    {
      "provider": "openai",
      "provider_label": "OpenAI",
      "models": ["GPT-4", "GPT-3.5"],
      "requests": 245678,
      "cost_usd": 4567.89,
      "cost_share": 78.0,
      "tokens": 890000
    },
    {
      "provider": "anthropic",
      "provider_label": "Anthropic",
      "models": ["Claude-3"],
      "requests": 65432,
      "cost_usd": 1234.56,
      "cost_share": 18.0,
      "tokens": 340000
    },
    {
      "provider": "cohere",
      "provider_label": "Cohere",
      "models": ["Command"],
      "requests": 12345,
      "cost_usd": 234.56,
      "cost_share": 4.0,
      "tokens": 170000
    }
  ],
  "alert_configuration": {
    "cost_alert": {
      "enabled": false,
      "threshold_usd": 100.0,
      "period": "daily"
    },
    "error_rate_alert": {
      "enabled": false,
      "threshold_percentage": 5.0,
      "period": "hourly"
    }
  },
  "api_keys_summary": {
    "active_count": 1,
    "total_count": 2,
    "keys": [
      {
        "id": "key_abc123",
        "key_prefix": "ak_live_proj_ghi...",
        "name": "Development Key",
        "status": "active",
        "last_used": "2026-03-08T00:45:00Z",
        "created_at": "2026-01-05T10:00:00Z"
      }
    ]
  }
}
```

**Error Responses:**

| Status | Error               | Description                              |
| ------ | ------------------- | ---------------------------------------- |
| 401    | `unauthorized`      | Missing or invalid JWT token             |
| 403    | `forbidden`         | User doesn't have access to this project |
| 404    | `project_not_found` | Project doesn't exist                    |

---

## **3. GET /api/v1/projects/{projectId}/providers-breakdown**

**Purpose:** Get detailed provider usage breakdown (if needed separately)

**Query Parameters:**

| Parameter    | Type   | Required | Default     | Description               |
| ------------ | ------ | -------- | ----------- | ------------------------- |
| `start_date` | DATE   | No       | month_start | Start date                |
| `end_date`   | DATE   | No       | today       | End date                  |
| `group_by`   | STRING | No       | provider    | Group by: provider, model |

**Request Example:**

```http
GET /api/v1/projects/proj_ghi789/providers-breakdown?start_date=2026-03-01&end_date=2026-03-08
Authorization: Bearer {jwt_token}
```

**Response (200):**

```json
{
  "project_id": "proj_ghi789",
  "period": {
    "start": "2026-03-01",
    "end": "2026-03-08"
  },
  "providers": [
    {
      "provider": "openai",
      "provider_label": "OpenAI",
      "models": [
        {
          "model": "gpt-4",
          "model_label": "GPT-4",
          "requests": 185000,
          "cost_usd": 3890.0,
          "tokens": 680000
        },
        {
          "model": "gpt-3.5-turbo",
          "model_label": "GPT-3.5",
          "requests": 60678,
          "cost_usd": 677.89,
          "tokens": 210000
        }
      ],
      "total_requests": 245678,
      "total_cost_usd": 4567.89,
      "cost_share": 78.0,
      "total_tokens": 890000
    },
    {
      "provider": "anthropic",
      "provider_label": "Anthropic",
      "models": [
        {
          "model": "claude-3-opus",
          "model_label": "Claude-3 Opus",
          "requests": 45432,
          "cost_usd": 989.56,
          "tokens": 240000
        },
        {
          "model": "claude-3-sonnet",
          "model_label": "Claude-3 Sonnet",
          "requests": 20000,
          "cost_usd": 245.0,
          "tokens": 100000
        }
      ],
      "total_requests": 65432,
      "total_cost_usd": 1234.56,
      "cost_share": 18.0,
      "total_tokens": 340000
    }
  ],
  "total_requests": 323455,
  "total_cost_usd": 6037.01,
  "total_tokens": 1400000
}
```

---

## **4. POST /api/v1/organizations/{orgId}/projects**

**Purpose:** Create new project

**Request Body:**

```json
{
  "name": "Development",
  "description": "Development environment for testing new features",
  "environment": "development",
  "generate_api_key": true,
  "alerts": {
    "cost_alert": {
      "enabled": false,
      "threshold_usd": 100.0
    },
    "error_rate_alert": {
      "enabled": false,
      "threshold_percentage": 5.0
    }
  }
}
```

**Request Example:**

```http
POST /api/v1/organizations/org_abc123/projects
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "name": "Development",
  "description": "Development environment for testing new features",
  "environment": "development",
  "generate_api_key": true
}
```

**Response (201):**

```json
{
  "project": {
    "id": "proj_ghi789",
    "name": "Development",
    "description": "Development environment for testing new features",
    "environment": "development",
    "icon": "code",
    "created_at": "2026-03-08T01:00:00Z",
    "updated_at": "2026-03-08T01:00:00Z"
  },
  "api_key": {
    "id": "key_abc123",
    "key": "ak_live_proj_ghi789_xyz789randomstring",
    "key_prefix": "ak_live_proj_ghi...",
    "name": "Default Key",
    "created_at": "2026-03-08T01:00:00Z",
    "warning": "⚠️ Save this key now. You won't be able to see it again."
  }
}
```

---

## **5. PUT /api/v1/projects/{projectId}**

**Purpose:** Update project (from Settings modal)

**Request Body:**

```json
{
  "name": "Development v2",
  "description": "Updated development environment",
  "environment": "development",
  "alerts": {
    "cost_alert": {
      "enabled": true,
      "threshold_usd": 50.0,
      "period": "daily"
    },
    "error_rate_alert": {
      "enabled": true,
      "threshold_percentage": 3.0,
      "period": "hourly"
    }
  }
}
```

**Response (200):**

```json
{
  "id": "proj_ghi789",
  "name": "Development v2",
  "description": "Updated development environment",
  "environment": "development",
  "updated_at": "2026-03-08T01:05:00Z",
  "message": "Project updated successfully"
}
```

---

## **6. DELETE /api/v1/projects/{projectId}**

**Purpose:** Delete project

**Query Parameters:**

| Parameter           | Type    | Required | Description                     |
| ------------------- | ------- | -------- | ------------------------------- |
| `confirm`           | BOOLEAN | Yes      | Must be true                    |
| `confirmation_text` | STRING  | Yes      | Must match project name exactly |

**Request Example:**

```http
DELETE /api/v1/projects/proj_ghi789?confirm=true&confirmation_text=Development
Authorization: Bearer {jwt_token}
```

**Response (200):**

```json
{
  "id": "proj_ghi789",
  "name": "Development",
  "deleted": true,
  "deleted_at": "2026-03-08T01:10:00Z",
  "message": "Project deleted successfully",
  "side_effects": {
    "api_keys_revoked": 2,
    "alerts_deleted": 2,
    "usage_data_archived": true
  }
}
```

---

## **7. GET /api/v1/projects/{projectId}/errors**

**Purpose:** Get error details (when user clicks "2 errors" link)

**Query Parameters:**

| Parameter | Type    | Required | Default | Description                |
| --------- | ------- | -------- | ------- | -------------------------- |
| `date`    | DATE    | No       | today   | Date to filter errors      |
| `limit`   | INTEGER | No       | 50      | Number of errors to return |

**Request Example:**

```http
GET /api/v1/projects/proj_ghi789/errors?date=2026-03-08&limit=50
Authorization: Bearer {jwt_token}
```

**Response (200):**

```json
{
  "project_id": "proj_ghi789",
  "date": "2026-03-08",
  "total_errors": 2,
  "error_rate": 1.3,
  "errors": [
    {
      "id": "err_abc123",
      "timestamp": "2026-03-08T00:45:12Z",
      "provider": "openai",
      "model": "gpt-4",
      "endpoint": "/v1/chat/completions",
      "status_code": 429,
      "error_type": "rate_limit_exceeded",
      "error_message": "Rate limit reached for requests",
      "duration_ms": 0,
      "cost_usd": 0.0
    },
    {
      "id": "err_def456",
      "timestamp": "2026-03-08T00:32:45Z",
      "provider": "anthropic",
      "model": "claude-3",
      "endpoint": "/v1/messages",
      "status_code": 500,
      "error_type": "internal_server_error",
      "error_message": "Internal server error",
      "duration_ms": 0,
      "cost_usd": 0.0
    }
  ]
}
```

---

## **Database Queries**

### **Get Project Details:**

```sql
SELECT
    p.id,
    p.organization_id,
    p.name,
    p.description,
    p.environment,
    p.created_at,
    p.updated_at,
    -- Today stats
    COALESCE(today.total_requests, 0) as today_requests,
    COALESCE(today.successful_requests, 0) as today_successful,
    COALESCE(today.failed_requests, 0) as today_failed,
    COALESCE(today.total_cost_usd, 0) as today_cost,
    COALESCE(today.total_tokens, 0) as today_tokens,
    COALESCE(today.avg_duration_ms, 0) as today_avg_duration,
    -- Month stats
    COALESCE(month.total_requests, 0) as month_requests,
    COALESCE(month.total_cost_usd, 0) as month_cost,
    COALESCE(month.total_tokens, 0) as month_tokens,
    COALESCE(month.avg_duration_ms, 0) as month_avg_duration,
    -- API keys
    COALESCE(keys_active.count, 0) as active_keys,
    COALESCE(keys_total.count, 0) as total_keys
FROM projects p
LEFT JOIN project_daily_stats today ON
    p.id = today.project_id
    AND today.date = $2
LEFT JOIN (
    SELECT
        project_id,
        SUM(total_requests) as total_requests,
        SUM(total_cost_usd) as total_cost_usd,
        SUM(total_tokens) as total_tokens,
        AVG(avg_duration_ms)::INTEGER as avg_duration_ms
    FROM project_daily_stats
    WHERE date >= date_trunc('month', $2)
      AND date <= $2
    GROUP BY project_id
) month ON p.id = month.project_id
LEFT JOIN (
    SELECT project_id, COUNT(*) as count
    FROM api_keys
    WHERE revoked_at IS NULL
    GROUP BY project_id
) keys_active ON p.id = keys_active.project_id
LEFT JOIN (
    SELECT project_id, COUNT(*) as count
    FROM api_keys
    GROUP BY project_id
) keys_total ON p.id = keys_total.project_id
WHERE p.id = $1;
```

### **Get Provider Breakdown:**

```sql
SELECT
    pds.provider,
    SUM(pds.total_requests) as total_requests,
    SUM(pds.total_cost_usd) as total_cost_usd,
    SUM(pds.total_tokens) as total_tokens,
    ROUND((SUM(pds.total_cost_usd) / (
        SELECT SUM(total_cost_usd)
        FROM provider_daily_stats
        WHERE project_id = $1
          AND date BETWEEN $2 AND $3
    ) * 100)::numeric, 1) as cost_share
FROM provider_daily_stats pds
WHERE pds.project_id = $1
  AND pds.date BETWEEN $2 AND $3
GROUP BY pds.provider
ORDER BY total_cost_usd DESC;
```

---

## **Frontend React Components**

### **ProjectDetailModal.tsx:**

```typescript
const ProjectDetailModal = ({ projectId, onClose }) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchProjectDetails();
  }, [projectId]);

  const fetchProjectDetails = async () => {
    const response = await fetch(
      `/api/v1/projects/${projectId}/details`,
      {
        headers: { 'Authorization': `Bearer ${token}` }
      }
    );
    const json = await response.json();
    setData(json);
    setLoading(false);
  };

  if (loading) return <ModalSkeleton />;

  return (
    <Modal onClose={onClose}>
      <ModalHeader>
        <h2>{data.project.name}</h2>
        <EnvironmentBadge>{data.project.environment}</EnvironmentBadge>
        <p>{data.project.description}</p>
      </ModalHeader>

      <StatsGrid>
        <StatCard
          title="Today's Requests"
          value={data.today_stats.requests.total}
          subtitle={`${data.today_stats.requests.failed} errors (${data.today_stats.requests.error_rate}%)`}
        />
        <StatCard
          title="Today's Cost"
          value={`$${data.today_stats.cost_usd}`}
        />
        <StatCard
          title="Tokens Used"
          value={formatNumber(data.today_stats.tokens.total)}
          subtitle={`Month: ${formatNumber(data.month_stats.tokens.total)}`}
        />
        <StatCard
          title="Avg Duration"
          value={`${data.today_stats.performance.avg_duration_ms}ms`}
          subtitle={`Month avg: ${data.month_stats.performance.avg_duration_ms}ms`}
        />
      </StatsGrid>

      <ProvidersSection providers={data.providers_usage} />

      <MonthlyOverview stats={data.month_stats} />

      <AlertConfiguration config={data.alert_configuration} />

      <APIKeysSummary summary={data.api_keys_summary} projectId={projectId} />

      <ModalFooter>
        <Button variant="secondary" onClick={onClose}>Close</Button>
        <Button variant="primary" onClick={() => openSettings(projectId)}>
          ⚙️ Settings
        </Button>
      </ModalFooter>
    </Modal>
  );
};
```

---

## **API Summary Table**

| #   | Endpoint                                 | Purpose              | Trigger               |
| --- | ---------------------------------------- | -------------------- | --------------------- |
| 1   | `GET /projects`                          | List projects        | Projects page load    |
| 2   | `GET /projects/{id}/details`             | Project detail modal | "View Details" button |
| 3   | `GET /projects/{id}/providers-breakdown` | Provider details     | Modal data            |
| 4   | `POST /projects`                         | Create project       | "New Project" button  |
| 5   | `PUT /projects/{id}`                     | Update project       | Settings → Save       |
| 6   | `DELETE /projects/{id}`                  | Delete project       | Delete confirmation   |
| 7   | `GET /projects/{id}/errors`              | Error details        | Click "2 errors" link |

---

## **Rate Limits**

```
GET /projects: 100/hour per user
GET /projects/{id}/details: 200/hour per user
POST /projects: 10/hour per organization
PUT /projects: 50/hour per user
DELETE /projects: 10/hour per organization
```

---

## **Caching Strategy**

```
GET /projects/{id}/details
Cache-Control: max-age=60 (1 minute)
ETag: "abc123xyz"

GET /projects/{id}/providers-breakdown
Cache-Control: max-age=300 (5 minutes)

