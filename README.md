# RabbitMQ for Event-Driven Architecture in Golang with Microservices

## RabbitMQ

- Sending data protocol: `AMQP`
- Expose ports: `5672` is AMQP port connection for RabbitMQ and `15672` is port used by Admin UI/Management UI for RabbitMQ.

---

Users:

RabbitMQ comes with a default guest user pre-installed.
In production we do not want to use this default guest user, first add new user, and remove guest pre-installed user.

Users are usually used to limit and manage what permissions your users has.

---

Virtual Hosts:

in rabbit you have resources like channels, exchanges, queues etc. this resources are contained in something called Virtual Hosts, Virtual Hosts is sort of a namespace, you use Virtual Hosts to kind of limit and separate resources in a logical way and it's called a Virtual because it's done in The Logical layer it's a soft restriction between what resources can reach what resource etc.

- in the Management UI in right top corner Virtual host: All and / is for global one.

Virtual Host used to group certain resource to grid together and restrict access on those Virtual Hosts.

Also we need permissions to `Virtual Host = customers` to communicate with resources inside Virtual Host.

---

- Queues
- Producers
- Exchanges
- Consumers

```text
(Anyone sending messages is)
Producer 1 -- Message 1
Producer 2 -- Message 1
--->
--->
(Decides where Messages Goes)
Exchange
--->
(A Queue is a Message Buffer)
Queue
Message 1, Message 2, Message 3
--->
(Anyone receiving messages)
Consumer
```

Producers is any piece of software that is sending messages, they produce messages.
Consumers who is any piece of software that is receiving those messages.
between them:
Exchanges - producers send messages to a exchange, is like a broker or router,
exchange knows which queues are bound to the exchange,
to bind something to an exchane we use something called a binding,
a binding is basically a rule or set of rules, if queue should receive messages.
Queues - is basically a buffer for messages, it's usually a `First In First Out` (FIFO) queue, the messages comes in and comes out in the correct order.

Exchange as bound to a queue by a set of rules so a producer sends a message say that the message is one the topic `Customers_Registred` now exchange will know if a certain queue is bound on that topic and send it further along,
it's real;y important to understand you don't send messages to the queue, you send messages to The Exchange which then routes the messages where they should go.

### RabbitMQ Docker commands log

#### Users

First way:

> `docker exec -it rabbitmq bash`

- add new user: `rabbitmqctl add_user <newusername> <secretpassword>` (in my case just username: `admin` password: `admin`)
- add permissions for new added user named admin: `rabbitmqctl set_user_tags admin administrator`
- remove default pre-installed `guest` user: `rabbitmqctl delete_user guest`

or Second way:

`docker exec -it rabbitmq rabbitmqctl <...command>`

#### Virtual Hosts

> `docker exec -it rabbitmq bash`

- add new virtual host: `rabbitmqctl add_vhost customers`
- add permissions to `user` (admin) to communication with this Virtual host: `rabbitmqctl set_permissions -p customers admin "^customers.*" ".*" ".*"` (structure: `rabbitmqctl set_permissions -p <vhost_name> <user> <configurations_vhost> <write_regex> <read_regex>`) (configurations, write, read) (using regex pattern `"^customer.*"` for customer virtual host, if we want for all virtual hosts then `".*"`)

## Event-Driven Architecture

ED-A is used for communication of Microservices with events that will be sended and received in two side of microservices between them will be like Exchange and it will be Asynchronously.

## Docker

### docker run of RabbitMQ Container instance on Windows

How to run Docker RabbitMQ container instance:

- basic: `docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management`
- with persistent volume: `docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 -v rabbitmq_data:/data rabbitmq:3-management`

Go to Management UI:

- `http://127.0.0.1:15672/`
- user: `guest`, password: `guest` (for my env: user: `admin`, password: `admin`)

---

There was an issue with accesing docker run rabbitmq management interface on Windows.
This issue was caused by Windows Firewall (`docker run` and `docker compose` act differently! `docker compose` uses defauly `bridge` network!).

- in powershell with administrator (it will add new inbound firewall rule for rabbitmq ports): `New-NetFirewallRule -DisplayName "Allow RabbitMQ Ports" -Direction Inbound -Protocol TCP -LocalPort 5672,15672 -Action Allow`

To stop and remove RabbitMQ docker run container instance:

- `docker stop rabbitmq`
- `docker rm rabbitmq`

