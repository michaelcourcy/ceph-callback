FROM alpine

COPY ceph-callback /ceph-callback

# health endpoint
EXPOSE 8080

CMD [ "/ceph-callback" ]