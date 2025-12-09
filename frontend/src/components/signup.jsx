import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { CreateTenant, LoginTenant } from "./services/auth";

const SignUp = () => {
  let navigate = useNavigate();
  const [isShow, setIsShow] = useState(false);
  const [loading, setLoading] = useState(false);
  const [formInputs, setFormInputs] = useState({
    email: "",
    password: "",
    name: "",
    organization: ""
  });
  let route = sessionStorage.getItem("button");
  const formData = {
    heading: "Login",
    subheading: "Enter your email and password to login",
    buttonDetails: "LogIn",
    callheading: "Not registered ?",
    buttoncall: "Create Account",
  };
  useEffect(() => {
    if (route === "register") {
      setIsShow(true);
    }
  }, [route]);
  if (route === "register") {
    formData.buttonDetails = "SignIn";
    formData.heading = "SignIn";
    formData.subheading = "Plug in, power up—secure access made simple.";
    formData.callheading = "Already a user ?";
    formData.buttoncall = "Login";
  }
  function pageReload() {
    location.reload();
  }

  const handleInputChange = (e) => {
    setFormInputs({
      ...formInputs,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    
    try {
      if (route === "register") {
        // Create new tenant
        const response = await CreateTenant(
          formInputs.email,
          formInputs.password,
          formInputs.name,
          formInputs.organization
        );
        console.log("Tenant created successfully:", response);
        // Navigate to login or dashboard
        sessionStorage.setItem("button", "login");
        pageReload();
      } else {
        // Login existing tenant
        const response = await LoginTenant(
          formInputs.email,
          formInputs.password
        );
        console.log("Login successful:", response.data.token);
        // Store auth token if needed
        if (response.data.token) {
          sessionStorage.setItem("authToken", response.data.token);
        }
        // Navigate to dashboard
        navigate("/home");
      }
    } catch (error) {
      console.error("Authentication error:", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="antialiased">
      <div className="container px-6 mx-auto">
        <div className="flex flex-col text-center md:text-left md:flex-row h-screen justify-evenly md:items-center">
          <div className="flex flex-col w-full">
            <div>
              <svg
                className="w-20 h-20 mx-auto md:float-left fill-stroke text-gray-800"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                ></path>
              </svg>
            </div>
            <h1 className="text-5xl text-gray-800 font-bold">TrustKit</h1>
            <p className="w-full mx-auto md:mx-0 text-gray-500">
              Your auth, your rules—managed from one place.
            </p>
          </div>
          <div className="w-full md:w-full lg:w-9/12 mx-auto md:mx-0">
            <div className="grid place-items-center min-w-screen min-h-screen p-4">
              <div className="w-full max-w-lg lg:max-w-md mx-auto p-6 sm:p-16 lg:p-0">
                <h2 className="font-sans antialiased text-stone-800 font-bold text-xl md:text-2xl lg:text-3xl  mb-2">
                  {formData.heading}
                </h2>
                <p className="font-sans antialiased text-base md:text-lg text-stone-600">
                  {formData.subheading}
                </p>

                <form className="mt-8" onSubmit={handleSubmit}>
                  <div className="mb-6 space-y-1.5">
                    <label className="block mb-2 text-sm font-semibold antialiased text-stone-800">
                      Email
                    </label>
                    <div className="relative w-full">
                      <input
                        name="email"
                        value={formInputs.email}
                        onChange={handleInputChange}
                        placeholder="someone@example.com"
                        type="email"
                        required
                        className="w-full aria-disabled:cursor-not-allowed outline-none focus:outline-none text-stone-800  placeholder:text-stone-600/60 ring-transparent border border-stone-200 transition-all ease-in disabled:opacity-50 disabled:pointer-events-none select-none text-sm py-2 px-2.5 ring shadow-sm bg-white rounded-lg duration-100 hover:border-stone-300 hover:ring-none focus:border-stone-400 focus:ring-none peer"
                      />
                    </div>
                  </div>

                  <div className="mb-6 space-y-1.5">
                    <label className="block mb-2 text-sm font-semibold antialiased text-stone-800">
                      Password
                    </label>
                    <div className="relative w-full">
                      <input
                        id="current-password"
                        name="password"
                        value={formInputs.password}
                        onChange={handleInputChange}
                        placeholder="enter your password"
                        type="password"
                        required
                        className="w-full aria-disabled:cursor-not-allowed outline-none focus:outline-none text-stone-800 placeholder:text-stone-600/60 ring-transparent border border-stone-200 transition-all ease-in disabled:opacity-50 disabled:pointer-events-none select-none text-sm py-2 px-2.5 ring shadow-sm bg-white rounded-lg duration-100 hover:border-stone-300 hover:ring-none focus:border-stone-400 focus:ring-none peer"
                      />
                    </div>
                  </div>
                  {route === "register" ? (
                    <>
                      <div className="mb-6 space-y-1.5">
                        <label className="block mb-2 text-sm font-semibold antialiased text-stone-800">
                          Username
                        </label>
                        <div className="relative w-full">
                          <input
                            id="username"
                            name="name"
                            value={formInputs.name}
                            onChange={handleInputChange}
                            placeholder="enter your username"
                            type="text"
                            required
                            className="w-full aria-disabled:cursor-not-allowed outline-none focus:outline-none text-stone-800 placeholder:text-stone-600/60 ring-transparent border border-stone-200 transition-all ease-in disabled:opacity-50 disabled:pointer-events-none select-none text-sm py-2 px-2.5 ring shadow-sm bg-white rounded-lg duration-100 hover:border-stone-300 hover:ring-none focus:border-stone-400 focus:ring-none peer"
                          />
                        </div>
                      </div>
                      <div className="mb-6 space-y-1.5">
                        <label className="block mb-2 text-sm font-semibold antialiased text-stone-800">
                          Organisation
                        </label>
                        <div className="relative w-full">
                          <input
                            id="organisation"
                            name="organization"
                            value={formInputs.organization}
                            onChange={handleInputChange}
                            placeholder="enter your organisation"
                            type="text"
                            required
                            className="w-full aria-disabled:cursor-not-allowed outline-none focus:outline-none text-stone-800 placeholder:text-stone-600/60 ring-transparent border border-stone-200 transition-all ease-in disabled:opacity-50 disabled:pointer-events-none select-none text-sm py-2 px-2.5 ring shadow-sm bg-white rounded-lg duration-100 hover:border-stone-300 hover:ring-none focus:border-stone-400 focus:ring-none peer"
                          />
                        </div>
                      </div>
                    </>
                  ) : (
                    <></>
                  )}

                  <div className="flex items-center flex-wrap gap-4 mb-6 justify-between">
                    <div className="flex items-center gap-2">
                      <div className="inline-flex items-center">
                        <label
                          className="flex items-center cursor-pointer relative"
                          htmlFor="check-2"
                        >
                          <input
                            type="checkbox"
                            className="peer h-5 w-5 cursor-pointer transition-all appearance-none rounded shadow-sm hover:shadow border border-stone-200 checked:bg-stone-800 checked:border-stone-800"
                            id="check-2"
                          />
                          <span className="absolute text-white opacity-0 peer-checked:opacity-100 top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
                            <svg
                              strokeWidth="1.5"
                              className="w-3.5 h-3.5"
                              viewBox="0 0 24 24"
                              fill="none"
                              xmlns="http://www.w3.org/2000/svg"
                              color="#ffffff"
                            >
                              <path
                                d="M5 13L9 17L19 7"
                                stroke="#ffffff"
                                strokeWidth="1.5"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                              ></path>
                            </svg>
                          </span>
                        </label>
                        <label
                          className="cursor-pointer ml-2 text-stone-600 text-sm"
                          htmlFor="check-2"
                        >
                          Remember Me
                        </label>
                      </div>
                    </div>
                    {route === "login" ? (
                      <button 
                      onClick={()=>{navigate("/auth/reset")}}
                      className="cursor-pointer font-sans antialiased text-sm text-stone-800 font-semibold">
                        Forgot Password?
                      </button>
                    ) : (
                      <></>
                    )}
                  </div>

                  <button 
                    type="submit"
                    disabled={loading}
                    className="inline-flex items-center justify-center border align-middle select-none font-sans font-medium text-center duration-300 ease-in disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed focus:shadow-none text-sm py-2 px-4 shadow-sm hover:shadow-md bg-stone-800 hover:bg-stone-700 border-stone-900 text-stone-50 rounded-lg transition antialiased w-full"
                  >
                    {loading ? "Loading..." : formData.buttonDetails}
                  </button>
                </form>

                <p className="font-sans antialiased flex items-center justify-center gap-1 text-stone-600 text-sm mt-2">
                  {formData.callheading}
                  <button
                    onClick={() => {
                      if (route === "login") {
                        sessionStorage.setItem("button","register");
                      } else {
                        sessionStorage.setItem("button","login");
                      }
                      pageReload();
                    }}
                    className="font-sans antialiased text-stone-800 font-semibold text-sm cursor-pointer"
                  >
                    {formData.buttoncall}
                  </button>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SignUp;
