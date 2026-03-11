package game

import (
	"fmt"
	"math/rand"
	"time"
)

// GetAliveFamily returns only alive family members
func (m *Model) GetAliveFamily() []FamilyMember {
	alive := make([]FamilyMember, 0)
	for _, member := range m.Family {
		if member.Alive {
			alive = append(alive, member)
		}
	}
	return alive
}

// AddMessage adds a message to the message log
func (m *Model) AddMessage(msg string) {
	m.Messages = append(m.Messages, msg)
	if len(m.Messages) > 10 {
		m.Messages = m.Messages[1:]
	}
	m.MessageTime = time.Now()
}

// AdvanceDay advances the game by one day
func (m *Model) AdvanceDay() {
	m.Day++
	m.Date = m.Date.AddDate(0, 0, 1)

	// Update season
	month := m.Date.Month()
	switch {
	case month >= 3 && month <= 5:
		m.Season = Spring
	case month >= 6 && month <= 8:
		m.Season = Summer
	case month >= 9 && month <= 11:
		m.Season = Fall
	default:
		m.Season = Winter
	}

	// Random weather change
	m.UpdateWeather()

	// Consume food
	aliveCount := len(m.GetAliveFamily())
	foodConsumed := aliveCount * m.Rations * 2 // 2-6 pounds per person
	m.Food -= foodConsumed
	if m.Food < 0 {
		m.Food = 0
		// Starving - health damage
		for i := range m.Family {
			if m.Family[i].Alive {
				m.Family[i].Health -= 10
				if m.Family[i].Health <= 0 {
					m.Family[i].Alive = false
					m.AddMessage(fmt.Sprintf("💔 %s 因饥饿而死", m.Family[i].Name))
				}
			}
		}
	}

	// Update diseases
	m.UpdateDiseases()

	// Travel
	milesTraveled := m.CalculateTravelDistance()
	m.MilesTraveled += milesTraveled
	m.MilesToNext = m.GetNextLandmark().Distance - m.MilesTraveled

	// Check if reached landmark
	if m.MilesToNext <= 0 {
		m.ReachLandmark()
	}

	// Check for random events
	if event := m.GenerateRandomEvent(); event != nil {
		m.CurrentEvent = event
		m.PrevState = m.State
		m.State = StateEvent
	}

	// Check win/lose conditions
	m.CheckGameState()
}

// UpdateWeather changes weather based on season and randomness
func (m *Model) UpdateWeather() {
	roll := rand.Float64()

	switch m.Season {
	case Spring:
		if roll < 0.5 {
			m.Weather = WeatherClear
		} else if roll < 0.7 {
			m.Weather = WeatherRain
		} else if roll < 0.85 {
			m.Weather = WeatherFog
		} else {
			m.Weather = WeatherStorm
		}
	case Summer:
		if roll < 0.6 {
			m.Weather = WeatherClear
		} else if roll < 0.75 {
			m.Weather = WeatherHot
		} else if roll < 0.9 {
			m.Weather = WeatherRain
		} else {
			m.Weather = WeatherStorm
		}
	case Fall:
		if roll < 0.4 {
			m.Weather = WeatherClear
		} else if roll < 0.6 {
			m.Weather = WeatherRain
		} else if roll < 0.8 {
			m.Weather = WeatherFog
		} else if roll < 0.95 {
			m.Weather = WeatherSnow
		} else {
			m.Weather = WeatherBlizzard
		}
	case Winter:
		if roll < 0.2 {
			m.Weather = WeatherClear
		} else if roll < 0.4 {
			m.Weather = WeatherSnow
		} else if roll < 0.7 {
			m.Weather = WeatherBlizzard
		} else {
			m.Weather = WeatherStorm
		}
	}
}

// CalculateTravelDistance calculates miles traveled in one day
func (m *Model) CalculateTravelDistance() int {
	baseMiles := 20

	// Pace modifier
	switch m.Pace {
	case 1: // steady
		baseMiles = 20
	case 2: // strenuous
		baseMiles = 30
	case 3: // grueling
		baseMiles = 40
	}

	// Weather modifier
	switch m.Weather {
	case WeatherRain, WeatherFog:
		baseMiles -= 5
	case WeatherSnow:
		baseMiles -= 10
	case WeatherBlizzard, WeatherStorm:
		baseMiles -= 15
	}

	// Ox health modifier
	baseMiles = baseMiles * m.OxHealth / 100

	// Wagon health modifier
	baseMiles = baseMiles * m.WagonHealth / 100

	// Random variation
	baseMiles += rand.Intn(10) - 5

	if baseMiles < 5 {
		baseMiles = 5
	}

	return baseMiles
}

// GetNextLandmark returns the next landmark
func (m *Model) GetNextLandmark() Landmark {
	for i, landmark := range m.Landmarks {
		if landmark.Distance > m.MilesTraveled {
			m.CurrentLandmark = i
			return landmark
		}
	}
	return m.Landmarks[len(m.Landmarks)-1]
}

// ReachLandmark handles reaching a landmark
func (m *Model) ReachLandmark() {
	landmark := m.Landmarks[m.CurrentLandmark]
	m.AddMessage(fmt.Sprintf("📍 到达 %s", landmark.Name))

	if landmark.HasRiver {
		m.PrevState = m.State
		m.State = StateRiver
	}

	if landmark.HasStore {
		m.AddMessage("🏪 这里可以购买补给")
	}
}

