import React from "react";
import {
  FiUsers,
  FiShield,
  FiActivity,
  FiCode,
  FiCheckCircle,
} from "react-icons/fi";
import { FaChartLine, FaServer, FaLock } from "react-icons/fa";
import { useNavigate } from "react-router-dom";

const About = () => {
  let navigate = useNavigate();
  // Features data
  const features = [
    {
      icon: <FaLock className="h-8 w-8 text-stone-800" />,
      title: "Secure Authentication",
      description:
        "Industry-standard JWT tokens with customizable expiration and refresh mechanisms.",
    },
    {
      icon: <FiUsers className="h-8 w-8 text-stone-800" />,
      title: "User Management",
      description:
        "Full CRUD operations for user accounts with detailed activity tracking.",
    },
    {
      icon: <FiShield className="h-8 w-8 text-stone-800" />,
      title: "Role-Based Access",
      description:
        "Flexible role and permission system to control access to your application resources.",
    },
    {
      icon: <FiActivity className="h-8 w-8 text-stone-800" />,
      title: "Activity Monitoring",
      description:
        "Real-time tracking of login/logout events and user sessions.",
    },
    {
      icon: <FaServer className="h-8 w-8 text-stone-800" />,
      title: "High Availability",
      description:
        "99.9% uptime with Redis caching for optimal performance at scale.",
    },
    {
      icon: <FiCode className="h-8 w-8 text-stone-800" />,
      title: "Easy Integration",
      description:
        "RESTful API with comprehensive documentation and client libraries.",
    },
  ];

  // Stats data
  const stats = [
    { value: "2+", label: "Years in Production" },
    { value: "10M+", label: "Daily Authentications" },
    { value: "99.9%", label: "Uptime" },
    { value: "24/7", label: "Support" },
  ];

  // How it works steps
  const steps = [
    {
      step: "1",
      title: "Register Your Application",
      description:
        "Get your API keys and configure your authentication settings.",
    },
    {
      step: "2",
      title: "Integrate Our API",
      description: "Use our direct Rest API calls to implement auth flows.",
    },
    {
      step: "3",
      title: "Manage Users & Roles",
      description:
        "Create users, assign roles, and set permissions through our dashboard.",
    },
    {
      step: "4",
      title: "Monitor Activity",
      description:
        "Track login attempts, sessions, and security events in real-time.",
    },
  ];

  return (
    <div className="bg-gray-50">
      <section className="relative bg-stone-800 text-white py-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto text-center">
          <h1 className="text-4xl md:text-5xl font-bold mb-6">
            Powering Secure Authentication
          </h1>
          <p className="text-xl max-w-3xl mx-auto opacity-90">
            We provide developers with enterprise-grade authentication solutions
            that scale with your application.
          </p>
          <div className="mt-10 flex flex-col sm:flex-row justify-center gap-4">
            <button className="px-6 py-3 bg-white cursor-pointer text-stone-800 font-medium rounded-md hover:bg-gray-100 transition-colors">
              View Documentation
            </button>
            <button
              onClick={() => {
                navigate("/contact");
              }}
              className="px-6 py-3 border cursor-pointer border-white text-white font-medium rounded-md hover:bg-white hover:text-black hover:bg-opacity-10 transition-colors"
            >
              Contact Sales
            </button>
          </div>
        </div>
      </section>
      {/* Features Section */}
      <section className="py-16 px-4 sm:px-6 lg:px-8 bg-white">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              Comprehensive Authentication Features
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              Everything you need to implement secure authentication in your
              application.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <div
                key={index}
                className="bg-gray-50 p-6 rounded-lg hover:shadow-md transition-shadow"
              >
                <div className="flex items-center mb-4">
                  <div className="mr-4">{feature.icon}</div>
                  <h3 className="text-xl font-semibold text-gray-900">
                    {feature.title}
                  </h3>
                </div>
                <p className="text-gray-600">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* How It Works Section */}
      <section className="py-16 px-4 sm:px-6 lg:px-8 bg-gray-50">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              Simple Integration
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              Get up and running with our authentication service in just a few
              steps.
            </p>
          </div>

          <div className="space-y-10">
            {steps.map((step, index) => (
              <div
                key={index}
                className="flex flex-col md:flex-row items-start"
              >
                <div className="flex items-center justify-center bg-stone-800 text-white rounded-full h-12 w-12 text-xl font-bold mb-4 md:mb-0 md:mr-8 flex-shrink-0">
                  {step.step}
                </div>
                <div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">
                    {step.title}
                  </h3>
                  <p className="text-gray-600">{step.description}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Code Example Section */}
      <section className="py-16 px-4 sm:px-6 lg:px-8 bg-white">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              Quick Start
            </h2>
            <p className="text-xl text-gray-600 max-w-3xl mx-auto">
              Here's how easy it is to implement user registration with our API.
            </p>
          </div>

          <div className="bg-gray-800 rounded-lg overflow-hidden">
            <div className="px-6 py-4 border-b border-gray-700 flex items-center">
              <div className="flex space-x-2">
                <div className="w-3 h-3 rounded-full bg-red-500"></div>
                <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div className="w-3 h-3 rounded-full bg-green-500"></div>
              </div>
              <div className="ml-4 text-sm text-gray-400">
                JavaScript Example
              </div>
            </div>
            <div className="p-6">
              <pre className="text-gray-300 overflow-x-auto">
                <code>
                  {`// Register a new user
const registerUser = async () => {
  try {
    const response = await fetch('https://api.yourservice.com/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'API-Key': 'your_api_key_here'
      },
      body: JSON.stringify({
        email: 'user@example.com',
        password: 'securePassword123',
        name: 'John Doe'
      })
    });

    const data = await response.json();
    console.log('Registration successful:', data);
  } catch (error) {
    console.error('Registration failed:', error);
  }
};

registerUser();`}
                </code>
              </pre>
            </div>
          </div>

          <div className="mt-8 text-center">
            <a
              href="/docs"
              className="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-stone-800 hover:bg-stone-700"
            >
              View Full API Documentation
              <FiCode className="ml-2" />
            </a>
          </div>
        </div>
      </section>
    </div>
  );
};

export default About;
