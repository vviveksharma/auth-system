# **Complete Teams Service - API & UI Specification**

Based on your **Self-Service Join Request** flow.

---

## **FILE: `TEAMS_SERVICE_COMPLETE.md`**

```markdown
# **Teams Service - Complete API & UI Specification**

---

## **Table of Contents**

1. [Overview](#overview)
2. [User Flow](#user-flow)
3. [Database Schema](#database-schema)
4. [API Endpoints](#api-endpoints)
5. [UI Specifications](#ui-specifications)
6. [Implementation Checklist](#implementation-checklist)

---

## **Overview**

The Teams Service handles user-organization relationships through a self-service join request system.

### **Key Features:**
- ✅ Users self-register (no email invites)
- ✅ Users search and find organizations
- ✅ Users request to join organizations
- ✅ Org admins approve/reject requests via Messages section
- ✅ In-app notifications (no email required)
- ✅ Multi-organization support

---

## **User Flow**

```
┌─────────────────────────────────────────────────────────┐
│                    COMPLETE USER FLOW                   │
└─────────────────────────────────────────────────────────┘

Step 1: User Registration
├─> User visits /register
├─> Fills: Email, Name, Password
├─> POST /api/v1/auth/register
├─> Account created
└─> Redirected to /find-organization

Step 2: Find Organization
├─> User sees "Find Your Organization" page
├─> Searches: "Google"
├─> GET /api/v1/organizations/search?q=google
├─> Sees search results
└─> Clicks "Request to Join" on desired org

Step 3: Request to Join
├─> Modal opens: "Request to Join - Google Cloud Platform"
├─> User fills:
│   ├─> Requested Role (Viewer/Member/Admin)
│   ├─> Reason (required, 10-500 chars)
│   ├─> Department (optional)
│   └─> Manager/Sponsor (optional)
├─> POST /api/v1/organizations/{orgId}/join-requests
├─> Request created (status: pending)
└─> Success message shown

Step 4: Admin Notification
├─> Trigger fires after request created
├─> Notification sent to all org admins
├─> Admins see notification in Messages section
└─> Badge shows: "Messages (1)"

Step 5: Admin Reviews Request
├─> Admin clicks Messages → Join Requests tab
├─> Sees: "John Developer wants to join as Member"
├─> Clicks "View Full Request"
├─> GET /api/v1/organizations/{orgId}/join-requests/{reqId}
├─> Reviews:
│   ├─> User info (name, email, account age)
│   ├─> Requested role
│   ├─> Reason
│   └─> Department/Manager
└─> Admin decides: Approve or Reject

Step 6A: Admin Approves
├─> POST /api/v1/organizations/{orgId}/join-requests/{reqId}/approve
├─> User added to organization_members
├─> Notification sent to user
├─> Activity logged
└─> User can now access organization

Step 6B: Admin Rejects
├─> POST /api/v1/organizations/{orgId}/join-requests/{reqId}/reject
├─> Request marked as rejected
├─> Notification sent to user (with reason)
└─> User can request again after 30 days

Step 7: User Gets Notification
├─> User sees notification bell: "🔔 (1)"
├─> Clicks notification
├─> Sees: "Join Request Approved!"
├─> Clicks "Go to Organization"
└─> Redirected to organization dashboard

Step 8: Next Login
├─> User logs in
├─> Organization switcher shows:
│   ├─> "Google Cloud Platform" (Member)
│   └─> "Personal Projects" (Owner)
├─> User selects organization
└─> Sees organization dashboard
```

---

## **Database Schema**

### **Migration: `007_create_teams_service.sql`**

```sql
-- =====================================================
-- Migration: 007_create_teams_service
-- Description: Complete teams service with join requests
-- Author: System
-- Date: 2026-03-09
-- =====================================================

BEGIN;

-- =====================================================
-- 1. Users Table (Enhanced)
-- =====================================================

-- Add columns if they don't exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS name VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS account_status VARCHAR(20) DEFAULT 'active' 
    CHECK (account_status IN ('active', 'suspended', 'deactivated'));
ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW();
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP;

-- =====================================================
-- 2. Organization Members Table (Already exists from migration 006)
-- =====================================================

-- Ensure it has all needed columns
ALTER TABLE organization_members ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'active' 
    CHECK (status IN ('active', 'suspended', 'left'));
ALTER TABLE organization_members ADD COLUMN IF NOT EXISTS last_active_at TIMESTAMP;

-- =====================================================
-- 3. Organization Join Requests Table
-- =====================================================

CREATE TABLE IF NOT EXISTS organization_join_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    
    -- Core relationships
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Request details
    requested_role VARCHAR(20) NOT NULL 
        CHECK (requested_role IN ('admin', 'member', 'viewer')),
    reason TEXT NOT NULL,
    department VARCHAR(255),
    manager_sponsor VARCHAR(255),
    
    -- Status tracking
    status VARCHAR(20) NOT NULL DEFAULT 'pending' 
        CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled', 'expired')),
    
    -- Review details
    reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMP,
    admin_notes TEXT,
    rejection_reason TEXT,
    approved_role VARCHAR(20) CHECK (approved_role IN ('admin', 'member', 'viewer')),
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '30 days'),
    
    -- Constraints
    CONSTRAINT valid_expiry CHECK (expires_at > created_at),
    CONSTRAINT one_pending_request_per_user_org UNIQUE (organization_id, user_id, status)
);

-- Indexes
CREATE INDEX idx_join_requests_org_id ON organization_join_requests(organization_id);
CREATE INDEX idx_join_requests_user_id ON organization_join_requests(user_id);
CREATE INDEX idx_join_requests_status ON organization_join_requests(status);
CREATE INDEX idx_join_requests_created_at ON organization_join_requests(created_at DESC);

-- =====================================================
-- 4. Organization Discovery Settings
-- =====================================================

CREATE TABLE IF NOT EXISTS organization_discovery_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    organization_id UUID NOT NULL UNIQUE REFERENCES organizations(id) ON DELETE CASCADE,
    
    -- Visibility settings
    is_discoverable BOOLEAN DEFAULT true,
    allow_join_requests BOOLEAN DEFAULT true,
    auto_approve BOOLEAN DEFAULT false,
    
    -- Request requirements
    require_reason BOOLEAN DEFAULT true,
    require_department BOOLEAN DEFAULT false,
    require_manager BOOLEAN DEFAULT false,
    
    -- Display settings
    display_name VARCHAR(255),
    description TEXT,
    show_member_count BOOLEAN DEFAULT true,
    show_plan BOOLEAN DEFAULT true,
    
    -- Approval settings
    approval_required_from VARCHAR(20) DEFAULT 'admin' 
        CHECK (approval_required_from IN ('owner', 'admin', 'any_admin')),
    request_expiry_days INTEGER DEFAULT 30,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create default settings for existing organizations
INSERT INTO organization_discovery_settings (organization_id)
SELECT id FROM organizations
ON CONFLICT (organization_id) DO NOTHING;

-- =====================================================
-- 5. User Notifications Table
-- =====================================================

