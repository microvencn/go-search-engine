<template>
  <el-row>
    <el-col :xs="{span: 22, offset: 1}" :sm="{span: 16, offset: 4}" :md="{span: 12, offset: 6}">
      <el-form @submit.prevent>
        <el-form-item style="height: 32px">
          <template #label>
            <router-link :to="{name: 'home'}"><h2 style="margin: 0">Go Search Engine</h2></router-link>
          </template>
          <SearchInput ref="si" @submit="redirect"/>
        </el-form-item>

        <el-form-item>
          <el-input v-model="filter" placeholder="请输入需要过滤的关键词，以英文逗号分隔"/>
        </el-form-item>
        <el-form-item class="filter">
          <el-button type="primary" @click="redirect($refs.si.getQuery())">Go</el-button>
        </el-form-item>
      </el-form>

      <!--图片搜图片-->
      <el-collapse v-model="openingImgSearch">
        <el-collapse-item title="图片搜索" name="1">
          <el-form>
            <el-form-item>
              <el-upload
                  list-type="picture"
                  action=''
                  accept=".jpg, .png"
                  :limit="1"
                  :auto-upload="false"
                  @change="pic"
                  :on-exceed="handleExceed"
                  ref="upload"
              >
                <el-button type="primary" :loading="imgSearching">图片搜索</el-button>
              </el-upload>
            </el-form-item>
          </el-form>
          <div style="display: flex; width: 100%; overflow:auto;flex-wrap: wrap;justify-content: center">
            <div v-for="img in displayingImg" style="width: 200px;margin-right: 1rem;margin-bottom: 1rem">
              <el-card :body-style="{ padding: 0, height: '300px', fontSize: '0.7rem' }">
                <el-image
                    style="height: 200px; width: 200px"
                    :src="img.url"
                    :preview-src-list="[img.url]"
                    fit="contain"
                />
                <div style="padding: 14px">
                  <span>{{ img.title }}<br><span
                      style="color: var(--el-color-primary)">得分：{{ (img.score * 100 + "").substring(0, 5) }}</span></span>
                </div>
              </el-card>
            </div>
          </div>
          <div v-if="imgStore.length > 0"
               style="justify-content: center;display: flex; margin-top: 1rem; margin-bottom: 1rem">
            <el-pagination background layout="prev, pager, next" :page-count="imgTotalPages"
                           v-model:current-page="imgCurrentPage"
                           @current-change="picPage"/>
          </div>
        </el-collapse-item>
      </el-collapse>


      <div class="guess">
        相关搜索
        <router-link v-for="w in related" :to="{name: 'result', query: {query: w.text}}">
          <span class="unit">
            {{ w.text }}
            <span class="score">
              {{ (w.sim * 100 + "").substring(0, 5) }}
            </span>
          </span>
        </router-link>
      </div>
      <div class="result-container" v-loading="searching">
        <Result v-for="result in results" :title="result.doc" :content="result.doc" :score="result.score"></Result>
      </div>
      <div style="justify-content: center;display: flex; margin-top: 1rem">
        <el-pagination background layout="prev, pager, next" :page-count="total_page" v-model:current-page="page"
                       @current-change="search"/>
      </div>
    </el-col>
  </el-row>
</template>

<script>
import Result from "@/components/Result";
import request, {picRequest} from "@/request";
import SearchInput from "@/components/SearchInput";
import axios from "axios";

export default {
  name: "ResultView",
  components: {SearchInput, Result},
  data() {
    return {
      query: "",
      results: [],
      searching: false,
      related: [],
      relating: false,
      page: 1,
      pageSize: 10,
      total_page: 1,
      filter: [],

      displayingImg: [],
      imgStore: [],
      imgTotalPages: 0,
      imgCurrentPage: 0,
      imgSearching: false,
      openingImgSearch: false
    }
  },
  mounted() {
    this.query = this.$route.query.query
    if (this.$route.query.hasOwnProperty("page")) {
      this.page = this.$route.query.page
    }
    if (this.$route.query.hasOwnProperty("filter")) {
      this.filter = this.$route.query.filter
    }
    this.$refs.si.setQuery(this.query)
    this.search()
    this.relate()
  },
  methods: {
    redirect(str) {
      this.$router.push({name: "result", query: {query: str, filter: this.filter}})
    },
    search() {
      this.searching = true
      request({
        method: 'get',
        url: '/search',
        params: {
          query: this.query,
          page_size: this.pageSize,
          page: this.page,
          filter: this.filter
        }
      }).then((data) => {
        this.searching = false
        this.results = data.data
        this.total_page = Math.ceil(data.total / this.pageSize)
      }).catch((err) => {
        this.$message.error("搜索失败请重试")
      })
    },
    relate() {
      this.relating = true
      request({
        method: 'get',
        url: '/related',
        params: {
          query: this.query
        }
      }).then((data) => {
        this.relating = false
        this.related = data.data
      })
    },
    base64(file) {
      return new Promise(function (resolve, reject) {
        let reader = new FileReader();
        let imgResult = "";
        reader.readAsDataURL(file);
        reader.onload = function () {
          imgResult = reader.result;
        };
        reader.onerror = function (error) {
          reject(error);
        };
        reader.onloadend = function () {
          resolve(imgResult);
        };
      });
    },
    handleExceed(file, fileList) {
      this.$refs.upload.clearFiles()
      this.$refs.upload.handleStart(file[0])
    },
    pic(file) {
      const isJPG = file.raw.type === 'image/jpeg';
      const isPNG = file.raw.type === 'image/png';
      const isLt5M = file.raw.size / 1024 / 1024 < 5;
      if (!isJPG && !isPNG) {
        alert("图片格式必须为 JPG 或 PNG")
        return
      }
      if (!isLt5M) {
        alert("图片必须小于 5 MB")
        return
      }
      this.base64(file.raw).then(base => {
        base = base.replace(/^data:image\/\w+;base64,/, "")
        this.imgSearching = true
        axios.get("http://110.42.237.111:9999/data/v1.0/getimages", {
          params: {
            image_encode: base
          }
        }).then((res) => {
          this.imgStore = res.data.data
          this.imgCurrentPage = 1
          this.imgTotalPages = Math.ceil(this.imgStore.length / 6)
          this.picPage(1)
          this.imgSearching = false
        }).catch((res) => {
          this.$message.error("搜索失败请重试")
        })
      });
    },
    picPage(currentPage) {
      this.displayingImg = []
      let offset = (currentPage - 1) * 6
      for (let i = offset; i < offset + 6 && i < this.imgStore.length; i++) {
        this.displayingImg.push(this.imgStore[i])
      }
    },
  },
  watch: {
    $route() {
      this.query = this.$route.query.query
      this.filter = this.$route.query.filter
      this.$refs.si.setQuery(this.query)
      this.page = 1
      this.search()
      this.relate()
    }
  }
}
</script>

<style scoped lang="scss">
.result-container {
  min-height: 100px;
}

.guess {
  line-break: strict;
  text-align: left;
  min-height: 3rem;
  padding-top: 0.5rem;
  padding-bottom: 0.5rem;
  line-height: 3rem;

  .unit {
    color: #1a1a1a;
    line-height: 2rem;
    display: inline-block;
    line-break: strict;
    background-color: #c4e1ff;
    margin-left: 0.5rem;
    height: 2rem;
    font-size: 1rem;
    padding: 0.2rem 0.5rem;
    border-radius: 10px;

    .score {
      padding-left: 0.5rem;
      color: var(--el-color-primary);
    }
  }
}

.filter :deep(.el-form-item__content) {
  justify-content: center !important;
}
</style>
