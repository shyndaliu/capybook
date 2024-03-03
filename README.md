
# Capybook API
a Golang project | [Letterboxd](https://letterboxd.com/) analog for books
 



## API Reference

#### Healthcheck

```http
  GET /api/v1/healthcheck
```
#### Create a new book

```http
  POST /api/v1/books
```

| Path parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to create |


| Body parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `title`      | `string` | **Required**. title of the book |
| `author`      | `string` | **Required**. author of the book |
| `year`      | `string` | **Required**. published year |
| `description`      | `string` | **Required**. smth about a book |
| `genres`      | `string[]` | **Required**. genres of the new book |


#### Get a book by id

```http
  GET /api/v1/books/${id}
```

| Path parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to fetch |


#### Update book by id

```http
  PATCH /api/v1/books/${id}
```

| Path parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to update |


| Body parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `title`      | `string` | title of the book |
| `author`      | `string` | author of the book |
| `year`      | `string` |  published year |
| `description`      | `string` |  smth about a book |
| `genres`      | `string[]` | genres of the new book |

#### Delete a book by id

```http
  DELETE /api/v1/books/${id}
```

| Path parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Id of item to delete |

## DB structure

```dbml
Table books {
  id bigserial [primary key]
  title text [not null]
  author text [not null]
  year int [not null]
  description text [not null]
  genres text[] [not null]
}

Table users {
  id bigserial [primary key]
  username varchar(50) [not null, unique]
  password text [not null]
}

// many-to-many
Table reviews {
  id bigserial [primary key]
  created_at timestamp [not null, default: `now()`],
  user_id bigserial [not null],
  book_id bigserial [not null],
  content text,
  rating integer,
}

Ref: reviews.book_id < books.id
Ref: reviews.user_id < users.id
```


## Authors

- Uldana Shyndali 22B030473 [@shyndaliu](https://www.github.com/shyndaliu)

