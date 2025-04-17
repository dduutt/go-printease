<template>
    <div class="header-content">
        <div class="logo">
            <img src="../assets/vue.svg" alt="Logo" style="height: 36px">
            <span>PrintEase</span>
        </div>
        <el-menu mode="horizontal" router :default-active="activeIndex" background-color="#fff" text-color="#606266"
            active-text-color="#409EFF" class="nav-menu">
            <el-menu-item index="/label-print">
                <router-link to="/label-print">标签打印</router-link>
            </el-menu-item>
            <el-menu-item index="/print-history">
                <router-link to="/print-history">打印记录</router-link>
            </el-menu-item>
            <el-menu-item index="/template-manage">
                <router-link to="/template-manage">模板管理</router-link>
            </el-menu-item>
        </el-menu>
    </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const activeIndex = computed(() => route.path);
</script>

<style scoped>

.header-content {
    display: flex;
    align-items: center;
    max-width: 1200px;
    margin: 0 auto;
}

.logo {
    display: flex;
    align-items: center;
    margin-right: 40px;
    font-size: 20px;
    font-weight: bold;
    color: #303133;
}

.logo img {
    margin-right: 10px;
}

.nav-menu {
    flex: 1;
    border-bottom: none;
}

.nav-menu :deep(.el-menu-item) {
    font-size: 16px;
    padding: 0 20px !important;
    height: 60px;
    line-height: 60px;
    transition: all 0.3s;
}

.nav-menu :deep(.el-menu-item.is-active) {
    border-bottom: 2px solid #409EFF;
}


a {
    text-decoration: none;
    color: inherit;
}
</style>