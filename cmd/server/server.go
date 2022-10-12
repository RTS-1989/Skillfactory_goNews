package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/mongoDB"
	"GoNews/pkg/storage/postgres"
	"context"
	"log"
	"net/http"
	"time"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Создаём объекты баз данных.
	//БД в памяти.
	db := memdb.New()

	// Реляционная БД PostgreSQL.
	db2, err := postgres.New(ctx, "postgres://admin:!qaz@wsx@localhost:5432/posts")
	if err != nil {
		log.Fatal(err)
	}
	defer db2.Db.Close()

	// Документная БД MongoDB.
	db3, err := mongoDB.New(ctx, "mongodb://localhost:27017/")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = db3.Db.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	_, _ = db2, db3

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}
