# Go Chat Application

A simple chat application built in Go with a **WebSocket server** and **CLI client**.  
The project supports direct messaging between users, password-protected users, and basic command handling.

---

## 📂 Project Structure

```
cmd/
  client/
     main.go        # Client entry point
  server/
     main.go        # Server entry point
internal/
   channel/         # Channel-related logic
   client/          # Client-side logic
   dto/             # Data transfer objects
   server/          # Server-side logic
   user/            # User model and helpers
.env                # Environment variables (server/client configs)
.gitignore
client.sh           # Build script for client
server.sh           # Build script for server
go.mod
go.sum
README.md
```

---

## ⚙️ Setup

1. **Clone the repo**

   ```sh
   git clone https://github.com/yourusername/go-chat.git
   cd go-chat
   ```

2. **Configure environment**
   Create a `.env` file in the root:

   ```env
   SERVER_URL=ws://localhost:8080
   HTTP_URL=http://localhost:8080
   ```

3. **Install dependencies**
   ```sh
   go mod tidy
   ```

---

## 🚀 Running the App

### Start Server

```sh
./server.sh
```

or manually:

```sh
go run cmd/server/main.go
```

### Start Client

```sh
./client.sh
```

or manually:

```sh
go run cmd/client/main.go
```

---

## 💬 Usage (Client Commands)

Once connected, you can use these commands:

- **Send a direct message**

  ```
  /msg <username> <message>
  ```

  Example:

  ```
  /msg alice Hello Alice!
  ```

- **Create a channel**

  ```
  /create <channel_name>
  ```

- **Join a channel**

  ```
  /join <channel_name>
  ```

- **Password-protected users**
  If you send a message to a user with a password, you’ll be prompted:
  ```
  Password of user:
  ```
  Then type the password (input hidden).

---

## 🛠 Development

- Client and Server share code via the `internal/` folder.
- Scripts (`client.sh`, `server.sh`) simplify building binaries.
- To add new commands, extend `internal/dto` and update handlers in `internal/server`.

---

## 📦 Building Binaries

### Build Client

```sh
./client.sh
```

Binary will be in `./bin/client`

### Build Server

```sh
./server.sh
```

Binary will be in `./bin/server`

---

## 📜 License

MIT
