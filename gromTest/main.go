package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	// not null
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null"`
	Password     string `gorm:"not null"`
	Age          uint8  `gorm:"not null"`
	Birthday     *time.Time
	MemberNumber sql.NullString
	CompanyID    uint
	Company      Company
}

type Company struct {
	gorm.Model
	Name string
}

var db *gorm.DB

// 构造函数
func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	// gorm 连接数据库
	var err error
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:hxc@tcp(127.0.0.1:13306)/gorm_test?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                             // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                            // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                            // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                            // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,
		// 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	// 执行迁移
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}

type Student struct {
	Name   string `json:"name" struct:"name"`
	Age    int    `json:"age"`
	gender int    `json:"gender"`
}

func main() {
	var st2 Student
	stu1 := Student{
		Name:   "hxc",
		Age:    18,
		gender: 1,
	}
	fmt.Println(stu1.Name)
	s2Str := `{"Name":"resr", "age": 22, "gender": 21}`
	_ = json.Unmarshal([]byte(s2Str), &st2)
	fmt.Println(st2)

	a := 22
	if true {
		a = 44
		fmt.Println(a)
	}
	// 数据库测试
	fmt.Println("test")

	var m1 map[string]interface{}
	m1 = make(map[string]interface{})
	m1["name"] = "hxc"
	fmt.Println(m1)

	//birthday := time.Now()
	//fmt.Println(birthday)
	//userOne := User{
	//	Name:     "hxc",
	//	Email:    "173315279",
	//	Password: "123456",
	//	Age:      18,
	//	Birthday: &birthday,
	//	MemberNumber: sql.NullString{
	//		String: "123456",
	//		Valid:  true,
	//	},
	//}
	//// 创建单个
	//res := db.Create(&userOne)
	//fmt.Println(userOne.ID)
	//fmt.Println(res.RowsAffected)
	//fmt.Println(res.Error)
	//// 创建多个
	//users := []*User{
	//	{
	//		Name:     "hxc",
	//		Email:    "<EMAIL>",
	//		Password: "1234563",
	//		Age:      19,
	//		Birthday: &birthday,
	//		MemberNumber: sql.NullString{
	//			String: "123456",
	//			Valid:  true,
	//		},
	//	},
	//	{
	//		Name:     "hxc",
	//		Email:    "<EMAIL>",
	//		Password: "1234565",
	//		Age:      20,
	//		Birthday: &birthday,
	//		MemberNumber: sql.NullString{
	//			String: "123456",
	//			Valid:  true,
	//		},
	//	},
	//}
	//db.Create(users)
	//
	//var userss []User
	//// 设置默认值
	//userss = []User{
	//	{
	//		Name:     "hxc22",
	//		Email:    "<EMAIL>",
	//		Password: "1234563",
	//		Age:      19,
	//		Birthday: &birthday,
	//		MemberNumber: sql.NullString{
	//			String: "123456",
	//			Valid:  true,
	//		},
	//	},
	//	{
	//		Name:     "hxc33",
	//		Email:    "<EMAIL>",
	//		Password: "1234565",
	//		Age:      20,
	//		Birthday: &birthday,
	//		MemberNumber: sql.NullString{
	//			String: "123456",
	//			Valid:  true,
	//		},
	//	},
	//}
	//db.Create(userss)
	//user := User{}
	//db.First(&user)
	//fmt.Println(user.Name)
	//
	//// 用map返回
	//result := map[string]interface{}{}
	//db.Model(&User{}).First(&result)
	//fmt.Println(result)
	//fmt.Println(result["name"])
	//
	//// 用表名
	//result1 := map[string]interface{}{}
	//
	//db.Table("users").First(&result1)
	//fmt.Println(result1)
	//
	////
	//userTwo := []User{}
	//
	//res1 := db.Find(&userTwo, []int{1, 2, 3})
	//if res1.Error != nil {
	//	panic(res1.Error)
	//}
	//
	//// 遍历res1
	//for _, u := range userTwo {
	//	fmt.Println(u.Age)
	//}
	//
	//whereUser := User{}
	//// 带条件
	//db.Where("name = ?", "hxc").First(&whereUser)
	//
	//fmt.Println(whereUser.Name)
	//
	//whereUsers := []User{}
	//db.Where("name = ?", "hxc").Find(&whereUsers)
	//
	//for _, u := range whereUsers {
	//	fmt.Println(u.Age)
	//}
	//
	//// 查询条件 结构会忽略0值查询。map不会
	//userZero := []User{}
	//db.Where(&User{Name: "hxc", Age: 0}).Find(&userZero)
	//for _, u := range userZero {
	//	fmt.Println(u.Age)
	//}
	//userZeroMap := []User{}
	//
	//db.Where(map[string]interface{}{"name": "hxc", "age": 0}).Find(&userZeroMap)
	//for _, u := range userZeroMap {
	//	fmt.Println(u.Age)
	//}

	var users User
	//db.Select("name, age, id").Where("name = ?", "hxc33").First(&users)
	db.Select("name, age, id").Where("name = ?", "hxc33").First(&users)
	fmt.Println("--------test")
	fmt.Println(users)
	fmt.Println(users.ID)
	//for _, u := range users {
	//	fmt.Println(u)
	//	fmt.Println(u.Age)
	//	fmt.Println(u.ID)
	//}
	db.Model(&users).Update("age", 9)
	db.Delete(&users)
	db.Where("name = ?", "hxc33").Delete(&User{})
	//

	res := db.Model(&User{}).Update("age", 10)
	if res.RowsAffected > 0 {
		fmt.Println("update success")
	}
	if res.Error != nil {
		fmt.Println(res.Error)
	}
	resDel := db.Delete(&User{})
	if resDel.Error != nil {
		fmt.Println(resDel.Error)
	}

}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Age == 9 {
		fmt.Println("BeforeUpdate")
	}

	if tx.Statement.Changed("Age") {
		fmt.Println("Age changed")
	}
	return
}
