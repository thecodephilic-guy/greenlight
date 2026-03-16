# Greenlight API Server
API Server for movies

## API Reference

The Greenlight API accepts and returns JSON payloads. Base URL: `https://api.greenlight.sohail.world`

### System
| Method | Endpoint | Description | Auth Required | Example Usage |
| :--- | :--- | :--- | :---: | :--- |
| `GET` | `/v1/healthcheck` | Returns the system status, environment, and version. | No | `curl -i base_url/v1/healthcheck` |

### Users & Authentication
| Method | Endpoint | Description | Auth Required | Example Payload |
| :--- | :--- | :--- | :---: | :--- |
| `POST` | `/v1/users` | Registers a new user. Triggers an activation email. | No | `{"name": "Sohail", "email": "sohail@example.com", "password": "securepassword123"}` |
| `PUT` | `/v1/users/activated` | Activates a user account using the token sent via email. | No | `{"token": "YOUR_ACTIVATION_TOKEN_HERE"}` |
| `POST` | `/v1/tokens/authentication` | Authenticates a user and returns a 24-hour Bearer token. | No | `{"email": "sohail@example.com", "password": "securepassword123"}` |

### Movies
| Method | Endpoint | Description | Auth Required | Example Payload / Usage |
| :--- | :--- | :--- | :---: | :--- |
| `GET` | `/v1/movies` | Lists all movies. Supports pagination, sorting, and filtering. | Yes | `GET /v1/movies?title=godfather&genres=crime,drama&page=1&sort=-year` |
| `POST` | `/v1/movies` | Creates a new movie record in the database. | Yes (Admin) | `{"title": "Moana", "year": 2016, "runtime": "107 mins", "genres": ["animation", "adventure"]}` |
| `GET` | `/v1/movies/:id` | Retrieves the details of a specific movie by its ID. | Yes | `curl -i base_url/v1/movies/1` |
| `PATCH` | `/v1/movies/:id` | Partially updates a specific movie. | Yes (Admin) | `{"year": 2017, "runtime": "115 mins"}` |
| `DELETE` | `/v1/movies/:id` | Deletes a specific movie from the database. | Yes (Admin) | `curl -X DELETE base_url/v1/movies/1` |

> **Note on Authentication:** For routes requiring authentication, you must include the token received from `/v1/tokens/authentication` in the header of your request like so:
> `Authorization: Bearer <your-authentication-token>`
