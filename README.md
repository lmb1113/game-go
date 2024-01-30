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
![image](https://github.com/lmb1113/game-go/assets/39643887/5e7ad864-dcd8-466e-85c2-c904b1248829)
技能截图
![image](https://github.com/lmb1113/game-go/assets/39643887/8d54721e-6296-4b59-bad0-deb6ed13a160)
服务器列表
![image](https://github.com/lmb1113/game-go/assets/39643887/6b1a564f-f626-4c2d-9b28-8e106349e2dc)
