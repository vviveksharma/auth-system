import React, { useState, useEffect } from "react";
import {
  RiUserAddFill,
  RiEditLine,
  RiDeleteBinLine,
  RiMore2Fill,
} from "react-icons/ri";
import { FiUserX, FiUserCheck } from "react-icons/fi";
import { ListUsers } from "../services/user";

const Users = () => {
  const [users, setUsers] = useState([]);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [statusFilter, setStatusFilter] = useState("All");
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [showDropdown, setShowDropdown] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [totalUsers, setTotalUsers] = useState(0);

  // Fetch users from API
  useEffect(() => {
    fetchUsers();
  }, [currentPage, rowsPerPage, statusFilter]);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Convert statusFilter to lowercase for API (Active -> active, Inactive -> inactive)
      const apiStatus = statusFilter === "All" ? "" : statusFilter.toLowerCase();
      
      const response = await ListUsers(apiStatus, currentPage, rowsPerPage);
      
      // Extract data from nested response structure
      const userData = response.data?.data || [];
      const paginationData = response.data?.pagination || {};
      
      // Transform API response to match the component's data structure
      const transformedUsers = userData.map((user, index) => ({
        id: user.id || user.user_id || `user-${index}`,
        name: user.name || user.username || "Unknown User",
        email: user.email,
        role: user.roles && user.roles.length > 0 ? user.roles.join(", ") : (user.role || "User"),
        created: user.created_at 
          ? new Date(user.created_at).toLocaleDateString('en-US', { 
              year: 'numeric', 
              month: 'short', 
              day: 'numeric' 
            })
          : new Date().toLocaleDateString(),
        status: user.status ? "Active" : "Inactive",
      }));
      
      setUsers(transformedUsers);
      setTotalUsers(paginationData.total_items || transformedUsers.length);
    } catch (err) {
      console.error("Error fetching users:", err);
      setError("Failed to load users. Please try again.");
      setUsers([]);
    } finally {
      setLoading(false);
    }
  };

  const toggleDropdown = (userId) => {
    setShowDropdown(showDropdown === userId ? null : userId);
  };

  const deleteUser = (userId) => {
    // TODO: Add API call to delete user
    setUsers(users.filter((user) => user.id !== userId));
    setShowDropdown(null);
    // Optionally refetch after delete
    // fetchUsers();
  };

  const toggleUserStatus = (userId) => {
    // TODO: Add API call to update user status
    setUsers(
      users.map((user) =>
        user.id === userId
          ? {
              ...user,
              status: user.status === "Active" ? "Inactive" : "Active",
            }
          : user
      )
    );
    setShowDropdown(null);
    // Optionally refetch after status change
    // fetchUsers();
  };

  // Client-side search filtering (searches within current page results)
  const filteredUsers = users.filter((user) => {
    if (searchQuery === "") return true;
    const matchesSearch =
      user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      user.email.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesSearch;
  });

  const totalPages = Math.ceil(totalUsers / rowsPerPage);

  return (
    <div className="antialiased font-sans bg-gray-100 min-h-screen">
      <div className="container mx-auto px-2 sm:px-8">
        <div className="py-4">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-2xl font-semibold leading-tight">
              User Management
            </h2>
            <button className="flex items-center gap-2 justify-center border align-middle select-none font-sans font-medium text-center duration-300 ease-in disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed focus:shadow-none text-sm py-2 px-4 shadow-sm hover:shadow-md bg-stone-800 hover:bg-stone-700 border-stone-900 text-stone-50 rounded-lg transition antialiased cursor-pointer">
              <RiUserAddFill className="text-lg" />
              Add User
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
                  <option value="Active">Active</option>
                  <option value="Inactive">Inactive</option>
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
                placeholder="Search users..."
                className="appearance-none rounded-r rounded-l sm:rounded-l-none border border-gray-400 border-b block pl-8 pr-6 py-2 w-full bg-white text-sm placeholder-gray-400 text-gray-700 focus:bg-white focus:placeholder-gray-600 focus:text-gray-700 focus:outline-none"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
          </div>

          {/* Users Table */}
          <div className="-mx-2 sm:-mx-8 px-2 sm:px-8 py-4 overflow-x-auto">
            <div className="inline-block min-w-full shadow rounded-lg overflow-hidden">
              <table className="min-w-full leading-normal">
                <thead>
                  <tr>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      User
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Role
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Created at
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Status
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
                        colSpan="5"
                        className="px-5 py-10 border-b border-gray-200 bg-white text-sm text-center"
                      >
                        <div className="flex justify-center items-center">
                          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
                          <span className="ml-3">Loading users...</span>
                        </div>
                      </td>
                    </tr>
                  ) : error ? (
                    <tr>
                      <td
                        colSpan="5"
                        className="px-5 py-10 border-b border-gray-200 bg-white text-sm text-center"
                      >
                        <div className="text-red-600">
                          <p>{error}</p>
                          <button
                            onClick={fetchUsers}
                            className="mt-2 text-sm text-blue-600 hover:text-blue-800 underline"
                          >
                            Retry
                          </button>
                        </div>
                      </td>
                    </tr>
                  ) : filteredUsers.length > 0 ? (
                    filteredUsers.map((user) => (
                      <tr key={user.id} className="hover:bg-gray-50">
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <div className="flex items-center">
                            <div className="flex-shrink-0 w-10 h-10 rounded-full bg-stone-200 flex items-center justify-center">
                              <span className="text-stone-600 font-medium">
                                {user.name.charAt(0)}
                              </span>
                            </div>
                            <div className="ml-3">
                              <p className="text-gray-900 whitespace-no-wrap font-medium">
                                {user.name}
                              </p>
                              <p className="text-gray-500 whitespace-no-wrap text-xs">
                                {user.email}
                              </p>
                            </div>
                          </div>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <p className="text-gray-900 whitespace-no-wrap">
                            {user.role}
                          </p>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <p className="text-gray-900 whitespace-no-wrap">
                            {user.created}
                          </p>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <span
                            className={`relative inline-block px-3 py-1 font-semibold leading-tight ${
                              user.status === "Active"
                                ? "text-green-900"
                                : "text-gray-900"
                            }`}
                          >
                            <span
                              aria-hidden
                              className={`absolute inset-0 rounded-full opacity-50 ${
                                user.status === "Active"
                                  ? "bg-green-200"
                                  : "bg-gray-300"
                              }`}
                            ></span>
                            <span className="relative">{user.status}</span>
                          </span>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <div className="relative">
                            <button
                              onClick={() => toggleDropdown(user.id)}
                              className="text-gray-500 hover:text-gray-700 focus:outline-none"
                            >
                              <RiMore2Fill className="h-5 w-5" />
                            </button>
                            {showDropdown === user.id && (
                              <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-10 border border-gray-200">
                                <div className="py-1">
                                  <button
                                    onClick={() => toggleUserStatus(user.id)}
                                    className="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 w-full text-left"
                                  >
                                    {user.status === "Active" ? (
                                      <FiUserX className="mr-2" />
                                    ) : (
                                      <FiUserCheck className="mr-2" />
                                    )}
                                    {user.status === "Active"
                                      ? "Deactivate"
                                      : "Activate"}
                                  </button>
                                  <button
                                    onClick={() => deleteUser(user.id)}
                                    className="flex items-center px-4 py-2 text-sm text-red-600 hover:bg-gray-100 w-full text-left"
                                  >
                                    <RiDeleteBinLine className="mr-2" />
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
                        colSpan="5"
                        className="px-5 py-5 border-b border-gray-200 bg-white text-sm text-center"
                      >
                        No users found matching your criteria
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>

              {/* Pagination */}
              <div className="px-5 py-5 bg-white border-t flex flex-col sm:flex-row items-center justify-between gap-2">
                <span className="text-xs xs:text-sm text-gray-900">
                  Showing {Math.min((currentPage - 1) * rowsPerPage + 1, totalUsers)} to{" "}
                  {Math.min(currentPage * rowsPerPage, totalUsers)} of{" "}
                  {totalUsers} users
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

export default Users;
