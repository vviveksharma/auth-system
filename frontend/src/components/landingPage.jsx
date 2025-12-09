import rolesUI from "../assets/roles.png";
import usersUI from "../assets/users.png";
import dashboardUI from "../assets/dashboard.png"
import { useNavigate } from "react-router-dom";
const LandingPage = () => {
  let navigate = useNavigate();
  return (
    <section className=" bg-white">
      <div className="flex flex-col px-8 mx-auto space-y-12 max-w-7xl xl:px-12">
        <div className="min-h-screen flex flex-col items-center justify-center">
          <h2 className="w-full text-3xl font-bold text-center sm:text-4xl md:text-5xl">
            No more auth headaches. Just ship.
          </h2>
          <p className="w-full py-8 mx-auto mt-2 text-lg text-center text-gray-700 intro sm:max-w-3xl">
            A plug-and-play authentication and authorization SDK with tenant
            support, JWT, RBAC, and more.
          </p>

          <button
            onClick={() => {
              sessionStorage.setItem("button", "register")
              navigate("/auth")
            }}
            className="inline-flex items-center justify-center border align-middle select-none font-sans font-medium text-center duration-300 ease-in disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed focus:shadow-none text-sm py-2 px-4 shadow-sm hover:shadow-md relative bg-gradient-to-b from-stone-700 to-stone-800 border-stone-900 text-stone-50 rounded-lg hover:bg-gradient-to-b hover:from-stone-800 hover:to-stone-800 hover:border-stone-900 after:absolute after:inset-0 after:rounded-[inherit] after:box-shadow after:shadow-[inset_0_1px_0px_rgba(255,255,255,0.25),inset_0_-2px_0px_rgba(0,0,0,0.35)] after:pointer-events-none transition antialiased cursor-pointer"
          >
            Get Started
          </button>
        </div>
        <div className="flex flex-col mb-8 animated fadeIn sm:flex-row">
          <div className="flex items-center mb-8 sm:w-1/2 md:w-5/12 sm:order-last">
            <img
              className="rounded-lg shadow-xl"
              src={dashboardUI}
              alt="Design Made Easy"
            />
          </div>
          <div className="flex flex-col justify-center mt-5 mb-8 md:mt-0 sm:w-1/2 md:w-7/12 sm:pr-16">
            <p className="mb-2 text-sm font-semibold leading-none text-left text-stone-800 uppercase">
              🚀 Launch Secure Login Faster Than Ever
            </p>
            <h3 className="mt-2 text-2xl sm:text-left md:text-4xl">
              Authentication in Minutes
            </h3>
            <p className="mt-5 text-lg text-gray-700 text md:text-left">
              No more writing login flows from scratch. Use our SDKs to
              integrate secure authentication in under 10 minutes.
            </p>
          </div>
        </div>
        <div className="flex flex-col mb-8 animated fadeIn sm:flex-row">
          <div className="flex items-center mb-8 sm:w-1/2 md:w-5/12">
            <img
              className="rounded-lg shadow-xl"
              src={usersUI}
              alt="Optimized For Conversions"
            />
          </div>
          <div className="flex flex-col justify-center mt-5 mb-8 md:mt-0 sm:w-1/2 md:w-7/12 sm:pl-16">
            <p className="mb-2 text-sm font-semibold leading-none text-left text-stone-800 uppercase">
              👥 Manage Users Like a Pro, No Backend Required
            </p>
            <h3 className="mt-2 text-2xl sm:text-left md:text-4xl">
              Effortless User Management
            </h3>
            <p className="mt-5 text-lg text-gray-700 text md:text-left">
              View, update, and delete users with a sleek dashboard. Built-in
              support for password resets, profile edits, and user search.
            </p>
          </div>
        </div>
        <div className="flex flex-col mb-8 animated fadeIn sm:flex-row">
          <div className="flex items-center mb-8 sm:w-1/2 md:w-5/12 sm:order-last">
            <img
              className="rounded-lg shadow-xl"
              src={rolesUI}
              alt="Make It Your Own"
            />
          </div>
          <div className="flex flex-col justify-center mt-5 mb-8 md:mt-0 sm:w-1/2 md:w-7/12 sm:pr-16">
            <p className="mb-2 text-sm font-semibold leading-none text-left text-stone-800 uppercase">
              🔐 Lock Down Access with Precision
            </p>
            <h3 className="mt-2 text-2xl sm:text-left md:text-4xl">
              Role-Based Access Control (RBAC)
            </h3>
            <p className="mt-5 text-lg text-gray-700 text md:text-left">
              Define roles and permissions for every API endpoint. Your data
              stays secure—only the right people see the right things.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};

export default LandingPage;
