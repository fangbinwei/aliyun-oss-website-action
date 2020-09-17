# aliyun-oss-website-action

deploy website on aliyun OSS(Alibaba Cloud OSS)

将静态网站部署在阿里云OSS

## 概览
- 在阿里云OSS创建一个存放网站的bucket
- 准备一个域名, 可能需要备案(bucket选择非大陆区域, 可以不备案, 但是如果CDN加速区域包括大陆, 仍然需要备案)
- 在你的网站repo中, 配置github action, 当action触发, 会清除bucket中原有的文件, 上传网站repo生成的资源文件到bucket中
- 通过阿里云OSS的CDN, 可以很方便地加速网站的访问, 支持HTTPS

## Usage

```yml
    - name: upload files to OSS
      uses: fangbinwei/aliyun-oss-website-action@v1
      with:
          accessKeyId: ${{ secrets.ACCESS_KEY_ID }}
          accessKeySecret: ${{ secrets.ACCESS_KEY_SECRET }}
          bucket: your-bucket-name
          # e.g. "oss-cn-shanghai.aliyuncs.com", 也可以填写自定义域名(需要配置cname 为 true)
          endpoint: ali-oss-endpoint
          folder: your-website-output-folder
```
### 配置项
- `accessKeyId`: **必填**
- `accessKeySecret`: **必填**
- `endpoint`: **必填**
- `folder`: **必填**, repo打包输出的资源文件夹
- `bucket`: **必填**,部署网站的bucket, 用于存放网站的资源
- `cname`: 默认`false`. 若`endpoint`填写自定义域名/bucket域名, 需设置为`true`.(若自定义域名解析到了CDN, 则不要使用该配置)
- `indexPage`: 默认`index.html`.网站首页(用于[静态页面配置](#静态页面配置))
- `notFoundPage`: 默认`404.html`.网站404页面(用于[静态页面配置](#静态页面配置))
- `skipSetting`: 默认`false`, 是否跳过设置[静态页面配置](#静态页面配置)
- `htmlCacheControl`: 默认`no-cache`
- `imageCacheControl`: 默认`max-age=864000`
- `otherCacheControl`: 默认`max-age=2592000`
- `exclude`: 不上传`folder`下的某些文件/文件夹

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

若不需要action来设置, 可以配置`skipSetting`为`true`

## exclude
如果`folder`下的某些文件不需要上传


```yml
    - name: exclude some files
      uses: fangbinwei/aliyun-oss-website-action@v1
      with:
        folder: dist
        exclude: |
          tmp.txt
          tmp/
          tmp2/*.txt
          tmp2/*/*.txt
      # match dist/tmp.txt
      # match dist/tmp/
      # match dist/tmp2/a.txt
      # match dist/tmp2/a/b.txt not match dist/tmp2/tmp3/b.txt
```
> 不支持`**`

或者
```yml
- name: Clean files before upload
  run: rm -f dist/tmp.txt
```

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
      uses: fangbinwei/aliyun-oss-website-action@v1
      with:
          accessKeyId: ${{ secrets.ACCESS_KEY_ID }}
          accessKeySecret: ${{ secrets.ACCESS_KEY_SECRET }}
          bucket: "your-bucket-name"
          endpoint: "oss-cn-shanghai.aliyuncs.com" 
          folder: ".vuepress/dist"
```
具体可以参考本项目的[workflow](.github/workflows/test.yml), npm/yarn配合`action/cache`加速依赖安装

### Vue

[see here](https://github.com/fangbinwei/oss-website-demo-spa-vue)

```yml
- name: upload files to OSS
      uses: fangbinwei/aliyun-oss-website-action@v1
      with:
          accessKeyId: ${{ secrets.ACCESS_KEY_ID }}
          accessKeySecret: ${{ secrets.ACCESS_KEY_SECRET }}
          bucket: website-spa-vue-demo
          endpoint: oss-spa-demo.fangbinwei.cn
          cname: true
          folder: dist
          notFoundPage: index.html
          htmlCacheControl: no-cache
          imageCacheControl: max-age=864001
          otherCacheControl: max-age=2592001
```

### Hugo
TODO
### Hexo
TODO