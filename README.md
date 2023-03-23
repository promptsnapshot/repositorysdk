# Table of Contents
1. [Base Entity](#about-entity)
2. [Base Dto](#about-dto)
3. [Gorm Repository](#about-newbie-repository)
4. [Redis Repository](#about-redis-repository)

# About Entity
The entity is the object that we interested in database

# Getting Start

## Types

### Base Entity
The fundamental entity

#### Structure
The base entity is contains the essential attributes that entity should have

```go
type Base struct {
	ID        *uuid.UUID     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime:nano"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime:nano"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;type:timestamp"`
}
```

#### Hook
The `Base` entity will create new uuid everytime when save to database if the id is **blank**

```go
func (b *Base) BeforeCreate(_ *gorm.DB) error {
	if b.ID == nil {
		b.ID = UUIDAdr(uuid.New())
	}

	return nil
}
```

#### Usage
When you want to define new entity you need to embed this entity

```go
type NewEntity struct{
	repositorysdk.Base
	// other fields
}
```

### PaginationMetadata
The entity for collect the metadata of pagination

#### Structure
The metadata of pagination

```go
type PaginationMetadata struct {
    ItemsPerPage int
    ItemCount    int
    TotalItem    int
    CurrentPage  int
    TotalPage    int
}
```

#### Methods

##### GetOffset
The method for get the offset value

```go
offset := meta.GetOffset()
```

#### Return

| name   | description                                                               | example |
|--------|---------------------------------------------------------------------------|---------|
| offset | the offset value (calculate from ItemPerPage value and CurrentPage value) | 10      |

##### GetItemPerPage
The method for get the item per page value

```go
itemPerPage := meta.GetItemPerPage()
```

#### Return

| name         | description                                 | example |
|--------------|---------------------------------------------|---------|
| itemPerPage  | the item per page value (min: 10, max: 100) | 10      |


##### GetCurrentPage
The method for get the current page value

```go
currentPage := meta.GetItemPerPage()
```

#### Return

| name         | description                     | example |
|--------------|---------------------------------|--------|
| currentPage  | the current page value (min: 1) | 1      |


##### ToProto
Convert to proto type

```go
metaProto := meta.ToProto()
```

#### Return

| name       | description                            | example |
|------------|----------------------------------------|---------|
| metaProto  | metadata in `*pb.PaginationMetadata`   |         |


# About DTO
Data Transfer Object is the object use for represent the attribute between the service

# Getting Start

## Types

### QueryResult
The base query result entity for opensearch

#### Structure
The base entity is contains the essential attributes that entity should have

```go
type QueryResult struct {
    Took    uint  `json:"took"`
    Timeout bool  `json:"timeout"`
    Shards  Shard `json:"_shards"`
}
```

#### Usage
When you want to define new entity you need to embed this entity

```go
type NewQueryResult struct{
	gosdk.QueryResult
	// other fields
}
```


### Shard
The stats of shards

#### Structure
The value of the statistic of shard

```go
type Shard struct {
    Total      uint `json:"total"`
    Successful uint `json:"successful"`
    Skipped    uint `json:"skipped"`
    Failed     uint `json:"failed"`
}
```

# About Newbie Repository
Gorm repository is the instance that has the code set for simple use case of gorm and can be extends by insert the scope

# Getting Start

## Connection
### PostgreSQL

return `*gorm.DB` when successfully

```go
db, err := repositorysdk.InitDatabase(PostgresDatabaseConfig, Debug)
if err != nil {
    // handle error
}
```

**Parameters**

| name                      | description        |
|---------------------------|--------------------|
| Postgres Database Config  | Postgres config    |
| Debug                     | Enable debug mode  |


**Configuration**

```go
type PostgresDatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Name     string `mapstructure:"name"`
    SSL      string `mapstructure:"ssl"`
}
```

| name     | description              | example   |
|----------|--------------------------|-----------|
| Host     | Hostname of the postgres | localhost | 
| Port     | Port of database         | 5432      |
| User     | Postgres username        | postgres  |
| Password | Postgres password        | root      |
| Name     | The database name        | postgres  |
| SSL      | SSL mode                 | disable   |

## Initialize

```go
repo := repositotysdk.NewGormRepository(*GormDB)
```

## Configuration
### Parameters

| name    | description              |
|---------|--------------------------|
| Gorm DB | the client of the gorm   |

### Return

| name | description                  | example |
|------|------------------------------|---------|
| repo | the gorm repository instance |         |

## Types

```go
type Entity interface {
	TableName() string
}
```

| name        | description                     |
|-------------|---------------------------------|
| TableName() | name of the table of the entity |

## Scopes

### Pagination

the gorm scope for pagination query

```go
if err := r.db.GetDB().
    Scopes(repositotysdk.Pagination[Entity](entityList, metadata, gormDB, ...Scope)).
    Find(&entityList).
    Error; err != nil {
    return err
}

metadata.ItemCount = len(*entityList)
metadata.CalItemPerPage()

return nil
```

#### Parameters
| name             | description              | example |
|------------------|--------------------------|---------|
| entityList       | list of entities         |         |
| metadata         | pagination metadata      |         |
| gormDB           | gorm client              |         |
| Scope            | extends scope (optional) |         |


**Example Basic**

```go
if err := r.db.GetDB().
    Scopes(repositotysdk.Pagination[Entity](entity, metadata, r.db.GetDB())).
    Find(&entity).
    Error; err != nil {
    return err
}

metadata.ItemCount = len(*entity)
metadata.CalItemPerPage()

return nil
```

**Example with Scope**

```go
if err := r.db.GetDB().
    Preload("Relationship")
    Scopes(repositotysdk.Pagination[Entity](entity, metadata, r.db.GetDB(), func(db *gorm.DB) *gorm.DB{
        return db.Where("something = ?", something)	
    })).
    Find(&entity, "something = ?", something).
    Error; err != nil {
    return err
}

metadata.ItemCount = len(*entity)
metadata.CalItemPerPage()

return nil
```

## Usage

### GetDB
return the gorm db instance

```go
db = gorm.GetDB()
```

#### Returns
| name             | description              | example |
|------------------|--------------------------|---------|
| gormDB           | gorm client              |         |

### FindAll

findAll with pagination

```go
var entityList []Entity

if err := repo.FindOne(metadata, &entityList, ...scope); err != nil{
	// handle error
}
```

#### Parameters
| name             | description              | example |
|------------------|--------------------------|---------|
| entityList       | list of entities         |         |
| metadata         | pagination metadata      |         |
| Scope            | extends scope (optional) |         |


### FindOne

findOne entity

```go
entity := Entity{}

if err := repo.FindOne(id, &entity, ...scope); err != nil{
	// handle error
}
```

#### Parameters
| name   | description                   | example |
|--------|-------------------------------|---------|
| id     | id of entity                  |         |
| entity | empty entity for receive data |         |
| Scope  | extends scope (optional)      |         |

**Example with Scope**

```go
	if err := repo.FindOne(id, &entity, func(db *gorm.DB)*gorm.DB{
	    return db.Preload("Relationship").Where("something = ?")	
    }).
    Error; err != nil {
		return err
	}
```

### Create

create entity

```go
entity := Entity{}

if err := repo.Create(&entity, ...scope); err != nil{
	// handle error
}
```

#### Parameters
| name   | description              | example |
|--------|--------------------------|---------|
| entity | entity with data         |         |
| Scope  | extends scope (optional) |         |

### Create

update entity

```go
entity := Entity{}

if err := repo.Update(id, &entity, ...scope); err != nil{
	// handle error
}
```

#### Parameters
| name   | description              | example |
|--------|--------------------------|---------|
| id     | id of entity             |         |
| entity | entity with data         |         |
| Scope  | extends scope (optional) |         |

### Delete

delete entity

```go
entity := Entity{}

if err := repo.Delete(id, &entity, ...scope); err != nil{
	// handle error
}
```

#### Parameters
| name   | description              | example |
|--------|--------------------------|---------|
| id     | id of entity             |         |
| entity | entity with data         |         |
| Scope  | extends scope (optional) |         |

# About Redis Repository
Redis repository is the repository interface for using redis work on-top of [go-redis](https://github.com/redis/go-redis)

# Getting Start

## Connection
### Redis

return `*redis.Client` when successfully

```go
 cache, err := gosdk.InitRedisConnect(RedisConfig)
 if err != nil {
    // handle error
 }
```

**Parameters**

| name          | description       |
|---------------|-------------------|
| Redis Config  | Redis config      |


**Configuration**

```go
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
```
| name     | description                                     | example        |
|----------|-------------------------------------------------|----------------|
| Host     | The host of the redis in format `hostname:port` | localhost:6379 |
| Password | Redis password                                  | password       |
| DB       | The database number                             | 0              |


## Initialization
Redis repository can be initialized by **NewRedisRepository** method

```go
repo := repositorysdk.NewRedisRepository(*RedisClient)
```

## Configuration
### Parameters

| name         | description                                 |
|--------------|---------------------------------------------|
| Redis Client | the client  of the redis for calling an API |

### Return

| name | description                   | example |
|------|-------------------------------|---------|
| repo | the redis repository instance |         |


## Usage

### SaveCache

```go
if err := repo.SaveCache(key, value, ttl); err != nil{
    // handle error
}
```

#### Parameters
| name               | description                     | example                 |
|--------------------|---------------------------------|-------------------------|
| key                | key of cache (must be `string`) | "key"                   |
| value              | value of cache (any type)       | "value", 1, &struct{}{} |
| ttl                | expiration time of cache        | 3600                    |

### SaveHashCache

```go
if err := repo.SaveHashCache(key, field, value, ttl); err != nil{
    // handle error
}
```

#### Parameters
| name  | description                            | example |
|-------|----------------------------------------|---------|
| key   | key of cache (must be `string`)        | "user"  |
| field | field of hash cache (must be `string`) | "name"  |
| value | value of hash cache (must be `string`) | "alice" |
| ttl   | expiration time of cache               | 3600    |

### SaveAllHashCache

```go
if err := repo.SaveAllHashCache(key, value, ttl); err != nil{
    // handle error
}
```

#### Parameters
| name  | description                                     | example                            |
|-------|-------------------------------------------------|------------------------------------|
| key   | key of cache (must be `string`)                 | "user"                             |
| value | map of hash cache (must be `map[string]string`) | map[string]string{"name": "alice"} |
| ttl   | expiration time of cache                        | 3600                               |

### GetCache

```go
type User struct{
	name string
}

result := User{}

if err := repo.GetCache(key, &result); err != nil{
    // handle error
}
```

#### Parameters
| name   | description                                       | example |
|--------|---------------------------------------------------|---------|
| key    | key of cache (must be `string`)                   | "user"  |
| result | the result point `struct{}` for receive the cache |         |

### GetHashCache

```go
value ,err := repo.GetHashCache(key, field)
if err != nil{
    // handle error
}
```

#### Parameters
| name   | description                             | example |
|--------|-----------------------------------------|---------|
| key    | key of cache (must be `string`)         | "user"  |
| field  | field of hash cache (must be `string`)  | "name"  |

#### Return
| name  | description                | example |
|-------|----------------------------|---------|
| value | value of cache in `string` | "alice" |

### GetAllHashCache

```go
values ,err := repo.GetAllHashCache(key, field)
if err != nil{
    // handle error
}
```

#### Parameters
| name   | description                             | example |
|--------|-----------------------------------------|---------|
| key    | key of cache (must be `string`)         | "user"  |
| field  | field of hash cache (must be `string`)  | "name"  |

#### Return
| name   | description                            | example                           |
|--------|----------------------------------------|-----------------------------------|
| values | values of cache in `map[string]string` | map[string]string{"name":"alice"} |

### RemoveCache

```go
if err := repo.RemoveCache(key); err != nil{
    // handle error
}
```

#### Parameters
| name               | description                     | example                 |
|--------------------|---------------------------------|-------------------------|
| key                | key of cache (must be `string`) | "key"                   |

### SetExpire

```go
if err := repo.SetExpire(key, ttl); err != nil{
    // handle error
}
```

#### Parameters
| name     | description                       | example   |
|----------|-----------------------------------|-----------|
| key      | key of cache (must be `string`)   | "key"     |
| ttl      | expiration time of cache          | 3600      |

