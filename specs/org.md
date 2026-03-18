# **Organization API Specs - Crisp Reference**

---

## **1. Create Organization**

```
POST /api/v1/organizations
```

**Request:**
```json
{
  "name": "Acme Corp",
  "slug": "acme-corp"
}
```

**Response (201):**
```json
{
  "id": "org_abc123",
  "name": "Acme Corp",
  "slug": "acme-corp",
  "owner_id": "user_xyz789",
  "plan": "free",
  "created_at": "2026-02-22T10:30:00Z"
}
```

**Errors:** 400 (validation), 401 (unauthorized), 409 (slug taken), 429 (rate limit)

---

## **2. List My Organizations**

```
GET /api/v1/organizations?page=1&limit=20
```

**Response (200):**
```json
{
  "organizations": [
    {
      "id": "org_abc123",
      "name": "Acme Corp",
      "slug": "acme-corp",
      "role": "owner",
      "plan": "pro",
      "member_count": 8,
      "project_count": 5,
      "current_month_cost": 2847.50
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 1
  }
}
```

**Errors:** 401 (unauthorized)

---

## **3. Get Single Organization**

```
GET /api/v1/organizations/{orgId}
```

**Response (200):**
```json
{
  "id": "org_abc123",
  "name": "Acme Corp",
  "slug": "acme-corp",
  "owner_id": "user_xyz789",
  "plan": "pro",
  "created_at": "2026-02-22T10:30:00Z",
  "updated_at": "2026-02-22T10:30:00Z",
  "stats": {
    "member_count": 8,
    "project_count": 5,
    "api_keys_count": 12,
    "current_month_requests": 125430,
    "current_month_cost": 2847.50
  },
  "subscription": {
    "plan": "pro",
    "status": "active",
    "billing_cycle": "monthly",
    "next_billing_date": "2026-03-01"
  }
}
```

**Errors:** 401, 403 (no access), 404 (not found)

---

## **4. Update Organization**

```
PUT /api/v1/organizations/{orgId}
```

**Request:**
```json
{
  "name": "Acme Corporation",
  "slug": "acme-corp"
}
```

**Response (200):**
```json
{
  "id": "org_abc123",
  "name": "Acme Corporation",
  "slug": "acme-corp",
  "owner_id": "user_xyz789",
  "plan": "pro",
  "updated_at": "2026-02-22T15:45:00Z"
}
```

**Errors:** 400 (validation), 401, 403, 404, 409 (slug taken)

---

## **5. Delete Organization**

```
DELETE /api/v1/organizations/{orgId}?confirm=true
```

**Response (200):**
```json
{
  "id": "org_abc123",
  "deleted": true,
  "deleted_at": "2026-02-22T15:50:00Z",
  "message": "Organization deleted. All projects, API keys, and members removed."
}
```

**Errors:** 400 (missing confirm), 401, 403 (only owner), 404

---

## **6. Invite Team Member**

```
POST /api/v1/organizations/{orgId}/members
```

**Request:**
```json
{
  "email": "developer@acme.com",
  "role": "developer"
}
```

**Response (201):**
```json
{
  "invitation_id": "inv_xyz123",
  "email": "developer@acme.com",
  "role": "developer",
  "status": "pending",
  "expires_at": "2026-03-01T10:30:00Z",
  "invited_by": "user_xyz789",
  "created_at": "2026-02-22T10:30:00Z"
}
```

**Errors:** 400, 401, 403, 409 (already member/invited)

---

## **7. List Team Members**

```
GET /api/v1/organizations/{orgId}/members
```

**Response (200):**
```json
{
  "members": [
    {
      "user_id": "user_xyz789",
      "email": "owner@acme.com",
      "name": "John Doe",
      "role": "owner",
      "joined_at": "2026-02-22T10:30:00Z",
      "last_active": "2026-02-22T14:25:00Z"
    }
  ],
  "pending_invitations": [
    {
      "invitation_id": "inv_xyz123",
      "email": "newdev@acme.com",
      "role": "developer",
      "status": "pending",
      "expires_at": "2026-03-01T10:30:00Z"
    }
  ]
}
```

