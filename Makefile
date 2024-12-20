# ---
# Golang Applications Management
producer:
	@echo "Running Producer..."
	@go run ./cmd/producer/main.go

consumer:
	@echo "Running Consumer..."
	@go run ./cmd/consumer/main.go


# ---
# RabbitMQ Docker Containers Management
RABBIT_CONTAINER_NAME = rabbitmq
RABBIT_IMAGE = rabbitmq:3-management
RABBIT_PORTS = -p 5672:5672 -p 15672:15672
RABBIT_VOLUME = -v rabbitmq_data:/var/lib/rabbitmq

# Run a new RabbitMQ container with persistence
rabbit-build:
	@echo "Starting a new RabbitMQ container with persistence..."
	@docker run -d --name $(RABBIT_CONTAINER_NAME) $(RABBIT_PORTS) $(RABBIT_VOLUME) $(RABBIT_IMAGE)

# Start an existing RabbitMQ container
rabbit-run:
	@echo "Starting RabbitMQ container..."
	@docker start $(RABBIT_CONTAINER_NAME)

# Stop the RabbitMQ container
rabbit-down:
	@echo "Stopping RabbitMQ container..."
	@docker stop $(RABBIT_CONTAINER_NAME)

# Restart the RabbitMQ container
rabbit-restart:
	@echo "Restarting RabbitMQ container..."
	@docker restart $(RABBIT_CONTAINER_NAME)

# Remove the RabbitMQ container
rabbit-rm:
	@echo "Removing RabbitMQ container..."
	@docker rm $(RABBIT_CONTAINER_NAME)

# Check logs of the RabbitMQ container
rabbit-logs:
	@echo "Showing logs of RabbitMQ container..."
	@docker logs -f $(RABBIT_CONTAINER_NAME)

# Access the RabbitMQ container shell
rabbit-exec:
	@echo "Accessing RabbitMQ container shell..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) bash


# ---
# RabbitMQ Docker Container Environments Setup Management

# Usage:
# 1. make rabbit-ctl-new_admin_user
# 2. make rabbit-ctl-new_customers_vhost
# 3. make rabbit-admin-new_customer_events_exchange
# 4. make rabbit-info

# RabbitMQ Admin credentials
RABBIT_ADMIN_USER = admin
RABBIT_ADMIN_PASS = admin

# RabbitMQ Virtual Host
VHOST_CUSTOMERS = customers

# Exchange
EXCHANGE_NAME = customer_events
EXCHANGE_TYPE = topic

# Adding new Admin user:
# - Create a new user with administrator privileges.
# - Delete the default 'guest' user for security.
rabbit-ctl-new_admin_user:
	@echo "Adding new admin user..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl add_user $(RABBIT_ADMIN_USER) $(RABBIT_ADMIN_PASS)
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl set_user_tags $(RABBIT_ADMIN_USER) administrator
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl delete_user guest
	@echo "Listing current users..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl list_users

# Adding new Virtual Host named 'customers' and setting permissions:
# - Create the virtual host.
# - Set permissions for the admin user to have full access.
rabbit-ctl-new_customers_vhost:
	@echo "Creating new virtual host '$(VHOST_CUSTOMERS)'..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl add_vhost $(VHOST_CUSTOMERS)
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl set_permissions -p $(VHOST_CUSTOMERS) $(RABBIT_ADMIN_USER) ".*" ".*" ".*"
	@echo "Listing permissions for '$(VHOST_CUSTOMERS)'..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl list_permissions -p $(VHOST_CUSTOMERS)

# Declaring a new exchange named 'customer_events' in the 'customers' virtual host:
# - Declare the exchange type of topic.
# - Set exchange topic permissions for admin.
rabbit-admin-new_customer_events_exchange:
	@echo "Declaring new exchange '$(EXCHANGE_NAME)' in virtual host '$(VHOST_CUSTOMERS)'..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqadmin declare exchange \
		--vhost=$(VHOST_CUSTOMERS) \
		--name=$(EXCHANGE_NAME) \
		--type=$(EXCHANGE_TYPE) \
		--user=$(RABBIT_ADMIN_USER) \
		--password=$(RABBIT_ADMIN_USER) \
		--durable=true
	@echo "Setting topic permissions for exchange '$(EXCHANGE_NAME)'..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl set_topic_permissions \
		-p $(VHOST_CUSTOMERS) \
		$(RABBIT_ADMIN_USER) \
		$(EXCHANGE_NAME) \
		".*" ".*"
	@echo "Listing exchanges in virtual host '$(VHOST_CUSTOMERS)'..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl list_exchanges -p $(VHOST_CUSTOMERS)

# Check RabbitMQ configuration and resources:
rabbit-info:
	@echo "Listing all virtual hosts..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl list_vhosts
	@echo "Listing users and their permissions..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl list_user_permissions $(RABBIT_ADMIN_USER)
	@echo "Listing exchanges in virtual host '$(VHOST_CUSTOMERS)'..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl list_exchanges -p $(VHOST_CUSTOMERS)
	@echo "Listing queues in virtual host '$(VHOST_CUSTOMERS)'..."
	@docker exec -it $(RABBIT_CONTAINER_NAME) rabbitmqctl list_queues -p $(VHOST_CUSTOMERS)


# ---
# Additional RabbitMQ Docker container utilities

# Check RabbitMQ container status
# in Windows use Select-String instead grep (or setup alias powershell: conf t, `new-alias grep Select-String`)
rabbit-status:
	@echo "Checking RabbitMQ container status..."
	@docker ps -a | grep $(RABBIT_CONTAINER_NAME)

# List all Docker volumes
# in Windows use Select-String instead grep (or setup alias powershell: conf t, `new-alias grep Select-String`)
docker-volumes:
	@echo "Listing all Docker volumes..."
	@docker volume ls | grep $(RABBIT_CONTAINER_NAME)

# Inspect the RabbitMQ container
rabbit-inspect:
	@echo "Inspecting RabbitMQ container..."
	@docker inspect $(RABBIT_CONTAINER_NAME)

# List all Docker containers
docker-containers:
	@echo "Listing all Docker containers..."
	@docker ps -a

# List all Docker containers
docker-running:
	@echo "Listing all Running Docker containers..."
	@docker ps
