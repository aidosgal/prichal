package group

type Group struct {
  ID           int    `json:"id"`
  Title        string `json:"name"`
  Description  string `json:"description"`
  ChatID       int    `json:"chat_id"`
  ImageURL     string `json:"image_url"`
  TarifID      int    `json:"tarif_id"`
  CreaterID    int    `json:"creater_id"`
}
