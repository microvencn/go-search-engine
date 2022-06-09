import numpy as np 
import pandas as pd 
import gc
from PIL import Image
import torch
torch.manual_seed(0)
import torchvision.models as models
import torchvision.transforms as transforms
import torchvision.datasets as datasets
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim
from torch.autograd import Variable
from torch.utils.data.dataset import Dataset
from sklearn.preprocessing import normalize

imagefeat = np.load("imagefeat_wukong50k.npy", allow_pickle=True)
imagefeat = torch.from_numpy(imagefeat)

# 图片入口在这
data_pos = ["images/3.jpg"]

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

DEVICE = 'cpu'

imgmodel_test = WukongImageEmbeddingNet()
imgmodel_test = imgmodel_test.to(DEVICE)
imagefeat_test = []
with torch.no_grad():
    for data in imageloader_test:
        data = data.to(DEVICE)
        feat = imgmodel_test(data)
        feat = feat.reshape(feat.shape[0], feat.shape[1])
        feat = feat.data.cpu().numpy()
        imagefeat_test.append(feat)
        
# l2 norm to kill all the sim in 0-1
imagefeat_test = np.vstack(imagefeat_test)
imagefeat_test = normalize(imagefeat_test)

imagefeat_test = torch.from_numpy(imagefeat_test)
# imagefeat_test = imagefeat_test.cuda()

preds_test = []
print('Finding similar images...')

a = 0
b = len(imagefeat)

DATA_PATH = ''
IMAGE_PATH = 'images/'
train = pd.read_csv(DATA_PATH + 'wukong.csv')

distances_test = torch.matmul(imagefeat, imagefeat_test.T).T
distances_test = distances_test.data.cpu().numpy()

IDX = np.where(distances_test[0,]>0.85)[0][:]

# print(IDX, len(IDX))
o = train.iloc[IDX].data_id.values
# print(o, len(o))

# print(distances_test[0,])
# print(len(distances_test[0,]))
# print(a, b, b-a)
# print(o)
distances_test_sort = []
for i in range(0, len(o)):
    distances_test_sort.append((o[i], distances_test[0,][IDX[i]]))

# distances_test = distances_test[0,].reshape(-1)
distances_test_sort = sorted(distances_test_sort, key = lambda x:x[1], reverse=True)
print(distances_test_sort)