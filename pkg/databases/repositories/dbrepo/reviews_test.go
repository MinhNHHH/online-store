package dbrepo

import (
	"testing"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
	"github.com/stretchr/testify/assert"
)

func TestReviewsByProductID(t *testing.T) {
	categoryID, err := testRepo.InsertCategory(&schema.Category{
		Name:        "Test Category",
		Description: "Test Description",
	})
	assert.NoError(t, err)
	assert.Greater(t, categoryID, 0)

	productID, err := testRepo.InsertProduct(&schema.Product{
		Name:          "Test Product",
		Description:   "Test Description",
		Price:         100,
		StockQuantity: 10,
		Status:        "in_stock",
		CategoryID:    categoryID,
	})
	assert.NoError(t, err)
	assert.Greater(t, productID, 0)

	_, err = testRepo.InsertReview(&schema.Review{
		ProductID: productID,
		UserID:    1,
		Rating:    5,
		Comment:   "Great product!",
	})

	assert.NoError(t, err)

	tests := []struct {
		name        string
		productID   int
		wantCount   int
		description string
	}{
		{
			name:        "Get reviews for product 1",
			productID:   productID,
			wantCount:   1,
			description: "Should return one review for product 1",
		},
		{
			name:        "Get reviews for non-existent product",
			productID:   999,
			wantCount:   0,
			description: "Should return empty slice for non-existent product",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviews, err := testRepo.ReviewsByProductID(tt.productID)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCount, len(reviews))

			if tt.wantCount > 0 {
				// Verify review data structure
				review := reviews[0]
				assert.NotEmpty(t, review.ProductName)
				assert.NotEmpty(t, review.UserName)
				assert.GreaterOrEqual(t, review.Rating, 1)
				assert.LessOrEqual(t, review.Rating, 5)
			}
		})
	}
}

func TestInsertReview(t *testing.T) {
	tests := []struct {
		name        string
		review      *schema.Review
		wantErr     bool
		description string
	}{
		{
			name: "Valid review",
			review: &schema.Review{
				ProductID: 1,
				UserID:    1,
				Rating:    5,
				Comment:   "Great product!",
			},
			wantErr:     false,
			description: "Should successfully insert a valid review",
		},
		{
			name: "Invalid rating too high",
			review: &schema.Review{
				ProductID: 1,
				UserID:    1,
				Rating:    6,
				Comment:   "Invalid rating",
			},
			wantErr:     true,
			description: "Should fail when rating is above 5",
		},
		{
			name: "Invalid rating too low",
			review: &schema.Review{
				ProductID: 1,
				UserID:    1,
				Rating:    0,
				Comment:   "Invalid rating",
			},
			wantErr:     true,
			description: "Should fail when rating is below 1",
		},
		{
			name: "Non-existent product",
			review: &schema.Review{
				ProductID: 999,
				UserID:    1,
				Rating:    5,
				Comment:   "Non-existent product",
			},
			wantErr:     true,
			description: "Should fail when product does not exist",
		},
		{
			name: "Non-existent user",
			review: &schema.Review{
				ProductID: 1,
				UserID:    999,
				Rating:    5,
				Comment:   "Non-existent user",
			},
			wantErr:     true,
			description: "Should fail when user does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := testRepo.InsertReview(tt.review)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0, id)
			} else {
				assert.NoError(t, err)
				assert.Greater(t, id, 0)
			}
		})
	}
}

func TestDeleteReview(t *testing.T) {
	// First insert a test review
	review := &schema.Review{
		ProductID: 1,
		UserID:    1,
		Rating:    5,
		Comment:   "Test review to delete",
	}
	id, err := testRepo.InsertReview(review)
	assert.NoError(t, err)
	assert.Greater(t, id, 0)

	tests := []struct {
		name        string
		id          int
		wantErr     bool
		description string
	}{
		{
			name:        "Valid delete",
			id:          id,
			wantErr:     false,
			description: "Should successfully delete review",
		},
		{
			name:        "Non-existent ID",
			id:          999,
			wantErr:     false,
			description: "Should not error on non-existent ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testRepo.DeleteReview(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
