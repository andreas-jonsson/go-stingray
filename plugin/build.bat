go build -buildmode=c-archive loader\loader.go
gcc -shared -pthread -o loader.dll loader\loader.c loader.a -lWinMM -lntdll -lWS2_32
