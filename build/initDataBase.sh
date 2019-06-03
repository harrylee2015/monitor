#!/usr/bin/env bash
# ----------------------------------------------------------------------
# name:         initDataBase.sh
# version:      1.0
# createTime:   2019-05-20
# description:  初始化monitor项目用到的相关数据表，运行前提确保机器能够连上外网
# author:       harry lee
# ----------------------------------------------------------------------
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
# Check the network status
NET_NUM=`ping -c 4 www.baidu.com |awk '/packet loss/{print $6}' |sed -e 's/%//'`
if [ -z "$NET_NUM" ] || [ $NET_NUM -ne 0 ];then
        echo "Please check your internet"
        exit 1
fi
sqlite3 -version
if [ $? -eq 0 ];then
    echo "This machine has installed sqlite3!"
else
    echo "This machine has not installed sqlite3,now start intall..."
    sudo apt install sqlite3
fi
#初始化数据库
mkdir -p datadir
sqlite3  datadir/monitor.db <install.sql
echo "monitor.db init sucessfully."