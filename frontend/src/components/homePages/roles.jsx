import React, { useState, useEffect, useRef } from "react";
import { IoAddCircleSharp, IoTrashOutline } from "react-icons/io5";
import {
  FiMoreVertical,
  FiEdit2,
  FiXCircle,
  FiCheckCircle,
} from "react-icons/fi";
import { IoMdClose } from "react-icons/io";
import { ListRoles, GetRolePermissions, AddRoles, EnableRole, DisableRole, DeleteRole } from "../services/roles";

const Roles = () => {
  const [roles, setRoles] = useState([]);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [statusFilter, setStatusFilter] = useState("All");
  const [roleTypeFilter, setRoleTypeFilter] = useState("All");
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [showDropdown, setShowDropdown] = useState(null);
  const [showAddModal, setShowAddModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showPermissionsModal, setShowPermissionsModal] = useState(false);
  const [selectedPermissions, setSelectedPermissions] = useState([]);
  const [editingRole, setEditingRole] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [totalRoles, setTotalRoles] = useState(0);
  const [loadingPermissions, setLoadingPermissions] = useState(false);
  const [permissionsError, setPermissionsError] = useState(null);
  const [selectedRole, setSelectedRole] = useState(null);
  const [addingRole, setAddingRole] = useState(false);
  const [addRoleError, setAddRoleError] = useState(null);
  const [editablePermissions, setEditablePermissions] = useState(new Set());
  const [editingRoleDetails, setEditingRoleDetails] = useState(false);
  const [updatingRole, setUpdatingRole] = useState(false);
  const [updateRoleError, setUpdateRoleError] = useState(null);
  const [newRole, setNewRole] = useState({
    name: "",
    displayName: "",
    description: "",
    permissions: [{
      route: "",
      methods: [],
      description: ""
    }],
  });
  const dropdownRef = useRef(null);
  
  const httpMethods = ["GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"];

  // Fetch roles from API
  useEffect(() => {
    fetchRoles();
  }, [currentPage, rowsPerPage, statusFilter, roleTypeFilter]);

  const fetchRoles = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Convert filters to API format (all lowercase)
      // Status: "All" -> "", "Enabled" -> "active", "Disabled" -> "inactive"
      let apiStatus = "";
      if (statusFilter === "Enabled") {
        apiStatus = "active";
      } else if (statusFilter === "Disabled") {
        apiStatus = "inactive";
      } else {
        apiStatus = ""; // "All" selected - empty parameter
      }
      
      // Role Type: "All" -> "", "Default" -> "system", "Custom" -> "custom"
      let apiRoleType = "";
      if (roleTypeFilter === "Default") {
        apiRoleType = "system";
      } else if (roleTypeFilter === "Custom") {
        apiRoleType = "custom";
      } else {
        apiRoleType = ""; // "All" selected - empty parameter
      }
      
      const response = await ListRoles(apiStatus, currentPage, rowsPerPage, apiRoleType);
      
      // Extract data from nested response structure
      const roleData = response.data?.data || [];
      const paginationData = response.data?.pagination || {};
      
      // Transform API response to match the component's data structure
      const transformedRoles = roleData.map((role) => ({
        id: role.id,
        name: role.string || role.name, // API uses 'string' field for role name
        displayName: role.display_name || role.displayName,
        description: role.description || "",
        roleType: role.role_type || "default",
        status: role.status ? "Enabled" : "Disabled",
        permissions: role.permissions || [],
      }));
      
      setRoles(transformedRoles);
      setTotalRoles(paginationData.total_items || transformedRoles.length);
    } catch (err) {
      console.error("Error fetching roles:", err);
      setError("Failed to load roles. Please try again.");
      setRoles([]);
    } finally {
      setLoading(false);
    }
  };

  const toggleDropdown = (roleId) => {
    setShowDropdown(showDropdown === roleId ? null : roleId);
  };

  useEffect(() => {
    function handleClickOutside(event) {
      if (
        showDropdown !== null &&
        !event.target.closest(`[data-dropdown-id="${showDropdown}"]`)
      ) {
        setShowDropdown(null);
      }
    }
    if (showDropdown !== null) {
      document.addEventListener("mousedown", handleClickOutside);
    }
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [showDropdown]);

  // Add a new permission field
  const addPermissionField = (isEdit = false) => {
    if (isEdit && editingRole) {
      setEditingRole({
        ...editingRole,
        permissions: [...editingRole.permissions, { route: "", methods: [], description: "" }],
      });
    } else {
      setNewRole({ 
        ...newRole, 
        permissions: [...newRole.permissions, { route: "", methods: [], description: "" }] 
      });
    }
  };

  // Remove a permission field
  const removePermissionField = (index, isEdit = false) => {
    if (isEdit && editingRole) {
      const updatedPermissions = [...editingRole.permissions];
      updatedPermissions.splice(index, 1);
      setEditingRole({ ...editingRole, permissions: updatedPermissions });
    } else {
      const updatedPermissions = [...newRole.permissions];
      updatedPermissions.splice(index, 1);
      setNewRole({ ...newRole, permissions: updatedPermissions });
    }
  };

  // Handle permission input changes
  const handlePermissionChange = (index, field, value, isEdit = false) => {
    if (isEdit && editingRole) {
      const updatedPermissions = [...editingRole.permissions];
      updatedPermissions[index] = {
        ...updatedPermissions[index],
        [field]: value
      };
      setEditingRole({ ...editingRole, permissions: updatedPermissions });
    } else {
      const updatedPermissions = [...newRole.permissions];
      updatedPermissions[index] = {
        ...updatedPermissions[index],
        [field]: value
      };
      setNewRole({ ...newRole, permissions: updatedPermissions });
    }
  };

  // Toggle HTTP method selection
  const toggleMethod = (permIndex, method, isEdit = false) => {
    if (isEdit && editingRole) {
      const updatedPermissions = [...editingRole.permissions];
      const currentMethods = updatedPermissions[permIndex].methods;
      
      if (currentMethods.includes(method)) {
        updatedPermissions[permIndex].methods = currentMethods.filter(m => m !== method);
      } else {
        updatedPermissions[permIndex].methods = [...currentMethods, method];
      }
      
      setEditingRole({ ...editingRole, permissions: updatedPermissions });
    } else {
      const updatedPermissions = [...newRole.permissions];
      const currentMethods = updatedPermissions[permIndex].methods;
      
      if (currentMethods.includes(method)) {
        updatedPermissions[permIndex].methods = currentMethods.filter(m => m !== method);
      } else {
        updatedPermissions[permIndex].methods = [...currentMethods, method];
      }
      
      setNewRole({ ...newRole, permissions: updatedPermissions });
    }
  };

  // Submit new role
  const handleAddRole = async () => {
    if (!newRole.name.trim() || !newRole.displayName.trim()) return;

    const filteredPermissions = newRole.permissions.filter(
      perm => perm.route.trim() && perm.methods.length > 0
    );
    
    if (filteredPermissions.length === 0) return;

    try {
      setAddingRole(true);
      setAddRoleError(null);

      // Call API to add role
      await AddRoles(
        newRole.name,
        newRole.displayName,
        newRole.description,
        filteredPermissions
      );

      // Reset form
      setNewRole({
        name: "",
        displayName: "",
        description: "",
        permissions: [{
          route: "",
          methods: [],
          description: ""
        }],
      });
      
      // Close modal
      setShowAddModal(false);
      
      // Refetch roles to update the list
      await fetchRoles();
      
    } catch (err) {
      console.error("Error adding role:", err);
      setAddRoleError("Failed to add role. Please try again.");
    } finally {
      setAddingRole(false);
    }
  };

  const isSystemRole = (role) => {
    // System roles have roleType "default"
    return role.roleType === "default";
  };

  // Open edit modal - fetch permissions first
  const openEditModal = async (role) => {
    try {
      setShowDropdown(null);
      setEditingRole({ ...role });
      setShowEditModal(true);
      setEditablePermissions(new Set());
      setEditingRoleDetails(false);
      setUpdateRoleError(null);
      
      // Fetch fresh permissions data
      setLoadingPermissions(true);
      const response = await GetRolePermissions(role.id);
      const permissions = response.data?.permissions || [];
      
      setEditingRole({ 
        ...role, 
        permissions: permissions.length > 0 ? permissions : role.permissions 
      });
    } catch (err) {
      console.error("Error fetching permissions for edit:", err);
      setUpdateRoleError("Failed to load permissions. Please try again.");
    } finally {
      setLoadingPermissions(false);
    }
  };

  // Toggle permission editable state
  const togglePermissionEditable = (index) => {
    const newSet = new Set(editablePermissions);
    if (newSet.has(index)) {
      newSet.delete(index);
    } else {
      newSet.add(index);
    }
    setEditablePermissions(newSet);
  };

  // Submit edited role - smart update
  const handleEditRole = async () => {
    if (!editingRole.name.trim() || !editingRole.displayName.trim()) return;

    const filteredPermissions = editingRole.permissions.filter(
      perm => perm.route.trim() && perm.methods.length > 0
    );
    
    if (filteredPermissions.length === 0) return;

    try {
      setUpdatingRole(true);
      setUpdateRoleError(null);

      // Check what changed
      const originalRole = roles.find(r => r.id === editingRole.id);
      const roleDetailsChanged = 
        originalRole.name !== editingRole.name ||
        originalRole.displayName !== editingRole.displayName ||
        originalRole.description !== editingRole.description;

      const permissionsChanged = editablePermissions.size > 0;

      if (roleDetailsChanged || permissionsChanged) {
        // TODO: Call UpdateRole API with only changed data
        // For now, update optimistically
        setRoles(
          roles.map((role) =>
            role.id === editingRole.id
              ? { ...editingRole, permissions: filteredPermissions }
              : role
          )
        );
        
        console.log("Changes to update:", {
          roleDetailsChanged,
          permissionsChanged,
          editedPermissionIndexes: Array.from(editablePermissions),
          data: {
            name: editingRole.name,
            displayName: editingRole.displayName,
            description: editingRole.description,
            permissions: permissionsChanged ? filteredPermissions : null
          }
        });
      }

      setShowEditModal(false);
      setEditingRole(null);
      setEditablePermissions(new Set());
      setEditingRoleDetails(false);
    } catch (err) {
      console.error("Error updating role:", err);
      setUpdateRoleError("Failed to update role. Please try again.");
    } finally {
      setUpdatingRole(false);
    }
  };

  // View permissions for a role
  const viewRolePermissions = async (role) => {
    try {
      console.log('Role Details:', role);
      setSelectedRole(role);
      setShowPermissionsModal(true);
      setLoadingPermissions(true);
      setPermissionsError(null);
      
      // Fetch permissions from API
      const response = await GetRolePermissions(role.id);
      
      // Extract permissions from response
      const permissions = response.data?.permissions || [];
      
      setSelectedPermissions(permissions);
    } catch (err) {
      console.error("Error fetching permissions:", err);
      setPermissionsError("Failed to load permissions. Please try again.");
      setSelectedPermissions([]);
    } finally {
      setLoadingPermissions(false);
    }
  };

  const deleteRole = async (roleId) => {
    try {
      setShowDropdown(null);
      
      // Optimistically remove from UI immediately
      setRoles(roles.filter((role) => role.id !== roleId));
      
      // Call API to delete role
      await DeleteRole(roleId);
      
    } catch (err) {
      console.error("Error deleting role:", err);
      // Refetch on error to restore the list
      await fetchRoles();
      alert("Failed to delete role. Please try again.");
    }
  };

  const toggleRoleStatus = async (roleId) => {
    try {
      setShowDropdown(null);
      const role = roles.find(r => r.id === roleId);
      if (!role) return;
      const updatedStatus = role.status === "Enabled" ? "Disabled" : "Enabled";
      setRoles(
        roles.map((r) =>
          r.id === roleId
            ? { ...r, status: updatedStatus }
            : r
        )
      );
      if (role.status === "Enabled") {
        await DisableRole(roleId);
      } else {
        await EnableRole(roleId);
      }
      
    } catch (err) {
      console.error("Error toggling role status:", err);
      setRoles(
        roles.map((r) =>
          r.id === roleId
            ? { ...r, status: role.status }
            : r
        )
      );
      alert("Failed to update role status. Please try again.");
    }
  };

  // Client-side search filtering (searches within current page results)
  const filteredRoles = roles.filter((role) => {
    if (searchQuery === "") return true;
    const matchesSearch = role.name
      .toLowerCase()
      .includes(searchQuery.toLowerCase()) || 
      role.displayName.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesSearch;
  });

  const totalPages = Math.ceil(totalRoles / rowsPerPage);

  return (
    <div className="antialiased font-sans bg-gray-50 min-h-screen">
      {/* Add Role Modal */}
      {showAddModal && (
        <div className="fixed inset-0 bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
            <div className="flex justify-between items-center border-b border-gray-200 px-6 py-4 sticky top-0 bg-white z-10">
              <h3 className="text-lg font-medium text-gray-900">
                Add New Role
              </h3>
              <button
                onClick={() => {
                  setShowAddModal(false);
                  setAddRoleError(null);
                }}
                className="text-gray-400 hover:text-gray-500 cursor-pointer"
                disabled={addingRole}
              >
                <IoMdClose className="h-6 w-6 cursor-pointer" />
              </button>
            </div>
            <div className="px-6 py-4 space-y-4">
              {addRoleError && (
                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
                  <p className="text-sm">{addRoleError}</p>
                </div>
              )}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Role Name *
                  </label>
                  <input
                    type="text"
                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500"
                    value={newRole.name}
                    onChange={(e) =>
                      setNewRole({ ...newRole, name: e.target.value })
                    }
                    placeholder="Enter role name (e.g., editor)"
                    disabled={addingRole}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Display Name *
                  </label>
                  <input
                    type="text"
                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500"
                    value={newRole.displayName}
                    onChange={(e) =>
                      setNewRole({ ...newRole, displayName: e.target.value })
                    }
                    placeholder="Enter display name (e.g., Content Editor)"
                    disabled={addingRole}
                  />
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <textarea
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500"
                  value={newRole.description}
                  onChange={(e) =>
                    setNewRole({ ...newRole, description: e.target.value })
                  }
                  placeholder="Enter role description"
                  rows="2"
                  disabled={addingRole}
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Permissions *
                </label>
                <div className="space-y-4">
                  {newRole.permissions.map((permission, index) => (
                    <div key={index} className="border border-gray-200 rounded-md p-4 bg-gray-50">
                      <div className="flex justify-between items-center mb-3">
                        <span className="text-sm font-medium text-gray-700">Permission #{index + 1}</span>
                        {newRole.permissions.length > 1 && (
                          <button
                            type="button"
                            onClick={() => removePermissionField(index)}
                            className="text-red-500 hover:text-red-700 text-sm flex items-center cursor-pointer"
                          >
                            <IoMdClose className="h-4 w-4 mr-1" />
                            Remove
                          </button>
                        )}
                      </div>
                      
                      <div className="grid grid-cols-1 gap-3">
                        <div>
                          <label className="block text-xs text-gray-600 mb-1">Route Path *</label>
                          <input
                            type="text"
                            className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500 text-sm"
                            value={permission.route}
                            onChange={(e) => handlePermissionChange(index, 'route', e.target.value)}
                            placeholder="Enter route path (e.g., /api/users)"
                          />
                        </div>
                        
                        <div>
                          <label className="block text-xs text-gray-600 mb-1">HTTP Methods *</label>
                          <div className="flex flex-wrap gap-2">
                            {httpMethods.map((method) => (
                              <button
                                key={method}
                                type="button"
                                onClick={() => toggleMethod(index, method)}
                                className={`px-2 py-1 text-xs rounded-md border cursor-pointer ${
                                  permission.methods.includes(method)
                                    ? 'bg-stone-800 text-white border-stone-800'
                                    : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-100'
                                }`}
                              >
                                {method}
                              </button>
                            ))}
                          </div>
                        </div>
                        
                        <div>
                          <label className="block text-xs text-gray-600 mb-1">Description</label>
                          <input
                            type="text"
                            className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500 text-sm"
                            value={permission.description}
                            onChange={(e) => handlePermissionChange(index, 'description', e.target.value)}
                            placeholder="Describe what this permission allows"
                          />
                        </div>
                      </div>
                    </div>
                  ))}
                  
                  <button
                    type="button"
                    onClick={() => addPermissionField(false)}
                    className="mt-2 text-sm text-stone-600 hover:text-stone-800 flex items-center cursor-pointer"
                  >
                    <IoAddCircleSharp className="mr-1" /> Add another permission
                  </button>
                </div>
              </div>
            </div>
            <div className="bg-gray-50 px-6 py-3 flex justify-end border-t border-gray-200 sticky bottom-0">
              <button
                type="button"
                onClick={() => {
                  setShowAddModal(false);
                  setAddRoleError(null);
                }}
                disabled={addingRole}
                className="mr-3 px-4 py-2 text-sm text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={handleAddRole}
                disabled={
                  addingRole ||
                  !newRole.name.trim() || 
                  !newRole.displayName.trim() ||
                  newRole.permissions.every(p => !p.route.trim() || p.methods.length === 0)
                }
                className="px-4 py-2 text-sm text-white bg-stone-800 border border-transparent rounded-md hover:bg-stone-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
              >
                {addingRole ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                    Adding...
                  </>
                ) : (
                  'Add Role'
                )}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Edit Role Modal */}
      {showEditModal && editingRole && (
        <div className="fixed inset-0 bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
            <div className="flex justify-between items-center border-b border-gray-200 px-6 py-4 sticky top-0 bg-white z-10">
              <h3 className="text-lg font-medium text-gray-900">Edit Role</h3>
              <button
                onClick={() => {
                  setShowEditModal(false);
                  setEditingRole(null);
                  setEditablePermissions(new Set());
                  setEditingRoleDetails(false);
                  setUpdateRoleError(null);
                }}
                disabled={updatingRole}
                className="text-gray-400 hover:text-gray-500 cursor-pointer disabled:opacity-50"
              >
                <IoMdClose className="h-6 w-6 cursor-pointer" />
              </button>
            </div>
            <div className="px-6 py-4 space-y-4">
              {updateRoleError && (
                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
                  <p className="text-sm">{updateRoleError}</p>
                </div>
              )}
              
              {loadingPermissions ? (
                <div className="flex justify-center items-center py-10">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
                  <span className="ml-3">Loading role details...</span>
                </div>
              ) : (
                <>
                  {/* Role Details Section with Toggle */}
                  <div className="border border-gray-200 rounded-md p-4 bg-gray-50">
                    <div className="flex items-center justify-between mb-3">
                      <h4 className="text-sm font-medium text-gray-700">Role Details</h4>
                      <label className="flex items-center cursor-pointer">
                        <input
                          type="checkbox"
                          checked={editingRoleDetails}
                          onChange={(e) => setEditingRoleDetails(e.target.checked)}
                          className="mr-2 h-4 w-4 text-stone-600 focus:ring-stone-500 border-gray-300 rounded"
                        />
                        <span className="text-xs text-gray-600">Enable editing</span>
                      </label>
                    </div>
                    
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                          Role Name *
                        </label>
                        <input
                          type="text"
                          className={`w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500 ${
                            !editingRoleDetails ? 'bg-gray-100 cursor-not-allowed' : ''
                          }`}
                          value={editingRole.name}
                          onChange={(e) =>
                            setEditingRole({ ...editingRole, name: e.target.value })
                          }
                          placeholder="Enter role name"
                          disabled={!editingRoleDetails || updatingRole}
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                          Display Name *
                        </label>
                        <input
                          type="text"
                          className={`w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500 ${
                            !editingRoleDetails ? 'bg-gray-100 cursor-not-allowed' : ''
                          }`}
                          value={editingRole.displayName}
                          onChange={(e) =>
                            setEditingRole({ ...editingRole, displayName: e.target.value })
                          }
                          placeholder="Enter display name"
                          disabled={!editingRoleDetails || updatingRole}
                        />
                      </div>
                    </div>
                    
                    <div className="mt-4">
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Description
                      </label>
                      <textarea
                        className={`w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500 ${
                          !editingRoleDetails ? 'bg-gray-100 cursor-not-allowed' : ''
                        }`}
                        value={editingRole.description}
                        onChange={(e) =>
                          setEditingRole({ ...editingRole, description: e.target.value })
                        }
                        placeholder="Enter role description"
                        rows="2"
                        disabled={!editingRoleDetails || updatingRole}
                      />
                    </div>
                  </div>
                  
                  {/* Permissions Section */}
                  <div>
                    <div className="flex items-center justify-between mb-2">
                      <label className="block text-sm font-medium text-gray-700">
                        Permissions *
                      </label>
                      <span className="text-xs text-gray-500">
                        Check the box to edit a permission
                      </span>
                    </div>
                    <div className="space-y-4">
                      {editingRole.permissions && editingRole.permissions.length > 0 ? (
                        editingRole.permissions.map((permission, index) => {
                          const isEditable = editablePermissions.has(index);
                          return (
                            <div 
                              key={index} 
                              className={`border rounded-md p-4 ${
                                isEditable ? 'border-stone-500 bg-stone-50' : 'border-gray-200 bg-gray-50'
                              }`}
                            >
                              <div className="flex justify-between items-center mb-3">
                                <span className="text-sm font-medium text-gray-700">
                                  Permission #{index + 1}
                                </span>
                                <div className="flex items-center gap-2">
                                  <label className="flex items-center cursor-pointer">
                                    <input
                                      type="checkbox"
                                      checked={isEditable}
                                      onChange={() => togglePermissionEditable(index)}
                                      disabled={updatingRole}
                                      className="mr-2 h-4 w-4 text-stone-600 focus:ring-stone-500 border-gray-300 rounded"
                                    />
                                    <span className="text-xs text-gray-600">
                                      {isEditable ? 'Editing' : 'Edit'}
                                    </span>
                                  </label>
                                  {editingRole.permissions.length > 1 && isEditable && (
                                    <button
                                      type="button"
                                      onClick={() => removePermissionField(index, true)}
                                      disabled={updatingRole}
                                      className="text-red-500 hover:text-red-700 text-sm flex items-center cursor-pointer disabled:opacity-50"
                                    >
                                      <IoMdClose className="h-4 w-4 mr-1" />
                                      Remove
                                    </button>
                                  )}
                                </div>
                              </div>
                              
                              <div className="grid grid-cols-1 gap-3">
                                <div>
                                  <label className="block text-xs text-gray-600 mb-1">Route Path *</label>
                                  <input
                                    type="text"
                                    className={`w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500 text-sm ${
                                      !isEditable ? 'bg-gray-100 cursor-not-allowed' : ''
                                    }`}
                                    value={permission.route}
                                    onChange={(e) => handlePermissionChange(index, 'route', e.target.value, true)}
                                    placeholder="Enter route path"
                                    disabled={!isEditable || updatingRole}
                                  />
                                </div>
                                
                                <div>
                                  <label className="block text-xs text-gray-600 mb-1">HTTP Methods *</label>
                                  <div className="flex flex-wrap gap-2">
                                    {httpMethods.map((method) => (
                                      <button
                                        key={method}
                                        type="button"
                                        onClick={() => toggleMethod(index, method, true)}
                                        disabled={!isEditable || updatingRole}
                                        className={`px-2 py-1 text-xs rounded-md border cursor-pointer transition-colors ${
                                          permission.methods.includes(method)
                                            ? 'bg-stone-800 text-white border-stone-800'
                                            : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-100'
                                        } ${
                                          !isEditable || updatingRole ? 'opacity-50 cursor-not-allowed' : ''
                                        }`}
                                      >
                                        {method}
                                      </button>
                                    ))}
                                  </div>
                                </div>
                                
                                <div>
                                  <label className="block text-xs text-gray-600 mb-1">Description</label>
                                  <input
                                    type="text"
                                    className={`w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500 text-sm ${
                                      !isEditable ? 'bg-gray-100 cursor-not-allowed' : ''
                                    }`}
                                    value={permission.description}
                                    onChange={(e) => handlePermissionChange(index, 'description', e.target.value, true)}
                                    placeholder="Describe what this permission allows"
                                    disabled={!isEditable || updatingRole}
                                  />
                                </div>
                              </div>
                            </div>
                          );
                        })
                      ) : (
                        <div className="text-center py-6 text-gray-500 border border-gray-200 rounded-md bg-gray-50">
                          No permissions found. Add a new permission below.
                        </div>
                      )}
                      
                      <button
                        type="button"
                        onClick={() => addPermissionField(true)}
                        disabled={updatingRole}
                        className="mt-2 text-sm text-stone-600 hover:text-stone-800 flex items-center cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <IoAddCircleSharp className="mr-1" /> Add another permission
                      </button>
                    </div>
                  </div>
                </>
              )}
            </div>
            <div className="bg-gray-50 px-6 py-3 flex justify-end border-t border-gray-200 sticky bottom-0">
              <button
                type="button"
                onClick={() => {
                  setShowEditModal(false);
                  setEditingRole(null);
                  setEditablePermissions(new Set());
                  setEditingRoleDetails(false);
                  setUpdateRoleError(null);
                }}
                disabled={updatingRole}
                className="mr-3 px-4 py-2 text-sm text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={handleEditRole}
                disabled={
                  updatingRole ||
                  loadingPermissions ||
                  (!editingRoleDetails && editablePermissions.size === 0) ||
                  !editingRole.name.trim() || 
                  !editingRole.displayName.trim() ||
                  editingRole.permissions.every(p => !p.route.trim() || p.methods.length === 0)
                }
                className="px-4 py-2 text-sm text-white bg-stone-800 border border-transparent rounded-md hover:bg-stone-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
              >
                {updatingRole ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                    Updating...
                  </>
                ) : (
                  'Save Changes'
                )}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* View Permissions Modal */}
      {showPermissionsModal && (
        <div className="fixed inset-0 bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
            <div className="flex justify-between items-center border-b border-gray-200 px-6 py-4 sticky top-0 bg-white">
              <div>
                <h3 className="text-lg font-medium text-gray-900">Role Permissions</h3>
                {selectedRole && (
                  <p className="text-sm text-gray-600 mt-1">
                    <span className="font-medium">{selectedRole.displayName}</span> 
                    <span className="text-gray-400 mx-2">•</span>
                    <span className="text-gray-500">{selectedRole.name}</span>
                  </p>
                )}
              </div>
              <button
                onClick={() => {
                  setShowPermissionsModal(false);
                  setSelectedRole(null);
                  setSelectedPermissions([]);
                  setPermissionsError(null);
                }}
                className="text-gray-400 hover:text-gray-500 cursor-pointer"
              >
                <IoMdClose className="h-6 w-6" />
              </button>
            </div>
            <div className="px-6 py-4">
              {loadingPermissions ? (
                <div className="flex justify-center items-center py-10">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
                  <span className="ml-3">Loading permissions...</span>
                </div>
              ) : permissionsError ? (
                <div className="text-red-600 text-center py-10">
                  <p>{permissionsError}</p>
                  <button
                    onClick={() => selectedRole && viewRolePermissions(selectedRole)}
                    className="mt-2 text-sm text-blue-600 hover:text-blue-800 underline"
                  >
                    Retry
                  </button>
                </div>
              ) : selectedPermissions.length > 0 ? (
                <div className="space-y-4">
                  {selectedPermissions.map((permission, index) => (
                    <div key={index} className="border border-gray-200 rounded-md p-4">
                      <div className="mb-2">
                        <span className="text-sm font-medium text-gray-700">Route:</span>
                        <span className="ml-2 bg-gray-100 rounded-md px-2 py-1 text-sm font-mono">
                          {permission.route}
                        </span>
                      </div>
                      
                      <div className="mb-2">
                        <span className="text-sm font-medium text-gray-700">Methods:</span>
                        <div className="flex flex-wrap gap-1 mt-1">
                          {permission.methods.map((method) => (
                            <span 
                              key={method} 
                              className="px-2 py-1 text-xs rounded-md bg-stone-800 text-white"
                            >
                              {method}
                            </span>
                          ))}
                        </div>
                      </div>
                      
                      {permission.description && (
                        <div>
                          <span className="text-sm font-medium text-gray-700">Description:</span>
                          <p className="mt-1 text-sm text-gray-600">{permission.description}</p>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-10 text-gray-500">
                  No permissions found for this role
                </div>
              )}
            </div>
            <div className="bg-gray-50 px-6 py-3 flex justify-end border-t border-gray-200 sticky bottom-0">
              <button
                type="button"
                onClick={() => setShowPermissionsModal(false)}
                className="px-4 py-2 text-sm text-white bg-stone-800 border border-transparent rounded-md hover:bg-stone-700 cursor-pointer"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Main Content */}
      <div className="container mx-auto px-2 sm:px-8">
        <div className="py-4">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-2xl font-semibold leading-tight">Roles</h2>
            <button
              onClick={() => setShowAddModal(true)}
              className="flex items-center gap-2 justify-center border align-middle select-none font-sans font-medium text-center duration-300 ease-in disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed focus:shadow-none text-sm py-2 px-4 shadow-sm hover:shadow-md bg-stone-800 hover:bg-stone-700 border-stone-900 text-stone-50 rounded-lg transition antialiased cursor-pointer"
            >
              <IoAddCircleSharp className="text-lg" />
              Add Role
            </button>
          </div>

          {/* Filters and Search */}
          <div className="my-2 flex flex-col sm:flex-row gap-2 items-stretch">
            <div className="flex flex-col sm:flex-row mb-1 sm:mb-0 gap-2 flex-1">
              <div className="relative">
                <select
                  className="appearance-none h-full rounded border block w-full bg-white border-gray-400 text-gray-700 py-2 px-4 pr-8 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                  value={rowsPerPage}
                  onChange={(e) => {
                    setRowsPerPage(Number(e.target.value));
                    setCurrentPage(1); // Reset to first page on page size change
                  }}
                >
                  <option value={5}>5 per page</option>
                  <option value={10}>10 per page</option>
                  <option value={20}>20 per page</option>
                </select>
                <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700">
                  <svg
                    className="fill-current h-4 w-4"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 20 20"
                  >
                    <path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z" />
                  </svg>
                </div>
              </div>
              <div className="relative">
                <select
                  className="appearance-none h-full rounded-md sm:rounded-none sm:border border block w-full bg-white border-gray-400 text-gray-700 py-2 px-4 pr-8 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                  value={statusFilter}
                  onChange={(e) => {
                    setStatusFilter(e.target.value);
                    setCurrentPage(1); // Reset to first page on filter change
                  }}
                >
                  <option value="All">All Status</option>
                  <option value="Enabled">Enabled</option>
                  <option value="Disabled">Disabled</option>
                </select>
                <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700">
                  <svg
                    className="fill-current h-4 w-4"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 20 20"
                  >
                    <path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z" />
                  </svg>
                </div>
              </div>
              <div className="relative">
                <select
                  className="appearance-none h-full rounded-md sm:rounded-none sm:border border block w-full bg-white border-gray-400 text-gray-700 py-2 px-4 pr-8 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                  value={roleTypeFilter}
                  onChange={(e) => {
                    setRoleTypeFilter(e.target.value);
                    setCurrentPage(1); // Reset to first page on filter change
                  }}
                >
                  <option value="All">All Types</option>
                  <option value="Default">System Roles</option>
                  <option value="Custom">Custom Roles</option>
                </select>
                <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700">
                  <svg
                    className="fill-current h-4 w-4"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 20 20"
                  >
                    <path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z" />
                  </svg>
                </div>
              </div>
            </div>
            <div className="block relative w-full sm:w-auto flex-1">
              <span className="h-full absolute inset-y-0 left-0 flex items-center pl-2">
                <svg
                  viewBox="0 0 24 24"
                  className="h-4 w-4 fill-current text-gray-500"
                >
                  <path d="M10 4a6 6 0 100 12 6 6 0 000-12zm-8 6a8 8 0 1114.32 4.906l5.387 5.387a1 1 0 01-1.414 1.414l-5.387-5.387A8 8 0 012 10z"></path>
                </svg>
              </span>
              <input
                placeholder="Search roles..."
                className="appearance-none rounded-r rounded-l sm:rounded-l-none border border-gray-400 border-b block pl-8 pr-6 py-2 w-full bg-white text-sm placeholder-gray-400 text-gray-700 focus:bg-white focus:placeholder-gray-600 focus:text-gray-700 focus:outline-none"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
          </div>

          {/* Roles Table */}
          <div className="-mx-2 sm:-mx-8 px-2 sm:px-8 py-4 overflow-x-auto">
            <div className="inline-block min-w-full shadow rounded-lg overflow-hidden">
              <table className="min-w-full leading-normal">
                <thead>
                  <tr>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Name
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Display Name
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Role Type
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Permissions
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {loading ? (
                    <tr>
                      <td
                        colSpan="6"
                        className="px-5 py-10 border-b border-gray-200 bg-white text-sm text-center"
                      >
                        <div className="flex justify-center items-center">
                          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
                          <span className="ml-3">Loading roles...</span>
                        </div>
                      </td>
                    </tr>
                  ) : error ? (
                    <tr>
                      <td
                        colSpan="6"
                        className="px-5 py-10 border-b border-gray-200 bg-white text-sm text-center"
                      >
                        <div className="text-red-600">
                          <p>{error}</p>
                          <button
                            onClick={fetchRoles}
                            className="mt-2 text-sm text-blue-600 hover:text-blue-800 underline"
                          >
                            Retry
                          </button>
                        </div>
                      </td>
                    </tr>
                  ) : filteredRoles.length > 0 ? (
                    filteredRoles.map((role) => (
                      <tr key={role.id} className="hover:bg-gray-50">
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <div className="flex items-center">
                            <div className="ml-3">
                              <p className="text-gray-900 whitespace-no-wrap font-medium">
                                {role.name}
                              </p>
                            </div>
                          </div>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <p className="text-gray-900 whitespace-no-wrap">
                            {role.displayName}
                          </p>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <span
                            className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                              role.roleType === "default"
                                ? "bg-blue-100 text-blue-800"
                                : "bg-purple-100 text-purple-800"
                            }`}
                          >
                            {role.roleType === "default"
                              ? "System Role"
                              : "Custom Role"}
                          </span>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <span
                            className={`relative inline-block px-3 py-1 font-semibold leading-tight ${
                              role.status === "Enabled"
                                ? "text-green-900"
                                : "text-gray-900"
                            }`}
                          >
                            <span
                              aria-hidden
                              className={`absolute inset-0 rounded-full opacity-50 ${
                                role.status === "Enabled"
                                  ? "bg-green-200"
                                  : "bg-gray-300"
                              }`}
                            ></span>
                            <span className="relative">{role.status}</span>
                          </span>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <button
                            onClick={() => viewRolePermissions(role)}
                            className="inline-flex items-center justify-center cursor-pointer border align-middle select-none font-sans font-medium text-center transition-all ease-in disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed focus:shadow-none text-sm py-2 px-4 shadow-sm bg-transparent relative text-stone-700 hover:text-stone-700 border-stone-500 hover:bg-transparent duration-150 hover:border-stone-600 rounded-lg hover:opacity-60 hover:shadow-none"
                          >
                            View Permissions
                          </button>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <div className="relative" data-dropdown-id={role.id}>
                            {!isSystemRole(role) && (
                              <button
                                onClick={() => toggleDropdown(role.id)}
                                className="text-gray-500 hover:text-gray-700 focus:outline-none"
                              >
                                <FiMoreVertical className="h-5 w-5" />
                              </button>
                            )}
                            {showDropdown === role.id && (
                              <div
                                data-dropdown-id={role.id}
                                className="absolute right-0 w-48 bg-white rounded-md shadow-lg z-10 border border-gray-200"
                              >
                                <div className="py-1">
                                  <button
                                    onClick={() => openEditModal(role)}
                                    className="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 w-full text-left"
                                  >
                                    <FiEdit2 className="mr-2" />
                                    Edit
                                  </button>
                                  <button
                                    onClick={() => toggleRoleStatus(role.id)}
                                    className="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 w-full text-left"
                                  >
                                    {role.status === "Enabled" ? (
                                      <FiXCircle className="mr-2" />
                                    ) : (
                                      <FiCheckCircle className="mr-2" />
                                    )}
                                    {role.status === "Enabled"
                                      ? "Disable"
                                      : "Enable"}
                                  </button>
                                  <button
                                    onClick={() => deleteRole(role.id)}
                                    className="flex items-center px-4 py-2 text-sm text-red-600 hover:bg-gray-100 w-full text-left"
                                  >
                                    <IoTrashOutline className="mr-2" />
                                    Delete
                                  </button>
                                </div>
                              </div>
                            )}
                          </div>
                        </td>
                      </tr>
                    ))
                  ) : (
                    <tr>
                      <td
                        colSpan="6"
                        className="px-5 py-5 border-b border-gray-200 bg-white text-sm text-center"
                      >
                        No roles found matching your criteria
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>

              {/* Pagination */}
              <div className="px-5 py-5 bg-white border-t flex flex-col sm:flex-row items-center justify-between gap-2">
                <span className="text-xs xs:text-sm text-gray-900">
                  Showing {Math.min((currentPage - 1) * rowsPerPage + 1, totalRoles)} to{" "}
                  {Math.min(currentPage * rowsPerPage, totalRoles)} of{" "}
                  {totalRoles} roles
                </span>
                <div className="inline-flex mt-2 xs:mt-0">
                  <button
                    onClick={() =>
                      setCurrentPage((prev) => Math.max(prev - 1, 1))
                    }
                    disabled={currentPage === 1}
                    className={`text-sm py-2 px-4 rounded-l cursor-pointer ${
                      currentPage === 1
                        ? "bg-gray-200 text-gray-500 cursor-not-allowed"
                        : "bg-gray-300 hover:bg-gray-400 text-gray-800"
                    }`}
                  >
                    Prev
                  </button>
                  <button
                    onClick={() =>
                      setCurrentPage((prev) => Math.min(prev + 1, totalPages))
                    }
                    disabled={currentPage === totalPages || totalPages === 0}
                    className={`text-sm py-2 px-4 rounded-r cursor-pointer ${
                      currentPage === totalPages || totalPages === 0
                        ? "bg-gray-200 text-gray-500 cursor-not-allowed"
                        : "bg-gray-300 hover:bg-gray-400 text-gray-800"
                    }`}
                  >
                    Next
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Roles;