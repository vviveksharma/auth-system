# GuardRail

- This project delivers a highly modular authentication system engineered for effortless integration into diverse organizational applications.
- It features plug-and-play capabilities, enabling organizations to easily incorporate or remove authentication functionalities as their needs evolve.
- Key features include secure login and logout mechanisms, along with comprehensive route-based access control to ensure that only authorized users can access protected resources.
- Designed with scalability and security as top priorities, the system is ideal for production environments that demand reliability, maintainability, and adherence to industry best practices.

## API Endpoints

| Endpoint             | Method | Description                                                | Access        |
| -------------------- | ------ | ---------------------------------------------------------- | ------------- |
| `/auth/register`     | POST   | Register a new user (admin only).                          | Admin         |
| `/auth/login`        | POST   | Authenticate a user and issue a session/token.             | Guest         |
| `/auth/logout`       | POST   | Log out the current user and invalidate the session/token. | Authenticated |
| `/auth/logout`       | PUT    | To log out the existing user any role that is logged in    | Authenticated |
| `/users/`            | GET    | List all users (admin only).                               | Admin         |
| `/users/me`          | GET    | Retrieve the profile of the currently authenticated user.  | Authenticated |
| `/users/me`          | PUT    | Update the logged in user ( any role)                      | Authenticated |
| `/users/:id/roles`   | PUT    | Update the specific user role.                             | Admin         |
| `/roles`             | GET    | Get all the present roles ( default and custom )           | Admin         |
| `/roles`             | POST   | Create custom role (e.g., support_agent)                   | Admin         |
| `/roles`             | PUT    | Update the custom role routes                              | Admin         |
| `/roles/permissions` | PUT    | Update permissions for a role                              | Admin         |
| `/admin/logs`        | GET    | View system access logs                                    | Admin         |
| `/admin/stats`       | GET    | Get RBAC usage statistics                                  | Tenant-Admin  |

---

### `/auth/register`

Allows an `admin`-privileged user to register new users.

### `/auth/login`

Authenticates a user and provides access to the system as per tenant requirements.

### `/auth/logout`

Logs out the currently authenticated user, invalidating their session or token.

### `/users/`

Lists all users in the system. Accessible only by admin users.

### `/users/me`

Returns the profile information of the currently logged-in user.

### `/users/me`

Retrieve the profile of the currently authenticated user.

### `/users/:id/roles`

Update the specific user role.

### `/roles`

Get all the present roles ( default and custom ).

### `/roles`

Create custom role (e.g., support_agent).

### `/roles/permissions`

Update permissions for a role ( Custom roles ).

### `/admin/logs`

View system access logs. User based logs accessed by the admin

### `/admin/audit`

Get the tenant level audit logs.

---
