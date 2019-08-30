package server

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/chat/model"
	"go-chat/app/chat/util/helper"
	"net/http"
	"time"
)

const (
	_tokenExpire  = 30 * 24 * 3600
)


func Login(ctx *gin.Context)  {


	var (
		err   error
		user  model.User
		token helper.Token
	)

	params := new(struct {
		Email   string `form:"email" json:"email" binding:"required"`
		Password   string `form:"password" json:"password" binding:"required"`
	})

	if err = ctx.ShouldBindJSON(&params); err != nil {

		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	db := DaoInstance.ChatDb.Model(model.User{}).Where("email = ?", params.Email)
	if err = db.Find(&user).Error; err != nil {
		// user not found
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	if user.Password != helper.SetPassword(params.Password, user.Salt) {
		// wrong password
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "error password",
			"data": make(map[string]interface{}),
		})
		return
	}

	// create a token
	token, err = helper.GenerateToken(user.Uid, _tokenExpire)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": token,
	})

}


func Register(ctx *gin.Context)  {

	var (
		err   error
		user  model.User
		token helper.Token
	)


	// request params
	params := new(struct {
		Email   string `form:"email" json:"email" binding:"required"`
		Password   string `form:"password" json:"password" binding:"required"`
	})

	if err = ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	DaoInstance.ChatDb.Model(model.User{}).Where("email = ?", params.Email).Find(&user)

	if user.Uid > 0 {
		// user exist
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "the email has been used",
			"data": make(map[string]interface{}),
		})
		return
	}

	salt := helper.GetRandomString(6)
	user = model.User{
		Nickname:   "user_" + helper.GetRandomString(4),
		Avatar:     "https://wx.qlogo.cn/mmopen/vi_32/KMzUH2aj5qvOostTYcJC1AgQSQzPTKRT80U0WxyFN3TmLqHGveFicTXs0W8jq94avzwvsvI84jLoFowMNYtL7zg/132",
		Email:      params.Email,
		Salt:       salt,
		Password:   helper.SetPassword(params.Password, salt),
		CreateIp:   ctx.ClientIP(),
		CreateTime: int(time.Now().Unix()),
	}

	if err = DaoInstance.ChatDb.Model(model.User{}).Create(&user).Error; err != nil {
		// create user error
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	// create a token
	token, err = helper.GenerateToken(user.Uid, _tokenExpire)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": token,
	})
}