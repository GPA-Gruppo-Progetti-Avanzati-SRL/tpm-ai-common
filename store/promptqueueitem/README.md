# README

```json
[
  {
    "$match":
    {
      "status": {
        "$in": ["ready", "pippo"]
      }
    }
  },
  {
    "$group": {
      "_id": {
        "domain": "$domain",
        "site": "$site",
        "category": "$category",
        "status": "$status"
      },
      "count": {
        "$sum": 1
      },
      "total_weight": {
        "$sum": "$weight"
      }
    }
  },
  {
    "$project": {
      "_id": 0,
      "domain": "$_id.domain",
      "site": "$_id.site",
      "category": "$_id.category",
      "weight": "$total_weight",
      "count": "$count",
      "status": "$_id.status"
    }
  }
]
```