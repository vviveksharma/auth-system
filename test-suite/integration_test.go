package testsuite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite tests both API and UI servers
type IntegrationTestSuite struct {
	suite.Suite
	apiBaseURL string
	uiBaseURL  string
	httpClient *http.Client

	// Test credentials and tokens
	tenantEmail string
	tenantPass  string
	tenantToken string

	userEmail    string
	userPass     string
	userJWTToken string
	userID       string
	appKey       string
	roleID       string
	messageID    string
	customRoleID string

	// Second user for role assignment tests
	secondUserID    string
	secondUserEmail string
	secondUserPass  string
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.apiBaseURL = "http://localhost:8080"
	s.uiBaseURL = "http://localhost:8081"
	s.httpClient = &http.Client{
		Timeout: 15 * time.Second,
	}

	// Generate unique test data for full flow testing
	timestamp := time.Now().Unix()
	s.tenantEmail = fmt.Sprintf("tenant_%d@test.com", timestamp)
	s.tenantPass = "TenantPass123!"
	s.userEmail = fmt.Sprintf("user_%d@test.com", timestamp)
	s.userPass = "UserPass123!"

	fmt.Println("\n🧪 Starting Integration Test Suite")
	fmt.Println("==================================") // Wait for servers
	s.waitForServer(s.apiBaseURL, "API Server")
	s.waitForServer(s.uiBaseURL, "UI Server")

	fmt.Println("\n✅ Both servers are ready!")
}

func (s *IntegrationTestSuite) waitForServer(baseURL, name string) {
	fmt.Printf("⏳ Waiting for %s at %s...\n", name, baseURL)
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := s.httpClient.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			fmt.Printf("✅ %s is ready!\n", name)
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	s.T().Fatalf("❌ %s did not start in time", name)
}

func (s *IntegrationTestSuite) makeRequest(method, url string, body interface{}, headers map[string]string) *http.Response {
	var bodyReader io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	assert.NoError(s.T(), err)

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		fmt.Println("the error : ", err)
	}
	assert.NoError(s.T(), err)
	return resp
}

func (s *IntegrationTestSuite) parseJSON(resp *http.Response) map[string]interface{} {
	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("⚠️  Failed to parse JSON response: %v\n", err)
		fmt.Printf("   Raw response: %s\n", string(body))
		return map[string]interface{}{"error": true, "message": "Failed to parse response"}
	}
	return result
}

// ==================== TEST SCENARIOS ====================

// SCENARIO 1: UI Server - Tenant Registration & Login
func (s *IntegrationTestSuite) Test01_UI_TenantRegistration() {
	fmt.Println("\n📝 Test 1: Tenant Registration")

	tenantData := map[string]interface{}{
		"name":     "Test Tenant Org",
		"email":    s.tenantEmail,
		"password": s.tenantPass,
		"campany":  "Test Company Inc",
	}

	resp := s.makeRequest("POST", s.uiBaseURL+"/tenant/", tenantData, nil)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should successfully register tenant")

	result := s.parseJSON(resp)
	fmt.Printf("   ✅ Tenant registered: %s\n", s.tenantEmail)
	// Note: tenant_id is NOT exposed for security reasons - it's encoded in the login token
	assert.NotNil(s.T(), result["message"], "Should have success message")
}

func (s *IntegrationTestSuite) Test02_UI_TenantLogin() {
	fmt.Println("\n🔐 Test 2: Tenant Login")

	loginData := map[string]interface{}{
		"email":    s.tenantEmail,
		"password": s.tenantPass,
	}

	resp := s.makeRequest("POST", s.uiBaseURL+"/tenant/login", loginData, nil)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should successfully login")

	result := s.parseJSON(resp)
	if data, ok := result["data"].(map[string]interface{}); ok {
		// Try "access_token" first, then fall back to "token" (UUID-based auth)
		if token, exists := data["access_token"]; exists && token != nil {
			s.tenantToken = token.(string)
		} else if token, exists := data["token"]; exists && token != nil {
			s.tenantToken = token.(string)
		}
		if s.tenantToken != "" && len(s.tenantToken) > 8 {
			fmt.Printf("   ✅ Tenant logged in successfully (token: %s...)\n", s.tenantToken[:8])
		} else if s.tenantToken != "" {
			fmt.Printf("   ✅ Tenant logged in successfully\n")
		}
	}

	// Note: tenant_id is NOT exposed - it's securely encoded in the token UUID
	assert.NotEmpty(s.T(), s.tenantToken, "Tenant token should not be empty")
}

