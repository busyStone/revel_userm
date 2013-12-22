package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"labix.org/v2/mgo/bson"
)

func (u *LoginUser) IsRegistered() (bool, error) {
	dal, err := NewDal()
	if err != nil {
		return false, err
	}

	defer dal.Close()

	uc := dal.session.DB(DbName).C(UserCollection)

	result := User{}
	err = uc.Find(bson.M{"email": u.Email}).One(&result)
	if err != nil {
		return false, err
	}

	if err = bcrypt.CompareHashAndPassword(result.Password, []byte(u.Password)); err == nil {
		return true, nil
	}

	return false, err
}

func (mu *MockUser) Register() error {
	dal, err := NewDal()
	if err != nil {
		return err
	}

	defer dal.Close()

	uc := dal.session.DB(DbName).C(UserCollection)

	// 检查 email、nickname 是否已经被使用
	i, _ := uc.Find(bson.M{"email": mu.Email}).Count()
	if i != 0 {
		return errors.New("邮箱已被使用！")
	}

	i, _ = uc.Find(bson.M{"nickname": mu.Nickname}).Count()
	if i != 0 {
		return errors.New("用户名已被使用！")
	}

	var u User

	u.Email = mu.Email
	u.Nickname = mu.Nickname
	u.Password, _ = bcrypt.GenerateFromPassword([]byte(mu.Password), bcrypt.DefaultCost)

	err = uc.Insert(u)

	return err
}