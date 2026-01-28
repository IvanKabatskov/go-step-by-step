package daysteps

import (
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

func parsePackage(data string) (int, time.Duration, error) {
	// Деление строки на слайсы строк
	sliceStr := strings.Split(data, ",")
	// Проверка наличия двух слайсов
	if len(sliceStr) != 2 {
		return 0, 0, fmt.Errorf("ожидалось 2 элемента, получено %d", len(sliceStr))
	}
	// Преобразование string в int
	steps, err := strconv.Atoi(sliceStr[0])
	// Проверка на ошибку
	if err != nil {
		return 0, 0, err
	}
	// Проверка на кол-во шагов
	if steps <= 0 {
		return 0, 0, fmt.Errorf("некорректное значение шагов %d", steps)
	}
	// Преобразование string в Duration
	duration, err := time.ParseDuration(sliceStr[1])
	// Проверка на ошибку
	if err != nil {
		return 0, 0, err
	}
	// Проверка длительности активности
	if duration <= 0 {
		return 0, 0, fmt.Errorf("некорректное значение интервала времени %v", duration)
	}
	// Возврат корректных значений
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получение данных их parsePackage
	steps, duration, err := parsePackage(data)
	// Вывод на экран ошибки
	if err != nil {
		fmt.Println("Произошла ошибка:", err)
		return ""
	}
	//Вычисление дистанции в метрах
	distance := float64(steps) * stepLength
	// Вычисление дистанции в км
	distance /= mInKm
	// Вычисление кол-ва калорий
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("Произошла ошибка:", err)
		return ""
	}
	//Создание сообщения вывода
	message := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
	return message
}
