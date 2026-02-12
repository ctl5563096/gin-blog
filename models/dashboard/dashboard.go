package dashboard

import (
	"gin-blog/models/blog"
	"gin-blog/models/resource"
	"gin-blog/models/user"
)

// DashboardStats 面板统计数据结构体
type DashboardStats struct {
	ArticleCount int `json:"article_count"`
	PhotoCount   int `json:"photo_count"`
	MusicCount   int `json:"music_count"`
	UserCount    int `json:"user_count"`
}

// GetDashboardStats 获取面板统计数据
func GetDashboardStats() DashboardStats {
	stats := DashboardStats{}

	// 获取文章数量 (is_delete = 0)
	dbBlog := blog.GetDB()
	var articleCount int
	dbBlog.Table("t_go_article").Where("is_delete = ?", 0).Count(&articleCount)
	stats.ArticleCount = articleCount

	// 获取图片数量 (is_delete = 1)
	dbResource := resource.GetDB()
	var photoCount int
	dbResource.Table("t_go_photos").Where("is_delete = ?", 1).Count(&photoCount)
	stats.PhotoCount = photoCount

	// 获取音乐数量 (is_delete = 1)
	var musicCount int
	dbResource.Table("t_go_music").Where("is_delete = ?", 1).Count(&musicCount)
	stats.MusicCount = musicCount

	// 获取用户数量
	dbUser := user.GetDB()
	var userCount int
	dbUser.Table("ci_admin_user").Count(&userCount)
	stats.UserCount = userCount

	return stats
}
