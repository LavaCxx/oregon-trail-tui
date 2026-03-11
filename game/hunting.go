package game

import (
	"fmt"
	"math/rand"
)

// StartHunting initializes a hunting session
func (m *Model) StartHunting() bool {
	if m.Ammunition < 1 {
		m.AddMessage("没有弹药，无法狩猎！")
		return false
	}

	m.HuntingAnimals = make([]Animal, 0)
	m.HuntingShots = min(m.Ammunition, 20) // Max 20 shots
	m.HuntingFood = 0
	m.HuntingDay = 0

	// Generate random animals
	numAnimals := rand.Intn(5) + 3
	for i := 0; i < numAnimals; i++ {
		animal := HuntingAnimals[rand.Intn(len(HuntingAnimals))]
		m.HuntingAnimals = append(m.HuntingAnimals, animal)
	}

	m.PrevState = m.State
	m.State = StateHunting
	return true
}

// Shoot attempts to shoot an animal
func (m *Model) Shoot(animalIndex int) string {
	if m.HuntingShots <= 0 {
		return "没有弹药了！"
	}

	if animalIndex < 0 || animalIndex >= len(m.HuntingAnimals) {
		return "无效的目标！"
	}

	animal := m.HuntingAnimals[animalIndex]
	m.HuntingShots--
	m.Ammunition--

	// Calculate hit chance based on animal speed
	hitChance := 100 - animal.Speed
	if hitChance < 20 {
		hitChance = 20
	}

	// Random roll
	roll := rand.Intn(100)
	if roll < hitChance {
		// Hit!
		food := animal.Weight
		// Can only carry back 100-200 lbs
		maxCarry := 150
		if m.HuntingFood+food > maxCarry {
			food = maxCarry - m.HuntingFood
		}
		m.HuntingFood += food

		// Remove the animal
		m.HuntingAnimals = append(m.HuntingAnimals[:animalIndex], m.HuntingAnimals[animalIndex+1:]...)

		return fmt.Sprintf("命中！获得 %d 磅食物！", food)
	}

	// Miss
	return fmt.Sprintf("打偏了！%s 跑掉了。", animal.Name)
}

// EndHunting finishes the hunting session
func (m *Model) EndHunting() {
	m.Food += m.HuntingFood
	m.State = StateHuntingResult
}

// GetHuntingStatus returns current hunting status
func (m *Model) GetHuntingStatus() string {
	return fmt.Sprintf("剩余弹药: %d | 已获得: %d 磅食物", m.HuntingShots, m.HuntingFood)
}

// GenerateNewAnimals spawns new animals for the next round
func (m *Model) GenerateNewAnimals() {
	if m.HuntingDay >= 3 {
		return // Max 3 hunting rounds per session
	}

	m.HuntingDay++
	m.HuntingAnimals = make([]Animal, 0)

	// Fewer animals each round
	numAnimals := rand.Intn(4-m.HuntingDay) + 2
	if numAnimals < 1 {
		numAnimals = 1
	}

	for i := 0; i < numAnimals; i++ {
		animal := HuntingAnimals[rand.Intn(len(HuntingAnimals))]
		m.HuntingAnimals = append(m.HuntingAnimals, animal)
	}
}

// HuntingScene generates ASCII art for hunting
func (m *Model) HuntingScene() string {
	scene := `
    ╔══════════════════════════════════════════════════════════════╗
    ║                      🏹 狩猎场                                ║
    ╠══════════════════════════════════════════════════════════════╣
    ║                                                              ║
`

	// Display animals with positions
	for i, animal := range m.HuntingAnimals {
		pos := rand.Intn(40) + 5
		spaces := ""
		for j := 0; j < pos; j++ {
			spaces += " "
		}
		scene += fmt.Sprintf("    ║ [%d] %s%s%-6s║\n", i+1, spaces, animal.Symbol, animal.Name)
	}

	scene += `    ║                                                              ║
    ╚══════════════════════════════════════════════════════════════╝
`
	return scene
}
