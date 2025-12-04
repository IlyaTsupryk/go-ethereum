import { ethers } from "ethers";
import { schnorr } from "@noble/curves/secp256k1";
import { generateUserKey } from "./utils.js";
import StringStorageArtifact from "../artifacts/contracts/StringStorage.sol/StringStorage.json";

const CONTRACT_ADDRESS = "0xE74A3C7427CDA785e0000D42a705B1f3fD371E09";

async function main() {
  const provider = new ethers.JsonRpcProvider("http://127.0.0.1:8545");

  const charlieKey = await generateUserKey(provider, 'charlie');
  const charliePubKey = '0x' + charlieKey.publicKey.slice(4); // remove 0x04 prefix
  console.log(`Charlie's address: ${charlieKey.address}`);

  const signer = new ethers.Wallet(charlieKey.privateKey, provider);
  const stringStorage = new ethers.Contract(CONTRACT_ADDRESS, StringStorageArtifact.abi, signer);

  const inputString = "some Charlie string";
  console.log(`String to add: "${inputString}"\n`);

  const messageHash = ethers.keccak256(ethers.toUtf8Bytes(inputString));
  console.log(`Message hash: ${messageHash}`);

  const privateKeyBytes = ethers.getBytes(charlieKey.privateKey);
  const messageHashBytes = ethers.getBytes(messageHash);

  const signature = schnorr.sign(messageHashBytes, privateKeyBytes);
  console.log(`Schnorr signature: 0x${Buffer.from(signature).toString('hex')}`);

  // ---------------------------------------------------
  // Properly signed string should be saved in contract
  // ---------------------------------------------------
  console.log("Calling addStringSecure...");
  const tx = await stringStorage.addStringSecure(
    inputString,
    signature,
    charliePubKey,
  );
  await tx.wait();

  console.log("Retrieving the added string...");
  const str = await stringStorage.getString(charlieKey.address, 0);
  console.log(`String retrieved from contract: "${str}"`);

  // ---------------------------------------------------
  // Try to add a string which was not signed properly
  // ---------------------------------------------------
  const invalidInputString = "another invalid string";
  console.log("\n\nCalling addStringSecure with invalid string...");
  try {
    const invalidTx = await stringStorage.addStringSecure(
        invalidInputString,
        signature,
        charliePubKey,
    );
    await invalidTx.wait();
    console.error("ERROR: Transaction should have failed but succeeded!");
  } catch (error: any) {
    if (error.message && error.message.includes("Invalid Schnorr signature")) {
      console.log("Transaction correctly failed with 'Invalid Schnorr signature' error");
    } else {
      console.error("ERROR: Transaction failed with unexpected error:");
      console.error(error.message || error);
    }
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });