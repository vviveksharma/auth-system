-- =============================================================
-- Migration 000014: Seed gr.* system roles + route-role mapping
--
-- UUIDs are hardcoded UUIDv7 (time-ordered, sequential).
-- Timestamp prefix 0195948d-c000 ≈ 2026-03-19 UTC
--
-- System tenant: dae760ab-0a7f-4cbd-8603-def85ad8e430
--   (no FK constraint on role_tbl, tenant row not required)
--
-- Replaces: initsetup/initsetup.go GORM-based seeding.
-- Idempotent: all INSERTs use ON CONFLICT (id) DO NOTHING.
-- =============================================================

BEGIN;

-- =============================================================
-- 1.  role_tbl
-- =============================================================
INSERT INTO role_tbl
    (id, role_id, tenant_id, role, display_name, description, role_type, status, created_at, updated_at)
VALUES
  -- ── Tenant-level ─────────────────────────────────────────
  ('0195948d-c000-7001-8001-000000000001'::UUID, '0195948d-c000-7001-8001-000000000001'::UUID,
   'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
   'gr.tenant.owner', 'Account Owner',
   'Full account control: users, orgs, billing, deletion, and all administrative actions.',
   'default', true, now(), now()),

  ('0195948d-c000-7002-8001-000000000002'::UUID, '0195948d-c000-7002-8001-000000000002'::UUID,
   'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
   'gr.tenant.admin', 'Account Administrator',
   'Manages users, orgs, roles and tokens. No billing or account deletion.',
   'default', true, now(), now()),

  ('0195948d-c000-7003-8001-000000000003'::UUID, '0195948d-c000-7003-8001-000000000003'::UUID,
   'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
   'gr.tenant.auditor', 'Account Auditor',
   'Read-only visibility across the full tenant. Cannot modify anything.',
   'default', true, now(), now()),

  -- ── User-level ───────────────────────────────────────────
  ('0195948d-c000-7010-8001-000000000010'::UUID, '0195948d-c000-7010-8001-000000000010'::UUID,
   'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
   'gr.user', 'User',
   'Standard authenticated user. Self-service access, can request to join organizations.',
   'default', true, now(), now()),

  ('0195948d-c000-7011-8001-000000000011'::UUID, '0195948d-c000-7011-8001-000000000011'::UUID,
   'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
   'gr.guest', 'Guest',
   'Unauthenticated visitor. Can only register and login.',
   'default', true, now(), now()),

  -- ── Org-level ────────────────────────────────────────────
  ('0195948d-c001-7001-8002-000000000021'::UUID, '0195948d-c001-7001-8002-000000000021'::UUID,
   'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
   'gr.org.owner', 'Organization Owner',
   'Created the organization. Full control including deletion and all member management.',
   'default', true, now(), now()),

  ('0195948d-c001-7002-8002-000000000022'::UUID, '0195948d-c001-7002-8002-000000000022'::UUID,
   'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
   'gr.org.lead', 'Organization Lead',
   'Approves/rejects join requests, manages members, updates org settings. Cannot delete org.',
   'default', true, now(), now()),

  ('0195948d-c001-7003-8002-000000000023'::UUID, '0195948d-c001-7003-8002-000000000023'::UUID,
   'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
   'gr.org.member', 'Member',
   'Standard organization contributor.',
   'default', true, now(), now()),

  ('0195948d-c001-7004-8002-000000000024'::UUID, '0195948d-c001-7004-8002-000000000024'::UUID,
   'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
   'gr.org.observer', 'Observer',
   'Read-only access to the organization. Cannot modify members or settings.',
   'default', true, now(), now()),

  -- ── Platform root ─────────────────────────────────────────
  ('0195948d-c002-7001-8003-000000000031'::UUID, '0195948d-c002-7001-8003-000000000031'::UUID,
   'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
   'gr.system.root', 'Platform Root',
   'Internal GuardRail system identity. Never assignable to human users.',
   'default', true, now(), now())

ON CONFLICT (id) DO NOTHING;

-- =============================================================
-- 2.  route_role_tbl
-- =============================================================

