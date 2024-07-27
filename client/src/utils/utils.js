export const GenerateSignature = async (data) => {
  const encoder = new TextEncoder();
  const dataUint8 = encoder.encode(data);
  // await the promise returned by crypto.subtle.digest
  const hashBuffer = await crypto.subtle.digest("SHA-256", dataUint8);
  // convert the buffer to a byte array
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  // convert the byte array to a hex string
  const hashHex = hashArray.map((b) => b.toString(16).padStart(2, "0")).join("");
  return hashHex;
};
