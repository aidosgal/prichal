package user

type User struct {
  ID       int    `json:"id"`
  Username string `json:"username"`
  ChatID   int    `json:"chat_id"`
  Name     string `json:"name"`
  ImageURL string `json:"image_url"`
  TarifID  int    `json:"tarif_id"`
}
