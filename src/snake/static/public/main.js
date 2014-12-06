// main.js
(function() {
    var userid = document.getElementById('userid').innerHTML,
        bodyW = 7,
        gap = 1,
        mapW,
        mapH,
        start = false;



    var connectInfoEl = document.getElementById('connectInfo'),
        canvasEl = document.getElementById('mainCanvas'),
        ctx = canvasEl.getContext("2d");


    // 连接
    var connect = new WebSocket("ws://" + location.host + "/ws");
    connect.onopen = function() {
        connectInfoEl.innerHTML = "已连接websocket";
        this.send("getMap"); // 获取地图
    }
    connect.onclose = function() {
        connectInfoEl.innerHTML = "连接已断开";
    }

    // 获取数据
    connect.onmessage = function(e) {
        connectInfoEl.innerHTML = e.data;

        var data = JSON.parse(e.data)
        
        switch (data.Type){
            case "s&f":
                // 绘蛇和食物
                drawSnakeFood(data);
                break;
            
            case "map":
                setMap(data);
                start = true;
                break;
        }
        
    }

    // 按键
    document.onkeydown = function(event) {
        switch (event.keyCode) {
            case 37:
                connect.send(4);
                break;
            case 38:
                connect.send(1);
                break;
            case 39:
                connect.send(2);
                break;
            case 40:
                connect.send(3);
                break;
        }
    };

    // 触摸动作
    (function () {
        var startX,
            srartY,
            endX,
            endY;
        document.addEventListener('touchstart', touchStart, false);
        document.addEventListener('touchmove', touchMove, false);
        document.addEventListener('touchend', touchEnd, false);
        function touchStart (event) {
            event.preventDefault();
            startX = event.touches[0].pageX;
            startY = event.touches[0].pageY;
        }
        function touchMove (event) {
            endX = event.touches[0].pageX;
            endY = event.touches[0].pageY;

        }
        function touchEnd (event) {
            var x = Math.abs(startX-endX),
                y = Math.abs(startY-endY);

            if (x>y) {
                if (startX>endX) {connect.send(4)} else{connect.send(2)};
            }else {
                if (startY>endY) {connect.send(1)} else{connect.send(3)};
            }
        }
    })();

    // 绘图方法
    function drawSnakeFood(data) {
        if (!start) {
            return;
        };
        // 清除画布
        ctx.clearRect(0, 0, canvasEl.width, canvasEl.height);

        // 绘制蛇
        var Snakes = data.Snakes;
        for (var i = Snakes.length - 1; i >= 0; i--) {
            var snake = Snakes[i],
                body = snake.Body,
                name = snake.Name,
                staust = snake.Staust;


            if (name==userid && staust ==1) {
                ctx.fillStyle = "#000";
            } else if (staust == 0) {
                ctx.fillStyle = "#800"
            } else{
                ctx.fillStyle = "#666";
            };

            for (var j = body.length - 1; j >= 0; j--) {
                var unit = body[j];
                drawRect(unit);
            };

        };

        // 绘制食物
        var Foods = data.Foods.Foods;
        for (var i = Foods.length - 1; i >= 0; i--) {
            var food = Foods[i];
            ctx.strokeStyle = 'orange';
            drawStrokeRect(food);
        };
    }

    function setMap(data) {
        mapW = data.Map.Width;
        mapH = data.Map.Height;
        canvasEl.width  =(gap + bodyW)*mapW + gap;
        canvasEl.style.width = canvasEl.width + "px";
        canvasEl.height  =(gap + bodyW)*mapH + gap;
        canvasEl.style.height = canvasEl.height + "px";
    }


    // 绘制方格方法
    function drawRect(unit) {
        ctx.fillRect((unit[0]-1) * (bodyW+gap) +gap , (unit[1]-1) * (bodyW+gap) +gap , bodyW, bodyW);
    }

    function drawStrokeRect(unit) {
        ctx.strokeRect((unit[0]-1) * (bodyW+gap) +gap , (unit[1]-1) * (bodyW+gap) +gap , bodyW, bodyW);
    }


})();