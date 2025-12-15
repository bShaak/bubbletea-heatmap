package bubbleteaheatmap

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	selectedX int
	selectedY int
	calData   []CalDataPoint
	viewData  [][7]viewDataPoint
	Weeks     int
	EndDate   time.Time // The reference date for the heatmap (usually "today")
}

var scaleColors = []string{
	// Light Theme
	// #ebedf0 - Less
	// #9be9a8
	// #40c463
	// #30a14e
	// #216e39 - More

	// Dark Theme
	"#161b22", // Less
	"#0e4429",
	"#006d32",
	"#26a641",
	"#39d353", // - More

}

type CalDataPoint struct {
	Date  time.Time
	Value float64
}

func (m *Model) addCalData(date time.Time, val float64) {
	// Create new cal data point and add to cal data
	newPoint := CalDataPoint{date, val}
	m.calData = append(m.calData, newPoint)
}

func getIndexDate(x int, y int, now time.Time, weeks int) time.Time {
	// compare the x,y to today and subtract
	today := now
	todayX, todayY := getDateIndex(today, now, weeks)

	diffX := todayX - x
	diffY := todayY - y

	diffDays := diffX*7 + diffY

	targetDate := today.AddDate(0, 0, -diffDays)
	return targetDate
}

func weeksAgo(date time.Time, now time.Time) int {
	today := truncateToDate(now)
	thisWeek := today.AddDate(0, 0, -int(today.Weekday())) // Most recent Sunday

	compareDate := truncateToDate(date)
	compareWeek := compareDate.AddDate(0, 0, -int(compareDate.Weekday()))

	result := thisWeek.Sub(compareWeek).Hours() / 24 / 7
	return int(result)
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Local().Year(), t.Local().Month(), t.Local().Day(), 0, 0, 0, 0, t.Local().Location())
}

func getDateIndex(date time.Time, now time.Time, weeks int) (int, int) {
	// Max index - number of weeks ago
	x := (weeks - 1) - weeksAgo(date, now)

	y := int(date.Local().Weekday())

	return x, y
}

func parseCalToView(calData []CalDataPoint, now time.Time, weeks int) [][7]viewDataPoint {
	viewData := make([][7]viewDataPoint, weeks)

	for _, v := range calData {
		x, y := getDateIndex(v.Date, now, weeks)
		// Check if in range
		if x > -1 && y > -1 &&
			x < weeks && y < 7 {
			viewData[x][y].actual += v.Value
		}
	}
	viewData = normalizeViewData(viewData)
	return viewData
}

func normalizeViewData(data [][7]viewDataPoint) [][7]viewDataPoint {
	var min float64
	var max float64

	// Find min/max
	min = data[0][0].actual
	max = data[0][0].actual

	for _, row := range data {
		for _, val := range row {

			if val.actual < min {
				min = val.actual
			}
			if val.actual > max {
				max = val.actual
			}
		}
	}

	// Normalize the data
	for i, row := range data {
		for j, val := range row {
			data[i][j].normalized = (val.actual - min) / (max - min)
		}
	}
	return data
}

type viewDataPoint struct {
	actual     float64
	normalized float64
}

