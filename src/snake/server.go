package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

// 地图类型
type maper struct {
	Width  int
	Height int
}

const (
	PORT       = "80"
	TICKERTIME = 150
	MAP_WIDTH  = 150
	MAP_HEIGHT = 80
)

var (
	users   []user            //用户库
	Snakes  []*snake          //蛇库
	conns   []*websocket.Conn //连接库
	foods   fooder            //食物库
	gameMap maper             // 地图
	config  struct {          // 设置
		port       string
		tickerTime int
	}
)

func main() {
	setConfig()

	ip := "localhost"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		if addr.String() != "0.0.0.0" {
			ip = addr.String()
		}
	}
	fmt.Printf("服务地址：%v:%v\n", ip, config.port)

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

// 设置运行参数
func setConfig() {
	config.port = PORT
	config.tickerTime = TICKERTIME
	gameMap.Width = MAP_WIDTH
	gameMap.Height = MAP_HEIGHT

	var in string
	fmt.Println("是否使用默认配置？")
	fmt.Println("端口：80")
	fmt.Println("地图宽：150")
	fmt.Println("地图高：80")
	fmt.Println("y/n(y)")
S:
	fmt.Scanln(&in)
	if in == "y" || in == "" {
		fmt.Println("使用默认配置。")
	} else if in != "y" && in != "n" {
		fmt.Println("输入错误")
		goto S
	} else {
		for {
			fmt.Println("请输入端口(" + PORT + ")：")
			i, _ := fmt.Scanln(&in)
			if i == 0 {
				break
			}
			_, err := strconv.Atoi(in)
			if err != nil {
				fmt.Println("输入错误")
			} else {
				config.port = in
				break
			}

		}
		for {
			fmt.Println("请输入心跳时间(" + strconv.Itoa(TICKERTIME) + ")：")
			i, _ := fmt.Scanln(&in)
			if i == 0 {
				break
			}
			v, err := strconv.Atoi(in)
			if err != nil {
				fmt.Println("输入错误")
			} else {
				config.tickerTime = v
				break
			}
		}

		for {
			fmt.Println("请输入地图宽度(" + strconv.Itoa(MAP_WIDTH) + ")：")
			i, _ := fmt.Scanln(&in)
			if i == 0 {
				break
			}
			v, err := strconv.Atoi(in)
			if err != nil {
				fmt.Println("输入错误")
			} else {
				gameMap.Width = v
				break
			}

		}
		for {
			fmt.Println("请输入地图高度(" + strconv.Itoa(MAP_HEIGHT) + ")：")
			i, _ := fmt.Scanln(&in)
			if i == 0 {
				break
			}
			v, err := strconv.Atoi(in)
			if err != nil {
				fmt.Println("输入错误")
			} else {
				gameMap.Height = v
				break
			}

		}
	}
}

// html服务
func homeServer(rw http.ResponseWriter, req *http.Request) {

	user := newUser()

	cookitUserId := &http.Cookie{
		Name:  "userId",
		Value: strconv.Itoa(user.id),
	}
	http.SetCookie(rw, cookitUserId) //设置cookit

	// user库更新
	users = append(users, user)

	tp, _ := template.ParseFiles("home.html")
	tp.Execute(rw, cookitUserId.Value)
}

// websockit
func websocketServer(rw http.ResponseWriter, req *http.Request) {
	cookit, _ := req.Cookie("userId")
	uid, _ := strconv.Atoi(cookit.Value)

	var u user
	for _, u = range users {
		if u.id == uid {
			break
		}
	}

	// fmt.Println(u)

	ug := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
	conn, _ := ug.Upgrade(rw, req, nil)
	defer fmt.Println("用户" + cookit.Value + "断开连接")
	defer conn.Close()
	if conn != nil {
		fmt.Println("建立一个websockit连接")
		fmt.Println("存入连接库")
		u.conn = conn //存连接
	}

	// 接收操作
	for {
		t, p, err := u.conn.ReadMessage()
		if err != nil {
			u.conn = nil
			return
		}
		fmt.Println("读取")
		fmt.Println(string(p))
		if t == websocket.TextMessage {
			switch string(p) {
			case "getMap": // 地图请求
				writerMap(conn)
			case "newSnake": // 新蛇请求 加入游戏
				// todo
				u.newSnake()
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
