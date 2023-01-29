package models

import "gobooklibrary/config"

type ProfileInfo struct {
	Login string `json:"login" binding:"required"`
}

func (p *ProfileInfo) TableName() string {
	return "Users"
}

func GetProfileInfo(profile ProfileInfo) (err error) {
	//profileinfo, err := config.DB.Query("EXECUTE showUserInfoBasedOnLogin USING ?", profile.Login)

	if err = config.ORMDB.Where("login = ?", profile.Login).Find(&profile).Error; err != nil {
		return err
	}

	return nil
}
