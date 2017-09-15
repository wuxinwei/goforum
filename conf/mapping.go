package conf

const PostMapping = `
{
  "settings": {
    "index": {
      "number_of_shards": 3,
      "number_of_replicas": 2
    }
  },
  "mappings": {
    "user": {
      "properties": {
        "user_id": {
          "type": "integer"
        },
        "created_at": {
          "type": "date"
        },
        "updated_at": {
          "type": "date"
        },
        "deleted_at": {
          "type": "date"
        },
        "is_deleted": {
          "type": "boolean",
          "store": false
        },
        "username": {
          "type": "keyword"
        },
        "name": {
          "type": "text"
        },
        "gender": {
          "type": "keyword"
        },
        "birthday": {
          "type": "date"
        },
        "country": {
          "type": "keyword"
        },
        "province": {
          "type": "keyword"
        },
        "city": {
          "type": "keyword"
        },
        "coin": {
          "type": "integer"
        }
      }
    },
    "post": {
      "properties": {
        "post_id": {
          "type": "integer"
        },
        "created_at": {
          "type": "date"
        },
        "updated_at": {
          "type": "date"
        },
        "deleted_at": {
          "type": "date"
        },
        "is_deleted": {
          "type": "boolean",
          "store": false
        },
        "user_id": {
          "type": "integer"
        },
        "title": {
          "type": "text"
        },
        "author": {
          "type": "keyword"
        },
        "content": {
          "type": "text"
        }
      }
    },
    "comment": {
      "properties": {
        "comment_id": {
          "type": "integer"
        },
        "created_at": {
          "type": "date"
        },
        "updated_at": {
          "type": "date"
        },
        "deleted_at": {
          "type": "date"
        },
        "is_deleted": {
          "type": "boolean",
          "store": false
        },
        "user_id": {
          "type": "integer"
        },
        "post_id": {
          "type": "integer"
        },
        "target_id": {
          "type": "integer"
        },
        "content": {
          "type": "text"
        }
      }
    }
  }
}
`
