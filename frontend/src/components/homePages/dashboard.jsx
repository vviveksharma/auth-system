import React, { useState, useEffect } from "react";
import {
  FiKey,
  FiShield,
  FiEdit2,
  FiChevronRight,
  FiRefreshCw,
  FiUsers,
} from "react-icons/fi";
import { GetDashBoardDetails, GetActiveRoles } from "../services/dashboard";
const Dashboard = () => {
  function pageReload() {
    location.reload();
  }
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    users: 0,
    tokens: 0,
    roles: 0,
  });
  const [users, setUsers] = useState([]);
  const [ActiveToken, setActiveToken] = useState(1)

  // Simulate data loading
  useEffect(() => {
    const fetchData = async () => {
      // In a real app, you'd fetch from your API
      await new Promise((resolve) => setTimeout(resolve, 800));
      const response = await GetDashBoardDetails()
      console.log(response)
      setStats({
        users: response.data.user_count,
        tokens: response.data.token_count,
        roles: response.data.role_count,
      });

      const TokensResponse = await GetActiveRoles(1,5)
      console.log("the token response: ",TokensResponse.data.pagination.total_items)
      setActiveToken(TokensResponse.data.pagination.total_items)
      setUsers([
        {
          id: 1,
          name: "John Doe",
          email: "john@example.com",
          title: "Software Engineer",
          department: "Web Development",
          status: "active",
          role: "Owner",
          lastActive: "2 hours ago",
        },
        {
          id: 2,
          name: "Jane Smith",
          email: "jane@example.com",
          title: "Product Manager",
          department: "Product",
          status: "active",
          role: "Admin",
          lastActive: "30 minutes ago",
        },
        {
          id: 3,
          name: "Robert Johnson",
          email: "robert@example.com",
          title: "QA Engineer",
          department: "Testing",
          status: "inactive",
          role: "Member",
          lastActive: "3 days ago",
        },
      ]);

      setLoading(false);
    };

    fetchData();
  }, []);

  const navigateTo = (section) => {
    sessionStorage.setItem("activeSection", section);
    window.location.reload();
  };

  const refreshData = () => {
    setLoading(true);
    // In a real app, you would refetch data here
    setTimeout(() => setLoading(false), 500);
  };

  const formatNumber = (num) => {
    return new Intl.NumberFormat().format(num);
  };

  return (
    <div className="flex min-h-screen bg-gray-100">
      <div className="flex flex-col flex-1 overflow-hidden">
        <main className="flex-1 overflow-x-hidden overflow-y-auto bg-gray-100">
          <div className="container px-6 py-8 mx-auto">
            {/* Header with refresh button */}
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-3xl font-medium text-gray-700">Dashboard</h3>
              <button
                onClick={refreshData}
                className="flex items-center px-3 py-2 text-sm text-gray-600 bg-white rounded-md shadow-sm hover:bg-gray-50 focus:outline-none"
                disabled={loading}
              >
                <FiRefreshCw
                  className={`mr-2 ${loading ? "animate-spin" : ""}`}
                />
                Refresh
              </button>
            </div>

            {/* Stats Cards */}
            <div className="mt-4">
              <div className="flex flex-wrap -mx-6">
                {/* Users Card */}
                <div className="w-full px-6 sm:w-1/2 xl:w-1/3">
                  <div className="relative flex items-center px-5 py-6 bg-white rounded-md shadow-sm hover:shadow-md transition-shadow duration-200">
                    <button
                      onClick={() => {
                        sessionStorage.setItem("activeSection", "users");
                        pageReload();
                      }}
                      className="absolute top-3 cursor-pointer right-3 flex items-center text-stone-400 hover:text-stone-700 transition"
                      aria-label="View users"
                    >
                      <span className="sr-only">View users</span>
                      <FiChevronRight className="w-5 h-5" />
                    </button>
                    <div className="p-3 bg-stone-800 bg-opacity-75 rounded-full">
                      <FiUsers className="w-8 h-8 text-white" />
                    </div>
                    <div className="mx-5">
                      <h4 className="text-2xl font-semibold text-gray-700">
                        {loading ? "--" : formatNumber(stats.users)}
                      </h4>
                      <div className="text-gray-500">Total Users</div>
                      <div className="text-xs text-gray-400 mt-1">
                        {loading ? "Loading..." : "Last updated just now"}
                      </div>
                    </div>
                  </div>
                </div>

                {/* Tokens Card */}
                <div className="w-full px-6 mt-6 sm:w-1/2 xl:w-1/3 sm:mt-0">
                  <div className="relative flex items-center px-5 py-6 bg-white rounded-md shadow-sm hover:shadow-md transition-shadow duration-200">
                    <button
                      onClick={() => {
                        sessionStorage.setItem("activeSection", "tokens");
                        pageReload();
                      }}
                      className="absolute top-3 right-3 cursor-pointer flex items-center text-stone-400 hover:text-stone-700 transition"
                      aria-label="View tokens"
                    >
                      <span className="sr-only">View tokens</span>
                      <FiChevronRight className="w-5 h-5" />
                    </button>
                    <div className="p-3 bg-stone-800 bg-opacity-75 rounded-full">
                      <FiKey className="w-8 h-8 text-white" />
                    </div>
                    <div className="mx-5">
                      <h4 className="text-2xl font-semibold text-gray-700">
                        {loading ? "--" : formatNumber(stats.tokens)}
                      </h4>
                      <div className="text-gray-500">Total Tokens</div>
                      <div className="text-xs text-gray-400 mt-1">
                        {loading ? "Loading..." : `Active:${ActiveToken}`}
                      </div>
                    </div>
                  </div>
                </div>

                {/* Roles Card */}
                <div className="w-full px-6 mt-6 sm:w-1/2 xl:w-1/3 xl:mt-0">
                  <div className="relative flex items-center px-5 py-6 bg-white rounded-md shadow-sm hover:shadow-md transition-shadow duration-200">
                    <button
                      onClick={() => {
                        sessionStorage.setItem("activeSection", "roles");
                        pageReload();
                      }}
                      className="absolute top-3 right-3 flex cursor-pointer items-center text-stone-400 hover:text-stone-700 transition"
                      aria-label="View roles"
                    >
                      <span className="sr-only">View roles</span>
                      <FiChevronRight className="w-5 h-5" />
                    </button>
                    <div className="p-3 bg-stone-800 bg-opacity-75 rounded-full">
                      <FiShield className="w-8 h-8 text-white" />
                    </div>
                    <div className="mx-5">
                      <h4 className="text-2xl font-semibold text-gray-700">
                        {loading ? "--" : stats.roles}
                      </h4>
                      <div className="text-gray-500">Custom Roles</div>
                      <div className="text-xs text-gray-400 mt-1">
                        {loading ? "Loading..." : "5 system roles"}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Recent Users Table */}
            <div className="flex flex-col mt-8">
              <div className="py-2 -my-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8">
                <div className="inline-block min-w-full overflow-hidden align-middle border-b border-gray-200 shadow sm:rounded-lg">
                  <div className="flex items-center justify-between px-6 py-4 bg-white border-b">
                    <h3 className="text-xl font-medium text-gray-700">
                      Recent Users
                    </h3>
                    <button
                      onClick={() => navigateTo("users")}
                      className="text-sm text-stone-600 hover:text-stone-900"
                    >
                      View All
                    </button>
                  </div>

                  <div className="overflow-x-auto">
                    <table className="min-w-full">
                      <thead>
                        <tr>
                          <th className="px-6 py-3 text-xs font-medium leading-4 tracking-wider text-left text-gray-500 uppercase bg-gray-50">
                            User
                          </th>
                          <th className="px-6 py-3 text-xs font-medium leading-4 tracking-wider text-left text-gray-500 uppercase bg-gray-50">
                            Details
                          </th>
                          <th className="px-6 py-3 text-xs font-medium leading-4 tracking-wider text-left text-gray-500 uppercase bg-gray-50">
                            Status
                          </th>
                          <th className="px-6 py-3 text-xs font-medium leading-4 tracking-wider text-left text-gray-500 uppercase bg-gray-50">
                            Role
                          </th>
                          <th className="px-6 py-3 text-xs font-medium leading-4 tracking-wider text-left text-gray-500 uppercase bg-gray-50">
                            Last Active
                          </th>
                          <th className="px-6 py-3 bg-gray-50"></th>
                        </tr>
                      </thead>

                      <tbody className="bg-white divide-y divide-gray-200">
                        {loading ? (
                          <tr>
                            <td
                              colSpan="6"
                              className="px-6 py-4 text-center text-gray-500"
                            >
                              Loading user data...
                            </td>
                          </tr>
                        ) : (
                          users.map((user) => (
                            <tr key={user.id} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap">
                                <div className="flex items-center">
                                  <div className="flex-shrink-0 w-10 h-10 rounded-full bg-stone-200 flex items-center justify-center">
                                    <span className="text-stone-600 font-medium">
                                      {user.name.charAt(0)}
                                    </span>
                                  </div>
                                  <div className="ml-4">
                                    <div className="text-sm font-medium text-gray-900">
                                      {user.name}
                                    </div>
                                    <div className="text-sm text-gray-500">
                                      {user.email}
                                    </div>
                                  </div>
                                </div>
                              </td>

                              <td className="px-6 py-4 whitespace-nowrap">
                                <div className="text-sm text-gray-900">
                                  {user.title}
                                </div>
                                <div className="text-sm text-gray-500">
                                  {user.department}
                                </div>
                              </td>

                              <td className="px-6 py-4 whitespace-nowrap">
                                <span
                                  className={`inline-flex px-2 text-xs font-semibold leading-5 rounded-full ${
                                    user.status === "active"
                                      ? "text-green-800 bg-green-100"
                                      : "text-gray-800 bg-gray-100"
                                  }`}
                                >
                                  {user.status.charAt(0).toUpperCase() +
                                    user.status.slice(1)}
                                </span>
                              </td>

                              <td className="px-6 py-4 text-sm text-gray-500 whitespace-nowrap">
                                {user.role}
                              </td>

                              <td className="px-6 py-4 text-sm text-gray-500 whitespace-nowrap">
                                {user.lastActive}
                              </td>

                              <td className="px-6 py-4 text-sm font-medium text-right whitespace-nowrap">
                                <button
                                  onClick={() => {
                                    sessionStorage.setItem(
                                      "activeSection",
                                      "profile"
                                    );
                                    sessionStorage.setItem(
                                      "selectedUserId",
                                      user.id
                                    );
                                    window.location.reload();
                                  }}
                                  className="flex items-center text-stone-600 hover:text-stone-900"
                                  aria-label={`Edit ${user.name}`}
                                >
                                  <FiEdit2 className="mr-1" />
                                  <span>Edit</span>
                                </button>
                              </td>
                            </tr>
                          ))
                        )}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};

export default Dashboard;
