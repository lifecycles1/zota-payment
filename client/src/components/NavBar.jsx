import { Link } from "react-router-dom";

const NavBar = () => {
  return (
    <nav className="bg-blue-200">
      <div className="container mx-auto px-6 py-3">
        <div className="flex items-center justify-between">
          <div className="w-full text-gray-600 md:flex md:items-center">
            <Link to="/" className="text-sm text-gray-600 underline">
              Home
            </Link>
            <Link to="/v1/account/deposit/:uid" className="text-sm text-gray-600 underline ml-6">
              Profile
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default NavBar;
