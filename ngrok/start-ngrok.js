const { exec } = require("child_process");
const fs = require("fs");
const path = require("path");

const startNgrok = () => {
  return new Promise((resolve, reject) => {
    const process = exec(`ngrok start --all --config=ngrok.yml`, (error, stdout, stderr) => {
      if (error) {
        reject(error);
        return;
      }
      if (stderr) {
        console.error(`Ngrok stderr: ${stderr}`);
      }
    });

    console.log("Ngrok process started. Waiting for tunnels to initialize...");

    setTimeout(async () => {
      try {
        // Check the status endpoint first
        const statusResponse = await fetch("http://localhost:4040/api/status");
        if (!statusResponse.ok) {
          reject(new Error("Ngrok status endpoint not reachable"));
          return;
        }
        const statusResponseJson = await statusResponse.json();
        console.log("Ngrok status response:", statusResponseJson);

        // Fetch tunnels information
        const response = await fetch("http://localhost:4040/api/tunnels");
        if (!response.ok) {
          reject(new Error("Ngrok tunnels endpoint not reachable"));
          return;
        }
        const responseJson = await response.json();
        // console.log("Ngrok tunnels response:", responseJson);
        const tunnels = responseJson.tunnels;
        const frontendTunnel = tunnels.find((tunnel) => tunnel.name === "frontend");
        const backendTunnel = tunnels.find((tunnel) => tunnel.name === "backend");
        if (frontendTunnel && backendTunnel) {
          console.log("Ngrok tunnels found:", frontendTunnel.public_url, backendTunnel.public_url);
          resolve({ frontendUrl: frontendTunnel.public_url, backendUrl: backendTunnel.public_url });
        } else {
          reject(new Error("Ngrok tunnels not found"));
        }
      } catch (error) {
        reject(error);
      }
    }, 10000);
  });
};

const updateEnvFile = (frontendUrl, backendUrl) => {
  const envFilePath = path.join(__dirname, "..", "client", ".env");

  // read existing .env file
  let existingVars = {};
  if (fs.existsSync(envFilePath)) {
    const envContent = fs.readFileSync(envFilePath, "utf8");
    existingVars = envContent
      .split("\n")
      .filter((line) => line.trim() && !line.startsWith("#"))
      .reduce((acc, line) => {
        const [key, value] = line.split("=");
        if (key && value) acc[key.trim()] = value.trim();
        return acc;
      }, {});
  }

  // add new or update existing variables
  const envVars = {
    ...existingVars,
    VITE_REDIRECT_URL: `${frontendUrl}/v1/payment/return`,
    VITE_CALLBACK_URL: `${backendUrl}/api/v1/payment-callback/`,
  };

  // write updated .env file
  fs.writeFileSync(
    envFilePath,
    Object.entries(envVars)
      .map(([key, value]) => `${key}=${value}`)
      .join("\n")
  );

  console.log("Environment variables updated:", envVars);
};

(async () => {
  try {
    const { frontendUrl, backendUrl } = await startNgrok();

    console.log("frontendUrl", frontendUrl);
    console.log("backendUrl", backendUrl);

    updateEnvFile(frontendUrl, backendUrl);

    console.log("Ngrok tunnels started:", frontendUrl, backendUrl);
  } catch (error) {
    console.error("Error starting ngrok:", error);
  }
})();
