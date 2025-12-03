import { ethers } from "ethers";

export async function generateUserKey(provider: ethers.JsonRpcProvider, name: string) {
  const privateKey = ethers.id(name);
  const wallet = new ethers.Wallet(privateKey);

  const funder = await provider.getSigner(0);
  const tx = await funder.sendTransaction({
    to: wallet.address,
    value: ethers.parseEther("0.01")
  });
  await tx.wait();

  return {
    address: wallet.address,
    privateKey: privateKey
  };
}
