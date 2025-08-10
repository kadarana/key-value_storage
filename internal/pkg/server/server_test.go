package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myproj/internal/pkg/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	store, err := storage.NewStorage()
	if err != nil {
		t.Errorf("Initialize error")
	}

	serve := New(&store)

	w := httptest.NewRecorder() // имитирует интрефейс ResponseWriter, позволяет ловить ответ сервера без реального сетевого запроса

	req, _ := http.NewRequest(http.MethodGet, "/health", nil) // создает http запрос без реального сетевого взаимодействия
	serve.newAPI().ServeHTTP(w, req)                          // выполняет имитацию http-запроса к api для тестирования

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestSET(t *testing.T) {
	store, err := storage.NewStorage()
	if err != nil {
		t.Errorf("Initialize error")
	}

	serve := New(&store)

	testKey := []string{"key1", "key2", "key3"}
	testVal := []any{"hello", 1221, 1221.07}

	expected := []any{http.StatusOK, http.StatusOK, http.StatusBadGateway}

	for i, k := range testKey {
		testVal := Entry{
			Value: testVal[i],
		}

		jsonVal, _ := json.MarshalIndent(testVal, "", "\t")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/scalar/set/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)
		fmt.Print(w)
		assert.Equal(t, expected[i], w.Code)
	}
}

func TestGET(t *testing.T) {
	store, err := storage.NewStorage()
	if err != nil {
		t.Errorf("Initialize error")
	}

	serve := New(&store)

	testKeys := []string{"key1", "key2", "key3"}
	testVals := []any{float64(111), "val2", float64(1234)}

	for i, k := range testKeys {
		testVal := Entry{
			Value: testVals[i],
		}

		jsonVal, _ := json.MarshalIndent(testVal, "", "\t")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/scalar/set/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)

		req, _ = http.NewRequest(http.MethodGet, "/scalar/get/"+k, nil)
		serve.newAPI().ServeHTTP(w, req)

		var val Entry
		json.Unmarshal(w.Body.Bytes(), val)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, val.Value, testVals[i])
	}

}

func TestLPUSH(t *testing.T) {
	store, err := storage.NewStorage()
	if err != nil {
		t.Errorf("Initialize error")
	}

	serve := New(&store)

	testKeys := []string{"key1", "key2", "key3", "key4"}
	testVals := [][]any{
		{1, 2, 3, 4, 5, 6},
		{"1", "2", "arr", "3"},
		{"1", "val2", 33, "99.9"},
		{1, 2, 3, 5, 24.8},
	}

	expected := []any{http.StatusOK, http.StatusOK, http.StatusOK, http.StatusBadGateway}

	for i, k := range testKeys {
		testVal := EntryArray{
			Value: testVals[i],
		}

		jsonVal, _ := json.MarshalIndent(testVal, "", "\t")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/array/lpush/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)
		fmt.Print(w)

		assert.Equal(t, expected[i], w.Code)
	}
}

func TestLPOP(t *testing.T) {
	store, err := storage.NewStorage()
	if err != nil {
		t.Errorf("Initialize error")
	}

	serve := New(&store)
	testKeys := []string{"key1", "key2", "key3"}
	testVals := [][]any{
		{1, 2, 3, 4, 5, 6},
		{"1", "2", "arr", "3", "1234", "2345"},
		{1, 2, "3", "4.04", "5", 6},
	}

	testSlices := [][]any{
		{2},
		{2, -2},
		{},
	}

	expected := [][]any{
		{float64(1), float64(2)},
		{"1", "2"},
		{float64(1)},
	}

	for i, k := range testKeys {
		testVal := EntryArray{
			Value: testVals[i],
		}

		testSlice := EntryArray{
			Value: testSlices[i],
		}

		jsonVal, _ := json.MarshalIndent(testVal, "", "\t")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/array/rpush/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		jsonVal, _ = json.MarshalIndent(testSlice, "", "\t")
		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, "/array/lpop/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)

		var val Entry
		json.Unmarshal(w.Body.Bytes(), &val)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected[i], val.Value)

	}
}

func TestRPUSH(t *testing.T) {
	store, err := storage.NewStorage()
	if err != nil {
		t.Errorf("Initialize error")
	}

	serve := New(&store)
	testKeys := []string{"key1", "key2", "key3", "key4"}
	testVals := [][]any{
		{1, 2, 3, 4, 5, 6},
		{"1", "2", "arr", "3"},
		{"1", "val2", 33, "99.9"},
		{1, 2, 3, 5, 24.8},
	}

	expected := []any{http.StatusOK, http.StatusOK, http.StatusOK, http.StatusBadGateway}

	for i, k := range testKeys {
		testVal := EntryArray{
			Value: testVals[i],
		}

		jsonVal, _ := json.MarshalIndent(testVal, "", "\t")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/array/rpush/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)
		fmt.Print(w)

		assert.Equal(t, expected[i], w.Code)
	}
}

