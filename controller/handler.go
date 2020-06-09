package controller

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/zhuge20100104/gin_blogger/util"

	"github.com/zhuge20100104/gin_blogger/logic"
)

var (
	uploadConfig map[string]interface{}
)

func IndexHandle(c *gin.Context) {
	// 首页列表
	// 取出结果
	articleRecordList, err := logic.GetArticleRecordList(0, 15)
	if err != nil {
		fmt.Printf("get article failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	allCategoryList, err := logic.GetAllCategoryList()
	if err != nil {
		fmt.Printf("get category list failed, err:%v\n", err)
	}
	// 定义模板变量
	var data map[string]interface{} = make(map[string]interface{}, 10)
	data["article_list"] = articleRecordList
	data["category_list"] = allCategoryList
	// 模板变量返回给模板进行渲染
	c.HTML(http.StatusOK, "views/index.html", data)
}

func CategoryList(c *gin.Context) {
	// 首页文章分类信息列表
	categoryIdStr := c.Query("category_id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	articleRecordList, err := logic.GetArticleRecordListById(int(categoryId), 0, 15)
	if err != nil {
		fmt.Printf("get article failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	allCategoryList, err := logic.GetAllCategoryList()
	if err != nil {
		fmt.Printf("get category list failed, err:%v\n", err)
	}

	var data map[string]interface{} = make(map[string]interface{}, 10)
	data["article_list"] = articleRecordList
	data["category_list"] = allCategoryList

	c.HTML(http.StatusOK, "views/index.html", data)
}

func NewArticle(c *gin.Context) {
	// 新增文章
	categoryList, err := logic.GetAllCategoryList()
	// 获取所有文章分类
	if err != nil {
		fmt.Printf("get article failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.HTML(http.StatusOK, "views/post_article.html", categoryList)
}

func LeaveNew(c *gin.Context) {
	// 获取留言列表
	leaveList, err := logic.GetLeaveList()
	if err != nil {
		fmt.Printf("get leave failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.HTML(http.StatusOK, "views/gbook.html", leaveList)
}

func AboutMe(c *gin.Context) {
	c.HTML(http.StatusOK, "views/about.html", gin.H{
		"title": "Posts",
	})
}

func ArticleSubmit(c *gin.Context) {
	// 文章提交
	content := c.PostForm("content")
	author := c.PostForm("author")
	categoryIdStr := c.PostForm("category_id")
	title := c.PostForm("title")
	// 转整数，10进制64位
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	err = logic.InsertArticle(content, author, title, categoryId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}
	// 重定向 301
	c.Redirect(http.StatusMovedPermanently, "/")
}

func ArticleDetail(c *gin.Context) {
	// 获取文章详情
	// 获取查询参数，文章id
	articleIdStr := c.Query("article_id")
	articleId, err := strconv.ParseInt(articleIdStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}
	// 获取文章详情
	articleDetail, err := logic.GetArticleDetail(articleId)
	if err != nil {
		fmt.Printf("get article detail failed,article_id:%d err:%v\n", articleId, err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	fmt.Printf("article detail:%#v\n", articleDetail)
	// 获取相关文章
	relativeArticle, err := logic.GetRelativeAricleList(articleId)
	if err != nil {
		fmt.Printf("get relative article failed, err:%v\n", err)
	}
	// 获取上一篇和下一篇文章
	prevArticle, nextArticle, err := logic.GetPrevAndNextArticleInfo(articleId)
	if err != nil {
		fmt.Printf("get prev or next article failed, err:%v\n", err)
	}
	// 获取所有分类列表
	allCategoryList, err := logic.GetAllCategoryList()
	if err != nil {
		fmt.Printf("get all category failed, err:%v\n", err)
	}
	// 获取评论列表
	commentList, err := logic.GetCommentList(articleId)
	if err != nil {
		fmt.Printf("get comment list failed, err:%v\n", err)
	}

	fmt.Printf("relative article size:%d article_id:%d\n", len(relativeArticle), articleId)
	// 使用map传递多个对象实例给模板
	var m map[string]interface{} = make(map[string]interface{}, 10)
	m["detail"] = articleDetail
	m["relative_article"] = relativeArticle
	m["prev"] = prevArticle
	m["next"] = nextArticle
	m["category"] = allCategoryList
	// 评论使用的隐藏字段，用于确定评论属于哪一个文章
	m["article_id"] = articleId
	m["comment_list"] = commentList

	c.HTML(http.StatusOK, "views/detail.html", m)
}

func UploadFile(c *gin.Context) {
	// single file
	// 接受发表文章中的图片
	file, err := c.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println(file.Filename)
	rootPath := util.GetRootDir()
	u2 := uuid.NewV4()
	// if err != nil {
	// 	return
	// }

	ext := path.Ext(file.Filename)
	url := fmt.Sprintf("/static/upload/%s%s", u2, ext)
	dst := filepath.Join(rootPath, url)
	// Upload the file to specific dst.
	c.SaveUploadedFile(file, dst)
	c.JSON(http.StatusOK, gin.H{
		"uploaded": true,
		"url":      url,
	})
}

func CommentSubmit(c *gin.Context) {
	// 评论提交
	comment := c.PostForm("comment")
	author := c.PostForm("author")
	email := c.PostForm("email")
	articleIdStr := c.PostForm("article_id")

	articleId, err := strconv.ParseInt(articleIdStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	err = logic.InsertComment(comment, author, email, articleId)
	if err != nil {
		fmt.Printf("insert comment failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	url := fmt.Sprintf("/article/detail/?article_id=%d", articleId)
	c.Redirect(http.StatusMovedPermanently, url)
}

func LeaveSubmit(c *gin.Context) {
	// 发表留言
	comment := c.PostForm("comment")
	author := c.PostForm("author")
	email := c.PostForm("email")

	err := logic.InsertLeave(author, email, comment)
	if err != nil {
		fmt.Printf("insert leave failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	url := fmt.Sprintf("/leave/new/")
	c.Redirect(http.StatusMovedPermanently, url)
}
