FROM elastic/kibana:8.18.0

USER root
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh
USER kibana

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["kibana"]