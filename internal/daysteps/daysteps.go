package daysteps

// Sprint 4: Final Version
import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65 // Длина одного шага в метрах
	mInKm      = 1000 // Количество метров в одном километре
)

// parsePackage парсит строку вида "6000,1h30m" (шаги и часы с минутами)
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("Неверный формат данных")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.New("неверный формат количества шагов")
	}
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, errors.New("Неверный формат времени")
	}
	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть больше нуля")
	}

	return steps, duration, nil
}

// Возвращает отчёт о дневной активности (шаги + дистанция + калории)
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err) // логируем ошибку — требование тестов
		return ""
	}

	distance := float64(steps) * stepLength / mInKm // Дистанция с фиксированной длиной шага

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distance, calories)
}
