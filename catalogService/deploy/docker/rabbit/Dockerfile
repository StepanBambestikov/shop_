FROM rabbitmq:3.11.2-management

RUN apt-get update && \
  apt-get install -y curl && \
  rm -rf /var/lib/apt/lists/*

RUN curl -L https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/3.11.1/rabbitmq_delayed_message_exchange-3.11.1.ez > $RABBITMQ_HOME/plugins/rabbitmq_delayed_message_exchange-3.11.1.ez

RUN rabbitmq-plugins enable --offline rabbitmq_delayed_message_exchange

COPY rabbitmq.config /etc/rabbitmq/
COPY definitions.json /etc/rabbitmq/

EXPOSE 5672
EXPOSE 15672