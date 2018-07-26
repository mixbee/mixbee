
FROM tianon/ubuntu-core:14.04
ENV MIXBEE_PATH /var/mixbee
RUN mkdir -p $MIXBEE_PATH
RUN apt-get -y install wget curl
COPY mixbee /usr/local/bin
# ADD config.tgz $MIXBEE_PATH
EXPOSE 20334 20335 20336 20337 20338 20339 20340 20341 20000
WORKDIR $MIXBEE_PATH
CMD ["mixbee"]

