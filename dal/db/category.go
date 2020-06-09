package db

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zhuge20100104/gin_blogger/model"
)

func InsertCategory(category *model.Category) (categoryId int64, err error) {

	sqlstr := "insert into category(category_name, category_no)values(?,?)"
	result, err := DB.Exec(sqlstr, category.CategoryName, category.CategoryNo)
	if err != nil {
		return
	}

	categoryId, err = result.LastInsertId()
	return
}

func GetCategoryList(categoryIds []int64) (categoryList []*model.Category, err error) {
	// 获取对应分类列表
	sqlstr, args, err := sqlx.In("select id, category_name, category_no from category where id in(?)", categoryIds)
	if err != nil {
		return
	}
	// 传入时必须将切片展开
	err = DB.Select(&categoryList, sqlstr, args...)
	return
}

func GetAllCategoryList() (categoryList []*model.Category, err error) {
	// 获取所有分类列表
	sqlstr := "select id, category_name, category_no from category order by category_no asc"
	err = DB.Select(&categoryList, sqlstr)
	return
}

func GetCategoryById(id int64) (category *model.Category, err error) {
	// 根据id返回分类
	category = &model.Category{}
	sqlstr := "select id, category_name, category_no from category where id=?"
	err = DB.Get(category, sqlstr, id)
	return
}
