package tag_service

import (
	"io"
	"time"
	"encoding/json"
	"go-gin-ex/pkg/models"
	"go-gin-ex/pkg/logging"
	"go-gin-ex/pkg/export"

)

type Tag struct {
	ID int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistByTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] t.ModifiedBy
	data["name"] = t.Name

	if t.State >= 0 {
		data["state"] = t.State
	}

	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	reutrn models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]modles.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)
	cache := cache_service.Tag{
		State: t.State,
		PageNum: t.PageNum,
		PageSize: t.PageSize,
	}

	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, tags, 3600)
	return tags, nil
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}


}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	} 
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}