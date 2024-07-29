import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import App from "../../src/App";

describe("Home view renders and user can navigate between all pages", () => {
  beforeEach(() => {
    const wrapper = render(<App />);
    expect(wrapper).toBeTruthy();
    const homeLink = screen.getByRole("link", { name: /Home/i });
    expect(homeLink).toBeInTheDocument();
    fireEvent.click(homeLink);
    expect(window.location.pathname).toBe("/");
  });

  test("navigate to the profile view", async () => {
    const profileLink = screen.getByRole("link", { name: /Profile/i });
    expect(profileLink).toBeInTheDocument();
    fireEvent.click(profileLink);
    await waitFor(() => expect(window.location.pathname).toBe("/v1/account/deposit/:uid"));
  });

  test("navigate to the TestEndpoints view", async () => {
    const testEndpointsLink = screen.getByRole("link", { name: /Test Endpoints/i });
    expect(testEndpointsLink).toBeInTheDocument();
    fireEvent.click(testEndpointsLink);
    await waitFor(() => expect(window.location.pathname).toBe("/v1/test-endpoints"));
  });

  test("navigate to the TestFlows view and back to home view", async () => {
    const testFlowsLink = screen.getByRole("link", { name: /Test Flows/i });
    expect(testFlowsLink).toBeInTheDocument();
    fireEvent.click(testFlowsLink);
    await waitFor(() => expect(window.location.pathname).toBe("/v1/flows"));

    const homeLink = screen.getByRole("link", { name: /Home/i });
    expect(homeLink).toBeInTheDocument();
    fireEvent.click(homeLink);
    await waitFor(() => expect(window.location.pathname).toBe("/"));
  });
});
