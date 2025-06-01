package component

type Expiry struct {
	TicksRemaining uint
}

func NewExpiryComponent(ticks uint) Expiry {
	return Expiry{
		ticks,
	}
}
