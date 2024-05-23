package tarif

type Tarif struct {
  ID          int    `json:"id"`
  Title       string `json:"title"`
  Description string `json:"description"`
  Price       int    `json:"price"`
  ImageURL    string `json:"image_url"`
}
