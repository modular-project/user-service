package model

import (
	"strings"
)

type Sort int

type By int

type Status int

const (
	ASC Sort = iota
	DES
)

const (
	NAME By = iota
	EMAIL
	ROL
	EST
)

const (
	ACTIVE Status = iota
	NOACTVIE
	ANY
)

type OrderBy struct {
	By   By
	Sort Sort
}

func (o OrderBy) get() string {
	var sort string
	if o.Sort == DES {
		sort = " DESC"
	}
	var order string
	switch o.By {
	case NAME:
		order = "name"
	case EMAIL:
		order = "email"
	case ROL:
		order = "role_id"
	case EST:
		order = "establishment_id"
	default:
		return ""
	}
	var b strings.Builder
	b.WriteString(order)
	b.WriteString(sort)
	return b.String()
}

type Search struct {
	OrderBys []OrderBy

	Limit  int
	Offset int
}

func (s Search) Query() string {
	var q strings.Builder
	for i, o := range s.OrderBys {
		if i != 0 {
			q.WriteString(",")
		}
		q.WriteString(o.get())
	}
	return q.String()
}

type SearchEMPL struct {
	Search
	Status
	Rols []uint
	Ests []uint
}
