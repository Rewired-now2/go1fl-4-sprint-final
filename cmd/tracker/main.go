package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/daysteps"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"github.com/golang/protobuf/ptypes/duration"
)

const (
	stepLength = 0.65
	mInKm = 1000 
	lenStep = 0.65
	minInH = 60
	spetLengthCoefficient = 0.45
	walkingCaloriesCoefficient = 0.5
)

func parsePackage(data string) (int, time.Duration, error){
	fragments := strings.Split(data, ",")
		if len(fragments) !==2 {
			return 0, 0, fmt.Errorf("Неверный формат данных")
}
	steps1 := strings.TrimSpace(fragments[0])
	steps, err := strconv.Atoi(steps1)
		if err != nil {
		return 0,0,err
}
		if steps <= 0 {

		return 0,0,fmt.Errorf("Шаги должны быть больше 0")
}

	duration1 := strings.TrimSpace(fragments[1])
	duration, err := time.ParseDuration(duration1)
		if err != nil {
		return 0,0,err

}

		return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
		if err != nil {
			return " "
}
		if steps <= 0{
			return " "
}
	distanceM := float64(steps) * stepLength
	distanceKm := distanceM / float64(mInKm)
	calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return " "
}
		return fmt.Sprintf("Количество шагов: %d.\n Дистанция составила :.2f км.\n Вы сожгли: %.2f ккал.", steps, distanceKm, calories) 
}


func parseTraining (data string) (int, string, time.Duration, error) {
	fragments := strings.Split(data, ",")
		if len(fragments) !=3{
			return 0, " ", 0, fmt.Errorf("Неверный формат данных")
}
	steps1 := strings.TrimSpace(fragments[0])
	steps, err := strconv.Atoi(steps1)
		if err != nil {
			return 0, " ", 0, err
}
	activity := strings.TrimSpace(fragments[1])
	duration1 := strings.TrimSpace(fragments[2])
	duration , err := time.ParseDuration(duration1)
		if err != nil {
			return 0, " ", err
}
		return steps, activity, duration, nil
}

func distance (steps int, height float64) float64{
	stepLengthCalc := height * spetLengthCoefficient
	distanceM := float64(steps) * stepLengthCalc
		return distanceM / float64(mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0{
		return 0
}
	dist := distance(steps, height)
	averSpeed := duration.Hours()
	if averSpeed == 0{
		return 0
}
		return dist/ averSpeed
}

func walkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error){

		if steps <= 0 {
			return 0, fmt.Errorf("Шаги должны быть больше 0")
}
		if weight <= 0{
			return o, fmt.Errorf("Вес должен быть больше 0")

}
		if height <=0 {
			return 0, fmt.Errorf("Рост должен быть больше 0")
}
		if duration <= 0 {
			return 0, fmt.Errorf("Продолжительность должна быть больше 0")
}
		speed := meanSpeed(steps, height, duration)
		if speed == 0 {
			return 0, fmt.Errorf("Не удалось рассчитать скорость")
}
	durationMinutes := duration.Minutes()
	calories := (weight*speed*durationMinutes) / float64(minInH)
	calories *= walkingCaloriesCoefficient
	return calories, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error){
	return  walkingSpentCalories(steps, weight, height, duration)
}

func TrainingInfo(data string, weight, height float64) (string, error){

	steps, activity, duration, err := parseTraining(data)
		if err != nil {
		log.Println(err)
			return " ", err
}
	var calories float64
	var dist, speed float64

	switch 
	strings.ToLower(activity) {
		case "хотьба", "walking" :
			dist = distance(steps, height)
			speed = meanSpeed(steps, height, duration)
			calories, err = walkingSpentCalories(steps, weight, height, duration)
				if err != nil {
					log.Println(err)
						return " ", err
	} 
		case "бег", "running" :
			dist = distance(steps, height)
			speed = meanSpeed(speed, height, duration)
	
			calories, err = RunningSpentCalories(steps, weight, height,duration)
				if err != nil{
					fmt.Errorf("Неизвестный тип тренировки")
					log.Println(err)
					return " ", err
	}

		default:
			err := fmt.Errorf("неизвестный тип тренировки")
			log.Println(err)
				return " ", err
}

	durationHours := duration.Hours()
	result := fmt.Sprintf("Тип тренировки: %s\n" + "Длительность: %.2f ч.\n" + "Дистанция: %.2f км.\n" + "Скорость:  %.2f км/ч\n" + "Сожгли калорий: %.2f"б activity, durationHours, dist, speed, calories )
		return result, nil
}


func main() {
	weight := 84.6
	height := 1.87

	// дневная активность
	input := []string{
		"678,0h50m",
		"792,1h14m",
		"1078,1h30m",
		"7830,2h40m",
		",3456",
		"12:40:00, 3456",
		"something is wrong",
	}

	fmt.Println("Активность в течение дня")

	var (
		dayActionsInfo string
		dayActionsLog  []string
	)

	for _, v := range input {
		dayActionsInfo = daysteps.DayActionInfo(v, weight, height)
		dayActionsLog = append(dayActionsLog, dayActionsInfo)
	}

	for _, v := range dayActionsLog {
		fmt.Println(v)
	}

	// тренировки
	trainings := []string{
		"3456,Ходьба,3h00m",
		"something is wrong",
		"678,Бег,0h5m",
		"1078,Бег,0h10m",
		",3456 Ходьба",
		"7892,Ходьба,3h10m",
		"15392,Бег,0h45m",
	}

	var trainingLog []string

	for _, v := range trainings {
		trainingInfo, err := spentcalories.TrainingInfo(v, weight, height)
		if err != nil {
			log.Printf("не получилось получить информацию о тренировке: %v", err)
			continue
		}
		trainingLog = append(trainingLog, trainingInfo)
	}

	fmt.Println("Журнал тренировок")

	for _, v := range trainingLog {
		fmt.Println(v)
	}
}
