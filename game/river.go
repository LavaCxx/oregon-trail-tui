package game

import (
	"fmt"
	"math/rand"
)

// RiverCrossing handles river crossing logic
type RiverCrossing struct {
	Depth  int
	Width  int
	Weather Weather
}

// CrossingResult represents the result of a river crossing
type CrossingResult struct {
	Success     bool
	Message     string
	FoodLost    int
	OxenLost    int
	PeopleLost  int
	PartsLost   int
}

// NewRiverCrossing creates a new river crossing scenario
func NewRiverCrossing(depth, width int, weather Weather) *RiverCrossing {
	return &RiverCrossing{
		Depth:  depth,
		Width:  width,
		Weather: weather,
	}
}

// CrossFerry uses a ferry to cross (costs money)
func (r *RiverCrossing) CrossFerry(m *Model) CrossingResult {
	cost := r.Depth * 5 // $5 per foot of depth

	if m.Money < cost {
		return CrossingResult{
			Success: false,
			Message: fmt.Sprintf("渡船费用需要 $%d，但你只有 $%d！", cost, m.Money),
		}
	}

	m.Money -= cost
	return CrossingResult{
		Success: true,
		Message: fmt.Sprintf("支付 $%d 后，安全渡过了河流。", cost),
	}
}

// CrossSwim attempts to swim across (dangerous)
func (r *RiverCrossing) CrossSwim(m *Model) CrossingResult {
	result := CrossingResult{Success: true}

	// Base success chance
	successChance := 70

	// Modify based on depth
	if r.Depth > 6 {
		successChance -= (r.Depth - 6) * 10
	}

	// Modify based on weather
	if r.Weather == WeatherRain || r.Weather == WeatherStorm {
		successChance -= 20
	}
	if r.Weather == WeatherSnow || r.Weather == WeatherBlizzard {
		successChance -= 30
	}

	// Ensure minimum chance
	if successChance < 10 {
		successChance = 10
	}

	// Roll for each ox
	oxenLost := 0
	for i := 0; i < m.Oxen; i++ {
		if rand.Intn(100) > successChance {
			oxenLost++
		}
	}
	m.Oxen -= oxenLost

	// Roll for food loss
	if rand.Intn(100) > successChance {
		foodLost := rand.Intn(50) + 20
		if foodLost > m.Food {
			foodLost = m.Food
		}
		m.Food -= foodLost
		result.FoodLost = foodLost
	}

	// Roll for each family member
	for i := range m.Family {
		if m.Family[i].Alive && rand.Intn(100) > successChance+20 { // People are better at swimming
			m.Family[i].Health -= rand.Intn(30) + 20
			if m.Family[i].Health <= 0 {
				m.Family[i].Health = 0
				m.Family[i].Alive = false
				result.PeopleLost++
			}
		}
	}

	// Build message
	if oxenLost > 0 || result.FoodLost > 0 || result.PeopleLost > 0 {
		result.Message = "渡河时遇到了麻烦！"
		if oxenLost > 0 {
			result.Message += fmt.Sprintf(" 损失了 %d 头牛。", oxenLost)
		}
		if result.FoodLost > 0 {
			result.Message += fmt.Sprintf(" 损失了 %d 磅食物。", result.FoodLost)
		}
		if result.PeopleLost > 0 {
			result.Message += fmt.Sprintf(" %d 人溺水身亡！", result.PeopleLost)
		}
	} else {
		result.Message = "虽然惊险，但所有人都安全渡过了河流！"
	}

	return result
}

// CrossWait waits for better conditions (costs time)
func (r *RiverCrossing) CrossWait(m *Model) CrossingResult {
	// Wait 1-3 days
	days := rand.Intn(3) + 1

	// Consume food for waiting days
	foodNeeded := len(m.GetAliveFamily()) * 3 * days
	if foodNeeded > m.Food {
		foodNeeded = m.Food
	}
	m.Food -= foodNeeded

	// Advance time
	for i := 0; i < days; i++ {
		m.Day++
		m.Date = m.Date.AddDate(0, 0, 1)
	}

	// Improve conditions slightly
	improvedDepth := r.Depth - rand.Intn(2)
	if improvedDepth < 2 {
		improvedDepth = 2
	}

	// Now try to cross with better conditions
	r.Depth = improvedDepth
	return r.CrossSwim(m)
}

// CrossGuide hires an Indian guide
func (r *RiverCrossing) CrossGuide(m *Model) CrossingResult {
	cost := 20 // Fixed cost for guide

	if m.Money < cost {
		return CrossingResult{
			Success: false,
			Message: fmt.Sprintf("雇佣向导需要 $%d，但你只有 $%d！", cost, m.Money),
		}
	}

	m.Money -= cost

	// Very high success rate with guide
	successChance := 95

	// Small chance of minor losses
	if rand.Intn(100) > successChance {
		foodLost := rand.Intn(20)
		m.Food -= foodLost
		return CrossingResult{
			Success: true,
			Message: fmt.Sprintf("在向导的帮助下安全渡河，但损失了 %d 磅食物。", foodLost),
		}
	}

	return CrossingResult{
		Success: true,
		Message: "在向导的带领下，所有人安全渡过了河流！",
	}
}

// GetRiverInfo returns formatted river information
func (m *Model) GetRiverInfo() string {
	if m.CurrentLandmark >= len(m.Landmarks) {
		return ""
	}

	landmark := m.Landmarks[m.CurrentLandmark]
	if !landmark.HasRiver {
		return ""
	}

	info := fmt.Sprintf(`
    ╔══════════════════════════════════════════════════════════════╗
    ║                      🌊 河流穿越                              ║
    ╠══════════════════════════════════════════════════════════════╣
    ║  河流: %-52s║
    ║  深度: %-2d 英尺                                              ║
    ║  宽度: %-3d 英尺                                              ║
    ║  天气: %-52s║
    ╚══════════════════════════════════════════════════════════════╝
`, landmark.Name, landmark.RiverDepth, landmark.RiverWidth, m.Weather.String())

	return info
}

// GetCrossingOptions returns available crossing options
func (m *Model) GetCrossingOptions() []string {
	landmark := m.Landmarks[m.CurrentLandmark]
	ferryCost := landmark.RiverDepth * 5
	guideCost := 20

	options := []string{
		fmt.Sprintf("[1] 渡船 - 花费 $%d (最安全)", ferryCost),
		"[2] 游泳 - 免费 (危险)",
		"[3] 等待 - 消耗时间等待水位下降",
		fmt.Sprintf("[4] 雇佣向导 - 花费 $%d (较安全)", guideCost),
	}

	return options
}

// HandleRiverCrossing processes the player's river crossing choice
func (m *Model) HandleRiverCrossing(choice int) string {
	landmark := m.Landmarks[m.CurrentLandmark]
	crossing := NewRiverCrossing(landmark.RiverDepth, landmark.RiverWidth, m.Weather)

	var result CrossingResult

	switch choice {
	case 1:
		result = crossing.CrossFerry(m)
	case 2:
		result = crossing.CrossSwim(m)
	case 3:
		result = crossing.CrossWait(m)
	case 4:
		result = crossing.CrossGuide(m)
	default:
		return "无效的选择！"
	}

	m.RiverResult = result.Message
	return result.Message
}
