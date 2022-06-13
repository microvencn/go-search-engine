import base64
# base64加密
with open("images/333.jpg", "rb") as f:
    content = f.read()
    image_encode = base64.b64encode(content)
import requests
# res = requests.get(url='http://43.138.211.125:9999/data/v1.0/getimages',
#                   params={"image_encode": image_encode})
# print(res.text)
import requests
res = requests.get(url='http://110.42.237.111:9999/data/v1.0/getimages',
                  params={"image_encode": image_encode})
print(res.text)