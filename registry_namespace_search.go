package client

import (
	"context"

	"github.com/cmskitdev/notion/types"
)

// SearchNamespace provides fluent access to search operations.
type SearchNamespace struct {
	registry *Registry
}

// Query performs a search query across the workspace.
//
// Arguments:
//   - ctx: Context for cancellation and timeouts.
//   - req: Search request.
//
// Returns:
//   - <-chan Result[SearchResult]: Channel of search results.
func (ns *SearchNamespace) Query(ctx context.Context, req SearchRequest) <-chan Result[types.SearchResult] {
	paginatedOp := NewPaginatedOperator[types.SearchResult](ns.registry.httpClient, DefaultOperatorConfig())

	r := &SearchRequest{
		SearchRequest: types.SearchRequest{
			Filter: &types.SearchFilter{
				Property: "object",
				Value:    "database",
			},
		},
	}

	return StreamPaginated(paginatedOp, ctx, r)
}

// QueryAll performs a search query and returns all results as a slice.
//
// Arguments:
//   - ctx: Context for cancellation and timeouts.
//   - req: Search request.
//
// Returns:
//   - Result[[]SearchResult]: All search results as a slice.
func (ns *SearchNamespace) QueryAll(ctx context.Context, req types.SearchRequest) Result[[]types.SearchResult] {
	// Check if search operator is available
	_, err := GetTyped[*SearchOperator[types.SearchResult]](ns.registry, "search")
	if err != nil {
		return Error[[]types.SearchResult](err)
	}

	searchOp := NewSearchOperator[types.SearchResult](ns.registry.httpClient, DefaultOperatorConfig())

	return searchOp.Execute(ctx, req)
}

// Stream performs a search query with pagination and streams results.
//
// Arguments:
//   - ctx: Context for cancellation and timeouts.
//   - req: Search request.
//
// Returns:
//   - <-chan Result[types.SearchResult]: Channel of search results.
func (ns *SearchNamespace) Stream(ctx context.Context, req types.SearchRequest) <-chan Result[types.SearchResult] {
	// Check if search operator is available
	_, err := GetTyped[*SearchOperator[types.SearchResult]](ns.registry, "search")
	if err != nil {
		resultCh := make(chan Result[types.SearchResult], 1)
		resultCh <- Error[types.SearchResult](err)
		close(resultCh)
		return resultCh
	}

	searchOp := NewSearchOperator[types.SearchResult](ns.registry.httpClient, DefaultOperatorConfig())
	return searchOp.Stream(ctx, req)
}
