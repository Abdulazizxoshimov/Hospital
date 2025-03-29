package time

import (
	"fmt"
	"strings"
	"time"
)

func ParseWorkTime(workTime string) (time.Time, time.Time, error) {
	parts := strings.Split(workTime, "-")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("Noto‘g‘ri vaqt formati: %s", workTime)
	}

	// Bugungi sanani olish
	currentDate := time.Now().Format("2006-01-02")

	// Vaqtlarni to'g'ri formatda yaratish
	startTime, err1 := time.Parse("2006-01-02 15:04", currentDate+" "+parts[0])
	endTime, err2 := time.Parse("2006-01-02 15:04", currentDate+" "+parts[1])

	if err1 != nil || err2 != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("Vaqtni parsing qilishda xatolik")
	}

	return startTime, endTime, nil
}
