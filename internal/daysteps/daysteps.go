package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go1fl-4-sprint-final/internal/spentcalories"
)

const (
	stepLength = 0.65 // длина шага в метрах
	mInKm      = 1000 // метров в километре
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || steps <= 0 {
		return 0, 0, fmt.Errorf("неверное количество шагов")
	}

	durationStr := strings.TrimSpace(parts[1])
	if durationStr == "" {
		return 0, 0, fmt.Errorf("неверная продолжительность")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат продолжительности: %v", err)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка:", err)
		return ""
	}

	if steps <= 0 || duration <= 0 {
		log.Println("Ошибка: шаги или продолжительность должны быть больше 0")
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / float64(mInKm)
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Ошибка при вычислении калорий:", err)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)
	return result
}
