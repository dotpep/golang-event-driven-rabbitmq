# RabbitMQ for Event-Driven Architecture in Golang with Microservices

## Issues

> Windows 10 - Enterprice Environment (hard settings, security rules)

- run as administrator docker desktop
- disable `public network` firewall rules, for each network or add inbound rule (also try to disable `private network`, `domain network`)

Important Note (latest try, temporary solution): try to just open settings of firewall called `Firewall & network protection` using Serach bar, and try to `on/off`, `Public network` or others like (Private network and Domain network) if doesn't works (in my Windows 10 Enterprice/Pro with special Server rule settings)

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

---

You should recreate the channel for each concurrent task,
but reuse the connection, always have one connection for your particular service,
and spawn channels from that.
Reason why we want to do that is because if you spawn connections instead,
you will create so many TCP connections and that does not scale very well.

---

Sending data:

- Creating Queue

```golang
type RabbitClient struct {
    conn *amqp.Connection
    ch *amqp.Channel
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autoDelete bool) error {
    _, err := rc.ch.QueueDeclare(
        queueName, durable, autoDelete,
        false, false, nil,
    )
    return err
}
```

- `queueName` string is queue name like `customers_created` or `customers_test`
- `durable` bool is for making queue durable and persistence/persisted (for survive messages between restarts or crashes) (whenever the broker restarts it will be saved/persisted Queue of RabbitMQ)
- `autoDelete` bool is for automatically deleting (if Queue which is set to automatically delete will be deleted whenever the software that created it shuts down so whenever the producer shuts down after the 10 second timeout it will delete the queue, this is very common when you have Dynamic queues being created by service, that you don't know to maybe they respond differently and you don't want to clog everything here with a million queues, you have to delete them so Auto delete is good for that)

---

You're not sending messages on queue we are sending messages to exchanges.

Exchanges.

It is just router but there is a few different exchanges:

- Direct Exchange (Exact Key match) just directly sends messages throught Routing Key
for example Producer:`customer_created` --> Excahnge:`customer_events` --> Queue:`customer_created` --> Consumer:messagge Received | Queue:`customer_emailed` --> Consumer:messagge Not Received and Queue do not have message
- Fanout Exchange (Ignores Routing Key) and sends messages to all Queues, no matter what
for example Producer:`customer_created` --> Excahnge:`customer_events` --> Queue:`customer_created` --> Consumer:messagge Received | Queue:`customer_emailed` --> Consumer:messagge Received
- Topic Exchange (Rules on Routing Key delimited by ./dot)
for example Producer:`customers.created.february` --> Excahnge:`customer_events` --> Rule:`customers.created.#` or `customers.*.february` (like `customers.deleted.february`) --> Queue:`customer_created` --> Consumer:messagge Received | Rule:`customers.created.march` --> Queue:`customer_emailed` --> Consumer:messagge Not Received and message do not pass to second Queue because of Rule dismatching.
Very Dynamic routes for some specific Rules and Topics.
- Header Exchange - (Rules based on Extra Header) basically key value fields, routing based on header
for example Producer:`browser = Linux` --> Excahnge:`customer_events` --> Rule:`browser = Linux` --> Queue:`customer_linux` --> Consumer:messagge Received | Rule:`browser = Windows` --> Queue:`customer_windows` --> Consumer:messagge Not Received and message do not pass to second Queue because of Rule dismatching.

To start receiving or sending messages on a Queue,
you need to bind that Queue to an Exchange this is called The Binding.

Binding is basically a routeing route.

Queue can be bound to multiple Exchanes, also you can have Exchanes being bound to Exchanges.

Whenever you send a Messages on Message Queue you have to add a Routing key,
and the Routing key is sometimes referred to as The Topic will be used by The Exchange.

---

Creating Exchange - to create Exchange we can use RabbitMQ Admin Command Line tool `rabbitmqadmin` instead `rabbitmqctl`, also you can create Queues inside of code, it depends on what you like to do.

Our admin user doesn't have permissions to send data to the cumtomers `rabbitmqctl set_topic_permissions`.

### RabbitMQ Docker commands log

#### Users (by `rabbitmqctl`)

First way:

> `docker exec -it rabbitmq bash`

