package entry

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	TaskQuantity int    `json:"TaskQuantity"`
	TwoCapKey    string `json:"2CapKey"`
}

func LoadConfig() (*Config, error) {

	file, err := os.Open("data/config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	if config.TaskQuantity <= 0 {
		return nil, fmt.Errorf("task quantity must be greater than 0")
	}

	if len(config.TwoCapKey) == 0 {
		return nil, fmt.Errorf("capSolver API key not provided")
	}

	return &config, nil
}

func LoadTextFile(filepath string) ([]string, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil

}

func AppendToFile(filepath, str string) error {

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the string to the file
	if _, err := file.WriteString(str + "\n"); err != nil {
		return err
	}

	return nil
}

func (s *Session) testProxy() (int, error) {
	res, err := s.Client.R().Get("https://ipecho.net/plain")
	// res.Close
	if err != nil {

		s.Log.Errorf("error with proxy %s, retrying...", err.Error())
		time.Sleep(3 * time.Second)
		s.setRandomProxy()
		s.testProxy()

	}
	return res.StatusCode(), nil

}
