# micro
#

update user SET role_ids = CONVERT('[1]', BINARY) where id = 1;

html 文件添加
<script>
     var head= document.getElementsByTagName('head')[0];  var script= document.createElement('script');  script.type= 'text/javascript';  script.src= 'https://res.zvo.cn/translate/inspector.js';  head.appendChild(script); 
</script>


docker run --restart=always -p 6379:6379 --name myredis -d redis:6.2.1  --requirepass smd013012



  #  micro-api micro-web 与micro v2 不兼容  需重新构建
  # micro-api:
  #   container_name: micro-api
  #   image: scg130/micro
  #   ports:
  #     - "8088:8080"
  #   command: --registry=etcd --registry_address=172.22.0.2:2379 --api_namespace=go.micro.api api --handler=api
  #   networks:
  #     hx_net:
  #       ipv4_address: 172.22.0.7

  # micro-web:
  #   container_name: micro-web
  #   image: scg130/micro
  #   ports:
  #     - "8082:8082"
  #   command: --registry=etcd --registry_address=172.22.0.2:2379 --web_namespace=go.micro.web --client=grpc --enable_stats web
  #   networks:
  #     hx_net:
  #       ipv4_address: 172.22.0.8