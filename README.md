# aliyun-oss-website-action(WIP)

deploy website on aliyun OSS(Alibaba Cloud OSS)

将静态网站部署在阿里云OSS

## 概览
- 在阿里云OSS创建一个存放网站的bucket
- 准备一个域名, 可能需要备案
- 在你的网站repo中, 配置github action, 当action触发, 会清除bucket中原有的文件, 上传网站repo生成的资源文件到bucket中
- 通过阿里云OSS的CDN, 可以很方便地加速网站的访问, 支持HTTPS

## Usage

```yml
    - name: upload files to OSS
      uses: fangbinwei/aliyun-oss-website-action@master
      with:
          accessKeyId: ${{ secrets.ACCESS_KEY_ID }}
          accessKeySecret: ${{ secrets.ACCESS_KEY_SECRET }}
          bucket: "your-bucket-name"
          # like "oss-cn-shanghai.aliyuncs.com"
          endpoint: "ali-oss-endpoint"
          folder: "your-website-output-folder"
```
### 配置项
- `accessKeyId`: **必填**
- `accessKeySecret`: **必填**
- `endpoint`: **必填**
- `folder`: **必填**, repo打包输出的资源文件夹
- `bucket`: **必填**,部署网站的bucket, 用于存放网站的资源
- `indexPage`: 默认`index.html`.网站首页
- `notFoundPage`: 默认`404.html`.网站404页面

## Cache-Control
为上传的资源默认设置的`Cache-Control`如下
|资源类型 | Cache-Control|
|----| ----|
|.html|no-cache|
|.png/jpg...(图片资源)|max-age=864000(10days)|
|other|max-age=2592000(30days)|

## 静态页面配置
默认的, action会将阿里云OSS的静态页面配置成如下
![2020-08-06-03-18-25](https://image.fangbinwei.cn/github/aliyun-oss-website-action/2020-08-06-03-18-25_05d556d8.png)


## Demo
### 部署VuePress项目

```yml

name: deploy vuepress

on:
  push:
    branches:
      - master

jobs:
  build:

    runs-on: ubuntu-latest
    steps:
      # load repo to /github/workspace
    - uses: actions/checkout@v2
      with:
          repository: fangbinwei/blog
          fetch-depth: 0
    - name: Use Node.js
      uses: actions/setup-node@v1
      with:
        node-version: '12'
    - run: npm install yarn@1.22.4 -g
    - run: yarn install
    # 打包文档命令
    - run: yarn docs:build
    - name: upload files to OSS
      uses: fangbinwei/aliyun-oss-website-action@master
      with:
          accessKeyId: ${{ secrets.ACCESS_KEY_ID }}
          accessKeySecret: ${{ secrets.ACCESS_KEY_SECRET }}
          bucket: "your-bucket-name"
          endpoint: "oss-cn-shanghai.aliyuncs.com" 
          folder: ".vuepress/dist"
```
具体可以参考本项目的[workflow](.github/workflows/test.yml), npm/yarn配合`action/cache`加速依赖安装

### Hugo
TODO
### Hexo
TODO