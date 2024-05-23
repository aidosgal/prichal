package review

type Review struct {
  ID       int    `json:"id"`
  UserID   int    `json:"user_id"`
  Description string `json:"description"`
  AuthorID int `json:"author_id"`
  Status   string `json:"status"`
}


