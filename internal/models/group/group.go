package group

type Group struct {
  ID           int    `json:"id"`
  Title        string `json:"name"`
  Description  string `json:"description"`
  TelegramLink string `json:"telegram_link"`
  ImageURL     string `json:"image_url"`
  TarifID      int    `json:"tarif_id"`
  CreaterID    int    `json:"creater_id"`
}
