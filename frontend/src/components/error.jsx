import { useNavigate } from "react-router-dom";
import { FiHome, FiArrowLeft, FiAlertTriangle } from "react-icons/fi";

const NotFound = () => {
  let navigate = useNavigate();
  
  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto flex flex-col lg:flex-row items-center gap-12 lg:gap-16">
        {/* Illustration Section */}
        <div className="flex-1 text-center lg:text-left">
          <div className="relative">
            <div className="mb-8 lg:mb-12">
              <div className="inline-flex items-center justify-center w-16 h-16 bg-stone-100 rounded-full mb-6">
                <FiAlertTriangle className="w-8 h-8 text-stone-600" />
              </div>
              <h1 className="text-6xl font-bold text-stone-800 mb-4">404</h1>
              <h2 className="text-2xl font-semibold text-stone-700 mb-4">
                Oops! Page not found
              </h2>
              <p className="text-stone-600 mb-8 text-lg">
                Looks like you've found the doorway to the great nothing. 
                Sorry about that! Please visit our homepage to get where you need to go.
              </p>
              
              <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
                <button 
                  onClick={() => navigate(-1)}
                  className="flex items-center justify-center px-6 py-3 border border-stone-300 text-stone-700 font-medium rounded-md hover:bg-stone-50 transition-colors"
                >
                  <FiArrowLeft className="mr-2" />
                  Go Back
                </button>
                <button 
                  onClick={() => navigate("/")}
                  className="flex items-center justify-center px-6 py-3 bg-stone-800 text-white font-medium rounded-md hover:bg-stone-700 transition-colors"
                >
                  <FiHome className="mr-2" />
                  Take Me Home
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Image Section */}
        <div className="flex-1">
          <div className="relative">
            <img 
              src="https://i.ibb.co/G9DC8S0/404-2.png" 
              alt="404 Illustration" 
              className="w-full max-w-md mx-auto lg:max-w-lg"
            />
            {/* Decorative elements */}
            <div className="absolute -top-4 -right-4 w-24 h-24 bg-stone-100 rounded-full opacity-50"></div>
            <div className="absolute -bottom-4 -left-4 w-16 h-16 bg-stone-200 rounded-full opacity-30"></div>
          </div>
        </div>
      </div>

      {/* Background decorative elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/4 left-10 w-32 h-32 bg-stone-100 rounded-full opacity-20"></div>
        <div className="absolute bottom-1/4 right-10 w-40 h-40 bg-stone-200 rounded-full opacity-10"></div>
        <div className="absolute top-10 right-1/4 w-20 h-20 bg-stone-100 rounded-full opacity-15"></div>
      </div>
    </div>
  );
};

export default NotFound;