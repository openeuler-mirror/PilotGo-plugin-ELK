<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Thu Jun 20 10:05:48 2024 +0800
-->
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