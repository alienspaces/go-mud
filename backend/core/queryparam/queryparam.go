package queryparam

type QueryParams struct {
	Params      map[string][]string
	SortColumns []SortColumn
	PageSize    int
	PageNumber  int
}

type SortColumn struct {
	Col          string
	IsDescending bool
}
