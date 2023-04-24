package shortcut

import (
	"context"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/wolif/gosaber/pkg/mysql/helper"
	"github.com/wolif/gosaber/pkg/ref"
)

func Create(_ context.Context, record interface{}, txDB ...*gorm.DB) error {
	db, err := resolveDB(txDB...)
	if err != nil {
		return err
	}

	res := db.Create(record)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func Save(_ context.Context, record interface{}, txDB ...*gorm.DB) error {
	db, err := resolveDB(txDB...)
	if err != nil {
		return err
	}

	res := db.Save(record)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func Update(_ context.Context, model interface{}, cond interface{}, data interface{}, txDB ...*gorm.DB) (int64, error) {
	db, err := resolveDB(txDB...)
	if err != nil {
		return 0, err
	}

	dataUpdate := make(map[string]interface{})
	rd := ref.New(data)
	if rd.IsStruct() {
		fields, _ := rd.StructFieldsName()
		for _, fieldName := range fields {
			jsonTag, _ := rd.StructTagGet(fieldName, "json")
			if jsonTag == "" {
				jsonTag = fieldName
			}
			jsonTagSegs := strings.Split(jsonTag, ",")
			dataUpdate[strings.TrimSpace(jsonTagSegs[0])], _ = rd.StructValueGet(fieldName)
		}
	} else if rd.IsMap() {
		for _, k := range rd.GetValue().MapKeys() {
			dataUpdate[k.String()], _ = rd.MapGet(k.Interface())
		}
	} else {
		return 0, nil
	}

	res := helper.WithWhere(db, cond).Model(model).Updates(dataUpdate)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func Delete(_ context.Context, cond interface{}, model interface{}, txDB ...*gorm.DB) (int64, error) {
	db, err := resolveDB(txDB...)
	if err != nil {
		return 0, err
	}

	res := helper.WithWhere(db, cond).Delete(model)
	if res.Error != nil {
		return 0, res.Error
	}

	return res.RowsAffected, nil
}

func FindOne(_ context.Context, cond interface{}, model interface{}, result interface{}, txDB ...*gorm.DB) error {
	db, err := resolveDB(txDB...)
	if err != nil {
		return err
	}

	res := helper.WithWhere(db, cond).Model(model).Find(result)
	if res.Error != nil {
		if gorm.IsRecordNotFoundError(res.Error) {
			return nil
		}
		return res.Error
	}
	return nil
}

func Find(ctx context.Context, model interface{}, cond interface{}, page, pageSize int, result interface{}, sort []string, txDB ...*gorm.DB) (count int64, err error) {
	db, err := resolveDB(txDB...)
	if err != nil {
		return 0, err
	}

	// 查数量
	count, err = Count(ctx, model, cond, txDB...)
	if err != nil {
		return
	}

	// 确定找不到数据的话, 直接返回
	page, pageSize = resolvePagi(page, pageSize)
	if count == 0 || int64((page-1)*pageSize) > count {
		return
	}

	res := helper.WithPaginate(helper.WithSort(db, sort...), page, pageSize).Find(result)
	if res.Error != nil {
		if gorm.IsRecordNotFoundError(res.Error) {
			return 0, nil
		}
		return 0, res.Error
	}
	return
}

func Count(_ context.Context, model interface{}, cond interface{}, txDB ...*gorm.DB) (count int64, err error) {
	db, err := resolveDB(txDB...)
	if err != nil {
		return 0, err
	}

	stat := helper.WithWhere(db, cond).Model(model)
	cntRes := stat.Count(&count)
	if cntRes.Error != nil {
		return 0, cntRes.Error
	}
	return
}
