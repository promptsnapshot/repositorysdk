package repositorysdk

type QueryResult struct {
	Took    uint  `json:"took"`
	Timeout bool  `json:"timeout"`
	Shards  Shard `json:"_shards"`
}

type Shard struct {
	Total      uint `json:"total"`
	Successful uint `json:"successful"`
	Skipped    uint `json:"skipped"`
	Failed     uint `json:"failed"`
}
