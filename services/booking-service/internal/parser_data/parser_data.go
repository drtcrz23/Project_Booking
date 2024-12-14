package parser_data

import (
	"fmt"
	"time"
)

func ParseAndCalculateDays(startDate, endDate string) (float64, error) {
	const dateFormat = "2006-01-02" // Пример формата YYYY-MM-DD

	// Преобразуем строки в объекты типа time.Time
	start, err := time.Parse(dateFormat, startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date: %v", err)
	}

	end, err := time.Parse(dateFormat, endDate)
	if err != nil {
		return 0, fmt.Errorf("invalid end date: %v", err)
	}

	// Проверка, что дата начала не позже даты окончания
	if start.After(end) {
		return 0, fmt.Errorf("start date cannot be after end date")
	}

	// Преобразуем обе даты в начало дня (00:00:00)
	start = start.Truncate(24 * time.Hour)
	end = end.Truncate(24 * time.Hour)

	// Вычисляем разницу в днях
	duration := end.Sub(start)
	days := duration.Hours() / 24

	// Если даты одинаковые, возвращаем 0, так как нет полного промежутка
	if days < 0 {
		return 0, fmt.Errorf("invalid date range: no full days")
	}

	// Возвращаем количество полных суток
	return days, nil
}
