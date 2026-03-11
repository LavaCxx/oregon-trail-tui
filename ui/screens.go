package ui

import (
	"fmt"
	"strings"

	"github.com/LavaCxx/oregon-trail-tui/game"
	tea "github.com/charmbracelet/bubbletea"
)

// Model wraps the game model for UI
type Model struct {
	game *game.Model
}

// NewModel creates a new UI model
func NewModel(g *game.Model) Model {
	return Model{game: g}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.game.State == game.StateTitle {
				return m, tea.Quit
			}
		case "enter", " ":
			return m.handleSelect()
		case "up", "k":
			if m.game.SelectedIndex > 0 {
				m.game.SelectedIndex--
			}
		case "down", "j":
			m.game.SelectedIndex++
		case "left", "h":
			return m.handleLeft()
		case "right", "l":
			return m.handleRight()
		case "backspace":
			if m.game.State == game.StateNaming && len(m.game.InputBuffer) > 0 {
				m.game.InputBuffer = m.game.InputBuffer[:len(m.game.InputBuffer)-1]
			}
		default:
			// Handle text input for naming
			if m.game.State == game.StateNaming && len(msg.String()) == 1 {
				m.game.InputBuffer += msg.String()
			}
		}

	case tea.WindowSizeMsg:
		m.game.Width = msg.Width
		m.game.Height = msg.Height
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	switch m.game.State {
	case game.StateTitle:
		return m.viewTitle()
	case game.StateProfession:
		return m.viewProfession()
	case game.StateNaming:
		return m.viewNaming()
	case game.StateShopping:
		return m.viewShopping()
	case game.StateTraveling:
		return m.viewTraveling()
	case game.StateRiver:
		return m.viewRiver()
	case game.StateHunting:
		return m.viewHunting()
	case game.StateHuntingResult:
		return m.viewHuntingResult()
	case game.StateEvent:
		return m.viewEvent()
	case game.StateDeath:
		return m.viewDeath()
	case game.StateVictory:
		return m.viewVictory()
	case game.StateTombstone:
		return m.viewTombstone()
	default:
		return "未知状态"
	}
}

// ==================== Title Screen ====================

func (m Model) viewTitle() string {
	title := `
    ╔══════════════════════════════════════════════════════════════╗
    ║                                                              ║
    ║     ██████╗  █████╗ ███████╗████████╗██████╗  ██████╗██╗  ██╗ ║
    ║    ██╔════╝ ██╔══██╗██╔════╝╚══██╔══╝╚════██╗██╔════╝██║ ██╔╝ ║
    ║    ██║  ███╗███████║███████╗   ██║     █████╔╝██║     █████╔╝  ║
    ║    ██║   ██║██╔══██║╚════██║   ██║    ██╔═══╝ ██║     ██╔═██╗  ║
    ║    ╚██████╔╝██║  ██║███████║   ██║    ███████╗╚██████╗██║ ╚██╗ ║
    ║     ╚═════╝ ╚═╝  ╚═╝╚══════╝   ╚═╝    ╚══════╝ ╚═════╝╚═╝  ╚═╝ ║
    ║                                                              ║
    ║              ██████╗ ██████╗ ███████╗ █████╗ ████████╗██╗██╗   ██╗██╗██████╗ ███████╗ ║
    ║             ██╔════╝ ██╔══██╗██╔════╝██╔══██╗╚══██╔══╝██║╚██╗ ██╔╝██║██╔══██╗██╔════╝ ║
    ║             ██║  ███╗██████╔╝█████╗  ███████║   ██║   ██║ ╚████╔╝ ██║██████╔╝█████╗   ║
    ║             ██║   ██║██╔══██╗██╔══╝  ██╔══██║   ██║   ██║  ╚██╔╝  ██║██╔══██╗██╔══╝   ║
    ║             ╚██████╔╝██║  ██║███████╗██║  ██║   ██║   ██║   ██║   ██║██║  ██║███████╗ ║
    ║              ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝   ╚═╝   ╚═╝   ╚═╝   ╚═╝╚═╝  ╚═╝╚══════╝ ║
    ║                                                              ║
    ║                    🤠 俄勒冈之旅 TUI 版 🤠                    ║
    ║                                                              ║
    ║          1848年，从密苏里州到俄勒冈州的2000英里旅程           ║
    ║                                                              ║
    ╚══════════════════════════════════════════════════════════════╝
`

	menu := []string{
		"🎮 开始新游戏",
		"📖 游戏说明",
		"🚪 退出游戏",
	}

	var menuItems string
	for i, item := range menu {
		if i == m.game.SelectedIndex%len(menu) {
			menuItems += SelectedStyle.Render(fmt.Sprintf("  ➤ %s  ", item)) + "\n\n"
		} else {
			menuItems += MenuItemStyle.Render(fmt.Sprintf("    %s  ", item)) + "\n\n"
		}
	}

	return title + "\n" + CenterText(menuItems, 70) + "\n\n" + 
		MutedStyle.Render("      使用 ↑↓ 选择，Enter 确认，Q 退出      ")
}

