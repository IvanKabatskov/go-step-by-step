package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Деление строки на слайсы строк
	sliceStr := strings.Split(data, ",")
	// Проверка наличия трех слайсов
	if len(sliceStr) != 3 {
		return 0, "", 0, fmt.Errorf("ожидалось 3 элемента, получено %d", len(sliceStr))
	}
	// Преобразование string в int
	steps, err := strconv.Atoi(sliceStr[0])
	// Проверка на ошибку
	if err != nil {
		return 0, "", 0, err
	}
	// Проверка на кол-во шагов
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("некорректное значение %d", steps)
	}
	// Преобразование string в Duration
	duration, err := time.ParseDuration(sliceStr[2])
	// Проверка на ошибку
	if err != nil {
		return 0, "", 0, err
	}
	// Проверка длительности активности
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("некорректное значение %v", duration)
	}
	// Вывод активности
	activity := sliceStr[1]
	// Проверка данных об актинвости
	if activity != "Бег" && activity != "Ходьба" {
		return 0, "", 0, fmt.Errorf("неизвестный тип тренировки")
	}
	// Возврат корректных значений
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	distance := 0.0
	if height > 0 {
		// Расчет дистанции c известным ростом
		distance = (height * stepLengthCoefficient) * float64(steps) / mInKm
	} else {
		// Pасчет дистанции с неизвестным ростом
		distance = lenStep * float64(steps) / mInKm
	}
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверка длительности активности
	if duration <= 0 {
		return 0
	}
	// Вычисление дистанции
	distance := distance(steps, height)
	// Вычисление средней скорости
	meanSpeed := distance / duration.Hours()
	return meanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получение данных из parseTraining
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}
	message := ""
	switch activity {
	case "Бег":
		//Расчет дистанции
		distance := distance(steps, height)
		//Рассчет стредней скорости
		meanSpeed := meanSpeed(steps, height, duration)
		//Рассчет калорий
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		//Создание сообщения вывода
		message = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distance, meanSpeed, calories)
	case "Ходьба":
		//Расчет дистанции
		distance := distance(steps, height)
		//Рассчет стредней скорости
		meanSpeed := meanSpeed(steps, height, duration)
		//Рассчет калорий
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		//Создание сообщения вывода
		message = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distance, meanSpeed, calories)
	}
	return message, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка параметра вес
	if weight <= 0 {
		return 0, fmt.Errorf("некорректное значение %.2f", weight)
	}
	// Проверка параметра рост
	if height <= 0 {
		return 0, fmt.Errorf("некорректное значение %.2f", height)
	}
	// Проверка продолжительности активности
	if duration <= 0 {
		return 0, fmt.Errorf("некорректное значение %v", duration)
	}
	// Проверка кол-ва шагов
	if steps <= 0 {
		return 0, fmt.Errorf("некорректное значение %d", steps)
	}
	// Рассчет средней скорости
	meanSpeed := meanSpeed(steps, height, duration)
	// Рассчет калорий
	calories := (weight * meanSpeed * duration.Minutes()) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка параметра вес
	if weight <= 0 {
		return 0, fmt.Errorf("некорректное значение %.2f", weight)
	}
	// Проверка параметра рост
	if height <= 0 {
		return 0, fmt.Errorf("некорректное значение %.2f", height)
	}
	// Рассчет средней скорости
	meanSpeed := meanSpeed(steps, height, duration)
	// Рассчет калорий
	calories := ((weight * meanSpeed * duration.Minutes()) / minInH) * walkingCaloriesCoefficient
	return calories, nil
}
