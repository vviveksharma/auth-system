import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { FiUsers, FiKey, FiUser, FiShield } from "react-icons/fi";
import { SlDocs } from "react-icons/sl";
import { IoLogOut } from "react-icons/io5";
import { RxDashboard } from "react-icons/rx";
import { TiMessages } from "react-icons/ti";
import Users from "./homePages/users";
import Tokens from "./homePages/tokens";
import Profile from "./homePages/profile";
import Roles from "./homePages/roles";
import logo from "../assets/logo.png";
import Messages from "./homePages/messages";
import Dashboard from "./homePages/dashboard";
import Docs from "./homePages/docs";

const Home = () => {
  const [activeSection, setActiveSection] = useState(
    () => sessionStorage.getItem("activeSection") || "dashboard"
  );
  const [sidebarOpen, setSidebarOpen] = useState(false);
  let navigate = useNavigate();

  useEffect(() => {
    sessionStorage.setItem("activeSection", activeSection);
  }, [activeSection]);

  const handleSectionClick = (section) => {
    setActiveSection(section);
    setSidebarOpen(false);
  };

  function pageReload() {
    location.reload();
  }

  return (
    <div className="min-h-screen flex flex-row">
      <button
        className="md:hidden fixed top-4 left-4 z-50 bg-white rounded-full p-2 shadow"
        onClick={() => setSidebarOpen(!sidebarOpen)}
        aria-label="Open sidebar"
      >
        <svg width="24" height="24" fill="none" stroke="currentColor">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M4 6h16M4 12h16M4 18h16"
          />
        </svg>
      </button>
      {/* Sidebar: collapsible only on small screens, static on md+ */}
      <div
        className={`bg-white rounded-r-3xl overflow-hidden w-56 flex flex-col z-40 transition-transform duration-300
          fixed top-0 left-0 h-full md:sticky md:top-0 md:h-screen md:translate-x-0
          ${sidebarOpen ? "translate-x-0" : "-translate-x-full"} md:flex`}
      >
        <div className="flex items-center justify-center h-20">
          <span
            onClick={() => {
              sessionStorage.setItem("activeSection", "dashboard");
              pageReload();
            }}
            className="font-sans cursor-pointer antialiased text-base md:text-lg text-stone-800 block py-1 font-semibold"
          >
            TrustKit
          </span>
        </div>
        <ul className="flex flex-col py-4">
          <li>
            <button
              onClick={() => handleSectionClick("dashboard")}
              className={`flex flex-row items-center h-12 px-4 gap-2 transform hover:translate-x-2 font-sans transition-transform ease-in duration-200 text-stone-800 hover:text-gray-800 w-full cursor-pointer ${
                activeSection === "dashboard" ? "" : ""
              }`}
            >
              <RxDashboard className="text-lg text-stone-800" />
              <span className="text-sm font-medium">Dashboard</span>
            </button>
          </li>
          <li>
            <button
              onClick={() => handleSectionClick("tokens")}
              className={`flex flex-row items-center h-12 px-4 gap-2 transform hover:translate-x-2 font-sans transition-transform ease-in duration-200 text-stone-800 hover:text-gray-800 w-full cursor-pointer ${
                activeSection === "tokens" ? "" : ""
              }`}
            >
              <FiKey className="text-lg text-stone-800" />
              <span className="text-sm font-medium">Tokens</span>
            </button>
          </li>
          <li>
            <button
              onClick={() => handleSectionClick("users")}
              className={`flex flex-row items-center h-12 px-4 gap-2 transform hover:translate-x-2 transition-transform ease-in duration-200 text-stone-800 hover:text-gray-800 w-full cursor-pointer font-sans ${
                activeSection === "users" ? "" : ""
              }`}
            >
              <FiUsers className="text-lg text-stone-800" />
              <span className="text-sm font-medium">User Management</span>
            </button>
          </li>
          <li>
            <button
              onClick={() => handleSectionClick("docs")}
              className={`flex flex-row items-center h-12 px-4 gap-2 transform hover:translate-x-2 font-sans transition-transform ease-in duration-200 text-stone-800 hover:text-gray-800 w-full cursor-pointer ${
                activeSection === "docs" ? "" : ""
              }`}
            >
              <SlDocs className="text-lg text-stone-800" />
              <span className="text-sm font-medium">Docs</span>
            </button>
          </li>
          <li>
            <button
              onClick={() => handleSectionClick("roles")}
              className={`flex flex-row items-center h-12 px-4 gap-2 transform hover:translate-x-2 font-sans transition-transform ease-in duration-200 text-stone-800 hover:text-gray-800 w-full cursor-pointer ${
                activeSection === "roles" ? "" : ""
              }`}
            >
              <FiShield className="text-lg text-stone-800" />
              <span className="text-sm font-medium">Roles</span>
            </button>
          </li>
          <li>
            <button
              onClick={() => handleSectionClick("messages")}
              className={`flex flex-row items-center h-12 px-4 gap-2 transform hover:translate-x-2 font-sans transition-transform ease-in duration-200 text-stone-800 hover:text-gray-800 w-full cursor-pointer ${
                activeSection === "messages" ? "" : ""
              }`}
            >
              <TiMessages className="text-lg text-stone-800" />
              <span className="text-sm font-medium">Messages</span>
            </button>
          </li>
          <li>
            <button
              onClick={() => handleSectionClick("profile")}
              className={`flex flex-row items-center h-12 px-4 gap-2 transform hover:translate-x-2 font-sans transition-transform ease-in duration-200 text-stone-800 hover:text-gray-800 w-full cursor-pointer ${
                activeSection === "profile" ? "" : ""
              }`}
            >
              <FiUser className="text-lg text-stone-800" />
              <span className="text-sm font-medium">Profile</span>
            </button>
          </li>
          <li>
            <button
              onClick={() => {
                sessionStorage.clear();
                navigate("/");
              }}
              className={`group flex flex-row items-center h-12 px-4 gap-2 transform hover:translate-x-2 transition-transform ease-in duration-200 text-red-400 hover:text-gray-800 w-full cursor-pointer ${
                activeSection === "logout" ? "" : ""
              }`}
            >
              <IoLogOut className="text-lg text-red-400 group-hover:text-gray-800 transition-colors" />
              <span className="text-sm font-medium group-hover:text-gray-800 transition-colors">
                Logout
              </span>
            </button>
          </li>
        </ul>
      </div>
      {/* Overlay for mobile when sidebar is open */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-30 z-30 md:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}
      {/* Main content */}
      <div className="flex-1 p-8">
        {activeSection === "tokens" && <Tokens />}
        {activeSection === "users" && <Users />}
        {activeSection === "profile" && <Profile />}
        {activeSection === "dashboard" && <Dashboard />}
        {activeSection === "roles" && <Roles />}
        {activeSection === "docs" && <Docs />}
        {activeSection === "messages" && <Messages />}
      </div>
    </div>
  );
};

export default Home;
