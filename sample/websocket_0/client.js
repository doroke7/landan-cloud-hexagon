const ws = new WebSocket(
    "ws://localhost:8080"
);


ws.onopen = () => {

    console.log("connected");

    ws.send("hello server");

};

    // 原生的 websocket 非常陽春，他只能
    // 1. 不知道這個是某人 第幾個 發的消息 
    // 2. 沒有帶 ack機制
    // 3. server 接受消息，不知道這個是什麼用途
    // 4. 客戶端沒有重試



    // 1. 不知道這個是某人 第幾個 發的消息 => 加上 request-id
    // 2. 沒有帶 ack機制 =》 透過 server 反向 request 實現 response + client 執行 callback（需要搭配 request-id, type=ack）
    // 3. server 接受消息，不知道這個是什麼用途 =》 server 加上（method | action），
    // 4. 客戶端沒有重試 =》 客戶端優化，如果  沒有收到 ack 就定時發送

    // 以上就完備 client 對 server 呼叫

    // 1. client 這邊加上 路由機制，讓server 主動呼叫
    // 2. sever 那邊增加 conn.Send 方法，能夠主動發非ack消息的 event 消息
    // 3. server 這邊暫時不做ack callback

ws.onmessage = (event) => {

    console.log(
        "receive:",
        event.data
    );

};


ws.onclose = () => {

    console.log("closed");

};