-- gr.tenant.owner (priority 100)
INSERT INTO route_role_tbl (id, role_name, tenant_id, role_id, permissions, routes) VALUES (
  '0195948d-c003-7001-8004-000000000041'::UUID, 'gr.tenant.owner',
  'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
  '0195948d-c000-7001-8001-000000000001'::UUID,
  '{
    "role_info": {"name":"gr.tenant.owner","display_name":"Account Owner","role_type":"default","priority":100,"is_system":true},
    "permissions": [
      {"route":"/auth/",                                          "methods":["POST"],               "description":"Register user"},
      {"route":"/auth/login",                                     "methods":["POST"],               "description":"Login"},
      {"route":"/auth/logout",                                    "methods":["PUT"],                "description":"Logout"},
      {"route":"/auth/refresh",                                   "methods":["PUT"],                "description":"Refresh token"},
      {"route":"/users/me",                                       "methods":["GET","PUT"],          "description":"Own profile"},
      {"route":"/users/:id",                                      "methods":["GET","PUT","DELETE"], "description":"Full user management"},
      {"route":"/users",                                          "methods":["GET","POST"],         "description":"List and create users"},
      {"route":"/users/:id/roles",                                "methods":["PUT"],                "description":"Assign roles to users"},
      {"route":"/user/resetpassword",                             "methods":["POST"],               "description":"Reset passwords"},
      {"route":"/user/setpassword",                               "methods":["PUT"],                "description":"Set passwords"},
      {"route":"/roles",                                          "methods":["GET","POST"],         "description":"List and create roles"},
      {"route":"/roles/:id",                                      "methods":["GET","DELETE"],       "description":"Role detail and delete"},
      {"route":"/roles/:id/permissions",                          "methods":["GET","PUT"],          "description":"Manage role permissions"},
      {"route":"/roles/enable/:id",                               "methods":["PUT"],                "description":"Enable role"},
      {"route":"/roles/disable/:id",                              "methods":["PUT"],                "description":"Disable role"},
      {"route":"/request",                                        "methods":["POST","GET"],         "description":"Role assignment requests"},
      {"route":"/request/status",                                 "methods":["GET"],                "description":"Request status"},
      {"route":"/organizations",                                  "methods":["GET","POST"],         "description":"List and create orgs"},
      {"route":"/organizations/search",                           "methods":["GET"],                "description":"Search discoverable orgs"},
      {"route":"/organizations/:id",                              "methods":["GET","PUT","DELETE"], "description":"Full org CRUD"},
      {"route":"/organizations/:id/switch",                       "methods":["POST"],               "description":"Switch active org"},
      {"route":"/organizations/:id/join-requests",                "methods":["GET","POST"],         "description":"Submit and list join requests"},
      {"route":"/organizations/:id/join-requests/:reqId",         "methods":["GET"],                "description":"Join request detail"},
      {"route":"/organizations/:id/join-requests/:reqId/approve", "methods":["POST"],               "description":"Approve join request"},
      {"route":"/organizations/:id/join-requests/:reqId/reject",  "methods":["POST"],               "description":"Reject join request"},
      {"route":"/organizations/:id/discovery-settings",           "methods":["PUT"],                "description":"Update discovery settings"},
      {"route":"/me/join-requests",                               "methods":["GET"],                "description":"Own join requests"},
      {"route":"/join-requests/:id",                              "methods":["DELETE"],             "description":"Cancel own join request"},
      {"route":"/me/notifications",                               "methods":["GET","PUT"],          "description":"Notifications"}
    ]
  }'::jsonb,
  ARRAY['/auth/','/auth/login','/auth/logout','/auth/refresh','/users/me','/users/:id','/users','/users/:id/roles','/user/resetpassword','/user/setpassword','/roles','/roles/:id','/roles/:id/permissions','/roles/enable/:id','/roles/disable/:id','/request','/request/status','/organizations','/organizations/search','/organizations/:id','/organizations/:id/switch','/organizations/:id/join-requests','/organizations/:id/join-requests/:reqId','/organizations/:id/join-requests/:reqId/approve','/organizations/:id/join-requests/:reqId/reject','/organizations/:id/discovery-settings','/me/join-requests','/join-requests/:id','/me/notifications']
) ON CONFLICT (id) DO NOTHING;

