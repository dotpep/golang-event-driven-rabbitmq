producer:
	@go run .\cmd\producer\main.go

consumer:
	@go run .\cmd\consumer\main.go

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

rabbit-restart:
	@docker restart rabbitmq

# Remove container
rabbit-rm:
	@docker rm rabbitmq

# Check logs of container
rabbit-logs:
	@docker logs -f rabbitmq

# Access to shell container:
rabbit-exec:
	@docker exec -it rabbitmq bash

# Seting Up RabbitMQ Environment

# 1. make rabbit-ctl-new_admin_user
# 2. make rabbit-ctl-new_customers_vhost
# 3. make rabbit-admin-new_customer_events_exchange

# Adding new Admin user,
# with password: admin,
# setting Administrator privilege,
# and Deleting Guest default user
# 1. (rabbitmqctl, add_user) docker exec -it <container_name> rabbitmqctl add_user <newusername> <secretpassword>
# 2. (rabbitmqctl, set_user_tags) docker exec -it <container_name> rabbitmqctl set_user_tags <username> <privilege>
# 2. (rabbitmqctl, delete_user) docker exec -it <container_name> rabbitmqctl delete_user <username>
rabbit-ctl-new_admin_user:
	@docker exec -it rabbitmq rabbitmqctl add_user admin admin
	@docker exec -it rabbitmq rabbitmqctl set_user_tags admin administrator
	@docker exec -it rabbitmq rabbitmqctl delete_user guest

# Adding new Virtual Host named Customers,
# for Admin user, with regex of all permissions (configurations, write, read)
# (rabbitmqctl, add_vhost) 1. docker exec -it <container_name> rabbitmqctl add_vhost <new_vhost_name>
# (rabbitmqctl, set_permissions) 2. docker exec -it <container_name> rabbitmqctl set_permissions -p <vhost_name> <user> <configurations_vhost> <write_regex> <read_regex>
rabbit-ctl-new_customers_vhost:
	@docker exec -it rabbitmq rabbitmqctl add_vhost customers
	@docker exec -it rabbitmq rabbitmqctl set_permissions -p customers admin ".*" ".*" ".*"

# Declaring new Exchange named customer_events,
# for Virtual Host of Customer,
# type of Topic Exchange,
# for Admin user.
# (rabbitmqadmin, declare exchange) 1. docker exec -it <container_name> rabbitmqadmin declare exchange --vhost=<vhost_name> name=<exchange_name> type=<exchange_type> -u <user> -p <user_password> durable=<durable_bool_param>
rabbit-admin-new_customer_events_exchange:
	@docker exec -it rabbitmq rabbitmqadmin declare exchange --vhost=customers name=customer_events type=topic -u admin -p admin durable=true
	@docker exec -it rabbitmq rabbitmqctl set_topic_permissions -p customers admin customer_event ".*" ".*"