// SCENARIO 2: UI Server - Application Token Management
func (s *IntegrationTestSuite) Test03_UI_CreateApplicationToken() {
	fmt.Println("\n🔑 Test 3: Get Application Token")

	// The tenant gets a default token upon creation
	// Let's fetch it to use for API tests
	headers := map[string]string{
		"Authorization": "Bearer " + s.tenantToken,
	}

	resp := s.makeRequest("GET", s.uiBaseURL+"/tenant/tokens", nil, headers)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should list tokens")

	result := s.parseJSON(resp)
	fmt.Println("the result : ", result)
	// Extract the first token's ID (which is the application key UUID)
	if data, ok := result["data"].(map[string]interface{}); ok {
		if dataArray, ok := data["data"].([]interface{}); ok && len(dataArray) > 0 {
			if token, ok := dataArray[0].(map[string]interface{}); ok {
				if tokenID, exists := token["token_id"]; exists && tokenID != nil {
					s.appKey = tokenID.(string)
					fmt.Printf("   ✅ Using application key: %s\n", s.appKey)
				}
			}
		}
	}

	assert.NotEmpty(s.T(), s.appKey, "Should have retrieved application key")
}

func (s *IntegrationTestSuite) Test04_UI_ListApplicationTokens() {
	fmt.Println("\n📋 Test 4: List Application Tokens")

	headers := map[string]string{
		"Authorization": "Bearer " + s.tenantToken,
	}

	resp := s.makeRequest("GET", s.uiBaseURL+"/tenant/tokens", nil, headers)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should list tokens")

	result := s.parseJSON(resp)

	// ✅ Fixed: Handle nested data structure
	if data, ok := result["data"].(map[string]interface{}); ok {
		if dataArray, ok := data["data"].([]interface{}); ok {
			fmt.Printf("   ✅ Found %v token(s)\n", len(dataArray))
		} else {
			fmt.Printf("   ⚠️  Tokens data array not found\n")
		}

		// Also show pagination info
		if pagination, ok := data["pagination"].(map[string]interface{}); ok {
			totalItems := pagination["total_items"]
			fmt.Printf("   📊 Total tokens: %v\n", totalItems)
		}
	} else if tokens, ok := result["tokens"].([]interface{}); ok {
		// Fallback for different response format
		fmt.Printf("   ✅ Found %v token(s)\n", len(tokens))
	} else {
		fmt.Printf("   ⚠️  Tokens response format unexpected: %v\n", result)
	}
}

func (s *IntegrationTestSuite) Test05_UI_GetTokenStatus() {
	fmt.Println("\n🔍 Test 5: Get Token Status")

	headers := map[string]string{
		"Authorization": "Bearer " + s.tenantToken,
	}

	url := fmt.Sprintf("%s/tenant/tokens/status?status=active", s.uiBaseURL)
	resp := s.makeRequest("GET", url, nil, headers)
	fmt.Println("the token response", resp)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should get token status")

	fmt.Println("   ✅ Token status retrieved")
}

// SCENARIO 3: API Server - User Registration & Authentication
func (s *IntegrationTestSuite) Test06_API_UserRegistration() {
	fmt.Println("\n👤 Test 6: User Registration")

	userData := map[string]interface{}{
		"name":     "Test User",
		"email":    s.userEmail,
		"password": s.userPass,
	}

	url := fmt.Sprintf("%s/auth/?application_key=%s", s.apiBaseURL, s.appKey)
	resp := s.makeRequest("POST", url, userData, nil)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should register user")

	result := s.parseJSON(resp)
	if data, ok := result["data"].(map[string]interface{}); ok {
		if userID, exists := data["user_id"]; exists && userID != nil {
			s.userID = userID.(string)
		}
	} else if userID, exists := result["user_id"]; exists && userID != nil {
		s.userID = userID.(string)
	}

	fmt.Printf("   ✅ User registered: %s\n", s.userEmail)
	if s.userID != "" {
		fmt.Printf("   User ID: %s\n", s.userID)
	}
}

