# 使用 CentOS 7 作为基础镜像
FROM centos:7

# 更新系统并安装必要的工具
RUN yum -y update && yum clean all
RUN yum -y install wget gcc make

# 下载并编译安装 Redis
RUN wget http://download.redis.io/releases/redis-7.2.6.tar.gz
RUN tar xzf redis-7.2.6.tar.gz
RUN cd redis-7.2.6 && make

# 将 Redis 的二进制文件复制到系统路径中
RUN cp redis-7.2.6/src/redis-server /usr/local/bin/
RUN cp redis-7.2.6/src/redis-cli /usr/local/bin/

# 清理安装文件
RUN rm -rf redis-7.2.6.tar.gz redis-7.2.6

# 设置 Redis 的默认配置文件
COPY redis.conf /etc/redis.conf

# 设置容器启动时运行的命令
CMD [ "redis-server", "/etc/redis.conf" ]