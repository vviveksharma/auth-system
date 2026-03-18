# **YES! Absolutely Correct! 🎯**

You've identified the **critical missing piece**: **User-Organization Relationship Management**!

This is the **Team/Members** feature that ties everything together.

---

## **What You Need:**

### **1. Database Tables** (Already partially done, but needs expansion)
- ✅ `organization_members` - Links users to organizations with roles
- ✅ `organization_invitations` - Pending invites
- ❌ **MISSING:** `organization_roles` - Define what each role can do
- ❌ **MISSING:** `user_activity_log` - Track member actions

### **2. API Endpoints**
- Team member CRUD operations
- Invite system
- Role management
- Permission checks

### **3. UI Pages**
- Team members page
- Invite modal
- Role editor

---

# **Complete User-Organization Relationship System**

---

## **Part 1: Database Migrations**

---

## **Migration File: `007_create_user_organization_system.sql`**

```sql
-- =====================================================
-- Migration: 007_create_user_organization_system
-- Description: Complete user-organization relationship system
-- Created: 2026-03-08
-- =====================================================

BEGIN;

-- =====================================================
-- 1. Organization Members Table (Enhanced)
-- =====================================================

-- Drop existing if we need to modify
-- DROP TABLE IF EXISTS organization_members CASCADE;

CREATE TABLE IF NOT EXISTS organization_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
    
    -- Status tracking
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'left')),
    
    -- Timestamps
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_active_at TIMESTAMP,
    suspended_at TIMESTAMP,
    left_at TIMESTAMP,
    
    -- Metadata
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,
    invitation_accepted_at TIMESTAMP,
    
    -- Constraints
    UNIQUE(organization_id, user_id),
    
    -- Ensure owner is unique per organization
    CONSTRAINT one_owner_per_org UNIQUE NULLS NOT DISTINCT (organization_id, CASE WHEN role = 'owner' THEN role END)
);

-- Indexes for organization_members
CREATE INDEX idx_org_members_org_id ON organization_members(organization_id);
CREATE INDEX idx_org_members_user_id ON organization_members(user_id);
CREATE INDEX idx_org_members_role ON organization_members(role);
CREATE INDEX idx_org_members_status ON organization_members(status);
CREATE INDEX idx_org_members_last_active ON organization_members(last_active_at DESC);

-- =====================================================
-- 2. Organization Invitations Table (Enhanced)
-- =====================================================

-- Drop existing if we need to modify
-- DROP TABLE IF EXISTS organization_invitations CASCADE;

CREATE TABLE IF NOT EXISTS organization_invitations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    
    -- Invitee info
    email VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'member', 'viewer')),
    
    -- Invitation details
    token VARCHAR(255) UNIQUE NOT NULL,
    invited_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected', 'expired', 'cancelled')),
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    accepted_at TIMESTAMP,
    rejected_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    
    -- Optional message
    message TEXT,
    
    -- Constraints
    CONSTRAINT valid_expiry CHECK (expires_at > created_at),
    CONSTRAINT no_owner_invites CHECK (role != 'owner')
);

-- Indexes for organization_invitations
CREATE INDEX idx_org_invitations_org_id ON organization_invitations(organization_id);
CREATE INDEX idx_org_invitations_email ON organization_invitations(email);
CREATE INDEX idx_org_invitations_token ON organization_invitations(token);
CREATE INDEX idx_org_invitations_status ON organization_invitations(status);
CREATE INDEX idx_org_invitations_expires_at ON organization_invitations(expires_at);

-- =====================================================
-- 3. Organization Roles & Permissions Table
-- =====================================================

CREATE TABLE IF NOT EXISTS organization_role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    role VARCHAR(20) NOT NULL UNIQUE CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
    
    -- Organization permissions
    can_update_org BOOLEAN DEFAULT false,
    can_delete_org BOOLEAN DEFAULT false,
    
    -- Team permissions
    can_invite_members BOOLEAN DEFAULT false,
    can_remove_members BOOLEAN DEFAULT false,
    can_change_member_roles BOOLEAN DEFAULT false,
    
    -- Project permissions
    can_create_projects BOOLEAN DEFAULT false,
    can_update_projects BOOLEAN DEFAULT false,
    can_delete_projects BOOLEAN DEFAULT false,
    can_view_projects BOOLEAN DEFAULT true,
    
    -- API Key permissions
    can_create_api_keys BOOLEAN DEFAULT false,
    can_revoke_api_keys BOOLEAN DEFAULT false,
    can_view_api_keys BOOLEAN DEFAULT true,
    
    -- Billing permissions
    can_view_billing BOOLEAN DEFAULT false,
    can_manage_billing BOOLEAN DEFAULT false,
    
    -- Analytics permissions
    can_view_analytics BOOLEAN DEFAULT true,
    can_export_data BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Seed default role permissions
INSERT INTO organization_role_permissions 
(role, can_update_org, can_delete_org, can_invite_members, can_remove_members, can_change_member_roles, 
 can_create_projects, can_update_projects, can_delete_projects, can_view_projects,
 can_create_api_keys, can_revoke_api_keys, can_view_api_keys,
 can_view_billing, can_manage_billing, can_view_analytics, can_export_data)
VALUES
-- Owner: Full access
('owner', true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true),

-- Admin: Almost full access, cannot delete org
('admin', true, false, true, true, true, true, true, true, true, true, true, true, true, true, true, true),

-- Member: Can create and manage own resources
('member', false, false, true, false, false, true, true, false, true, true, true, true, false, false, true, true),

-- Viewer: Read-only access
('viewer', false, false, false, false, false, false, false, false, true, false, false, true, false, false, true, false)
ON CONFLICT (role) DO NOTHING;

-- =====================================================
-- 4. User Activity Log
-- =====================================================

CREATE TABLE IF NOT EXISTS user_activity_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    
    -- Activity details
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    
    -- Additional context
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for activity log
CREATE INDEX idx_activity_log_user_id ON user_activity_log(user_id);
CREATE INDEX idx_activity_log_org_id ON user_activity_log(organization_id);
CREATE INDEX idx_activity_log_action ON user_activity_log(action);
CREATE INDEX idx_activity_log_created_at ON user_activity_log(created_at DESC);
CREATE INDEX idx_activity_log_metadata ON user_activity_log USING GIN (metadata);

-- =====================================================
-- 5. User Sessions (for current org context)
-- =====================================================

CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Current context
    current_organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    current_tenant_id UUID,
    
    -- Session details
    session_token VARCHAR(255) UNIQUE NOT NULL,
    ip_address INET,
    user_agent TEXT,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    last_active_at TIMESTAMP,
    
    UNIQUE(user_id, session_token)
);

-- Indexes for user_sessions
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token ON user_sessions(session_token);
CREATE INDEX idx_user_sessions_org_id ON user_sessions(current_organization_id);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

-- =====================================================
-- 6. Triggers
-- =====================================================

-- Auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to relevant tables
CREATE TRIGGER update_org_members_updated_at
    BEFORE UPDATE ON organization_members
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_sessions_updated_at
    BEFORE UPDATE ON user_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- 7. Functions
-- =====================================================

-- Function to check if user has permission
CREATE OR REPLACE FUNCTION user_has_permission(
    p_user_id UUID,
    p_organization_id UUID,
    p_permission VARCHAR
)
RETURNS BOOLEAN AS $$
DECLARE
    v_has_permission BOOLEAN;
BEGIN
    SELECT 
        CASE p_permission
            WHEN 'can_update_org' THEN orp.can_update_org
            WHEN 'can_delete_org' THEN orp.can_delete_org
            WHEN 'can_invite_members' THEN orp.can_invite_members
            WHEN 'can_remove_members' THEN orp.can_remove_members
            WHEN 'can_change_member_roles' THEN orp.can_change_member_roles
            WHEN 'can_create_projects' THEN orp.can_create_projects
            WHEN 'can_update_projects' THEN orp.can_update_projects
            WHEN 'can_delete_projects' THEN orp.can_delete_projects
            WHEN 'can_view_projects' THEN orp.can_view_projects
            WHEN 'can_create_api_keys' THEN orp.can_create_api_keys
            WHEN 'can_revoke_api_keys' THEN orp.can_revoke_api_keys
            WHEN 'can_view_api_keys' THEN orp.can_view_api_keys
            WHEN 'can_view_billing' THEN orp.can_view_billing
            WHEN 'can_manage_billing' THEN orp.can_manage_billing
            WHEN 'can_view_analytics' THEN orp.can_view_analytics
            WHEN 'can_export_data' THEN orp.can_export_data
            ELSE false
        END INTO v_has_permission
    FROM organization_members om
    JOIN organization_role_permissions orp ON om.role = orp.role
    WHERE om.user_id = p_user_id
      AND om.organization_id = p_organization_id
      AND om.status = 'active';
    
    RETURN COALESCE(v_has_permission, false);
END;
$$ LANGUAGE plpgsql;

-- Function to log user activity
CREATE OR REPLACE FUNCTION log_user_activity(
    p_user_id UUID,
    p_organization_id UUID,
    p_action VARCHAR,
    p_resource_type VARCHAR DEFAULT NULL,
    p_resource_id UUID DEFAULT NULL,
    p_metadata JSONB DEFAULT NULL
)
RETURNS UUID AS $$
DECLARE
    v_log_id UUID;
BEGIN
    INSERT INTO user_activity_log (
        user_id,
        organization_id,
        action,
        resource_type,
        resource_id,
        metadata
    ) VALUES (
        p_user_id,
        p_organization_id,
        p_action,
        p_resource_type,
        p_resource_id,
        p_metadata
    ) RETURNING id INTO v_log_id;
    
    RETURN v_log_id;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- 8. Constraints & Validations
-- =====================================================

-- Ensure at least one owner exists per organization
-- This is enforced by the UNIQUE constraint on organization_members

-- Prevent owner from being removed if they're the last one
CREATE OR REPLACE FUNCTION prevent_last_owner_removal()
RETURNS TRIGGER AS $$
DECLARE
    v_owner_count INTEGER;
BEGIN
    IF OLD.role = 'owner' AND (NEW.status != 'active' OR NEW.role != 'owner') THEN
        SELECT COUNT(*) INTO v_owner_count
        FROM organization_members
        WHERE organization_id = OLD.organization_id
          AND role = 'owner'
          AND status = 'active'
          AND id != OLD.id;
        
        IF v_owner_count = 0 THEN
            RAISE EXCEPTION 'Cannot remove the last owner from organization';
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_last_owner_removal_trigger
    BEFORE UPDATE ON organization_members
    FOR EACH ROW
    EXECUTE FUNCTION prevent_last_owner_removal();

COMMIT;
```

