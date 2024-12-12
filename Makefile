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
