package config

const (
	PopulateIndex      = "4b5d8cc8c8e73033d2ba4c4e9eef90b009c3afc6"
	RemoveIndexEntries = "4f50bc2d90fe90f501acdc9e349d86faef2d329b"
	SearchListIndex    = "daacb0d0f789a25feb9be0d9c4891b2b51fd1913"
	GetArticle         = "20b7119e6ad77eaea9d81942b880d58adb4daba3"
	SearchZIndex       = "f81ae105ab5d928ac14921cc2b7e0e4011bfd0e1"
	/*
		populate_full_text_forward_index.lua
		ARGS:
		1. The hit
		2. The prefix to prepend to the hit when adding to the index
		3. The forward index prefix.
		[4:3+len(KEYS)]
	*/
	PopulateFullTextForwardIndex = "eea5b70f4781952e66a2c98c7f289b1b72cf47fd"
	/*
		remove_fts_hits.lua
		ARGS:
		1: The "hit" (articleID)
		2: The forward index prefix (this + word stores a sorted set containing hits and their scores)
		3: The reverse index prefix (this + hit stores a sorted set containing words and their frequencies)
		4: Hit prefix. This is what the hit is prefixed with in the sorted set. Beautiful.
		KEYS:
		The forward index keys at which the hit to remove is stored.
	*/
	RemoveHitFromFTS            = "2e69e815a3ab35f4a2589e73274008c5cfa3ed50"
	SnippetLength               = 200
	ArticlePrefix               = "articles."
	SnippetPrefix               = "snippets."
	FullTextSearchPrefix        = "fts."
	FullTextSearchForwardPrefix = FullTextSearchPrefix + "forward."
	FullTextSearchReversePrefix = FullTextSearchPrefix + "reverse."
	DateIndexKey                = "__date_index__"
	PositiveInfinity            = "+inf"
	NegativeInfinity            = "-inf"
	HitPrefix                   = SnippetPrefix
	DefaultLimit                = 3
	IDKey                       = "__id__"
	TopicSeparator              = ", "
)