// ==================== Profession Screen ====================

func (m Model) viewProfession() string {
	header := TitleStyle.Render("🧑‍🌾 选择你的职业")

	professions := []struct {
		name        string
		money       int
		description string
	}{
		{"🏦 银行家", 1600, "资金充裕，但最终得分较低"},
		{"🔨 木匠", 800, "资金适中，最终得分翻倍"},
		{"🌾 农夫", 400, "资金紧张，最终得分三倍"},
	}

	var items string
	for i, p := range professions {
		if i == m.game.SelectedIndex%len(professions) {
			items += SelectedStyle.Render(fmt.Sprintf("  ➤ %-10s $%-5d - %s  ", p.name, p.money, p.description)) + "\n\n"
		} else {
			items += MenuItemStyle.Render(fmt.Sprintf("    %-10s $%-5d - %s  ", p.name, p.money, p.description)) + "\n\n"
		}
	}

	info := BoxStyle.Render("💡 提示：资金越少难度越高，但最终得分倍率越大")

	return "\n" + header + "\n\n" + items + "\n" + info + "\n\n" + 
		MutedStyle.Render("使用 ↑↓ 选择，Enter 确认")
}

// ==================== Naming Screen ====================

func (m Model) viewNaming() string {
	header := TitleStyle.Render("👨‍👩‍👧‍👦 命名你的家庭成员")

	var currentNaming string
	if m.game.PlayerName == "" {
		currentNaming = "领队"
	} else {
		currentNaming = fmt.Sprintf("家庭成员 %d", len(m.game.Family)+1)
	}

	var namedList string
	if m.game.PlayerName != "" {
		namedList = "\n已命名成员：\n"
		namedList += fmt.Sprintf("  👤 领队：%s\n", m.game.PlayerName)
		for _, member := range m.game.Family {
			namedList += fmt.Sprintf("  👤 成员：%s\n", member.Name)
		}
	}

	inputBox := InputStyle.Render(fmt.Sprintf("请输入 %s 的名字：%s", currentNaming, m.game.InputBuffer+"█"))

	info := "\n" + MutedStyle.Render("提示：最多5名家庭成员（领队 + 4名成员），输入完毕按 Enter 确认，留空结束命名")

	return "\n" + header + namedList + "\n\n" + inputBox + info + "\n"
}

// ==================== Shopping Screen ====================

