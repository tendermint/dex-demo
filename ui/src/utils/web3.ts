import Web3 from "web3";

let cachedWeb3: Web3;

export default function getWeb3() {
  const { web3, ethereum } = window as any;

  if (cachedWeb3) {
    return cachedWeb3;
  }

  if (ethereum) {
    cachedWeb3 = new Web3(ethereum);
  } else if (web3) {
    cachedWeb3 = new Web3(web3.currentProvider);
  } else {
    cachedWeb3 = new Web3(new Web3.providers.HttpProvider('https://rinkeby.infura.io/v3/2fb7c637ca404ddb9407385391f7232d'));
  }

  return cachedWeb3;
}

export async function getCurrentBlock() {
  const web3 = getWeb3();
  return await web3.eth.getBlock('latest');
}