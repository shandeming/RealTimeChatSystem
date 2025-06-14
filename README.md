# RealTimeChatSystem

This is a real-time chat system implemented in Go.

## Features

- User online and offline notifications
- Public chat mode: All online users can send and receive messages.
- Private chat mode: Users can choose a specific online user for one-on-one chat.
- Modify username: Users can change their display name.
- Query online users: Users can view a list of all currently online users.
- Timeout handling: Users who are inactive for an extended period will be disconnected by the server.

## Feature Demonstration

### Start Server

![start server](images/image-4.png)

### Client Chat

![client chat](images/image-2.png)
![client chat](images/image-3.png)

## Project Structure

```
.
├── client
├── client.go     # Client main program
└── serverp/
    └── user.go       # User-related logic
    ├── main.go       # Server main program entry point
    ├── server
    ├── server.go     # Server core logic
```

## How to Run

### Server

```bash
./server
```

### Client

```bash
./client
```

Or run with parameters:

```bash
./client -Ip <Server IP Address> -Port <Server Port Number>
```

Connects to `127.0.0.1:8888` by default.

## Usage Instructions

After the client starts, the following menu will be displayed:

```
1. Public Chat Mode
2. Private Chat Mode
3. Update Username
0. Exit
```

Enter the number corresponding to the desired function according to the prompts.

- **Public Chat Mode**: Enter a message and press Enter to send; visible to all online users. Enter `exit` to leave public chat mode.
- **Private Chat Mode**:
  1. A list of currently online users will be displayed first.
  2. Enter the target user's name.
  3. Enter the message to be sent.
  4. Enter `exit` to leave the current private chat or return to the main menu.
- **Update Username**: Enter the new username.
- **Exit**: Closes the client.

## Notes

- The server listens on `127.0.0.1:8888` by default.
- The client and server need to be in the same network environment, or the client must be able to access the IP address and port where the server is running.
