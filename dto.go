package repositorysdk

// QueryResult is a struct that holds the result of an Opensearch query, including the time taken to execute the query,
// whether the query timed out, and information about the shards.
type QueryResult struct {
	Took    uint  `json:"took"`
	Timeout bool  `json:"timeout"`
	Shards  Shard `json:"_shards"`
}

// Shard is a struct that holds information about the shards in an Opensearch cluster, including the total number of shards,
// the number of successful shards, the number of skipped shards, and the number of failed shards.
type Shard struct {
	Total      uint `json:"total"`
	Successful uint `json:"successful"`
	Skipped    uint `json:"skipped"`
	Failed     uint `json:"failed"`
}
