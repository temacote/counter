package counter

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"sber_cloud/tw/definition/config"
	"sber_cloud/tw/redis"
)

type (
	Counter interface {
		SaveDataFromIMDBToFile() error
		LoadFromFileToIMDB() error
		AddToHistory(request *http.Request) error
	}

	counter struct {
		rds     redis.Redis
		conf    config.Config
		history []*RequestHistory
	}

	RequestHistory struct {
		RequestTime time.Time
		Method      string
		RequestURI  string
	}
)

func NewCounter(conf config.Config, rds redis.Redis) Counter {
	return &counter{
		rds:  rds,
		conf: conf,
	}
}

// SaveToDisk сохранение истории на диск
func (c *counter) SaveDataFromIMDBToFile() error {
	values, err := c.rds.GetAllValues()
	if err != nil {
		return err
	}

	f, err := os.Create(c.conf.GetString("storage_file.path"))
	if err != nil {
		return err
	}
	defer f.Close()
	for _, value := range values {
		f.WriteString(string(value.GetJson()) + "\n")
	}

	return nil
}

// LoadFromDisk загрузка истории с диска
func (c *counter) LoadFromFileToIMDB() error {
	file, err := os.Open(c.conf.GetString("storage_file.path"))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c.rds.ParseAndSave(scanner.Bytes(), time.Minute)
	}

	err = scanner.Err()
	return err
}

// AddToHistory добавление запроса в историю
func (c *counter) AddToHistory(request *http.Request) error {
	var requestHistory = &RequestHistory{RequestTime: time.Now(), Method: request.Method, RequestURI: request.RequestURI}
	byteHistory, err := json.Marshal(requestHistory)
	if err != nil {
		return err
	}
	err = c.rds.SaveValue(string(byteHistory), strconv.Itoa(int(time.Now().UnixNano())), 60*time.Second)
	if err != nil {
		return err
	}
	return nil
}
