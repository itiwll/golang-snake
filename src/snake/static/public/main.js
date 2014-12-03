// main.js
(function() {
    var bodyW = 5;


    var connectInfoEl = document.getElementById('connectInfo'),
        canvasEl = document.getElementById('mainCanvas'),
        ctx = canvasEl.getContext("2d"),
        w = canvasEl.offsetWidth,
        h = canvasEl.offsetHeight;

    // 选择画布大小
    if (w > 1000) {
        canvasEl.width = w = w / 2;
        canvasEl.height = h = h / 2;
    } else {
        canvasEl.width = w;
        canvasEl.height = h;
    }

    // 连接
    var connect = new WebSocket("ws://" + location.host + "/ws");
    connect.onopen = function() {
        connectInfoEl.innerHTML = "已连接websocket";
    }
    connect.onclose = function() {
        connectInfoEl.innerHTML = "连接已断开";
    }
    connect.onmessage = function(e) {
        connectInfoEl.innerHTML = e.data;
        // 绘图
        draw(JSON.parse(e.data))
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
            console.log(event)
            startX = event.touches[0].pageX;
            startY = event.touches[0].pageY;
        }
        function touchMove (event) {
            console.log(event)
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
    function draw(data) {
        ctx.lineWidth = 1;

        // 清除画布
        ctx.clearRect(0, 0, w, h);

        // 绘制地地图边框
        var Map = data.Map;
        for (var i = 0; i < Map.length; i++) {
            Map[i] = Map[i]*bodyW;
        };

        ctx.strokeRect(Map[0]-bodyW,Map[1]-bodyW,Map[2]+bodyW,Map[3]+bodyW);
        

        // 绘制蛇
        var Snakes = data.Snakes;
        ctx.fillStyle = "#000";
        ctx.strokeStyle = "#000";
        for (var i = Snakes.length - 1; i >= 0; i--) {
            var snake = Snakes[i],
                body = snake.Body,
                name = snake.Name,
                staust = snake.Status;

            for (var j = body.length - 1; j >= 0; j--) {
                var unit = body[j];
                drawRect(unit);
            };

        };

        // 绘制食物
        var Foods = data.Foods.Foods;
        for (var i = Foods.length - 1; i >= 0; i--) {
            var food = Foods[i];
            drawStrokeRect(food);
        };
    }

    // 绘制方格方法
    function drawRect(unit) {
        ctx.fillRect(unit[0] * bodyW, unit[1] * bodyW, bodyW - 1, bodyW - 1);
    }
    function drawStrokeRect(unit) {
        ctx.strokeRect(unit[0] * bodyW, unit[1] * bodyW, bodyW - 1, bodyW - 1);
    }


})();