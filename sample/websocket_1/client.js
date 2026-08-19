const ws = new WebSocket(
    "ws://localhost:8080"
);


let id = 0;

const callbacks = new Map();


ws.onopen = () => {

    console.log("connected");


    emit(
        "login",
        {
            username: "test",
            password: "123"
        },
        (response) => {

            console.log(
                "ACK:",
                response
            );

        }
    );

};



function emit(action, data, callback) {

    const requestId = ++id;


    // 保存 callback
    callbacks.set(
        requestId,
        callback
    );


    ws.send(
        JSON.stringify({
            type: "event",
            id: requestId,
            action,
            data
        })
    );

}



ws.onmessage = (event) => {

    const packet = JSON.parse(
        event.data
    );


    if (packet.type === "ack") {

        const callback =
            callbacks.get(packet.id);


        if (callback) {

            callback(packet.data);

            callbacks.delete(packet.id);
        }

    }

};