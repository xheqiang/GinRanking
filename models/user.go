package models

import (
	"fmt"
	"ginRanking/util/logger"
	"time"
)

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	/* CreatedAt CustomTime `json:"created_at"`
	UpdatedAt CustomTime `json:"updated_at"` */
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomTime time.Time

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	formatted := fmt.Sprintf(`"%s"`, t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (User) TableName() string {
	return "user"
}

func GetUserInfoById(id int) (User, error) {
	var user User

	if DB == nil {
		logger.Error(map[string]interface{}{
			"mysql connect error": "database connection is not initialized",
		})
		return User{}, nil
	}

	err := DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, err
}

func GetAllUserList() ([]User, error) {

	var users []User
	err := DB.Where("1 = ?", 1).Find(&users).Error

	if err != nil {
		return users, err
	}

	return users, err
}

func GetUserInfoByUserName(userName string) (User, error) {
	var user User

	err := DB.Where("user_name = ?", userName).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, err
}

func GetUserByUserName(userName string) (User, error) {
	var user User

	err := DB.Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, err
}

func AddUser(user_name, password string) (int, error) {
	createdAt := time.Now()
	updatedAt := time.Now()

	user := User{
		UserName:  user_name,
		Password:  password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	err := DB.Create(&user).Error
	return user.Id, err
}

func UpdateUserName(id int, userName string) error {
	err := DB.Model(&User{}).Where("id = ?", id).Update("user_name", userName).Error
	return err
}

func DeleteUserById(id int) (int, error) {
	var user User
	err := DB.Where("id = ?", id).Delete(&user).Error
	// 其它方法删除
	//err := DB.Delete(&User{}, 10).Error
	return id, err
}

/*
 *
 * - 时间格式原样输出 2024-07-18 00:00:00 而非 2024-07-18T00:00:00+08:00
 * - 创建自定义类型 type CustomTime time.Time
 * - CreatedAt time.Time `json:"created_at"` => CreatedAt CustomTime `json:"created_at"`
 * - 创建CustomTime格式化方法
 *   func (ct CustomTime) MarshalJSON() ([]byte, error) {
 *       t := time.Time(ct)
 *       formatted := fmt.Sprintf(`"%s"`, t.Format("2006-01-02 15:04:05"))
 *       return []byte(formatted), nil
 *   }
 */

// 其它方法处理 user.created_at 和 user.updated_at
/* func formatUser(u User) map[string]interface{} {
	return map[string]interface{}{
		"id":         u.Id,
		"user_name":  u.UserName,
		"password":   u.Password,
		"created_at": u.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
} */

/* func GetUserInfoById(id int) (map[string]interface{}, error) {
	var user User

	if DB == nil {
		logger.Error(map[string]interface{}{
			"mysql connect error": "database connection is not initialized",
		})
		return map[string]interface{}{}, nil
	}

	err := DB.Where("id = ?", id).First(&user).Error
	formatUser := formatUser(user)
	return formatUser, err
} */

/* func GetAllUserList() ([]map[string]interface{}, error) {
	var users []User
	err := DB.Where("1 = ?", 1).Find(&users).Error

	var userList []map[string]interface{}
	for _, user := range users {
		formatUser := formatUser(user)
		userList = append(userList, formatUser)
	}

	return userList, err
} */

/* func GetUserInfoByUserName(username string) (map[string]interface{}, error) {
	var user User

	err := DB.Where("user_name = ?", username).First(&user).Error
	formatUser := formatUser(user)

	return formatUser, err
} */
