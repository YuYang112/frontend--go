## elasticsearch操作，gin框架
### elasticsearch安装
```
    1、docker pull docker.elastic.co/elasticsearch/elasticsearch:6.3.2
       
    2、docker run -p 9200:9200 -d --name elasticsearch 96dd1575de0f

    3、docker pull containerize/elastichd

    4、run -p 9800:9800 -d --link elasticsearch:demo containerize/elastichd

    5、在浏览器运行127.0.0.1:9800

    6、输入:http://demo:9200进行连接数据库

```
### docker-compose
```
    1、docker-compose up -d
       
    2、在浏览器运行IP:9800

    3、输入:http://IP:9200进行连接数据库

```
