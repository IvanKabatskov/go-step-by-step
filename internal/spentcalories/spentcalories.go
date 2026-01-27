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
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

var (
	ErrWeiHei    = errors.New("ошибка значений роста и веса")
	ErrDuration  = errors.New("ошибка значения времени")
	ErrWorkOut   = errors.New("неизвестный тип тренировки")
	ErrConvToInt = errors.New("ошибка преобразования в целое число")
	ErrNumSteps  = errors.New("кол-во шагов равно 0")
	ErrConvTime  = errors.New("ошибка преобразования в интервал времени")
	ErrLenSlice  = errors.New("не хватает данных для обработки")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// TODO: реализовать функцию
	list := strings.Split(data, ",") // list содержит слайс подстрок
	if len(list) == 3 {              // Проверка на наличие трех значений в слайсе
		steps, err := strconv.Atoi(list[0]) // Изменение типа string на int
		if err != nil {                     // Проверка наличия ошибки
			return 0, "", 0, ErrConvToInt
		}
		period, err := time.ParseDuration(list[2]) // Изменение типа string на time
		if err != nil {                            // Проверка наличия ошибки при измененеии типа на time
			return 0, "", 0, ErrConvTime
		}
		return steps, list[1], period, nil
	}
	return 0, "", 0, ErrLenSlice // возврат в случае если длина слайса != 3
}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	var way float64
	if height <= 0.0 {
		way = (lenStep * float64(steps)) / float64(mInKm)
	}
	way = (height * stepLengthCoefficient * float64(steps)) / float64(mInKm)
	return way
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	if duration <= 0 {
		return 0.0
	}
	way := distance(steps, height)  // way содержит дистанцию
	speed := way / duration.Hours() // Рассчет скорости
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	steps, workOut, period, err := parseTraining(data)
	if err != nil { // Проверка возврата ошибки из предыдущей функции
		log.Println(err)
		return "", err
	}
	if steps < 1 { // Проверка числа шагов
		log.Println(err)
		return "", ErrNumSteps
	}
	if workOut == "" { // Проверка указания вида тренировки
		log.Println(ErrWorkOut)
		return "", ErrWorkOut
	}
	if period <= 0 { // Проверка значения интервала времени
		log.Println(ErrDuration)
		return "", ErrDuration
	}
	if weight <= 0.0 && height <= 0.0 { // Проверка значений веса и роста
		return "", ErrWeiHei
	}
	var text string
	switch workOut { // вывод на экран информации в зависимости от типа тренировки
	case "Бег":
		way := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, period)
		calories, err := RunningSpentCalories(steps, weight, height, period)
		if err != nil {
			return "", err
		}
		text = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", workOut, period.Hours(), way, meanSpeed, calories)
	case "Ходьба":
		way := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, period)
		calories, err := WalkingSpentCalories(steps, weight, height, period)
		if err != nil {
			return "", err
		}
		text = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", workOut, period.Hours(), way, meanSpeed, calories)
	default:
		return "", ErrWorkOut
	}
	return text, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if steps < 1 { // Проверка кол-ва шагов
		return 0.0, ErrNumSteps
	}
	if weight <= 0.0 && height <= 0.0 { // Проверка значений веса и роста
		return 0.0, ErrWeiHei
	}
	if duration <= 0 {
		return 0.0, ErrDuration // Проверка значения интервал времени
	}
	speed := meanSpeed(steps, height, duration)
	calories := (weight * duration.Minutes() * speed) / minInH // calories содержит расчет калорий
	return calories, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if steps < 1 { // Проверка кол-ва шагов
		return 0.0, ErrNumSteps
	}
	if weight <= 0.0 && height <= 0.0 { // Проверка значений веса и роста
		return 0.0, ErrWeiHei
	}
	if duration <= 0.0 {
		return 0.0, ErrDuration // Проверка значения интервал времени
	}
	speed := meanSpeed(steps, height, duration)
	calories := (weight * duration.Minutes() * speed) / minInH // calories содержит расчет калорий
	calories *= walkingCaloriesCoefficient                     // calories содержит расчет калорий с корректирующим коэф-ом
	return calories, nil
}
