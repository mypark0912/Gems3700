import path from 'path';
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  define: {
    'process.env': process.env
  },
  plugins: [vue()],
  resolve: {
    alias: {
      '@tailwindConfig': path.resolve(__dirname, 'tailwind.config.js'),
      '@': path.resolve(__dirname, 'src'),
    },
  },
  optimizeDeps: {
    include: [
      '@tailwindConfig',
      'pinia' // ✅ Pinia를 미리 번들링하도록 추가
    ],
    exclude: [] // ✅ Pinia가 번들에서 빠지지 않도록 설정
  }, 
  css: {
    devSourcemap: true,
  },
  server: {
    hmr: {
      overlay: true,  // 에러 오버레이 표시
      port: 24678,    // HMR 전용 포트
    },
    // CSS 변경 감지 개선
    watch: {
      usePolling: true, // 파일 변경 감지 개선
      interval: 100,
    }
  }, 
  build: {
    rollupOptions: {
      external: [], // ✅ Pinia가 외부 모듈로 빠지지 않도록 설정,
      output: {
        manualChunks: {
          vue: ['vue'],
          vendor: ['axios', 'moment'], // 많이 쓰는 라이브러리 묶기
          'vendor-css': ['tailwindcss']
        },
        assetFileNames: (assetInfo) => {
          if (assetInfo.name && assetInfo.name.endsWith('.css')) {
            return 'assets/css/[name]-[hash][extname]';
          }
          return 'assets/[name]-[hash][extname]';
        }
      },
    },
    commonjsOptions: {
      transformMixedEsModules: true,
    },
    // 🔧 CSS 최적화 설정
    cssMinify: 'esbuild',
    
    // 🔧 빌드 캐시 설정
    target: 'esnext',
    minify: 'esbuild',
  }
});
