package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65 // Средняя длина шага
	mInKm                      = 1000 // Количество метров в километре
	minInH                     = 60   // Количество минут в часе
	stepLengthCoefficient      = 0.45 // Коэффициент для расчета длины шага
	walkingCaloriesCoefficient = 0.5  // Коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных, должно быть 3 компонента: шаги, активность и продолжительность")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("неверное количество шагов")
	}
	activity := strings.TrimSpace(parts[1])

	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, errors.New("неверный формат продолжительности")
	}

	return steps, activity, duration, err
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return float64(steps) * stepLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("неверные входные параметры")
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("неверные входные параметры")
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var speed float64
	var calories float64

	switch activity {
	case "Ходьба":

		speed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":

		speed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, duration.Hours(), distance(steps, height), speed, calories), nil
}
