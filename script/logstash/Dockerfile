FROM elastic/logstash:8.18.0

COPY ./pipeline/logstash.conf /etc/logstash/conf.d/

CMD logstash -f /etc/logstash/conf.d/logstash.conf