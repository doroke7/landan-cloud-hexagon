const WebSocket = require("ws");

const server = new WebSocket.Server({
    port: 8080
});


server.on("connection", (socket) => {

    console.log("client connected");


    socket.on("message", (message) => {

        const packet = JSON.parse(message);


        console.log(
            "action:",
            packet.action
        );


        // 模擬業務處理
        if (packet.action === "login") {

            const oResponse = {
                type: "ack",
                id: packet.id,
                result: {
                    success: true,
                    userId: 100
                }
            };


            socket.send(
                JSON.stringify(oResponse)
            );

        }

    });

});