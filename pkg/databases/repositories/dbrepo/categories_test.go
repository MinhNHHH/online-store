package dbrepo

import (
	"testing"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
	"github.com/stretchr/testify/assert"
)

func TestAllCategories(t *testing.T) {
	tests := []struct {
		name        string
		category    *schema.Category
		wantErr     bool
		description string
	}{
		{
			name: "Test category 1",
			category: &schema.Category{
				Name:        "Test Category 1",
				Description: "Test Description 1",
			},
			wantErr:     false,
			description: "Should successfully insert a Test category 1",
		},
		{
			name: "Test category 2",
			category: &schema.Category{
				Name:        "Test Category 2",
				Description: "Test Description 2",
			},
			wantErr:     false,
			description: "Should successfully insert a Test category 2",
		},
	}

	for _, tt := range tests {
		id, err := testRepo.InsertCategory(tt.category)
		assert.NoError(t, err)
		assert.Greater(t, id, 0)
	}

	categories, total, err := testRepo.AllCategories("", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.LessOrEqual(t, len(categories), 10)
}

func TestGetCategoryByFilter(t *testing.T) {
	categories, total, err := testRepo.AllCategories("test category 1", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, "Test Category 1", categories[0].Name)
}

func TestInsertCategory(t *testing.T) {
	tests := []struct {
		name        string
		category    *schema.Category
		wantErr     bool
		description string
	}{
		{
			name: "Valid category",
			category: &schema.Category{
				Name:        "Test Category",
				Description: "Test Description",
			},
			wantErr:     false,
			description: "Should successfully insert a valid category",
		},
		{
			name: "Empty name",
			category: &schema.Category{
				Name:        "",
				Description: "Test Description",
			},
			wantErr:     true,
			description: "Should fail when name is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := testRepo.InsertCategory(tt.category)
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

func TestUpdateCategory(t *testing.T) {
	// Insert test category
	category := &schema.Category{
		Name:        "Original Name",
		Description: "Original Description",
	}
	id, err := testRepo.InsertCategory(category)
	assert.NoError(t, err)
	assert.Greater(t, id, 0)

	tests := []struct {
		name        string
		category    *schema.Category
		wantErr     bool
		description string
	}{
		{
			name: "Valid update",
			category: &schema.Category{
				ID:          id,
				Name:        "Updated Name",
				Description: "Updated Description",
			},
			wantErr:     false,
			description: "Should successfully update category",
		},
		{
			name: "Non-existent ID",
			category: &schema.Category{
				ID:          999,
				Name:        "Test Name",
				Description: "Test Description",
			},
			wantErr:     false,
			description: "Should not error on non-existent ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testRepo.UpdateCategory(tt.category)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteCategory(t *testing.T) {
	// Insert test category
	category := &schema.Category{
		Name:        "To Delete",
		Description: "Will be deleted",
	}
	id, err := testRepo.InsertCategory(category)
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
			description: "Should successfully delete category",
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
			err := testRepo.DeleteCategory(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
