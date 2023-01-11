# comchat

Simple CLI chat app using gRPC

### Servers

* Go

### Clients

* Go
* Node.js

### Prerequisites

* Go
* Node.js (>= 17.0.0)

### Usage

1. Start server

    ```
    $ git clone https://github.com/CRaLFa/comchat.git
    $ cd comchat/go
    $ go run server/server.go
    ```
1. Start client

    * Go

    ```
    $ cd comchat/go
    $ go run client/client.go
    Enter your name: <name>
    ```

    * Node.js

    ```
    $ cd comchat/node
    $ npm install
    $ npm start
    Enter your name: <name>
    ```

1. Stop client

    Type `exit` and press [Enter].

    ```
    Enter your name: hoge
    2022/11/06 17:44:10 [SYSTEM] : hoge has entered.
     ︙
    exit
    ```
