package dto

type ClassResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RaceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type QuestLevelResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type OptionReq struct {
	Name string `json:"name" validate:"required,max=100"`
}
