package spentcalories

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

func parseTraining(data string) (int, string, time.Duration, error) {
	// TODO: реализовать функцию
	// Парсит строки с данными о тренировке в формате "<кол-во шагов>,<тип тренировки>,<длительность>"

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных ожидается 3 части, получено %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("некорректное количество шагов: %v", err)
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("неккоректная длительность: %v", err)
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	// Расчет дистанции в метрах и километрах на основе количества шагов и роста пользователя

	stepLength := stepLengthCoefficient * height // расчет длины шага

	distanceMeters := float64(steps) * stepLength // расчет пройденной дистанции в метрах

	return distanceMeters / mInKm // возврат дистанции в километрах
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	// Расчет средней скорости в км/ч

	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	// Обрабатывает строку с данными о тренировке и возвращает отформатированный отчёт
	// (тип тренировки, длительность, дистанция, скорость, сожженные калории)

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

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		trainingType, duration.Hours(), distance(steps, height), meanSpeed(steps, height, duration), calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	// Проверка корректности входных данных, расчет и возврат количества сожженных калорий при беге

	if steps <= 0 {
		return 0, errors.New("Количество шагов обязан быть больше нуля - проверь данные")
	}
	if weight <= 0 {
		return 0, errors.New("Вес обязан быть больше нуля - проверь данные")
	}
	if height <= 0 {
		return 0, errors.New("Рост обязан быть больше нуля - проверь данные")
	}
	if duration <= 0 {
		return 0, errors.New("Время тренировки обязано быть больше нуля - проверь данные")
	}

	avgSpeed := meanSpeed(steps, height, duration) // Средняя скорость

	if avgSpeed <= 0 {
		return 0, errors.New("Ошибка в входных параметров для расчета средней скорости - проверь данные")
	}

	durationInMinutes := duration.Minutes()

	runCal := (weight * avgSpeed * durationInMinutes) / minInH
	return runCal, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	// Проверка корректности входных данных, расчет и возврат количества сожженных калорий при ходьбе

	if steps <= 0 {
		return 0, errors.New("Количество шагов обязан быть больше нуля - проверь данные")
	}
	if weight <= 0 {
		return 0, errors.New("Вес обязан быть больше нуля - проверь данные")
	}
	if height <= 0 {
		return 0, errors.New("Рост обязан быть больше нуля - проверь данные")
	}
	if duration <= 0 {
		return 0, errors.New("Время тренировки обязано быть больше нуля - проверь данные")
	}

	avgSpeed := meanSpeed(steps, height, duration) // Средняя скорость

	if avgSpeed <= 0 {
		return 0, errors.New("Ошибка в входных параметров для расчета средней скорости - проверь данные")
	}
	durationInMinutes := duration.Minutes()

	calWalk := (weight * avgSpeed * durationInMinutes) / minInH * walkingCaloriesCoefficient
	return calWalk, nil
}
