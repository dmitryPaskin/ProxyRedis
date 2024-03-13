package main

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
)

type SomeRepository interface {
	GetData(request string) (string, error)
}

type SomeRepositoryImpl struct {
	db *sql.DB
}

func (r *SomeRepositoryImpl) GetData(request string) (string, error) {
	// Здесь происходит запрос к базе данных

	// Примерный SQL - запрос
	query := fmt.Sprintf("SELECT %s FROM table_name WHERE condition", request)

	//Выполнение запроса к базе данных
	rows, err := r.db.Query(query)
	if err != nil {
		return "", err
	}

	var data string
	for rows.Next() {
		err := rows.Scan(&data)
		if err != nil {
			return "", err
		}
	}
	return data, nil
}

type SomeRepositoryProxy struct {
	repository SomeRepository
	cache      redis.Client
}

func (r *SomeRepositoryProxy) GetData(request string) (string, error) {
	// Здесь происходит проверка наличия данных в кэше
	data, err := r.cache.Get(request).Result()
	if err == nil && data != "" {
		// Если данные есть в кэше, то они возвращаются
		return data, err
	} else if err != nil {
		return "", err
	}
	// Если данных нет в кэше, то они запрашиваются у оригинального объекта и сохраняются в кэш
	data, err = r.repository.GetData(request)
	if err != nil {
		return "", err
	}
	// Сохраняем данные в кэше
	err = r.cache.Set(request, data, 0).Err()
	if err != nil {
		return "", err
	}

	return data, nil
}
