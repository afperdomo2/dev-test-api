package notifications

import "github.com/felipe/dev-test-api/internal/common"

// El listado no expone sortBy: el orden es parte del producto (created_at DESC).
var sortConfig = common.SortConfig{
	Allowed: []string{},
	Default: "created_at DESC",
}

type ListNotificationsParams struct {
	common.PaginationParams
}
