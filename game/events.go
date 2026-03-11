package game

import (
	"fmt"
	"math/rand"
)

// GenerateRandomEvent creates a random event based on conditions
func (m *Model) GenerateRandomEvent() *GameEvent {
	// Base chance of event per day
	eventChance := 0.3

	// Modify based on conditions
	if m.Weather == WeatherBlizzard || m.Weather == WeatherStorm {
		eventChance += 0.2
	}
	if m.Food < 50 {
		eventChance += 0.1
	}
	if m.Clothing < len(m.GetAliveFamily()) {
		eventChance += 0.15
	}

	if rand.Float64() > eventChance {
		return nil
	}

	// Choose event type
	eventTypes := []EventType{
		EventDisease,
		EventAccident,
		EventThief,
		EventOxDeath,
		EventWagonDamage,
		EventFound,
	}

	// Weight events based on conditions
	weights := []int{
		30, // disease
		15, // accident
		10, // thief
		15, // ox death
		15, // wagon damage
		15, // found
	}

	// Increase disease chance if low clothing
	if m.Clothing < len(m.GetAliveFamily()) {
		weights[0] += 20
	}

	// Increase ox death if ox health low
	if m.OxHealth < 50 {
		weights[3] += 15
	}

	// Increase wagon damage if wagon health low
	if m.WagonHealth < 50 {
		weights[4] += 15
	}

	// Select weighted random
	totalWeight := 0
	for _, w := range weights {
		totalWeight += w
	}
	r := rand.Intn(totalWeight)
	cumulative := 0
	var selectedType EventType
	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			selectedType = eventTypes[i]
			break
		}
	}

	return m.createEvent(selectedType)
}

func (m *Model) createEvent(eventType EventType) *GameEvent {
	switch eventType {
	case EventDisease:
		return m.createDiseaseEvent()
	case EventAccident:
		return m.createAccidentEvent()
	case EventThief:
		return m.createThiefEvent()
	case EventOxDeath:
		return m.createOxDeathEvent()
	case EventWagonDamage:
		return m.createWagonDamageEvent()
	case EventFound:
		return m.createFoundEvent()
	default:
		return nil
	}
}

func (m *Model) createDiseaseEvent() *GameEvent {
	aliveFamily := m.GetAliveFamily()
	if len(aliveFamily) == 0 {
		return nil
	}

	victim := aliveFamily[rand.Intn(len(aliveFamily))]
	diseaseKeys := []string{"dysentery", "typhoid", "cholera", "measles", "fever"}
	diseaseKey := diseaseKeys[rand.Intn(len(diseaseKeys))]
	disease := Diseases[diseaseKey]

	// Check if already has disease
	for i := range m.Family {
		if m.Family[i].Name == victim.Name && m.Family[i].Disease == nil {
			m.Family[i].Disease = &disease
			m.Family[i].DiseaseDay = 0

			return &GameEvent{
				Type:        EventDisease,
				Title:       "疾病爆发！",
				Description: fmt.Sprintf("%s 患上了 %s！\n%s", victim.Name, disease.Name, disease.Description),
			}
		}
	}
	return nil
}

func (m *Model) createAccidentEvent() *GameEvent {
	aliveFamily := m.GetAliveFamily()
	if len(aliveFamily) == 0 {
		return nil
	}

	victim := aliveFamily[rand.Intn(len(aliveFamily))]
	accidents := []struct {
		name   string
		damage int
	}{
		{"被蛇咬伤", 25},
		{"摔断了腿", 30},
		{"中暑", 20},
		{"被马踢伤", 35},
		{"从车上摔下", 25},
	}

	accident := accidents[rand.Intn(len(accidents))]

	// Apply damage
	for i := range m.Family {
		if m.Family[i].Name == victim.Name {
			m.Family[i].Health -= accident.damage
			if m.Family[i].Health <= 0 {
				m.Family[i].Health = 0
				m.Family[i].Alive = false
				return &GameEvent{
					Type:        EventAccident,
					Title:       "致命意外！",
					Description: fmt.Sprintf("%s %s，不幸身亡！", victim.Name, accident.name),
				}
			}
			break
		}
	}

	return &GameEvent{
		Type:        EventAccident,
		Title:       "意外事故！",
		Description: fmt.Sprintf("%s %s，健康受损！", victim.Name, accident.name),
	}
}

