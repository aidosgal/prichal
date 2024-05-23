package postgre

import (
  "github.com/jackc/pgx/v4/pgxpool"
  "context"
  "fmt"
  "github.com/aidosgal/prichal/internal/models/user"
  "github.com/aidosgal/prichal/internal/models/tarif"
  "github.com/aidosgal/prichal/internal/models/group"
  "github.com/aidosgal/prichal/internal/models/activity"
  "github.com/aidosgal/prichal/internal/models/request"
  "github.com/aidosgal/prichal/internal/models/category"
  "github.com/aidosgal/prichal/internal/models/subcategory"
  "github.com/aidosgal/prichal/internal/models/specialization"
  "github.com/aidosgal/prichal/internal/models/review"
  "github.com/aidosgal/prichal/internal/models/raport"
)

type Postgre struct {
  conn *pgxpool.Pool
}

func (p *Postgre) Conn() *pgxpool.Pool {
	return p.conn
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
        onboarding BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    ); 
  `
  
  createTarifTable := `
    CREATE TABLE IF NOT EXISTS tarifs (
        id SERIAL PRIMARY KEY,
        title VARCHAR(50) NOT NULL,
        description TEXT NOT NULL,
        price INT NOT NULL,
        image_url VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
  `

  createSubscribeTable := `
    CREATE TABLE IF NOT EXISTS subscribes (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        subscriber_id INT NOT NULL,
        status VARCHAR(50) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
  `

  createGroupTable := `
    CREATE TABLE IF NOT EXISTS groups (
        id SERIAL PRIMARY KEY,
        title VARCHAR(50) NOT NULL,
        description TEXT NOT NULL,
        creator_id INT NOT NULL,
        image_url VARCHAR(255) NOT NULL,
        telegram_link VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
  `
  
  createSubscribeGroupTable := `
    CREATE TABLE IF NOT EXISTS subscribe_groups (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        group_id INT NOT NULL,
        status VARCHAR(50) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
  `

  createActivityTable := `
    CREATE TABLE IF NOT EXISTS activities (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        title VARCHAR(50) NOT NULL,
        description TEXT NOT NULL,
        location VARCHAR(50) NOT NULL,
        category_id INT NOT NULL,
        image_url VARCHAR(255) NOT NULL,
        subcategory_id INT NOT NULL,
        specialization_id INT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
  `

  createRequestTable := `
    CREATE TABLE IF NOT EXISTS requests (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        title VARCHAR(50) NOT NULL,
        description TEXT NOT NULL,
        location VARCHAR(50) NOT NULL,
        category_id INT NOT NULL,
        subcategory_id INT NOT NULL,
        specialization_id INT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
  `

  createCategoryTable := `
    CREATE TABLE IF NOT EXISTS categories (
        id SERIAL PRIMARY KEY,
        title VARCHAR(50) NOT NULL
    );
  `

  createSubcategoryTable := `
    CREATE TABLE IF NOT EXISTS subcategories (
        id SERIAL PRIMARY KEY,
        title VARCHAR(50) NOT NULL,
        category_id INT NOT NULL
    );
  `

  createSpecializationTable := `
    CREATE TABLE IF NOT EXISTS specializations (
        id SERIAL PRIMARY KEY,
        title VARCHAR(50) NOT NULL,
        subcategory_id INT NOT NULL
    );
  `

  createReviewTable := `
    CREATE TABLE IF NOT EXISTS reviews (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        description TEXT NOT NULL,
        author_id INT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
  `

  createRaportTable := `
    CREATE TABLE IF NOT EXISTS raports (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        description TEXT NOT NULL,
        author_id INT NOT NULL,
        status VARCHAR(50) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
  `

  // Execute table creation queries
  tables := []string{
    createUserTable, createTarifTable, createSubscribeTable,
    createGroupTable, createSubscribeGroupTable, createActivityTable,
    createRequestTable, createCategoryTable, createSubcategoryTable,
    createSpecializationTable, createReviewTable, createRaportTable,
  }

  for _, tableQuery := range tables {
    _, err = conn.Exec(context.Background(), tableQuery)
    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }
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

func (p *Postgre) GetTarifs() ([]tarif.Tarif, error) {
  const op = "storage.postgre.GetTarifs"
  query := `
    SELECT id, title, description, price, image_url
    FROM tarifs
  `

  rows, err := p.conn.Query(context.Background(), query)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var tarifs []tarif.Tarif

  for rows.Next() {
    var tarif tarif.Tarif
    err := rows.Scan(&tarif.ID, &tarif.Title, &tarif.Description, &tarif.Price, &tarif.ImageURL)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    tarifs = append(tarifs, tarif)
  }

  return tarifs, nil
}

func (p *Postgre) GetTarifByID(id int) (tarif.Tarif, error) {
  const op = "storage.postgre.GetTarifByID"
  query := `
    SELECT id, title, description, price, image_url
    FROM tarifs
    WHERE id = $1
  `

  var tarif tarif.Tarif
  err := p.conn.QueryRow(context.Background(), query, id).Scan(&tarif.ID, &tarif.Title, &tarif.Description, &tarif.Price, &tarif.ImageURL)

  if err != nil {
    return tarif, fmt.Errorf("%s: %w", op, err)
  }

  return tarif, nil
}

func (p *Postgre) CreateTarif(title string, description string, price int, image_url string) error {
  const op = "storage.postgre.CreateTarif"
  query := `
    INSERT INTO tarifs (title, description, price, image_url)
    VALUES ($1, $2, $3, $4)
  `

  _, err := p.conn.Exec(context.Background(), query, title, description, price, image_url)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) Subscribe(userID int, subscriberID int) error {
  const op = "storage.postgre.Subscribe"
  query := `
    INSERT INTO subscribes (user_id, subscriber_id, status)
    VALUES ($1, $2, $3)
  `

  _, err := p.conn.Exec(context.Background(), query, userID, subscriberID, "Знакомый")

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) UpdateStatus(userID int, subscriberID, status string) error {
  const op = "storage.postgre.UpdateStatus"
  query := `
    UPDATE subscribes
    SET status = $1
    WHERE user_id = $2 AND subscriber_id = $3
  `
  _, err := p.conn.Exec(context.Background(), query, status, userID, subscriberID)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }
  
  return nil
}

func (p *Postgre) Unsubscribe(userID int, subscriberID int) error {
  const op = "storage.postgre.Unsubscribe"
  query := `
    DELETE FROM subscribes
    WHERE user_id = $1 AND subscriber_id = $2
  `

  _, err := p.conn.Exec(context.Background(), query, userID, subscriberID)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) GetSubscribers(userID int) ([]user.User, error) {
  const op = "storage.postgre.GetSubscribers"
  query := `
    SELECT u.id, u.username, u.chat_id, u.name, u.image_url, u.tarif_id
    FROM users u
    JOIN subscribes s ON u.id = s.subscriber_id
    WHERE s.user_id = $1
  `

  rows, err := p.conn.Query(context.Background(), query, userID)

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

func (p *Postgre) GetSubscriptions(userID int) ([]user.User, error) {
  const op = "storage.postgre.GetSubscriptions"
  query := `
    SELECT u.id, u.username, u.chat_id, u.name, u.image_url, u.tarif_id
    FROM users u
    JOIN subscribes s ON u.id = s.user_id
    WHERE s.subscriber_id = $1
  `

  rows, err := p.conn.Query(context.Background(), query, userID)

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

func (p *Postgre) CreateGroup(title string, description string, createrID int, imageUrl string, telegramLink string) error {
  const op = "storage.postgre.CreateGroup"
  query := `
    INSERT INTO groups (title, description, creater_id, image_url, telegram_link)
    VALUES ($1, $2, $3, $4, $5)
  `

  _, err := p.conn.Exec(context.Background(), query, title, description, createrID, imageUrl, telegramLink)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) SubscribeGroup(userID int, groupID int) error {
  const op = "storage.postgre.SubscribeGroup"
  query := `
    INSERT INTO subscribe_groups (user_id, group_id, status)
    VALUES ($1, $2, $3)
  `

  _, err := p.conn.Exec(context.Background(), query, userID, groupID, "active")

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) UnsubscribeGroup(userID int, groupID int) error {
  const op = "storage.postgre.UnsubscribeGroup"
  query := `
    DELETE FROM subscribe_groups
    WHERE user_id = $1 AND group_id = $2
  `

  _, err := p.conn.Exec(context.Background(), query, userID, groupID)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) GetGroups() ([]group.Group, error) {
  const op = "storage.postgre.GetGroups"
  query := `
    SELECT id, title, description, creater_id, image_url, telegram_link
    FROM groups
  `

  rows, err := p.conn.Query(context.Background(), query)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var groups []group.Group

  for rows.Next() {
    var group group.Group
    err := rows.Scan(&group.ID, &group.Title, &group.Description, &group.CreaterID, &group.ImageURL, &group.TelegramLink)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    groups = append(groups, group)
  }

  return groups, nil
}

func (p *Postgre) GetGroupByID(id int) (group.Group, error) {
  const op = "storage.postgre.GetGroupByID"
  query := `
    SELECT id, title, description, creater_id, image_url, telegram_link
    FROM groups
    WHERE id = $1
  `

  var group group.Group
  err := p.conn.QueryRow(context.Background(), query, id).Scan(&group.ID, &group.Title, &group.Description, &group.CreaterID, &group.ImageURL, &group.TelegramLink)

  if err != nil {
    return group, fmt.Errorf("%s: %w", op, err)
  }

  return group, nil
}

func (p *Postgre) GetGroupSubscribers(groupID int) ([]user.User, error) {
  const op = "storage.postgre.GetGroupSubscribers"
  query := `
    SELECT u.id, u.username, u.chat_id, u.name, u.image_url, u.tarif_id
    FROM users u
    JOIN subscribe_groups s ON u.id = s.user_id
    WHERE s.group_id = $1
  `

  rows, err := p.conn.Query(context.Background(), query, groupID)

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

func (p *Postgre) GetGroupsByCreaterID(createrID int) ([]group.Group, error) {
  const op = "storage.postgre.GetGroupsByCreaterID"
  query := `
    SELECT id, title, description, creater_id, image_url, telegram_link
    FROM groups
    WHERE creater_id = $1
  `

  rows, err := p.conn.Query(context.Background(), query, createrID)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var groups []group.Group

  for rows.Next() {
    var group group.Group
    err := rows.Scan(&group.ID, &group.Title, &group.Description, &group.CreaterID, &group.ImageURL, &group.TelegramLink)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    groups = append(groups, group)
  }

  return groups, nil
}

func (p *Postgre) CreateActivity(userID int, title string, description string, location string, categoryID int, imageUrl string, subcategoryID int, specializationID int) error {
  const op = "storage.postgre.CreateActivity"
  query := `
    INSERT INTO activities (user_id, title, description, location, category_id, image_url, subcategory_id, specialization_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
  `

  _, err := p.conn.Exec(context.Background(), query, userID, title, description, location, categoryID, imageUrl, subcategoryID, specializationID)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) GetActivities() ([]activity.Activity, error) {
  const op = "storage.postgre.GetActivities"
  query := `
    SELECT id, user_id, title, description, location, category_id, image_url, subcategory_id, specialization_id
    FROM activities
  `

  rows, err := p.conn.Query(context.Background(), query)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var activities []activity.Activity

  for rows.Next() {
    var activity activity.Activity
    err := rows.Scan(&activity.ID, &activity.UserID, &activity.Title, &activity.Description, &activity.Location, &activity.CategoryID, &activity.ImageURL, &activity.SubCategoryID, &activity.SpecializationID)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    activities = append(activities, activity)
  }

  return activities, nil
}

func (p *Postgre) GetActivityByID(id int) (activity.Activity, error) {
  const op = "storage.postgre.GetActivityByID"
  query := `
    SELECT id, user_id, title, description, location, category_id, image_url, subcategory_id, specialization_id
    FROM activities
    WHERE id = $1
  `

  var activity activity.Activity
  err := p.conn.QueryRow(context.Background(), query, id).Scan(&activity.ID, &activity.UserID, &activity.Title, &activity.Description, &activity.Location, &activity.CategoryID, &activity.ImageURL, &activity.SubCategoryID, &activity.SpecializationID)

  if err != nil {
    return activity, fmt.Errorf("%s: %w", op, err)
  }

  return activity, nil
}

func (p *Postgre) GetActivitiesByUserID(userID int) ([]activity.Activity, error) {
  const op = "storage.postgre.GetActivitiesByUserID"
  query := `
    SELECT id, user_id, title, description, location, category_id, image_url, subcategory_id, specialization_id
    FROM activities
    WHERE user_id = $1
  `

  rows, err := p.conn.Query(context.Background(), query, userID)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var activities []activity.Activity

  for rows.Next() {
    var activity activity.Activity
    err := rows.Scan(&activity.ID, &activity.UserID, &activity.Title, &activity.Description, &activity.Location, &activity.CategoryID, &activity.ImageURL, &activity.SubCategoryID, &activity.SpecializationID)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    activities = append(activities, activity)
  }

  return activities, nil
}

func (p *Postgre) UpdateActivity(id int, title string, description string, location string, categoryID int, imageUrl string, subcategoryID int, specializationID int) error {
  const op = "storage.postgre.UpdateActivity"

  query := `
    UPDATE activities
    SET title = $1, description = $2, location = $3, category_id = $4, image_url = $5, subcategory_id = $6, specialization_id = $7
    WHERE id = $8
  `
  _, err := p.conn.Exec(context.Background(), query, title, description, location, categoryID, imageUrl, subcategoryID, specializationID, id)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) CreateRequest(userID int, title string, description string, location string, categoryID int, subcategoryID int, specializationID int) error {
  const op = "storage.postgre.CreateRequest"
  query := `
    INSERT INTO requests (user_id, title, description, location, category_id, subcategory_id, specialization_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
  `

  _, err := p.conn.Exec(context.Background(), query, userID, title, description, location, categoryID, subcategoryID, specializationID)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) GetRequests() ([]request.Request, error) {
  const op = "storage.postgre.GetRequests"
  query := `
    SELECT id, user_id, title, description, location, category_id, subcategory_id, specialization_id
    FROM requests
  `

  rows, err := p.conn.Query(context.Background(), query)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var requests []request.Request

  for rows.Next() {
    var request request.Request
    err := rows.Scan(&request.ID, &request.UserID, &request.Title, &request.Description, &request.Location, &request.CategoryID, &request.SubCategoryID, &request.SpecializationID)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    requests = append(requests, request)
  }

  return requests, nil
}

func (p *Postgre) GetRequestByID(id int) (request.Request, error) {
  const op = "storage.postgre.GetRequestByID"
  query := `
    SELECT id, user_id, title, description, location, category_id, subcategory_id, specialization_id
    FROM requests
    WHERE id = $1
  `

  var request request.Request
  err := p.conn.QueryRow(context.Background(), query, id).Scan(&request.ID, &request.UserID, &request.Title, &request.Description, &request.Location, &request.CategoryID, &request.SubCategoryID, &request.SpecializationID)

  if err != nil {
    return request, fmt.Errorf("%s: %w", op, err)
  }

  return request, nil
}

func (p *Postgre) GetRequestsByUserID(userID int) ([]request.Request, error) {
  const op = "storage.postgre.GetRequestsByUserID"
  query := `
    SELECT id, user_id, title, description, location, category_id, subcategory_id, specialization_id
    FROM requests
    WHERE user_id = $1
  `

  rows, err := p.conn.Query(context.Background(), query, userID)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var requests []request.Request

  for rows.Next() {
    var request request.Request
    err := rows.Scan(&request.ID, &request.UserID, &request.Title, &request.Description, &request.Location, &request.CategoryID, &request.SubCategoryID, &request.SpecializationID)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    requests = append(requests, request)
  }

  return requests, nil
}

func (p *Postgre) UpdateRequest(id int, title string, description string, location string, categoryID int, subcategoryID int, specializationID int) error {
  const op = "storage.postgre.UpdateRequest"

  query := `
    UPDATE requests
    SET title = $1, description = $2, location = $3, category_id = $4, subcategory_id = $5, specialization_id = $6
    WHERE id = $7
  `
  _, err := p.conn.Exec(context.Background(), query, title, description, location, categoryID, subcategoryID, specializationID, id)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) CreateCategory(title string) error {
  const op = "storage.postgre.CreateCategory"
  query := `
    INSERT INTO categories (title)
    VALUES ($1)
  `

  _, err := p.conn.Exec(context.Background(), query, title)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) GetCategories() ([]category.Category, error) {
  const op = "storage.postgre.GetCategories"
  query := `
    SELECT id, title
    FROM categories
  `

  rows, err := p.conn.Query(context.Background(), query)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var categories []category.Category

  for rows.Next() {
    var category category.Category
    err := rows.Scan(&category.ID, &category.Title)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    categories = append(categories, category)
  }

  return categories, nil
}

func (p *Postgre) GetCategoryByID(id int) (category.Category, error) {
  const op = "storage.postgre.GetCategoryByID"
  query := `
    SELECT id, title
    FROM categories
    WHERE id = $1
  `

  var category category.Category
  err := p.conn.QueryRow(context.Background(), query, id).Scan(&category.ID, &category.Title)

  if err != nil {
    return category, fmt.Errorf("%s: %w", op, err)
  }

  return category, nil
}

func (p *Postgre) CreateSubcategory(title string, categoryID int) error {
  const op = "storage.postgre.CreateSubcategory"
  query := `
    INSERT INTO subcategories (title, category_id)
    VALUES ($1, $2)
  `

  _, err := p.conn.Exec(context.Background(), query, title, categoryID)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) GetSubcategories() ([]subcategory.SubCategory, error) {
  const op = "storage.postgre.GetSubcategories"
  query := `
    SELECT id, title, category_id
    FROM subcategories
  `

  rows, err := p.conn.Query(context.Background(), query)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var subcategories []subcategory.SubCategory

  for rows.Next() {
    var subcategory subcategory.SubCategory
    err := rows.Scan(&subcategory.ID, &subcategory.Title, &subcategory.CategoryID)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    subcategories = append(subcategories, subcategory)
  }

  return subcategories, nil
}

func (p *Postgre) GetSubcategoryByID(id int) (subcategory.SubCategory, error) {
  const op = "storage.postgre.GetSubcategoryByID"
  query := `
    SELECT id, title, category_id
    FROM subcategories
    WHERE id = $1
  `

  var subcategory subcategory.SubCategory
  err := p.conn.QueryRow(context.Background(), query, id).Scan(&subcategory.ID, &subcategory.Title, &subcategory.CategoryID)

  if err != nil {
    return subcategory, fmt.Errorf("%s: %w", op, err)
  }

  return subcategory, nil
}

func (p *Postgre) CreateSpecialization(title string, subcategoryID int) error {
  const op = "storage.postgre.CreateSpecialization"
  query := `
    INSERT INTO specializations (title, subcategory_id)
    VALUES ($1, $2)
  `

  _, err := p.conn.Exec(context.Background(), query, title, subcategoryID)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) GetSpecializations() ([]specialization.Specializtion, error) {
  const op = "storage.postgre.GetSpecializations"
  query := `
    SELECT id, title, subcategory_id
    FROM specializations
  `

  rows, err := p.conn.Query(context.Background(), query)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var specializations []specialization.Specializtion

  for rows.Next() {
    var specialization specialization.Specializtion
    err := rows.Scan(&specialization.ID, &specialization.Title, &specialization.SubCategoryID)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    specializations = append(specializations, specialization)
  }

  return specializations, nil
}

func (p *Postgre) GetSpecializationByID(id int) (specialization.Specializtion, error) {
  const op = "storage.postgre.GetSpecializationByID"
  query := `
    SELECT id, title, subcategory_id
    FROM specializations
    WHERE id = $1
  `

  var specialization specialization.Specializtion
  err := p.conn.QueryRow(context.Background(), query, id).Scan(&specialization.ID, &specialization.Title, &specialization.SubCategoryID)

  if err != nil {
    return specialization, fmt.Errorf("%s: %w", op, err)
  }

  return specialization, nil
}

func (p *Postgre) CreateReview(userID int, description string, authorID int) error {
  const op = "storage.postgre.CreateReview"
  query := `
    INSERT INTO reviews (user_id, description, author_id)
    VALUES ($1, $2, $3)
  `

  _, err := p.conn.Exec(context.Background(), query, userID, description, authorID)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) GetReviews() ([]review.Review, error) {
  const op = "storage.postgre.GetReviews"
  query := `
    SELECT id, user_id, description, author_id
    FROM reviews
  `

  rows, err := p.conn.Query(context.Background(), query)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var reviews []review.Review

  for rows.Next() {
    var review review.Review
    err := rows.Scan(&review.ID, &review.UserID, &review.Description, &review.AuthorID)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    reviews = append(reviews, review)
  }

  return reviews, nil
}

func (p *Postgre) GetReviewByID(id int) (review.Review, error) {
  const op = "storage.postgre.GetReviewByID"
  query := `
    SELECT id, user_id, description, author_id
    FROM reviews
    WHERE id = $1
  `

  var review review.Review
  err := p.conn.QueryRow(context.Background(), query, id).Scan(&review.ID, &review.UserID, &review.Description, &review.AuthorID)

  if err != nil {
    return review, fmt.Errorf("%s: %w", op, err)
  }

  return review, nil
}

func (p *Postgre) GetReviewsByUserID(userID int) ([]review.Review, error) {
  const op = "storage.postgre.GetReviewsByUserID"
  query := `
    SELECT id, user_id, description, author_id
    FROM reviews
    WHERE user_id = $1
  `

  rows, err := p.conn.Query(context.Background(), query, userID)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var reviews []review.Review

  for rows.Next() {
    var review review.Review
    err := rows.Scan(&review.ID, &review.UserID, &review.Description, &review.AuthorID)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    reviews = append(reviews, review)
  }

  return reviews, nil
}

func (p *Postgre) CreateRaport(userID int, description string, authorID int, status string) error {
  const op = "storage.postgre.CreateRaport"
  query := `
    INSERT INTO raports (user_id, description, author_id, status)
    VALUES ($1, $2, $3, $4)
  `

  _, err := p.conn.Exec(context.Background(), query, userID, description, authorID, status)

  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  return nil
}

func (p *Postgre) GetRaports() ([]raport.Raport, error) {
  const op = "storage.postgre.GetRaports"
  query := `
    SELECT id, user_id, description, author_id, status
    FROM raports
  `

  rows, err := p.conn.Query(context.Background(), query)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var raports []raport.Raport

  for rows.Next() {
    var raport raport.Raport
    err := rows.Scan(&raport.ID, &raport.UserID, &raport.Description, &raport.AuthorID, &raport.Status)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    raports = append(raports, raport)
  }

  return raports, nil
}

func (p *Postgre) GetRaportByID(id int) (raport.Raport, error) {
  const op = "storage.postgre.GetRaportByID"
  query := `
    SELECT id, user_id, description, author_id, status
    FROM raports
    WHERE id = $1
  `

  var raport raport.Raport
  err := p.conn.QueryRow(context.Background(), query, id).Scan(&raport.ID, &raport.UserID, &raport.Description, &raport.AuthorID, &raport.Status)

  if err != nil {
    return raport, fmt.Errorf("%s: %w", op, err)
  }

  return raport, nil
}

func (p *Postgre) GetRaportsByUserID(userID int) ([]raport.Raport, error) {
  const op = "storage.postgre.GetRaportsByUserID"
  query := `
    SELECT id, user_id, description, author_id, status
    FROM raports
    WHERE user_id = $1
  `

  rows, err := p.conn.Query(context.Background(), query, userID)

  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  defer rows.Close()

  var raports []raport.Raport

  for rows.Next() {
    var raport raport.Raport
    err := rows.Scan(&raport.ID, &raport.UserID, &raport.Description, &raport.AuthorID, &raport.Status)

    if err != nil {
      return nil, fmt.Errorf("%s: %w", op, err)
    }

    raports = append(raports, raport)
  }

  return raports, nil
}
