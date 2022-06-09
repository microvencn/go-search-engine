import pandas as pd
import requests
import threading
import time
import concurrent.futures
import os


df=pd.read_csv('wukong50k_release.csv')
df['data_id'] = ''
for i in range(len(df)):
    df['data_id'][i] = i
urls = df['url']

path = r'images/'

delete_images = []
headers = {
    'User-Agent':'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36 QIHU 360SE'
}

lock = threading.Lock()

def save(img_url, img_name):
    r = requests.request('get', img_url, headers=headers)  #获取网页
    if r.status_code != 200:
        lock.acquire()
        print(img_name, "Error", r.status_code)
        delete_images.append(int(img_name))
        print(delete_images)
        lock.release()
    else:
        with open(path + img_name + '.jpg','wb') as f:
            f.write(r.content)
        f.close()
        # 这里不用加锁
        # lock.acquire()
        print(img_name, "ok")
        # lock.release()

time_1 = time.time()
exe = concurrent.futures.ThreadPoolExecutor(max_workers=20)
for i in range(len(urls)):
    exe.submit(save, urls[i], str(i))
exe.shutdown()
time_2 = time.time()
use_time = int(time_2) - int(time_1)
print(f'总计耗时:{use_time}秒')

# 这个wukong2.csv和wukong.csv含义是一样的
for i in range(len(delete_images)):
    delete_images[i] = int(delete_images[i])
df.drop(delete_images, inplace=True)
df.to_csv("wukong2.csv", index=False)