package subscribe

type Subscribe struct {
  ID           int    `json:"id"`
  UserID       int    `json:"user_id"`
  SubscribeID  int    `json:"subscribe_id"`
  Status       string `json:"status"`
  CreatedAt    string `json:"created_at"`
}
