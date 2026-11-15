# WASAText

**Full-stack messaging web app: a Go REST backend specified with OpenAPI 3, SQLite persistence and a Vue 3 single-page frontend, deployed with Docker Compose.**

![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)
![Vue.js](https://img.shields.io/badge/Vue.js%203-4FC08D?logo=vuedotjs&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-003B57?logo=sqlite&logoColor=white)
![OpenAPI](https://img.shields.io/badge/OpenAPI-3.0-6BA539?logo=openapiinitiative&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white)

Project for the _Web and Software Architecture_ course (B.Sc. in Computer Science, Sapienza University of Rome, a.y. 2024/2025).

---

## Features

- **Login** by username (accounts are created on first login)
- **Profile:** change username and profile photo
- **One-to-one chats and groups:** create groups, add members, leave, rename, set the group photo
- **Messages:** text and photos, **reply**, **forward**, delete your own messages
- **Reactions:** comment a message with an emoji, and remove it
- **Delivery status** with "received" and "read" checkmarks, tracked per recipient and aggregated per message
- **User search** to start new conversations

## Architecture

- **API design first:** the whole REST interface is specified in [`doc/api.yaml`](doc/api.yaml) (OpenAPI 3.0.3), with request/response schemas, parameters and error codes for about 20 endpoints.
- **Backend (Go):** `service/api` has one handler per endpoint with input validation (names, Base64 images) and authorisation checks (for example, only conversation members can read or post). `service/database` contains all SQL access, behind an `AppDatabase` interface.
- **Database:** six tables (users, conversations, members, messages, comments, per-recipient message status) created on start-up.
- **Frontend (Vue 3 + Vue Router):** views for login, home, conversation, new chat, new group and profile, built from reusable components (message list, input box, forward/reply popups, group info overlay).

> **Note on authentication:** as required by the course specification, login is _simplified_: the API returns the user identifier, which the client sends back as `Authorization: Bearer <id>`. This is not a secure authentication scheme and is not meant for real use.

### Main endpoints

| Area              | Endpoints                                                                                                                                                              |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Session & profile | `POST /session`, `GET /profile`, `PUT /profile/username`, `PUT /profile/photo`, `GET /users`                                                                           |
| Conversations     | `GET /conversations`, `POST /conversations/create`, `GET /privateChat/{userId}`, `GET /conversation-info/{id}`, `DELETE /conversation-delete/{id}`                     |
| Messages          | `POST /conversation-messages/{id}`, `DELETE /conversations/{id}/messages/{messageId}`, `POST /messages/{id}/forwards`, `POST /messages/{id}/replies`                   |
| Reactions         | `POST /messages/{id}/comments`, `DELETE /messages/{id}/comments`                                                                                                       |
| Groups            | `GET /conversation-members/{id}`, `POST /conversation-add/{id}`, `DELETE /conversation-members/{id}/me`, `PUT /conversation-name/{id}`, `PUT /conversation-photo/{id}` |

Full details are in the [OpenAPI specification](doc/api.yaml).

---

## Running it

### With Docker Compose (recommended)

```bash
git clone https://github.com/GiovanniCiarra/WASA-Project-2425.git
cd WASA-Project-2425
docker compose up --build
```

- Frontend: <http://localhost>
- Backend API: <http://localhost:3000>

Both images use multi-stage builds: Go is compiled in a `golang` image and copied into a slim Debian image; the Vue app is built with Node and served by `nginx`.

### For development

```bash
# backend (Go 1.19+)
go run ./cmd/webapi

# frontend, in another terminal
cd webui
corepack enable
yarn run dev
```

The frontend expects the API at `http://localhost:3000`. Go dependencies are vendored in `vendor/` and frontend packages are kept in the Yarn PnP cache, as required by the course template.

---

## Repository structure

```
cmd/webapi/          entry point, configuration, CORS, graceful shutdown
cmd/healthcheck/     container health check
service/api/         HTTP handlers (one file per endpoint) and request context
service/database/    SQLite schema and queries
doc/api.yaml         OpenAPI 3 specification
webui/               Vue 3 + Vite frontend
Dockerfile.backend   multi-stage build for the Go API
Dockerfile.frontend  multi-stage build for the Vue app (nginx)
docker-compose.yml
```

## Credits

The project skeleton (build setup, configuration loading, vendoring) comes from the course template _Fantastic Coffee (decaffeinated)_ by Enrico Bassetti, released under the MIT License. The application (API design, backend logic, database and frontend) was developed by me.
