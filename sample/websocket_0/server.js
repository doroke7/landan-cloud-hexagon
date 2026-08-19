const WebSocket = require("ws");

const server = new WebSocket.Server({
    port: 8080
});


server.on("connection", (socket) => {

    console.log("client connected");



    // 原生的 websocket 非常陽春，他只能
    // 1. 不知道這個是某人 第幾個 發的消息 => 加上 request-id
    // 2. 沒有帶 ack機制 =》 透過 server 反向 request 實現 + client 執行 callback（需要搭配 request-id）
    // 3. server 接受消息，不知道這個是什麼用途 =》 加上（method | action），

    socket.on("message", (message) => {

        console.log(
            "recv:",
            message.toString()
        );


        // 回覆 client
        socket.send(
            "server reply: "
        );

    });


    socket.on("close", () => {
        console.log("client closed");
    });

});


console.log("WebSocket server :8080");