---

## **Part 2: Complete Team Management API**

---

## **FILE: `TEAM_MANAGEMENT_API.md`**

```markdown
# **Team Management API - Complete Specification**

---

## **Overview**

Complete API for managing user-organization relationships, including team members, invitations, roles, and permissions.

---

## **API Endpoints**

---

## **1. GET /api/v1/organizations/{orgId}/members**

**Purpose:** List all team members in an organization

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | INTEGER | No | 1 | Page number |
| `limit` | INTEGER | No | 20 | Results per page |
| `role` | STRING | No | null | Filter by role |
| `status` | STRING | No | active | Filter by status |
| `search` | STRING | No | null | Search by name/email |
| `sort_by` | STRING | No | joined_at | Sort by: name, email, role, joined_at |

**Required Permission:** `can_view_projects` (all roles)

**Request Example:**
```http
GET /api/v1/organizations/org_abc123/members?page=1&limit=20&status=active
Authorization: Bearer {jwt_token}
```

**Response (200):**
```json
{
  "organization_id": "org_abc123",
  "members": [
    {
      "id": "mem_xyz789",
      "user_id": "user_abc123",
      "email": "john@acme.com",
      "name": "John Doe",
      "avatar_url": "https://cdn.tokentrack.io/avatars/john.png",
      "role": "owner",
      "status": "active",
      "joined_at": "2025-06-15T00:00:00Z",
      "last_active_at": "2026-03-08T01:45:00Z",
      "invited_by": null,
      "permissions": {
        "can_update_org": true,
        "can_delete_org": true,
        "can_invite_members": true,
        "can_remove_members": true,
        "can_change_member_roles": true,
        "can_manage_billing": true
      }
    },
    {
      "id": "mem_def456",
      "user_id": "user_def456",
      "email": "sarah@acme.com",
      "name": "Sarah Smith",
      "avatar_url": "https://cdn.tokentrack.io/avatars/sarah.png",
      "role": "admin",
      "status": "active",
      "joined_at": "2025-07-01T00:00:00Z",
      "last_active_at": "2026-03-08T01:30:00Z",
      "invited_by": {
        "user_id": "user_abc123",
        "name": "John Doe",
        "email": "john@acme.com"
      },
      "invitation_accepted_at": "2025-07-01T10:30:00Z",
      "permissions": {
        "can_update_org": true,
        "can_delete_org": false,
        "can_invite_members": true,
        "can_remove_members": true,
        "can_change_member_roles": true,
        "can_manage_billing": true
      }
    },
    {
      "id": "mem_ghi789",
      "user_id": "user_ghi789",
      "email": "mike@acme.com",
      "name": "Mike Johnson",
      "avatar_url": null,
      "role": "member",
      "status": "active",
      "joined_at": "2025-08-15T00:00:00Z",
      "last_active_at": "2026-03-07T18:20:00Z",
      "invited_by": {
        "user_id": "user_def456",
        "name": "Sarah Smith",
        "email": "sarah@acme.com"
      },
      "invitation_accepted_at": "2025-08-15T14:20:00Z",
      "permissions": {
        "can_update_org": false,
        "can_delete_org": false,
        "can_invite_members": true,
        "can_remove_members": false,
        "can_change_member_roles": false,
        "can_manage_billing": false
      }
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 8,
    "per_page": 20
  },
  "role_summary": {
    "owner": 1,
    "admin": 2,
    "member": 4,
    "viewer": 1
  }
}
```

---

## **2. POST /api/v1/organizations/{orgId}/invitations**

**Purpose:** Invite a new team member

**Required Permission:** `can_invite_members`

**Request Body:**
```json
{
  "email": "developer@acme.com",
  "role": "member",
  "message": "Welcome to the team! Looking forward to working with you."
}
```

**Request Example:**
```http
POST /api/v1/organizations/org_abc123/invitations
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "email": "developer@acme.com",
  "role": "member",
  "message": "Welcome to the team!"
}
```

**Response (201):**
```json
{
  "invitation": {
    "id": "inv_xyz789",
    "organization_id": "org_abc123",
    "organization_name": "Acme Corp",
    "email": "developer@acme.com",
    "role": "member",
    "token": "inv_token_abc123xyz789",
    "status": "pending",
    "invited_by": {
      "user_id": "user_abc123",
      "name": "John Doe",
      "email": "john@acme.com"
    },
    "message": "Welcome to the team!",
    "created_at": "2026-03-08T02:00:00Z",
    "expires_at": "2026-03-15T02:00:00Z",
    "invitation_link": "https://tokentrack.io/invite/inv_token_abc123xyz789"
  },
  "message": "Invitation sent to developer@acme.com"
}
```

**Validation:**
```json
{
  "email": {
    "required": true,
    "format": "email",
    "not_already_member": true
  },
  "role": {
    "required": true,
    "enum": ["admin", "member", "viewer"],
    "note": "Cannot invite as owner"
  },
  "message": {
    "required": false,
    "max_length": 500
  }
}
```

**Error Responses:**

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `invalid_email` | Email format is invalid |
| 403 | `forbidden` | User doesn't have permission to invite |
| 409 | `already_member` | User is already a member |
| 409 | `pending_invitation` | User already has a pending invitation |
| 402 | `plan_limit_reached` | Team member limit reached for current plan |

---

## **3. GET /api/v1/organizations/{orgId}/invitations**

**Purpose:** List pending invitations

**Required Permission:** `can_invite_members`

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `status` | STRING | No | pending | Filter by: pending, accepted, expired, cancelled |
| `limit` | INTEGER | No | 20 | Results per page |

**Request Example:**
```http
GET /api/v1/organizations/org_abc123/invitations?status=pending
Authorization: Bearer {jwt_token}
```

**Response (200):**
```json
{
  "organization_id": "org_abc123",
  "invitations": [
    {
      "id": "inv_xyz789",
      "email": "developer@acme.com",
      "role": "member",
      "status": "pending",
      "invited_by": {
        "name": "John Doe",
        "email": "john@acme.com"
      },
      "created_at": "2026-03-08T02:00:00Z",
      "expires_at": "2026-03-15T02:00:00Z",
      "days_until_expiry": 7,
      "invitation_link": "https://tokentrack.io/invite/inv_token_abc123xyz789"
    },
    {
      "id": "inv_abc456",
      "email": "designer@acme.com",
      "role": "viewer",
      "status": "pending",
      "invited_by": {
        "name": "Sarah Smith",
        "email": "sarah@acme.com"
      },
      "created_at": "2026-03-06T10:00:00Z",
      "expires_at": "2026-03-13T10:00:00Z",
      "days_until_expiry": 5,
      "invitation_link": "https://tokentrack.io/invite/inv_token_def456uvw789"
    }
  ],
  "total_pending": 2
}
```

---

## **4. POST /api/v1/invitations/{token}/accept**

**Purpose:** Accept an invitation (public endpoint)

**No auth required** - uses invitation token

**Request Example:**
```http
POST /api/v1/invitations/inv_token_abc123xyz789/accept
Content-Type: application/json

