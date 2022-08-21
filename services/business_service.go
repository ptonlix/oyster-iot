package services

import (
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/core/logs"
)

type BusinessService struct {
}

// 插入一条业务数据
func (b *BusinessService) Add(business *models.Business) (int64, error) {
	id, err := mysql.Mydb.Insert(business)
	if err != nil {
		logs.Warn(err)
		return id, err
	}
	logs.Info("Insert business successful! ID:", id)
	return id, err
}

// 更新一条业务数据
func (b *BusinessService) Update(business *models.Business) error {
	id, err := mysql.Mydb.Update(business)
	if err != nil {
		logs.Warn(err)
		return err
	}
	logs.Info("Update business successful! ID:", id)
	return err
}

// 删除一条业务数据
func (b *BusinessService) Delete(business *models.Business) (err error) {
	// 获取该业务关联的设备列表
	var devices []*models.Device
	qs := mysql.Mydb.QueryTable(&models.Device{})
	_, err = qs.Filter("business_id", business.Id).All(&devices)
	if err != nil {
		logs.Error(err)
		return
	}
	// 开启事务
	to, err := mysql.Mydb.Begin()
	if err != nil {
		logs.Error("start the transaction failed")
		return

	}

	defer func() {
		if err != nil {
			err = to.Rollback()
			if err != nil {
				logs.Error("roll back transaction failed", err)
			}
		} else {
			err = to.Commit()
			if err != nil {
				logs.Error("commit transaction failed.", err)
			}
		}
	}()

	for _, v := range devices {
		v.BusinessId = 0 //更新设备的业务ID为默认值
		_, err = to.Update(v, "BusinessId")
		if err != nil {
			logs.Error("execute transaction's select sql fail, rollback.", err)
			return
		}
	}
	id, err := to.Delete(business)
	if err != nil {
		logs.Warn(err)
		return
	}
	logs.Info(" Delete business successful! ID:", id)
	return
}

// 获取全部设备
func (*BusinessService) GetBusinessByPage(userId, pageSize, pageNum int) (int, int, []*models.Business, error) {

	var business []*models.Business
	qs := mysql.Mydb.QueryTable(&models.Business{})

	totalRecord, err := qs.Count()
	if err != nil {
		logs.Warn(err)
		return -1, -1, nil, err

	}
	num, err := qs.Filter("user_id", userId).Limit(pageSize, pageSize*(pageNum-1)).All(&business)

	if err != nil {
		logs.Warn(err)
		return -1, -1, nil, err
	}

	totalPageNum := (int(totalRecord) + pageSize - 1) / pageSize

	logs.Info("Get Business successful! Totalcount: %v TotalPages: %v Returned Rows Num: %#v", totalRecord, totalPageNum, num)

	return int(totalRecord), totalPageNum, business, err
}

// 获取全部设备
func (*BusinessService) GetBusinessAllNum() (int, error) {

	qs := mysql.Mydb.QueryTable(&models.Business{})

	totalRecord, err := qs.Count()
	if err != nil {
		logs.Warn(err)
		return -1, err

	}

	logs.Info("Get Business successful! Totalcount: %v ", totalRecord)

	return int(totalRecord), err
}