func (s *IntegrationTestSuite) Test07_API_UserLogin() {
	fmt.Println("\n👤 Test 7: User Login with the guest role")

	loginData := map[string]interface{}{
		"email":    s.userEmail,
		"password": s.userPass,
		"role":     "guest",
	}

	loginURL := fmt.Sprintf("%s/auth/login?application_key=%s", s.apiBaseURL, s.appKey)
	loginResp := s.makeRequest("POST", loginURL, loginData, nil)
	assert.Equal(s.T(), http.StatusOK, loginResp.StatusCode, "Should login successfully")

	loginResult := s.parseJSON(loginResp)
	if data, ok := loginResult["data"].(map[string]interface{}); ok {
		if token, exists := data["jwt"]; exists && token != nil {
			s.userJWTToken = token.(string)
		}
	}
	assert.NotEmpty(s.T(), s.userJWTToken, "Access token should not be empty")
	fmt.Printf("   ✅ User Login successfully \n")
}

func (s *IntegrationTestSuite) Test08_API_UserRequestRole() {
	fmt.Println("\n👤 Test 8: User request Role as admin")

	userData := map[string]interface{}{
		"email":          s.userEmail,
		"requested_role": "admin",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/request?application_key=%s", s.apiBaseURL, s.appKey)
	resp := s.makeRequest("POST", url, userData, headers)
	requestResult := s.parseJSON(resp)
	fmt.Println("the result: ", requestResult)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "message sent to the tenant-admin")

	// ✅ Wait for queue to process the message
	fmt.Println("   ⏳ Waiting for queue to process message...")
	time.Sleep(2 * time.Second)
}

func (s *IntegrationTestSuite) Test09_API_GetMessages() {
	fmt.Println("\n👤 Test 9: Getting the Messages information")

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/request?application_key=%s&email=%s", s.apiBaseURL, s.appKey, s.userEmail)
	resp := s.makeRequest("GET", url, nil, headers)
	requestResult := s.parseJSON(resp)

	// ✅ Extract and store message_id
	if data, ok := requestResult["data"].([]interface{}); ok && len(data) > 0 {
		fmt.Printf("   ✅ Messages count: %d\n", len(data))

		// Get the first message and extract its ID
		if msgMap, ok := data[0].(map[string]interface{}); ok {
			if msgID, exists := msgMap["message_id"]; exists && msgID != nil {
				s.messageID = msgID.(string)
				fmt.Printf("   ✅ Message ID stored: %s\n", s.messageID)
			}

			// Print all message details
			fmt.Printf("   📊 Status: %v\n", msgMap["status"])
			fmt.Printf("   👤 Requested Role: %v\n", msgMap["requested_role"])
		}
	} else {
		fmt.Printf("   ⚠️  No messages found or data is not an array\n")
		fmt.Printf("   Debug - Full response: %+v\n", requestResult)
	}

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "messages listed for this user")
	assert.NotEmpty(s.T(), s.messageID, "Should have extracted message ID")
}

