/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Thu Jun 20 10:05:48 2024 +0800
 */
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { resolve } from 'path';

import basicSsl from '@vitejs/plugin-basic-ssl'

// https://vitejs.dev/config/
export default defineConfig({
  // base: '/plugin/modeule_name/', // 生产环境下的公共路径
  server: {
    host:'10.41.107.33',
    port: 8080,
    strictPort: false, // 设为 true 时若端口已被占用则会直接退出，而不是尝试下一个可用端口
    cors: true,
    open: true, // 启动后是否浏览器自动打开
    hmr: true, // 为开发服务启用热更新，默认是不启用热更新的
    proxy: {
      '/plugin/elk/api': {
        target: 'https://10.41.107.29:9993',
        secure: false,
        changeOrigin: true,
        rewrite: path => path.replace(/^\//, '')
      },
    },
  },
  plugins: [vue(),basicSsl()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@assetImg':resolve(__dirname,'src/assets/images') 
    },
  },
  build: {
    outDir: 'dist',
    assetsDir: 'static/', //指定静态资源存放路径
    sourcemap: false, //是否构建source map 文件
    chunkSizeWarningLimit: 1500, //chunk 大小警告的限制，默认500KB
  },
});
