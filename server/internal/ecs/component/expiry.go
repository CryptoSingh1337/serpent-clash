package component

type Expiry struct {
	TicksRemaining uint32
}

func NewExpiryComponent(ticks uint32) Expiry {
	return Expiry{
		ticks,
	}
}
