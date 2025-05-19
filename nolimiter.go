package main

type NoLimiter struct {
}

func (nl *NoLimiter) Run() (*LimiterResult, error) {

	return nil, nil
}
