# aliyun-oss-website-action

deploy website on aliyun OSS(Alibaba Cloud OSS)

将静态网站部署在阿里云OSS

## 概览
- 在阿里云OSS创建一个存放网站的bucket
- 准备一个域名, 可能需要备案(bucket选择非大陆区域, 可以不备案, 但是如果CDN加速区域包括大陆, 仍然需要备案)
- 在你的网站repo中, 配置github action, action 触发则**增量上传**网站repo生成的资源文件到bucket中
- 通过阿里云OSS的CDN, 可以很方便地加速网站的访问, 支持HTTPS
> 阿里云HTTPS免费证书停止自动续签, 但是可以自己[申请免费的证书](https://help.aliyun.com/document_detail/156645.htm), 具体解决方案参考[该公告](https://help.aliyun.com/document_detail/479351.html)

## Usage

```yml
    - name: upload files to OSS
      uses: fangbinwei/aliyun-oss-website-action@v1
      with:
          accessKeyId: ${{ secrets.ACCESS_KEY_ID }}
          accessKeySecret: ${{ secrets.ACCESS_KEY_SECRET }}
          bucket: your-bucket-name
          # use your own endpoint
          endpoint: oss-cn-shanghai.aliyuncs.com
          folder: your-website-output-folder
```
> 如果你使用了environment secret请[查看这里](#配置了environment-secret怎么不生效)
### 配置项
- `accessKeyId`: **必填**
- `accessKeySecret`: **必填**
- `endpoint`: **必填**, 支持指定protocol, 例如`https://example.org`或者`http://example.org`
- `folder`: **必填**, repo打包输出的资源文件夹
- `bucket`: **必填**,部署网站的bucket, 用于存放网站的资源
- `indexPage`: 默认`index.html`.网站首页(用于[静态页面配置](#静态页面配置))
- `notFoundPage`: 默认`404.html`.网站404页面(用于[静态页面配置](#静态页面配置))
- `errorDocumentHTTPCode`: 默认`404`.错误页面返回HTTP Code(用于[静态页面配置](#静态页面配置))
- `incremental`: 默认`true`. 使用增量上传.
- `skipSetting`: 默认`false`, 是否跳过设置[静态页面配置](#静态页面配置)
- `htmlCacheControl`: 默认`no-cache`
- `imageCacheControl`: 默认`max-age=864000`
- `pdfCacheControl`: 默认`max-age=2592000`
- `otherCacheControl`: 默认`max-age=2592000`
- `exclude`: 不上传`folder`下的某些文件/文件夹
- `cname`: 默认`false`. 若`endpoint`填写自定义域名/bucket域名, 需设置为`true`. (使用CDN的场景下, 不推荐使用自定义域名)

## incremental
**开启`incremental`**
上传文件到OSS后, 还会将文件的`ContentMD5`和`Cache-Control`收集到名为`.actioninfo`的私有文件中. 当再次触发action的时候, 会将待上传的文件信息与`.actioninfo`中记录的信息比对, 信息未发生变化的文件将跳过上传步骤, 只进行增量上传. 且在上传之后, 根据`.actioninfo`和已上传的文件信息, 将OSS中多余的文件进行删除.

> `.actioninfo` 记录了上一次action执行时, 所上传的文件信息. 私有, 不可公共读写.

**关闭`incremental`** 或 OSS中不存在`.actioninfo`文件

会执行如下步骤
1. 清除所有OSS中已有的文件
2. 上传新的文件到OSS中

> **计划未来优化这个步骤, 优化后, 先上传新的文件到OSS中, 再diff删除多余的文件.** 

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
      # match dist/tmp2/a/b.txt, not match dist/tmp2/tmp3/a/b.txt
```
> 不支持`**`

或者
```yml
- name: Clean files before upload
  run: rm -f dist/tmp.txt
```

## Docker image
直接使用已经build好的docker image
```yml
    - name: upload files to OSS
      uses: docker://fangbinwei/aliyun-oss-website-action:v1
      # 使用env而不是with, 参数可以见本项目的action.yml
      env:
          ACCESS_KEY_ID: ${{ secrets.ACCESS_KEY_ID }}
          ACCESS_KEY_SECRET: ${{ secrets.ACCESS_KEY_SECRET }}
          BUCKET: your-bucket-name
          ENDPOINT: ali-oss-endpoint
          FOLDER: your-website-output-folder
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
          errorDocumentHTTPCode: '404'
          htmlCacheControl: no-cache
          imageCacheControl: max-age=864001
          otherCacheControl: max-age=2592001
```

## FAQ

### 配合CDN使用时, OSS更新后, CDN未刷新

开启OSS提供的CDN缓存自动刷新功能, 将触发操作配置为`PutObject`, `DeleteObject`.

![2020-12-13-23-51-28](https://image.fangbinwei.cn/github/aliyun-oss-website-action/2020-12-13-23-51-28_2c310155.png)

![2020-12-13-23-51-55](https://image.fangbinwei.cn/github/aliyun-oss-website-action/2020-12-13-23-51-55_5fe79a54.png)

### `endpoint`使用自定义域名, 但是无法上传
1. 如果`endpoint`的域名CNAME记录为阿里云CDN, CDN是否配置了http强制跳转https? 若配置了, 需要在`endpoint`中指定https, 即`endpoint`为`https://example.org`

2. 如果`endpoint`的域名CNAME记录为阿里云CDN, 在CDN为加速范围为全球时有遇到过如下报错`The bucket you are attempting to access must be addressed using the specified endpoint. Please send all future requests to this endpoint.`, 则`endpoint`不能使用自定义域名, 使用OSS源站的endpoint.

### 配置了environment secret怎么不生效

![2021-05-21-16-47-59](https://image.fangbinwei.cn/github/aliyun-oss-website-action/2021-05-21-16-47-59_affec2b0.png)

如果使用environment secret, 那么需要如下类似的配置

```diff

jobs:
  build:
    runs-on: ubuntu-latest
+    environment: your-environment-name

```