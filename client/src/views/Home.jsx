import { Link } from "react-router-dom";

const Home = () => {
  return (
    <div className="home-container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Home</h1>
      <p>Welcome to the home page.</p>
      <section className="bg-gray-100 py-8">
        <div className="container mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="p-4 border rounded-md shadow-md">
              <h2 className="text-xl font-semibold mb-2">Test Endpoints</h2>
              <p className="mb-4">Test the deposit request and order status request endpoints.</p>
              <Link to="/v1/test-endpoints" className="bg-blue-500 text-white py-2 px-4 rounded-md">
                Test Endpoints
              </Link>
            </div>
            <div className="p-4 border rounded-md shadow-md">
              <h2 className="text-xl font-semibold mb-2">Test Flows</h2>
              <p className="mb-4">Test the client and backend flows.</p>
              <Link to="/v1/flows" className="bg-blue-500 text-white py-2 px-4 rounded-md">
                Test Flows
              </Link>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
};

export default Home;
