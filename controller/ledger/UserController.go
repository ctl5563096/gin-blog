package ledger

import (
	ledgerUpload "gin-blog/controller/upload"
	"gin-blog/models/ledger"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	uploadOss "gin-blog/pkg/oss"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"path"
	"strings"
)

func GetProfile(c *gin.Context) {
	ownerId := getLedgerOwnerId(c)
	if ownerId == "" {
		app.MissToken(c)
		return
	}

	user, err := ledger.GetUser(ownerId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user, err = ledger.UpsertUser(ownerId, "")
			if err == nil {
				app.OkWithData(user, c)
				return
			}
		}
		app.FailWithMessage("用户信息获取失败", 1, c)
		return
	}
	app.OkWithData(user, c)
}

func UpdateProfile(c *gin.Context) {
	var r ledger.UpdateProfile
	if !bindAndValidate(c, &r) {
		return
	}
	r.OwnerId = getLedgerOwnerId(c)

	user, err := ledger.UpdateUserProfile(&r)
	if err != nil {
		app.FailWithMessage("用户资料保存失败", 1, c)
		return
	}
	app.OkWithData(user, c)
}

func UploadAvatar(c *gin.Context) {
	ownerId := getLedgerOwnerId(c)
	if ownerId == "" {
		app.MissToken(c)
		return
	}

	f, err := c.FormFile("avatar")
	if err != nil {
		app.FailWithMessage(e.GetMsg(e.MISS_PARAMS), e.MISS_PARAMS, c)
		return
	}

	ext := strings.ToLower(path.Ext(f.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".gif" && ext != ".jpeg" {
		app.FailWithMessage(e.GetMsg(e.IMAGE_TYPE_ERROR), e.IMAGE_TYPE_ERROR, c)
		return
	}

	localPath, randomName, ok := ledgerUpload.GetNewFileName(ext, "resource/public/pic/ledger-avatar")
	if !ok {
		app.Fail(c)
		return
	}

	if err = c.SaveUploadedFile(f, localPath); err != nil {
		app.FailWithMessage(e.GetMsg(e.UPLOAD_FAIL), e.UPLOAD_FAIL, c)
		return
	}

	ossUrl := uploadOss.UploadALiYunOss(c, "ledger/avatar", randomName, ext, localPath)
	if ossUrl == nil {
		return
	}

	app.OkWithData(gin.H{
		"path": ossUrl,
	}, c)
}
