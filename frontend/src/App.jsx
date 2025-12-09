import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import LandingPage from "./components/LandingPage.jsx";
import Navbar from "./components/navbar.jsx";
import Articles from "./components/articles.jsx";
import { useLocation } from "react-router-dom";
import SignUp from "./components/signup.jsx";
import Home from "./components/home.jsx";
import Tokens from "./components/homePages/tokens.jsx";
import Users from "./components/homePages/users.jsx";
import NotFound from "./components/error.jsx";
import Profile from "./components/homePages/profile.jsx";
import Roles from "./components/homePages/roles.jsx";
import Temp from "./components/homePages/temp.jsx";
import Reset from "./components/reset.jsx";
import OTP from "./components/otp.jsx";
import Dashboard from "./components/homePages/dashboard.jsx";
import Contact from "./components/contact.jsx";
import About from "./components/about.jsx";
export default function App() {
  return (
    <>
      <Router>
        <ConditionalNavbar />
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/auth" element={<SignUp />} />
          <Route path="/auth/reset" element={<Reset />} />
          <Route path="/auth/set" element={<OTP />} />
          <Route path="/home" element={<Home />} />
          <Route path="/home/tokens" element={<Tokens />} />
          <Route path="/home/users" element={<Users />} />
          <Route path="/home/profile" element={<Profile />} />
          <Route path="/home/dashboard" element={<Dashboard />} />
          <Route path="/about" element={<About />} />
          <Route path="/contact" element={<Contact />} />
          <Route path="/*" element={<NotFound />} />
          <Route path="/home/roles" element={<Roles />} />
          <Route path="/temp" element={<Temp />} />
        </Routes>
      </Router>
    </>
  );
}

function ConditionalNavbar() {
  const routes = ["/", "/contact", "/about"];
  const location = useLocation();
  if (routes.includes(location.pathname)) {
    return <Navbar />;
  } else if (!routes.includes(location.pathname)) {
    return null;
  }
}
