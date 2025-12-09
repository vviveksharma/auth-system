import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
const Navbar = () => {
  const [isOpen, setIsOpen] = useState(false);
  let navigate = useNavigate();
  return (
    <nav className="w-full shadow-lg overflow-hidden p-2 border-transparent shadow-transparent rounded-none py-3 border-0 bg-transparent">
      <div className="container mx-auto">
        <div className="flex items-center relative">
          <button
            onClick={() => {navigate("/")}}
            className="font-sans cursor-pointer antialiased text-base md:text-lg text-stone-800 block py-1 font-semibold"
          >
            TrustKit
          </button>
          <div className="hidden lg:block absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2">
            <ul className="mt-4 flex flex-col gap-x-6 gap-y-1.5 lg:mt-0 lg:flex-row lg:items-center">
              <li>
                <button
                  onClick={() => (navigate("/about"))}
                  className="font-sans cursor-pointer antialiased text-base text-stone-800 p-1"
                >
                  About Us
                </button>
              </li>
              <li>
                <button
                  onClick={() => {navigate("/contact")}}
                  className="font-sans antialiased text-base text-stone-800 p-1 cursor-pointer"
                >
                  Contact Us
                </button>
              </li>
            </ul>
          </div>
          {/* Sign In Button */}
          <button
            onClick={() => {
              sessionStorage.setItem("button", "login");
              navigate("/auth");
            }}
            className="inline-flex items-center justify-center border align-middle select-none font-sans font-medium text-center transition-all ease-in disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed focus:shadow-none text-sm py-2 px-4 shadow-sm bg-transparent relative text-stone-700 hover:text-stone-700 border-stone-500 hover:bg-transparent duration-150 hover:border-stone-600 rounded-lg hover:opacity-60 hover:shadow-none lg:ml-auto lg:inline-block cursor-pointer"
          >
            Login
          </button>
          {/* Components Dropdown for Mobile */}
          <button
            onClick={() => setIsOpen(!isOpen)}
            className="inline-grid place-items-center border font-sans font-medium text-center transition-all duration-300 ease-in disabled:opacity-50 disabled:shadow-none disabled:pointer-events-none text-sm min-w-[38px] min-h-[38px] rounded-md shadow-sm hover:shadow-md bg-stone-800 hover:bg-stone-700 border-stone-900 text-stone-50 ml-auto lg:hidden"
          >
            <svg
              width="1.5em"
              height="1.5em"
              strokeWidth="1.5"
              viewBox="0 0 24 24"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              color="currentColor"
              className="h-5 w-5"
            >
              <path
                d="M3 5H21"
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
              ></path>
              <path
                d="M3 12H21"
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
              ></path>
              <path
                d="M3 19H21"
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
              ></path>
            </svg>
          </button>
        </div>
        {/* Mobile Menu (collapsed by default) */}
        <div
          className={`overflow-hidden transition-[max-height] duration-300 ease-in-out ${
            isOpen ? "max-h-40" : "max-h-0"
          }`}
          id="collapseList"
        >
          <ul className="flex flex-col gap-0.5 min-w-60">
            <li className="flex items-center cursor-pointer py-1.5 px-2.5 rounded-md font-sans transition-all duration-300 ease-in bg-transparent text-stone-600 hover:text-stone-800 dark:text-stone-300 dark:hover:text-white hover:bg-stone-200 dark:hover:bg-stone-700 focus:bg-stone-200 dark:focus:bg-stone-700 focus:text-stone-800 dark:focus:text-white">
              <a
                href="#"
                className="font-sans antialiased text-base text-stone-800 p-1"
              >
                About Us
              </a>
            </li>
            
            <li className="flex items-center cursor-pointer py-1.5 px-2.5 rounded-md font-sans transition-all duration-300 ease-in bg-transparent text-stone-600 hover:text-stone-800 dark:text-stone-300 dark:hover:text-white hover:bg-stone-200 dark:hover:bg-stone-700 focus:bg-stone-200 dark:focus:bg-stone-700 focus:text-stone-800 dark:focus:text-white">
              <a
                href="#"
                className="font-sans antialiased text-base text-stone-800 p-1"
              >
                Contact Us
              </a>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