func (m Model) viewShopping() string {
	header := TitleStyle.Render("🏪 独立城商店")

	// Money display
	moneyDisplay := HeaderStyle.Render(fmt.Sprintf("💰 你的资金：$%d", m.game.Money))

	// Shop items
	var items string
	shopItems := game.ShopItems
	for i, item := range shopItems {
		quantity := m.getQuantityForItem(i)
		line := fmt.Sprintf("%-8s | $%-3d/%-8s | %s | 已购买：%d", 
			item.Name, item.Price, item.Unit, item.Description, quantity)
		
		if i == m.game.SelectedIndex%len(shopItems) {
			items += SelectedStyle.Render(fmt.Sprintf("  ➤ %s  ", line)) + "\n\n"
		} else {
			items += MenuItemStyle.Render(fmt.Sprintf("    %s  ", line)) + "\n\n"
		}
	}

	// Controls
	controls := BoxStyle.Render("⬅️➡️ 调整数量 | Enter 完成购物")

	// Current supplies
	supplies := fmt.Sprintf("\n已购物资：牛%d对 | 食物%d磅 | 衣服%d套 | 弹药%d发 | 零件%d个",
		m.game.Oxen, m.game.Food, m.game.Clothing, m.game.Ammunition, m.game.SpareParts)
	suppliesDisplay := InfoStyle.Render(supplies)

	return "\n" + header + "\n\n" + moneyDisplay + "\n\n" + items + suppliesDisplay + "\n" + controls + "\n"
}

func (m Model) getQuantityForItem(index int) int {
	switch index {
	case 0:
		return m.game.Oxen
	case 1:
		return m.game.Food / 20
	case 2:
		return m.game.Clothing
	case 3:
		return m.game.Ammunition / 20
	case 4:
		return m.game.SpareParts
	}
	return 0
}

// ==================== Traveling Screen ====================

func (m Model) viewTraveling() string {
	// Status header
	status := m.game.GetStatus()
	statusBox := StatusBoxStyle.Render(status)

	// Family status
	familyStatus := m.game.GetFamilyStatus()
	familyBox := StatusBoxStyle.Render(familyStatus)

	// Messages
	var messages string
	if len(m.game.Messages) > 0 {
		messages = "\n📋 事件记录：\n"
		for i, msg := range m.game.Messages {
			if i >= len(m.game.Messages)-5 {
				messages += fmt.Sprintf("  %s\n", msg)
			}
		}
	}

	// Menu
	menu := []string{
		"🚶 继续前进",
		"🎯 狩猎",
		"😴 休息",
		"⚙️ 调整速度",
		"🍖 调整配给",
		"📊 查看状态",
	}

	var menuItems string
	for i, item := range menu {
		if i == m.game.SelectedIndex%len(menu) {
			menuItems += SelectedStyle.Render(fmt.Sprintf("  ➤ %s  ", item)) + "\n"
		} else {
			menuItems += MenuItemStyle.Render(fmt.Sprintf("    %s  ", item)) + "\n"
		}
	}

	// Progress bar
	progress := ProgressBar(m.game.MilesTraveled, 2000, 40)
	progressDisplay := "\n" + InfoStyle.Render(fmt.Sprintf("旅程进度：%s %d%%", progress, m.game.MilesTraveled/20))

	return "\n" + statusBox + "\n" + familyBox + messages + "\n" + menuItems + progressDisplay + "\n"
}

// ==================== River Screen ====================

func (m Model) viewRiver() string {
	header := TitleStyle.Render("🌊 河流穿越")

	// River info
	riverInfo := m.game.GetRiverInfo()

	// Crossing options
	options := m.game.GetCrossingOptions()
	var optionsDisplay string
	for i, opt := range options {
		if i == m.game.SelectedIndex%len(options) {
			optionsDisplay += SelectedStyle.Render(fmt.Sprintf("  %s  ", opt)) + "\n\n"
		} else {
			optionsDisplay += MenuItemStyle.Render(fmt.Sprintf("  %s  ", opt)) + "\n\n"
		}
	}

	// Result
	var result string
	if m.game.RiverResult != "" {
		result = "\n" + BoxStyle.Render(m.game.RiverResult)
	}

	// Current status
	status := fmt.Sprintf("💰 金钱：$%d | 🐮 牛：%d对 | 🍖 食物：%d磅", 
		m.game.Money, m.game.Oxen, m.game.Food)
	statusDisplay := InfoStyle.Render(status)

	return "\n" + header + "\n" + riverInfo + "\n" + optionsDisplay + statusDisplay + result + "\n"
}

// ==================== Hunting Screen ====================

