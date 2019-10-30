import MetamaskIcon from "../assets/icons/metamask.svg";

export const ETH_CLEARINGHOUSE_ADDRESS = '0x526aaAAACF82B100FF0590B905Ed1554B6a9047B';
export const ETH_CLEARINGHOUSE_ABI = [
  {
    "constant": false,
    "name": "deposit",
    "inputs": [
      {
        "name": "recipient",
        "type": "bytes20"
      }
    ],
    "payable": true,
    "stateMutability": "payable",
    "type": "function"
  },
  {
    "constant": false,
    "name": "withdraw",
    "inputs": [
      { "name": "root", "type": "bytes32" },
      { "name": "assetId", "type": "uint256" },
      { "name": "burnId", "type": "uint256" },
      { "name": "beneficiary", "type": "address" },
      { "name": "amount", "type": "uint256" },
      { "name": "proof", "type": "bytes" },
    ],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  },
];

export enum FundingSources {
  metamask = 'metamask',
}

export type FundingSourceDataType = {
  name: string
  logoUrl: string
}

const FundingSourcesData: { [k: string]: FundingSourceDataType } = {
  [FundingSources.metamask]: {
    name: 'Metamask',
    logoUrl: MetamaskIcon,
  }
};

export function getFundingSourceData(source: FundingSources | null): FundingSourceDataType {
  if (!source) {
    return {
      name: '',
      logoUrl: '',
    };
  }
  return FundingSourcesData[source];
}
