# 'caching_sha2_password' cannot be loaded：
# mysql 8.0 默认使用 caching_sha2_password 身份验证机制；客户端不支持新的加密方式。
# select host,user,plugin,authentication_string from mysql.user;
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '123456';
# select host,user,plugin,authentication_string from mysql.user;