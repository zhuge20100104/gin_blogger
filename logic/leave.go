package logic

import (
	"fmt"

	"github.com/zhuge20100104/gin_blogger/model"

	"github.com/zhuge20100104/gin_blogger/dal/db"
)

func InsertLeave(username, email, content string) (err error) {
	// 插入留言
	var leave model.Leave
	leave.Content = content
	leave.Email = email
	leave.Username = username

	err = db.InsertLeave(&leave)
	if err != nil {
		fmt.Printf("insert leave failed, err:%v, leave:%#v\n", err, leave)
		return
	}

	return
}

func GetLeaveList() (leaveList []*model.Leave, err error) {
	// 获取留言列表
	leaveList, err = db.GetLeaveList()
	if err != nil {
		fmt.Printf("get leave list failed, err:%v\n", err)
		return
	}
	return
}
