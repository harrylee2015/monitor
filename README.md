# monitor

* chain33 para chain monitor

## 使用方法

 1. 编译monitor 可执行文件

    前提是要有go基础环境安装好了
    在本项目目录下执行
    ```bash
    make build
    ```
 
 2. 将build 目录下的所有文件都拷贝到安装目录下
    ```bash
    mv build/*   targetDir/
    ```
 
 3. 启动guess服务

    > 第一次使用需要安装和初始化数据库，以后就不需要了
     
     ```bash
     cd targetDir
     bash initDataBase.sh
     ```
    > 启动
    
     ```bash
     nohup ./monitor -f monitor.toml > console.log 2>&1 &
     ```