import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const StringStorageModule = buildModule("StringStorageModule", (m) => {
  const stringStorage = m.contract("StringStorage");
  return { stringStorage };
});

export default StringStorageModule;