package spentcalories

// Sprint 4: Final Version
import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	// lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining парсит строку вида "7500,Бег,1h15m" (шаги/тип тренировки/ время)
// Возвращает шаги, тип тренировки, длительность и ошибку
// Все некорректные события (ноль/отрицательные значения) — возвращаем ошибку сразу на этапе парсинга
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных ожидается 3 части, получено %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("некорректное количество шагов: %v", err)
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("неккоректная длительность: %v", err)
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("длительность должна быть больше нуля")
	}

	return steps, parts[1], duration, nil
}

// distance возвращает дистанцию в километрах по росту и количеству шагов
func distance(steps int, height float64) float64 {
	stepLength := stepLengthCoefficient * height
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

// meanSpeed - средняя скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

// TrainingInfo — функция формирует отчёт о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	var calcErr error

	switch trainingType {
	case "Бег":
		calories, calcErr = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, calcErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if calcErr != nil {
		log.Println(calcErr)
		return "", calcErr
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainingType, duration.Hours(), distance(steps, height), meanSpeed(steps, height, duration), calories), nil
}

// RunningSpentCalories — расчёт калорий при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов обязан быть больше нуля - проверь данные")
	}
	if weight <= 0 {
		return 0, errors.New("вес обязан быть больше нуля - проверь данные")
	}
	if height <= 0 {
		return 0, errors.New("рост обязан быть больше нуля - проверь данные")
	}
	if duration <= 0 {
		return 0, errors.New("время тренировки обязано быть больше нуля - проверь данные")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	if avgSpeed <= 0 {
		return 0, errors.New("ошибка в входных параметров для расчета средней скорости - проверь данные")
	}

	durationInMinutes := duration.Minutes()
	runCal := (weight * avgSpeed * durationInMinutes) / minInH
	return runCal, nil
}

// WalkingSpentCalories — то же самое, что и бег, но * 0.5
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов обязан быть больше нуля - проверь данные")
	}
	if weight <= 0 {
		return 0, errors.New("вес обязан быть больше нуля - проверь данные")
	}
	if height <= 0 {
		return 0, errors.New("рост обязан быть больше нуля - проверь данные")
	}
	if duration <= 0 {
		return 0, errors.New("время тренировки обязано быть больше нуля - проверь данные")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	if avgSpeed <= 0 {
		return 0, errors.New("ошибка в входных параметров для расчета средней скорости - проверь данные")
	}
	durationInMinutes := duration.Minutes()

	calWalk := (weight * avgSpeed * durationInMinutes) / minInH * walkingCaloriesCoefficient
	return calWalk, nil
}