To view logs of RabbitMQ docker run container instance:

- `docker logs -f rabbitmq`

To access the shell of RabbitMQ container:

- `docker exec -it rabbitmq bash`
- in bash (to check status of rabbitmq service): `rabbitmqctl status`

---

#### Setting Up a RabbitMQ Cluster in Docker

1. docker network: `docker network create rabbitmq_cluster`
2. rabbitmq instances on network: `docker run -d --name rabbitmq1 --hostname rabbitmq1 --network rabbitmq_cluster rabbitmq:3-management`, `docker run -d --name rabbitmq2 --hostname rabbitmq2 --network rabbitmq_cluster rabbitmq:3-management`
3. cluster the nodes: (access to the shell) `docker exec -it rabbitmq1 bash` and then (cluster the nodes using `rabbitmqctl`) `rabbitmqctl stop_app`, `rabbitmqctl reset`, `rabbitmqctl join_cluster rabbit@rabbitmq2`, `rabbitmqctl start_app`
4. verify cluster status: `rabbitmqctl cluster_status`
5. persisting data (to persist data, you can mount a volume to the RabbitMQ container): `docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 -v /path/to/data:/var/lib/rabbitmq rabbitmq:3-management`

---

#### Other docker commands

- `docker ps` check runned containers
- `docker ps -a` check stopped containers
- `docker network --help`
- `docker network ls`
- `docker network inspect bridge`
- `docker run -d --hostname rabbitmq --name rabbitmq rabbitmq:3-management` (provided `hostname` and with skipped port mapping)
- `docker network inspect rabbitmq`
- `docker system prune -f` (clean up unused docker resources)
- `docker network prune -f` (clean up unused docker resources)
- `docker inspect rabbitmq`
- `docker container inspect rabbitmq`

---

#### Dockerfile

```yml
FROM rabbitmq:3-management
RUN rabbitmq-plugins enable --offline rabbitmq_mqtt rabbitmq_federation_management rabbitmq_stomp
WORKDIR /usr/src/app
ENV RABBITMQ_ERLANG_COOKIE: 'secret cookie here'
VOLUME ~/.docker-conf/rabbitmq/data/:/var/lib/rabbitmq/mnesia/
EXPOSE 5672 15672
```

#### Makefile

```makefile
# Run RabbitMQ container instance
# with persistence volume
rabbit-build:
	@docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 -v rabbitmq_data:/data rabbitmq:3-management

# Start the Existing container
# docker ps -a
rabbit-run:
	@docker start rabbitmq

# Stop container
rabbit-down:
	@docker stop rabbitmq

# Remove container
rabbit-rm:
	@docker rm rabbitmq

# Check logs of container
rabbit-logs:
	@docker logs -f rabbitmq

# Access to shell container:
rabbit-exec:
	@docker exec -it rabbitmq bash
```

### docker compose

```yml
version: "3"
services:
 rabbitmq:
    image: "rabbitmq:3-management"
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - 'rabbitmq_data:/data'

volumes:
  rabbitmq_data:
```

- `docker compose up --build -d` (with build)
- `docker compose up -d` (without)
- `docker compose logs -f` (logs)
- `docker compose down` (down)
- `docker compose down -v` (remove db volume)

## Golang

- in powershell: `mkdir cmd/producer`, `mkdir internal`, `ni cmd\producer\main.go`, `ni internal/rabbitmq.go` (`ni` is alias for `New-Item`) (in bash `touch .\cmd\producer\main.go`)
- `go mod init github.com/dotpep/golang-event-driven-rabbitmq`
- `go get github.com/rabbitmq/amqp091-go`

## Resources/Links

- [Running RabbitMQ in Docker: A Comprehensive Guide](https://www.svix.com/resources/guides/rabbitmq-docker-setup-guide/#:~:text=Step-by-Step%20Guide%20with%20Code%20Samples%201%20Step%201%3A,network%3A%20...%205%20Step%205%3A%20Persisting%20Data%20)
- [How to open rabbitmq in browser using docker container?](https://stackoverflow.com/questions/47290108/how-to-open-rabbitmq-in-browser-using-docker-container#:~:text=Please%20you%20need%20to%20enable%20the%20management%20plugins%2C,go%20to%20http%3A%2F%2Flocalhost%3A8085%2F%2C%20to%20access%20the%20management%20console.) (It can be because of Firewall, that was in my case issue to accessing rabbitmq management interface)
