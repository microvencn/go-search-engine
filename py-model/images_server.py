import numpy as np
import pandas as pd
from PIL import Image
import torch
torch.manual_seed(0)
import torchvision.models as models
import torchvision.transforms as transforms
import torch.nn as nn
from torch.utils.data.dataset import Dataset
from sklearn.preprocessing import normalize
import base64
from flask import Flask
from flask import request
import json
import pandas as pd
from flask_cors import CORS, cross_origin

app = Flask(__name__)
CORS(app, supports_credentials=True)

class WukongImageDataset(Dataset):
    def __init__(self, img_path, transform):
        self.img_path = img_path
        self.transform = transform

    def __getitem__(self, index):
        img = Image.open(self.img_path[index]).convert('RGB')
        img = self.transform(img)
        return img

    def __len__(self):
        return len(self.img_path)


class WukongImageEmbeddingNet(nn.Module):
    def __init__(self):
        super(WukongImageEmbeddingNet, self).__init__()

        model = models.resnet18(True)
        model.avgpool = nn.AdaptiveMaxPool2d(output_size=(1, 1))
        model = nn.Sequential(*list(model.children())[:-1])
        model.eval()
        self.model = model

    def forward(self, img):
        out = self.model(img)
        return out




# app.route（）方法是路由控制器，第一个参数是浏览器访问的路径，methods参数是请求的方式，常用GET、POST
@app.route('/data/v1.0/getimages', methods=['GET'])
@cross_origin(supports_credentials=True)
def start():
    # 参数类型强制为str, 如果收到的不是str类型,则结果是None
    image_encode = request.args.get('image_encode', type=str)

    # 对接收的3个参数简单验证下,只有每一个参数都传入了数值才会返回正确的值,get_data方法会从数据库中查询并返回
    if image_encode:
        code = 0
        msg = '成功'
        # data = image_encode
        data = get_data(image_encode)
        # data = [{'data_id': 3, 'score': 1.0000001, 'url': 'https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fa.vpimg2.com%2Fupload%2Fmerchandise%2F227958%2FLYQ-S314186413-3.jpg&refer=http%3A%2F%2Fa.vpimg2.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=jpeg?sec=1632501843&t=b0a3b843f9ecebd71fe6f27643c17486', 'title': '女童黄色连衣裙'}, {'data_id': 23817, 'score': 0.9024606, 'url': 'https://pics5.baidu.com/feed/b8014a90f603738d06fa397ff91ebd59f919eccc.jpeg?token=59bcb50c19c01df1cb8ce81b4a6a3d80&s=B581DB11DD0670CE009D98C80300F031', 'title': '优雅女神孙允珠:一身灰色运动装出席,满满的运动气息'}]
        #print(data)
        #print(type(data))

        # data = get_data(service_type, city_id, car_type_id)
        # if not data:  # 当数据库中找不到对应的记录需要返回的信息
        #     code = 2
        #     msg = '请求的数据不存在'
        #     data = None
    else:
        code = 1
        msg = '参数有误'
        data = None

    # 将响应体转成json数据
    return_data = json.dumps({
        "code": code,
        "msg": msg,
        'data': data
    }, ensure_ascii=False)
    return return_data

def get_data(image_encode):
    imagefeat = np.load("imagefeat_wukong50k.npy", allow_pickle=True)
    imagefeat = torch.from_numpy(imagefeat)

    # base64加密
    # with open("images/3.jpg", "rb") as f:
    #     content = f.read()
    #     image_encode = base64.b64encode(content)

    import random
    import string
    image_decode = base64.b64decode(image_encode)
    str = random.sample(string.ascii_letters + string.digits, 16)
    str = ''.join(str)

    image_pos = "images/" + str + ".jpg"
    file = open(image_pos, 'wb')
    file.write(image_decode)
    file.close()

    # 图片入口在这
    data_pos = [image_pos]

    imagedataset_test = WukongImageDataset(
        data_pos,
        transforms.Compose([
            transforms.Resize((512, 512)),
            transforms.ToTensor(),
            transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
        ]))

    imageloader_test = torch.utils.data.DataLoader(
        imagedataset_test,
        batch_size=40, shuffle=False, num_workers=4
    )

    imgmodel_test = WukongImageEmbeddingNet()
    imagefeat_test = []
    with torch.no_grad():
        for data in imageloader_test:
            feat = imgmodel_test(data)
            feat = feat.reshape(feat.shape[0], feat.shape[1])
            feat = feat.data.cpu().numpy()
            imagefeat_test.append(feat)

    imagefeat_test = np.vstack(imagefeat_test)
    imagefeat_test = normalize(imagefeat_test)

    imagefeat_test = torch.from_numpy(imagefeat_test)

    preds_test = []

    # print('Finding similar images...')

    a = 0
    b = len(imagefeat)

    DATA_PATH = ''
    IMAGE_PATH = 'images/'
    train = pd.read_csv(DATA_PATH + 'wukong.csv')

    distances_test = torch.matmul(imagefeat, imagefeat_test.T).T
    distances_test = distances_test.data.cpu().numpy()

    IDX = np.where(distances_test[0,] > 0.5)[0][:]

    # print(IDX, len(IDX))
    o = train.iloc[IDX].data_id.values
    urls = train.iloc[IDX].url.values
    titles = train.iloc[IDX].title.values
    # print(o, len(o))

    # print(distances_test[0,])
    # print(len(distances_test[0,]))
    # print(a, b, b-a)
    # print(o)
    distances_test_sort = []
    for i in range(0, len(o)):
        distances_test_sort.append((int(o[i]), float(distances_test[0,][IDX[i]]), urls[i], titles[i]))

    # distances_test = distances_test[0,].reshape(-1)
    distances_test_sort = sorted(distances_test_sort, key=lambda x: x[1], reverse=True)[:50]

    jsonList = []
    for i in range(0, len(distances_test_sort)):
        tmp = {}
        tmp["data_id"] = distances_test_sort[i][0]
        tmp["score"] = distances_test_sort[i][1]
        tmp["url"] = distances_test_sort[i][2]
        tmp["title"] = distances_test_sort[i][3]
        jsonList.append(tmp)
    # print(distances_test_sort)
    return jsonList


# def get_data(service_type, city_id, car_type_id):
#     sql = f"select * from table " \
#           f"where service_type={service_type} " \
#           f"and city_id={city_id} " \
#           f"and car_type_id={car_type_id} " \
#           f"and date(update_time)=curdate()"
#     data = pd.read_sql(sql, engine)  # 通过pandas读取数据库，非常方便有木有
#     json_data = data.to_json(orient='table', force_ascii=False, index=False)
#     result = json.loads(json_data)['data']
#     return result


if __name__ == '__main__':
    # 创建一个flask api 对象

    app.run(host='0.0.0.0', port=9999, debug=False)  # 80端口在浏览器地址中可以省略不填，其他端口就不行