CREATE TABLE IF NOT EXISTS user_notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Notification content
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    
    -- Related resource
    related_type VARCHAR(50),
    related_id UUID,
    
    -- Action
    action_url TEXT,
    action_label VARCHAR(100),
    
    -- Read status
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_notifications_user_id ON user_notifications(user_id);
CREATE INDEX idx_notifications_is_read ON user_notifications(is_read);
CREATE INDEX idx_notifications_created_at ON user_notifications(created_at DESC);
CREATE INDEX idx_notifications_user_unread ON user_notifications(user_id, is_read) 
    WHERE is_read = false;

-- =====================================================
-- 6. User Sessions Table (for org switching)
-- =====================================================

CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Current context
    current_organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    current_tenant_id UUID,
    
    -- Session tracking
    session_token VARCHAR(255) UNIQUE NOT NULL,
    ip_address INET,
    user_agent TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    last_active_at TIMESTAMP,
    
    UNIQUE(user_id, session_token)
);

CREATE INDEX idx_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_sessions_token ON user_sessions(session_token);
CREATE INDEX idx_sessions_expires_at ON user_sessions(expires_at);

-- =====================================================
-- 7. Functions
-- =====================================================

-- Update timestamp trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to tables
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_join_requests_updated_at
    BEFORE UPDATE ON organization_join_requests
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_discovery_settings_updated_at
    BEFORE UPDATE ON organization_discovery_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sessions_updated_at
    BEFORE UPDATE ON user_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create notification function
CREATE OR REPLACE FUNCTION create_notification(
    p_user_id UUID,
    p_type VARCHAR,
    p_title VARCHAR,
    p_message TEXT,
    p_related_type VARCHAR DEFAULT NULL,
    p_related_id UUID DEFAULT NULL,
    p_action_url TEXT DEFAULT NULL,
    p_action_label VARCHAR DEFAULT NULL
)
RETURNS UUID AS $$
DECLARE
    v_notification_id UUID;
BEGIN
    INSERT INTO user_notifications (
        user_id, type, title, message,
        related_type, related_id, action_url, action_label
    ) VALUES (
        p_user_id, p_type, p_title, p_message,
        p_related_type, p_related_id, p_action_url, p_action_label
    ) RETURNING id INTO v_notification_id;
    
    RETURN v_notification_id;
END;
$$ LANGUAGE plpgsql;

-- Notify admins when join request created
CREATE OR REPLACE FUNCTION notify_admins_of_join_request()
RETURNS TRIGGER AS $$
DECLARE
    v_admin RECORD;
    v_org_name VARCHAR;
    v_user_name VARCHAR;
    v_user_email VARCHAR;
BEGIN
    -- Get organization name
    SELECT name INTO v_org_name
    FROM organizations
    WHERE id = NEW.organization_id;
    
    -- Get user details
    SELECT name, email INTO v_user_name, v_user_email
    FROM users
    WHERE id = NEW.user_id;
    
    -- Notify all admins and owners
    FOR v_admin IN
        SELECT om.user_id
        FROM organization_members om
        WHERE om.organization_id = NEW.organization_id
          AND om.role IN ('owner', 'admin')
          AND om.status = 'active'
    LOOP
        PERFORM create_notification(
            v_admin.user_id,
            'join_request',
            'New Join Request',
            v_user_name || ' (' || v_user_email || ') wants to join ' || v_org_name || ' as ' || NEW.requested_role,
            'join_request',
            NEW.id,
            '/messages/join-requests/' || NEW.id::TEXT,
            'Review Request'
        );
    END LOOP;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_admins_on_join_request
    AFTER INSERT ON organization_join_requests
    FOR EACH ROW
    WHEN (NEW.status = 'pending')
    EXECUTE FUNCTION notify_admins_of_join_request();

-- Notify user when request is reviewed
CREATE OR REPLACE FUNCTION notify_user_of_request_decision()
RETURNS TRIGGER AS $$
DECLARE
    v_org_name VARCHAR;
    v_title VARCHAR;
    v_message TEXT;
