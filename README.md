## Запуск
```
    docker compose up postgres server
```

## CURL

POST
```
curl --location 'http://localhost:8080/banner' \
--header 'accept: application/json' \
--header 'token: admin' \
--header 'Content-Type: application/json' \
--data '{
  "tag_ids": [
    6, 66, 666
  ],
  "feature_id": 17,
  "content": {
    "title": "banner_title",
    "text": "banner_description",
    "url": "some_url"
  },
  "is_active": true
}'
```

PATCH
```
curl --location --request PATCH 'http://localhost:8080/banner/16' \
--header 'accept: */*' \
--header 'token: admin' \
--header 'Content-Type: application/json' \
--data '{
  "tag_ids": [
    6, 66, 666
  ],
  "feature_id": 17,
  "content": {
    "title": "updated_title",
    "text": "updated_text",
    "url": "updated_url"
  },
  "is_active": true
}'
```

GET
```
curl --location 'http://localhost:8080/user_banner?tag_id=666&feature_id=17&use_last_revision=true' \
--header 'accept: application/json' \
--header 'token: admin'
```

```
"{\"text\":\"updated_text\",\"title\":\"updated_title\",\"url\":\"updated_url\"}\n"
```

GET
```
curl --location 'http://localhost:8080/banner?feature_id=17&tag_id=66&limit=10&offset=0' \
--header 'accept: application/json' \
--header 'token: user'
```

```
[
    {
        "id": 16,
        "tags": [
            6,
            66,
            666
        ],
        "feature_id": 17,
        "content": "{\"text\":\"updated_text\",\"ttle\":\"updated_title\",\"url\":\"updated_url\"}\n",
        "is_active": true,
        "created_at": "2024-04-14T20:13:54.828126Z",
        "updated_at": "2024-04-14T20:17:33.383296Z"
    }
]
```

DELETE
```
curl --location --request DELETE 'http://localhost:8080/banner/16' \
--header 'accept: */*' \
--header 'token: admin'
```
