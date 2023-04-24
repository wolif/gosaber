package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

type Pl struct {
	data []bson.D
}

var Pipeline = new(Pl)

func (p *Pl) New() *Pl {
	return &Pl{data: make([]bson.D, 0)}
}

func (p *Pl) AddFields(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$addFields",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Bucket(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$bucket",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) ColStats(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$colStats",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Count(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$count",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) CurrentOp(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$currentOp",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Facet(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$facet",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) GeoNear(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$geoNear",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) GraphLookup(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$graphLookup",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Group(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$group",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) IndexStats(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$indexStats",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Limit(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$limit",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) ListLocalSessions(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$listLocalSessions",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) ListSessions(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$listSessions",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Lookup(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$lookup",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Match(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$match",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Out(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$out",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Project(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$project",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Redact(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$redact",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Sample(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$sample",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) ReplaceRoot(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$replaceRoot",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Skip(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$skip",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Sort(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$sort",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) SortByCount(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$sortByCount",
			Value: cond,
		},
	})
	return p
}

func (p *Pl) Unwind(cond interface{}) *Pl {
	p.data = append(p.data, bson.D{
		bson.E{
			Key:   "$unwind",
			Value: cond,
		},
	})
	return p
}

// --------------------------------------------------------

func (p *Pl) Paginate(page, pageSize int) *Pl {
	return p.Skip((page - 1) * pageSize).Limit(pageSize)
}

func (p *Pl) Extract() mongo2.Pipeline {
	pl := mongo2.Pipeline{}
	for _, d := range p.data {
		pl = append(pl, d)
	}
	return pl
}
