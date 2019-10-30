package embedded

type BlockInclusion struct {
	BlockNumber     int64  `json:"block_number"`
	TransactionHash string `json:"block_hash"`
	BlockTimestamp  string `json:"block_timestamp"`
}