**Errors:** 401, 403, 404

---

## **8. Update Member Role**

```
PUT /api/v1/organizations/{orgId}/members/{userId}
```

**Request:**
```json
{
  "role": "admin"
}
```

**Response (200):**
```json
{
  "user_id": "user_abc456",
  "email": "developer@acme.com",
  "role": "admin",
  "updated_at": "2026-02-22T15:45:00Z"
}
```

**Errors:** 400, 401, 403 (can't change owner), 404

---

## **9. Remove Team Member**

```
DELETE /api/v1/organizations/{orgId}/members/{userId}
```

**Response (200):**
```json
{
  "user_id": "user_abc456",
  "removed": true,
  "removed_at": "2026-02-22T15:50:00Z"
}
```

**Errors:** 401, 403 (can't remove owner/self), 404

---

## **10. Transfer Ownership**

```
POST /api/v1/organizations/{orgId}/transfer-ownership
```

**Request:**
```json
{
  "new_owner_id": "user_abc456"
}
```

**Response (200):**
```json
{
  "id": "org_abc123",
  "previous_owner_id": "user_xyz789",
  "new_owner_id": "user_abc456",
  "transferred_at": "2026-02-22T16:00:00Z"
}
```

**Errors:** 400, 401, 403 (only owner), 404

---

## **Database Tables**

### **organizations**
```sql
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    plan VARCHAR(50) DEFAULT 'free',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_organizations_owner ON organizations(owner_id);
CREATE INDEX idx_organizations_slug ON organizations(slug);
```

### **organization_members**
```sql
CREATE TABLE organization_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,
    joined_at TIMESTAMP DEFAULT NOW(),
    last_active_at TIMESTAMP,
    UNIQUE(organization_id, user_id)
);

CREATE INDEX idx_org_members_org ON organization_members(organization_id);
CREATE INDEX idx_org_members_user ON organization_members(user_id);
```

### **organization_invitations**
```sql
CREATE TABLE organization_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    invited_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'pending',
    expires_at TIMESTAMP NOT NULL,
    accepted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_org_invitations_org ON organization_invitations(organization_id);
CREATE INDEX idx_org_invitations_token ON organization_invitations(token);
```

---

## **Roles**

| Role | Permissions |
|------|------------|
| **owner** | Full control, delete org, transfer ownership |
| **admin** | Manage projects, API keys, invite members |
| **developer** | Create API keys, view analytics |
| **viewer** | Read-only access |

---

## **Validation Rules**

**Organization name:**
- 2-255 characters
- Any characters allowed

**Slug:**
- 3-100 characters
- Only: lowercase letters, numbers, hyphens
- Must be globally unique
- Format: `^[a-z0-9-]+$`

**Role:**
- Must be one of: `owner`, `admin`, `developer`, `viewer`
- Cannot invite as `owner` (only one owner per org)

**Email:**
- Valid email format
- Cannot invite existing member
- Cannot have duplicate pending invitations

---

## **Common Errors**

| Code | Error | When |
|------|-------|------|
| 400 | `validation_error` | Invalid input |
| 401 | `unauthorized` | Missing/invalid token |
| 403 | `forbidden` | Insufficient permissions |
| 404 | `not_found` | Resource doesn't exist |
| 409 | `conflict` | Duplicate slug/email |
| 429 | `rate_limit_exceeded` | Too many requests |
| 500 | `internal_server_error` | Server error |

---

## **Rate Limits**

- Create org: 10/hour per user
- Invite member: 20/hour per org
- List/Get: 1000/hour per user
- Update/Delete: 100/hour per user

---

## **Quick Reference**

**Headers (all endpoints):**
```
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

**Pagination (list endpoints):**
```
?page=1&limit=20
```

**Role hierarchy:**
```
owner > admin > developer > viewer
```

**Invitation flow:**
```
1. POST /organizations/{id}/members → RabbitMQ
2. Worker sends email
3. User clicks link → GET /invite/{token}
4. Accept → Creates organization_member
5. Status updated to 'accepted'
```

---

**Save this file as `ORGANIZATION_API.md` for quick reference!**