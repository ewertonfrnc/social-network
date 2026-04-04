package store

const (
	defaultFeedLimit         = 20
	defaultFeedOffset        = 0
	defaultFeedSortDirection = "desc"
)

type PaginatedFeedQuery struct {
	Limit         int      `json:"limit" validate:"gte=1,lte=100"`
	Offset        int      `json:"offset" validate:"gte=0"`
	SortDirection string   `json:"sort" validate:"oneof=asc desc"`
	Tags          []string `json:"tags" validate:"max=5"`
	Search        string   `json:"search" validate:"max=100"`
}

func NewPaginatedFeedQuery() PaginatedFeedQuery {
	return PaginatedFeedQuery{
		Limit:         defaultFeedLimit,
		Offset:        defaultFeedOffset,
		SortDirection: defaultFeedSortDirection,
	}
}
