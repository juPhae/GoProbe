<template>
  <el-card>
    <div slot="header" class="clearfix">
      <span>系统状态</span>
    </div>
    <div class="card-body">
      <div class="card-item">
        <i class="el-icon-monitor"></i>
        <div class="card-value">{{ systemStatus.device }}</div>
        <div class="card-label">设备名称</div>
      </div>
      <div class="card-item">
        <i class="el-icon-cpu"></i>
        <div class="card-value">{{ systemStatus.cpu }}%</div>
        <div class="card-label">CPU 使用率</div>
      </div>
      <div class="card-item">
        <i class="el-icon-memory"></i>
        <div class="card-value">{{ systemStatus.memory }}%</div>
        <div class="card-label">内存使用率</div>
      </div>
      <div class="card-item">
        <i class="el-icon-sold-out"></i>
        <div class="card-value">{{ systemStatus.disk }}%</div>
        <div class="card-label">磁盘使用率</div>
      </div>
      <div class="card-item">
        <i class="el-icon-arrow-down"></i>
        <div class="card-value">{{ systemStatus.net_in | formatBytes }}/s</div>
        <div class="card-label">网络下载速度</div>
      </div>
      <div class="card-item">
        <i class="el-icon-arrow-up"></i>
        <div class="card-value">{{ systemStatus.net_out | formatBytes }}/s</div>
        <div class="card-label">网络上传速度</div>
      </div>
    </div>
  </el-card>
</template>

<script>
export default {
  name: 'SystemStatusCard',
  props: {
    systemStatus: {
      type: Object,
      required: true,
    },
  },
  filters: {
    // 格式化字节数
    formatBytes(bytes) {
      if (bytes < 1024) {
        return bytes.toFixed(2) + 'B';
      } else if (bytes < 1024 * 1024) {
        return (bytes / 1024).toFixed(2) + 'KB';
      } else if (bytes < 1024 * 1024 * 1024) {
        return (bytes / 1024 / 1024).toFixed(2) + 'MB';
      } else {
        return (bytes / 1024 / 1024 / 1024).toFixed(2) + 'GB';
      }
    },
  },
};
</script>
  
  <style scoped>
  .card-body {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    grid-gap: 20px;
    justify-items: center;
    align-items: center;
    padding: 20px;
  }
  
  .card-item {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .card-item i {
    font-size: 48px;
    margin-bottom: 10px;
  }
  
  .card-value {
    font-size: 24px;
    font-weight: bold;
    margin-bottom: 5px;
  }
  
  .card-label {
    font-size: 14px;
    color: #999;
  }
  
  </style>