-- gr.tenant.admin (priority 90)
INSERT INTO route_role_tbl (id, role_name, tenant_id, role_id, permissions, routes) VALUES (
  '0195948d-c003-7002-8004-000000000042'::UUID, 'gr.tenant.admin',
  'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
  '0195948d-c000-7002-8001-000000000002'::UUID,
  '{
    "role_info": {"name":"gr.tenant.admin","display_name":"Account Administrator","role_type":"default","priority":90,"is_system":true},
    "permissions": [
      {"route":"/auth/",                                          "methods":["POST"],               "description":"Register user"},
      {"route":"/auth/login",                                     "methods":["POST"],               "description":"Login"},
      {"route":"/auth/logout",                                    "methods":["PUT"],                "description":"Logout"},
      {"route":"/auth/refresh",                                   "methods":["PUT"],                "description":"Refresh token"},
      {"route":"/users/me",                                       "methods":["GET","PUT"],          "description":"Own profile"},
      {"route":"/users/:id",                                      "methods":["GET","PUT","DELETE"], "description":"Full user management"},
      {"route":"/users",                                          "methods":["GET","POST"],         "description":"List and create users"},
      {"route":"/users/:id/roles",                                "methods":["PUT"],                "description":"Assign roles to users"},
      {"route":"/user/resetpassword",                             "methods":["POST"],               "description":"Reset passwords"},
      {"route":"/user/setpassword",                               "methods":["PUT"],                "description":"Set passwords"},
      {"route":"/roles",                                          "methods":["GET","POST"],         "description":"List and create roles"},
      {"route":"/roles/:id",                                      "methods":["GET","DELETE"],       "description":"Role detail and delete"},
      {"route":"/roles/:id/permissions",                          "methods":["GET","PUT"],          "description":"Manage role permissions"},
      {"route":"/roles/enable/:id",                               "methods":["PUT"],                "description":"Enable role"},
      {"route":"/roles/disable/:id",                              "methods":["PUT"],                "description":"Disable role"},
      {"route":"/request",                                        "methods":["POST","GET"],         "description":"Role assignment requests"},
      {"route":"/request/status",                                 "methods":["GET"],                "description":"Request status"},
      {"route":"/organizations",                                  "methods":["GET","POST"],         "description":"List and create orgs"},
      {"route":"/organizations/search",                           "methods":["GET"],                "description":"Search orgs"},
      {"route":"/organizations/:id",                              "methods":["GET","PUT"],          "description":"View and update org (no delete)"},
      {"route":"/organizations/:id/switch",                       "methods":["POST"],               "description":"Switch active org"},
      {"route":"/organizations/:id/join-requests",                "methods":["GET","POST"],         "description":"Submit and list join requests"},
      {"route":"/organizations/:id/join-requests/:reqId",         "methods":["GET"],                "description":"Join request detail"},
      {"route":"/organizations/:id/join-requests/:reqId/approve", "methods":["POST"],               "description":"Approve join request"},
      {"route":"/organizations/:id/join-requests/:reqId/reject",  "methods":["POST"],               "description":"Reject join request"},
      {"route":"/organizations/:id/discovery-settings",           "methods":["PUT"],                "description":"Update discovery settings"},
      {"route":"/me/join-requests",                               "methods":["GET"],                "description":"Own join requests"},
      {"route":"/join-requests/:id",                              "methods":["DELETE"],             "description":"Cancel own join request"},
      {"route":"/me/notifications",                               "methods":["GET","PUT"],          "description":"Notifications"}
    ]
  }'::jsonb,
  ARRAY['/auth/','/auth/login','/auth/logout','/auth/refresh','/users/me','/users/:id','/users','/users/:id/roles','/user/resetpassword','/user/setpassword','/roles','/roles/:id','/roles/:id/permissions','/roles/enable/:id','/roles/disable/:id','/request','/request/status','/organizations','/organizations/search','/organizations/:id','/organizations/:id/switch','/organizations/:id/join-requests','/organizations/:id/join-requests/:reqId','/organizations/:id/join-requests/:reqId/approve','/organizations/:id/join-requests/:reqId/reject','/organizations/:id/discovery-settings','/me/join-requests','/join-requests/:id','/me/notifications']
) ON CONFLICT (id) DO NOTHING;