- add new user: `rabbitmqctl add_user <newusername> <secretpassword>` (in my case just username: `admin` password: `admin`)
- add permissions for new added user named admin: `rabbitmqctl set_user_tags admin administrator`
- remove default pre-installed `guest` user: `rabbitmqctl delete_user guest`

or Second way:

`docker exec -it rabbitmq rabbitmqctl <...command>`

#### Virtual Hosts (by `rabbitmqctl`)

> `docker exec -it rabbitmq bash`

- add new virtual host: `rabbitmqctl add_vhost customers`
- add permissions to `user` (admin) to communication with this Virtual host: `rabbitmqctl set_permissions -p customers admin ".*" ".*" ".*"` (structure: `rabbitmqctl set_permissions -p <vhost_name> <user> <configurations_vhost> <write_regex> <read_regex>`) (configurations, write, read) (using regex pattern `"^customer.*"` for customer virtual host, if we want for all virtual hosts then `".*"`)

#### Exchange (using rabbitmq cmd `rabbitmqadmin`)

Declaring new Exchange named customer_events, for Virtual Host of Customer, type of Topic Exchange, for Admin user `rabbitmqadmin declare exchange --vhost=<vhost_name> name=<exchange_name> type=<exchange_type> durable=<durable_bool_param> -u <user> -p <user_password>`

- `docker exec -it rabbitmq rabbitmqadmin declare exchange --vhost=customers name=customer_events type=topic durable=true -u admin -p admin`

Set Topic permissions for admin user, for access of sending data to the customer. `rabbitmqctl set_topic_permissions -p <vhost_name> <user_name> <exchange_name> <write_permission_regex> <read_permission_regex>`

- `rabbitmqctl set_topic_permissions -p customers admin customer_event ".*" ".*"`

#### RabbitMQ Commands

> `docker exec -it rabbitmq bash`

- `rabbitmqctl add_user admin admin_password`
- `rabbitmqctl set_user_tags admin administrator`
- `rabbitmqctl delete_user guest`

---

- `rabbitmqctl add_vhost customers`
- `rabbitmqctl set_permissions -p customers admin ".*" ".*" ".*"`
- `rabbitmqctl list_vhosts`
- `rabbitmqctl list_permissions -p customers`
- `rabbitmqctl list_user_permissions admin`

---

- `rabbitmqadmin declare exchange --vhost=customers name=customer_events type=topic durable=true -u admin -p admin`
- `rabbitmqadmin delete exchange name='customer_events' --vhost=customers -u admin -p admin`
- `rabbitmqctl list_exchanges --vhost=customers`

---

- `rabbitmqctl list_exchanges`

---

- `tail -f /var/log/rabbitmq/rabbit@<hostname>.log` Debug Logs
- `systemctl restart rabbitmq-server` Restart RabbitMQ

## Event-Driven Architecture

ED-A is used for communication of Microservices with events that will be sended and received in two side of microservices between them will be like Exchange and it will be Asynchronously.

## Makefile Step-By-Step Instruction

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

- in powershell: `mkdir cmd/producer`, `mkdir internal`, `ni cmd\producer\main.go`, `mkdir cmd/consumer`, `ni cmd\consumer\main.go`, `ni internal/rabbitmq.go`(`ni` is alias for `New-Item`) (in bash `touch .\cmd\producer\main.go`)
- `go mod init github.com/dotpep/golang-event-driven-rabbitmq`
- `go get github.com/rabbitmq/amqp091-go`

## Resources/Links

- [Running RabbitMQ in Docker: A Comprehensive Guide](https://www.svix.com/resources/guides/rabbitmq-docker-setup-guide/#:~:text=Step-by-Step%20Guide%20with%20Code%20Samples%201%20Step%201%3A,network%3A%20...%205%20Step%205%3A%20Persisting%20Data%20)
- [How to open rabbitmq in browser using docker container?](https://stackoverflow.com/questions/47290108/how-to-open-rabbitmq-in-browser-using-docker-container#:~:text=Please%20you%20need%20to%20enable%20the%20management%20plugins%2C,go%20to%20http%3A%2F%2Flocalhost%3A8085%2F%2C%20to%20access%20the%20management%20console.) (It can be because of Firewall, that was in my case issue to accessing rabbitmq management interface)
- [Learn RabbitMQ for Event-Driven Architecture (EDA), Golang - Percy](https://programmingpercy.tech/blog/event-driven-architecture-using-rabbitmq/)
