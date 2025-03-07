package application

import (
	"context"
	"strconv"

	"encore.app/developer_area/internal/utils"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type ApplicationListParams struct {
	Page    int `mod:"default=1" query:"page" validate:"min=1"`
	PerPage int `query:"per_page" validate:"min=1,max=100" mod:"default=10"`
}

func (params *ApplicationListParams) Validate() error {
	err := utils.ValidateTransform(context.Background(), params)
	return err
}

type ApplicationListResponse struct {
	Applications []ApplicationListItem `json:"applications"`
	TotalCount   int64                 `json:"total_count"`
	CurrentPage  int                   `json:"current_page,omitempty"`
	TotalPages   int                   `json:"total_pages,omitempty"`
}

type ApplicationListItem struct {
	Name string `json:"name"`
}

//encore:api auth method=GET path=/applications
func (s *Service) List(ctx context.Context, params *ApplicationListParams) (*ApplicationListResponse, error) {
	userID, hasUserId := auth.UserID()
	if !hasUserId {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "unauthenticated"}
	}

	ownerID, err := strconv.Atoi(string(userID))
	if err != nil {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "unauthenticated"}
	}

	page := params.Page
	perPage := params.PerPage

	query := s.db.Model(&Application{}).Where("owner_id = $1", ownerID)

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		rlog.Error("Database error counting applications.", "err", err, "userID", userID)
		return nil, &errs.Error{Code: errs.Unknown, Message: "Unknown error"}
	}

	var applications []Application
	err = query.
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&applications).Error
	if err != nil {
		rlog.Error("Database error fetching applications.", "err", err, "userID", userID, "page", page, "perPage", perPage)
		return nil, &errs.Error{Code: errs.Unknown, Message: "Unknown error"}
	}

	totalPages := int((totalCount + int64(perPage) - 1) / int64(perPage))

	response := ApplicationListResponse{
		Applications: make([]ApplicationListItem, len(applications)),
		CurrentPage:  page,
		TotalCount:   totalCount,
		TotalPages:   totalPages,
	}

	for i, app := range applications {
		response.Applications[i] = ApplicationListItem{
			Name: app.Name,
		}
	}

	return &response, nil
}
