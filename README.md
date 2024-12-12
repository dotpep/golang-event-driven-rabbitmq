# RabbitMQ for Event-Driven Architecture in Golang with Microservices

## RabbitMQ

- Sending data protocol: `AMQP`
- Expose ports: `5672` is AMQP port connection for RabbitMQ and `15672` is port used by Admin UI/Management UI for RabbitMQ.

---

RabbitMQ comes with a default guest user pre-installed.
In production we do not want to use this default guest user, first add new user, and remove guest pre-installed user.

## Event-Driven Architecture

ED-A is used for communication of Microservices with events that will be sended and received in two side of microservices between them will be like Exchange and it will be Asynchronously.

## Docker

### docker run of RabbitMQ Container instance on Windows

How to run Docker RabbitMQ container instance:

- basic: `docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management`
- with persistent volume: `docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 -v rabbitmq_data:/data rabbitmq:3-management`

Go to Management UI:

- `http://127.0.0.1:15672/`
- user: `guest`, password: `guest`

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

## Resources/Links

- [Running RabbitMQ in Docker: A Comprehensive Guide](https://www.svix.com/resources/guides/rabbitmq-docker-setup-guide/#:~:text=Step-by-Step%20Guide%20with%20Code%20Samples%201%20Step%201%3A,network%3A%20...%205%20Step%205%3A%20Persisting%20Data%20)
- [How to open rabbitmq in browser using docker container?](https://stackoverflow.com/questions/47290108/how-to-open-rabbitmq-in-browser-using-docker-container#:~:text=Please%20you%20need%20to%20enable%20the%20management%20plugins%2C,go%20to%20http%3A%2F%2Flocalhost%3A8085%2F%2C%20to%20access%20the%20management%20console.) (It can be because of Firewall, that was in my case issue to accessing rabbitmq management interface)
