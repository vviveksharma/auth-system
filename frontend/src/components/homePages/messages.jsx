import React, { useState, useEffect } from "react";
import { FiCheck, FiX } from "react-icons/fi";
import { ListMessages, ApproveMessage, RejectMessage } from "../services/messages";

const Messages = () => {
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState("All");
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalMessages, setTotalMessages] = useState(0);

  // Fetch messages from API
  useEffect(() => {
    fetchMessages();
  }, [currentPage, rowsPerPage, statusFilter]);

  const fetchMessages = async () => {
    try {
      setLoading(true);
      setError(null);

      // Convert status filter to API format
      let apiStatus = "";
      if (statusFilter === "Pending") {
        apiStatus = "pending";
      } else if (statusFilter === "Approved") {
        apiStatus = "approved";
      } else if (statusFilter === "Rejected") {
        apiStatus = "rejected";
      }

      const response = await ListMessages(currentPage, rowsPerPage, apiStatus);

      // Extract data from response
      const messageData = response.data?.data || [];
      const paginationData = response.data?.pagination || {};

      // Transform API response to match component structure
      const transformedMessages = messageData.map((msg) => ({
        id: msg.message_id,
        email: msg.email || msg.user_email,
        currentRole: msg.current_role || msg.currentRole,
        requestedRole: msg.requested_role || msg.requestedRole,
        status: msg.status.charAt(0).toUpperCase() + msg.status.slice(1),
        requestedAt: msg.request_at || msg.requested_at || msg.requestedAt || msg.created_at,
      }));

      setRequests(transformedMessages);
      setTotalMessages(paginationData.total_items || transformedMessages.length);
    } catch (err) {
      console.error("Error fetching messages:", err);
      setError("Failed to load messages. Please try again.");
      setRequests([]);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (requestId) => {
    try {
      // Optimistically update UI
      setRequests(
        requests.map((request) =>
          request.id === requestId
            ? { ...request, status: "Approved" }
            : request
        )
      );

      // Call API
      await ApproveMessage(requestId);
    } catch (err) {
      console.error("Error approving message:", err);
      // Revert on error
      await fetchMessages();
      alert("Failed to approve request. Please try again.");
    }
  };

  const handleReject = async (requestId) => {
    try {
      // Optimistically update UI
      setRequests(
        requests.map((request) =>
          request.id === requestId
            ? { ...request, status: "Rejected" }
            : request
        )
      );

      // Call API
      await RejectMessage(requestId);
    } catch (err) {
      console.error("Error rejecting message:", err);
      // Revert on error
      await fetchMessages();
      alert("Failed to reject request. Please try again.");
    }
  };

  // Client-side search filtering
  const filteredRequests = requests.filter((request) => {
    if (searchQuery === "") return true;
    const matchesSearch =
      request.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
      request.requestedRole.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesSearch;
  });

  const totalPages = Math.ceil(totalMessages / rowsPerPage);

  const getStatusColor = (status) => {
    switch (status) {
      case "Pending": 
        return "text-yellow-900 bg-yellow-200";
      case "Approved": 
        return "text-green-900 bg-green-200";
      case "Rejected": 
        return "text-red-900 bg-red-200";
      default: 
        return "text-gray-900 bg-gray-200";
    }
  };

  return (
    <div className="antialiased font-sans bg-gray-100 min-h-screen">
      <div className="container mx-auto px-2 sm:px-8">
        <div className="py-4">
          {/* Header */}
          <div className="flex justify-between items-center mb-5">
            <h2 className="text-2xl font-semibold leading-tight">
              Messages
            </h2>
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
                    setCurrentPage(1);
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
                    setCurrentPage(1);
                  }}
                >
                  <option value="All">All Status</option>
                  <option value="Pending">Pending</option>
                  <option value="Approved">Approved</option>
                  <option value="Rejected">Rejected</option>
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
                placeholder="Search by email"
                className="appearance-none rounded-r rounded-l sm:rounded-l-none border border-gray-400 border-b block pl-8 pr-6 py-2 w-full bg-white text-sm placeholder-gray-400 text-gray-700 focus:bg-white focus:placeholder-gray-600 focus:text-gray-700 focus:outline-none"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
          </div>

          {/* Requests Table */}
          <div className="-mx-2 sm:-mx-8 px-2 sm:px-8 py-4 overflow-x-auto">
            <div className="inline-block min-w-full shadow rounded-lg overflow-hidden">
              <table className="min-w-full leading-normal">
                <thead>
                  <tr>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      User
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Current Role
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Requested Role
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Requested At
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
                          <span className="ml-3">Loading requests...</span>
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
                            onClick={fetchMessages}
                            className="mt-2 text-sm text-blue-600 hover:text-blue-800 underline"
                          >
                            Retry
                          </button>
                        </div>
                      </td>
                    </tr>
                  ) : filteredRequests.length > 0 ? (
                    filteredRequests.map((request) => (
                      <tr key={request.id} className="hover:bg-gray-50">
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <div className="flex items-center">
                            <div className="flex-shrink-0 w-10 h-10 rounded-full bg-stone-200 flex items-center justify-center">
                              <span className="text-stone-600 font-medium">
                                {request.email.charAt(0).toUpperCase()}
                              </span>
                            </div>
                            <div className="ml-3">
                              <p className="text-gray-900 whitespace-no-wrap font-medium">
                                {request.email}
                              </p>
                            </div>
                          </div>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <p className="text-gray-900 whitespace-no-wrap">
                            {request.currentRole}
                          </p>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <p className="text-gray-900 whitespace-no-wrap font-medium">
                            {request.requestedRole}
                          </p>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <span
                            className={`relative inline-block px-3 py-1 font-semibold rounded-full leading-tight ${getStatusColor(request.status)}`}
                          >
                            <span
                              aria-hidden
                              className={`absolute inset-0 rounded-full opacity-50 ${getStatusColor(request.status).split(' ')[1]}`}
                            ></span>
                            <span className="relative">{request.status}</span>
                          </span>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <p className="text-gray-900 whitespace-no-wrap">
                            {new Date(request.requestedAt).toLocaleDateString('en-US', {
                              year: 'numeric',
                              month: 'short',
                              day: 'numeric'
                            })}
                          </p>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          {request.status === "Pending" ? (
                            <div className="flex items-center gap-2">
                              <button
                                onClick={() => handleApprove(request.id)}
                                className="flex items-center text-sm text-green-600 hover:text-green-800 transition-colors cursor-pointer"
                              >
                                <FiCheck className="mr-1" />
                                Approve
                              </button>
                              <button
                                onClick={() => handleReject(request.id)}
                                className="flex items-center text-sm text-red-600 hover:text-red-800 transition-colors cursor-pointer"
                              >
                                <FiX className="mr-1" />
                                Reject
                              </button>
                            </div>
                          ) : (
                            <span className="text-gray-400 text-sm">—</span>
                          )}
                        </td>
                      </tr>
                    ))
                  ) : (
                    <tr>
                      <td
                        colSpan="6"
                        className="px-5 py-5 border-b border-gray-200 bg-white text-sm text-center"
                      >
                        No requests found matching your criteria
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>

              {/* Pagination */}
              <div className="px-5 py-5 bg-white border-t flex flex-col sm:flex-row items-center justify-between gap-2">
                <span className="text-xs xs:text-sm text-gray-900">
                  Showing {Math.min((currentPage - 1) * rowsPerPage + 1, totalMessages)} to{" "}
                  {Math.min(currentPage * rowsPerPage, totalMessages)} of{" "}
                  {totalMessages} requests
                </span>
                <div className="inline-flex mt-2 xs:mt-0">
                  <button
                    onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
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
                    onClick={() => setCurrentPage((prev) => Math.min(prev + 1, totalPages))}
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

export default Messages;