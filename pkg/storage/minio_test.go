package storage

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testStore *Storage

func TestMain(m *testing.M) {
	testStore = New(
		Config{
			Bucket:   "fiber-bucket",
			Endpoint: "localhost:9000",
			Credentials: Credentials{
				AccessKeyID:     "minio-user",
				SecretAccessKey: "minio-password",
			},
			Reset: true,
		},
	)

	code := m.Run()

	_ = testStore.Close()
	os.Exit(code)
}

func Test_Get(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStore.Set(key, val, 0)
	require.NoError(t, err)

	result, err := testStore.Get(key)
	require.NoError(t, err)
	require.Equal(t, val, result)

	result, err = testStore.Get("doe")
	require.Error(t, err)
	require.Zero(t, len(result))
}

func Test_Get_Empty_Key(t *testing.T) {
	var (
		key = ""
	)

	_, err := testStore.Get(key)
	require.Error(t, err)
	require.EqualError(t, err, "the key value is required")

}

func Test_Get_Not_Exists_Key(t *testing.T) {
	var (
		key = "not-exists"
	)

	_, err := testStore.Get(key)
	require.Error(t, err)
	require.EqualError(t, err, "The specified key does not exist.")

}

func Test_Get_Not_Exists_Bucket(t *testing.T) {
	var (
		key = "john"
	)

	// random bucket name
	testStore.cfg.Bucket = strconv.FormatInt(time.Now().UnixMicro(), 10)

	result, err := testStore.Get(key)
	require.Error(t, err)
	require.Zero(t, len(result))
	require.EqualError(t, err, "The specified bucket does not exist")

	testStore.cfg.Bucket = "fiber-bucket"

}

func Test_Set(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStore.Set(key, val, 0)
	require.NoError(t, err)
}

func Test_Set_Empty_Key(t *testing.T) {
	var (
		key = ""
		val = []byte("doe")
	)

	err := testStore.Set(key, val, 0)

	require.Error(t, err)
	require.EqualError(t, err, "the key value is required")

}

func Test_Set_Not_Exists_Bucket(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	// random bucket name
	testStore.cfg.Bucket = strconv.FormatInt(time.Now().UnixMicro(), 10)

	err := testStore.Set(key, val, 0)
	require.Error(t, err)
	require.EqualError(t, err, "The specified bucket does not exist")

	testStore.cfg.Bucket = "fiber-bucket"
}

func Test_Delete(t *testing.T) {
	var (
		key = "john"
		val = []byte("doe")
	)

	err := testStore.Set(key, val, 0)
	require.NoError(t, err)

	err = testStore.Delete(key)
	require.NoError(t, err)

}

func Test_Delete_Empty_Key(t *testing.T) {
	var (
		key = ""
		val = []byte("doe")
	)

	err := testStore.Set(key, val, 0)

	require.Error(t, err)
	require.EqualError(t, err, "the key value is required")

}

func Test_Delete_Not_Exists_Bucket(t *testing.T) {
	var (
		key = "john"
	)

	// random bucket name
	testStore.cfg.Bucket = strconv.FormatInt(time.Now().UnixMicro(), 10)

	err := testStore.Delete(key)

	require.Error(t, err)
	require.EqualError(t, err, "The specified bucket does not exist")

	testStore.cfg.Bucket = "fiber-bucket"

}

func Test_Reset(t *testing.T) {
	var (
		val = []byte("doe")
	)

	err := testStore.Set("john1", val, 0)
	require.NoError(t, err)

	err = testStore.Set("john2", val, 0)
	require.NoError(t, err)

	err = testStore.Reset()
	require.NoError(t, err)

	result, err := testStore.Get("john1")
	require.Error(t, err)
	require.Zero(t, len(result))

}

func Test_Close(t *testing.T) {
	require.NoError(t, testStore.Close())
}

func Benchmark_Minio_Set(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var err error
	for i := 0; i < b.N; i++ {
		err = testStore.Set("john", []byte("doe"), 0)
	}

	require.NoError(b, err)
}

func Benchmark_Minio_Get(b *testing.B) {
	err := testStore.Set("john", []byte("doe"), 0)
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = testStore.Get("john")
	}

	require.NoError(b, err)
}

func Benchmark_Minio_SetAndDelete(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var err error
	for i := 0; i < b.N; i++ {
		_ = testStore.Set("john", []byte("doe"), 0)
		err = testStore.Delete("john")
	}

	require.NoError(b, err)
}
