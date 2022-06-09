package model

import (
	"errors"
	"go-search-engine/src/database"
	types "go-search-engine/src/global"
	"strconv"
)

type Member struct {
	UserID    int
	Nickname  string
	Username  string
	Password  string
	hobby     string
	UserType  types.UserType
	IsDeleted bool
}

func (Member) TableName() string {
	return "member"
}

/*
func (member *Member) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4().String()
	return scope.SetColumn("user_id", uuid)
}*/

func (member *Member) CreateMember() (string, error) {

	if err := db.Create(&member).Error; err != nil {
		return "", err
	}
	return strconv.Itoa(member.UserID), nil
}

func (model *Member) GetMember(user_id string) (Member, error) {
	var result Member
	err := database.MySqlDb.First(&Member{}, "user_id = ?", user_id).Scan(&result).Error
	return result, err
}

// GetAllMembers 返回所有成员
func (member *Member) GetAllMembers(offset, limit int) ([]Member, error) {
	var ans []Member
	err := database.MySqlDb.Limit(limit).Offset(offset).Find(&ans).Error
	if err != nil {
		return ans, err
	}
	return ans, nil
}

func GetMemberByUsername(username string) (Member, error) {
	var ans = Member{}
	err := db.Where("username = ? ", username).First(&ans).Error
	return ans, err
}

func UpdateMember(user_id string, nickname string) error {
	id, _ := strconv.Atoi(user_id)
	var result = Member{}
	db.Where(&Member{UserID: id}).First(&result)
	if result.Nickname == "" {
		return errors.New("未找到该用户")
	}
	if result.IsDeleted == true {
		return errors.New("用户已删除")
	}
	db.Model(&Member{}).Where("user_id = ?", user_id).Update("nickname", nickname)
	return nil
}

func DeleteMember(user_id string) error {
	var result = Member{}
	db.Where("user_id = ? ", user_id).First(&result)
	if result.Nickname == "" {
		return errors.New("未找到该用户")
	}
	if result.IsDeleted == true {
		return errors.New("用户已删除")
	}

	id, _ := strconv.Atoi(user_id)
	db.Model(&Member{}).Where("user_id = ?", id).Update("is_deleted", true)
	return nil
}
