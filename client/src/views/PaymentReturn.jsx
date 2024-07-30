import { useEffect, useState } from "react";
import axios from "axios";
import backArrow from "../assets/back-arrow-svgrepo-com.svg";
import { Link } from "react-router-dom";
import { GenerateSignature } from "../utils/utils";

const PaymentReturn = () => {
  const baseURL = "http://localhost:8080";
  const orderStatusReqEndpointURL = "/api/v1/query/order-status/";
  const [data, setData] = useState({
    type: "",
    status: "",
    errorMessage: "",
    endpointID: "",
    processorTransactionID: "",
    merchantOrderID: "",
    orderID: "",
    amount: "",
    currency: "",
    customerEmail: "",
    customParam: "",
    extraData: {
      amountChanged: "",
      amountRounded: "",
      amountManipulated: "",
      dcc: "",
      originalAmount: "",
      paymentMethod: "",
      selectedBankCode: "",
      selectedBankName: "",
    },
    request: {
      merchantID: "",
      merchantOrderID: "",
      orderID: "",
      timestamp: "",
    },
  });
  const [message, setMessage] = useState("");

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const status = urlParams.get("status");
    const orderID = urlParams.get("orderID");
    const merchantOrderID = urlParams.get("merchantOrderID");
    const paramsSignature = urlParams.get("signature");

    const generateSignature = async () => {
      const signatureString = `${status}${orderID}${merchantOrderID}${import.meta.env.VITE_MERCHANT_SECRET_KEY}`;
      const signature = await GenerateSignature(signatureString);
      return signature;
    };

    const fetchOrderStatus = async () => {
      try {
        const queryString = `merchantID=${import.meta.env.VITE_MERCHANT_ID}&merchantOrderID=${merchantOrderID}&orderID=${orderID}&timestamp=${localStorage.getItem("timestamp")}&signature=${localStorage.getItem("signature")}`;
        const response = await axios.get(`${baseURL}${orderStatusReqEndpointURL}?${queryString}`);

        const data = response?.data;
        setMessage(data?.message);
        setData(data?.data);

        if (data.data?.status === "APPROVED" || data.data?.status === "DECLINED" || data.data?.status === "FILTERED" || data.data?.status === "ERROR") {
          setMessage((prevMessage) => `The order status is now final.\n\n${prevMessage}`);
        } else {
          // Continue polling if the status is not final
          setTimeout(fetchOrderStatus, 13000); // Poll every 13 seconds
        }
      } catch (error) {
        if (error.response) {
          // server responded with a status other than 2xx
          console.error("Error response:", error.response.data);
          setMessage(error.response.data.message);
        } else {
          // something else happened while setting up the request
          console.error("Error:", error.message);
          setMessage(`Error: ${error.message}`);
        }
      }
    };

    const initiateFetch = async () => {
      try {
        const signature = await generateSignature();
        if (signature) {
          if (signature === paramsSignature) {
            console.log("Signatures is verified");
            fetchOrderStatus();
          } else {
            console.error("Signatures do not match");
            setMessage("Signatures do not match");
          }
        }
      } catch (error) {
        if (error.response) {
          // server responded with a status other than 2xx
          console.error("Error response:", error.response.data);
          setMessage(error.response.data.message);
        } else {
          // something else happened while setting up the request
          console.error("Error:", error.message);
          setMessage(`Error: ${error.message}`);
        }
      }
    };

    initiateFetch();
  }, [data?.data?.status]);

  return (
    <div className="payment-return container mx-auto p-4">
      <div className="flex justify-between items-center pb-5">
        <Link to="/" onClick={() => localStorage.clear()}>
          <img src={backArrow} alt="back arrow" className="w-8 h-8" />
        </Link>
        <h1 className="text-2xl font-bold mb-4">Payment Status</h1>
        {message && <p className="text-red-500">{message}</p>}
      </div>
      {data && (
        <div className="p-4 border rounded-md shadow-md">
          <div className="flex justify-between items-center">
            <h2 className="text-xl font-semibold mb-2">Order Status Response</h2>
            <p className="ml-2 text-sm text-gray-500">
              <span className="bg-gray-200 p-1 rounded-md">{data?.status}</span>Polling every 13 seconds until a final status is received
            </p>
          </div>
          <pre>
            <code>{JSON.stringify(data, null, 2)}</code>
          </pre>
        </div>
      )}
    </div>
  );
};

export default PaymentReturn;
