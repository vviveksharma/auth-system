import React, { useEffect, useState } from "react";
import { CgOrganisation } from "react-icons/cg";
import { GetProfile, DeleteProfile } from "../services/profile";
import { useNavigate } from "react-router-dom";
import {
  FiMail,
  FiUser,
  FiLock,
  FiRefreshCw,
  FiTrash2,
  FiEdit2,
  FiCheck,
  FiEye,
  FiEyeOff,
  FiX,
} from "react-icons/fi";

const Profile = () => {
  const navigate = useNavigate();
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [password, setPassword] = useState("");
  const [deletingAccount, setDeletingAccount] = useState(false);
  const [userDetails, setUserDetails] = useState({
    name: "",
    email: "",
    organization: "",
  });
  const [passwordRequirements, setPasswordRequirements] = useState({
    hasMinLength: false,
    hasUppercase: false,
    hasLowercase: false,
    hasNumber: false,
    hasSpecialChar: false,
  });

  const checkPasswordRequirements = (password) => {
    const requirements = {
      hasMinLength: password.length >= 8,
      hasUppercase: /[A-Z]/.test(password),
      hasLowercase: /[a-z]/.test(password),
      hasNumber: /[0-9]/.test(password),
      hasSpecialChar: /[!@#$%^&*(),.?":{}|<>]/.test(password),
    };
    setPasswordRequirements(requirements);
    return Object.values(requirements).every((req) => req);
  };

  const handlePasswordChange = (e) => {
    const newPassword = e.target.value;
    setPassword(newPassword);
    checkPasswordRequirements(newPassword);
  };

  const isPasswordValid = Object.values(passwordRequirements).every(
    (req) => req
  );

   useEffect(() => {
    const fetchProfile = async () => {
      const response = await GetProfile();
      setUserDetails({
        name: response.data.name,
        email: response.data.email,
        organization: response.data.organisation
      })
      console.log(response);
    };
    fetchProfile();
  }, []);

  // Handle account deletion
  const handleDeleteAccount = async () => {
    try {
      setDeletingAccount(true);
      
      // Call API to delete account
      const response = await DeleteProfile();
      console.log("Account deleted successfully:", response);
      
      // Clear session storage
      sessionStorage.clear();
      
      // Close modal
      setShowDeleteModal(false);
      
      // Redirect to login/signup page
      navigate("/");
      
    } catch (error) {
      console.error("Error deleting account:", error);
      alert("Failed to delete account. Please try again.");
    } finally {
      setDeletingAccount(false);
    }
  };

  return (
    <section className="bg-gray-50 min-h-screen py-8 px-4 sm:px-6 lg:px-8">
      {showDeleteModal && (
        <div className="fixed inset-0  bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-md mx-4">
            <div className="flex justify-between items-center border-b border-gray-200 px-6 py-4">
              <h3 className="text-lg font-medium text-gray-900">
                Confirm Account Deletion
              </h3>
              <button
                onClick={() => setShowDeleteModal(false)}
                className="text-gray-400 hover:text-gray-500"
              >
                <FiX className="h-6 w-6" />
              </button>
            </div>
            <div className="px-6 py-4">
              <p className="text-gray-700">
                Are you sure you want to delete your account? This action cannot
                be undone. All your data will be permanently removed.
              </p>
            </div>
            <div className="bg-gray-50 px-6 py-3 flex justify-end border-t border-gray-200">
              <button
                onClick={() => setShowDeleteModal(false)}
                className="mr-3 px-4 py-2 text-sm text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                disabled={deletingAccount}
              >
                Cancel
              </button>
              <button
                onClick={handleDeleteAccount}
                disabled={deletingAccount}
                className="px-4 py-2 text-sm text-white bg-red-600 border border-transparent rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {deletingAccount ? "Deleting..." : "Delete Account"}
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="max-w-3xl mx-auto">
        <div className="bg-white shadow-sm rounded-lg mb-6 p-6">
          <div className="flex items-center space-x-4">
            <div className="flex-shrink-0">
              <div className="w-16 h-16 rounded-full bg-stone-100 flex items-center justify-center">
                <span className="text-stone-600 text-2xl font-medium">
                  {userDetails.name.charAt(0)}
                </span>
              </div>
            </div>
            <div>
              <h1 className="text-2xl font-bold text-gray-800">
                {userDetails.name}
              </h1>
              <p className="text-gray-600">{userDetails.organization}</p>
            </div>
          </div>
        </div>

        {/* Account Information Section */}
        <div className="bg-white shadow-sm rounded-lg mb-6 overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-lg font-medium text-gray-800">
              Account Information
            </h2>
          </div>

          <div className="divide-y divide-gray-200">
            {/* Email Field */}
            <div className="px-6 py-4 flex flex-col md:flex-row">
              <div className="md:w-1/3 mb-2 md:mb-0">
                <label className="block text-sm font-medium text-gray-700">
                  Email
                </label>
                <p className="text-xs text-gray-500 mt-1">
                  Your primary email address
                </p>
              </div>
              <div className="md:w-2/3 flex items-center space-x-2">
                <div className="relative rounded-md shadow-sm flex-grow">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <FiMail className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    type="email"
                    className="block w-full pl-10 sm:text-sm focus:outline-none border-gray-300 rounded-md p-2 border bg-gray-50"
                    value={userDetails.email}
                    readOnly
                  />
                </div>
                <button className="inline-flex items-center px-3 py-1.5 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none">
                  <FiEdit2 className="mr-1 h-4 w-4" />
                  Update
                </button>
              </div>
            </div>

            {/* Name Field */}
            <div className="px-6 py-4 flex flex-col md:flex-row">
              <div className="md:w-1/3 mb-2 md:mb-0">
                <label className="block text-sm font-medium text-gray-700">
                  Full Name
                </label>
                <p className="text-xs text-gray-500 mt-1">Your display name</p>
              </div>
              <div className="md:w-2/3 flex items-center space-x-2">
                <div className="relative rounded-md shadow-sm flex-grow">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <FiUser className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    type="text"
                    className="block w-full pl-10 focus:outline-none sm:text-sm border-gray-300 rounded-md p-2 border bg-gray-50"
                    value={userDetails.name}
                    readOnly
                  />
                </div>
                <button className="inline-flex items-center px-3 py-1.5 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50">
                  <FiEdit2 className="mr-1 h-4 w-4" />
                  Update
                </button>
              </div>
            </div>

            {/* Organization Field */}
            <div className="px-6 py-4 flex flex-col md:flex-row">
              <div className="md:w-1/3 mb-2 md:mb-0">
                <label className="block text-sm font-medium text-gray-700">
                  Organization
                </label>
                <p className="text-xs text-gray-500 mt-1">
                  Your company or team
                </p>
              </div>
              <div className="md:w-2/3 flex items-center space-x-2">
                <div className="relative rounded-md shadow-sm flex-grow">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <CgOrganisation className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    type="text"
                    className="block w-full pl-10 focus:outline-none sm:text-sm border-gray-300 rounded-md p-2 border bg-gray-50"
                    value={userDetails.organization}
                    readOnly
                  />
                </div>
                <button className="inline-flex items-center px-3 py-1.5 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 ">
                  <FiEdit2 className="mr-1 h-4 w-4" />
                  Update
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Password Section */}
        <div className="bg-white shadow-sm rounded-lg mb-6 overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-lg font-medium text-gray-800">
              Change Password
            </h2>
          </div>

          <div className="px-6 py-4">
            <div className="flex flex-col sm:flex-row sm:items-end sm:space-x-4">
              <div className="flex-grow">
                <label
                  htmlFor="new-password"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  New Password
                </label>
                <div className="relative rounded-md shadow-sm">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <FiLock className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    id="new-password"
                    type={showPassword ? "text" : "password"}
                    className="block w-full pl-10 pr-10 focus:outline-none sm:text-sm border-gray-300 rounded-md p-2 border"
                    placeholder="Enter new password"
                    value={password}
                    onChange={handlePasswordChange}
                  />
                  <button
                    type="button"
                    className="absolute inset-y-0 right-0 pr-3 flex items-center"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? (
                      <FiEyeOff className="h-5 w-5 text-gray-400" />
                    ) : (
                      <FiEye className="h-5 w-5 text-gray-400" />
                    )}
                  </button>
                </div>
              </div>
              <button
                className="flex items-center gap-2 justify-center border align-middle select-none font-sans font-medium text-center duration-300 ease-in disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed focus:shadow-none text-sm py-2 px-4 shadow-sm hover:shadow-md bg-stone-800 hover:bg-stone-700 border-stone-900 text-stone-50 rounded-lg transition antialiased cursor-pointer mt-4 sm:mt-0"
                disabled={!isPasswordValid}
              >
                <FiCheck className="mr-2 h-4 w-4" />
                Update
              </button>
            </div>

            {/* Password Requirements */}
            <div className="mt-4 p-3 bg-gray-50 rounded-md">
              <p className="text-sm font-medium text-gray-700 mb-2">
                Password Requirements:
              </p>
              <ul className="text-sm text-gray-600 space-y-1">
                <li
                  className={`flex items-center ${
                    passwordRequirements.hasMinLength ? "text-green-600" : ""
                  }`}
                >
                  <FiCheck
                    className={`mr-2 h-4 w-4 ${
                      passwordRequirements.hasMinLength
                        ? "text-green-600"
                        : "text-gray-400"
                    }`}
                  />
                  At least 8 characters
                </li>
                <li
                  className={`flex items-center ${
                    passwordRequirements.hasUppercase ? "text-green-600" : ""
                  }`}
                >
                  <FiCheck
                    className={`mr-2 h-4 w-4 ${
                      passwordRequirements.hasUppercase
                        ? "text-green-600"
                        : "text-gray-400"
                    }`}
                  />
                  One uppercase letter
                </li>
                <li
                  className={`flex items-center ${
                    passwordRequirements.hasLowercase ? "text-green-600" : ""
                  }`}
                >
                  <FiCheck
                    className={`mr-2 h-4 w-4 ${
                      passwordRequirements.hasLowercase
                        ? "text-green-600"
                        : "text-gray-400"
                    }`}
                  />
                  One lowercase letter
                </li>
                <li
                  className={`flex items-center ${
                    passwordRequirements.hasNumber ? "text-green-600" : ""
                  }`}
                >
                  <FiCheck
                    className={`mr-2 h-4 w-4 ${
                      passwordRequirements.hasNumber
                        ? "text-green-600"
                        : "text-gray-400"
                    }`}
                  />
                  One number
                </li>
                <li
                  className={`flex items-center ${
                    passwordRequirements.hasSpecialChar ? "text-green-600" : ""
                  }`}
                >
                  <FiCheck
                    className={`mr-2 h-4 w-4 ${
                      passwordRequirements.hasSpecialChar
                        ? "text-green-600"
                        : "text-gray-400"
                    }`}
                  />
                  One special character
                </li>
              </ul>
            </div>
          </div>
        </div>

        {/* Danger Zone */}
        <div className="bg-white shadow-sm rounded-lg overflow-hidden border border-red-100">
          <div className="px-6 py-4 bg-red-50 border-b border-red-100">
            <h2 className="text-lg font-medium text-red-800">Danger Zone</h2>
          </div>

          <div className="px-6 py-4 flex justify-between items-center">
            <div>
              <h3 className="text-sm font-medium text-gray-800">
                Delete Account
              </h3>
              <p className="text-sm text-gray-500 mt-1">
                Once you delete your account, there is no going back. Please be
                certain.
              </p>
            </div>
            <button
              onClick={() => setShowDeleteModal(true)}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none cursor-pointer"
            >
              <FiTrash2 className="mr-2 h-4 w-4" />
              Delete Account
            </button>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Profile;
