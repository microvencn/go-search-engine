<template>
  <el-row class="gse_home-input">
    <el-col :xs="{span: 22, offset: 1}" :sm="{span: 18, offset: 3}" :md="{span: 12, offset: 6}">
      <h1>Go Search Engine</h1>
      <el-form @submit.prevent>
        <el-form-item>
          <template #label>
          </template>
          <SearchInput @submit="redirect" ref="i" />
        </el-form-item>
      </el-form>
    </el-col>
  </el-row>
</template>

<script>

import Result from "@/components/Result";
import SearchInput from "@/components/SearchInput";
export default {
  name: 'HomeView',
  components: {SearchInput, Result},
  data() {
    return {
      query: ""
    }
  },
  mounted() {
    this.$refs.i.focus()
  },
  methods: {
    redirect(str) {
      this.$router.push({name: "result", query: {query: str}})
    },
  }
}
</script>

<style lang="scss">
.gse_home-input {
  margin-top: 30vh;
  input {
    height: 3rem;
    font-size: 1.5rem;
  }
}
</style>
