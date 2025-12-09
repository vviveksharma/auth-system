package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// Model definitions
type RoleInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	RoleType    string `json:"role_type"`
	Priority    int    `json:"priority"`
	IsSystem    bool   `json:"is_system"`
}

type Permission struct {
	Route       string   `json:"route"`
	Methods     []string `json:"methods"`
	Description string   `json:"description"`
}

type SimplePermission struct {
	Route       string `json:"route"`
	Description string `json:"description"`
}

type RoleData struct {
	RoleInfo    RoleInfo     `json:"role_info"`
	Permissions []Permission `json:"permissions"`
}


func ConvertDBData(db string) (*RoleData, error) {
	var roleData RoleData
	err := json.Unmarshal([]byte(db), &roleData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}
	return &roleData, nil
}

func IsSystemRole(roleID uuid.UUID) bool {
	roleIDStr := roleID.String()
	systemRoleMap := map[string]bool{
		"f47ac10b-58cc-4372-a567-0e02b2c3d479": true,
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8": true,
		"1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed": true,
		"550e8400-e29b-41d4-a716-446655440000": true,
	}
	
	return systemRoleMap[roleIDStr]
}

// Classify permissions by HTTP method and return as JSON-ready map
func ClassifyPermissionsByMethod(permissions []Permission) map[string][]SimplePermission {
	methodMap := make(map[string][]SimplePermission)

	for _, perm := range permissions {
		for _, method := range perm.Methods {
			simplePermission := SimplePermission{
				Route:       perm.Route,
				Description: perm.Description,
			}
			methodMap[method] = append(methodMap[method], simplePermission)
		}
	}

	return methodMap
}

func CreateFastLookup(permissions []Permission) map[string]map[string]bool {
	fastLookup := make(map[string]map[string]bool)

	for _, perm := range permissions {
		for _, method := range perm.Methods {
			if fastLookup[method] == nil {
				fastLookup[method] = make(map[string]bool)
			}
			fastLookup[method][perm.Route] = true
		}
	}

	return fastLookup
}

func FindMethodRoute(method string, route string, fastLookup map[string]map[string]bool) bool {
	if methodRoutes, exists := fastLookup[method]; exists {
		return methodRoutes[route]
	}
	return false
}

func routeToRegex(routePattern string) *regexp.Regexp {
	// Start with the original pattern
	regexPattern := routePattern

	// Replace parameter patterns with regex equivalents
	// :id becomes [^/]+ (matches any non-slash character)
	paramRegex := regexp.MustCompile(`:([a-zA-Z_][a-zA-Z0-9_]*)`)
	regexPattern = paramRegex.ReplaceAllString(regexPattern, `[^/]+`)

	// Escape special regex characters after parameter replacement
	regexPattern = regexp.QuoteMeta(regexPattern)

	// Restore the parameter regex patterns that got escaped
	regexPattern = strings.ReplaceAll(regexPattern, `\[`, `[`)
	regexPattern = strings.ReplaceAll(regexPattern, `\]`, `]`)
	regexPattern = strings.ReplaceAll(regexPattern, `\^`, `^`)
	regexPattern = strings.ReplaceAll(regexPattern, `\+`, `+`)

	// Anchor the pattern to match the complete string
	regexPattern = "^" + regexPattern + "$"

	return regexp.MustCompile(regexPattern)
}

// Check if an actual route matches a route pattern
func matchesPattern(actualRoute, routePattern string) bool {
	// If no parameters, do exact match
	if !strings.Contains(routePattern, ":") {
		return actualRoute == routePattern
	}

	// Use regex for parameterized routes
	regex := routeToRegex(routePattern)
	return regex.MatchString(actualRoute)
}

func FindMethodWithPatterns(method string, actualRoute string, patternLookup map[string][]SimplePermission) (*SimplePermission, bool) {
	methodRoutes, exists := patternLookup[method]
	if !exists {
		return nil, false
	}

	// Try to find a matching route pattern
	for _, permission := range methodRoutes {
		if matchesPattern(actualRoute, permission.Route) {
			return &permission, true
		}
	}

	return nil, false
}
