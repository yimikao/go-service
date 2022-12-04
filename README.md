## Design Database

1. start with the models
```go

type User struct {
	Name string `json:"name"`
}
```

2. define actions on model (i.e) repository
	```go
type UserRepository interface {
	GetById(id int) (*User, error)
	All() ([]*User, error)
	Create(cm *User) (*User, error)
}	
```
3. define real model  layer and database source and implement repositories
	```go
type commentLayer struct {
	dbConnection *pg.DB
}
```
4. connect to your database and test connections
	```go
connect and ping
```

## Create server and setup its routes
1. routes have paths that lead to handlers
2. handlers call repositories to work with/access real or fake data layer (CRUD)# go-class
