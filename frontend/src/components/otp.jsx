import { useNavigate } from "react-router-dom";
import React from "react";
const OTP = () => {
  let navigate = useNavigate();
  const [showNewPassword, setShowNewPassword] = React.useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = React.useState(false);
  const [newPassword, setNewPassword] = React.useState("");
  const [confirmPassword, setConfirmPassword] = React.useState("");

  const passwordsMatch = newPassword === confirmPassword && newPassword !== "";

  return (
    <>
      <div className="grid place-items-center min-w-screen min-h-screen p-4">
        <div className="w-full max-w-lg lg:max-w-md mx-auto p-6 sm:p-16 lg:p-0">
          <h2 className="font-sans antialiased text-stone-800 font-bold text-xl md:text-2xl lg:text-3xl mb-2">
            Set New Password
          </h2>

          <form className="mt-8">
            <div className="mb-6 space-y-1.5">
              <label className="block mb-2 text-sm font-semibold antialiased text-stone-800">
                New Password
              </label>
              <div className="relative w-full">
                <input
                  type={showNewPassword ? "text" : "password"}
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  className="w-full aria-disabled:cursor-not-allowed outline-none focus:outline-none text-stone-800  placeholder:text-stone-600/60 ring-transparent border border-stone-200 transition-all ease-in disabled:opacity-50 disabled:pointer-events-none select-none text-sm py-2 px-2.5 ring shadow-sm bg-white rounded-lg duration-100 hover:border-stone-300 hover:ring-none focus:border-stone-400 focus:ring-none peer"
                />
                <button
                  type="button"
                  tabIndex={-1}
                  className="absolute right-2 top-2 text-stone-500 hover:text-stone-700"
                  onClick={() => setShowNewPassword((prev) => !prev)}
                  aria-label={
                    showNewPassword ? "Hide password" : "Show password"
                  }
                >
                  {showNewPassword ? "👁️" : "🙈"}
                </button>
              </div>
            </div>
            <div className="mb-6 space-y-1.5">
              <label className="block mb-2 text-sm font-semibold antialiased text-stone-800">
                Confirm Password
              </label>
              <div className="relative w-full">
                <input
                  type={showConfirmPassword ? "text" : "password"}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className="w-full aria-disabled:cursor-not-allowed outline-none focus:outline-none text-stone-800  placeholder:text-stone-600/60 ring-transparent border border-stone-200 transition-all ease-in disabled:opacity-50 disabled:pointer-events-none select-none text-sm py-2 px-2.5 ring shadow-sm bg-white rounded-lg duration-100 hover:border-stone-300 hover:ring-none focus:border-stone-400 focus:ring-none peer"
                />
                <button
                  type="button"
                  tabIndex={-1}
                  className="absolute right-2 top-2 text-stone-500 hover:text-stone-700"
                  onClick={() => setShowConfirmPassword((prev) => !prev)}
                  aria-label={
                    showConfirmPassword ? "Hide password" : "Show password"
                  }
                >
                  {showConfirmPassword ? "👁️" : "🙈"}
                </button>
              </div>
              {!passwordsMatch && confirmPassword && (
                <p className="text-red-500 text-xs mt-1">
                  Passwords do not match
                </p>
              )}
            </div>
            <button
              type="submit"
              onClick={()=>{navigate("/auth")}}
              disabled={!passwordsMatch}
              className={`inline-flex items-center justify-center border align-middle select-none font-sans font-medium text-center duration-300 ease-in focus:shadow-none text-sm py-2 px-4 shadow-sm hover:shadow-md border-stone-900 rounded-lg transition antialiased w-full ${
                passwordsMatch
                  ? "bg-stone-800 hover:bg-stone-700 cursor-pointer text-stone-50"
                  : "bg-stone-200 cursor-not-allowed text-black"
              }`}
            >
              Set Password
            </button>
          </form>
        </div>
      </div>
    </>
  );
};

export default OTP;