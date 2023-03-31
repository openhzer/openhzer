package ws

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type UserWs struct {
	Coon *websocket.Conn
	Mt   int
	Id   int
}

var (
	users  = make(map[int]*UserWs)
	id     int
	writes = make(chan []byte)
)

//设置websocket
//CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func testws(c *gin.Context) {
	//升级get请求为webSocket协议
	//创建一个新的userWs
	var userWs UserWs
	var err error
	userWs.Coon, err = upGrader.Upgrade(c.Writer, c.Request, nil)

	//如果有错结束
	if err != nil {
		return
	}
	//如果没有错将 *websocket.Conn 赋予userWs.Coon

	id++
	userWs.Id = id
	users[id] = &userWs
	fmt.Println("建立连接:", id)
	//设置连接断开事件
	userWs.Coon.SetCloseHandler(func(code int, text string) error {
		//关闭链接
		fmt.Println("关闭连接")
		_ = userWs.Coon.Close()
		fmt.Println("删除map", userWs.Id)
		//删除map
		delete(users, userWs.Id)

		return nil
	})

	_ = userWs.Coon.WriteMessage(1, []byte("ok"))

	////读取
	go read(&userWs)
}

func Write() {
	for {
		msg := <-writes
		//data,_ :=json.Marshal(msg)
		for _, userWs := range users {
			err := userWs.Coon.WriteMessage(1, msg)
			if err != nil {
				break
			}
		}

	}
}

//读取队列
func read(WsCoon *UserWs) {
	//读取ws中的数据

	for {
		_, message, err := WsCoon.Coon.ReadMessage()
		fmt.Printf("%#v\n", message)
		if err != nil {
			break
		} else if message != nil {
			fmt.Printf("%#v\n", "进入群发")
			fmt.Println(message)
			writes <- message
		}
	}

}
