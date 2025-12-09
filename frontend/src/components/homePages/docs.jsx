import { FiLock, FiUser, FiShield } from "react-icons/fi";

const Docs = () => {
  // Base URL for Swagger UI
  const swaggerBaseUrl = "http://localhost:8080/swagger/index.html";

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-stone-800 text-white py-6 rounded-xl">
        <div className="container mx-auto px-4">
          <h1 className="text-3xl font-bold">TrustKit API Documentation</h1>
          <p className="mt-2 text-stone-100">
            Secure, scalable authentication and role management for your
            applications
          </p>
        </div>
      </header>
      <div className="container mx-auto px-4 py-8 grid grid-cols-1 lg:grid-cols-4 gap-8">
        {/* Sidebar Navigation */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-lg shadow p-6 sticky top-8">
            <h2 className="text-xl font-semibold mb-4">API Endpoints</h2>
            <nav>
              <ul className="space-y-2">
                <li>
                  <button
                    onClick={() => {
                      const el = document.getElementById("authentication");
                      if (el) el.scrollIntoView({ behavior: "smooth" });
                    }}
                    className="flex items-center text-stone-500 hover:text-stone-900 cursor-pointer"
                  >
                    <FiLock className="mr-2" /> Authentication
                  </button>
                </li>
                <li>
                  <button
                    onClick={() => {
                      const el = document.getElementById("user-management");
                      if (el) el.scrollIntoView({ behavior: "smooth" });
                    }}
                    className="flex items-center text-stone-500 hover:text-stone-800 cursor-pointer"
                  >
                    <FiUser className="mr-2" /> User Management
                  </button>
                </li>
                <li>
                  <button
                    onClick={() => {
                      const el = document.getElementById("role-management");
                      if (el) el.scrollIntoView({ behavior: "smooth" });
                    }}
                    className="flex items-center text-stone-500 hover:text-stone-800 cursor-pointer"
                  >
                    <FiShield className="mr-2" /> Role Management
                  </button>
                </li>
                <li className="mt-6 pt-4 border-t border-gray-200">
                  <a
                    href={swaggerBaseUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-2 justify-center border align-middle select-none font-sans font-medium text-center duration-300 ease-in disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed focus:shadow-none text-sm py-2 px-4 shadow-sm hover:shadow-md bg-stone-800 hover:bg-stone-700 border-stone-900 text-stone-50 rounded-lg transition antialiased cursor-pointer"
                  >
                    View Swagger Docs
                  </a>
                </li>
              </ul>
            </nav>

            <div className="mt-8 pt-6 border-t border-gray-200">
              <h3 className="font-medium mb-2">System Features</h3>
              <ul className="text-sm space-y-1">
                <li className="flex items-center">
                  <span className="text-green-500 mr-2">✓</span> JWT
                  Authentication
                </li>
                <li className="flex items-center">
                  <span className="text-green-500 mr-2">✓</span> Role-Based
                  Access Control
                </li>
                <li className="flex items-center">
                  <span className="text-green-500 mr-2">✓</span> Redis Caching
                  Layer
                </li>
                <li className="flex items-center">
                  <span className="text-green-500 mr-2">✓</span> Password Reset
                  Flow
                </li>
                <li className="flex items-center">
                  <span className="text-green-500 mr-2">✓</span> Custom Role
                  Creation
                </li>
              </ul>
            </div>
          </div>
        </div>

        {/* Documentation Content */}
        <div className="lg:col-span-3 space-y-8">
          {/* Introduction */}
          <section className="bg-white rounded-lg shadow p-6">
            <h2 className="text-2xl font-bold mb-4">TrustKit API</h2>
            <p className="text-gray-700 mb-4">
              TrustKit provides enterprise-grade authentication and
              authorization services that can be easily integrated into any
              application.
            </p>
            <div className="bg-blue-50 border-l-4 border-blue-400 p-4">
              <p className="text-blue-700">
                <strong>Pro Tip:</strong> All endpoints require proper headers.
                Check the
                <a
                  href={swaggerBaseUrl}
                  className="text-blue-600 underline ml-1 mr-1"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Swagger documentation
                </a>
                for complete details.
              </p>
            </div>
          </section>

          {/* Authentication Section */}
          <section
            id="authentication"
            className="bg-white rounded-lg shadow p-6"
          >
            <h2 className="text-xl font-bold mb-4 flex items-center">
              <FiLock className="mr-2 text-stone-600" /> Authentication
            </h2>

            <div className="space-y-6">
              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">POST /auth/login</h3>
                    <p className="text-gray-600 mt-1">
                      Authenticate user credentials and receive JWT tokens
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                    Public
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Example Request
                  </h4>
                  <pre className="mt-1 bg-gray-100 p-3 rounded text-sm overflow-x-auto">
                    {`{
  "email": "user@example.com",
  "password": "securePassword123"
}`}
                  </pre>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/Auth/post_login`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">PUT /auth/refresh</h3>
                    <p className="text-gray-600 mt-1">
                      Refresh expired access tokens using a valid refresh token
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                    JWT Required
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">Headers</h4>
                  <pre className="mt-1 bg-gray-100 p-3 rounded text-sm overflow-x-auto">
                    Authorization: Bearer [refresh_token]
                  </pre>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/Auth/post_refresh_token`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>
            </div>
          </section>

          {/* User Management Section */}
          <section
            id="user-management"
            className="bg-white rounded-lg shadow p-6"
          >
            <h2 className="text-xl font-bold mb-4 flex items-center">
              <FiUser className="mr-2 text-stone-600" /> User Management
            </h2>

            <div className="space-y-6">
              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">GET /user/me</h3>
                    <p className="text-gray-600 mt-1">
                      Retrieve comprehensive profile information for the
                      authenticated user
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                    JWT Required
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Returns complete user profile details including personal
                      information, assigned roles, and account status
                    </li>
                    <li>
                      Automatically validates JWT token and retrieves current
                      user context from the authorization header
                    </li>
                    <li>
                      Ideal for populating user profiles, dashboards, and
                      personalization features
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/User/get_user_details`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">PUT /user/me</h3>
                    <p className="text-gray-600 mt-1">
                      Update profile information for the currently authenticated
                      user
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                    JWT Required
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Allows users to modify their own profile information
                      including name, contact details, and preferences
                    </li>
                    <li>
                      Supports partial updates - only included fields will be
                      modified
                    </li>
                    <li>
                      Automatically validates input data and applies business
                      logic constraints
                    </li>
                    <li>
                      Returns the updated user object with applied changes
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/User/put_user_details`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">GET /user/{'{id}'}</h3>
                    <p className="text-gray-600 mt-1">
                      Retrieve detailed user information by unique identifier
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    Admin | Moderator
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Provides comprehensive user details for administrative
                      purposes
                    </li>
                    <li>
                      Includes sensitive information such as account status,
                      role assignments, and audit trail data
                    </li>
                    <li>
                      Validates user permissions to ensure authorized access to
                      user records
                    </li>
                    <li>
                      Returns 404 error if the specified user ID does not exist
                      in the system
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/User/get_user__id_`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">DELETE /user/{'{id}'}</h3>
                    <p className="text-gray-600 mt-1">
                      Permanently remove a user account from the system
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    Admin | Moderator
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Completely removes user account and associated data from
                      the database
                    </li>
                    <li>
                      Performs cascading deletion of user-related records while
                      maintaining referential integrity
                    </li>
                    <li>
                      Includes safety checks to prevent accidental deletion of
                      critical system accounts
                    </li>
                    <li>
                      Returns confirmation message with details of the deletion
                      operation
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/User/delete_user__id_`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">
                      PUT /user/{'{id}'}/roles
                    </h3>
                    <p className="text-gray-600 mt-1">
                      Manage role assignments for specific users within the
                      system
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    Admin
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Request Body Example
                  </h4>
                  <pre className="mt-1 bg-gray-100 p-3 rounded text-sm overflow-x-auto">
                    {`{
  "role": "content-manager"
}`}
                  </pre>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Enables administrators to assign or modify roles for
                      specific users identified by their unique ID
                    </li>
                    <li>
                      Validates role existence and user permissions before
                      applying changes
                    </li>
                    <li>
                      Supports atomic role updates with immediate effect across
                      the system
                    </li>
                    <li>
                      Returns the updated user object reflecting the new role
                      assignments
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/User/post_user__id__role`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">
                      POST /user/resetpassword
                    </h3>
                    <p className="text-gray-600 mt-1">
                      Initiate password reset process for user accounts
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                    Public
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Request Body Example
                  </h4>
                  <pre className="mt-1 bg-gray-100 p-3 rounded text-sm overflow-x-auto">
                    {`{
  "email": "user@TrustKits.com"
}`}
                  </pre>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Initiates secure password reset workflow by generating and
                      delivering a time-sensitive OTP to the registered email
                      address
                    </li>
                    <li>
                      Implements rate limiting and security measures to prevent
                      abuse
                    </li>
                    <li>
                      Validates email existence while maintaining user privacy
                      through generic success responses
                    </li>
                    <li>
                      OTP tokens are securely handled and include expiration
                      timestamps
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/User/post_user_password_reset`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">
                      PUT /user/setpassword
                    </h3>
                    <p className="text-gray-600 mt-1">
                      Complete password reset process with OTP verification
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                    Public
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Request Body Example
                  </h4>
                  <pre className="mt-1 bg-gray-100 p-3 rounded text-sm overflow-x-auto">
                    {`{
  "email": "user@TrustKits.com",
  "otp": "123456",
  "new_password": "new_secure_password",
  "confirm_password": "new_secure_password"
}`}
                  </pre>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Verifies OTP authenticity and validity before allowing
                      password change
                    </li>
                    <li>
                      Enforces strong password policies and matches confirmation
                      fields
                    </li>
                    <li>
                      Invalidates the used OTP immediately after successful
                      password update
                    </li>
                    <li>
                      Securely hashes new password and updates user
                      authentication credentials
                    </li>
                    <li>
                      Can be combined with session invalidation for higher
                      security
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/User/post_user_password_set`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>
            </div>
          </section>

          {/* Role Management Section */}
          <section
            id="role-management"
            className="bg-white rounded-lg shadow p-6"
          >
            <h2 className="text-xl font-bold mb-4 flex items-center">
              <FiShield className="mr-2 text-stone-600" /> Role Management
            </h2>

            <div className="space-y-6">
              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">POST /roles/</h3>
                    <p className="text-gray-600 mt-1">
                      Create custom roles with granular permissions
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    Admin
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Request Body Example
                  </h4>
                  <pre className="mt-1 bg-gray-100 p-3 rounded text-sm overflow-x-auto">
                    {`{
  "name": "content-manager",
  "display_name": "Content Manager",
  "description": "Manages content creation and moderation",
  "permissions": [
    {
      "route": "/api/content",
      "methods": ["GET", "POST", "PUT"],
      "description": "Full content management access"
    }
  ]
}`}
                  </pre>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Creates a custom role with precisely defined permissions
                      for API endpoints
                    </li>
                    <li>
                      Users assigned to this role inherit all specified
                      permissions automatically
                    </li>
                    <li>
                      Returns the newly created role object with
                      system-generated ID
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/Roles/post_roles_`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">GET /roles</h3>
                    <p className="text-gray-600 mt-1">
                      Retrieve all roles with detailed information
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    Admin | Moderator | User
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Query Parameters
                  </h4>
                  <pre className="mt-1 bg-gray-100 p-3 rounded text-sm overflow-x-auto">
                    {`?type=custom      // Filter by role type (system|custom)
?status=enabled   // Filter by status (enabled|disabled)
?page=1           // Pagination page number
?limit=20         // Results per page`}
                  </pre>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Returns a paginated list of roles with comprehensive
                      details including permissions
                    </li>
                    <li>
                      System roles can be distinguished from custom roles using
                      filters
                    </li>
                    <li>
                      Response can include metadata for pagination (total count,
                      current page, etc.)
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/Roles/get_roles`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">
                      PUT /roles/{'{id}'}/permissions
                    </h3>
                    <p className="text-gray-600 mt-1">
                      Modify permissions for a specific role
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    Admin | Moderator
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Request Body Example
                  </h4>
                  <pre className="mt-1 bg-gray-100 p-3 rounded text-sm overflow-x-auto">
                    {`{
  "add_permissions": [
    {
      "route": "/api/reports",
      "methods": ["GET"],
      "description": "Access to reporting dashboard"
    }
  ],
  "remove_permissions": [
    {
      "route": "/api/content",
      "methods": ["DELETE"]
    }
  ]
}`}
                  </pre>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Atomically adds and removes permissions from the specified
                      role
                    </li>
                    <li>
                      Changes take effect immediately for all users assigned to
                      this role
                    </li>
                    <li>
                      Returns the updated role object with the modified
                      permissions set
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/Roles/put_roles__id__permissions`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">
                      GET /roles/{'{id}'}/permissions
                    </h3>
                    <p className="text-gray-600 mt-1">
                      List all permissions for a specific role
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    Admin | Moderator | User
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Returns a comprehensive list of all permissions granted to
                      the specified role
                    </li>
                    <li>
                      Includes detailed information about each permission
                      (route, HTTP methods, description)
                    </li>
                    <li>
                      Useful for auditing role permissions or displaying them in
                      administration interfaces
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/Roles/get_roles__id__permissions`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">
                      PUT /roles/enable/{'{id}'}
                    </h3>
                    <p className="text-gray-600 mt-1">
                      Activate a disabled role
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    Admin | Moderator
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Changes the status of a role from &apos;disabled&apos; to
                      &apos;enabled&apos;
                    </li>
                    <li>
                      Once enabled, the role can be assigned to users and its
                      permissions become active
                    </li>
                    <li>Returns the updated role object with the new status</li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/Roles/put_roles__id__enable`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">
                      PUT /roles/disable/{'{id}'}
                    </h3>
                    <p className="text-gray-600 mt-1">
                      Deactivate a role without deleting it
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    Admin | Moderator
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>
                      Changes the status of a role from &apos;enabled&apos; to
                      &apos;disabled&apos;
                    </li>
                    <li>
                      Disabled roles cannot be assigned to new users, and
                      existing assignments become inactive
                    </li>
                    <li>
                      Useful for temporarily restricting access without removing
                      role definitions
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/Roles/put_roles__id__disable`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>

              <div className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-medium text-lg">DELETE /roles/{'{id}'}</h3>
                    <p className="text-gray-600 mt-1">
                      Permanently remove a custom role
                    </p>
                  </div>
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                    Admin
                  </span>
                </div>
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700">
                    Explanation
                  </h4>
                  <ul className="mt-1 text-sm text-gray-600 list-disc pl-5 space-y-1">
                    <li>Completely removes a custom role from the system</li>
                    <li>
                      Automatically revokes this role from all users who had it
                      assigned
                    </li>
                    <li>
                      System roles cannot be deleted - this endpoint only works
                      for custom roles
                    </li>
                    <li>
                      Returns a confirmation message upon successful deletion
                    </li>
                  </ul>
                </div>
                <a
                  href={`${swaggerBaseUrl}#/Roles/delete_roles__id_`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-3 text-sm text-stone-600 hover:underline"
                >
                  View in Swagger →
                </a>
              </div>
            </div>
          </section>

          {/* Integration Guide */}
          <section className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-bold mb-4">Implementation Guide</h2>
            <div className="prose max-w-none">
              <h3 className="text-lg font-medium">Getting Started</h3>
              <ol className="list-decimal pl-5 space-y-2">
                <li>
                  Configure your application to point to our TrustKit endpoints
                </li>
                <li>Implement the login flow to obtain JWT tokens</li>
                <li>
                  Store refresh tokens securely (httpOnly cookies recommended)
                </li>
                <li>
                  Add the Authorization header to all authenticated requests
                </li>
              </ol>

              <h3 className="text-lg font-medium mt-6">
                Why Choose TrustKit?
              </h3>
              <ul className="list-disc pl-5 space-y-2">
                <li>
                  <strong>Production-ready design:</strong> Built with patterns
                  inspired by real-world SaaS systems
                </li>
                <li>
                  <strong>High Performance:</strong> Redis caching layer reduces
                  database load
                </li>
                <li>
                  <strong>Flexible:</strong> Custom roles and permissions for
                  a wide range of use cases
                </li>
                <li>
                  <strong>Secure by design:</strong> JWT-based authentication,
                  strong password flows, and rate limiting
                </li>
              </ul>

              <div className="mt-6 bg-stone-50 p-4 rounded-lg border border-stone-100">
                <h4 className="font-medium text-stone-800">
                  Ready to integrate?
                </h4>
                <p className="mt-1 text-stone-700">
                  Check our{" "}
                  <a
                    href={swaggerBaseUrl}
                    className="font-semibold underline"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    complete API documentation
                  </a>{" "}
                  or contact our support team for implementation assistance.
                </p>
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
};

export default Docs;
