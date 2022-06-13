<template>
  <el-autocomplete
      v-model="query"
      @keyup.enter.native="submit"
      placeholder="Search Something..."
      :fetch-suggestions="querySearch"
      clearable
      style="width: 100%;"
      @select="submit"
      ref="i"
  />
</template>

<script>
import request from "@/request";

export default {
  name: "SearchInput",
  data() {
    return {
      query: ""
    }
  },
  methods: {
    focus() {
      this.$refs.i.focus()
    },
    submit() {
      this.$emit("submit", this.query)
    },
    querySearch(prefix, cb) {
      request({
        url: "/auto",
        params: {
          prefix: this.query
        }
      }).then((data) => {
        let result = []
        for (let i = 0; i < data.data.length; i ++) {
          result.push({ value: data.data[i] })
        }
        cb(result)
      })
    },
    setQuery(str) {
      this.query = str
    },
    getQuery() {
      return this.query
    }
  }
}
</script>

<style scoped lang="scss">

</style>
