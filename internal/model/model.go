package model

import (
	"fmt"
	"time"

	"github.com/go-web/global"
	"github.com/go-web/pkg/settings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

func NewDBEngine(databaseSetting *settings.DatabaseSettingS) (*gorm.DB, error) {
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}
	db.SingularTable(true)
	// 使用新函数`updateTimeStampForCreateCallback`替换回调`gorm:update_time_stamp`用于创建过程
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	// 设置最大连接数量
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	// 设置池最大空闲连接数
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}

// 处理model回调
// 新增行为的回调
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		// scope.FieldByName 方法，获取当前是否包含所需的字段
		if createTimeField, ok := scope.FieldByName("CreateOn"); ok {
			// Field.IsBlank 的值，可以得知该字段的值是否为空。
			if createTimeField.IsBlank {
				// 若为空，则会调用 Field.Set 方法给该字段设置值
				_ = createTimeField.Set(nowTime)
			}
		}

		if modifiedTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifiedTimeField.IsBlank {
				_ = modifiedTimeField.Set(nowTime)
			}
		}
	}
}

// 更新行为的回调
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	// scope.Get("gorm:update_column") 去获取当前设置了标识 gorm:update_column 的字段属性
	if _, ok := scope.Get("gorm:update_column"); !ok {
		// 若不存在，也就是没有自定义设置 update_column，那么将会在更新回调内设置默认字段 ModifiedOn 的值为当前的时间戳
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// 删除行为的回调
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		// 判断是否存在 DeletedOn 和 IsDel 字段，
		// 若存在则调整为执行 UPDATE 操作进行软删除（修改 DeletedOn 和 IsDel 的值），否则执行 DELETE 进行硬删除。
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")
		if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField {
			now := time.Now().Unix()
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v, %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