func (m Model) viewHunting() string {
	header := TitleStyle.Render("🏹 狩猎")

	// Hunting scene
	scene := m.game.HuntingScene()

	// Status
	status := m.game.GetHuntingStatus()
	statusDisplay := InfoStyle.Render(status)

	// Animals to shoot
	var animals string
	for i, animal := range m.game.HuntingAnimals {
		if i == m.game.SelectedIndex%max(1, len(m.game.HuntingAnimals)) {
			animals += SelectedStyle.Render(fmt.Sprintf("  ➤ [%d] %s %s (%d磅)  ", 
				i+1, animal.Symbol, animal.Name, animal.Weight)) + "\n"
		} else {
			animals += MenuItemStyle.Render(fmt.Sprintf("    [%d] %s %s (%d磅)  ", 
				i+1, animal.Symbol, animal.Name, animal.Weight)) + "\n"
		}
	}

	// Controls
	controls := MutedStyle.Render("按 1-9 选择目标射击 | Enter 结束狩猎 | 剩余射击次数会消耗弹药")

	return "\n" + header + "\n" + scene + "\n" + statusDisplay + "\n\n" + animals + "\n" + controls
}

// ==================== Hunting Result Screen ====================

func (m Model) viewHuntingResult() string {
	header := TitleStyle.Render("🎯 狩猎结束")

	result := fmt.Sprintf(`
    ╔══════════════════════════════════════════════════════════════╗
    ║                                                              ║
    ║     🏹 狩猎统计                                              ║
    ║                                                              ║
    ║     获得食物：%-6d 磅                                   ║
    ║     剩余弹药：%-6d 发                                    ║
    ║                                                              ║
    ╚══════════════════════════════════════════════════════════════╝
`, m.game.HuntingFood, m.game.Ammunition)

	foodDisplay := fmt.Sprintf("\n总食物储备：%d 磅", m.game.Food)
	
	continuePrompt := "\n" + MutedStyle.Render("按 Enter 继续旅程...")

	return "\n" + header + result + InfoStyle.Render(foodDisplay) + continuePrompt
}

// ==================== Event Screen ====================

func (m Model) viewEvent() string {
	if m.game.CurrentEvent == nil {
		return "没有事件"
	}

	event := m.game.CurrentEvent
	header := TitleStyle.Render(fmt.Sprintf("⚠️ %s", event.Title))

	// Event description box
	description := BoxStyle.Render(event.Description)

	// Event message if any
	var message string
	if m.game.EventMessage != "" {
		message = "\n" + WarningStyle.Render(m.game.EventMessage)
	}

	continuePrompt := "\n" + MutedStyle.Render("按 Enter 继续...")

	return "\n" + header + "\n" + description + message + continuePrompt
}

// ==================== Death Screen ====================

func (m Model) viewDeath() string {
	titleContent := fmt.Sprintf(`
    ╔══════════════════════════════════════════════════════════════╗
    ║                                                              ║
    ║     💀💀💀  全员死亡  💀💀💀                                   ║
    ║                                                              ║
    ║     你的旅程在第 %-4d 天结束                               ║
    ║     行进了 %-5d 英里                                       ║
    ║                                                              ║
    ╚══════════════════════════════════════════════════════════════╝
`, m.game.Day, m.game.MilesTraveled)
	title := DangerStyle.Render(titleContent)

	// Tombstone
	tombstone := AsciiStyle.Render(`
        _______________
       /               \
      /                 \
     |    R.I.P.        |
     |                   |
     |    这里长眠着     |
     |    你的旅队       |
     |                   |
     |    他们没能       |
     |    到达俄勒冈     |
     |___________________|
            |   |
    ________|___|________
   /                      \
  /________________________\
`)

	// Death causes
	var causes string
	causes += "\n死亡记录：\n"
	for _, member := range m.game.Family {
		if !member.Alive {
			cause := "未知原因"
			if member.Disease != nil {
				cause = member.Disease.Name
			}
			causes += fmt.Sprintf("  💀 %s - %s\n", member.Name, cause)
		}
	}

	return "\n\n" + title + causes + "\n" + tombstone + "\n\n" + 
		MutedStyle.Render("按 Enter 返回主菜单...")
}

// ==================== Victory Screen ====================

