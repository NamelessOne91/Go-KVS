# Go-KVS

An implementation of a distributed key-value store for study purposes.

## REST API endpoints - V1

API root: <strong>/v1<strong>

| Functionality | Method | Endpoint | Status codes
| :---|:--|:--:|:--:|
| Put a key-value in the store | PUT | /{key}  | 201 (created) |
| Read a key-value from the store | GET |  /{key}  | 200 (OK), 404 (Not Found)|
| Delete a key-value pair | DELETE |  /{key} | 200 (OK)