import React, { useState, useEffect, useRef } from "react";
import { IoAddCircleSharp, IoTrashOutline, IoClose } from "react-icons/io5";
import { FiMoreVertical, FiEdit2 } from "react-icons/fi";

const Temp = () => {
  const [roles, setRoles] = useState([
    {
      id: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      name: "Admin",
      roleType: "default",
      status: "Enabled",
      routes: ["/admin/logs", "/admin/views"],
    },
    {
      id: "1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed",
      name: "Moderator",
      roleType: "default",
      status: "Enabled",
      routes: ["/admin/logs", "/admin/views"],
    },
    {
      id: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      name: "User",
      roleType: "default",
      status: "Enabled",
      routes: ["/admin/logs", "/admin/views"],
    },
    {
      id: "550e8400-e29b-41d4-a716-446655440000",
      name: "Guest",
      roleType: "default",
      status: "Enabled",
      routes: ["/customer/logs", "/customer/views"],
    },
    {
      id: "7d75aee6-0bd7-4c59-a6d5-c31c51984246",
      name: "Manager",
      roleType: "custom",
      status: "Enabled",
      routes: ["/customer/logs", "/customer/views"],
    },
  ]);

  // State for modals and form
  const [showAddModal, setShowAddModal] = useState(false);
  const [showRoutesModal, setShowRoutesModal] = useState(false);
  const [selectedRoutes, setSelectedRoutes] = useState([]);
  const [newRole, setNewRole] = useState({
    name: "",
    routes: [""] // Start with one empty route
  });

  // Other existing state
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [statusFilter, setStatusFilter] = useState("All");
  const [roleTypeFilter, setRoleTypeFilter] = useState("All");
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [showDropdown, setShowDropdown] = useState(null);
  const dropdownRef = useRef(null);

  // Toggle dropdown with outside click handler
  const toggleDropdown = (roleId) => {
    setShowDropdown(showDropdown === roleId ? null : roleId);
  };

  useEffect(() => {
    function handleClickOutside(event) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setShowDropdown(null);
      }
    }
    if (showDropdown) {
      document.addEventListener("mousedown", handleClickOutside);
    }
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [showDropdown]);

  // Add a new route input field
  const addRouteField = () => {
    setNewRole({ ...newRole, routes: [...newRole.routes, ""] });
  };

  // Remove a route input field
  const removeRouteField = (index) => {
    const updatedRoutes = [...newRole.routes];
    updatedRoutes.splice(index, 1);
    setNewRole({ ...newRole, routes: updatedRoutes });
  };

  // Handle route input change
  const handleRouteChange = (index, value) => {
    const updatedRoutes = [...newRole.routes];
    updatedRoutes[index] = value;
    setNewRole({ ...newRole, routes: updatedRoutes });
  };

  // Submit new role
  const handleAddRole = () => {
    if (!newRole.name.trim()) return;
    
    const filteredRoutes = newRole.routes.filter(route => route.trim());
    if (filteredRoutes.length === 0) return;

    const newRoleObj = {
      id: crypto.randomUUID(),
      name: newRole.name,
      roleType: "custom",
      status: "Enabled",
      routes: filteredRoutes
    };

    setRoles([...roles, newRoleObj]);
    setNewRole({ name: "", routes: [""] });
    setShowAddModal(false);
  };

  // View routes for a role
  const viewRoleRoutes = (routes) => {
    setSelectedRoutes(routes);
    setShowRoutesModal(true);
  };

  // Delete role
  const deleteRole = (roleId) => {
    setRoles(roles.filter(role => role.id !== roleId));
    setShowDropdown(null);
  };

  // Filtering logic
  const filteredRoles = roles.filter(role => {
    const matchesStatus = statusFilter === "All" || role.status === statusFilter;
    const matchesRoleType = roleTypeFilter === "All" || role.roleType === roleTypeFilter.toLowerCase();
    const matchesSearch = role.name.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesStatus && matchesRoleType && matchesSearch;
  });

  const paginatedRoles = filteredRoles.slice(
    (currentPage - 1) * rowsPerPage,
    currentPage * rowsPerPage
  );

  const totalPages = Math.ceil(filteredRoles.length / rowsPerPage);

  return (
    <div className="antialiased font-sans bg-gray-50 min-h-screen">
      {/* Add Role Modal */}
      {showAddModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
            <div className="flex justify-between items-center border-b border-gray-200 px-6 py-4">
              <h3 className="text-lg font-medium text-gray-900">Add New Role</h3>
              <button 
                onClick={() => setShowAddModal(false)}
                className="text-gray-400 hover:text-gray-500"
              >
                <IoClose className="h-6 w-6" />
              </button>
            </div>
            <div className="px-6 py-4 space-y-4">
              <div>
                <label htmlFor="role-name" className="block text-sm font-medium text-gray-700 mb-1">
                  Role Name
                </label>
                <input
                  type="text"
                  id="role-name"
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500"
                  value={newRole.name}
                  onChange={(e) => setNewRole({ ...newRole, name: e.target.value })}
                  placeholder="Enter role name"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Routes
                </label>
                {newRole.routes.map((route, index) => (
                  <div key={index} className="flex items-center mb-2">
                    <input
                      type="text"
                      className="flex-1 border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500"
                      value={route}
                      onChange={(e) => handleRouteChange(index, e.target.value)}
                      placeholder={`Route path ${index + 1}`}
                    />
                    {newRole.routes.length > 1 && (
                      <button
                        type="button"
                        onClick={() => removeRouteField(index)}
                        className="ml-2 text-red-500 hover:text-red-700"
                      >
                        <IoClose className="h-5 w-5" />
                      </button>
                    )}
                  </div>
                ))}
                <button
                  type="button"
                  onClick={addRouteField}
                  className="mt-2 text-sm text-stone-600 hover:text-stone-800 flex items-center"
                >
                  <IoAddCircleSharp className="mr-1" /> Add another route
                </button>
              </div>
            </div>
            <div className="bg-gray-50 px-6 py-3 flex justify-end border-t border-gray-200">
              <button
                type="button"
                onClick={() => setShowAddModal(false)}
                className="mr-3 px-4 py-2 text-sm text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={handleAddRole}
                disabled={!newRole.name.trim() || newRole.routes.every(r => !r.trim())}
                className="px-4 py-2 text-sm text-white bg-stone-800 border border-transparent rounded-md hover:bg-stone-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Add Role
              </button>
            </div>
          </div>
        </div>
      )}

      {/* View Routes Modal */}
      {showRoutesModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
            <div className="flex justify-between items-center border-b border-gray-200 px-6 py-4">
              <h3 className="text-lg font-medium text-gray-900">Role Routes</h3>
              <button 
                onClick={() => setShowRoutesModal(false)}
                className="text-gray-400 hover:text-gray-500"
              >
                <IoClose className="h-6 w-6" />
              </button>
            </div>
            <div className="px-6 py-4">
              <ul className="space-y-2">
                {selectedRoutes.map((route, index) => (
                  <li key={index} className="flex items-center">
                    <span className="bg-gray-100 rounded-md px-3 py-1 text-sm font-mono">
                      {route}
                    </span>
                  </li>
                ))}
              </ul>
            </div>
            <div className="bg-gray-50 px-6 py-3 flex justify-end border-t border-gray-200">
              <button
                type="button"
                onClick={() => setShowRoutesModal(false)}
                className="px-4 py-2 text-sm text-white bg-stone-800 border border-transparent rounded-md hover:bg-stone-700"
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

          {/* Rest of your existing code (filters, table, pagination) */}
          {/* Only change needed is to update the "View Routes" button to call viewRoleRoutes */}
          <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
            <button 
              onClick={() => viewRoleRoutes(role.routes)}
              className="inline-flex items-center justify-center cursor-pointer border align-middle select-none font-sans font-medium text-center transition-all ease-in disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed focus:shadow-none text-sm py-2 px-4 shadow-sm bg-transparent relative text-stone-700 hover:text-stone-700 border-stone-500 hover:bg-transparent duration-150 hover:border-stone-600 rounded-lg hover:opacity-60 hover:shadow-none"
            >
              View Routes
            </button>
          </td>

          {/* ... rest of your existing component code ... */}
        </div>
      </div>
    </div>
  );
};

export default Temp;