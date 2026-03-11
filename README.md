# 俄勒冈之旅 TUI (Oregon Trail TUI)

🎮 经典游戏《俄勒冈之旅》的 TUI 复刻版，使用 Go + Bubble Tea 开发。

## 游戏简介

《俄勒冈之旅》是1971年开发的经典教育游戏，模拟1848年美国移民从密苏里州独立城到俄勒冈州威拉米特谷的2000英里旅程。

## 游戏特色

- 🧑‍🌾 **职业选择**：银行家、木匠、农夫，不同职业有不同起始资金和得分倍率
- 🏪 **商店系统**：购买牛、食物、衣服、弹药、零件
- 🗺️ **旅程系统**：穿越多个地标，体验天气和季节变化
- ⚠️ **随机事件**：疾病、意外、盗贼袭击等
- 🌊 **河流穿越**：选择不同方式渡河
- 🏹 **狩猎系统**：狩猎获取食物
- 📊 **资源管理**：管理食物配给和行进速度

## 运行方式

```bash
# 设置 Go 代理（中国大陆）
go env -w GOPROXY=https://goproxy.cn,direct

# 运行游戏
go run main.go
```

## 操作说明

- `↑/k` - 上移
- `↓/j` - 下移
- `←/h` - 左移/减少
- `→/l` - 右移/增加
- `Enter/Space` - 确认
- `Q/Ctrl+C` - 退出

## 技术栈

- **语言**: Go 1.21+
- **TUI 框架**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **样式**: [Lipgloss](https://github.com/charmbracelet/lipgloss)

## 致敬

本项目致敬1971年原版《The Oregon Trail》，由 MECC (Minnesota Educational Computing Consortium) 开发。

## License

MIT
