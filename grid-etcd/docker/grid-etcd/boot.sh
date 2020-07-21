echo 'etcd-master creating...'
/usr/local/bin/etcd --name etcd-master --data-dir /mnt/grid/runtime/etcd-data --listen-client-urls http://0.0.0.0:32379 --advertise-client-urls http://0.0.0.0:32379 --listen-peer-urls http://0.0.0.0:32380 --initial-advertise-peer-urls http://0.0.0.0:32380 --initial-cluster etcd-master=http://0.0.0.0:32380 --initial-cluster-token grid-etcd --initial-cluster-state new