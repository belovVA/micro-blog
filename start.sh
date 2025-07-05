docker build  -t micro-blog .


docker run -d \
  --name micro-blog \
  -p 8080:8080 \
  micro-blog
