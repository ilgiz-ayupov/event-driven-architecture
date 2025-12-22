package envloader

import (
	"log"
	"os"
	"strconv"
	"time"
)

// MustGetBool возвращает значение переменной окружения как bool.
// Если переменная не установлена или имеет неверный формат, программа завершается с ошибкой.
func MustGetBool(key string) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("переменная окружения %q имеет некорректное значение: ожидался тип bool, получено %q\n", key, value)
		}
		return boolValue
	}

	log.Fatalf("переменная окружения %q не задана", key)
	return false
}

// MustGetString возвращает значение переменной окружения как string.
// Если переменная не установлена, программа завершается с ошибкой.
func MustGetString(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Fatalf("переменная окружения %q не задана", key)
	return ""
}

// MustGetInt возвращает значение переменной окружения как int.
// Если переменная не установлена или не может быть преобразована в int,
// программа завершается с ошибкой.
func MustGetInt(key string) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("переменная окружения %q имеет некорректное значение: ожидался тип int, получено %q\n", key, value)
		}
		return intValue
	}

	log.Fatalf("переменная окружения %q не задана", key)
	return 0
}

// MustGetDuration возвращает значение переменной окружения как time.Duration.
// Если переменная не установлена или не может быть преобразована в time.Duration,
// программа завершается с ошибкой.
func MustGetDuration(key string) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		durationValue, err := time.ParseDuration(value)
		if err != nil {
			log.Fatalf("переменная окружения %q имеет некорректное значение: ожидался тип time.Duration, получено %q\n", key, value)
		}
		return durationValue
	}

	log.Fatalf("переменная окружения %q не задана", key)
	return 0
}
