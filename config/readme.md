# Configuration files for chatbot server
## Server
* is_production: set to true if the server is deployed on real server with https enabled, false otherwise
* port: port number of the server
* cert_file: location of certificate file
* key_file: location of key file

## Google Custom Search JSON API
### Prerequisites
* Please visit https://developers.google.com/custom-search/json-api/v1/overview
* Go to Prerequisites section and follow the steps to create new Custom Search Engine
* After that get the "Search Engine ID" for search_engine_id field
* Generate an API Key for api_key field