BEGIN
    -- Only trigger on status change
    IF NEW.status = OLD.status THEN
        RETURN NEW;
    END IF;
    
    -- Get organization name
    SELECT name INTO v_org_name
    FROM organizations
    WHERE id = NEW.organization_id;
    
    IF NEW.status = 'approved' THEN
        v_title := 'Join Request Approved';
        v_message := 'Your request to join ' || v_org_name || ' has been approved! You now have access as ' || COALESCE(NEW.approved_role, NEW.requested_role) || '.';
        
        PERFORM create_notification(
            NEW.user_id,
            'join_request_approved',
            v_title,
            v_message,
            'organization',
            NEW.organization_id,
            '/organizations/' || NEW.organization_id::TEXT,
            'Go to Organization'
        );
        
    ELSIF NEW.status = 'rejected' THEN
        v_title := 'Join Request Rejected';
        v_message := 'Your request to join ' || v_org_name || ' has been rejected.';
        
        IF NEW.rejection_reason IS NOT NULL AND NEW.rejection_reason != '' THEN
            v_message := v_message || ' Reason: ' || NEW.rejection_reason;
        END IF;
        
        PERFORM create_notification(
            NEW.user_id,
            'join_request_rejected',
            v_title,
            v_message,
            NULL,
            NULL,
            NULL,
            NULL
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_user_on_request_decision
    AFTER UPDATE ON organization_join_requests
    FOR EACH ROW
    EXECUTE FUNCTION notify_user_of_request_decision();

-- Expire old join requests
CREATE OR REPLACE FUNCTION expire_old_join_requests()
RETURNS INTEGER AS $$
DECLARE
    v_count INTEGER;
BEGIN
    UPDATE organization_join_requests
    SET status = 'expired',
        updated_at = NOW()
    WHERE status = 'pending'
      AND expires_at < NOW();
    
    GET DIAGNOSTICS v_count = ROW_COUNT;
    RETURN v_count;
END;
$$ LANGUAGE plpgsql;

COMMIT;
```

---

## **API Endpoints**

### **Authentication APIs**

---

### **1. POST /api/v1/auth/register**

**Purpose:** User self-registration

**Auth Required:** No

**Request Body:**
```json
{
  "email": "john@gmail.com",
  "name": "John Developer",
  "password": "SecurePassword123!",
  "confirm_password": "SecurePassword123!"
}
```

**Validation Rules:**
```json
{
  "email": {
    "required": true,
    "format": "email",
    "unique": true,
    "max_length": 255
  },
  "name": {
    "required": true,
    "min_length": 2,
    "max_length": 255
  },
  "password": {
    "required": true,
    "min_length": 8,
    "must_contain": ["uppercase", "lowercase", "number", "special_char"]
  },
  "confirm_password": {
    "required": true,
    "must_match": "password"
  }
}
```

**Response (201):**
```json
{
  "user": {
    "id": "user_abc123",
    "email": "john@gmail.com",
    "name": "John Developer",
    "account_status": "active",
    "created_at": "2026-03-09T10:00:00Z"
  },
  "auth": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "refresh_xyz789",
    "token_type": "Bearer",
    "expires_in": 3600
  },
  "next_step": "find_organization",
  "redirect_url": "/find-organization"
}
```

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 400 | `email_required` | Email is required |
| 400 | `email_invalid` | Invalid email format |
| 400 | `password_weak` | Password must be at least 8 characters with uppercase, lowercase, number, and special character |
| 400 | `passwords_mismatch` | Passwords do not match |
| 409 | `email_exists` | An account with this email already exists |

---

### **2. POST /api/v1/auth/login**

**Purpose:** User login

**Auth Required:** No

**Request Body:**
```json
{
  "email": "john@gmail.com",
  "password": "SecurePassword123!"
}
```

**Response (200):**
```json
{
  "user": {
    "id": "user_abc123",
    "email": "john@gmail.com",
    "name": "John Developer",
    "account_status": "active"
  },
  "auth": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "refresh_xyz789",
    "token_type": "Bearer",
    "expires_in": 3600
  },
  "organizations": [
    {
      "id": "org_abc123",
      "name": "Google Cloud Platform",
      "role": "member",
      "is_current": true
    },
    {
      "id": "org_def456",
      "name": "Personal Projects",
      "role": "owner",
      "is_current": false
    }
  ],
  "redirect_url": "/dashboard"
}
```

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 400 | `invalid_credentials` | Invalid email or password |
| 403 | `account_suspended` | Your account has been suspended |

---

### **Organization Discovery APIs**

---

### **3. GET /api/v1/organizations/search**

**Purpose:** Search for organizations to join

**Auth Required:** Yes

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `q` | STRING | Yes | - | Search query (tenant name, org name) |
| `limit` | INTEGER | No | 10 | Max results |
| `offset` | INTEGER | No | 0 | Pagination offset |

**Request Example:**
```http
GET /api/v1/organizations/search?q=google&limit=10
Authorization: Bearer {access_token}
```

**Response (200):**
```json
{
  "query": "google",
  "results": [
    {
      "organization_id": "org_abc123",
      "tenant_id": "tenant_google",
      "name": "Google Cloud Platform",
      "tenant_name": "Google",
      "description": "Engineering team for Google Cloud Platform",
      "member_count": 12,
      "plan": "Enterprise",
      "is_discoverable": true,
      "allow_join_requests": true,
      "user_membership_status": null,
      "pending_request": null,
      "can_request_to_join": true
    },
    {
      "organization_id": "org_def456",
      "tenant_id": "tenant_google",
      "name": "Google Workspace Team",
      "tenant_name": "Google",
      "description": "Workspace administration and support",
      "member_count": 45,
      "plan": "Enterprise",
      "is_discoverable": true,
      "allow_join_requests": true,
      "user_membership_status": null,
      "pending_request": {
        "id": "req_xyz789",
        "status": "pending",
        "created_at": "2026-03-09T09:00:00Z"
      },
      "can_request_to_join": false
    },
    {
      "organization_id": "org_ghi789",
      "tenant_id": "tenant_google",
      "name": "Google Marketing",
      "tenant_name": "Google",
      "description": null,
      "member_count": null,
      "plan": null,
      "is_discoverable": true,
      "allow_join_requests": false,
      "user_membership_status": null,
      "pending_request": null,
      "can_request_to_join": false,
      "join_disabled_reason": "This organization has disabled join requests"
    }
  ],
  "total": 3,
  "limit": 10,
  "offset": 0
}
```

**Database Query:**
```sql
SELECT 
    o.id as organization_id,
    o.tenant_id,
    o.name,
    COALESCE(ods.display_name, o.name) as display_name,
    ods.description,
    CASE WHEN ods.show_member_count THEN 
        (SELECT COUNT(*) FROM organization_members WHERE organization_id = o.id AND status = 'active')
    ELSE NULL END as member_count,
    CASE WHEN ods.show_plan THEN o.plan ELSE NULL END as plan,
    ods.is_discoverable,
    ods.allow_join_requests,
    -- Check if user is already a member
    (SELECT role FROM organization_members 
     WHERE organization_id = o.id AND user_id = $2 AND status = 'active') as user_membership_status,
    -- Check for pending request
    (SELECT jsonb_build_object('id', id, 'status', status, 'created_at', created_at)
     FROM organization_join_requests 
     WHERE organization_id = o.id AND user_id = $2 AND status = 'pending'
     LIMIT 1) as pending_request
FROM organizations o
JOIN organization_discovery_settings ods ON o.id = ods.organization_id
WHERE ods.is_discoverable = true
  AND (
    LOWER(o.name) LIKE LOWER($1) OR
    LOWER(o.slug) LIKE LOWER($1) OR
    LOWER(ods.description) LIKE LOWER($1)
  )
