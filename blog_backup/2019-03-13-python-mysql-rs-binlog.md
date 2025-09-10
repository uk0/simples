---
layout: post
categories: python
title: python 监控binlog实现hue一个小插件
date: 2019-03-13 21:43:56 +0800
description: python 实现hue用户自动创建插件
keywords: python binlog
---


### 代码Python可以使用


```python
from pymysqlreplication import BinLogStreamReader
from pymysqlreplication.row_event import (
    DeleteRowsEvent,
    UpdateRowsEvent,
    WriteRowsEvent,
)
import threading
import paramiko
import logging.handlers

LOG_FILE = r'hue_auto_create_user.log'

handler = logging.handlers.RotatingFileHandler(LOG_FILE, maxBytes=1024 * 1024, backupCount=5, encoding='utf-8')  # 实例化handler
fmt = '%(asctime)s - %(levelname)s - %(message)s'

formatter = logging.Formatter(fmt)  # 实例化formatter
handler.setFormatter(formatter)  # 为handler添加formatter

logger = logging.getLogger('hue_auto_create_user')  # 获取名为hue_auto_create_user的logger
logger.addHandler(handler)  # 为logger添加handler
logger.setLevel(logging.DEBUG)


MYSQL_SETTINGS = {
    "host": "cdh-m1.temp.online",
    "port": 3306,
    "user": "root",
    "passwd": "12345"
}

HOST_ARRAY = [
    '10.50.40.1',
    '10.50.40.2',
    '10.50.40.3',
    '10.50.40.4',
    '10.50.40.5',
    '10.50.40.6',
    '10.50.40.7',
    '10.50.40.8',
    '10.50.40.9',
    '10.50.40.10',
]


def ssh_host_createUserAndGroup(ip, username, passwd, cmd):
    try:
        # 指定本地的RSA私钥文件,如果建立密钥对时设置的有密码，password为设定的密码，如无不用指定password参数
        # pkey = paramiko.RSAKey.from_private_key_file('/home/super/.ssh/id_rsa')
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(ip, 22, username, passwd, timeout=5)
        for m in cmd:
            stdin, stdout, stderr = ssh.exec_command(m)
            out = stdout.readlines()
            # 屏幕输出
            for o in out:
                print(o)
        logger.info(u"ssh 远程执行 {} Success".format(ip))
        ssh.close()
    except:
        logger.error(u"ssh 远程执行 {} Error".format(ip))


def startTask(user):
    cmd = ['useradd {}'.format(user), 'groupadd {}'.format(user),
           'usermod  -a  -G {} {}'.format(user, user)]  # 你要执行的命令列表
    username = "root"  # 用户名
    passwd = "123456"  # 密码

    for ip in HOST_ARRAY:
        a = threading.Thread(target=ssh_host_createUserAndGroup, args=(ip, username, passwd, cmd))
        logger.info(u"Host {} 执行命令 {}".format(ip,cmd))
        a.start()


def updateMysql(uid,userName):
    import pymysql
    cnx = pymysql.connect(user='root',
                          password='123456',
                          host='cdh-m1.temop.online',
                          database='hue',
                          port=3306,
                          charset='utf8'
                          )
    cursor = cnx.cursor()
    try:
        cursor.execute("update auth_user set is_superuser=0 where id={}".format(uid))
        #INSERT IGNORE INTO
        cursor.execute("insert INTO hue.auth_group ( id, name) SELECT  (auth_group.id+1), '{}' FROM hue.auth_group order by id DESC limit 1".format(userName))
    except Exception:
        logger.error("已经存在 Group {}".format(userName))
        print("已经存在 Group {}".format(userName))


    cursor.execute("insert INTO hue.auth_user_groups(user_id,group_id)  SELECT {},auth_group.id FROM hue.auth_group where hue.auth_group.name='{}'".format(uid, userName))
    cnx.commit()
    cnx.close()

def main():
    stream = BinLogStreamReader(
        connection_settings=MYSQL_SETTINGS,
        server_id=100,
        blocking=True,
        resume_stream=True,
        only_events=[DeleteRowsEvent, WriteRowsEvent, UpdateRowsEvent])

    for binlogevent in stream:
        e_start_pos, last_pos = stream.log_pos, stream.log_pos
        for row in binlogevent.rows:
            event = {"schema": binlogevent.schema,
                     "table": binlogevent.table,
                     "type": type(binlogevent).__name__,
                     "row": row
                     };
            if binlogevent.table == "auth_user" and type(binlogevent).__name__ == "DeleteRowsEvent":
                logger.info(u"DELETE User " + row['values']['username'])
                print(u"DELETE User " + row['values']['username'])

            if binlogevent.table == "auth_user" and type(binlogevent).__name__ == "WriteRowsEvent":
                userName = row['values']['username'];
                uid = row['values']['id'];
                print(u"INSERT User {} Uid {}".format(row['values']['username'],uid))

                startTask(userName)
                logger.info(u"[INFO] 添加用户到集群所有节点 -> {}".format(userName) + "\n")
                updateMysql(uid,userName)
                logger.info(u"[INFO] 去掉Hue SuperUser {} Uid={}".format(userName, uid) + "\n")


if __name__ == "__main__":
    main()


```

转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