func (m Model) viewVictory() string {
	score := m.game.CalculateScore()

	title := TitleStyle.Render(`
    ╔══════════════════════════════════════════════════════════════╗
    ║                                                              ║
    ║     🎉🎉🎉  恭喜到达俄勒冈！  🎉🎉🎉                           ║
    ║                                                              ║
    ║     你成功完成了2000英里的旅程！                              ║
    ║                                                              ║
    ╚══════════════════════════════════════════════════════════════╝
`)

	// Statistics
	stats := fmt.Sprintf(`
    ╔══════════════════════════════════════════════════════════════╗
    ║                    📊 最终统计                               ║
    ╠══════════════════════════════════════════════════════════════╣
    ║  到达日期：%-48s║
    ║  旅程天数：%-48d║
    ║  存活人数：%-48d║
    ║  剩余资金：$%-47d║
    ║  剩余食物：%-48d磅
    ║  剩余牛：%-48d对
    ║                                                              ║
    ║  ⭐ 最终得分：%-46d⭐
    ║                                                              ║
    ╚══════════════════════════════════════════════════════════════╝
`, m.game.Date.Format("1848年1月2日"), m.game.Day, len(m.game.GetAliveFamily()),
		m.game.Money, m.game.Food, m.game.Oxen, score)

	// Surviving family
	var survivors string
	survivors = "\n存活的家庭成员：\n"
	for _, member := range m.game.Family {
		if member.Alive {
			survivors += fmt.Sprintf("  ✅ %s (健康度：%d%%)\n", member.Name, member.Health)
		}
	}

	return "\n\n" + title + stats + survivors + "\n" + 
		MutedStyle.Render("按 Enter 返回主菜单...")
}

// ==================== Tombstone Screen ====================

func (m Model) viewTombstone() string {
	return m.viewDeath()
}

// ==================== Event Handlers ====================

func (m Model) handleSelect() (tea.Model, tea.Cmd) {
	switch m.game.State {
	case game.StateTitle:
		return m.handleTitleSelect()
	case game.StateProfession:
		return m.handleProfessionSelect()
	case game.StateNaming:
		return m.handleNamingConfirm()
	case game.StateShopping:
		return m.handleShoppingConfirm()
	case game.StateTraveling:
		return m.handleTravelingSelect()
	case game.StateRiver:
		return m.handleRiverSelect()
	case game.StateHunting:
		return m.handleHuntingSelect()
	case game.StateHuntingResult:
		m.game.State = game.StateTraveling
		return m, nil
	case game.StateEvent:
		m.game.State = m.game.PrevState
		m.game.CurrentEvent = nil
		m.game.EventMessage = ""
		return m, nil
	case game.StateDeath, game.StateVictory:
		// Reset game
		*m.game = game.NewModel()
		return m, nil
	}
	return m, nil
}

