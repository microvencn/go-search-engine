<template>
  <el-row>
    <el-col :xs="{span: 18, offset: 3}" :sm="{span: 16, offset: 4}" :md="{span: 12, offset: 6}">
      <el-form @submit.prevent>
        <el-form-item style="height: 32px">
          <template #label>
            <router-link :to="{name: 'home'}"><h2 style="margin: 0">Go Search Engine</h2></router-link>
          </template>
          <SearchInput ref="si" />
        </el-form-item>

      </el-form>
      <Result v-for="result in results" :title="result.doc" :content="result.doc" :score="result.score"></Result>
    </el-col>
  </el-row>
</template>

<script>
import Result from "@/components/Result";
import request from "@/request";
import SearchInput from "@/components/SearchInput";

export default {
  name: "ResultView",
  components: {SearchInput, Result},
  data() {
    return {
      query: "",
      results: []
    }
  },
  mounted() {
    this.query = this.$route.query.query
    this.$refs.si.setQuery(this.query)
    this.search()
  },
  methods: {
    search() {
      request({
        method: 'get',
        url: '/search',
        params: {
          query: this.query,
          page_size: 20,
          page: 1
        }
      }).then((data) => {
        this.results = data.data
      })
    }
  },
  watch: {
    $route() {
      this.query = this.$route.query
      this.search()
    }
  }
}
</script>

<style scoped>

</style>
