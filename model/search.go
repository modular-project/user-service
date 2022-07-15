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
	By   By   `json:"by,omitempty"`
	Sort Sort `json:"sort,omitempty"`
}
type Search struct {
	OrderBys []OrderBy `json:"order,omitempty"`

	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type SearchEMPL struct {
	Search
	Status `json:"status"`
	Rols   []uint `json:"roles,omitempty"`
	Ests   []uint `json:"ests,omitempty"`
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

func (s Search) Query() string {
	var q strings.Builder
	for i, o := range s.OrderBys {
		g := o.get()
		if i != 0 && g != "" {
			q.WriteString(",")
		}
		q.WriteString(g)
	}
	return q.String()
}
