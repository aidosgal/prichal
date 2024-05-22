package postgre

import (
  "github.com/jackc/pgx/v4/pgxpool"
  "context"
  "fmt"
  "github.com/aidosgal/prichal/internal/models/user"
)

type Postgre struct {
  conn *pgxpool.Pool
}

func New() (*Postgre, error) {
  const op = "storage.postgre.New"
  conn, err := pgxpool.Connect(context.Background(), "postgres://postgres:aidos2004@localhost:5432/prichal")

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  createUserTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) NOT NULL,
        chat_id INT NOT NULL,
        name VARCHAR(50) NOT NULL,
        image_url VARCHAR(255) NOT NULL,
        tarif_id INT DEFAULT 1,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    ); 
  `

  _, err = conn.Exec(context.Background(), createUserTable)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }
  return &Postgre{conn: conn}, nil
}

func (p *Postgre) CreateUser(username string, chatId int, name string, imageUrl string) error {
  const op = "storage.postgre.CreateUser"
  query := `
    INSERT INTO users (username, chat_id, name, image_url)
    VALUES ($1, $2, $3, $4)
  `

  _, err := p.conn.Exec(context.Background(), query, username, chatId, name, imageUrl)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) GetUsers() ([]user.User, error) {
  const op = "storage.postgre.GetUsers"
  query := `
    SELECT id, username, chat_id, name, image_url, tarif_id
    FROM users
  `

  rows, err := p.conn.Query(context.Background(), query)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var users []user.User

  for rows.Next() {
    var user user.User
    err := rows.Scan(&user.ID, &user.Username, &user.ChatID, &user.Name, &user.ImageURL, &user.TarifID)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    users = append(users, user)
  }

  return users, nil
}

func (p *Postgre) GetUserByChatID(chatID int) (user.User, error) {
  const op = "storage.postgre.GetUserByChatID"
  query := `
    SELECT id, username, chat_id, name, image_url, tarif_id
    FROM users
    WHERE chat_id = $1
  `

  var user user.User
  err := p.conn.QueryRow(context.Background(), query, chatID).Scan(&user.ID, &user.Username, &user.ChatID, &user.Name, &user.ImageURL, &user.TarifID)

  if err != nil {
    return user, fmt.Errorf("%s: %w", op, err)
  }

  return user, nil
}

func (p *Postgre) UpdateUserTarif(chatID int, tarifID int) error {
  const op = "storage.postgre.UpdateUserTarif"
  query := `
    UPDATE users
    SET tarif_id = $1
    WHERE chat_id = $2
  `

  _, err := p.conn.Exec(context.Background(), query, tarifID, chatID)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) UpdateUser(chatid int, name string, image_url string) error {
  const op = "storage.postgre.UpdateUser"

  query := `
    UPDATE users
    SET name = $1, image_url = $2
    WHERE chat_id = $3
  `
  _, err := p.conn.Exec(context.Background(), query, name, image_url, chatid)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

// TODO: Add DeleteUser method
