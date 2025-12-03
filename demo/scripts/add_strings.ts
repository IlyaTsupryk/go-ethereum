import { ethers } from "ethers";
import StringStorageArtifact from "../artifacts/contracts/StringStorage.sol/StringStorage.json";
import {generateUserKey} from "./utils.js";

const CONTRACT_ADDRESS = "0x3A220f351252089D385b29beca14e27F204c296A";

async function main() {
  const provider = new ethers.JsonRpcProvider("http://127.0.0.1:8545");
  console.log("Genarate 2 users and save some data for them");

  const aliceKey = await generateUserKey(provider, 'alice');
  console.log(`Alice address: ${aliceKey.address}`);
  const aliceContract = new ethers.Contract(
      CONTRACT_ADDRESS,
      StringStorageArtifact.abi,
      new ethers.Wallet(aliceKey.privateKey, provider)
  );
  let tx = await aliceContract.addString('some Alice string');
  await tx.wait();
  console.log('Alice saved a string');

  const bobKey = await generateUserKey(provider, 'bob');
  console.log(`Bob address: ${aliceKey.address}`);
  const bobContract = new ethers.Contract(
      CONTRACT_ADDRESS,
      StringStorageArtifact.abi,
      new ethers.Wallet(bobKey.privateKey, provider)
  );
  tx = await bobContract.addString('some Bob string');
  await tx.wait();
  console.log('Bob saved a string');

  console.log("\nRetrieving saved data...");
  let str = await aliceContract.getString(aliceKey.address, 0);
  console.log(`  ["alice", 0]: "${str}"`);
  str = await aliceContract.getString(bobKey.address, 0);
  console.log(`  ["bob", 0]: "${str}"`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });