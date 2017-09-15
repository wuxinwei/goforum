package services

import (
	"fmt"

	"github.com/fpay/gopress"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"github.com/wuxinwei/goforum/conf"
	"github.com/wuxinwei/goforum/models"
)

const (
	// DbServiceName is the identity of db service
	DbServiceName = "db"
)

// Pair is a type that provide a search argument for query operation
// for instance, if you wanna to call query operation by db service, mysql clause like below:
// SELECT * FROM xxx WHERE a = xxx AND b = xxx;
// you should define this pair,
// Pair{Name: a, Value: xxx}, Pair{Name b, Value: xxxx}
type Pair struct {
	Name  string
	Value interface{}
}

// DbService type
type DbService struct {
	// c  *gopress.Container
	*gorm.DB
}

// NewDbService returns instance of db service
func NewDbService() *DbService {
	// TODO 设置时区
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&&parseTime=true",
		conf.GlobalConf.DB.Username,
		conf.GlobalConf.DB.Password,
		conf.GlobalConf.DB.Address,
		conf.GlobalConf.DB.Port,
		conf.GlobalConf.DB.Name)
	db, err := gorm.Open("mysql", addr)
	if err != nil {
		logrus.Errorf("mysql connection failed, address: %s, Err: %s", addr, err)
		return nil
	}
	db.LogMode(conf.GlobalConf.DB.Debug)
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Tag{}, &models.Comment{})
	db.Model(&models.User{}).Related(&models.Post{})
	db.Model(&models.User{}).Related(&models.Comment{})
	//db.Model(&models.Post{}).Related(&models.Comment{})
	//db.Model(&models.Post{}).Related(&models.Tag{})
	return &DbService{
		db,
	}
}

// ServiceName is used to implements gopress.Service
func (s *DbService) ServiceName() string {
	return DbServiceName
}

// RegisterContainer is used to implements gopress.Service
func (s *DbService) RegisterContainer(c *gopress.Container) {
	//s.c = c
}
