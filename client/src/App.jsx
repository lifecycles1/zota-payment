import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Home from "./views/Home";
import Profile from "./views/Profile";
import PaymentReturn from "./views/PaymentReturn";
import TestFlows from "./views/TestFlows";
import TestEndpoints from "./views/TestEndpoints";
import NavBar from "./components/NavBar";

function App() {
  return (
    <>
      <NavBar />
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/v1/test-endpoints" element={<TestEndpoints />} />
          <Route path="/v1/flows" element={<TestFlows />} />
          <Route path="/v1/account/deposit/:uid" element={<Profile />} />
          <Route path="/v1/payment/return" element={<PaymentReturn />} />
        </Routes>
      </Router>
    </>
  );
}

export default App;
