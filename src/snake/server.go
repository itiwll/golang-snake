package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// 地图类型
type maper struct {
	Width  int
	Height int
}

var (
	userIds []int
	Snakes  []*snake                                                      //蛇库
	conns   map[string]*websocket.Conn = make(map[string]*websocket.Conn) //连接库
	foods   fooder                                                        //食物库
	gameMap maper                                                         // 地图
	config  struct {                   // 设置
		port       string
		tickerTime int
	}
)

func main() {

	gameMap.Width = 100
	gameMap.Height = 50

	config.port = "80"
	config.tickerTime = 500

	// 静态文件
	http.Handle("/public/", http.FileServer(http.Dir("static")))

	// 动态页面
	http.HandleFunc("/", homeServer)

	// 数据接口
	http.HandleFunc("/ws", websocketServer)

	// 游戏服务启动
	go game()

	// 启动
	log.Fatal(http.ListenAndServe(":"+config.port, nil))

}

// html服务
func homeServer(rw http.ResponseWriter, req *http.Request) {

	cookitUserId, _ := req.Cookie("userId")

	// 是否新玩家
	newUser := true
	if cookitUserId == nil {
		newUser = true
	} else {
		// 遍历用户表
		for _, v := range userIds {
			if strconv.Itoa(v) == cookitUserId.Value {
				newUser = false
				break
			} else {
				newUser = true
			}
		}
	}

	if newUser {
		// 新用户
		id := len(userIds)
		fmt.Println("新用户id：" + strconv.Itoa(id))
		cookitUserId = &http.Cookie{
			Name:  "userId",
			Value: strconv.Itoa(id),
		}
		http.SetCookie(rw, cookitUserId) //设置cookit

		// id库更新
		userIds = append(userIds, id)

		// 分配新snake
		userSnake := snake{
			cookitUserId.Value,
			[][2]int{{1, 1}, {1, 2}, {1, 3}},
			3,
			3,
			1,
		}
		// 蛇库更新
		Snakes = append(Snakes, &userSnake)
	}

	fmt.Println("用户访问，ID：" + cookitUserId.Value)
	tp, _ := template.ParseFiles("home.html")
	tp.Execute(rw, cookitUserId.Value)
}

// websockit
func websocketServer(rw http.ResponseWriter, req *http.Request) {
	cookit, _ := req.Cookie("userId")
	uid, _ := strconv.Atoi(cookit.Value)
	ug := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
	conn, _ := ug.Upgrade(rw, req, nil)
	defer fmt.Println("用户" + cookit.Value + "断开连接")
	defer conn.Close()
	defer Snakes[uid].die()
	if conn != nil {
		fmt.Println("建立一个websockit连接")
		fmt.Println("存入连接库")
		conns[cookit.Value] = conn //存连接
	}

	// 接收操作
	for {
		t, p, err := conn.ReadMessage()
		if err != nil {
			delete(conns, cookit.Value)
			return
		}
		fmt.Println("读取")
		fmt.Println(string(p))
		if t == websocket.TextMessage {
			switch string(p) {
			case "getMap": // 地图请求
				writerMap(conn)

			case "1", "2", "3", "4": // 移动请求
				fmt.Println("收到操作 id:" + cookit.Value + ";value:" + string(p))
				i, _ := strconv.Atoi(cookit.Value)
				s := Snakes[i]
				s.direcion, _ = strconv.Atoi(string(p))
			}
		}
	}
}

// 游戏服务
func game() {
	fmt.Println("游戏服务启动")
	foods.produceFood(gameMap, Snakes) // 生成食物
	foods.produceFood(gameMap, Snakes)
	runTicker := time.NewTicker(time.Duration(config.tickerTime) * time.Millisecond)
	for {
		select {
		case <-runTicker.C:
			snakeDo()            // 动
			collide(Snakes)      // 计算碰撞
			go writerSnakeFood() // 向客户端传送数据
		}
	}
}

// 动作
func snakeDo() {
	// 移动、吃掉食物、创造食物
	for _, v := range Snakes {
		if v.Staust == 0 {
			continue
		}
		if v.Move(&foods) {
			fmt.Println("生成食物")
			foods.produceFood(gameMap, Snakes)
		}
	}
}

// 计算碰撞
func collide(Snakes []*snake) {
	// 遍历蛇库
	for i, s1 := range Snakes {
		// 蛇的状态
		if s1.Staust != 1 {
			continue
		}
		head := s1.Body[len(s1.Body)-1] // 头

		if head[0] <= 0 || head[1] <= 0 || head[0] > gameMap.Width || head[1] > gameMap.Height {
			fmt.Println("撞墙死掉了")
			s1.die()
		}
	H:
		for j, s2 := range Snakes {
			for k, unit := range s2.Body {
				if head[0] == unit[0] && head[1] == unit[1] {
					if i != j || k != len(s1.Body)-1 { // 排除本身
						fmt.Println("死掉了")
						s1.die()
						break H
					}
				}
			}
		}
	}
}

// 写到客户端 蛇库和食物库
func writerSnakeFood() {
	var json struct {
		Type   string
		Snakes []*snake
		Foods  fooder
	}
	json.Type = "s&f"
	json.Snakes = Snakes
	json.Foods = foods
	// fmt.Println("传送数据")
	for _, v := range conns {
		v.WriteJSON(&json)
	}
}

// 写到客户端 地图

func writerMap(c *websocket.Conn) {
	var json struct {
		Type string
		Map  maper
	}
	json.Type = "map"
	json.Map = gameMap
	c.WriteJSON(json)
}
