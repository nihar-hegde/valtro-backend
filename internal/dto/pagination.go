package dto

// PaginationRequest represents pagination parameters in requests
type PaginationRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

// PaginationMetadata represents pagination information in responses
type PaginationMetadata struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
	TotalItems   int64 `json:"total_items"`
	HasNextPage  bool  `json:"has_next_page"`
	HasPrevPage  bool  `json:"has_prev_page"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Message    string             `json:"message"`
	Data       interface{}        `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// NewPaginationMetadata creates pagination metadata from request and total count
func NewPaginationMetadata(req PaginationRequest, totalItems int64) PaginationMetadata {
	totalPages := int((totalItems + int64(req.PageSize) - 1) / int64(req.PageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationMetadata{
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		HasNextPage: req.Page < totalPages,
		HasPrevPage: req.Page > 1,
	}
}

// Validate validates and sets defaults for pagination request
func (p *PaginationRequest) Validate() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20 // Default page size
	}
	if p.PageSize > 100 {
		p.PageSize = 100 // Max page size
	}
}

// Offset calculates the database offset for the pagination
func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the limit for database queries
func (p PaginationRequest) Limit() int {
	return p.PageSize
}
