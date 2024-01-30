# game-go
用go写的联机小游戏，目前有点丑，打算是先实现功能再优化UI

##linux server build
```azure
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s'  -o ./build/game-server service/*go

```

## windows build 
```azure
go build -ldflags '-w -s'  -o ./build/game.exe main.go
```
## 1.0页面截图
![image](https://github.com/lmb1113/game-go/assets/39643887/feeb5a03-7883-424d-8f36-aa04509efc28)

## 服务器列表
![image](https://github.com/lmb1113/game-go/assets/39643887/d50c9e8d-404a-4437-8363-85dfb890cc71)

技能截图
![image](https://github.com/lmb1113/game-go/assets/39643887/cbca17aa-c55a-4e64-90eb-fc06049a0885)

