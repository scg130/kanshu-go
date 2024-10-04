FROM scg130/alpine
# RUN apk add --update ca-certificates && \
#     rm -rf /var/cache/apk/* /tmp/*
#ADD ./.env /micro/.env
ADD ./runapp /micro/runapp
ADD ./resource /micro/resource
RUN chmod -R 755 /micro/runapp
RUN chmod -R 755 /micro/*
WORKDIR /micro
CMD [ "./runapp" ]