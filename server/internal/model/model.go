package model

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Role      string `json:"role"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"createdAt"`
}

type Soul struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	BirthInfo        string `json:"birthInfo"`
	DeathInfo        string `json:"deathInfo"`
	Status           string `json:"status"`
	MeritScore       int    `json:"meritScore"`
	KarmaScore       int    `json:"karmaScore"`
	MemoryResidue    int    `json:"memoryResidue"`
	RelationshipRisk string `json:"relationshipRisk"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}
