// 原生的 websocket 非常陽春，他只能
// 1. 不知道這個是某人 第幾個 發的消息 => 加上 request-id
// 2. 沒有帶 ack機制 =》 透過 server 反向 request 實現 + client 執行 callback（需要搭配 request-id）
// 3. server 接受消息，不知道這個是什麼用途 =》 加上（method | action），

class WSClient {
    constructor(url) {
        this.ws = new WebSocket(url);
        this.requestId = 0;
        this.callbacks = new Map();

        this.ws.onmessage = (event) => {
            const packet = JSON.parse(event.data);

            // 收到 ACK
            if (packet.type === "ack") {
                const callback = this.callbacks.get(packet.id);

                if (callback) {
                    callback(packet.data);
                    this.callbacks.delete(packet.id);
                }

                return;
            }

            // 普通事件
            console.log("event:", packet);
        };
    }


    emit(sMethod, oRequest, fnAck) {

        const iRequestId = ++this.requestId;


        if (fnAck) {
            this.callbacks.set(iRequestId, fnAck);
        }


        this.ws.send(JSON.stringify({
            type: "event",
            id: iRequestId,
            method: sMethod,
            data: oRequest
        }));
    }
}


// 使用

const client = new WSClient(
    "ws://localhost:8080/ws"
);


client.ws.onopen = () => {

    client.emit(
        "login",
        {
            username: "test",
            password: "123"
        },
        (response) => {
            console.log("login ack:", response);
        }
    );

};