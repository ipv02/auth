package main

import (
	"context"
	"log"
	"time"

	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/ipv02/auth/pkg/user_v1"
	"github.com/jackc/pgx/v4/pgxpool"
)

const dbDSN = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password"

func main() {
	ctx := context.Background()

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Делаем запрос на вставку записи в таблицу auth
	builderInsert := sq.Insert("auth").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role", "password", "password_confirm").
		Values(gofakeit.Name(), gofakeit.Email(), 1, "password", "password").
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var authID int
	err = pool.QueryRow(ctx, query, args...).Scan(&authID)
	if err != nil {
		log.Fatalf("failed to insert auth: %v", err)
	}

	log.Printf("inserted auth with id: %d", authID)

	// Делаем запрос на выборку записей из таблицы auth
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to select auth: %v", err)
	}

	var id int
	var name, email string
	var createdAt time.Time
	var updatedAt sql.NullTime
	var role user_v1.UserRole

	for rows.Next() {
		err = rows.Scan(&id, &name, &email, &createdAt, &updatedAt, &role)
		if err != nil {
			log.Fatalf("failed to scan auth: %v", err)
		}

		log.Printf("id: %d, name: %s, email: %s,created_at: %v, updated_at: %v\n, role: %s", id, name, email, createdAt, updatedAt, role)
	}

	// Делаем запрос на обновление записи в таблице auth
	builderUpdate := sq.Update("auth").
		PlaceholderFormat(sq.Dollar).
		Set("name", gofakeit.Name()).
		Set("email", gofakeit.Email()).
		Set("updated_at", time.Now()).
		Set("role", role).
		Where(sq.Eq{"id": authID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update auth: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	// Делаем запрос на получение измененной записи из таблицы auth
	builderSelectOne := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": authID}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &createdAt, &updatedAt, &role)
	if err != nil {
		log.Fatalf("failed to select auth: %v", err)
	}

	log.Printf("id: %d, name: %s, email: %s,created_at: %v, updated_at: %v\n, role: %s", id, name, email, createdAt, updatedAt, role)
}
