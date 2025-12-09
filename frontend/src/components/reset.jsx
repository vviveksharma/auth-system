import { useNavigate } from "react-router-dom";

const Reset = () => {
  let navigate = useNavigate();
  return (
    <>
      <div class="grid place-items-center min-w-screen min-h-screen p-4">
        <div class="w-full max-w-lg lg:max-w-md mx-auto p-6 sm:p-16 lg:p-0">
          <h2 class="font-sans antialiased text-stone-800 font-bold text-xl md:text-2xl lg:text-3xl mb-2">
            Reset Password
          </h2>
          <p class="font-sans antialiased text-base md:text-lg text-stone-600">
            You will receive an e-mail please follow the directions given
          </p>

          <form class="mt-8">
            <div class="mb-6 space-y-1.5">
              <label
                for="email"
                class="block mb-2 text-sm font-semibold antialiased text-stone-800"
              >
                Email
              </label>
              <div class="relative w-full">
                <input
                  placeholder="someone@example.com"
                  type="email"
                  class="w-full aria-disabled:cursor-not-allowed outline-none focus:outline-none text-stone-800  placeholder:text-stone-600/60 ring-transparent border border-stone-200 transition-all ease-in disabled:opacity-50 disabled:pointer-events-none select-none text-sm py-2 px-2.5 ring shadow-sm bg-white rounded-lg duration-100 hover:border-stone-300 hover:ring-none focus:border-stone-400 focus:ring-none peer"
                />
              </div>
            </div>
            <div class="flex items-center flex-wrap gap-4 mb-6 justify-center">
              <div class="flex items-center gap-2 text-sm">
                Remember Password ?{" "}
              </div>
              <button
                onClick={() => {
                  sessionStorage.setItem("button", "login");
                  navigate("/auth");
                }}
                class="font-sans antialiased text-sm text-stone-800 font-semibold cursor-pointer"
              >
                Login
              </button>
            </div>

            <button class="inline-flex items-center justify-center border align-middle select-none font-sans font-medium text-center duration-300 ease-in disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed focus:shadow-none text-sm py-2 px-4 shadow-sm hover:shadow-md bg-stone-800 hover:bg-stone-700 border-stone-900 text-stone-50 rounded-lg transition antialiased w-full cursor-pointer">
              Send Email
            </button>
          </form>
        </div>
      </div>
    </>
  );
};

export default Reset;
