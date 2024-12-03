package repository

import (
	"github.com/gosimple/slug"
	"gocommerce/models/entity"
	"gocommerce/models/request"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetAll() ([]entity.Category, error)
	Create(request request.CategoryRequest) (entity.Category, error)
	GetBySlug(slug string) (entity.Category, error)
	GetById(id string) (entity.Category, error)
	Update(request request.CategoryUpdateRequest, id string) (entity.Category, error)
	Delete(id string) error
}

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func (c categoryRepositoryImpl) GetAll() ([]entity.Category, error) {
	var categories []entity.Category

	result := c.db.Debug().Find(&categories).Error

	if result != nil {
		return nil, result
	}

	return categories, nil
}

func (c categoryRepositoryImpl) Create(request request.CategoryRequest) (entity.Category, error) {
	newCategory := entity.Category{
		Name: request.Name,
		Slug: slug.Make(request.Name),
	}

	result := c.db.Debug().Create(&newCategory).Error
	if result != nil {
		return entity.Category{}, result
	}

	return newCategory, nil
}

func (c categoryRepositoryImpl) GetBySlug(slug string) (entity.Category, error) {
	var category entity.Category

	result := c.db.Debug().Where("slug = ?", slug).First(&category).Error
	if result != nil {
		return entity.Category{}, result
	}

	return category, nil
}

func (c categoryRepositoryImpl) GetById(id string) (entity.Category, error) {
	var category entity.Category

	result := c.db.Debug().Where("id = ?", id).First(&category).Error
	if result != nil {
		return entity.Category{}, result
	}

	return category, nil
}

func (c categoryRepositoryImpl) Update(request request.CategoryUpdateRequest, id string) (entity.Category, error) {
	category, err := c.GetById(id)
	if err != nil {
		return entity.Category{}, err
	}

	category.Name = request.Name
	category.Slug = slug.Make(request.Name)

	result := c.db.Debug().Save(&category).Error
	if result != nil {
		return entity.Category{}, result
	}

	return category, nil
}

func (c categoryRepositoryImpl) Delete(id string) error {
	category, err := c.GetById(id)
	if err != nil {
		return err
	}

	result := c.db.Debug().Delete(&category, id).Error
	if result != nil {
		return result
	}

	return nil
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}
