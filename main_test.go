package main

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis"
	"testing"
)

func TestSomeRepositoryImpl_GetData(t *testing.T) {
	// Создание мок-объекта для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	// Создание объекта SomeRepositoryImpl с мок-базой данных
	repo := &SomeRepositoryImpl{
		db: db,
	}

	// Создание строкового значения, которое будет возвращаться из базы данных
	expectedData := "test data"

	// Устанавливаем ожидание запроса к базе данных
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"data"}).AddRow(expectedData))

	// Вызываем метод GetData для получения данных
	data, err := repo.GetData("request")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем, что данные, полученные из базы данных, совпадают с ожидаемыми
	if data != expectedData {
		t.Fatalf("expected data: %s, got: %s", expectedData, data)
	}

	// Проверяем, что все ожидания были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestSomeRepositoryProxy_GetData(t *testing.T) {
	// Создаем мок для клиента Redis
	client := redis.NewClient(&redis.Options{})

	// Создаем объект SomeRepositoryProxy
	repoProxy := &SomeRepositoryProxy{
		repository: &SomeRepositoryImpl{}, // Можно использовать мок SomeRepositoryImpl для более точного тестирования
		cache:      *client,
	}

	// Проверка кэша на наличие данных
	key := "test_key"
	expectedData := "test_data"
	client.Set(key, expectedData, 0)

	// Вызываем GetData с существующим ключом
	data, err := repoProxy.GetData(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем, что данные, полученные из кэша, совпадают с ожидаемыми
	if data != expectedData {
		t.Fatalf("expected data: %s, got: %s", expectedData, data)
	}

	// Проверяем поведение, когда данных нет в кэше
	keyNotExists := "key_not_exists"
	dataNotExists, err := repoProxy.GetData(keyNotExists)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем, что полученные данные не пусты, так как они не были найдены в кэше
	if dataNotExists == "" {
		t.Fatalf("expected non-empty data, got empty")
	}
}
