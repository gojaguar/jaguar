package repository

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"gorm.io/gorm"
)

type Page struct {
	// Index is the page number.
	Index int
	// Size is the current page size.
	Size int
}

// Limit converts a human-readable page size into a database-readable limit.
func (p Page) Limit() int {
	if p.Size < 1 {
		return 1
	}
	return p.Size
}

// Offset converts a human-readable page number into a database-readable offset.
func (p Page) Offset() int {
	if p.Index < 1 {
		return 0
	}
	return (p.Index - 1) * p.Size
}

// OrderBy allows to order results from a database in a certain order by using a certain column as reference.
type OrderBy struct {
	// Column is the colum that will be used as reference when ordering results.
	Column string
	// Desc is set to true the results should be in descending order. If not, it defaults to ascending order.
	Desc bool
}

// SQL converts the current OrderBy to a SQL expression.
func (o OrderBy) SQL() string {
	order := "desc"
	if !o.Desc {
		order = "asc"
	}
	query := fmt.Sprintf("%s %s", o.Column, order)
	return query
}

// Firestore is a helper method to convert the current OrderBy to a firestore expression.
func (o OrderBy) Firestore() (string, firestore.Direction) {
	dir := firestore.Asc
	if o.Desc {
		dir = firestore.Desc
	}
	return o.Column, dir
}

type Query struct {
	page        Page
	orderBy     OrderBy
	groupColumn string
}

func (q Query) GORM(tx *gorm.DB) *gorm.DB {
	if q.page.Size > 0 {
		tx = tx.Limit(q.page.Limit())
	}
	if q.page.Index >= 0 {
		tx = tx.Offset(q.page.Offset())
	}
	if len(q.orderBy.Column) > 0 {
		tx = tx.Order(q.orderBy.SQL())
	}
	if len(q.groupColumn) > 0 {
		tx = tx.Group(q.groupColumn)
	}
	return tx
}

func (q Query) Firestore(qf firestore.Query) firestore.Query {
	if q.page.Size > 0 {
		qf = qf.Limit(q.page.Limit())
	}
	if q.page.Index >= 0 {
		qf = qf.Offset(q.page.Offset())
	}
	if len(q.orderBy.Column) > 0 {
		qf = qf.OrderBy(q.orderBy.Firestore())
	}
	return qf
}
