package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

var (
	ErrConvToInt = errors.New("ошибка преобразования в целое число")
	ErrNumSteps  = errors.New("кол-во шагов равно 0")
	ErrConvTime  = errors.New("ошибка преобразования в интервал времени")
	ErrLenSlice  = errors.New("не хватает данных для обработки")
)

func parsePackage(data string) (int, time.Duration, error) {
	// TODO: реализовать функцию
	list := strings.Split(data, ",") // list содержит слайс подстрок
	if len(list) == 2 {              // Проверка на наличие двух значений в слайсе
		steps, err := strconv.Atoi(list[0]) // Изменение типа string на int
		if err != nil {                     // Проверка наличия ошибки при изменении типа на int
			return 0, 0, ErrConvToInt
		}
		if steps < 1 { // Проверка числа шагов
			return 0, 0, ErrNumSteps
		}
		period, err := time.ParseDuration(list[1]) // Изменение типа string на time
		if err != nil {                            // Проверка наличия ошибки при измененеии типа на time
			return 0, 0, ErrConvTime
		}
		return steps, period, nil
	}
	return 0, 0, ErrLenSlice
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	steps, period, err := parsePackage(data) // steps - шаги, period - время, err - ошибка
	if err != nil {                          // проверка на наличие возврата ошибки из parsePackage
		fmt.Println(err)
		return " "
	}
	if steps < 1 { // проверка числа шагов
		return " "
	}
	wayLenghtM := float64(steps) * stepLength                                                // Дистанция в м
	wayLenghtKm := wayLenghtM / float64(mInKm)                                               // Дистанция в км
	caloriesOnWalk, err := spentcalories.WalkingSpentCalories(steps, weight, height, period) // Кол-во калорий на прогулке
	if err != nil {
		return " "
	}
	text := fmt.Sprintf("Количество шагов: %d.\n Дистанция составила: %.2f км.\n Вы сожгли: %.2f  ккал.\n", steps, wayLenghtKm, caloriesOnWalk)
	return text
}
