package utils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/wolif/gosaber/pkg/mysql"
	"github.com/wolif/gosaber/pkg/mysql/gormhelper"
	"github.com/wolif/gosaber/pkg/ref"
	"strings"
)

type record struct {
	dbConnName string
}

var Record = &record{dbConnName: "default"}

func (r *record) SetDBConnName(name string) {
	r.dbConnName = name
}

func (r *record) Create(_ context.Context, record interface{}) error {
	db, err := mysql.GetDbByName(r.dbConnName)
	if err != nil {
		return err
	}

	res := db.Create(record)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *record) Update(_ context.Context, record interface{}) error {
	db, err := mysql.GetDbByName(r.dbConnName)
	if err != nil {
		return err
	}

	res := db.Save(record)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *record) UpdateField(_ context.Context, model interface{}, cond interface{}, data interface{}) (int64, error) {
	db, err := mysql.GetDbByName(r.dbConnName)
	if err != nil {
		return 0, err
	}

	dataUpdate := make(map[string]interface{})
	dataRef := ref.New(data)
	if dataRef.IsStruct() {
		fields, _ := dataRef.GetStructFields()
		if fields == nil || len(fields) == 0 {
			return 0, err
		}
		for fn, _ := range fields {
			jsonTag, _ := dataRef.GetStructFieldTag(fn, "json")
			if jsonTag == "" {
				jsonTag = fn
			}
			jsonTagSegs := strings.Split(jsonTag, ",")
			dataUpdate[strings.TrimSpace(jsonTagSegs[0])], _ = dataRef.GetStructFieldValue(fn)
		}
	} else if dataRef.IsMap() {
		for _, k := range dataRef.GetValue().MapKeys() {
			dataUpdate[k.String()], _ = dataRef.MapValGet(k.Interface())
		}
	} else {
		return 0, nil
	}

	res := gormhelper.Scope(db, cond).Model(model).Updates(dataUpdate)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (r *record) Delete(_ context.Context, cond interface{}, model interface{}) (int64, error) {
	db, err := mysql.GetDbByName(r.dbConnName)
	if err != nil {
		return 0, err
	}

	res := gormhelper.Scope(db, cond).Delete(model)
	if res.Error != nil {
		return 0, res.Error
	}

	return res.RowsAffected, nil
}

func (r *record) FindOne(_ context.Context, cond interface{}, model interface{}, result interface{}) error {
	db, err := mysql.GetDbByName(r.dbConnName)
	if err != nil {
		return err
	}

	res := gormhelper.Scope(db, cond).Model(model).Find(result)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *record) Search(_ context.Context, cond interface{}, model interface{}, page, pageSize int, result interface{}, order ...map[string]string) (count int64, err error) {
	db, err := mysql.GetDbByName(r.dbConnName)
	if err != nil {
		return 0, err
	}

	stat := gormhelper.Scope(db, cond).Model(model)
	cntRes := stat.Count(&count)
	if cntRes.Error != nil {
		return 0, cntRes.Error
	}

	if len(order) > 0 {
		for field, sort := range order {
			stat = stat.Order(fmt.Sprintf("%s %s", field, sort))
		}
	}

	res := gormhelper.Paginate(stat, page, pageSize).Find(result)
	if res.Error != nil {
		return 0, res.Error
	}
	return
}

func (r *record) SearchRaw(_ context.Context, cond interface{}, model interface{}, pagination ...int) (cnt int64, rows *sql.Rows, err error) {
	db, err := mysql.GetDbByName(r.dbConnName)
	if err != nil {
		return 0, nil, err
	}

	stat := gormhelper.Scope(db, cond).Model(model)
	if len(pagination) == 2 && pagination[0] > 1 && pagination[1] > 1 {
		cntRes := stat.Count(&cnt)
		if cntRes.Error != nil {
			return 0, nil, cntRes.Error
		}
		stat = gormhelper.Paginate(stat, pagination[0], pagination[1])
	}

	rows, err = stat.Rows()
	if err != nil {
		return 0, nil, err
	}
	return
}
