import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import App from "../../src/App";

const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

describe("Trigger all endpoints and verify responses", () => {
  test("Trigger deposit request and order status request, then verify responses", async () => {
    // render App
    const wrapper = render(<App />);
    expect(wrapper).toBeTruthy();

    // navigate to TestEndpoints view
    const testEndpointsLink = screen.getByRole("link", { name: /Test Endpoints/i });
    expect(testEndpointsLink).toBeInTheDocument();
    fireEvent.click(testEndpointsLink);
    await waitFor(() => expect(window.location.pathname).toBe("/v1/test-endpoints"));

    // generate a unique random merchantOrderID
    let result = "";
    for (let i = 0; i < 16; i++) {
      result += characters.charAt(Math.floor(Math.random() * characters.length));
    }
    // get merchantOrderID input field, and update its value
    const merchantOrderIDInput = screen.getByPlaceholderText("merchantOrderID");
    expect(merchantOrderIDInput).toBeInTheDocument();
    fireEvent.change(merchantOrderIDInput, { target: { value: result } });
    expect(merchantOrderIDInput).toHaveValue(result);

    // get deposit request button, and click it
    const depositButton = screen.getByRole("button", { name: /deposit request/i });
    expect(depositButton).toBeInTheDocument();
    fireEvent.click(depositButton);

    // verify deposit response
    const depositResponseHeading = screen.getByRole("heading", { level: 2, name: /deposit request response/i });
    expect(depositResponseHeading).toBeInTheDocument();
    await waitFor(() => {
      const depositResponse = depositResponseHeading.nextSibling.querySelector("code");
      expect(depositResponse).toBeInTheDocument();
      const response = JSON.parse(depositResponse.textContent);
      expect(response.depositUrl).not.toBe("");
      expect(response).toHaveProperty("merchantOrderID", result);
      expect(response.orderID).not.toBe("");
      console.log("deposit response:", response);
    });

    // prepare for order status request

    // get "Generate Timestamp and Auth Signature" button, and click it
    const generateTimestampButton = screen.getByRole("button", { name: /generate timestamp and auth signature/i });
    expect(generateTimestampButton).toBeInTheDocument();

    await waitFor(() => fireEvent.click(generateTimestampButton));

    // verify timestamp and auth signature
    const timestampInput = screen.getByPlaceholderText("timestamp");
    expect(timestampInput).toBeInTheDocument();
    expect(timestampInput.value).toMatch(/^\d{10}$/);
    const signatureInput = screen.getByPlaceholderText("signature");
    expect(signatureInput).toBeInTheDocument();
    expect(signatureInput.value).toMatch(/^\w{64}$/);

    // get order status request button, and click it
    const orderStatusButton = screen.getByRole("button", { name: /order status request/i });
    expect(orderStatusButton).toBeInTheDocument();

    await waitFor(() => fireEvent.click(orderStatusButton));

    // verify order status response
    const orderStatusResponseHeading = screen.getByRole("heading", { level: 2, name: /order status response/i });
    expect(orderStatusResponseHeading).toBeInTheDocument();
    await waitFor(() => {
      const orderStatusResponse = orderStatusResponseHeading.nextSibling.querySelector("code");
      expect(orderStatusResponse).toBeInTheDocument();
      const response = JSON.parse(orderStatusResponse.textContent);
      expect(response.status).not.toBe("");
      expect(response.merchantOrderID).toBe(result);
      expect(response.orderID).not.toBe("");
      expect(response.amount).not.toBe("");
      expect(response.currency).not.toBe("");
      console.log("order status response:", response);
    });
  });
});
