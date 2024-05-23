package activity

type Activity struct {
  ID                int    `json:"id"`
  Title             string `json:"name"`
  Description       string `json:"description"`
  ImageURL          string `json:"image_url"`
  UserID            int    `json:"user_id"`
  CategoryID        int    `json:"category_id"`
  SubCategoryID     int    `json:"sub_category_id"`
  SpecializationID  int    `json:"specialization_id"`
  Location          string `json:"location"`
}