func (s *IntegrationTestSuite) Test10_API_GetRequestStatus() {
	fmt.Println("\n👤 Test 10: Get message status using messageId")

	if s.messageID == "" {
		s.T().Skip("No message ID available from previous test")
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/request/status?application_key=%s&id=%s", s.apiBaseURL, s.appKey, s.messageID)
	resp := s.makeRequest("GET", url, nil, headers)
	requestResult := s.parseJSON(resp)

	fmt.Printf("   📊 Status response: %+v\n", requestResult)

	// ✅ Fixed assert - use Equal instead of Contains for status code
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should get request status")

	// Extract and display status
	if message, ok := requestResult["message"].(string); ok {
		fmt.Printf("   ✅ Request status: %s\n", message)
	}
}

func (s *IntegrationTestSuite) Test11_UI_ApproveAdminRequest() {
	fmt.Println("\n🔐 Test 11: Tenant approves the user request for the admin role")

	headers := map[string]string{
		"Authorization": "Bearer " + s.tenantToken,
	}

	url := fmt.Sprintf("%s/tenant/messages/approve?id=%s", s.uiBaseURL, s.messageID)
	resp := s.makeRequest("PUT", url, nil, headers)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should approve message successfully")

	result := s.parseJSON(resp)
	fmt.Printf("   ✅ Approval response: %v\n", result["message"])

	// Wait for role assignment to process
	fmt.Println("   ⏳ Waiting for role assignment...")
	time.Sleep(1 * time.Second)
}

// SCENARIO 4: User Management - Admin Operations
func (s *IntegrationTestSuite) Test12_API_UserReloginWithAdminRole() {
	fmt.Println("\n🔄 Test 12: User re-login with newly granted admin role")

	loginData := map[string]interface{}{
		"email":    s.userEmail,
		"password": s.userPass,
		"role":     "admin",
	}

	// ✅ Fixed: Add application_key query parameter
	loginURL := fmt.Sprintf("%s/auth/login?application_key=%s", s.apiBaseURL, s.appKey)
	loginResp := s.makeRequest("POST", loginURL, loginData, nil)
	assert.Equal(s.T(), http.StatusOK, loginResp.StatusCode, "Should login successfully")

	loginResult := s.parseJSON(loginResp)
	fmt.Println(loginResult)
	if data, ok := loginResult["data"].(map[string]interface{}); ok {
		if token, exists := data["jwt"]; exists && token != nil {
			s.userJWTToken = token.(string)
			fmt.Printf("   ✅ User logged in successfully with admin role\n")
		}
	}
	assert.NotEmpty(s.T(), s.userJWTToken, "Access token should not be empty")
}

func (s *IntegrationTestSuite) Test13_API_ListUsers() {
	fmt.Println("\n👥 Test 13: List all users")

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	// List users endpoint with pagination
	url := fmt.Sprintf("%s/users?application_key=%s&page=1&page_size=10", s.apiBaseURL, s.appKey)
	resp := s.makeRequest("GET", url, nil, headers)
	result := s.parseJSON(resp)
	fmt.Printf("   📦 Full response: %+v\n", result)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should list users successfully")

	// Extract user ID from the list
	if data, ok := result["data"].(map[string]interface{}); ok {
		if users, ok := data["data"].([]interface{}); ok {
			fmt.Printf("   ✅ Found %d user(s)\n", len(users))
			// Find our test user and extract ID
			for _, user := range users {
				if userMap, ok := user.(map[string]interface{}); ok {
					if email, ok := userMap["email"].(string); ok && email == s.userEmail {
						if userID, exists := userMap["id"]; exists && userID != nil {
							s.userID = userID.(string)
							fmt.Printf("   ✅ Extracted User ID: %s for email: %s\n", s.userID, email)
						}
					}
				}
			}
		}
	}

	assert.NotEmpty(s.T(), s.userID, "Should extract user ID from users list")
}

func (s *IntegrationTestSuite) Test14_API_GetUserProfile() {
	fmt.Println("\n👤 Test 14: Get current user profile (/users/me)")

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/users/me?application_key=%s", s.apiBaseURL, s.appKey)
	resp := s.makeRequest("GET", url, nil, headers)
	result := s.parseJSON(resp)
	fmt.Printf("   📦 Response: %+v\n", result)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should get user profile")

	if data, ok := result["data"].(map[string]interface{}); ok {
		fmt.Printf("   ✅ User Email: %v\n", data["email"])
		fmt.Printf("   ✅ User Name: %v\n", data["name"])
		if roles, ok := data["role"].([]interface{}); ok {
			fmt.Printf("   ✅ User Roles: %v\n", roles)
			// Verify admin role is present
			hasAdmin := false
			for _, role := range roles {
				if role == "admin" {
					hasAdmin = true
					break
				}
			}
			assert.True(s.T(), hasAdmin, "User should have admin role")
		}
	}
}

func (s *IntegrationTestSuite) Test15_API_GetUserById() {
	fmt.Println("\n🔍 Test 15: Get user by ID (/users/:id)")

	if s.userID == "" {
		s.T().Skip("No user ID available from previous test")
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/users/%s?application_key=%s", s.apiBaseURL, s.userID, s.appKey)
	fmt.Printf("   🔗 Request URL: %s\n", url)
	resp := s.makeRequest("GET", url, nil, headers)
	result := s.parseJSON(resp)
	fmt.Printf("   📦 Response: %+v\n", result)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should get user by ID")

	if data, ok := result["data"].(map[string]interface{}); ok {
		fmt.Printf("   ✅ User ID: %v\n", data["id"])
		fmt.Printf("   ✅ User Email: %v\n", data["email"])
		fmt.Printf("   ✅ User Name: %v\n", data["name"])
	}
}

func (s *IntegrationTestSuite) Test16_API_UpdateUserProfile() {
	fmt.Println("\n✏️  Test 16: Update current user profile")

	updateData := map[string]interface{}{
		"name": "Updated Test User",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/users/me?application_key=%s", s.apiBaseURL, s.appKey)
	resp := s.makeRequest("PUT", url, updateData, headers)
	result := s.parseJSON(resp)
	fmt.Printf("   📦 Response: %+v\n", result)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should update user profile")
	fmt.Printf("   ✅ Profile updated: %v\n", result["message"])
}

func (s *IntegrationTestSuite) Test17_API_CreateCustomRole() {
	fmt.Println("\n➕ Test 18: Create a custom role")

	roleData := map[string]interface{}{
		"role":         "manager",
		"display_name": "Manager",
		"description":  "Manager role for integration tests",
		"permissions": []map[string]interface{}{
			{
				"route":       "/user",
				"methods":     []string{"GET"},
				"description": "View users",
			},
			{
				"route":       "/roles",
				"methods":     []string{"GET"},
				"description": "View roles",
			},
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/roles?application_key=%s", s.apiBaseURL, s.appKey)
	resp := s.makeRequest("POST", url, roleData, headers)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should create custom role")

	result := s.parseJSON(resp)
	fmt.Printf("   ✅ Custom role created: %v\n", result["message"])

	// Extract role ID
	if data, ok := result["data"].(map[string]interface{}); ok {
		if roleID, exists := data["role_id"]; exists && roleID != nil {
			s.customRoleID = roleID.(string)
			fmt.Printf("   ✅ Custom role ID: %s\n", s.customRoleID)
		}
	}
}

func (s *IntegrationTestSuite) Test18_API_GetRolePermissions() {
	fmt.Println("\n🔐 Test 19: Get role permissions")

	if s.customRoleID == "" {
		s.T().Skip("No custom role ID available")
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/roles/%s/permissions?application_key=%s", s.apiBaseURL, s.customRoleID, s.appKey)
	resp := s.makeRequest("GET", url, nil, headers)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should get role permissions")

	result := s.parseJSON(resp)
	if data, ok := result["data"].(map[string]interface{}); ok {
		if perms, ok := data["permissions"].([]interface{}); ok {
			fmt.Printf("   ✅ Role has %d permission(s)\n", len(perms))
		}
	}
}

func (s *IntegrationTestSuite) Test19_API_UpdateRolePermissions() {
	fmt.Println("\n✏️  Test 20: Update role permissions")

	if s.customRoleID == "" {
		s.T().Skip("No custom role ID available")
	}

	updateData := map[string]interface{}{
		"permissions": []map[string]interface{}{
			{
				"route":       "/user",
				"methods":     []string{"GET", "PUT"},
				"description": "View and update users",
			},
			{
				"route":       "/roles",
				"methods":     []string{"GET", "POST"},
				"description": "View and create roles",
			},
			{
				"route":       "/request",
				"methods":     []string{"GET"},
				"description": "View role requests",
			},
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/roles/%s/permissions?application_key=%s", s.apiBaseURL, s.customRoleID, s.appKey)
	resp := s.makeRequest("PUT", url, updateData, headers)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should update role permissions")

	result := s.parseJSON(resp)
	fmt.Printf("   ✅ Permissions updated: %v\n", result["message"])
}

func (s *IntegrationTestSuite) Test20_API_DisableCustomRole() {
	fmt.Println("\n🚫 Test 21: Disable custom role")

	if s.customRoleID == "" {
		s.T().Skip("No custom role ID available")
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/roles/disable/%s?application_key=%s", s.apiBaseURL, s.customRoleID, s.appKey)
	resp := s.makeRequest("PUT", url, nil, headers)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should disable role")

	result := s.parseJSON(resp)
	fmt.Printf("   ✅ Role disabled: %v\n", result["message"])
}

func (s *IntegrationTestSuite) Test21_API_EnableCustomRole() {
	fmt.Println("\n✅ Test 22: Enable custom role")

	if s.customRoleID == "" {
		s.T().Skip("No custom role ID available")
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/roles/enable/%s?application_key=%s", s.apiBaseURL, s.customRoleID, s.appKey)
	resp := s.makeRequest("PUT", url, nil, headers)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should enable role")

	result := s.parseJSON(resp)
	fmt.Printf("   ✅ Role enabled: %v\n", result["message"])
}

// SCENARIO 6: User Role Assignment
func (s *IntegrationTestSuite) Test22_API_RegisterSecondUser() {
	fmt.Println("\n👤 Test 23: Register second user for role assignment test")

	timestamp := time.Now().Unix()
	s.secondUserEmail = fmt.Sprintf("testuser2_%d@example.com", timestamp)
	s.secondUserPass = "TestPassword123!"

	userData := map[string]interface{}{
		"email":    s.secondUserEmail,
		"password": s.secondUserPass,
		"name":     "Test User 2",
	}

	url := fmt.Sprintf("%s/auth?application_key=%s", s.apiBaseURL, s.appKey)
	resp := s.makeRequest("POST", url, userData, nil)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should register second user")

	result := s.parseJSON(resp)

	// Store second user details
	if data, ok := result["data"].(map[string]interface{}); ok {
		if userID, exists := data["user_id"]; exists && userID != nil {
			s.secondUserID = userID.(string)
			fmt.Printf("   ✅ Second user registered: %s (ID: %s)\n", s.secondUserEmail, s.secondUserID)
		}
	}

	assert.NotEmpty(s.T(), s.secondUserID, "Second user ID should not be empty")
}

func (s *IntegrationTestSuite) Test23_API_AssignRoleToUser() {
	fmt.Println("\n🎭 Test 24: Assign custom role to second user")

	if s.customRoleID == "" {
		s.T().Skip("No custom role ID available")
	}

	if s.secondUserID == "" {
		s.T().Skip("No second user ID available")
	}

	roleData := map[string]interface{}{
		"role_id": s.customRoleID,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/users/%s/roles?application_key=%s", s.apiBaseURL, s.secondUserID, s.appKey)
	resp := s.makeRequest("PUT", url, roleData, headers)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should assign role to user")

	result := s.parseJSON(resp)
	fmt.Printf("   ✅ Custom role assigned to second user: %v\n", result["message"])
}

// SCENARIO 7: Password Management
func (s *IntegrationTestSuite) Test24_API_RequestPasswordReset() {
	fmt.Println("\n🔑 Test 25: Request password reset")

	resetData := map[string]interface{}{
		"email": s.userEmail,
	}

	url := fmt.Sprintf("%s/users/resetpassword?application_key=%s", s.apiBaseURL, s.appKey)
	resp := s.makeRequest("POST", url, resetData, nil)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should request password reset")

	result := s.parseJSON(resp)
	fmt.Printf("   ✅ Password reset requested: %v\n", result["message"])
	fmt.Println("   ℹ️  In production, check email for reset token")
}

// SCENARIO 8: Cleanup Tests
func (s *IntegrationTestSuite) Test25_API_DeleteCustomRole() {
	fmt.Println("\n🗑️  Test 26: Delete custom role")

	if s.customRoleID == "" {
		s.T().Skip("No custom role ID available")
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/roles/%s?application_key=%s", s.apiBaseURL, s.customRoleID, s.appKey)
	resp := s.makeRequest("DELETE", url, nil, headers)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should delete custom role")

	result := s.parseJSON(resp)
	fmt.Printf("   ✅ Custom role deleted: %v\n", result["message"])
}

func (s *IntegrationTestSuite) Test26_API_LogoutUser() {
	fmt.Println("\n👋 Test 27: User logout")

	headers := map[string]string{
		"Authorization": "Bearer " + s.userJWTToken,
	}

	url := fmt.Sprintf("%s/auth/logout?application_key=%s", s.apiBaseURL, s.appKey)
	resp := s.makeRequest("PUT", url, nil, headers)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Should logout successfully")

	result := s.parseJSON(resp)
	fmt.Printf("   ✅ User logged out: %v\n", result["message"])
	fmt.Println("\n🎉 All integration tests completed!")
}