ORDER BY o.name
LIMIT $3 OFFSET $4;
```

---

### **Join Request APIs**

---

### **4. POST /api/v1/organizations/{orgId}/join-requests**

**Purpose:** Create join request

**Auth Required:** Yes

**Request Body:**
```json
{
  "requested_role": "member",
  "reason": "I'm a contractor working on the authentication system project. I need access to track API usage for backend services.",
  "department": "Engineering - Backend",
  "manager_sponsor": "Sarah Smith"
}
```

**Validation:**
```json
{
  "requested_role": {
    "required": true,
    "enum": ["admin", "member", "viewer"],
    "note": "Cannot request 'owner' role"
  },
  "reason": {
    "required": true,
    "min_length": 10,
    "max_length": 500
  },
  "department": {
    "required": false,
    "max_length": 255
  },
  "manager_sponsor": {
    "required": false,
    "max_length": 255
  }
}
```

**Response (201):**
```json
{
  "join_request": {
    "id": "req_abc123",
    "organization_id": "org_abc123",
    "organization_name": "Google Cloud Platform",
    "user_id": "user_xyz789",
    "requested_role": "member",
    "reason": "I'm a contractor working on the authentication system project...",
    "department": "Engineering - Backend",
    "manager_sponsor": "Sarah Smith",
    "status": "pending",
    "created_at": "2026-03-09T10:30:00Z",
    "expires_at": "2026-04-08T10:30:00Z"
  },
  "message": "Your request has been sent to the organization admins. You'll be notified when they respond.",
  "next_steps": [
    "Organization admins will review your request",
    "You'll receive a notification when they respond",
    "Check 'My Requests' to track status"
  ]
}
```

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 400 | `already_member` | You are already a member of this organization |
| 400 | `pending_request_exists` | You already have a pending request for this organization |
| 400 | `reason_too_short` | Reason must be at least 10 characters |
| 403 | `join_requests_disabled` | This organization has disabled join requests |
| 404 | `organization_not_found` | Organization not found |
| 429 | `too_many_requests` | You've made too many requests recently. Please try again later. |

---

### **5. GET /api/v1/me/join-requests**

**Purpose:** Get user's own join requests

**Auth Required:** Yes

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `status` | STRING | No | all | Filter: `pending`, `approved`, `rejected`, `all` |
| `limit` | INTEGER | No | 20 | Results per page |
| `offset` | INTEGER | No | 0 | Pagination offset |

**Request Example:**
```http
GET /api/v1/me/join-requests?status=pending&limit=20
Authorization: Bearer {access_token}
```

**Response (200):**
```json
{
  "join_requests": [
    {
      "id": "req_abc123",
      "organization": {
        "id": "org_abc123",
        "name": "Google Cloud Platform",
        "tenant_name": "Google"
      },
      "requested_role": "member",
      "reason": "I'm a contractor working on...",
      "department": "Engineering - Backend",
      "manager_sponsor": "Sarah Smith",
      "status": "pending",
      "created_at": "2026-03-09T10:30:00Z",
      "expires_at": "2026-04-08T10:30:00Z",
      "days_until_expiry": 29
    },
    {
      "id": "req_def456",
      "organization": {
        "id": "org_def456",
        "name": "Microsoft Azure Team",
        "tenant_name": "Microsoft"
      },
      "requested_role": "viewer",
      "reason": "Need read-only access to analytics",
      "status": "approved",
      "approved_role": "viewer",
      "reviewed_by": {
        "name": "Admin User",
        "email": "admin@microsoft.com"
      },
      "reviewed_at": "2026-03-08T15:00:00Z",
      "created_at": "2026-03-08T10:00:00Z"
    }
  ],
  "pagination": {
    "total": 2,
    "limit": 20,
    "offset": 0
  },
  "summary": {
    "total": 2,
    "pending": 1,
    "approved": 1,
    "rejected": 0,
    "cancelled": 0,
    "expired": 0
  }
}
```

---

### **6. GET /api/v1/organizations/{orgId}/join-requests**

**Purpose:** Admin views join requests for their organization

**Auth Required:** Yes

**Required Permission:** `can_invite_members` (admin/owner only)

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `status` | STRING | No | pending | Filter by status |
| `limit` | INTEGER | No | 20 | Results per page |

**Request Example:**
```http
GET /api/v1/organizations/org_abc123/join-requests?status=pending
Authorization: Bearer {access_token}
```

**Response (200):**
```json
{
  "organization_id": "org_abc123",
  "organization_name": "Google Cloud Platform",
  "join_requests": [
    {
      "id": "req_abc123",
      "user": {
        "id": "user_xyz789",
        "email": "john@gmail.com",
        "name": "John Developer",
        "registered_at": "2026-03-09T10:00:00Z",
        "account_age_days": 0
      },
      "requested_role": "member",
      "reason": "I'm a contractor working on the authentication system project. I need access to track API usage for backend services.",
      "department": "Engineering - Backend",
      "manager_sponsor": "Sarah Smith",
      "status": "pending",
      "created_at": "2026-03-09T10:30:00Z",
      "expires_at": "2026-04-08T10:30:00Z",
      "days_since_request": 0,
      "hours_since_request": 2
    }
  ],
  "pagination": {
    "total": 1,
    "limit": 20,
    "offset": 0
  },
  "summary": {
    "total_pending": 1,
    "today": 1,
    "this_week": 1,
    "this_month": 1
  }
}
```

---

### **7. GET /api/v1/organizations/{orgId}/join-requests/{requestId}**

**Purpose:** View detailed join request

**Auth Required:** Yes

**Permission:** Admin of org OR requester themselves

**Request Example:**
```http
GET /api/v1/organizations/org_abc123/join-requests/req_abc123
Authorization: Bearer {access_token}
```

**Response (200):**
```json
{
  "id": "req_abc123",
  "organization": {
    "id": "org_abc123",
    "name": "Google Cloud Platform",
    "tenant_name": "Google"
  },
  "user": {
    "id": "user_xyz789",
    "email": "john@gmail.com",
    "name": "John Developer",
    "registered_at": "2026-03-09T10:00:00Z",
    "account_age_days": 0,
    "account_age_hours": 2
  },
  "requested_role": "member",
  "reason": "I'm a contractor working on the authentication system project. I need access to track API usage for backend services.",
  "department": "Engineering - Backend",
  "manager_sponsor": "Sarah Smith",
  "status": "pending",
  "admin_notes": null,
  "rejection_reason": null,
  "approved_role": null,
  "reviewed_by": null,
  "reviewed_at": null,
  "created_at": "2026-03-09T10:30:00Z",
  "updated_at": "2026-03-09T10:30:00Z",
  "expires_at": "2026-04-08T10:30:00Z"
}
```

---

### **8. POST /api/v1/organizations/{orgId}/join-requests/{requestId}/approve**

**Purpose:** Approve join request and add user to organization

**Auth Required:** Yes

**Required Permission:** `can_invite_members`

**Request Body:**
```json
{
  "role": "member",
  "admin_notes": "Verified with Sarah Smith - approved for backend team access"
}
```

**Request Example:**
```http
POST /api/v1/organizations/org_abc123/join-requests/req_abc123/approve
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "role": "member",
  "admin_notes": "Verified with Sarah - approved"
}
```

**Response (200):**
```json
{
  "join_request": {
    "id": "req_abc123",
    "status": "approved",
    "approved_role": "member",
    "reviewed_by": {
      "id": "user_admin123",
      "name": "Admin User",
      "email": "admin@google.com"
    },
    "reviewed_at": "2026-03-09T12:00:00Z",
    "admin_notes": "Verified with Sarah - approved"
  },
  "membership": {
    "id": "mem_new123",
    "user_id": "user_xyz789",
    "organization_id": "org_abc123",
    "role": "member",
    "status": "active",
    "joined_at": "2026-03-09T12:00:00Z"
  },
  "notification_sent": true,
  "message": "John Developer has been added to Google Cloud Platform as member"
}
```

**What Happens Backend:**
```go
// 1. Validate request exists and is pending
// 2. Check admin has permission
// 3. Create organization_member record
// 4. Update join_request status to 'approved'
// 5. Send notification to user
// 6. Log activity
// 7. Return response
```

---

### **9. POST /api/v1/organizations/{orgId}/join-requests/{requestId}/reject**

**Purpose:** Reject join request

**Auth Required:** Yes

**Required Permission:** `can_invite_members`

**Request Body:**
```json
{
  "rejection_reason": "We're not accepting contractors at this time. Please reapply in Q3 2026."
}
```

**Request Example:**
```http
POST /api/v1/organizations/org_abc123/join-requests/req_abc123/reject
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "rejection_reason": "We're not accepting contractors at this time."
}
```

**Response (200):**
```json
{
  "join_request": {
    "id": "req_abc123",
    "status": "rejected",
    "rejection_reason": "We're not accepting contractors at this time.",
    "reviewed_by": {
      "id": "user_admin123",
      "name": "Admin User"
    },
    "reviewed_at": "2026-03-09T12:05:00Z"
  },
  "notification_sent": true,
  "message": "Join request rejected"
}
```

---

### **10. DELETE /api/v1/join-requests/{requestId}**

**Purpose:** User cancels their own pending join request

**Auth Required:** Yes (must be request owner)

**Request Example:**
```http
DELETE /api/v1/join-requests/req_abc123
Authorization: Bearer {access_token}
```

**Response (200):**
```json
{
  "join_request_id": "req_abc123",
  "status": "cancelled",
  "cancelled_at": "2026-03-09T12:10:00Z",
  "message": "Join request cancelled successfully"
}
```

**Error Responses:**

| Status | Code | Message |
|--------|------|---------|
| 403 | `not_request_owner` | You can only cancel your own requests |
| 400 | `request_already_processed` | This request has already been approved/rejected |

---

### **Notification APIs**

---

### **11. GET /api/v1/me/notifications**

**Purpose:** Get user's notifications

**Auth Required:** Yes

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `is_read` | BOOLEAN | No | null | Filter: `true` (read), `false` (unread), `null` (all) |
| `type` | STRING | No | null | Filter by type |
| `limit` | INTEGER | No | 20 | Results per page |

**Request Example:**
```http
GET /api/v1/me/notifications?is_read=false&limit=20
Authorization: Bearer {access_token}
```

**Response (200):**
```json
{
  "notifications": [
    {
      "id": "notif_abc123",
      "type": "join_request_approved",
      "title": "Join Request Approved",
      "message": "Your request to join Google Cloud Platform has been approved! You now have access as member.",
      "is_read": false,
      "related_type": "organization",
      "related_id": "org_abc123",
      "action_url": "/organizations/org_abc123",
      "action_label": "Go to Organization",
      "created_at": "2026-03-09T12:00:00Z"
    },
    {
      "id": "notif_def456",
      "type": "join_request",
      "title": "New Join Request",
      "message": "John Developer (john@gmail.com) wants to join Google Cloud Platform as member",
      "is_read": false,
      "related_type": "join_request",
      "related_id": "req_xyz789",
      "action_url": "/messages/join-requests/req_xyz789",
      "action_label": "Review Request",
      "created_at": "2026-03-09T10:30:00Z"
    }
  ],
  "pagination": {
    "total": 2,
    "limit": 20,
    "offset": 0
  },
  "unread_count": 2
}
```

---

### **12. PUT /api/v1/me/notifications/{notificationId}/read**

**Purpose:** Mark notification as read

**Auth Required:** Yes

**Request Example:**
```http
PUT /api/v1/me/notifications/notif_abc123/read
Authorization: Bearer {access_token}
```

**Response (200):**
```json
{
  "notification_id": "notif_abc123",
  "is_read": true,
  "read_at": "2026-03-09T12:15:00Z"
}
```

---

### **13. PUT /api/v1/me/notifications/read-all**

**Purpose:** Mark all notifications as read

**Auth Required:** Yes

**Request Example:**
```http
PUT /api/v1/me/notifications/read-all
Authorization: Bearer {access_token}
```

**Response (200):**
```json
{
  "marked_read": 5,
  "message": "All notifications marked as read"
}
```

---

### **Team Member APIs**

---

### **14. GET /api/v1/organizations/{orgId}/members**

**Purpose:** List organization members

**Auth Required:** Yes

**Required Permission:** `can_view_projects` (all members)

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `role` | STRING | No | null | Filter by role |
| `status` | STRING | No | active | Filter by status |
| `limit` | INTEGER | No | 20 | Results per page |

**Request Example:**
```http
GET /api/v1/organizations/org_abc123/members?status=active&limit=50
Authorization: Bearer {access_token}
```

**Response (200):**
```json
{
  "organization_id": "org_abc123",
  "members": [
    {
      "id": "mem_owner123",
      "user_id": "user_owner123",
      "email": "owner@google.com",
      "name": "Organization Owner",
      "role": "owner",
      "status": "active",
      "joined_at": "2025-06-15T00:00:00Z",
      "last_active_at": "2026-03-09T11:30:00Z"
    },
    {
      "id": "mem_admin456",
      "user_id": "user_admin456",
      "email": "admin@google.com",
      "name": "Admin User",
      "role": "admin",
      "status": "active",
      "joined_at": "2025-07-01T00:00:00Z",
      "last_active_at": "2026-03-09T10:00:00Z"
    },
    {
      "id": "mem_new123",
      "user_id": "user_xyz789",
      "email": "john@gmail.com",
      "name": "John Developer",
      "role": "member",
      "status": "active",
      "joined_at": "2026-03-09T12:00:00Z",
      "last_active_at": null
    }
  ],
  "pagination": {
    "total": 3,
    "limit": 50,
    "offset": 0
  },
  "summary": {
    "total": 3,
    "owner": 1,
    "admin": 1,
    "member": 1,
    "viewer": 0
  }
}
```

---

### **15. PUT /api/v1/organizations/{orgId}/members/{memberId}/role**

**Purpose:** Change member's role

**Auth Required:** Yes

**Required Permission:** `can_change_member_roles`

**Request Body:**
```json
{
  "role": "admin"
}
```

**Response (200):**
```json
{
  "member_id": "mem_new123",
  "user_id": "user_xyz789",
  "old_role": "member",
  "new_role": "admin",
  "updated_at": "2026-03-09T13:00:00Z",
  "updated_by": {
    "id": "user_admin123",
    "name": "Admin User"
  },
  "message": "Role updated successfully"
}
```

---

### **16. DELETE /api/v1/organizations/{orgId}/members/{memberId}**

**Purpose:** Remove member from organization

**Auth Required:** Yes

**Required Permission:** `can_remove_members`

**Request Example:**
```http
DELETE /api/v1/organizations/org_abc123/members/mem_new123
Authorization: Bearer {access_token}
```

**Response (200):**
```json
{
  "member_id": "mem_new123",
  "user_id": "user_xyz789",
  "removed": true,
  "removed_at": "2026-03-09T13:05:00Z",
  "removed_by": {
    "id": "user_admin123",
    "name": "Admin User"
  },
  "message": "Member removed successfully"
}
```

---

### **Organization Settings APIs**

---

### **17. GET /api/v1/organizations/{orgId}/discovery-settings**

**Purpose:** Get organization discovery settings

**Auth Required:** Yes

**Required Permission:** `can_view_analytics`

**Request Example:**
```http
GET /api/v1/organizations/org_abc123/discovery-settings
Authorization: Bearer {access_token}
```

**Response (200):**
```json
{
  "organization_id": "org_abc123",
  "is_discoverable": true,
  "allow_join_requests": true,
  "auto_approve": false,
  "require_reason": true,
  "require_department": false,
  "require_manager": false,
  "display_name": "Google Cloud Platform",
  "description": "Engineering team for Google Cloud Platform",
  "show_member_count": true,
  "show_plan": true,
  "approval_required_from": "admin",
  "request_expiry_days": 30
}
```

---

### **18. PUT /api/v1/organizations/{orgId}/discovery-settings**

**Purpose:** Update organization discovery settings

**Auth Required:** Yes

**Required Permission:** `can_update_org`

**Request Body:**
```json
{
  "is_discoverable": true,
  "allow_join_requests": true,
  "auto_approve": false,
  "require_reason": true,
  "require_department": true,
  "require_manager": true,
  "display_name": "Google Cloud Platform",
  "description": "Engineering team for Google Cloud Platform - accepting new members!",
  "show_member_count": true,
  "show_plan": true
}
```

**Response (200):**
```json
{
  "organization_id": "org_abc123",
  "settings": {
    "is_discoverable": true,
    "allow_join_requests": true,
    "auto_approve": false,
    "require_reason": true,
    "require_department": true,
    "require_manager": true,
    "display_name": "Google Cloud Platform",
    "description": "Engineering team for Google Cloud Platform - accepting new members!",
    "show_member_count": true,
    "show_plan": true
  },
  "updated_at": "2026-03-09T13:10:00Z",
  "message": "Discovery settings updated successfully"
}
```

---

## **UI Specifications**

---

## **UI 1: User Registration Page**

**URL:** `/register`

**File:** `RegisterPage.jsx`

```jsx
┌─────────────────────────────────────────────────────────┐
│                                                         │
│              TrustKit Logo                              │
│                                                         │
│        Create Your Account                              │
│        Start tracking your API usage today              │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │ Full Name *                                     │   │
│  │ ┌─────────────────────────────────────────┐     │   │
│  │ │ John Developer                          │     │   │
│  │ └─────────────────────────────────────────┘     │   │
│  │                                                 │   │
│  │ Email Address *                                 │   │
│  │ ┌─────────────────────────────────────────┐     │   │
│  │ │ john@gmail.com                          │     │   │
│  │ └─────────────────────────────────────────┘     │   │
│  │                                                 │   │
│  │ Password *                                      │   │
│  │ ┌─────────────────────────────────────────┐     │   │
│  │ │ ••••••••••••••••            👁️         │     │   │
│  │ └─────────────────────────────────────────┘     │   │
│  │                                                 │   │
│  │ Password Requirements:                          │   │
│  │ ✅ At least 8 characters                        │   │
│  │ ✅ One uppercase letter                         │   │
│  │ ✅ One lowercase letter                         │   │
│  │ ✅ One number                                   │   │
│  │ ✅ One special character                        │   │
│  │                                                 │   │
│  │ Confirm Password *                              │   │
│  │ ┌─────────────────────────────────────────┐     │   │
│  │ │ ••••••••••••••••            👁️         │     │   │
│  │ └─────────────────────────────────────────┘     │   │
│  │                                                 │   │
│  │ ☑ I agree to the Terms of Service and          │   │
│  │   Privacy Policy                                │   │
│  │                                                 │   │
│  │ [        Create Account        ]                │   │
│  │                                                 │   │
│  │ Already have an account? [Sign In]             │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**Component Code:**
```jsx
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authAPI } from '@/api/auth';

export default function RegisterPage() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);

  const validatePassword = (password) => {
    const requirements = {
      minLength: password.length >= 8,
      hasUppercase: /[A-Z]/.test(password),
      hasLowercase: /[a-z]/.test(password),
      hasNumber: /[0-9]/.test(password),
      hasSpecial: /[!@#$%^&*]/.test(password)
    };
    return requirements;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setErrors({});

    try {
      const response = await authAPI.register(formData);
      
      // Store auth tokens
      localStorage.setItem('access_token', response.auth.access_token);
      localStorage.setItem('refresh_token', response.auth.refresh_token);
      
      // Redirect to find organization
      navigate('/find-organization');
    } catch (error) {
      setErrors(error.response?.data?.errors || { general: error.message });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full bg-white p-8 rounded-lg shadow">
        <h1 className="text-2xl font-bold text-center mb-2">Create Your Account</h1>
        <p className="text-gray-600 text-center mb-6">Start tracking your API usage today</p>
        
        <form onSubmit={handleSubmit}>
          {/* Form fields */}
        </form>
      </div>
    </div>
  );
}
```

