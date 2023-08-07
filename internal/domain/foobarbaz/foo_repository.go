package foobarbaz

//go:generate go run github.com/golang/mock/mockgen -source foo_repository.go -destination mock/foo_repository_mock.go -package foobarbaz_mock

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	fooQueries = struct {
		selectFoo                    string
		selectFooItem                string
		insertFoo                    string
		insertFooItemBulk            string
		insertFooItemBulkPlaceholder string
		updateFoo                    string
	}{
		selectFoo: `
			SELECT
				foo.entity_id,
				foo.name,
				foo.total_quantity,
				foo.total_price,
				foo.total_discount,
				foo.shipping_fee,
				foo.grand_total,
				foo.status,
				foo.created,
				foo.created_by,
				foo.updated,
				foo.updated_by,
				foo.deleted,
				foo.deleted_by
			FROM foo `,

		selectFooItem: `
			SELECT
				entity_id,
				foo_id,
				sku,
				product_name,
				quantity,
				unit_price,
				total_price,
				discount,
				grand_total
			FROM foo_item`,

		insertFoo: `
			INSERT INTO foo (
				entity_id,
				name,
				total_quantity,
				total_price,
				total_discount,
				shipping_fee,
				grand_total,
				status,
				created,
				created_by,
				updated,
				updated_by,
				deleted,
				deleted_by
			) VALUES (
				:entity_id,
				:name,
				:total_quantity,
				:total_price,
				:total_discount,
				:shipping_fee,
				:grand_total,
				:status,
				:created,
				:created_by,
				:updated,
				:updated_by,
				:deleted,
				:deleted_by)`,

		insertFooItemBulk: `
			INSERT INTO foo_item (
				entity_id,
				foo_id,
				sku,
				product_name,
				quantity,
				unit_price,
				total_price,
				discount,
				grand_total
			) VALUES `,

		insertFooItemBulkPlaceholder: `
			(:entity_id,
			:foo_id,
			:sku,
			:product_name,
			:quantity,
			:unit_price,
			:total_price,
			:discount,
			:grand_total)`,

		updateFoo: `
			UPDATE foo
			SET
				name = :name,
				total_quantity = :total_quantity,
				total_price = :total_price,
				total_discount = :total_discount,
				shipping_fee = :shipping_fee,
				grand_total = :grand_total,
				status = :status,
				created = :created,
				created_by = :created_by,
				updated = :updated,
				updated_by = :updated_by,
				deleted = :deleted,
				deleted_by = :deleted_by
			WHERE entity_id = :entity_id `,
	}
)

// FooRepository is the repository for Foo data.
type FooRepository interface {
	Create(foo Foo) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveByID(id uuid.UUID) (foo Foo, err error)
	ResolveItemsByFooIDs(ids []uuid.UUID) (fooItems []FooItem, err error)
	Update(foo Foo) (err error)
}

// FooRepositoryMySQL is the MySQL-backed implementation of FooRepository.
type FooRepositoryMySQL struct {
	DB *infras.MySQLConn
}

// ProvideFooRepositoryMySQL is the provider for this repository.
func ProvideFooRepositoryMySQL(db *infras.MySQLConn) *FooRepositoryMySQL {
	s := new(FooRepositoryMySQL)
	s.DB = db
	return s
}