-- gr.tenant.auditor (priority 70)
INSERT INTO route_role_tbl (id, role_name, tenant_id, role_id, permissions, routes) VALUES (
  '0195948d-c003-7003-8004-000000000043'::UUID, 'gr.tenant.auditor',
  'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
  '0195948d-c000-7003-8001-000000000003'::UUID,
  '{
    "role_info": {"name":"gr.tenant.auditor","display_name":"Account Auditor","role_type":"default","priority":70,"is_system":true},
    "permissions": [
      {"route":"/auth/login",                      "methods":["POST"], "description":"Login"},
      {"route":"/auth/logout",                     "methods":["PUT"],  "description":"Logout"},
      {"route":"/auth/refresh",                    "methods":["PUT"],  "description":"Refresh token"},
      {"route":"/users/me",                        "methods":["GET"],  "description":"View own profile"},
      {"route":"/users/:id",                       "methods":["GET"],  "description":"View any user"},
      {"route":"/users",                           "methods":["GET"],  "description":"List users"},
      {"route":"/roles",                           "methods":["GET"],  "description":"List roles"},
      {"route":"/roles/:id",                       "methods":["GET"],  "description":"View role detail"},
      {"route":"/roles/:id/permissions",           "methods":["GET"],  "description":"View permissions"},
      {"route":"/request",                         "methods":["GET"],  "description":"View role requests"},
      {"route":"/request/status",                  "methods":["GET"],  "description":"View request status"},
      {"route":"/organizations",                   "methods":["GET"],  "description":"List orgs"},
      {"route":"/organizations/search",            "methods":["GET"],  "description":"Search orgs"},
      {"route":"/organizations/:id",               "methods":["GET"],  "description":"View org detail"},
      {"route":"/organizations/:id/join-requests", "methods":["GET"],  "description":"View join requests"},
      {"route":"/me/join-requests",                "methods":["GET"],  "description":"Own join requests"},
      {"route":"/me/notifications",                "methods":["GET"],  "description":"View notifications"}
    ]
  }'::jsonb,
  ARRAY['/auth/login','/auth/logout','/auth/refresh','/users/me','/users/:id','/users','/roles','/roles/:id','/roles/:id/permissions','/request','/request/status','/organizations','/organizations/search','/organizations/:id','/organizations/:id/join-requests','/me/join-requests','/me/notifications']
) ON CONFLICT (id) DO NOTHING;

-- gr.user (priority 10)
INSERT INTO route_role_tbl (id, role_name, tenant_id, role_id, permissions, routes) VALUES (
  '0195948d-c003-7010-8004-000000000050'::UUID, 'gr.user',
  'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
  '0195948d-c000-7010-8001-000000000010'::UUID,
  '{
    "role_info": {"name":"gr.user","display_name":"User","role_type":"default","priority":10,"is_system":true},
    "permissions": [
      {"route":"/auth/login",                       "methods":["POST"],       "description":"Login"},
      {"route":"/auth/logout",                      "methods":["PUT"],        "description":"Logout"},
      {"route":"/auth/refresh",                     "methods":["PUT"],        "description":"Refresh token"},
      {"route":"/users/me",                         "methods":["GET","PUT"],  "description":"View and update own profile"},
      {"route":"/user/resetpassword",               "methods":["POST"],       "description":"Reset own password"},
      {"route":"/user/setpassword",                 "methods":["PUT"],        "description":"Set own password"},
      {"route":"/request",                          "methods":["POST","GET"], "description":"Submit role assignment request"},
      {"route":"/request/status",                   "methods":["GET"],        "description":"Check request status"},
      {"route":"/organizations/search",             "methods":["GET"],        "description":"Search discoverable orgs"},
      {"route":"/organizations/:id",                "methods":["GET"],        "description":"View org detail"},
      {"route":"/organizations/:id/join-requests",  "methods":["POST"],       "description":"Submit join request"},
      {"route":"/me/join-requests",                 "methods":["GET"],        "description":"View own join requests"},
      {"route":"/join-requests/:id",                "methods":["DELETE"],     "description":"Cancel own join request"},
      {"route":"/me/notifications",                 "methods":["GET","PUT"],  "description":"Notifications"}
    ]
  }'::jsonb,
  ARRAY['/auth/login','/auth/logout','/auth/refresh','/users/me','/user/resetpassword','/user/setpassword','/request','/request/status','/organizations/search','/organizations/:id','/organizations/:id/join-requests','/me/join-requests','/join-requests/:id','/me/notifications']
) ON CONFLICT (id) DO NOTHING;

-- gr.guest (priority 0)
INSERT INTO route_role_tbl (id, role_name, tenant_id, role_id, permissions, routes) VALUES (
  '0195948d-c003-7011-8004-000000000051'::UUID, 'gr.guest',
  'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
  '0195948d-c000-7011-8001-000000000011'::UUID,
  '{
    "role_info": {"name":"gr.guest","display_name":"Guest","role_type":"default","priority":0,"is_system":true},
    "permissions": [
      {"route":"/auth/",               "methods":["POST"], "description":"Register new account"},
      {"route":"/auth/login",          "methods":["POST"], "description":"Login"},
      {"route":"/organizations/search","methods":["GET"],  "description":"Browse orgs before registering"}
    ]
  }'::jsonb,
  ARRAY['/auth/','/auth/login','/organizations/search']
) ON CONFLICT (id) DO NOTHING;

