package service

import (
	"errors"
	"fmt"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userService struct {
}

var UserService = new(userService)

func (u *userService) CheckUserExistByName(name string) (bool, error) {
	db := dao.GetMysql()
	var dbUser model.User
	if err := db.First(&dbUser, "username = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			return false, nil
		} else {
			// 其他错误
			return false, err
		}
	}
	return true, nil
}

func (u *userService) GetUserByName(name string) (*model.User, error) {
	db := dao.GetMysql()
	var dbUser model.User
	if err := db.First(&dbUser, "username = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			return nil, fmt.Errorf("username(%s) not found", name)
		} else {
			// 其他错误
			return nil, err
		}
	}
	return &dbUser, nil
}

func (u *userService) GetUserByID(id uint) (*model.User, error) {
	db := dao.GetMysql()
	var dbUser model.User
	if err := db.First(&dbUser, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			return nil, fmt.Errorf("id(%v) not found", id)
		} else {
			// 其他错误
			return nil, err
		}
	}
	return &dbUser, nil
}

func (u *userService) Register(req *request.RegisterRequest) (*model.User, error) {
	ok, err := u.CheckUserExistByName(req.Username)
	if err != nil {
		return nil, err
	}
	if !ok {
		user := model.NewUser(uuid.New().String(), req.Username, req.Password)
		dao.GetMysql().Create(&user)
		return user, nil
	} else {
		return nil, fmt.Errorf("username(%s) already exists", req.Username)
	}
}

func (u *userService) Login(req *request.LoginRequest) (*model.User, error) {
	dbUser, err := u.GetUserByName(req.Username)
	if err != nil {
		return nil, err
	}
	if dbUser.Password != req.Password {
		return nil, fmt.Errorf("password error")
	}
	return dbUser, nil
}

func (u *userService) ChangePassword(req *request.ChangePasswordRequest) error {
	dbUser, err := u.GetUserByID(req.ID)
	if err != nil {
		return err
	}
	if dbUser.Password != req.Password {
		return fmt.Errorf("password error")
	}
	if dbUser.Password == req.NewPassword {
		return fmt.Errorf("same as old password")
	}
	dbUser.Password = req.NewPassword
	return dao.GetMysql().Save(&dbUser).Error
}
