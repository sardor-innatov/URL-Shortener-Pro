package service

import (
	"strings"
	"sync"

	"url_shortener_pro/src/services/link_service/click/model"

	"gorm.io/gorm"
)

type ClickService interface {
	Record(linkId int64, ipAddress string, userAgent string)
	Shutdown()
	worker(workerID int)
}

type clickService struct {
	db         *gorm.DB
	clicksChan chan model.Click
	wg         sync.WaitGroup
}

func NewClickService(db *gorm.DB, workerCount int) ClickService {
	service := &clickService{
		db:         db,
		clicksChan: make(chan model.Click, 1000),
	}

	for i := 0; i < workerCount; i++ {
		service.wg.Add(1)
		go service.worker(i)
	}

	return service
}

func (s *clickService) worker(workerID int) {
	defer s.wg.Done()
	println("Worker %d started", workerID)

	for click := range s.clicksChan {

		err := s.db.Create(&click).Error
		if err != nil {
			println("Worker %d error: %v", workerID, err)
		}
	}
	println("Worker %d stopped", workerID)
}

func (s *clickService) Record(linkId int64, ipAddress string, userAgent string) {

	country := getCountry(ipAddress)

	click := model.Click{
		LinkId: linkId,
		IpAddress: ipAddress,
		UserAgent: userAgent,
		Country: &country,
	} 

	select {
	case s.clicksChan <- click:
	default:
		// if chanel is full
		println("Click dropped: channel closed or full")
	}
}

func (s *clickService) Shutdown() {
    println("Stopping ClickService workers...")

    close(s.clicksChan)

    // waiting till all workers DONE
    s.wg.Wait()

    println("ClickService: All workers finished their jobs.")
}

func getCountry(ipAddress string) string {
	// 1. Убираем лишние пробелы, если они есть
	ip := strings.TrimSpace(ipAddress)

	// 2. Проверяем, с чего начинается IP
	if strings.HasPrefix(ip, "192.") {
		return "US"
	}
	if strings.HasPrefix(ip, "172.") {
		return "DE"
	}
	if strings.HasPrefix(ip, "10.") {
		return "FR"
	}

	// 3. Если не подошло ни под одно правило
	return "UNKNOWN"
}