package entry

import (
	"fmt"
	"math/rand"
	"time"

	"secondlap.go/discord"
)

func (s *Session) Ignite() (bool, *discord.DiscordData, error) {

	s.Log.Info("Starting Task...")

	if err := s.setRandomProxy(); err != nil {
		return false, nil, fmt.Errorf("failed to set proxy: %w", err)
	}

	if err := s.getHomepage(); err != nil {
		return false, nil, fmt.Errorf("failed to get homepage: %w", err)
	}

	if err := s.createCaptchaSolver(); err != nil {
		return false, nil, fmt.Errorf("failed to create captcha solver: %w", err)
	}

	if err := s.submitEntry(); err != nil {
		return false, nil, fmt.Errorf("failed submitting entry: %w", err)
	}

	return true, &s.Webhook, nil

}

func (s *Session) ProxyTest() (bool, error) {

	if err := s.setIndexProxy(s.ID - 1); err != nil {
		return false, fmt.Errorf("failed to set proxy: %w", err)
	}
	t := rand.Intn(10)
	s.Log.Infof("delaying for %d seconds", t)

	time.Sleep(time.Duration(t) * time.Second)

	if s, err := s.testProxy(); err != nil {
		return false, fmt.Errorf("proxy failed test: %w, %d", err, s)
	}

	return true, nil
}
