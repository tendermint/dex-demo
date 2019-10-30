### Aggregate Supply/Aggregate Demand (AS/AD) Curve

Figure 1 depicts the process flow for frequent
batch auctions. There are three components:
order submission, auction, and reporting. We
describe the design details for each component in
turn. This section focuses on a non-fragmented
market; we discuss how to augment frequent
batch auctions to accommodate fragmented
markets (e.g., US equities markets) below.
Throughout, design details are chosen to minimize the departure from current practice, subject
to realizing the benefits of frequent batching.
This is both to reduce transition costs and to
limit the scope for unintended consequences.
A. Order Submission
Orders in a batch auction consist of a direction (buy or sell), a price, and a quantity, just
like limit orders in a CLOB. During the order
submission stage, orders can be freely submitted, modified, or withdrawn. If an order is not
executed in the batch auction at time t, it automatically carries over for the next auction at
time t + 1, t + 2, etc., until it is either executed
or withdrawn.
Orders are not displayed during the order
submission stage.1
 This is important to prevent
gaming, and is why we describe the auction as
“sealed bid.” Orders are instead displayed in
aggregate at the reporting stage, as described
below.
An important open question is the optimal
tick size in a batch auction. The simplest policy
would be to mimic the tick size under the current CLOB, which in the United States is $0.01
by regulation ($0.0001 for stocks less than $1
per share). However, we note that one of the
arguments against finer tick sizes—the explosion in message traffic that arises from traders
outbidding each other by economically negligible amounts—is moot here, because orders are
opaque during the batch interval. We also note
that the coarser the tick size, the more important
a role rationing will play. For these reasons, we
conjecture that the optimal tick size in a frequent
batch auction is at least as fine as in the continuous market.
B. Auction
At the conclusion of the order submission
stage, the exchange batches all of the outstanding orders, and computes the aggregate demand
and supply functions from orders to buy and
sell, respectively. There are two cases: supply
and demand cross, or they do not. See Figure 2.