---

## **UI 2: Find Organization Page**

**URL:** `/find-organization`

**File:** `FindOrganizationPage.jsx`

```jsx
┌─────────────────────────────────────────────────────────┐
│                                                         │
│  Welcome, John! 👋                                      │
│  Let's find your organization                           │
│                                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │ 🔍  Search by tenant or organization name...    │  │
│  └──────────────────────────────────────────────────┘  │
│  [Search]                                              │
│                                                         │
│  Popular Tenants:                                       │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐        │
│  │ 🏢   │ │ 🏢   │ │ 🏢   │ │ 🏢   │ │ 🏢   │        │
│  │Google│ │Microsoft│ │Amazon│ │ Meta │ │Apple │        │
│  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘        │
│                                                         │
│  OR                                                     │
│                                                         │
│  [+ Create Your Own Organization]                      │
│                                                         │
│  Don't see your organization?                          │
│  You can create a new one or request to join          │
│  an existing organization.                             │
└─────────────────────────────────────────────────────────┘
```

---

## **UI 3: Search Results Page**

**URL:** `/find-organization?q=google`

**File:** `SearchResultsPage.jsx`

```jsx
┌─────────────────────────────────────────────────────────┐
│  ← Back                                                 │
│                                                         │
│  Search Results for "google"                           │
│  Found 3 organizations                                 │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ 🏢 Google Cloud Platform                          │ │
│  │    Tenant: Google                                 │ │
│  │    Engineering team for Google Cloud Platform     │ │
│  │    👥 12 members • 💼 Enterprise plan             │ │
│  │                                                   │ │
│  │    [Request to Join]                              │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ 🏢 Google Workspace Team                          │ │
│  │    Tenant: Google                                 │ │
│  │    Workspace administration and support           │ │
│  │    👥 45 members • 💼 Enterprise plan             │ │
│  │                                                   │ │
│  │    ⏳ Request Pending (2 days ago)                │ │
│  │    [View Request]                                 │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ 🏢 Google Marketing                               │ │
│  │    Tenant: Google                                 │ │
│  │    🔒 Private - Join requests disabled            │ │
│  │    💼 Pro plan                                    │ │
│  │                                                   │ │
│  │    Contact an admin to join                       │ │
│  └───────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

**Component Code:**
```jsx
import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { organizationAPI } from '@/api/organizations';

