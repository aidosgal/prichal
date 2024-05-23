package subcategory

type SubCategory struct {
  ID          int    `json:"id"`
  Title       string `json:"title"`
  CategoryID  int    `json:"category_id"`
}
