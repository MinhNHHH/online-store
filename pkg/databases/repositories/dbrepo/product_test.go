package dbrepo

import (
	"testing"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
	"github.com/stretchr/testify/assert"
)

func TestAllProducts(t *testing.T) {
	mockCategory := schema.Category{
		Name:        "Test Category 1",
		Description: "Test Description 1",
	}

	categoryID, err := testRepo.InsertCategory(&mockCategory)
	assert.NoError(t, err)
	assert.Greater(t, categoryID, 0)

	mockProduct := []*schema.Product{
		{
			Name:          "Test Product",
			Description:   "Test Description",
			Price:         99.99,
			StockQuantity: 10,
			Status:        "in_stock",
			CategoryID:    categoryID,
		},
		{
			Name:          "Test Product 2",
			Description:   "Test Description 2",
			Price:         99.99,
			StockQuantity: 10,
			Status:        "in_stock",
			CategoryID:    categoryID,
		},
		{
			Name:          "Test Product 3",
			Description:   "Test Description 3",
			Price:         99.99,
			StockQuantity: 10,
			Status:        "in_stock",
			CategoryID:    categoryID,
		},
		{
			Name:          "Test Product 4",
			Description:   "Test Description 4",
			Price:         99.99,
			StockQuantity: 10,
			Status:        "out_of_stock",
			CategoryID:    categoryID,
		},
	}

	for _, product := range mockProduct {
		id, err := testRepo.InsertProduct(product)
		assert.NoError(t, err)
		assert.Greater(t, id, 0)
	}

	tests := []struct {
		name         string
		filterName   string
		categoryName string
		status       string
		page         int
		pageSize     int
		wantTotal    int
		description  string
	}{
		{
			name:         "Get all products",
			filterName:   "",
			categoryName: "",
			status:       "",
			page:         1,
			pageSize:     10,
			wantTotal:    4,
			description:  "Should return all products with pagination",
		},
		{
			name:         "Filter by name",
			filterName:   "Test",
			categoryName: "",
			status:       "",
			page:         1,
			pageSize:     10,
			wantTotal:    4,
			description:  "Should return products matching name filter",
		},
		{
			name:         "Filter by category",
			filterName:   "",
			categoryName: "Test Category 1",
			status:       "",
			page:         1,
			pageSize:     10,
			wantTotal:    4,
			description:  "Should return products in specified category",
		},
		{
			name:         "Filter by status",
			filterName:   "",
			categoryName: "",
			status:       "out_of_stock",
			page:         1,
			pageSize:     10,
			wantTotal:    1,
			description:  "Should return products with specified status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products, total, err := testRepo.AllProducts(tt.filterName, tt.categoryName, tt.status, tt.page, tt.pageSize)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantTotal, total)
			assert.LessOrEqual(t, len(products), tt.pageSize)
		})
	}
}

func TestInsertProduct(t *testing.T) {
	tests := []struct {
		name        string
		product     *schema.Product
		wantErr     bool
		description string
	}{
		{
			name: "Valid product",
			product: &schema.Product{
				Name:          "Test Product",
				Description:   "Test Description",
				Price:         99.99,
				StockQuantity: 10,
				Status:        "in_stock",
				CategoryID:    1,
			},
			wantErr:     false,
			description: "Should successfully insert a valid product",
		},
		{
			name: "Invalid price",
			product: &schema.Product{
				Name:          "Invalid Product",
				Description:   "Test Description",
				Price:         -10.0,
				StockQuantity: 10,
				Status:        "in_stock",
				CategoryID:    1,
			},
			wantErr:     true,
			description: "Should fail when price is negative",
		},
		{
			name: "Invalid stock quantity",
			product: &schema.Product{
				Name:          "Invalid Product",
				Description:   "Test Description",
				Price:         99.99,
				StockQuantity: -5,
				Status:        "in_stock",
				CategoryID:    1,
			},
			wantErr:     true,
			description: "Should fail when stock quantity is negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := testRepo.InsertProduct(tt.product)
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

func TestUpdateProduct(t *testing.T) {
	// Insert test product
	product := &schema.Product{
		Name:          "Original Product",
		Description:   "Original Description",
		Price:         99.99,
		StockQuantity: 10,
		Status:        "in_stock",
		CategoryID:    1,
	}
	id, err := testRepo.InsertProduct(product)
	assert.NoError(t, err)
	assert.Greater(t, id, 0)

	tests := []struct {
		name        string
		product     *schema.Product
		wantErr     bool
		description string
	}{
		{
			name: "Valid update",
			product: &schema.Product{
				ID:            id,
				Name:          "Updated Product",
				Description:   "Updated Description",
				Price:         149.99,
				StockQuantity: 20,
				Status:        "in_stock",
			},
			wantErr:     false,
			description: "Should successfully update product",
		},
		{
			name: "Non-existent ID",
			product: &schema.Product{
				ID:            999,
				Name:          "Non-existent Product",
				Description:   "Test Description",
				Price:         99.99,
				StockQuantity: 10,
				Status:        "in_stock",
			},
			wantErr:     false,
			description: "Should not error on non-existent ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testRepo.UpdateProduct(tt.product)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	// Insert test product
	product := &schema.Product{
		Name:          "To Delete",
		Description:   "Will be deleted",
		Price:         99.99,
		StockQuantity: 10,
		Status:        "in_stock",
		CategoryID:    1,
	}
	id, err := testRepo.InsertProduct(product)
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
			description: "Should successfully delete product",
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
			err := testRepo.DeleteProduct(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