func getScaleColor(value float64) string {
	const numColors = 5
	// Assume it's normalized between 0.0-1.0
	const max = 1.0
	// const min = 0.0

	return scaleColors[int((value/max)*(numColors-1))]
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Create a new model with default settings.
func New(data []CalDataPoint, endDate time.Time, weeks int) Model {
	todayX, todayY := getDateIndex(endDate, endDate, weeks)

	parsedData := parseCalToView(data, endDate, weeks)
	return Model{
		selectedX: todayX,
		selectedY: todayY,
		calData:   data,
		viewData:  parsedData,
		Weeks:     weeks,
		EndDate:   endDate,
		// focus:     false, // TODO
	}
}

// func (m *Model) Focus() tea.Cmd { // TODO
// m.focus = true
// return m.Cursor.Focus()
// }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: ignore if not focused
	// if !m.focus { return m, nil }

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.selectedY > 0 {
				m.selectedY--
			} else if m.selectedY == 0 && m.selectedX > 0 {
				// Scroll to the end of the previous week if on Sunday
				// and not at the beginning of the calendar.
				m.selectedY = 6
				m.selectedX--
			}

		case "down", "j":
			// Don't allow user to scroll beyond today
			if m.selectedY < 6 &&
				(m.selectedX != m.Weeks-1 ||
					m.selectedY < int(m.EndDate.Weekday())) {
				m.selectedY++
			} else if m.selectedY == 6 && m.selectedX != m.Weeks-1 {
				// Scroll to the beginning of next week if on Saturday
				// and not at the end of the calendar.
				m.selectedY = 0
				m.selectedX++
			}
		case "right", "l":
			// Don't allow users to scroll beyond today from the previous column
			if m.selectedX < m.Weeks-2 ||
				(m.selectedX == m.Weeks-2 &&
					m.selectedY <= int(m.EndDate.Weekday())) {
				m.selectedX++
			}
		case "left", "h":
			if m.selectedX > 0 {
				m.selectedX--
			}
		case "enter", " ":
			// Hard coded to add a new entry with value `1.0`
			m.addCalData(
				getIndexDate(m.selectedX, m.selectedY, m.EndDate, m.Weeks),

				1.0)
			m.viewData = parseCalToView(m.calData, m.EndDate, m.Weeks)

		}
	}
	return m, nil
}

func (m Model) View() string {
	// The header

	theTime := getIndexDate(m.selectedX, m.selectedY, m.EndDate, m.Weeks) // time.Now()

	title, _ := glamour.Render(theTime.Format("# Monday, January 02, 2006"), "dark")
	s := title

	selectedDetail := "    Value: " +
		fmt.Sprint(m.viewData[m.selectedX][m.selectedY].actual) +
		" normalized: " +
		fmt.Sprint(m.viewData[m.selectedX][m.selectedY].normalized) +
		"\n\n"

	s += selectedDetail

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	boxStyle := lipgloss.NewStyle().
		PaddingRight(1).
		Foreground(lipgloss.Color(scaleColors[2]))

	boxSelectedStyle := boxStyle.Copy().
		Background(lipgloss.Color("#9999ff")).
		Foreground(lipgloss.Color(scaleColors[0]))

	// Month Labels
	var currMonth time.Month
	s += "  "
	for j := 0; j < m.Weeks; j++ {
		// Check the last day of the week for that column
		jMonth := getIndexDate(j, 6, m.EndDate, m.Weeks).Month()

		if currMonth != jMonth {
			currMonth = jMonth
			s += labelStyle.Render(getIndexDate(j, 6, m.EndDate, m.Weeks).Format("Jan") + " ")

			// Skip the length of the label we just added
			j += 1
		} else {
			s += "  "
		}
	}
	s += "\n"

	for j := 0; j < 7; j++ {
		// Add day of week labels
		switch j {
		case 0:
			s += labelStyle.Render("S ")
		case 1:
			s += labelStyle.Render("M ")
		case 2:
			s += labelStyle.Render("T ")
		case 3:
			s += labelStyle.Render("W ")
		case 4:
			s += labelStyle.Render("T ")
		case 5:
			s += labelStyle.Render("F ")
		case 6:
			s += labelStyle.Render("S ")
		}

		// Add calendar days
		for i := 0; i < m.Weeks; i++ {
			// Selected Item
			if m.selectedX == i && m.selectedY == j {
				s += boxSelectedStyle.Copy().Foreground(
					lipgloss.Color(
						getScaleColor(
							m.viewData[i][j].normalized))).
					Render("■")
			} else if i == m.Weeks-1 &&
				j > int(m.EndDate.Weekday()) {

				// In the future
				s += boxStyle.Render(" ")
			} else {

				// Not Selected Item and not in the future
				s += boxStyle.Copy().
					Foreground(
						lipgloss.Color(
							getScaleColor(
								m.viewData[i][j].normalized))).
					Render("■")
			}
		}
		s += "\n"
	}

	return s
}
