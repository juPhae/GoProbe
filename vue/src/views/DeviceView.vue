<template>
  <div >
    <div class="status-cards-wrapper">
      <div  style="margin-bottom:10px"  v-for="(systemStatus, index) in systemStatusData"
      :key="index">
      <SystemStatusCard
        :systemStatus="systemStatus"
      />
    <el-button>启动</el-button>
    <el-button>停止</el-button>
    </div>
    </div>
  </div>
 
</template>

<script>
import SystemStatusCard from "@/components/SystemStatusCard.vue";

export default {
  components: {
    SystemStatusCard,
  },
  data() {
    return {
      systemStatusData: [
        {
          device: "MI-Book",
          cpu: 0.7692307692307693,
          memory: 84,
          disk: 59.61408510745181,
          net_in: 5133058.1328125,
          net_out: 443597.2890625,
        },
        {
          device: "MacBook Pro",
          cpu: 0.5692307692307693,
          memory: 64,
          disk: 40.61408510745181,
          net_in: 3133058.1328125,
          net_out: 243597.2890625,
        },
        {
          device: "Surface Pro",
          cpu: 0.6692307692307693,
          memory: 72,
          disk: 49.61408510745181,
          net_in: 4133058.1328125,
          net_out: 343597.2890625,
        },
      ],
    };
  },
};
</script>

<style scoped>
.status-cards-wrapper {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}
</style>