func (m Model) handleTitleSelect() (tea.Model, tea.Cmd) {
	index := m.game.SelectedIndex % 3
	switch index {
	case 0: // New game
		m.game.State = game.StateProfession
		m.game.SelectedIndex = 0
	case 1: // Instructions
		// TODO: Show instructions
	case 2: // Quit
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleProfessionSelect() (tea.Model, tea.Cmd) {
	index := m.game.SelectedIndex % 3
	m.game.Profession = game.Profession(index)
	m.game.Money = m.game.Profession.StartingMoney()
	m.game.State = game.StateNaming
	m.game.SelectedIndex = 0
	m.game.InputBuffer = ""
	return m, nil
}

func (m Model) handleNamingConfirm() (tea.Model, tea.Cmd) {
	name := strings.TrimSpace(m.game.InputBuffer)
	
	if m.game.PlayerName == "" {
		// Naming the leader
		if name == "" {
			name = "领队"
		}
		m.game.PlayerName = name
		m.game.Family = append(m.game.Family, game.FamilyMember{
			Name:   name,
			Health: 100,
			Alive:  true,
		})
	} else {
		// Naming family member
		if name == "" || len(m.game.Family) >= 5 {
			// Done naming
			m.game.State = game.StateShopping
			m.game.SelectedIndex = 0
			return m, nil
		}
		m.game.Family = append(m.game.Family, game.FamilyMember{
			Name:   name,
			Health: 100,
			Alive:  true,
		})
	}
	
	m.game.InputBuffer = ""
	return m, nil
}

func (m Model) handleShoppingConfirm() (tea.Model, tea.Cmd) {
	// Check if has enough oxen
	if m.game.Oxen < 2 {
		m.game.AddMessage("⚠️ 你需要至少2对牛才能拉车！")
		return m, nil
	}
	
	m.game.State = game.StateTraveling
	m.game.SelectedIndex = 0
	m.game.MilesToNext = m.game.Landmarks[1].Distance
	m.game.AddMessage("🚀 旅程开始！向俄勒冈出发！")
	return m, nil
}

func (m Model) handleTravelingSelect() (tea.Model, tea.Cmd) {
	index := m.game.SelectedIndex % 6
	switch index {
	case 0: // Continue
		m.game.AdvanceDay()
	case 1: // Hunt
		if m.game.Ammunition > 0 {
			m.game.StartHunting()
		} else {
			m.game.AddMessage("没有弹药，无法狩猎！")
		}
	case 2: // Rest
		m.game.RestDay()
	case 3: // Change pace
		m.game.Pace = (m.game.Pace % 3) + 1
		paceNames := []string{"稳健", "努力", "拼命"}
		m.game.AddMessage(fmt.Sprintf("行进速度改为：%s", paceNames[m.game.Pace-1]))
	case 4: // Change rations
		m.game.Rations = (m.game.Rations % 3) + 1
		rationNames := []string{"充足", "节约", "极简"}
		m.game.AddMessage(fmt.Sprintf("食物配给改为：%s", rationNames[m.game.Rations-1]))
	case 5: // View status
		// Already shown in view
	}
	return m, nil
}

func (m Model) handleRiverSelect() (tea.Model, tea.Cmd) {
	if m.game.RiverResult != "" {
		// Already crossed, continue
		m.game.State = game.StateTraveling
		m.game.RiverResult = ""
		m.game.CurrentLandmark++
		m.game.MilesToNext = m.game.GetNextLandmark().Distance - m.game.MilesTraveled
		return m, nil
	}
	
	index := m.game.SelectedIndex % 4
	result := m.game.HandleRiverCrossing(index + 1)
	m.game.AddMessage(result)
	
	if m.game.State != game.StateDeath {
		m.game.RiverResult = result
	}
	return m, nil
}

func (m Model) handleHuntingSelect() (tea.Model, tea.Cmd) {
	if len(m.game.HuntingAnimals) == 0 {
		m.game.EndHunting()
		return m, nil
	}
	
	index := m.game.SelectedIndex % len(m.game.HuntingAnimals)
	if index < len(m.game.HuntingAnimals) {
		result := m.game.Shoot(index)
		m.game.AddMessage(result)
	}
	
	if m.game.HuntingShots <= 0 || len(m.game.HuntingAnimals) == 0 {
		m.game.EndHunting()
	}
	return m, nil
}

func (m Model) handleLeft() (tea.Model, tea.Cmd) {
	if m.game.State == game.StateShopping {
		// Decrease quantity
		index := m.game.SelectedIndex % len(game.ShopItems)
		item := game.ShopItems[index]
		if item.Price <= m.game.Money {
			// We can afford to decrease (meaning we've bought some)
			m.game.BuyItem(index, -1)
			// Actually, we need to implement a sell function or just not allow decreasing
			// For simplicity, let's not allow decreasing
		}
	}
	return m, nil
}

func (m Model) handleRight() (tea.Model, tea.Cmd) {
	if m.game.State == game.StateShopping {
		// Increase quantity
		index := m.game.SelectedIndex % len(game.ShopItems)
		item := game.ShopItems[index]
		if item.Price <= m.game.Money {
			m.game.BuyItem(index, 1)
		} else {
			m.game.AddMessage("资金不足！")
		}
	}
	return m, nil
}
