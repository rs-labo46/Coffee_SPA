package usecase

import "context"

type Pinger interface {
	PingContext(ctx context.Context) error // DBの確認する
}

type HealthUC struct {
	p Pinger // 依存
}

func NewHealthUC(p Pinger) HealthUC {
	return HealthUC{p: p} // DI
}

func (u HealthUC) Check(ctx context.Context) error {
	return u.p.PingContext(ctx) // DBにPingする
}
