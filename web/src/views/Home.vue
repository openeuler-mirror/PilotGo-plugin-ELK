<template>
  <div>
    home页面
  </div>
</template>

<script setup lang="ts">
import { testHttps } from '../api';

const test = () => {
  testHttps().then(res => {
    console.log(res)
  })
}
test();
</script>

<style scoped></style>