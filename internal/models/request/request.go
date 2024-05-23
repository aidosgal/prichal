package request

type Request struct {
  ID                int    `json:"id"`
  UserID            int    `json:"user_id"`
  Title             string `json:"title"`
  Description       string `json:"description"`
  Location          string `json:"location"`
  CategoryID        int    `json:"category_id"`
  SubCategoryID     int    `json:"sub_category_id"`
  SpecializationID  int    `json:"specialization_id"`
}
