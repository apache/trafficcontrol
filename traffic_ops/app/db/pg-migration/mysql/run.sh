docker run \
        --name mysql \
        -p 3306:3306 \
	-v $(pwd)/mysql/conf.d:/etc/mysql/conf.d \
	-v $(pwd)/mysql/initdb.d:/docker-entrypoint-initdb.d \
        -d mysql