export default function SearchResultsPage() {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('q');
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const searchOrganizations = async () => {
      setLoading(true);
      try {
        const data = await organizationAPI.search(query);
        setResults(data.results);
      } catch (error) {
        console.error('Search failed:', error);
      } finally {
        setLoading(false);
      }
    };

    if (query) {
      searchOrganizations();
    }
  }, [query]);

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1>Search Results for "{query}"</h1>
      <p>Found {results.length} organizations</p>
      
      <div className="space-y-4 mt-6">
        {results.map(org => (
          <OrganizationCard key={org.organization_id} organization={org} />
        ))}
      </div>
    </div>
  );
}
```

---

## **UI 4: Request to Join Modal**

**Component:** `RequestToJoinModal.jsx`

```jsx
┌─────────────────────────────────────────────────────────┐
│  Request to Join - Google Cloud Platform          ✕    │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Your Details:                                         │
│  Name: John Developer                                  │
│  Email: john@gmail.com                                 │
│                                                         │
│  What role are you requesting? *                       │
│  ○ Viewer   ● Member   ○ Admin                        │
│                                                         │
│  Viewer: Can view analytics and projects               │
│  Member: Can create projects and API keys              │
│  Admin: Can manage team and organization               │
│                                                         │
│  Why do you want to join this organization? *          │
│  ┌─────────────────────────────────────────────────┐  │
│  │ I'm a contractor working on the authentication  │  │
│  │ system project. I need access to track API      │  │
│  │ usage for the backend services.                 │  │
│  │                                        (124/500) │  │
│  └─────────────────────────────────────────────────┘  │
│  Minimum 10 characters required                        │
│                                                         │
│  Department/Team (optional)                            │
│  ┌─────────────────────────────────────────────────┐  │
│  │ Engineering - Backend                           │  │
│  └─────────────────────────────────────────────────┘  │
│                                                         │
│  Manager/Sponsor (optional)                            │
│  ┌─────────────────────────────────────────────────┐  │
│  │ Sarah Smith                                     │  │
│  └─────────────────────────────────────────────────┘  │
│                                                         │
│  ℹ️  Your request will be reviewed by organization     │
│     admins. You'll be notified when they respond.      │
│                                                         │
│              [Cancel]  [Send Request]                  │
└─────────────────────────────────────────────────────────┘
```

**Component Code:**
```jsx
export default function RequestToJoinModal({ organization, onClose, onSuccess }) {
  const [formData, setFormData] = useState({
    requested_role: 'member',
    reason: '',
    department: '',
    manager_sponsor: ''
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await organizationAPI.createJoinRequest(organization.id, formData);
      onSuccess();
      onClose();
    } catch (error) {
      console.error('Request failed:', error);
    }
  };

  return (
    <Modal onClose={onClose}>
      <form onSubmit={handleSubmit}>
        {/* Form fields */}
      </form>
    </Modal>
  );
}
```

---

## **UI 5: My Join Requests Page**

**URL:** `/me/join-requests`

**File:** `MyJoinRequestsPage.jsx`

```jsx
┌─────────────────────────────────────────────────────────┐
│  TrustKit                           [Profile ▼]         │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  My Join Requests                                      │
│  Track your organization access requests               │
│                                                         │
│  Filter: [All] [Pending] [Approved] [Rejected]        │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ 🏢 Google Cloud Platform                          │ │
│  │    Requested Role: Member                         │ │
│  │    Status: ⏳ Pending Review                      │ │
│  │    Submitted: 2 hours ago                         │ │
│  │    Expires in: 29 days                            │ │
│  │                                                   │ │
│  │    Reason: "I'm a contractor working on the      │ │
│  │    authentication system project..."              │ │
│  │                                                   │ │
│  │    [View Details]  [Cancel Request]              │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ 🏢 Microsoft Azure Team                           │ │
│  │    Requested Role: Viewer                         │ │
│  │    Status: ✅ Approved                            │ │
│  │    Approved Role: Viewer                          │ │
│  │    Reviewed: Yesterday at 3:30 PM                 │ │
│  │    Reviewed by: Admin User                        │ │
│  │                                                   │ │
│  │    [Go to Organization]                           │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ 🏢 Amazon AWS Team                                │ │
│  │    Requested Role: Member                         │ │
│  │    Status: ❌ Rejected                            │ │
│  │    Reviewed: 3 days ago                           │ │
│  │    Reason: "We're not accepting contractors       │ │
│  │    at this time."                                 │ │
│  │                                                   │ │
│  │    You can reapply in 27 days                     │ │
│  └───────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

