# buffered_request

Run app:
```
docker compose -f deployment/docker-compose.yml up 
```

Bufferd service wait a request on localhost:9000/.
Request should by a array of json objects like
```json
[
    {
        "period_start": "2024-05-01",
        "period_end": "2024-05-31",
        "period_key": "month",
        "indicator_to_mo_id": "227373",
        "indicator_to_mo_fact_id": "0",
        "value": "1",
        "fact_time": "2024-05-31",
        "is_plane": "0",
        "auth_user_id": "40",
        "comment": "buffer Last_name"
    }
]
```
after that service async write this request to kafka:9092. And 
consumer service wait for message in queue. Next step is to
send resquet to another API service.

