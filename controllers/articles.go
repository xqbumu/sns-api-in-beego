package controllers

import (
	"os"
	"pet/models"
	"pet/utils"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

// 文章相关操作
type ArticlesController struct {
	beego.Controller
}

func (this *ArticlesController) URLMapping() {
	this.Mapping("Post", this.Post)
	this.Mapping("GetOne", this.GetOne)
	this.Mapping("GetAll", this.GetAll)
	this.Mapping("Put", this.Put)
	this.Mapping("Delete", this.Delete)
}

// @Title 创建文章
// @Description 创建文章
// @Param	title		form	string	true	"文章标题"
// @Param	content		form	string	true	"文章内容"
// @Param	title_image	form	file	true	"文章标题配图"
// @Success 200 {int} models.Articles.Id
// @Failure 403 body is empty
// @router / [post]
func (this *ArticlesController) Post() {
	var v models.Articles
	valid := validation.Validation{}
	this.ParseForm(&v)
	passed, _ := valid.Valid(&v)
	if !passed {
		outPut := helper.Reponse(1, nil, valid.Errors[0].Key+" "+valid.Errors[0].Message)
		this.Data["json"] = outPut
	} else {
		todayDateDir := "/" + helper.GetTodayDate()
		if _, err := os.Stat(uploadPhotoPath + todayDateDir); os.IsNotExist(err) {
			os.Mkdir(uploadPhotoPath+todayDateDir, 0777)
		}
		currentUser := this.GetSession("user").(models.Users)
		photoName := helper.GetGuid(currentUser.Id)
		dateSubdir := "/" + string(photoName[0]) + string(photoName[1])

		if _, err := os.Stat(uploadPhotoPath + todayDateDir + dateSubdir); os.IsNotExist(err) {
			os.Mkdir(uploadPhotoPath+todayDateDir+dateSubdir, 0777)
		}

		imagePath := uploadPhotoPath + todayDateDir + dateSubdir + "/" + photoName + ".jpg"

		err := this.SaveToFile("title_image", imagePath)

		if err != nil {
			outPut := helper.Reponse(1, nil, err.Error())
			this.Data["json"] = outPut
		} else {
			v.TitleImage = imagePath
			v.CreatedAt = time.Now()
			v.UpdatedAt = time.Now()
			if id, err := models.AddArticles(&v); err == nil {
				v.Id = int(id)
				outPut := helper.Reponse(0, v, "创建成功")
				this.Data["json"] = outPut
			} else {
				outPut := helper.Reponse(1, nil, err.Error())
				this.Data["json"] = outPut
			}
		}
	}
	this.ServeJson()
}

// @Title 获取文章
// @Description 通过文章id获取详情
// @Param	id		path 	string	true		"文章Id"
// @Success 200 {object} models.Articles
// @Failure 403 :id is empty
// @router /:id [get]
func (this *ArticlesController) GetOne() {
	idStr := this.Ctx.Input.Params[":id"]
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetArticlesById(id)
	if err != nil {
		this.Data["json"] = err.Error()
	} else {
		this.Data["json"] = v
	}
	this.ServeJson()
}

// @Title 获取文章列表
// @Description 获取文章列表
// @Param	offset	query	string	false	"结果索引"
// @Success 200 {object} models.Articles
// @Failure 403
// @router / [get]
func (this *ArticlesController) GetAll() {

	if v, err := this.GetInt("offset"); err == nil {
		offset = int64(v)
	}
	l, err := models.GetAllArticles(query, fields, sortby, order, offset, limit)
	if err != nil {
		outPut := helper.Reponse(1, nil, err.Error())
		this.Data["json"] = outPut
	} else {
		outPut := helper.Reponse(0, l, "")
		this.Data["json"] = outPut
	}
	this.ServeJson()
}

// @Title 更新文章
// @Description 更新文章内容
// @Param	id		path 	string	true		"文章ID"
// @Param	title		form 	string	true		"文章标题"
// @Param	content		form 	string	true		"文章内容"
// @Param	title_image		form 	file	true		"标题图片"
// @Success 200 {object} models.Articles
// @Failure 403 :id is not int
// @router /:id [put]
func (this *ArticlesController) Put() {
	idStr := this.Ctx.Input.Params[":id"]
	id, _ := strconv.Atoi(idStr)

	v := models.Articles{Id: id}

	artical, err := models.GetArticlesById(id)
	if err != nil {
		outPut := helper.Reponse(1, nil, err.Error())
		this.Data["json"] = outPut
		this.ServeJson()
		return
	}
	valid := validation.Validation{}
	this.ParseForm(&v)
	passed, _ := valid.Valid(&v)
	if !passed {
		outPut := helper.Reponse(1, nil, valid.Errors[0].Key+" "+valid.Errors[0].Message)
		this.Data["json"] = outPut
	} else {

		todayDateDir := "/" + helper.GetTodayDate()
		if _, err := os.Stat(uploadPhotoPath + todayDateDir); os.IsNotExist(err) {
			os.Mkdir(uploadPhotoPath+todayDateDir, 0777)
		}
		currentUser := this.GetSession("user").(models.Users)
		photoName := helper.GetGuid(currentUser.Id)
		dateSubdir := "/" + string(photoName[0]) + string(photoName[1])

		if _, err := os.Stat(uploadPhotoPath + todayDateDir + dateSubdir); os.IsNotExist(err) {
			os.Mkdir(uploadPhotoPath+todayDateDir+dateSubdir, 0777)
		}

		imagePath := uploadPhotoPath + todayDateDir + dateSubdir + "/" + photoName + ".jpg"
		err := this.SaveToFile("title_image", imagePath)

		if err != nil {
			outPut := helper.Reponse(1, nil, err.Error())
			this.Data["json"] = outPut

		} else {
			artical.UpdatedAt = time.Now()
			artical.Title = v.Title
			artical.Content = v.Content
			artical.TitleImage = imagePath

			if err := models.UpdateArticlesById(artical); err == nil {
				v.Id = int(id)
				outPut := helper.Reponse(0, v, "创建成功")
				this.Data["json"] = outPut
			} else {
				outPut := helper.Reponse(1, nil, err.Error())
				this.Data["json"] = outPut
			}
		}
	}

	this.ServeJson()
}

// @Title Delete
// @Description delete the Articles
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]

func (this *ArticlesController) Delete() {
	//idStr := this.Ctx.Input.Params[":id"]
	//id, _ := strconv.Atoi(idStr)
	//if err := models.DeleteArticles(id); err == nil {
	//this.Data["json"] = "OK"
	//} else {
	//this.Data["json"] = err.Error()
	//}
	//this.ServeJson()
}