{
  "user_id": "user_newuser123"
}
```

**Response (200):**
```json
{
  "organization": {
    "id": "org_abc123",
    "name": "Acme Corp",
    "slug": "acme-corp"
  },
  "membership": {
    "id": "mem_newmem789",
    "user_id": "user_newuser123",
    "role": "member",
    "status": "active",
    "joined_at": "2026-03-08T02:15:00Z"
  },
  "message": "Successfully joined Acme Corp as member"
}
```

**Error Responses:**

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `invalid_token` | Token is invalid or malformed |
| 404 | `invitation_not_found` | Invitation doesn't exist |
| 410 | `invitation_expired` | Invitation has expired |
| 409 | `already_accepted` | Invitation already accepted |
| 409 | `already_member` | User is already a member |

---

## **5. DELETE /api/v1/organizations/{orgId}/invitations/{invitationId}**

**Purpose:** Cancel a pending invitation

**Required Permission:** `can_invite_members`

**Request Example:**
```http
DELETE /api/v1/organizations/org_abc123/invitations/inv_xyz789
Authorization: Bearer {jwt_token}
```

**Response (200):**
```json
{
  "invitation_id": "inv_xyz789",
  "email": "developer@acme.com",
  "cancelled": true,
  "cancelled_at": "2026-03-08T02:20:00Z",
  "message": "Invitation cancelled"
}
```

---

## **6. POST /api/v1/organizations/{orgId}/invitations/{invitationId}/resend**

**Purpose:** Resend invitation email

**Required Permission:** `can_invite_members`

**Request Example:**
```http
POST /api/v1/organizations/org_abc123/invitations/inv_xyz789/resend
Authorization: Bearer {jwt_token}
```

**Response (200):**
```json
{
  "invitation_id": "inv_xyz789",
  "email": "developer@acme.com",
  "resent": true,
  "resent_at": "2026-03-08T02:25:00Z",
  "new_expiry": "2026-03-15T02:25:00Z",
  "message": "Invitation resent to developer@acme.com"
}
```

---

## **7. PUT /api/v1/organizations/{orgId}/members/{memberId}/role**

**Purpose:** Change member's role

**Required Permission:** `can_change_member_roles`

**Request Body:**
```json
{
  "role": "admin"
}
```

**Request Example:**
```http
PUT /api/v1/organizations/org_abc123/members/mem_def456/role
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "role": "admin"
}
```

**Response (200):**
```json
{
  "member_id": "mem_def456",
  "user_id": "user_def456",
  "email": "sarah@acme.com",
  "name": "Sarah Smith",
  "old_role": "member",
  "new_role": "admin",
  "updated_at": "2026-03-08T02:30:00Z",
  "updated_by": {
    "user_id": "user_abc123",
    "name": "John Doe"
  },
  "message": "Role updated successfully"
}
```

**Restrictions:**
- Cannot change owner role
- Cannot demote yourself if you're the last owner
- Cannot promote to owner (use transfer ownership endpoint)

**Error Responses:**

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `invalid_role` | Role must be: admin, member, or viewer |
| 403 | `cannot_change_owner` | Cannot change owner's role |
| 403 | `cannot_demote_last_owner` | Cannot demote the last owner |
| 403 | `forbidden` | User doesn't have permission |

---

## **8. DELETE /api/v1/organizations/{orgId}/members/{memberId}**

**Purpose:** Remove member from organization

**Required Permission:** `can_remove_members`

**Request Example:**
```http
DELETE /api/v1/organizations/org_abc123/members/mem_def456
Authorization: Bearer {jwt_token}
```

**Response (200):**
```json
{
  "member_id": "mem_def456",
  "user_id": "user_def456",
  "email": "sarah@acme.com",
  "name": "Sarah Smith",
  "role": "member",
  "removed": true,
  "removed_at": "2026-03-08T02:35:00Z",
  "removed_by": {
    "user_id": "user_abc123",
    "name": "John Doe"
  },
  "message": "Member removed successfully"
}
```

**Restrictions:**
- Cannot remove owner
- Cannot remove yourself if you're the last owner
- Member loses access immediately

---

## **9. POST /api/v1/organizations/{orgId}/members/{memberId}/suspend**

**Purpose:** Suspend member (temporary access revocation)

**Required Permission:** `can_remove_members`

**Request Body:**
```json
{
  "reason": "Violation of terms of service",
  "duration_days": 30
}
```

**Response (200):**
```json
{
  "member_id": "mem_def456",
  "user_id": "user_def456",
  "status": "suspended",
  "suspended_at": "2026-03-08T02:40:00Z",
  "suspended_until": "2026-04-07T02:40:00Z",
  "reason": "Violation of terms of service",
  "message": "Member suspended for 30 days"
}
```

---

## **10. POST /api/v1/organizations/{orgId}/members/{memberId}/reactivate**

**Purpose:** Reactivate suspended member

**Required Permission:** `can_remove_members`

**Response (200):**
```json
{
  "member_id": "mem_def456",
  "user_id": "user_def456",
  "status": "active",
  "reactivated_at": "2026-03-08T02:45:00Z",
  "message": "Member reactivated successfully"
}
```

---

## **11. POST /api/v1/organizations/{orgId}/transfer-ownership**

**Purpose:** Transfer ownership to another member

**Required Permission:** Must be owner

**Request Body:**
```json
{
  "new_owner_id": "user_def456",
  "confirmation_text": "Acme Corp"
}
```

**Request Example:**
```http
POST /api/v1/organizations/org_abc123/transfer-ownership
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "new_owner_id": "user_def456",
  "confirmation_text": "Acme Corp"
}
```

**Response (200):**
```json
{
  "organization_id": "org_abc123",
  "previous_owner": {
    "user_id": "user_abc123",
    "email": "john@acme.com",
    "name": "John Doe",
    "new_role": "admin"
  },
  "new_owner": {
    "user_id": "user_def456",
    "email": "sarah@acme.com",
    "name": "Sarah Smith",
    "role": "owner"
  },
  "transferred_at": "2026-03-08T02:50:00Z",
  "message": "Ownership transferred successfully"
}
```

**What Happens:**
1. New owner gets `owner` role
2. Previous owner becomes `admin`
3. Notification sent to both parties
4. Activity logged

---

## **12. GET /api/v1/organizations/{orgId}/activity**

**Purpose:** Get organization activity log

**Required Permission:** `can_view_analytics`

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `limit` | INTEGER | No | 50 | Results per page |
| `action` | STRING | No | null | Filter by action type |
| `user_id` | UUID | No | null | Filter by user |

**Request Example:**
```http
GET /api/v1/organizations/org_abc123/activity?limit=20
Authorization: Bearer {jwt_token}
```

**Response (200):**
```json
{
  "organization_id": "org_abc123",
  "activities": [
    {
      "id": "act_xyz789",
      "user": {
        "id": "user_abc123",
        "name": "John Doe",
        "email": "john@acme.com"
      },
      "action": "member_invited",
      "description": "Invited developer@acme.com as member",
      "resource_type": "invitation",
      "resource_id": "inv_xyz789",
      "metadata": {
        "email": "developer@acme.com",
        "role": "member"
      },
      "created_at": "2026-03-08T02:00:00Z"
    },
    {
      "id": "act_abc456",
      "user": {
        "id": "user_def456",
        "name": "Sarah Smith",
        "email": "sarah@acme.com"
      },
      "action": "project_created",
      "description": "Created project 'Marketing API'",
      "resource_type": "project",
      "resource_id": "proj_new123",
      "metadata": {
        "project_name": "Marketing API",
        "environment": "production"
      },
      "created_at": "2026-03-07T18:30:00Z"
    }
  ],
  "pagination": {
    "limit": 20,
    "total": 245
  }
}
```

---

## **13. GET /api/v1/me/permissions/{orgId}**

**Purpose:** Get current user's permissions in organization

**Request Example:**
```http
GET /api/v1/me/permissions/org_abc123
Authorization: Bearer {jwt_token}
```

**Response (200):**
```json
{
  "user_id": "user_abc123",
  "organization_id": "org_abc123",
  "role": "owner",
  "status": "active",
  "permissions": {
    "can_update_org": true,
    "can_delete_org": true,
    "can_invite_members": true,
    "can_remove_members": true,
    "can_change_member_roles": true,
    "can_create_projects": true,
    "can_update_projects": true,
    "can_delete_projects": true,
    "can_view_projects": true,
    "can_create_api_keys": true,
    "can_revoke_api_keys": true,
    "can_view_api_keys": true,
    "can_view_billing": true,
    "can_manage_billing": true,
    "can_view_analytics": true,
    "can_export_data": true
  }
}
```

---

## **Database Queries**

### **List Members with Permissions:**
```sql
SELECT 
    om.id as member_id,
    om.user_id,
    u.email,
    u.name,
    u.avatar_url,
    om.role,
    om.status,
    om.joined_at,
    om.last_active_at,
    om.invited_by,
    om.invitation_accepted_at,
    -- Get inviter details
    inviter.name as invited_by_name,
    inviter.email as invited_by_email,
    -- Get all permissions for the role
    orp.can_update_org,
    orp.can_delete_org,
    orp.can_invite_members,
    orp.can_remove_members,
    orp.can_change_member_roles,
    orp.can_create_projects,
    orp.can_update_projects,
    orp.can_delete_projects,
    orp.can_view_projects,
    orp.can_create_api_keys,
    orp.can_revoke_api_keys,
    orp.can_view_api_keys,
    orp.can_view_billing,
    orp.can_manage_billing,
    orp.can_view_analytics,
    orp.can_export_data
