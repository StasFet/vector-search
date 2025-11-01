# Vector Search
This is a little project where I will experiment using MongoDB's Vector Search. 

For now, it has no purpose other than to experiment with MongoDB and vector search, but perhaps I will find an interesting use.

## API Endpoint
There is also an API endpoint for interfacing with this project. Here are the routes:

- `POST /api/:database/:collection/vectorstore/` - Accepts a JSON body with a `text` field, inserts that text into the DB. Automatically makes vector index.
- `POST /api/:database/:collection/vectorsearch/` - Accepts a JSOn body with a `text` and an `amount` field, responds with the `amount` closest matches. 
- `GET /api/:database/:collection/vectorstore/` - Returns all the elements in the DB (text only)

Note: Each operation is done on the `:collection` collection of the `:database` database.