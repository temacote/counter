FROM alpine:3.7

RUN apk update && apk add bash curl
RUN mkdir /wait-for-it
RUN curl https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh > /wait-for-it/wait-for-it.sh
RUN chmod +x /wait-for-it/wait-for-it.sh
ENV PATH="/wait-for-it:${PATH}"
CMD ["/docker/сonfigure.sh"]