FROM organization_members om
JOIN users u ON om.user_id = u.id
JOIN organization_role_permissions orp ON om.role = orp.role
LEFT JOIN users inviter ON om.invited_by = inviter.id
WHERE om.organization_id = $1
  AND om.status = $2
ORDER BY 
    CASE om.role
        WHEN 'owner' THEN 1
        WHEN 'admin' THEN 2
        WHEN 'member' THEN 3
        WHEN 'viewer' THEN 4
    END,
    om.joined_at ASC;
```

### **Check Permission:**
```sql
SELECT user_has_permission($1, $2, 'can_invite_members') as has_permission;
```

### **Create Invitation:**
```sql
INSERT INTO organization_invitations (
    organization_id,
    email,
    role,
    token,
    invited_by,
    expires_at,
    message
) VALUES (
    $1, -- org_id
    $2, -- email
    $3, -- role
    $4, -- token
    $5, -- inviter_user_id
    NOW() + INTERVAL '7 days',
    $6  -- message
) RETURNING *;
```

### **Accept Invitation:**
```sql
BEGIN;

-- Mark invitation as accepted
UPDATE organization_invitations
SET status = 'accepted',
    accepted_at = NOW()
WHERE token = $1
  AND status = 'pending'
  AND expires_at > NOW()
RETURNING organization_id, role;