// UpdateDiseases processes ongoing diseases
func (m *Model) UpdateDiseases() {
	for i := range m.Family {
		if m.Family[i].Alive && m.Family[i].Disease != nil {
			m.Family[i].DiseaseDay++

			// Check if disease is fatal
			if m.Family[i].DiseaseDay >= m.Family[i].Disease.Duration {
				if rand.Float64() < m.Family[i].Disease.Fatality {
					m.Family[i].Alive = false
					m.AddMessage(fmt.Sprintf("💀 %s 因%s去世", m.Family[i].Name, m.Family[i].Disease.Name))
				} else {
					// Recovered
					m.Family[i].Disease = nil
					m.Family[i].DiseaseDay = 0
					m.Family[i].Health = 50
					m.AddMessage(fmt.Sprintf("✅ %s 从%s中康复", m.Family[i].Name, m.Family[i].Disease.Name))
				}
			}
		}
	}
}

// CheckGameState checks for win/lose conditions
func (m *Model) CheckGameState() {
	// Check if everyone is dead
	aliveCount := len(m.GetAliveFamily())
	if aliveCount == 0 {
		m.State = StateDeath
		return
	}

	// Check if reached Oregon
	if m.MilesTraveled >= 2000 {
		m.State = StateVictory
	}
}

// BuyItem purchases an item from the store
func (m *Model) BuyItem(itemIndex, quantity int) bool {
	if itemIndex < 0 || itemIndex >= len(ShopItems) {
		return false
	}

	item := ShopItems[itemIndex]
	totalCost := item.Price * quantity

	if totalCost > m.Money {
		return false
	}

	m.Money -= totalCost

	switch itemIndex {
	case 0: // Oxen
		m.Oxen += quantity
	case 1: // Food
		m.Food += quantity * 20
	case 2: // Clothing
		m.Clothing += quantity
	case 3: // Ammunition
		m.Ammunition += quantity * 20
	case 4: // Spare parts
		m.SpareParts += quantity
	}

	return true
}

// RestDay rests for a day (no travel)
func (m *Model) RestDay() {
	m.Day++
	m.Date = m.Date.AddDate(0, 0, 1)

	// Restore some health
	for i := range m.Family {
		if m.Family[i].Alive && m.Family[i].Health < 100 {
			m.Family[i].Health += 10
			if m.Family[i].Health > 100 {
				m.Family[i].Health = 100
			}
		}
	}

	// Restore ox health
	if m.OxHealth < 100 {
		m.OxHealth += 15
		if m.OxHealth > 100 {
			m.OxHealth = 100
		}
	}

	m.AddMessage("😴 休息了一天")
}

// ChangePace changes the travel pace
func (m *Model) ChangePace(pace int) {
	if pace >= 1 && pace <= 3 {
		m.Pace = pace
		paceNames := []string{"稳健", "努力", "拼命"}
		m.AddMessage(fmt.Sprintf("🚶 行进速度改为: %s", paceNames[pace-1]))
	}
}

// ChangeRations changes the food consumption rate
func (m *Model) ChangeRations(rations int) {
	if rations >= 1 && rations <= 3 {
		m.Rations = rations
		rationNames := []string{"充足", "节约", "极简"}
		m.AddMessage(fmt.Sprintf("🍖 食物配给改为: %s", rationNames[rations-1]))
	}
}

// GetStatus returns a formatted status string
func (m *Model) GetStatus() string {
	landmark := m.GetNextLandmark()
	return fmt.Sprintf(
		"日期: %s | 第%d天 | 天气: %s | 季节: %s\n"+
			"里程: %d/%d 英里 | 距离%s: %d 英里\n"+
			"食物: %d磅 | 牛: %d对 | 弹药: %d发 | 零件: %d个 | 衣服: %d套 | 金钱: $%d",
		m.Date.Format("2006年1月2日"),
		m.Day,
		m.Weather.String(),
		m.Season.String(),
		m.MilesTraveled,
		2000,
		landmark.Name,
		m.MilesToNext,
		m.Food,
		m.Oxen,
		m.Ammunition,
		m.SpareParts,
		m.Clothing,
		m.Money,
	)
}

// GetFamilyStatus returns family member status
func (m *Model) GetFamilyStatus() string {
	result := "家庭成员:\n"
	for _, member := range m.Family {
		status := "健康"
		if !member.Alive {
			status = "已故"
		} else if member.Disease != nil {
			status = fmt.Sprintf("患病(%s)", member.Disease.Name)
		} else if member.Health < 50 {
			status = "虚弱"
		}
		result += fmt.Sprintf("  %s - %s (健康度: %d%%)\n", member.Name, status, member.Health)
	}
	return result
}

// CalculateScore calculates final score
func (m *Model) CalculateScore() int {
	score := 0

	// Base score for arriving
	score += 1000

	// Bonus for surviving family members
	aliveCount := len(m.GetAliveFamily())
	score += aliveCount * 200

	// Bonus for remaining supplies
	score += m.Food
	score += m.Oxen * 20
	score += m.Clothing * 10
	score += m.Ammunition
	score += m.SpareParts * 10
	score += m.Money / 2

	// Bonus for profession
	switch m.Profession {
	case Farmer:
		score *= 3
	case Carpenter:
		score *= 2
	case Banker:
		// No multiplier
	}

	// Bonus for arrival date (earlier is better)
	if m.Date.Before(time.Date(1848, 9, 1, 0, 0, 0, 0, time.UTC)) {
		score += 500
	} else if m.Date.Before(time.Date(1848, 11, 1, 0, 0, 0, 0, time.UTC)) {
		score += 200
	}

	return score
}