-- gr.org.owner (priority 80)
INSERT INTO route_role_tbl (id, role_name, tenant_id, role_id, permissions, routes) VALUES (
  '0195948d-c003-7021-8004-000000000061'::UUID, 'gr.org.owner',
  'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
  '0195948d-c001-7001-8002-000000000021'::UUID,
  '{
    "role_info": {"name":"gr.org.owner","display_name":"Organization Owner","role_type":"default","priority":80,"is_system":true},
    "permissions": [
      {"route":"/auth/login",                                     "methods":["POST"],               "description":"Login"},
      {"route":"/auth/logout",                                    "methods":["PUT"],                "description":"Logout"},
      {"route":"/auth/refresh",                                   "methods":["PUT"],                "description":"Refresh token"},
      {"route":"/users/me",                                       "methods":["GET","PUT"],          "description":"Own profile"},
      {"route":"/organizations/search",                           "methods":["GET"],                "description":"Search orgs"},
      {"route":"/organizations/:id",                              "methods":["GET","PUT","DELETE"], "description":"Full org CRUD incl. delete"},
      {"route":"/organizations/:id/switch",                       "methods":["POST"],               "description":"Switch active org"},
      {"route":"/organizations/:id/join-requests",                "methods":["GET"],                "description":"List join requests"},
      {"route":"/organizations/:id/join-requests/:reqId",         "methods":["GET"],                "description":"Join request detail"},
      {"route":"/organizations/:id/join-requests/:reqId/approve", "methods":["POST"],               "description":"Approve join request"},
      {"route":"/organizations/:id/join-requests/:reqId/reject",  "methods":["POST"],               "description":"Reject join request"},
      {"route":"/organizations/:id/discovery-settings",           "methods":["PUT"],                "description":"Update discovery settings"},
      {"route":"/me/join-requests",                               "methods":["GET"],                "description":"Own join requests"},
      {"route":"/join-requests/:id",                              "methods":["DELETE"],             "description":"Cancel own join request"},
      {"route":"/me/notifications",                               "methods":["GET","PUT"],          "description":"Notifications"}
    ]
  }'::jsonb,
  ARRAY['/auth/login','/auth/logout','/auth/refresh','/users/me','/organizations/search','/organizations/:id','/organizations/:id/switch','/organizations/:id/join-requests','/organizations/:id/join-requests/:reqId','/organizations/:id/join-requests/:reqId/approve','/organizations/:id/join-requests/:reqId/reject','/organizations/:id/discovery-settings','/me/join-requests','/join-requests/:id','/me/notifications']
) ON CONFLICT (id) DO NOTHING;

-- gr.org.lead (priority 60)
INSERT INTO route_role_tbl (id, role_name, tenant_id, role_id, permissions, routes) VALUES (
  '0195948d-c003-7022-8004-000000000062'::UUID, 'gr.org.lead',
  'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
  '0195948d-c001-7002-8002-000000000022'::UUID,
  '{
    "role_info": {"name":"gr.org.lead","display_name":"Organization Lead","role_type":"default","priority":60,"is_system":true},
    "permissions": [
      {"route":"/auth/login",                                     "methods":["POST"],       "description":"Login"},
      {"route":"/auth/logout",                                    "methods":["PUT"],        "description":"Logout"},
      {"route":"/auth/refresh",                                   "methods":["PUT"],        "description":"Refresh token"},
      {"route":"/users/me",                                       "methods":["GET","PUT"],  "description":"Own profile"},
      {"route":"/organizations/search",                           "methods":["GET"],        "description":"Search orgs"},
      {"route":"/organizations/:id",                              "methods":["GET","PUT"],  "description":"View and update org (no delete)"},
      {"route":"/organizations/:id/switch",                       "methods":["POST"],       "description":"Switch active org"},
      {"route":"/organizations/:id/join-requests",                "methods":["GET"],        "description":"List join requests"},
      {"route":"/organizations/:id/join-requests/:reqId",         "methods":["GET"],        "description":"Join request detail"},
      {"route":"/organizations/:id/join-requests/:reqId/approve", "methods":["POST"],       "description":"Approve join request"},
      {"route":"/organizations/:id/join-requests/:reqId/reject",  "methods":["POST"],       "description":"Reject join request"},
      {"route":"/organizations/:id/discovery-settings",           "methods":["PUT"],        "description":"Update discovery settings"},
      {"route":"/me/join-requests",                               "methods":["GET"],        "description":"Own join requests"},
      {"route":"/join-requests/:id",                              "methods":["DELETE"],     "description":"Cancel own join request"},
      {"route":"/me/notifications",                               "methods":["GET","PUT"],  "description":"Notifications"}
    ]
  }'::jsonb,
  ARRAY['/auth/login','/auth/logout','/auth/refresh','/users/me','/organizations/search','/organizations/:id','/organizations/:id/switch','/organizations/:id/join-requests','/organizations/:id/join-requests/:reqId','/organizations/:id/join-requests/:reqId/approve','/organizations/:id/join-requests/:reqId/reject','/organizations/:id/discovery-settings','/me/join-requests','/join-requests/:id','/me/notifications']
) ON CONFLICT (id) DO NOTHING;

