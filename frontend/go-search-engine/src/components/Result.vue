<template>
  <div class="result">
    <el-link class="title">{{title}}</el-link>
    <div class="content">{{content}}</div>
    <div class="score">得分：{{(score*100+"").substring(0, 5)}}</div>
  </div>
</template>

<script>
export default {
  name: "Result",
  props: {
    title: {
      type: String,
      default: ""
    },
    content: {
      type: String,
      default: ""
    },
    score: {
      type: Number,
      default: ""
    },
  },
  data() {
    return {

    }
  }
}
</script>

<style scoped>
.result {
  text-align: left;
  padding-top: 1.5rem;
}

.title {
  font-size: 1.3rem;
}

.content {
  font-size: 1rem;
}

.score {
  color: var(--el-color-primary);
}
</style>
