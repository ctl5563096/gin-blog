package resource

import (
	"encoding/json"
	"fmt"
)

const MusicTableName = "t_go_music"

// CreateMusicData 创建收藏音乐sql
type CreateMusicData struct {
	Id       int    `json:"id"`
	Title    string `json:"title" validate:"required"`
	Summary  string `json:"summary"`
	Lyric    string `json:"lyric"`
	Url      string `json:"url" validate:"required"`
	Author   string `json:"author" validate:"required"`
	Cover    string `json:"cover" validate:"required"`
	Thumb    string `json:"thumb"`
	FileName string `json:"file_name" validate:"required"`
	IsTop    int    `json:"is_top"`
}

// CreateMusicRecord 新建音乐
func CreateMusicRecord(params *CreateMusicData) (int, error) {
	err := db.Table(MusicTableName).Create(params).Error
	if err != nil {
		fmt.Println()
		return 0, err
	}
	m := make(map[string]interface{})
	j, _ := json.Marshal(params)
	err = json.Unmarshal(j, &m)
	if err != nil {
		return 0, err
	}
	recordId := int(m["id"].(float64))
	return recordId, err
}

// GetMusicList 获取音乐列表
func GetMusicList(page int, pageSize int, keyword string, order string) ([]CreateMusicData, error) {
	var r []CreateMusicData
	dbq := db.Table(MusicTableName).Select("id,title,summary,lyric,url,author,cover,thumb")
	if keyword != "" {
		keywords := "%" + keyword + "%"
		dbq = dbq.Where(" is_delete = ? and title LIKE ? ", 1, keywords)
	} else {
		dbq = dbq.Where(" is_delete = 1")
	}
	dbq.Order("updated_at " + order).Limit(pageSize).Offset((page - 1) * pageSize).Find(&r)
	if dbq.Error != nil {
		fmt.Println(dbq.Error.Error())
		return nil, dbq.Error
	}
	return r, nil
}

// GetMusicCount 获取总数
func GetMusicCount(keywords string) (int, error) {
	var count int
	dbq := db.Table(MusicTableName)
	if keywords != "" {
		keywords := "%" + keywords + "%"
		dbq = dbq.Where(" is_delete = ? and title LIKE ? ", 1, keywords)
	} else {
		dbq = dbq.Where(" is_delete = 1")
	}
	dbq.Count(&count)
	return count, nil
}

// DeleteMusicRecord 删除
func DeleteMusicRecord(id int) bool {
	values := map[string]interface{}{
		"is_delete": 2,
	}
	err := db.Table(MusicTableName).Where("id = ?", id).Updates(values)
	if err.Error != nil {
		return false
	}
	return true
}

// GetAudioDetail 获取音频详情
func GetAudioDetail(id int) (CreateMusicData, error) {
	var r CreateMusicData
	err := db.Table(MusicTableName).Where("id = ?", id).Where("is_delete = ?", 1).First(&r)
	if err.Error != nil {
		return r, err.Error
	}
	return r, nil
}

// UpdateAudioDetail 修改音频
func UpdateAudioDetail(id int, params CreateMusicData) bool {
	err := db.Table(MusicTableName).Where("id = ?", id).Update(params)
	if err.Error != nil {
		return false
	}
	return true
}
