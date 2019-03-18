Cachr
=====

S3 based Remote Caching for CI/CD

### Intro

A key based file store. While storing the files, you give the key at which to save the files. When retrieving, you give back again the key & it’ll fetch the cache from remote store.

Files are zipped before remote sync. And are unzipped after download.

Buckets at which it saves the file is configured via ENV: CACHR_BUCKET & region via AWS_REGION


### Usage

```
cachr exists 'andromeda/assets/sha-a12d221daa23dfgh2'

cachr save 'andromeda/assets/sha-a12d221daa23dfgh2' public/assets/ public/tmp/cache/assets/

cachr get 'andromeda/assets/sha-a12d221daa23dfgh2'
```
