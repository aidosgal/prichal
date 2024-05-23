package subscribe_group

type SubscribeGroup struct {
  ID           int    `json:"id"`
  UserID       int    `json:"user_id"`
  GroupID      int    `json:"group_id"`
  Status       string `json:"status"`
  CreatedAt    string `json:"created_at"`
}
