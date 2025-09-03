<template>
  <div class="bg-white dark:bg-gray-800 shadow-sm rounded-xl px-5 py-4">
      <div class="md:flex justify-between items-center space-y-4 md:space-y-0 space-x-2">
          <!-- Left side -->
          <div class="flex items-start space-x-3 md:space-x-4">
              <div>
                  <a href="#" @click.prevent="handleTitleClick" class="inline-flex font-semibold text-gray-800 dark:text-gray-100">
                    {{ item.title }}
                  </a>
                  <div class="text-sm">{{item.context}} </div>
                  <div class="text-sm">{{ getVersionText(item.utype, item.a_version, item.w_version, item.f_version) }}</div>
              </div>
          </div>
          <!-- Right side -->
          <div class="flex items-center space-x-4 pl-10 md:pl-0">
              <div class="text-sm text-gray-500 dark:text-gray-400 italic whitespace-nowrap">{{item.date}}</div>
              <div v-if="item.utype" class="text-xs inline-flex font-medium rounded-full text-center px-2.5 py-1 bg-green-500/20 text-green-700">{{item.utype}}</div>
              <svg v-if="item.mtype == 0" class="shrink-0 mr-2 stroke-current text-green-500 transition-colors" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                <polyline points="7 10 12 15 17 10" />
                <line x1="12" y1="15" x2="12" y2="3" />
              </svg>
              <svg v-else-if="item.mtype == 1" class="shrink-0 mr-2 stroke-current text-yellow-500 transition-colors" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="23 4 23 10 17 10" />
                <polyline points="1 20 1 14 7 14" />
                <path d="M3.51 9a9 9 0 0 1 14.13-3.36L23 10" />
                <path d="M20.49 15a9 9 0 0 1-14.13 3.36L1 14" />
            </svg>
            <svg v-else class="shrink-0 mr-2 stroke-current text-pink-500 transition-colors" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="16 3 21 3 21 8" />
                <line x1="4" y1="20" x2="21" y2="3" />
                <polyline points="8 21 3 21 3 16" />
            </svg>
          </div>
      </div>
  </div>
</template>

<script>

export default {
  name: 'maintenanceItem',
  props: ['item'],
  setup(props, { emit }){
    function getVersionText(utype, appVer, webVer, fwVer) {
      if (!utype) return '';

      const result = [];
      const types = utype.split(',').map(s => s.trim().toLowerCase());

      if (types.includes('app') && appVer) {
        result.push(`App Version : ${appVer}`);
      }
      if (types.includes('web') && webVer) {
        result.push(`Web Version : ${webVer}`);
      }
      if (types.includes('fw') && fwVer) {
        result.push(`FW Version : ${fwVer}`);
      }

      return result.join(' / ');
    }

    function handleTitleClick() {
      emit('editTitle', props.item);
    }

    return {
      getVersionText,
      handleTitleClick,
    }
  }
}
</script>