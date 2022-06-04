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
CHUNK = 1024*4

print('Finding similar images...')
CTS_test = len(imagefeat)//CHUNK
if len(imagefeat)%CHUNK!=0: CTS_test += 1

a = 0
b = len(imagefeat)
    
distances_test = torch.matmul(imagefeat_test, imagefeat[a:b].T).T
distances_test = distances_test.data.cpu().numpy()
distances_test = distances_test.reshape(-1)
distances_test_sort = sorted(enumerate(distances_test), key = lambda x:x[1], reverse=True)
print(distances_test_sort)