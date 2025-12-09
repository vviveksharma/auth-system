import React, { useEffect, useState } from "react";
import { IoAddCircleSharp, IoClose } from "react-icons/io5";
import { FiCopy } from "react-icons/fi";
import { MdOutlineBlock } from "react-icons/md";
import { ListTokensPaginted, CreateTokens, RevokeToken, ListTokensWithStatus } from "../services/tokens";
const Tokens = () => {
  const [tokens, setTokens] = useState([]);
  const [loading, setLoading] = useState(true);
  
  // Initialize pagination from sessionStorage or use defaults
  const [pagination, setPagination] = useState(() => {
    const savedPagination = sessionStorage.getItem('tokensPagination');
    if (savedPagination) {
      try {
        const parsed = JSON.parse(savedPagination);
        return {
          page: parsed.page || 1,
          page_size: parsed.page_size || 5,
          total_pages: 1,
          total_items: 0,
          has_next: false,
          has_prev: false,
        };
      } catch (e) {
        console.error('Error parsing saved pagination:', e);
      }
    }
    return {
      page: 1,
      page_size: 5,
      total_pages: 1,
      total_items: 0,
      has_next: false,
      has_prev: false,
    };
  });

  // Modal states
  const [showAddModal, setShowAddModal] = useState(false);
  const [showRevokeModal, setShowRevokeModal] = useState(false);
  const [tokenToRevoke, setTokenToRevoke] = useState(null);
  const [newToken, setNewToken] = useState({
    name: "",
    expiry: "",
  });
  const [addingToken, setAddingToken] = useState(false);
  const [revokingToken, setRevokingToken] = useState(false);

  // Other existing state
  const [statusFilter, setStatusFilter] = useState("All");
  const [searchQuery, setSearchQuery] = useState("");
  const [copiedTokenId, setCopiedTokenId] = useState(null);

  // Save pagination to sessionStorage whenever it changes
  useEffect(() => {
    sessionStorage.setItem('tokensPagination', JSON.stringify({
      page: pagination.page,
      page_size: pagination.page_size,
    }));
  }, [pagination.page, pagination.page_size]);

  // Function to reload page (will restore pagination from sessionStorage)
  const pageReload = () => {
    window.location.reload();
  };

  // Helper function to format date
  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      month: "short",
      day: "2-digit",
      year: "numeric",
    });
  };

  // Transform API response to match UI structure
  const transformApiTokens = (apiTokens) => {
    return apiTokens.map((token) => ({
      id: token.token_id,
      name: token.name,
      created: formatDate(token.created_at),
      expiry: formatDate(token.expiry_at),
      status: token.status ? "Active" : "Revoked",
    }));
  };

  useEffect(() => {
    const fetchTokens = async () => {
      try {
        setLoading(true);
        
        let response;
        // Use different API based on status filter
        if (statusFilter === "Active" || statusFilter === "Revoked") {
          // Convert "Active" to "active", "Revoked" to "revoked" for API
          const statusValue = statusFilter.toLowerCase();
          response = await ListTokensWithStatus(statusValue, pagination.page, pagination.page_size);
        } else {
          // For "All" or any other value, use regular paginated list
          response = await ListTokensPaginted(pagination.page, pagination.page_size);
        }
        
        console.log("API Response:", response);
        
        // Extract tokens from nested response structure
        const apiTokens = response.data?.data || response.data?.data || [];
        const paginationData = response?.data?.pagination || {};
        
        const transformedTokens = transformApiTokens(apiTokens);
        
        setTokens(transformedTokens);
        setPagination(prev => ({
          ...prev,
          total_pages: paginationData.total_pages ?? 1,
          total_items: paginationData.total_items ?? 0,
          has_next: paginationData.has_next ?? false,
          has_prev: paginationData.has_prev ?? false,
        }));
      } catch (error) {
        console.error('Error fetching tokens:', error);
        // Keep empty array on error or set fallback data
        setTokens([]);
      } finally {
        setLoading(false);
      }
    };

    fetchTokens();
  }, [pagination.page, pagination.page_size, statusFilter]); 

  // Copy token ID
  const copyTokenId = (tokenId) => {
    navigator.clipboard.writeText(tokenId);
    setCopiedTokenId(tokenId);
    setTimeout(() => setCopiedTokenId(null), 2000);
  };

  // Prepare to revoke token
  const prepareRevokeToken = (tokenId) => {
    setTokenToRevoke(tokenId);
    setShowRevokeModal(true);
  };

  // Confirm revoke token
  const confirmRevokeToken = async () => {
    if (!tokenToRevoke) return;

    try {
      setRevokingToken(true);
      
      // Call API to revoke token
      const response = await RevokeToken(tokenToRevoke);
      console.log("Token revoked successfully:", response);
      
      // Reload page to refresh data (pagination will be restored from sessionStorage)
      pageReload();
      
    } catch (error) {
      console.error("Error revoking token:", error);
      setRevokingToken(false);
    }
  };

  // Add new token
  const addToken = async () => {
    if (!newToken.name.trim() || !newToken.expiry) return;

    try {
      setAddingToken(true);
      
      // Convert date to YYYY-MM-DD format for API
      const expiryDate = new Date(newToken.expiry).toISOString().split('T')[0];
      
      // Call API to create token
      const response = await CreateTokens(newToken.name, expiryDate);
      console.log("Token created successfully:", response);
      
      // Reload page to refresh data (pagination will be restored from sessionStorage)
      pageReload();
      
    } catch (error) {
      console.error("Error creating token:", error);
      setAddingToken(false);
    }
  };

  // Filter tokens (client-side search only, status is handled by API)
  const filteredTokens = tokens.filter((token) => {
    const matchesSearch =
      token.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
      token.name.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesSearch;
  });

  // Use filtered tokens for display (server already handles pagination and status filtering)
  const displayTokens = filteredTokens;

  return (
    <div className="antialiased font-sans bg-gray-100 min-h-screen">
      {/* Add Token Modal */}
      {showAddModal && (
        <div className="fixed inset-0 bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
            <div className="flex justify-between items-center border-b border-gray-200 px-6 py-4">
              <h3 className="text-lg font-medium text-gray-900">
                Create New Token
              </h3>
              <button
                onClick={() => setShowAddModal(false)}
                className="text-gray-400 hover:text-gray-500"
              >
                <IoClose className="h-6 w-6 cursor-pointer" />
              </button>
            </div>
            <div className="px-6 py-4 space-y-4">
              <div>
                <label
                  htmlFor="token-name"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  Token Name
                </label>
                <input
                  type="text"
                  id="token-name"
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500"
                  value={newToken.name}
                  onChange={(e) =>
                    setNewToken({ ...newToken, name: e.target.value })
                  }
                  placeholder="Enter token name"
                />
              </div>
              <div>
                <label
                  htmlFor="token-expiry"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  Expiry Date
                </label>
                <input
                  type="date"
                  id="token-expiry"
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-1 focus:ring-stone-500"
                  value={newToken.expiry}
                  min={new Date().toISOString().split("T")[0]}
                  onChange={(e) =>
                    setNewToken({ ...newToken, expiry: e.target.value })
                  }
                />
              </div>
            </div>
            <div className="bg-gray-50 px-6 py-3 flex justify-end border-t border-gray-200">
              <button
                type="button"
                onClick={() => setShowAddModal(false)}
                className="mr-3 px-4 cursor-pointer py-2 text-sm text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={addToken}
                disabled={!newToken.name.trim() || !newToken.expiry || addingToken}
                className="px-4 py-2 text-sm text-white cursor-pointer bg-stone-800 border border-transparent rounded-md hover:bg-stone-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {addingToken ? "Creating..." : "Create Token"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Revoke Confirmation Modal */}
      {showRevokeModal && (
        <div className="fixed inset-0  bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
            <div className="flex justify-between items-center border-b border-gray-200 px-6 py-4">
              <h3 className="text-lg font-medium text-gray-900">
                Revoke Token
              </h3>
              <button
                onClick={() => setShowRevokeModal(false)}
                className="text-gray-400 cursor-pointer hover:text-gray-500"
              >
                <IoClose className="h-6 w-6" />
              </button>
            </div>
            <div className="px-6 py-4">
              <p className="text-gray-700">
                Are you sure you want to revoke this token? This action cannot
                be undone.
              </p>
            </div>
            <div className="bg-gray-50 px-6 py-3 flex justify-end border-t border-gray-200">
              <button
                type="button"
                onClick={() => setShowRevokeModal(false)}
                className="mr-3 px-4 py-2 text-sm text-gray-700 cursor-pointer bg-white border border-gray-300 rounded-md hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={confirmRevokeToken}
                disabled={revokingToken}
                className="px-4 py-2 text-sm text-white bg-red-600 cursor-pointer border border-transparent rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {revokingToken ? "Revoking..." : "Revoke Token"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Main Content */}
      <div className="container mx-auto px-2 sm:px-8">
        <div className="py-4">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-2xl font-semibold leading-tight">Tokens</h2>
            <button
              onClick={() => setShowAddModal(true)}
              className="flex items-center gap-2 justify-center border align-middle select-none font-sans font-medium text-center duration-300 ease-in disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed focus:shadow-none text-sm py-2 px-4 shadow-sm hover:shadow-md bg-stone-800 hover:bg-stone-700 border-stone-900 text-stone-50 rounded-lg transition antialiased cursor-pointer"
            >
              <IoAddCircleSharp className="text-lg" />
              Add Token
            </button>
          </div>

          {/* Filters and Search */}
          <div className="my-2 flex flex-col sm:flex-row gap-2 items-stretch">
            <div className="flex flex-col sm:flex-row mb-1 sm:mb-0 gap-2 flex-1">
              <div className="relative">
                <select
                  className="appearance-none h-full rounded border block w-full bg-white border-gray-400 text-gray-700 py-2 px-4 pr-8 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                  value={pagination.page_size}
                  onChange={(e) => {
                    const newPageSize = Number(e.target.value);
                    setPagination(prev => ({
                      ...prev,
                      page: 1, // Reset to page 1 when page size changes
                      page_size: newPageSize,
                    }));
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
                    // Reset to page 1 when status filter changes
                    setPagination(prev => ({
                      ...prev,
                      page: 1,
                    }));
                  }}
                >
                  <option value="All">All Status</option>
                  <option value="Active">Active</option>
                  <option value="Revoked">Revoked</option>
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
                placeholder="Search tokens..."
                className="appearance-none rounded-r rounded-l sm:rounded-l-none border border-gray-400 border-b block pl-8 pr-6 py-2 w-full bg-white text-sm placeholder-gray-400 text-gray-700 focus:bg-white focus:placeholder-gray-600 focus:text-gray-700 focus:outline-none"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
          </div>

          {/* Tokens Table */}
          <div className="-mx-2 sm:-mx-8 px-2 sm:px-8 py-4 overflow-x-auto">
            <div className="inline-block min-w-full shadow rounded-lg overflow-hidden">
              <table className="min-w-full leading-normal">
                <thead>
                  <tr>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Token ID
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Name
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Created at
                    </th>
                    <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Expiry at
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
                        colSpan="6"
                        className="px-5 py-5 border-b border-gray-200 bg-white text-sm text-center"
                      >
                        Loading tokens...
                      </td>
                    </tr>
                  ) : displayTokens.length > 0 ? (
                    displayTokens.map((token) => (
                      <tr key={token.id} className="hover:bg-gray-50">
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <div className="flex items-center group">
                            <p className="text-gray-900 whitespace-no-wrap font-mono text-sm">
                              {token.id.substring(0, 8)}...
                            </p>
                            <button
                              onClick={() => copyTokenId(token.id)}
                              className="ml-2 text-gray-400 hover:text-gray-600 opacity-0 group-hover:opacity-100 transition-opacity"
                              title="Copy Token ID"
                            >
                              <FiCopy className="h-4 w-4 cursor-pointer" />
                            </button>
                            {copiedTokenId === token.id && (
                              <span className="ml-2 text-xs text-green-600">
                                Copied!
                              </span>
                            )}
                          </div>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <p className="text-gray-900 whitespace-no-wrap">
                            {token.name}
                          </p>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <p className="text-gray-900 whitespace-no-wrap">
                            {token.created}
                          </p>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <p className="text-gray-900 whitespace-no-wrap">
                            {token.expiry}
                          </p>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          <span
                            className={`relative inline-block px-3 py-1 font-semibold leading-tight ${
                              token.status === "Active"
                                ? "text-green-900"
                                : "text-red-900"
                            }`}
                          >
                            <span
                              aria-hidden
                              className={`absolute inset-0 rounded-full opacity-50 ${
                                token.status === "Active"
                                  ? "bg-green-200"
                                  : "bg-red-200"
                              }`}
                            ></span>
                            <span className="relative">{token.status}</span>
                          </span>
                        </td>
                        <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                          {token.status === "Active" ? (
                            <button
                              onClick={() => prepareRevokeToken(token.id)}
                              className="flex items-center text-sm text-red-600 hover:text-red-800 transition-colors cursor-pointer"
                            >
                              <MdOutlineBlock className="mr-1" />
                              Revoke
                            </button>
                          ) : (
                            <span className="text-gray-400 text-sm">
                              Revoked
                            </span>
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
                        No tokens found matching your criteria
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>

              {/* Pagination */}
              <div className="px-5 py-5 bg-white border-t flex flex-col sm:flex-row items-center justify-between gap-2">
                <span className="text-xs xs:text-sm text-gray-900">
                  {loading ? (
                    "Loading..."
                  ) : (
                    <>
                      Showing page {pagination.page} of {pagination.total_pages} ({pagination.total_items} total tokens)
                    </>
                  )}
                </span>
                <div className="inline-flex mt-2 xs:mt-0">
                  <button
                    onClick={() =>
                      setPagination(prev => ({
                        ...prev,
                        page: prev.page - 1,
                      }))
                    }
                    disabled={!pagination.has_prev || loading}
                    className={`text-sm py-2 px-4 rounded-l cursor-pointer ${
                      !pagination.has_prev || loading
                        ? "bg-gray-200 text-gray-500 cursor-not-allowed"
                        : "bg-gray-300 hover:bg-gray-400 text-gray-800"
                    }`}
                  >
                    Prev
                  </button>
                  <button
                    onClick={() =>
                      setPagination(prev => ({
                        ...prev,
                        page: prev.page + 1,
                      }))
                    }
                    disabled={!pagination.has_next || loading}
                    className={`text-sm py-2 px-4 rounded-r cursor-pointer ${
                      !pagination.has_next || loading
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

export default Tokens;