-- gr.org.member (priority 30)
INSERT INTO route_role_tbl (id, role_name, tenant_id, role_id, permissions, routes) VALUES (
  '0195948d-c003-7023-8004-000000000063'::UUID, 'gr.org.member',
  'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
  '0195948d-c001-7003-8002-000000000023'::UUID,
  '{
    "role_info": {"name":"gr.org.member","display_name":"Member","role_type":"default","priority":30,"is_system":true},
    "permissions": [
      {"route":"/auth/login",               "methods":["POST"],       "description":"Login"},
      {"route":"/auth/logout",              "methods":["PUT"],        "description":"Logout"},
      {"route":"/auth/refresh",             "methods":["PUT"],        "description":"Refresh token"},
      {"route":"/users/me",                 "methods":["GET","PUT"],  "description":"Own profile"},
      {"route":"/organizations/search",     "methods":["GET"],        "description":"Search orgs"},
      {"route":"/organizations/:id",        "methods":["GET"],        "description":"View org detail"},
      {"route":"/organizations/:id/switch", "methods":["POST"],       "description":"Switch active org"},
      {"route":"/me/join-requests",         "methods":["GET"],        "description":"Own join requests"},
      {"route":"/join-requests/:id",        "methods":["DELETE"],     "description":"Cancel own join request"},
      {"route":"/me/notifications",         "methods":["GET","PUT"],  "description":"Notifications"}
    ]
  }'::jsonb,
  ARRAY['/auth/login','/auth/logout','/auth/refresh','/users/me','/organizations/search','/organizations/:id','/organizations/:id/switch','/me/join-requests','/join-requests/:id','/me/notifications']
) ON CONFLICT (id) DO NOTHING;

-- gr.org.observer (priority 10)
INSERT INTO route_role_tbl (id, role_name, tenant_id, role_id, permissions, routes) VALUES (
  '0195948d-c003-7024-8004-000000000064'::UUID, 'gr.org.observer',
  'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
  '0195948d-c001-7004-8002-000000000024'::UUID,
  '{
    "role_info": {"name":"gr.org.observer","display_name":"Observer","role_type":"default","priority":10,"is_system":true},
    "permissions": [
      {"route":"/auth/login",           "methods":["POST"], "description":"Login"},
      {"route":"/auth/logout",          "methods":["PUT"],  "description":"Logout"},
      {"route":"/auth/refresh",         "methods":["PUT"],  "description":"Refresh token"},
      {"route":"/users/me",             "methods":["GET"],  "description":"View own profile"},
      {"route":"/organizations/search", "methods":["GET"],  "description":"Search orgs"},
      {"route":"/organizations/:id",    "methods":["GET"],  "description":"View org detail"},
      {"route":"/me/notifications",     "methods":["GET"],  "description":"View notifications"}
    ]
  }'::jsonb,
  ARRAY['/auth/login','/auth/logout','/auth/refresh','/users/me','/organizations/search','/organizations/:id','/me/notifications']
) ON CONFLICT (id) DO NOTHING;

-- gr.system.root (priority 999 — never assign to humans)
INSERT INTO route_role_tbl (id, role_name, tenant_id, role_id, permissions, routes) VALUES (
  '0195948d-c003-7031-8004-000000000071'::UUID, 'gr.system.root',
  'dae760ab-0a7f-4cbd-8603-def85ad8e430'::UUID,
  '0195948d-c002-7001-8003-000000000031'::UUID,
  '{
    "role_info": {"name":"gr.system.root","display_name":"Platform Root","role_type":"default","priority":999,"is_system":true},
    "permissions": [
      {"route":"*","methods":["GET","POST","PUT","DELETE","PATCH"],"description":"Unrestricted platform access — internal use only"}
    ]
  }'::jsonb,
  ARRAY['*']
) ON CONFLICT (id) DO NOTHING;

COMMIT;
