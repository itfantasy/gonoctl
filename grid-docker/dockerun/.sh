docker run -v aa:/etc/grid/runtime -p 30089:30089/udp  itfantasy/grid etc/grid/grid-core -proj=/etc/grid/runtime -nodeid=aa -endpoints=["kcp://192.168.10.20:30089"]