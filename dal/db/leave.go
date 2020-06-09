package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zhuge20100104/gin_blogger/model"
)

func InsertLeave(leave *model.Leave) (err error) {
	// 插入留言
	// leave是关键字需要反引号
	sqlstr := "insert into `leave`(username,email,content)values(?,?,?)"
	_, err = DB.Exec(sqlstr, leave.Username, leave.Email, leave.Content)
	if err != nil {
		fmt.Printf("exec sql:%s failed, err:%v\n", sqlstr, err)
		return
	}

	return
}

func GetLeaveList() (leaveList []*model.Leave, err error) {
	// 获取留言列表
	sqlstr := "select id, username, email, content, create_time from `leave` order by id desc"
	err = DB.Select(&leaveList, sqlstr)
	if err != nil {
		fmt.Printf("exec sql:%s failed, err:%v\n", sqlstr, err)
		return
	}

	return
}
