import {
  estimateBatch,
  aggregateSupply,
  aggregateDemand,
} from './exchange-utils';
import {Order} from "../ducks/exchange";
import {bn} from "./bn";

type AggregateTestCase = {
  stop: number
  expected: number
}

describe('aggregateSupply', () => {
  const orders: Order[] = [
    mo(986391, 1410571704),
    mo(986417, 119892334),
    mo(986419, 85064912),
    mo(986422, 40000000),
    mo(986516, 1406452062),
    mo(986737, 159467720),
    mo(988878, 2017720),
    mo(990461, 22780730),
    mo(990532, 12254590),
    mo(990557, 248346),
    mo(990604, 13544580),
    mo(990820, 248343),
    mo(990900, 248346),
    mo(991199, 200000000),
  ];
  const testcases: AggregateTestCase[] = [
    { stop: 986390, expected: 0 },
    { stop: 986737, expected: 3221448732 },
    { stop: 990531, expected: 3246247182 },
    { stop: 990533, expected: 3258501772 },
  ];

  testcases.forEach(({ stop, expected }) => {
    it(`find aggregate supply at price: ${stop}`, () => {
      expect(aggregateSupply(orders, bn(stop)).toNumber()).toBe(expected);
    });
  });
});

describe('aggregateDemand', () => {
  const orders: Order[] = [
    mo(986391, 1410571704),
    mo(986417, 119892334),
    mo(986419, 85064912),
    mo(986422, 40000000),
    mo(986516, 1406452062),
    mo(986737, 159467720),
    mo(988878, 2017720),
    mo(990461, 22780730),
    mo(990532, 12254590),
    mo(990557, 248346),
    mo(990604, 13544580),
    mo(990820, 248343),
    mo(990900, 248346),
    mo(991199, 200000000),
  ];
  const testcases: AggregateTestCase[] = [
    { stop: 991200, expected: 0 },
    { stop: 991199, expected: 200000000 },
    { stop: 986737, expected: 410810375 },
    { stop: 990531, expected: 226544205 },
    { stop: 990533, expected: 214289615 },
  ];

  testcases.forEach(({ stop, expected }) => {
    it(`find aggregate supply at price: ${stop}`, () => {
      expect(aggregateDemand(orders, bn(stop)).toNumber()).toBe(expected);
    });
  });
});

type EstimateBatchTestCase = {
  bids: Order[]
  asks: Order[]
  expectedClearingPrice: number
  expectedAskRation: number
  expectedBidRation: number
}

describe.only('estimateBatch', () => {
  const testcases: EstimateBatchTestCase[] = [
    {
      bids: [
        mo(100, 10),
        mo(101, 5),
        mo(102, 5),
      ],
      asks: [
        mo(103, 5),
        mo(104, 15),
        mo(105, 15),
      ],
      expectedClearingPrice: 0,
      expectedAskRation: 0,
      expectedBidRation: 0,
    },
    {
      bids: [
        mo(100, 10),
        mo(101, 5),
        mo(102, 5),
      ],
      asks: [
        mo(102, 5),
        mo(103, 5),
        mo(104, 15),
        mo(105, 15),
      ],
      expectedClearingPrice: 102,
      expectedAskRation: 1,
      expectedBidRation: 1,
    },
    {
      bids: [
        mo(100, 10),
        mo(101, 5),
        mo(102, 5),
      ],
      asks: [
        mo(102, 4),
        mo(103, 5),
        mo(104, 15),
        mo(105, 15),
      ],
      expectedClearingPrice: 102,
      expectedAskRation: 1,
      expectedBidRation: .8,
    },
    {
      bids: [
        mo(100, 10),
        mo(101, 5),
        mo(102, 5),
      ],
      asks: [
        mo(101, 7),
        mo(103, 5),
        mo(104, 15),
        mo(105, 15),
      ],
      expectedClearingPrice: 101,
      expectedAskRation: 1,
      expectedBidRation: .7,
    },
    {
      bids: [
        mo(100, 10),
        mo(101, 5),
        mo(102, 5),
      ],
      asks: [
        mo(101, 10),
        mo(103, 5),
        mo(104, 15),
        mo(105, 15),
      ],
      expectedClearingPrice: 101,
      expectedAskRation: 1,
      expectedBidRation: 1,
    },
    {
      bids: [
        mo(100, 10),
        mo(103, 5),
      ],
      asks: [
        mo(101, 4),
        mo(104, 15),
        mo(105, 15),
      ],
      expectedClearingPrice: 103,
      expectedAskRation: 1,
      expectedBidRation: .8,
    },
    {
      bids: [
        mo(100, 10),
        mo(103, 5),
        mo(105, 5),
        mo(200, 5),
      ],
      asks: [
        mo(101, 4),
        mo(104, 15),
        mo(105, 15),
      ],
      expectedClearingPrice: 103,
      expectedAskRation: 1,
      expectedBidRation: 0.26666666666666666,
    },
    {
      bids: [
        mo(99, 5),
        mo(103, 5),
      ],
      asks: [
        mo(101, 5),
        mo(104, 10),
        mo(105, 10),
      ],
      expectedClearingPrice: 103,
      expectedAskRation: 1,
      expectedBidRation: 1,
    },
  ];

  testcases.forEach(({ bids, asks, expectedAskRation, expectedBidRation, expectedClearingPrice }) => {
    const { clearingPrice, askRation, bidRation } = estimateBatch(bids, asks);
    it(`expect clearing price at ${expectedClearingPrice}`, () => {
      expect(clearingPrice.toNumber()).toBe(expectedClearingPrice);
      expect(askRation.toNumber()).toBe(expectedAskRation);
      expect(bidRation.toNumber()).toBe(expectedBidRation);
    });
  })
});

function mo(price: number, quantity: number): Order {
  return {
    marketId: '1',
    price: bn(price),
    quantity: bn(quantity),
  };
}
