package game

import (
	"time"
)

// Profession represents player's starting profession
type Profession int

const (
	Banker Profession = iota
	Carpenter
	Farmer
)

func (p Profession) String() string {
	return []string{"银行家", "木匠", "农夫"}[p]
}

func (p Profession) StartingMoney() int {
	return []int{1600, 800, 400}[p]
}

// Weather represents current weather condition
type Weather int

const (
	WeatherClear Weather = iota
	WeatherRain
	WeatherSnow
	WeatherBlizzard
	WeatherStorm
	WeatherFog
	WeatherHot
)

func (w Weather) String() string {
	names := []string{"☀️ 晴朗", "🌧️ 下雨", "❄️ 下雪", "🌨️ 暴风雪", "⛈️ 暴风雨", "🌫️ 大雾", "🥵 酷热"}
	if int(w) < len(names) {
		return names[w]
	}
	return names[0]
}

// Season represents current season
type Season int

const (
	Spring Season = iota
	Summer
	Fall
	Winter
)

func (s Season) String() string {
	return []string{"🌸 春季", "☀️ 夏季", "🍂 秋季", "❄️ 冬季"}[s]
}

// Disease represents possible diseases
type Disease struct {
	Name        string
	Description string
	Fatality    float64 // 0.0 to 1.0
	Duration    int     // days
}

// FamilyMember represents a family member
type FamilyMember struct {
	Name       string
	Health     int // 0-100
	Alive      bool
	Disease    *Disease
	DiseaseDay int // days since infected
}

// RiverCrossingMethod represents how to cross a river
type RiverCrossingMethod int

const (
	CrossFerry RiverCrossingMethod = iota
	CrossSwim
	CrossWait
	CrossGuide
)

// Landmark represents a milestone location
type Landmark struct {
	Name        string
	Distance    int // miles from start
	HasRiver    bool
	RiverDepth  int // feet (for river crossings)
	RiverWidth  int // feet
	HasStore    bool
	IsFort      bool
	Description string
}

// GameState represents the current game state
type GameState int

const (
	StateTitle GameState = iota
	StateProfession
	StateNaming
	StateShopping
	StateTraveling
	StateEvent
	StateRiver
	StateHunting
	StateHuntingResult
	StateDeath
	StateVictory
	StateTombstone
)

// EventType for random events
type EventType int

const (
	EventDisease EventType = iota
	EventAccident
	EventWeather
	EventThief
	EventOxDeath
	EventWagonDamage
	EventFound
	EventNone
)

// GameEvent represents a random event
type GameEvent struct {
	Type        EventType
	Title       string
	Description string
	Choices     []EventChoice
}

// EventChoice represents a choice in an event
type EventChoice struct {
	Text        string
	Consequence func(*Model) string
}

// Animal for hunting
type Animal struct {
	Name   string
	Weight int   // pounds of food
	Speed  int   // affects hit chance
	Value  int   // points
	Symbol string // ASCII representation
}

// ShopItem represents an item in the store
type ShopItem struct {
	Name        string
	Price       int
	Unit        string
	Description string
}

// Model is the main game model
type Model struct {
	// Game state
	State       GameState
	PrevState   GameState // for returning from sub-screens

	// Player info
	Profession   Profession
	PlayerName   string
	Family       []FamilyMember
	Money        int

	// Supplies
	Oxen        int // yoke (pairs)
	Food        int // pounds
	Clothing    int // sets
	Ammunition  int // boxes
	SpareParts  int // wheels, axles, tongues

	// Journey
	MilesTraveled    int
	MilesToNext      int
	CurrentLandmark  int
	Date             time.Time // game date starting March 1, 1848
	Day              int       // day of journey

	// Conditions
	Weather       Weather
	Season        Season
	Temperature   int // -20 to 100 (Fahrenheit)
	Rations       int // 1=filling, 2=meager, 3=bare bones
	Pace          int // 1=steady, 2=strenuous, 3=grueling

	// Wagon
	WagonHealth   int // 0-100
	OxHealth      int // 0-100

	// Current event
	CurrentEvent  *GameEvent
	EventMessage  string

	// River crossing
	RiverChoice   RiverCrossingMethod
	RiverResult   string

	// Hunting
	HuntingAnimals    []Animal
	HuntingShots      int
	HuntingFood       int
	HuntingDay        int

	// UI
	Width         int
	Height        int
	InputBuffer   string
	SelectedIndex int
	ScrollOffset  int

	// Messages
	Messages      []string
	MessageTime   time.Time

	// Landmarks
	Landmarks []Landmark
}

