package services

import (
	"fmt"
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type OperLogService struct {
}

// 添加操作日志
type SearchOperlogParam struct {
	Starttime     string `json:"starttime,omitempy" valid:"MaxSize(255)"`
	Endtime       string `json:"endtime,omitempy" valid:"MaxSize(255)"`
	SearchImei    string `json:"imei,omitempy" valid:"MaxSize(255)"`
	SearchContext string `json:"context,omitempy" valid:"MaxSize(255)"`
}

func (*OperLogService) Add(operlog *models.Operlog) error {

	id, err := mysql.Mydb.Insert(operlog)
	if err != nil {
		logs.Warn(err)
		return err
	}
	logs.Info("Insert Operlog successful! ID:", id)
	return err
}

// 修改操作日志
func (*OperLogService) Update(operlog *models.Operlog) error {

	id, err := mysql.Mydb.Update(operlog)
	if err != nil {
		logs.Warn(err)
		return err
	}
	logs.Info("Update Operlog successful! ID:", id)
	return err
}

// 删除操作日志
func (*OperLogService) Delete(operlog *models.Operlog) error {

	id, err := mysql.Mydb.Delete(operlog)
	if err != nil {
		logs.Warn(err)
		return err
	}

	logs.Info("Delete Operlog successful! ID:", id)
	return err
}

// 可以支持关键字搜索查询
func (*OperLogService) GetOperlogBySearch(pageSize, pageNum int, params *SearchOperlogParam) (int, int, []*models.Operlog, error) {

	oper := []*models.Operlog{}

	SQL := ""
	if params.Starttime != "" && params.Endtime != "" {
		if params.SearchImei != "" {
			SQL = fmt.Sprintf("SELECT * FROM operlog WHERE createdat between '%s' and '%s' and INSTR(%s, '%s')>0 ", params.Starttime, params.Endtime, params.SearchImei, params.SearchContext)
		} else {
			SQL = fmt.Sprintf("SELECT * FROM operlog WHERE createdat between '%s' and '%s' ", params.Starttime, params.Endtime) //模糊查询
		}
	} else {
		if params.SearchImei != "" {
			SQL = fmt.Sprintf("SELECT * FROM operlog WHERE INSTR(%s, '%s')>0  ", params.SearchImei, params.SearchContext)
		} else {
			SQL = "SELECT * FROM operlog "
		}
	}

	totalRecord, err := mysql.Mydb.Raw(SQL).QueryRows(&oper)

	if err == orm.ErrNoRows {
		logs.Warn("Get SQL:%#v  ErrNoRows!", SQL)
		return 0, 0, nil, err
	} else if err != nil {
		logs.Warn("Get SQL:%#v  Failed! err:%#v", SQL, err)
		return 0, 0, nil, err
	}
	totalPageNum := (int(totalRecord) + pageSize - 1) / pageSize

	num, err := mysql.Mydb.Raw(SQL+"LIMIT ? OFFSET  ?", pageSize, pageSize*(pageNum-1)).QueryRows(&oper)

	if err == orm.ErrNoRows {
		logs.Warn("Get SQL:%#v  ErrNoRows!", SQL)
		return 0, 0, nil, err
	} else if err != nil {
		logs.Warn("Get SQL:%#v  Failed! err:%#v", SQL, err)
		return 0, 0, nil, err
	}

	logs.Info("Get operlog successful! Totalcount: %v TotalPages: %v Returned Rows Num: %#v", totalRecord, totalPageNum, num)
	return int(totalRecord), totalPageNum, oper, err
}

// 可以支持关键字搜索查询
func (*OperLogService) GetOperlogByPage(pageSize, pageNum int) (int, int, []*models.Operlog, error) {

	oper := []*models.Operlog{}

	SQL := "SELECT * FROM operlog "

	totalRecord, err := mysql.Mydb.Raw(SQL).QueryRows(&oper)

	if err == orm.ErrNoRows {
		logs.Warn("Get SQL:%#v  ErrNoRows!", SQL)
		return 0, 0, nil, err
	} else if err != nil {
		logs.Warn("Get SQL:%#v  Failed! err:%#v", SQL, err)
		return 0, 0, nil, err
	}
	totalPageNum := (int(totalRecord) + pageSize - 1) / pageSize

	num, err := mysql.Mydb.Raw(SQL+"LIMIT ? OFFSET  ?", pageSize, pageSize*(pageNum-1)).QueryRows(&oper)

	if err == orm.ErrNoRows {
		logs.Warn("Get SQL:%#v  ErrNoRows!", SQL)
		return 0, 0, nil, err
	} else if err != nil {
		logs.Warn("Get SQL:%#v  Failed! err:%#v", SQL, err)
		return 0, 0, nil, err
	}

	logs.Info("Get operlog successful! Totalcount: %v TotalPages: %v Returned Rows Num: %#v", totalRecord, totalPageNum, num)
	return int(totalRecord), totalPageNum, oper, err
}