---

## **UI 6: Messages Section (Admin View)**

**URL:** `/messages`

**File:** `MessagesPage.jsx`

**Sidebar:**
```
🏠 Dashboard
📊 Analytics
📁 Projects
🔑 API Keys
👥 Team
💬 Messages (1) ← Active, with badge
🏢 Organisation
👤 Profile
```

**Main Content:**
```jsx
┌─────────────────────────────────────────────────────────┐
│  Messages                                              │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Tabs: [Join Requests (1)] [Team Updates] [System]    │
│  ────────────────────────────────────────────          │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ 🔔 New Join Request                       2h ago  │ │
│  │                                                   │ │
│  │ 👤 John Developer                                 │ │
│  │    john@gmail.com                                 │ │
│  │    Wants to join as Member                        │ │
│  │                                                   │ │
│  │ 📝 Reason:                                        │ │
│  │ "I'm a contractor working on the authentication   │ │
│  │  system project. I need access to track API       │ │
│  │  usage for backend services."                     │ │
│  │                                                   │ │
│  │ 🏢 Department: Engineering - Backend              │ │
│  │ 👔 Manager: Sarah Smith                           │ │
│  │                                                   │ │
│  │ Account created: 2 hours ago (brand new)          │ │
│  │                                                   │ │
│  │ [View Full Request]  [✅ Approve]  [❌ Reject]    │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  No more pending requests                              │
└─────────────────────────────────────────────────────────┘
```

---

## **UI 7: Join Request Review Page**

**URL:** `/messages/join-requests/{requestId}`

**File:** `JoinRequestReviewPage.jsx`

```jsx
┌─────────────────────────────────────────────────────────┐
│  ← Back to Messages                                     │
│                                                         │
│  Join Request Review                                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  👤 Requestor Information                              │
│  ──────────────────────────────────────────            │
│  Name:           John Developer                        │
│  Email:          john@gmail.com                        │
│  Registered:     Mar 9, 2026 at 10:00 AM              │
│  Account Age:    2 hours (brand new account)           │
│                                                         │
│  🎯 Request Details                                    │
│  ──────────────────────────────────────────            │
│  Requested Role:  Member                               │
│  Submitted:       2 hours ago                          │
│  Expires:         In 29 days                           │
│                                                         │
│  📝 Reason for Joining                                 │
│  ──────────────────────────────────────────            │
│  "I'm a contractor working on the authentication       │
│   system project. I need access to track API usage     │
│   for the backend services."                           │
│                                                         │
│  🏢 Additional Information                             │
│  ──────────────────────────────────────────            │
│  Department:  Engineering - Backend                    │
│  Manager:     Sarah Smith                              │
│                                                         │
│  ⚙️  Admin Actions                                     │
│  ──────────────────────────────────────────            │
│                                                         │
│  Assign Role (defaults to requested):                  │
│  [Member ▼]                                            │
│                                                         │
│  Internal Notes (visible only to admins):              │
│  ┌─────────────────────────────────────────────────┐  │
│  │ Verified with Sarah - approved for backend      │  │
│  │ team access                                     │  │
│  └─────────────────────────────────────────────────┘  │
│                                                         │
│  ─── OR ───                                            │
│                                                         │
│  Rejection Reason (if rejecting):                      │
│  ┌─────────────────────────────────────────────────┐  │
│  │ We're not accepting contractors at this time.   │  │
│  │ Please reapply in Q3 2026.                      │  │
│  └─────────────────────────────────────────────────┘  │
│                                                         │
│         [❌ Reject Request]  [✅ Approve & Add]        │
└─────────────────────────────────────────────────────────┘
```

---

## **UI 8: Notification Bell**

**Component:** `NotificationBell.jsx`

**Location:** Top navbar