func TestRADDTOSET(t *testing.T) {
	store, err := storage.NewStorage()
	if err != nil {
		t.Errorf("Initialize error")
	}

	serve := New(&store)
	testKeys := []string{"key1", "key2", "key3"}
	testVals := [][]any{
		{1, 2, 3, 4, 5, 6},
		{"1", "2", "arr", "3", "1234", "2345"},
		{1, 2, "3", "4.04", "5", 6},
	}

	testSlices := [][]any{
		{2},
		{2, -2},
		{},
	}

	expected := [][]any{
		{float64(6), float64(5)},
		{"1234", "3", "arr"},
		{float64(6)},
	}

	for i, k := range testKeys {
		testVal := EntryArray{
			Value: testVals[i],
		}

		testSlice := EntryArray{
			Value: testSlices[i],
		}

		jsonVal, _ := json.MarshalIndent(testVal, "", "\t")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/array/rpush/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		jsonVal, _ = json.MarshalIndent(testSlice, "", "\t")
		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodPost, "/array/raddtoset/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)

		var val Entry
		json.Unmarshal(w.Body.Bytes(), &val)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected[i], val.Value)

	}
}

func TestRPOP(t *testing.T) {
	store, err := storage.NewStorage()
	if err != nil {
		t.Errorf("Initialize error")
	}

	serve := New(&store)
	testKeys := []string{"key1", "key2", "key3"}
	testVals := [][]any{
		{1, 2, 3, 4, 5, 6},
		{"1", "2", "arr", "3", "1234", "2345"},
		{1, 2, "3", "4.04", "5", 6},
	}

	testSlices := [][]any{
		{2},
		{2, -2},
		{},
	}

	expected := [][]any{
		{float64(1), float64(2)},
		{"1", "2"},
		{float64(1)},
	}

	for i, k := range testKeys {
		testVal := EntryArray{
			Value: testVals[i],
		}

		testSlice := EntryArray{
			Value: testSlices[i],
		}

		jsonVal, _ := json.MarshalIndent(testVal, "", "\t")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/array/rpush/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		jsonVal, _ = json.MarshalIndent(testSlice, "", "\t")
		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, "/array/rpop/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)

		var val Entry
		json.Unmarshal(w.Body.Bytes(), &val)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected[i], val.Value)

	}
}

func TestLSET(t *testing.T) {
	store, err := storage.NewStorage()
	if err != nil {
		t.Errorf("Initialize error")
	}

	serve := New(&store)
	testKeys := []string{"key1", "key2", "key3"}
	testVals := [][]any{
		{1, 2, 3, 4, 5, 6},
		{"1", "2", "arr", "3"},
		{1, 2, "3"},
	}

	testArgs := [][]any{
		{1, "1"},
		{3, 1},
		{2, 0},
	}

	for i, k := range testKeys {
		testVal := EntryArray{
			Value: testVals[i],
		}

		testArg := EntryArray{
			Value: testArgs[i],
		}

		jsonVal, _ := json.MarshalIndent(testVal, "", "\t")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/array/lpush/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		jsonVal, _ = json.MarshalIndent(testArg, "", "\t")
		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodPost, "/array/lset/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestLGET(t *testing.T) {
	store, err := storage.NewStorage()
	if err != nil {
		t.Errorf("Initialize error")
	}

	serve := New(&store)
	testKeys := []string{"key1", "key2", "key3"}
	testVals := [][]any{
		{1, 2, 3, 4, 5, 6},
		{"1", "2", "arr", "3"},
		{1, 2, "3"},
	}

	testArgs := [][]any{
		{1, "1"},
		{3, float64(1)},
		{2, float64(0)},
	}

	testArgsGet := [][]any{
		{1},
		{3},
		{2},
	}

	for i, k := range testKeys {
		testVal := EntryArray{
			Value: testVals[i],
		}

		testArg := EntryArray{
			Value: testArgs[i],
		}

		testArgGet := EntryArray{
			Value: testArgsGet[i],
		}

		jsonVal, _ := json.MarshalIndent(testVal, "", "\t")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/array/lpush/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		jsonVal, _ = json.MarshalIndent(testArg, "", "\t")
		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodPost, "/array/lset/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		jsonVal, _ = json.MarshalIndent(testArgGet, "", "\t")
		req, _ = http.NewRequest(http.MethodGet, "/array/lget/"+k, bytes.NewBuffer(jsonVal))
		serve.newAPI().ServeHTTP(w, req)

		var val Entry
		json.Unmarshal(w.Body.Bytes(), &val)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, testArgs[i][1], val.Value)
	}
}
