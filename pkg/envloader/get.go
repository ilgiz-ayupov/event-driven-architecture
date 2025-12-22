package envloader

import (
	"log"
	"os"
	"strconv"
	"time"
)

// GetBool возвращает значение переменной окружения как bool.
// Если переменная не установлена или имеет неверный формат, возвращается значение по умолчанию.
func GetBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("переменная окружения %q имеет некорректное значение: ожидался тип bool, получено %q; будет использоваться значение по умолчанию: %v\n", key, value, defaultValue)
			return defaultValue
		}
		return boolValue
	}

	log.Printf("переменная окружения %q не задана; будет использоваться значение по умолчанию: %v\n", key, defaultValue)
	return defaultValue
}

// GetString возвращает значение переменной окружения как string.
// Если переменная не установлена, возвращается значение по умолчанию.
func GetString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Printf("переменная окружения %q не задана; будет использоваться значение по умолчанию: %v\n", key, defaultValue)
	return defaultValue
}

// GetInt возвращает значение переменной окружения как int.
// Если переменная не установлена или не может быть преобразована в int,
// возвращается значение по умолчанию.
func GetInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("переменная окружения %q имеет некорректное значение: ожидался тип int, получено %q; будет использоваться значение по умолчанию: %v\n", key, value, defaultValue)
			return defaultValue
		}
		return intValue
	}

	log.Printf("переменная окружения %q не задана; будет использоваться значение по умолчанию: %v\n", key, defaultValue)
	return defaultValue
}

// GetDuration возвращает значение переменной окружения как time.Duration.
// Если переменная не установлена или не может быть преобразована в time.Duration,
// возвращается значение по умолчанию.
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		durationValue, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("переменная окружения %q имеет некорректное значение: ожидался тип time.Duration, получено %q; будет использоваться значение по умолчанию: %v\n", key, value, defaultValue)
			return defaultValue
		}
		return durationValue
	}

	log.Printf("переменная окружения %q не задана; будет использоваться значение по умолчанию: %v\n", key, defaultValue)
	return defaultValue
}
