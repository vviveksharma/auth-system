# Database Index Optimization Analysis

## Fields to Index in dbmodels.go

### DBUser
- `Email` (line 16) - Add `gorm:"index"` - Used in GetUserByEmail query
- `TenantId` (line 14) - Add `gorm:"index"` - Used in ListUsersPaginated, tenant-specific queries
- `Status` (line 19) - Add `gorm:"index"` - Used in ListUsersPaginated with status filtering

### DBRoles
- `TenantId` (line 36) - Add `gorm:"index"` - Used in GetAllRoles, GetRolesByTenant
- `RoleId` (line 35) - Add `gorm:"index"` - Used in ChangeStatus, UpdateRoleDetails, DeleteRole
- `Role` (line 32) - Add `gorm:"index"` - Used in FindRoleId query
- `Status` (line 38) - Add `gorm:"index"` - Used in status-based filtering

### DBLogin
- `UserId` (line 63) - Add `gorm:"index"` - Used in GetUserById, Logout queries
- `TenantId` (line 62) - Add `gorm:"index"` - Used in tenant-specific filtering
- `Revoked` (line 69) - Add `gorm:"index"` - Used in session validation
- `ExpiresAt` (line 68) - Add `gorm:"index"` - Used in token expiry checks

### DBToken
- `TenantId` (line 99) - Add `gorm:"index"` - Used in ListTokensPaginated, ListTokens
- `IsActive` (line 106) - Add `gorm:"index"` - Used in VerifyToken, ListTokensPaginated
- `ExpiresAt` (line 105) - Add `gorm:"index"` - Used in VerifyToken expiry checks
- `ApplicationKey` (line 107) - Add `gorm:"index"` - Used in UpdateLoginToken, token filtering

### DBTenantLogin
- `Email` (line 122) - Add `gorm:"index"` - Used in GetDetailsByEmail query
- `TenantId` (line 123) - Add `gorm:"index"` - Used in tenant-specific queries
- `IsActive` (line 125) - Add `gorm:"index"` - Used in active session filtering

### DBRouteRole
- `RoleId` (line 138) - Add `gorm:"index"` - Used in FindByRoleId, UpdateRouteRole, DeleteAndUpdateRole
- `TenantId` (line 137) - Add `gorm:"index"` - Used in tenant-specific queries

### DBMessage
- `TenantId` (line 168) - Add `gorm:"index"` - Used in GetStatus, GetMessageByConditions
- `Status` (line 171) - Add `gorm:"index"` - Used in status-based filtering
- `UserEmail` (line 167) - Add `gorm:"index"` - Used in user-specific message queries

## Repo Layer Updates Required

### user.go (internal/repo/user.go)
- Line 52: GetUserByEmail - Benefits from Email + TenantId indexes
- Line 140: ListUsersPaginated - Benefits from TenantId + Status indexes

### login.go (internal/repo/login.go)
- Line 30: GetUserById - Benefits from UserId index
- Line 87: Logout - Benefits from UserId index

### role.go (internal/repo/role.go)
- Line 27: GetAllRoles - Benefits from TenantId index
- Line 51: FindRoleId - Benefits from Role index
- Line 103: ChangeStatus - Benefits from RoleId index
- Line 127: GetRolesByTenant - Benefits from TenantId index
- Line 139: UpdateRoleDetails - Benefits from RoleId + TenantId indexes

### token.go (internal/repo/token.go)
- Line 40: UpdateLoginToken - Benefits from TenantId + ApplicationKey indexes
- Line 83: ListTokensPaginated - Benefits from TenantId + IsActive + ApplicationKey indexes
- Line 112: ListTokens - Benefits from TenantId index
- Line 151: VerifyToken - Benefits from ExpiresAt + IsActive + ApplicationKey indexes
- Line 204: VerifyApplicationToken - Benefits from ExpiresAt + IsActive + ApplicationKey indexes

### tenant.go (internal/repo/tenant.go)
- Line 97: DeleteTenant - Benefits from TenantId indexes across all related tables
- Line 202: ListUserPaginated - Benefits from TenantId + Status indexes

### tenant_login.go (internal/repo/tenant_login.go)
- Line 33: GetDetailsByEmail - Benefits from Email index

### messages.go (internal/repo/messages.go)
- Line 33: GetStatus - Benefits from TenantId index
- Line 41: GetMessageByConditions - Benefits from TenantId + Status indexes

### route_role.go (internal/repo/route_role.go)
- Line 37: FindByRoleId - Benefits from RoleId index
- Line 51: UpdateRouteRole - Benefits from RoleId index
- Line 73: DeleteAndUpdateRole - Benefits from RoleId index
- Line 116: GetRoleRouteMapping - Benefits from RoleId index
