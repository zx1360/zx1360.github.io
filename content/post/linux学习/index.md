+++
date = '2025-11-02T14:56:39+08:00'
draft = true
title = 'Linux学习'



categories=['编程']

tags=['linux', '云服务器']

+++



### 基础linux指令

- 登陆服务器

  `ssh 用户名@IP`

- 文件拷贝

  `scp 本地文件路径 用户名@IP:目标目录路径`

  `scp 用户名@IP:目标文件路径 本地目录路径`

- 文件操作

  - `mkdir dirname`
  - `cat aa.txt -n` (带行序号)
  - `less bb.txt` 更好阅读体验的查看文件内容
  - `vim cc.txt` 编辑

## 部署应用

- 编译linux程序使用scp上传
- 到对应目录添加可执行权限(首次运行必要)
  - `chmod +x myapp`

- 运行

  - 直接运行
    - 随终端关闭而终止

  - nohup (不随终端一同关闭)

    - `nohup ./myapp &`

    - ```bash
      ps -ef | grep myapp  # 找到程序的进程ID（第二列）
      kill -9 进程ID  # 强制终止
      ```

  - systemd (加入系统服务)

    1. 创建服务配置文件：

       ```bash
       sudo vim /etc/systemd/system/myapp.service  # 用vim编辑
       ```

    2. 按`i`进入编辑模式，输入以下内容（根据实际路径修改）：

       ```ini
       [Unit]
       Description=My Go Application  # 服务描述
       After=network.target  # 网络启动后再启动服务
       
       [Service]
       User=root  # 运行用户
       WorkingDirectory=/root/app  # 程序所在目录
       ExecStart=/root/app/myapp  # 程序路径
       Restart=always  # 程序异常退出时自动重启
       RestartSec=3  # 重启间隔3秒
       
       [Install]
       WantedBy=multi-user.target  # 多用户模式下开机启动
       ```

    3. 按`Esc`，输入`:wq`保存退出。

    4. 管理服务的命令：

       ```bash
       sudo systemctl daemon-reload  # 重新加载服务配置
       sudo systemctl start myapp    # 启动服务
       sudo systemctl status myapp   # 查看服务状态（是否运行）
       sudo systemctl stop myapp     # 停止服务
       sudo systemctl enable myapp   # 设置开机自启动
       ```

       查看程序日志：

       ```bash
       journalctl -u myapp -f  # 实时查看服务日志
       ```

- 查看网络端口

  `netstat -tuln | grep 8080  # 若有输出，说明程序正常监听`