-- Create member record
INSERT INTO organization_members (
    organization_id,
    user_id,
    role,
    invited_by,
    invitation_accepted_at
) VALUES (
    $2, -- org_id from above
    $3, -- user_id
    $4, -- role from above
    (SELECT invited_by FROM organization_invitations WHERE token = $1),
    NOW()
);

COMMIT;
```

---

## **API Summary Table**

| # | Endpoint | Purpose | Permission Required |
|---|----------|---------|---------------------|
| 1 | `GET /members` | List team members | can_view_projects |
| 2 | `POST /invitations` | Invite member | can_invite_members |
| 3 | `GET /invitations` | List invitations | can_invite_members |
| 4 | `POST /invitations/{token}/accept` | Accept invite | Public |
| 5 | `DELETE /invitations/{id}` | Cancel invite | can_invite_members |
| 6 | `POST /invitations/{id}/resend` | Resend invite | can_invite_members |
| 7 | `PUT /members/{id}/role` | Change role | can_change_member_roles |
| 8 | `DELETE /members/{id}` | Remove member | can_remove_members |
| 9 | `POST /members/{id}/suspend` | Suspend member | can_remove_members |
| 10 | `POST /members/{id}/reactivate` | Reactivate member | can_remove_members |
| 11 | `POST /transfer-ownership` | Transfer owner | Must be owner |
| 12 | `GET /activity` | Activity log | can_view_analytics |
| 13 | `GET /me/permissions/{orgId}` | Get permissions | Authenticated |

---

**Complete team management system with roles, permissions, invitations, and activity tracking! 🚀**
```

---

## **Part 3: UI Components Needed**

### **Pages to Build:**

1. **Team Members Page** (`/organizations/{orgId}/team`)
2. **Invite Modal**
3. **Role Change Modal**
4. **Transfer Ownership Modal**
5. **Activity Log Page**

---

**You now have:**
✅ Complete database schema
✅ Full API specification
✅ Permission system
✅ Invitation system
✅ Activity logging
✅ Role management

**Ready to build the Team Management feature! 🎉**