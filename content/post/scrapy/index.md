+++
date = '2025-10-30T16:31:53+08:00'
draft = true
title = 'Scrapy'



categories=['编程']

tags=['python']

+++

> 自用, 不完全由我自己写的. 用作备忘.

## 制作 Scrapy 爬虫 一共需要4步：

1. 新建项目 (scrapy startproject xxx)：新建一个新的爬虫项目
2. 明确目标 （编写items.py）：明确你想要抓取的目标
3. 制作爬虫 （spiders/xxspider.py）：制作爬虫开始爬取网页
4. 存储内容 （pipelines.py）：设计管道存储爬取内容
