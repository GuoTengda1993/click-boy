<template>
  <div class="common-layout">
    <el-container>
      <el-header>
        <el-menu class="el-menu-demo" mode="horizontal" :ellipsis="false" :router="true">
          <el-menu-item index="">
            <img style="width: 45px" alt="ClickBoy" src="./assets/images/logo-universal.png"/>
          </el-menu-item>
        </el-menu>
      </el-header>
      <el-main>
        <Click/>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import Click from './components/Click.vue'
</script>

<style>
.el-header {
    padding: 0;
    top: 0;
    width: 100%;
    position: fixed;
    z-index: 1000;
}

.el-main {
    margin-top: 60px;
}
</style>
