package component

type PlayerInfo struct {
	ID       string
	Username string
}

func NewPlayerInfoComponent(id string, username string) PlayerInfo {
	return PlayerInfo{
		id, username,
	}
}
