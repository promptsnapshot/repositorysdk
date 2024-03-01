package repositorysdk

import (
	"gorm.io/gorm"
	"math"
)

type Entity interface {
	TableName() string
}

// Pagination returns a function that can be used as a GORM scope to paginate results. It takes a pointer to a slice of the entity type, a pointer to a PaginationMetadata struct, a GORM database instance, and an optional list of additional GORM scopes. It calculates the total number of items that match the query, updates the provided PaginationMetadata struct with the total number of items, total number of pages, and current page number, and returns a GORM scope that can be used to fetch the results for the current page.
func Pagination(meta *PaginationMetadata, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalItems int64
	db.Count(&totalItems)

	meta.TotalItem = int(totalItems)
	totalPages := math.Ceil(float64(totalItems) / float64(meta.GetItemPerPage()))
	meta.TotalPage = int(totalPages)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(meta.GetOffset()).Limit(meta.ItemsPerPage)
	}
}

// FindOneByID returns a function that queries the entity with the given ID and returns the query result.
func FindOneByID[T Entity](id string, entity T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			First(&entity, "id = ?", id)
	}
}

// UpdateWithoutResult returns a function that updates the entity with the given ID using the given entity, but doesn't return the updated entity.
func UpdateWithoutResult[T Entity](id string, entity T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Where(id, "id = ?", id).
			Updates(&entity)
	}
}

// UpdateByIDWithResult returns a function that updates the entity with the given ID using the given entity, and returns the updated entity.
func UpdateByIDWithResult[T Entity](id string, entity T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Where(id, "id = ?", id).
			Updates(&entity).
			First(&entity, "id = ?", id)
	}
}

// DeleteWithResult returns a function that queries the entity with the given ID, deletes it, and returns the deleted entity.
func DeleteWithResult[T Entity](id string, entity T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			First(&entity, "id = ?", id).
			Delete(&entity, "id = ?", id)
	}
}

// DeleteWithoutResult returns a function that deletes the entity with the given ID using the given entity, but doesn't return the deleted entity.
func DeleteWithoutResult[T Entity](id string, entity T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Delete(&entity, "id = ?", id)
	}
}

type GormRepository[T Entity] interface {
	FindAll(metadata *PaginationMetadata, entities *[]T) error
	FindOne(id string, entity T, scope ...func(db *gorm.DB) *gorm.DB) error
	Create(entity T, scope ...func(db *gorm.DB) *gorm.DB) error
	Update(id string, entity T, scope ...func(db *gorm.DB) *gorm.DB) error
	Delete(id string, entity T, scope ...func(db *gorm.DB) *gorm.DB) error
	GetDB() *gorm.DB
}

type gormRepository[T Entity] struct {
	db *gorm.DB
}

// NewGormRepository function that create a new instance of gormRepository[T] with a GORM database connection
func NewGormRepository[T Entity](db *gorm.DB) GormRepository[T] {
	return &gormRepository[T]{
		db: db,
	}
}

func (r *gormRepository[T]) GetDB() *gorm.DB {
	return r.db
}

// FindAll the entities with pagination metadata and scopes.
// Pagination is achieved by using the Pagination function.
// The method updates the metadata to reflect the total number of items and the number of items on the current page.
func (r *gormRepository[T]) FindAll(metadata *PaginationMetadata, entities *[]T) error {
	if err := r.db.
		Scopes(Pagination(metadata, r.db)).
		Find(&entities).
		Error; err != nil {
		return err
	}

	metadata.ItemCount = len(*entities)
	return nil
}

// FindOne finds a single entity with the given id and optional scopes.
func (r *gormRepository[T]) FindOne(id string, entity T, scope ...func(db *gorm.DB) *gorm.DB) error {
	return r.db.
		Scopes(scope...).
		First(entity, "id = ?", id).
		Error
}

// Create a new entity in the database.
func (r *gormRepository[T]) Create(entity T, scope ...func(db *gorm.DB) *gorm.DB) error {
	return r.db.
		Scopes(scope...).
		Create(entity).
		Error
}

// Update an existing entity with the given id in the database.
// It returns an error if no entity with the given id is found.
func (r *gormRepository[T]) Update(id string, entity T, scope ...func(db *gorm.DB) *gorm.DB) error {
	return r.db.
		Scopes(scope...).
		Where(id, "id = ?", id).
		Updates(&entity).
		First(&entity, "id = ?", id).
		Error
}

// Delete an existing entity with the given id from the database.
// It returns an error if no entity with the given id is found.
func (r *gormRepository[T]) Delete(id string, entity T, scope ...func(db *gorm.DB) *gorm.DB) error {
	return r.db.
		Scopes(scope...).
		First(&entity, "id = ?", id).
		Delete(&entity).
		Error
}

// WithTransaction runs a list of functions inside a single transaction.
//
// Parameters:
// - fns: a list of functions that will be executed within a single transaction.
//
// Returns:
// - error: an error if any of the functions returns an error or the transaction commit fails, otherwise nil.
func (r *gormRepository[T]) WithTransaction(fns ...func(tx *gorm.DB) error) error {
	tx := r.db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		} else if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			panic(err)
		}
	}()

	for _, fn := range fns {
		if err := fn(tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}
