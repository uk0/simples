
 CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -o blog

 rsync -av \
--exclude='.git'  \
--exclude='.DS_Store'  \
 /Users/firshme/Desktop/work/simples/tpl root@vm2:/root/blog/


  rsync -av \
--exclude='.git'  \
--exclude='.DS_Store'  \
 /Users/firshme/Desktop/work/simples/blog root@vm2:/root/blog/





