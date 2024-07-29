import { render, screen } from "@testing-library/react";
import App from "../../src/App";

describe("App mounts properly, lands on Home view, and all links are rendered", () => {
  test("App mounts properly", () => {
    const wrapper = render(<App />);
    expect(wrapper).toBeTruthy();
  });

  beforeEach(() => {
    render(<App />);
  });

  test("Home view renders", () => {
    const home = screen.getByText(/Welcome to the home page/i);
    expect(home).toBeInTheDocument();
  });

  test("nav bar renders with links to Home and Profile", () => {
    const nav = screen.getByRole("navigation");
    expect(nav).toBeInTheDocument();
    const homeLink = screen.getByRole("link", { name: /Home/i });
    expect(homeLink).toBeInTheDocument();
    expect(homeLink).toHaveAttribute("href", "/");
    const profileLink = screen.getByRole("link", { name: /Profile/i });
    expect(profileLink).toBeInTheDocument();
    expect(profileLink).toHaveAttribute("href", "/v1/account/deposit/:uid");
  });

  test("Tile with link to TestEndpoints component renders", () => {
    const testEndpoints = screen.getByRole("heading", { level: 2, name: /Test Endpoints/i });
    expect(testEndpoints).toBeInTheDocument();
    const testEndpointsLink = screen.getByRole("link", { name: /Test Endpoints/i });
    expect(testEndpointsLink).toBeInTheDocument();
    expect(testEndpointsLink).toHaveAttribute("href", "/v1/test-endpoints");
  });

  test("Tile with link to TestFlows component renders", () => {
    const testFlows = screen.getByRole("heading", { level: 2, name: /Test Flows/i });
    expect(testFlows).toBeInTheDocument();
    const testFlowsLink = screen.getByRole("link", { name: /Test Flows/i });
    expect(testFlowsLink).toBeInTheDocument();
    expect(testFlowsLink).toHaveAttribute("href", "/v1/flows");
  });
});
