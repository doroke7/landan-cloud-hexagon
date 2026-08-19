#include <arpa/inet.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define PORT 9000

int main()
{
    int server_fd;
    int client_fd;
    struct sockaddr_in server_addr;
    struct sockaddr_in client_addr;
    socklen_t client_len = sizeof(client_addr);

    char buf[1024];

    // 建立 Socket
    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (server_fd < 0)
    {
        perror("socket");
        exit(1);
    }

    memset(&server_addr, 0, sizeof(server_addr));

    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(PORT);
    server_addr.sin_addr.s_addr = INADDR_ANY;

    // 綁定 Port
    if (bind(server_fd,
             (struct sockaddr *)&server_addr,
             sizeof(server_addr)) < 0)
    {
        perror("bind");
        exit(1);
    }

    // 開始監聽
    if (listen(server_fd, 128) < 0)
    {
        perror("listen");
        exit(1);
    }

    printf("TCP Server Start : %d\n", PORT);

    while (1)
    {
        printf("等待 Client...\n");

        client_fd = accept(
            server_fd,
            (struct sockaddr *)&client_addr,
            &client_len);

        if (client_fd < 0)
        {
            perror("accept");
            continue;
        }

        printf("Client Connected\n");

        while (1)
        {

            // 這是一個 經典 異步模式 reactor 模式： 
            // Reactor 就是把「監聽事件」委託給 Kernel（epoll/kqueue），當事件發生時，Kernel 通知我；收到通知後，我再主動去把資料讀出來。


            // 實現算法是 epoll 
            int n = read(client_fd, buf, sizeof(buf));

            if (n <= 0)
            {
                printf("Client Disconnect\n");
                break;
            }

            printf("Receive: %.*s\n", n, buf);

            write(client_fd, "OK", 2);
        }

        close(client_fd);
    }

    close(server_fd);

    return 0;
}