```jsx
┌──────────────────────────────────────┐
│  TrustKit    [Search]    🔔(1)  [👤]│
└──────────────────────────────────────┘
                            ↓
                   ┌──────────────────────────────────┐
                   │  Notifications                   │
                   ├──────────────────────────────────┤
                   │  ✅ Join Request Approved        │
                   │     Your request to join Google  │
                   │     Cloud Platform has been      │
                   │     approved!                    │
                   │     2 hours ago                  │
                   │     [Go to Organization]         │
                   │                                  │
                   │  ──────────────────────────      │
                   │                                  │
                   │  🔔 New Team Member              │
                   │     Sarah Smith joined your      │
                   │     organization                 │
                   │     Yesterday                    │
                   │                                  │
                   │  [Mark all as read]              │
                   │  [View all notifications]        │
                   └──────────────────────────────────┘
```

---

## **UI 9: Organization Switcher**

**Component:** `OrganizationSwitcher.jsx`

**Location:** Header (after user has organizations)

```jsx
┌────────────────────────────────────────────────┐
│  TrustKit   [Google Cloud Platform ▼]   🔔 👤│
└────────────────────────────────────────────────┘
                        ↓
            ┌──────────────────────────────────────┐
            │  Your Organizations                  │
            ├──────────────────────────────────────┤
            │  ✓ 🏢 Google Cloud Platform          │
            │     Member • 320K requests/month     │
            │                                      │
            │    🏢 Personal Projects              │
            │     Owner • 4.5K requests/month      │
            │                                      │
            │    🏢 Microsoft Azure Team           │
            │     Viewer • 1.2M requests/month     │
            │                                      │
            │  ──────────────────────────────      │
            │  [+ Find Organization]               │
            │  [⚙️ Organization Settings]          │
            └──────────────────────────────────────┘
```

---

## **UI 10: Team Members Page**

**URL:** `/organizations/{orgId}/team`

**File:** `TeamMembersPage.jsx`

```jsx
┌─────────────────────────────────────────────────────────┐
│  Team Members (3)                      [+ Add Member]   │
│  Manage your organization team                         │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Filter: [All Roles ▼] [Active ▼]  🔍 [Search...]     │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ 👤 Organization Owner              Owner          │ │
│  │    owner@google.com                               │ │
│  │    Joined 9 months ago • Active now               │ │
│  │                                                   │ │
│  │    Cannot remove owner                            │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ 👤 Admin User                      Admin          │ │
│  │    admin@google.com                               │ │
│  │    Joined 8 months ago • Active 10 mins ago       │ │
│  │                                                   │ │
│  │    [Change Role ▼]  [Remove]                     │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ 👤 John Developer                  Member         │ │
│  │    john@gmail.com                                 │ │
│  │    Joined 2 hours ago • Never active              │ │
│  │                                                   │ │
│  │    [Change Role ▼]  [Remove]                     │ │
│  └───────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

---

## **Implementation Checklist**

### **Backend (Go/Fiber)**

```
Database:
☐ Create migration 007_create_teams_service.sql
☐ Run migration
☐ Verify all tables created
☐ Test triggers and functions

API Endpoints - Auth:
☐ POST /api/v1/auth/register
☐ POST /api/v1/auth/login
☐ POST /api/v1/auth/logout
☐ POST /api/v1/auth/refresh-token

API Endpoints - Organizations:
☐ GET /api/v1/organizations/search
☐ GET /api/v1/organizations/{orgId}/discovery-settings
☐ PUT /api/v1/organizations/{orgId}/discovery-settings

API Endpoints - Join Requests:
☐ POST /api/v1/organizations/{orgId}/join-requests
☐ GET /api/v1/me/join-requests
☐ GET /api/v1/organizations/{orgId}/join-requests
☐ GET /api/v1/organizations/{orgId}/join-requests/{id}
☐ POST /api/v1/organizations/{orgId}/join-requests/{id}/approve
☐ POST /api/v1/organizations/{orgId}/join-requests/{id}/reject
☐ DELETE /api/v1/join-requests/{id}

API Endpoints - Notifications:
☐ GET /api/v1/me/notifications
☐ PUT /api/v1/me/notifications/{id}/read
☐ PUT /api/v1/me/notifications/read-all

API Endpoints - Team:
☐ GET /api/v1/organizations/{orgId}/members
☐ PUT /api/v1/organizations/{orgId}/members/{id}/role
☐ DELETE /api/v1/organizations/{orgId}/members/{id}

Middleware:
☐ JWT authentication middleware
☐ Permission checking middleware
☐ Rate limiting middleware
☐ Organization context middleware

Services:
☐ AuthService (register, login, tokens)
☐ OrganizationService (search, settings)
☐ JoinRequestService (create, approve, reject)
☐ NotificationService (create, send)
☐ TeamService (members CRUD)
```

### **Frontend (React)**

```
Pages:
☐ RegisterPage.jsx
☐ LoginPage.jsx
☐ FindOrganizationPage.jsx
☐ SearchResultsPage.jsx
☐ MyJoinRequestsPage.jsx
☐ MessagesPage.jsx (with tabs)
☐ JoinRequestReviewPage.jsx
☐ TeamMembersPage.jsx

Components:
☐ OrganizationCard.jsx
☐ RequestToJoinModal.jsx
☐ JoinRequestCard.jsx
☐ NotificationBell.jsx
☐ OrganizationSwitcher.jsx
☐ MemberCard.jsx
☐ ChangeRoleModal.jsx
☐ RemoveMemberModal.jsx

API Clients:
☐ api/auth.js (register, login, logout)
☐ api/organizations.js (search, settings)
☐ api/joinRequests.js (CRUD operations)
☐ api/notifications.js (list, mark read)
☐ api/team.js (members CRUD)

State Management:
☐ AuthContext (current user, tokens)
☐ OrganizationContext (current org)
☐ NotificationContext (unread count)

Routing:
☐ /register
☐ /login
☐ /find-organization
☐ /me/join-requests
☐ /messages
☐ /messages/join-requests/:id
☐ /organizations/:id/team

Utils:
☐ formatDate()
☐ formatRelativeTime()
☐ validateEmail()
☐ validatePassword()
☐ truncateText()
```

### **Testing**

```
Unit Tests:
☐ User registration validation
☐ Password hashing
☐ JWT token generation
☐ Join request creation
☐ Notification creation
☐ Permission checking

Integration Tests:
☐ Complete registration flow
☐ Search organizations flow
☐ Join request flow
☐ Approval flow
☐ Rejection flow
☐ Organization switching

E2E Tests:
☐ User registers → finds org → requests to join
☐ Admin reviews → approves request
☐ User sees notification → accesses org
☐ User switches between organizations
```

---

## **Summary**

This complete specification provides:

✅ **18 API Endpoints** for full teams functionality
✅ **10 UI Pages/Components** with detailed mockups
✅ **Complete database schema** with migrations
✅ **Self-service join request system** (no email required)
✅ **In-app notifications** (Messages section)
✅ **Multi-organization support** (org switcher)
✅ **Role-based permissions** (owner/admin/member/viewer)
✅ **Production-ready** error handling and validation

**Ready to implement! 🚀**
```

---

**This is your complete Teams Service specification! Everything you need to build the feature from database to UI!** 🎉