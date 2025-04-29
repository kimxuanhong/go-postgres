package postgres

import (
	"context"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Operations defines common database operations.
type Operations interface {
	Select(ctx context.Context, model interface{}, query interface{}, args ...interface{}) error
	SelectOne(ctx context.Context, model interface{}, query interface{}, args ...interface{}) error
	Insert(ctx context.Context, model interface{}) error
	Update(ctx context.Context, model interface{}, query interface{}, updates interface{}, args ...interface{}) error
	Delete(ctx context.Context, model interface{}, query interface{}, args ...interface{}) error
	DeleteOne(ctx context.Context, model interface{}, query interface{}, args ...interface{}) error
	Count(ctx context.Context, model interface{}, query interface{}, args ...interface{}) (int64, error)
	Exists(ctx context.Context, model interface{}, query interface{}, args ...interface{}) (bool, error)
	Close() error
}

type Postgres struct {
	db *gorm.DB
}

// NewPostgres initializes and returns a new Postgres instance.
//
// Example:
//
//	cfg := &Config{}
//	pg, err := NewPostgres(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer pg.Close()
func NewPostgres(cfg *Config) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect database: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to PostgresSQL database")
	return &Postgres{db: db}, nil
}

// Select retrieves multiple records.
//
// Example:
//
//	var users []User
//	err := pg.Select(ctx, &users, "age > ?", 18)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (p *Postgres) Select(ctx context.Context, model interface{}, query interface{}, args ...interface{}) error {
	return p.db.WithContext(ctx).Where(query, args...).Find(model).Error
}

// SelectOne retrieves a single record.
//
// Example:
//
//	var user User
//	err := pg.SelectOne(ctx, &user, "email = ?", "someone@example.com")
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//	    log.Println("Không tìm thấy user")
//	} else if err != nil {
//	    log.Fatal(err)
//	}
func (p *Postgres) SelectOne(ctx context.Context, model interface{}, query interface{}, args ...interface{}) error {
	return p.db.WithContext(ctx).Where(query, args...).First(model).Error
}

// Insert creates a new record.
//
// Example:
//
//	user := User{
//	    Name: "John Doe",
//	    Email: "john@example.com",
//	    Age: 25,
//	}
//	err := pg.Insert(ctx, &user)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (p *Postgres) Insert(ctx context.Context, model interface{}) error {
	return p.db.WithContext(ctx).Create(model).Error
}

// Update updates existing records.
//
// Example:
//
//	updates := map[string]interface{}{
//	    "name": "New Name",
//	    "age":  30,
//	}
//	err := pg.Update(ctx, &User{}, "id = ?", updates, 1)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (p *Postgres) Update(ctx context.Context, model interface{}, query interface{}, updates interface{}, args ...interface{}) error {
	return p.db.WithContext(ctx).Model(model).Where(query, args...).Updates(updates).Error
}

// Delete removes records matching the query.
//
// Example:
//
//	err := pg.Delete(ctx, &User{}, "age < ?", 18)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (p *Postgres) Delete(ctx context.Context, model interface{}, query interface{}, args ...interface{}) error {
	return p.db.WithContext(ctx).Where(query, args...).Delete(model).Error
}

// DeleteOne deletes the first record matching the query.
//
// Example:
//
//	var user User
//	err := pg.DeleteOne(ctx, &user, "email = ?", "someone@example.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (p *Postgres) DeleteOne(ctx context.Context, model interface{}, query interface{}, args ...interface{}) error {
	tx := p.db.WithContext(ctx).Where(query, args...).First(model)
	if tx.Error != nil {
		return tx.Error
	}
	return p.db.WithContext(ctx).Delete(model).Error
}

// Count counts the number of matching records.
//
// Example:
//
//	count, err := pg.Count(ctx, &User{}, "age > ?", 30)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	log.Printf("Có %d users trên 30 tuổi", count)
func (p *Postgres) Count(ctx context.Context, model interface{}, query interface{}, args ...interface{}) (int64, error) {
	var count int64
	err := p.db.WithContext(ctx).Model(model).Where(query, args...).Count(&count).Error
	return count, err
}

// Exists checks if any records match the query.
//
// Example:
//
//	exists, err := pg.Exists(ctx, &User{}, "email = ?", "someone@example.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if exists {
//	    log.Println("User tồn tại")
//	} else {
//	    log.Println("User không tồn tại")
//	}
func (p *Postgres) Exists(ctx context.Context, model interface{}, query interface{}, args ...interface{}) (bool, error) {
	count, err := p.Count(ctx, model, query, args...)
	return count > 0, err
}

// WithTransaction executes multiple operations in a single transaction.
//
// Example:
//
//	err := pg.WithTransaction(func(tx *gorm.DB) error {
//	    if err := tx.Create(&user).Error; err != nil {
//	        return err
//	    }
//	    if err := tx.Create(&order).Error; err != nil {
//	        return err
//	    }
//	    return nil
//	})
func (p *Postgres) WithTransaction(fn func(tx *gorm.DB) error) error {
	return p.db.Transaction(fn)
}

// Close closes the database connection.
//
// Example:
//
//	defer func() {
//	    if err := pg.Close(); err != nil {
//	        log.Printf("Error closing database: %v", err)
//	    }
//	}()
func (p *Postgres) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// DB returns the underlying *gorm.DB instance.
//
// Example:
//
//	db := pg.DB()
//	err := db.Exec("RAW SQL HERE").Error
func (p *Postgres) DB() *gorm.DB {
	return p.db
}