func (m *Model) createThiefEvent() *GameEvent {
	// Random item stolen
	items := []struct {
		name   string
		amount int
	}{
		{"食物", rand.Intn(50) + 20},
		{"衣服", rand.Intn(3) + 1},
		{"弹药", rand.Intn(20) + 10},
		{"零件", rand.Intn(2) + 1},
	}

	item := items[rand.Intn(len(items))]

	switch item.name {
	case "食物":
		stolen := min(item.amount, m.Food)
		m.Food -= stolen
		if stolen > 0 {
			return &GameEvent{
				Type:        EventThief,
				Title:       "盗贼袭击！",
				Description: fmt.Sprintf("盗贼偷走了 %d 磅食物！", stolen),
			}
		}
	case "衣服":
		stolen := min(item.amount, m.Clothing)
		m.Clothing -= stolen
		if stolen > 0 {
			return &GameEvent{
				Type:        EventThief,
				Title:       "盗贼袭击！",
				Description: fmt.Sprintf("盗贼偷走了 %d 套衣服！", stolen),
			}
		}
	case "弹药":
		stolen := min(item.amount, m.Ammunition)
		m.Ammunition -= stolen
		if stolen > 0 {
			return &GameEvent{
				Type:        EventThief,
				Title:       "盗贼袭击！",
				Description: fmt.Sprintf("盗贼偷走了 %d 发弹药！", stolen),
			}
		}
	case "零件":
		stolen := min(item.amount, m.SpareParts)
		m.SpareParts -= stolen
		if stolen > 0 {
			return &GameEvent{
				Type:        EventThief,
				Title:       "盗贼袭击！",
				Description: fmt.Sprintf("盗贼偷走了 %d 个零件！", stolen),
			}
		}
	}

	return nil
}

func (m *Model) createOxDeathEvent() *GameEvent {
	if m.Oxen <= 0 {
		return nil
	}

	// Chance of ox death increases with poor health
	deathChance := 0.3
	if m.OxHealth < 30 {
		deathChance = 0.6
	}

	if rand.Float64() < deathChance {
		m.Oxen--
		m.OxHealth = max(50, m.OxHealth)
		return &GameEvent{
			Type:        EventOxDeath,
			Title:       "牛死亡！",
			Description: "一头牛倒下了！你的牛群现在只剩 " + fmt.Sprintf("%d", m.Oxen) + " 对。",
		}
	}

	// Just injury
	m.OxHealth -= rand.Intn(20) + 10
	if m.OxHealth < 0 {
		m.OxHealth = 0
	}
	return &GameEvent{
		Type:        EventOxDeath,
		Title:       "牛受伤！",
		Description: fmt.Sprintf("一头牛受伤了，牛群健康降至 %d%%", m.OxHealth),
	}
}

func (m *Model) createWagonDamageEvent() *GameEvent {
	damage := rand.Intn(30) + 10
	m.WagonHealth -= damage

	if m.WagonHealth <= 0 {
		m.WagonHealth = 0
		if m.SpareParts > 0 {
			m.SpareParts--
			m.WagonHealth = 50
			return &GameEvent{
				Type:        EventWagonDamage,
				Title:       "车辆损坏！",
				Description: "车辆严重损坏！使用了一个零件进行紧急修理。",
			}
		}
		return &GameEvent{
			Type:        EventWagonDamage,
			Title:       "车辆报废！",
			Description: "车辆彻底损坏，没有零件无法继续前进！",
		}
	}

	return &GameEvent{
		Type:        EventWagonDamage,
		Title:       "车辆损坏！",
		Description: fmt.Sprintf("车辆受损，当前状况 %d%%", m.WagonHealth),
	}
}

func (m *Model) createFoundEvent() *GameEvent {
	finds := []struct {
		name   string
		amount int
	}{
		{"食物", rand.Intn(30) + 10},
		{"衣服", rand.Intn(2) + 1},
		{"弹药", rand.Intn(10) + 5},
	}

	find := finds[rand.Intn(len(finds))]

	switch find.name {
	case "食物":
		m.Food += find.amount
		return &GameEvent{
			Type:        EventFound,
			Title:       "发现补给！",
			Description: fmt.Sprintf("在路上发现了 %d 磅食物！", find.amount),
		}
	case "衣服":
		m.Clothing += find.amount
		return &GameEvent{
			Type:        EventFound,
			Title:       "发现补给！",
			Description: fmt.Sprintf("在路上发现了 %d 套衣服！", find.amount),
		}
	case "弹药":
		m.Ammunition += find.amount
		return &GameEvent{
			Type:        EventFound,
			Title:       "发现补给！",
			Description: fmt.Sprintf("在路上发现了 %d 发弹药！", find.amount),
		}
	}

	return nil
}

// ProcessDiseases handles ongoing disease effects
func (m *Model) ProcessDiseases() []string {
	var messages []string

	for i := range m.Family {
		if !m.Family[i].Alive || m.Family[i].Disease == nil {
			continue
		}

		m.Family[i].DiseaseDay++

		// Daily health loss from disease
		healthLoss := rand.Intn(5) + 3
		m.Family[i].Health -= healthLoss

		// Check if disease is cured (after duration)
		if m.Family[i].DiseaseDay >= m.Family[i].Disease.Duration {
			m.Family[i].Disease = nil
			m.Family[i].DiseaseDay = 0
			messages = append(messages, fmt.Sprintf("%s 的疾病痊愈了！", m.Family[i].Name))
			continue
		}

		// Check for death
		if m.Family[i].Health <= 0 {
			m.Family[i].Health = 0
			m.Family[i].Alive = false
			messages = append(messages, fmt.Sprintf("%s 因 %s 不幸去世！", m.Family[i].Name, m.Family[i].Disease.Name))
			m.Family[i].Disease = nil
		}
	}

	return messages
}

// CountAlive returns number of alive family members
func (m *Model) CountAlive() int {
	count := 0
	for _, member := range m.Family {
		if member.Alive {
			count++
		}
	}
	return count
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
