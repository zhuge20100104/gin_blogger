package logic

import (
	"fmt"

	"github.com/zhuge20100104/gin_blogger/model"

	"github.com/zhuge20100104/gin_blogger/dal/db"
)

func GetAllCategoryList() (categoryList []*model.Category, err error) {
	//1. 从数据库中，获取文章分类列表
	categoryList, err = db.GetAllCategoryList()
	if err != nil {
		fmt.Printf("1 get article list failed, err:%v\n", err)
		return
	}

	return
}