// NewModel creates a new game model
func NewModel() Model {
	m := Model{
		State:         StateTitle,
		Family:        make([]FamilyMember, 0),
		Date:          time.Date(1848, 3, 1, 0, 0, 0, 0, time.UTC),
		Day:           1,
		Weather:       WeatherClear,
		Season:        Spring,
		Rations:       2, // meager
		Pace:          1, // steady
		WagonHealth:   100,
		OxHealth:      100,
		Messages:      make([]string, 0),
		Landmarks:     initLandmarks(),
	}

	return m
}

// Diseases
var Diseases = map[string]Disease{
	"dysentery": {
		Name:        "痢疾",
		Description: "严重的肠道感染",
		Fatality:    0.3,
		Duration:    7,
	},
	"typhoid": {
		Name:        "伤寒",
		Description: "由细菌引起的严重疾病",
		Fatality:    0.4,
		Duration:    14,
	},
	"cholera": {
		Name:        "霍乱",
		Description: "致命的水传播疾病",
		Fatality:    0.6,
		Duration:    5,
	},
	"measles": {
		Name:        "麻疹",
		Description: "传染性病毒感染",
		Fatality:    0.15,
		Duration:    10,
	},
	"fever": {
		Name:        "发烧",
		Description: "高烧不退",
		Fatality:    0.1,
		Duration:    5,
	},
}

// Animals for hunting
var HuntingAnimals = []Animal{
	{Name: "野牛", Weight: 500, Speed: 30, Value: 100, Symbol: "🦬"},
	{Name: "鹿", Weight: 150, Speed: 50, Value: 50, Symbol: "🦌"},
	{Name: "熊", Weight: 400, Speed: 40, Value: 80, Symbol: "🐻"},
	{Name: "兔子", Weight: 5, Speed: 70, Value: 10, Symbol: "🐰"},
	{Name: "松鼠", Weight: 2, Speed: 80, Value: 5, Symbol: "🐿️"},
}

// Shop items
var ShopItems = []ShopItem{
	{Name: "牛", Price: 40, Unit: "对", Description: "拉车必需，建议购买3-6对"},
	{Name: "食物", Price: 20, Unit: "份(20磅)", Description: "每人每天消耗2-5磅"},
	{Name: "衣服", Price: 10, Unit: "套", Description: "御寒防病，每人至少2套"},
	{Name: "弹药", Price: 2, Unit: "盒(20发)", Description: "狩猎用，每盒可射击20次"},
	{Name: "零件", Price: 10, Unit: "个", Description: "修车用，建议备3-5个"},
}

func initLandmarks() []Landmark {
	return []Landmark{
		{Name: "独立城", Distance: 0, Description: "旅程的起点，密苏里州的边境小镇"},
		{Name: "堪萨斯河", Distance: 102, HasRiver: true, RiverDepth: 6, RiverWidth: 200, Description: "第一条大河，需要渡河"},
		{Name: "蓝河", Distance: 185, HasRiver: true, RiverDepth: 4, RiverWidth: 150, Description: "较浅的河流"},
		{Name: "福特要塞", Distance: 320, HasStore: true, IsFort: true, Description: "军事要塞，可以补给"},
		{Name: "拉勒米要塞", Distance: 550, HasStore: true, IsFort: true, Description: "重要的补给站"},
		{Name: "独立岩", Distance: 730, Description: "著名的地标岩石"},
		{Name: "南关", Distance: 830, Description: "穿越落基山脉的关口"},
		{Name: "福特布里杰", Distance: 950, HasStore: true, IsFort: true, Description: "最后的补给站"},
		{Name: "蛇河", Distance: 1100, HasRiver: true, RiverDepth: 8, RiverWidth: 300, Description: "危险的河流"},
		{Name: "蓝山", Distance: 1350, Description: "美丽的山脉"},
		{Name: "达尔斯", Distance: 1550, HasRiver: true, RiverDepth: 10, RiverWidth: 400, Description: "哥伦比亚河峡谷"},
		{Name: "俄勒冈城", Distance: 2000, Description: "旅程的终点！"},
	}
}