// Create creates a new Foo.
func (r *FooRepositoryMySQL) Create(foo Foo) (err error) {
	exists, err := r.ExistsByID(foo.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if exists {
		err = failure.Conflict("create", "foo", "already exists")
		logger.ErrorWithStack(err)
		return
	}

	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreate(tx, foo); err != nil {
			e <- err
			return
		}

		if err := r.txCreateItems(tx, foo.Items); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

// ExistsByID checks the existence of a Foo by its ID.
func (r *FooRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Read.Get(
		&exists,
		"SELECT COUNT(entity_id) FROM foo WHERE foo.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

// ResolveByID resolves a Foo by its ID
func (r *FooRepositoryMySQL) ResolveByID(id uuid.UUID) (foo Foo, err error) {
	err = r.DB.Read.Get(
		&foo,
		fooQueries.selectFoo+" WHERE foo.entity_id = ?",
		id.String())
	if err != nil && err == sql.ErrNoRows {
		err = failure.NotFound("foo")
		logger.ErrorWithStack(err)
		return
	}
	return
}

// ResolveItemsByFooIDs resolves FooItems based on a set of FooIDs.
func (r *FooRepositoryMySQL) ResolveItemsByFooIDs(ids []uuid.UUID) (fooItems []FooItem, err error) {
	if len(ids) == 0 {
		return
	}

	query, args, err := sqlx.In(fooQueries.selectFooItem+" WHERE foo_item.foo_id IN (?)", ids)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	err = r.DB.Read.Select(&fooItems, query, args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

// Update updates a Foo.
func (r *FooRepositoryMySQL) Update(foo Foo) (err error) {
	exists, err := r.ExistsByID(foo.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if !exists {
		err = failure.NotFound("foo")
		logger.ErrorWithStack(err)
		return
	}

	// transactionally update the Foo
	// strategy:
	// 1. delete all the Foo's items
	// 2. create a new set of Foo's items
	// 3. update the Foo
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txDeleteItems(tx, foo.ID); err != nil {
			e <- err
			return
		}

		if err := r.txCreateItems(tx, foo.Items); err != nil {
			e <- err
			return
		}

		if err := r.txUpdate(tx, foo); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

// internal methods

// composeBulkInsertItemQuery composes a bulk insert item query given a slice of FooItems.
func (r *FooRepositoryMySQL) composeBulkInsertItemQuery(fooItems []FooItem) (query string, params []interface{}, err error) {
	values := []string{}
	for _, fi := range fooItems {
		param := map[string]interface{}{
			"entity_id":    fi.ID,
			"foo_id":       fi.FooID,
			"sku":          fi.SKU,
			"product_name": fi.ProductName,
			"quantity":     fi.Quantity,
			"unit_price":   fi.UnitPrice,
			"total_price":  fi.TotalPrice,
			"discount":     fi.Discount,
			"grand_total":  fi.GrandTotal,
		}
		q, args, err := sqlx.Named(fooQueries.insertFooItemBulkPlaceholder, param)
		if err != nil {
			return query, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}
	query = fmt.Sprintf("%v %v", fooQueries.insertFooItemBulk, strings.Join(values, ","))
	return
}

// txCreate creates a Foo transactionally given the *sqlx.Tx param.
func (r *FooRepositoryMySQL) txCreate(tx *sqlx.Tx, foo Foo) (err error) {
	stmt, err := tx.PrepareNamed(fooQueries.insertFoo)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(foo)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

// txCreateItems create FooItems transactionally given the *sqlx.Tx param.
func (r *FooRepositoryMySQL) txCreateItems(tx *sqlx.Tx, fooItems []FooItem) (err error) {
	if len(fooItems) == 0 {
		return
	}

	query, args, err := r.composeBulkInsertItemQuery(fooItems)
	if err != nil {
		return
	}

	stmt, err := tx.Preparex(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Stmt.Exec(args...)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

// txDeleteeItems deletes FooItems based on their FooID transactionally given the *sqlx.Tx param.
func (r *FooRepositoryMySQL) txDeleteItems(tx *sqlx.Tx, fooID uuid.UUID) (err error) {
	_, err = tx.Exec("DELETE FROM foo_item WHERE foo_id = ?", fooID.String())
	return
}

// txUpdate updates a Foo transactionally, given the *sqlx.Tx param.
func (r *FooRepositoryMySQL) txUpdate(tx *sqlx.Tx, foo Foo) (err error) {
	stmt, err := tx.PrepareNamed(fooQueries.updateFoo)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(foo)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
