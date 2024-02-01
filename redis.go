package repositorysdk

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisRepository interface {
	SaveCache(string, interface{}, int) error
	SaveHashCache(string, string, string, int) error
	SaveAllHashCache(string, map[string]string, int) error
	AddSetMember(key string, ttl int, member ...interface{}) error
	GetCache(string, interface{}) error
	GetHashCache(string, string) (string, error)
	GetAllHashCache(string) (map[string]string, error)
	RemoveCache(string) error
	RemoveSetMember(key string, member interface{}) error
	RemoveHashCache(key string, field string) error
	SetExpire(string, int) error
	CheckSetMember(key string, member interface{}) (bool, error)
	Exist(key string) (bool, error)
}

const RedisKeepTTL = 0

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{client: client}
}

// SaveCache saves cache to redis by using the command `SET`.
// Zero expiration time means no expiration time for cache.
//
// Parameters:
// - key: the cache key.
// - value: the cache value to be saved.
// - ttl: the expiration time for cache in seconds, 0 means no expiration time.
//
// Returns:
// - err: an error if something goes wrong, otherwise nil.
func (r *redisRepository) SaveCache(key string, value interface{}, ttl int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	v, err := json.Marshal(value)
	if err != nil {
		return
	}

	return r.client.Set(ctx, key, v, time.Duration(ttl)*time.Second).Err()
}

// SaveHashCache saves a single field cache to redis.
//
// Parameters:
// - key: the cache key.
// - field: the cache field to be saved.
// - value: the cache value to be saved.
// - ttl: the expiration time for cache in seconds.
//
// Returns:
// - err: an error if something goes wrong, otherwise nil.
func (r *redisRepository) SaveHashCache(key string, field string, value string, ttl int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.client.HSet(ctx, key, field, value).Err(); err != nil {
		return err
	}

	if ttl > 0 {
		return r.client.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
	}

	return nil
}

// SaveAllHashCache saves multiple field cache to redis.
//
// Parameters:
// - key: the cache key.
// - value: a map containing the fields and values to be saved.
// - ttl: the expiration time for cache in seconds.
//
// Returns:
// - err: an error if something goes wrong, otherwise nil.
func (r *redisRepository) SaveAllHashCache(key string, value map[string]string, ttl int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.client.HSet(ctx, key, value).Err(); err != nil {
		return err
	}

	if ttl > 0 {
		return r.client.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()

	}

	return nil
}

// GetHashCache retrieves a single field cache from redis.
//
// Parameters:
// - key: the cache key.
// - field: the cache field to be retrieved.
//
// Returns:
// - string: the cache value if it exists.
// - err: an error if something goes wrong, otherwise nil.
func (r *redisRepository) GetHashCache(key string, field string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.client.HGet(ctx, key, field).Result()
}

// GetAllHashCache retrieves all fields of a hash cache from redis.
//
// Parameters:
// - key: the cache key.
//
// Returns:
// - map[string]string: a map containing all the fields and their values if the hash exists.
// - error: an error if something goes wrong, otherwise nil.
func (r *redisRepository) GetAllHashCache(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.client.HGetAll(ctx, key).Result()
}

// RemoveHashCache remove a single field of hash cache.
//
// Parameters:
// - key: the cache key.
// - field: the cache field to be saved.
//
// Returns:
// - err: an error if something goes wrong, otherwise nil.
func (r *redisRepository) RemoveHashCache(key string, field string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.client.HDel(ctx, key, field).Err()
}

// GetCache retrieves a cache from redis.
//
// Parameters:
// - key: the cache key.
// - value: a pointer to the object that will hold the unmarshalled cache value.
//
// Returns:
// - error: an error if something goes wrong, otherwise nil.
func (r *redisRepository) GetCache(key string, value interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	v, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return
	}

	return json.Unmarshal([]byte(v), value)
}

// RemoveCache removes a cache from redis.
//
// Parameters:
// - key: the cache key to be removed.
//
// Returns:
// - error: an error if something goes wrong, otherwise nil.
func (r *redisRepository) RemoveCache(key string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = r.client.Del(ctx, key).Result()
	return err
}

// CheckSetMember check is member existed in the set
//
// Parameters:
// - key: the member to check.
//
// Return values:
// - bool: true if the key exists, false otherwise.
// - error: if the Redis operation fails.
func (r *redisRepository) CheckSetMember(key string, member interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.client.SIsMember(ctx, key, member).Result()
}

// AddSetMember add member to set
//
// Parameters:
// - key: the member to check.
// - member: the member.
// - ttl: expiration time of this cache.
//
// Return values:
// - error: if the Redis operation fails.
func (r *redisRepository) AddSetMember(key string, ttl int, member ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.client.SAdd(ctx, key, member).Err(); err != nil {
		return err
	}

	if ttl > 0 {
		return r.client.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
	}

	return nil
}

// RemoveSetMember remove member from set
//
// Parameters:
// - key: the member to check.
// - member: the member.
//
// Return values:
// - error: if the Redis operation fails.
func (r *redisRepository) RemoveSetMember(key string, member interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.client.SRem(ctx, key, member).Err()
}

// SetExpire sets an expiration time for a cache in redis.
//
// Parameters:
// - key: the cache key to set expiration for.
// - ttl: the expiration time for cache in seconds.
//
// Returns:
// - error: an error if something goes wrong, otherwise nil.
func (r *redisRepository) SetExpire(key string, ttl int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.client.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
}

// Exist checks if a key exists in the Redis database.
// Parameters:
// - key: the key to check.
//
// Return values:
// - bool: true if the key exists, false otherwise.
// - error: if the Redis operation fails.
func (r *redisRepository) Exist(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := r.client.Exists(ctx, key)

	return res.Val() == 1, res.Err()
}
