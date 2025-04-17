<template>
  <div class="common-layout">
    <el-form :inline="true" :model="formData" label-position="top" class="printForm">
      <el-row>
        <el-col :span="8">
          <el-form-item label="打印模板" prop="labelTemplate">
            <el-select v-model="formData.labelTemplate" filterable placeholder="请选择打印模板">
              <el-option v-for="item in labelTemplates" :key="item.id" :label="item.name" :value="item.data" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="标签模板" prop="labelTemplate">
            <el-select v-model="formData.labelTemplates" remote :remote-method="searchRemoteLabelData" filterable
              placeholder="Select">
              <el-option v-for="item in labelTemplateDatas" :key="item.id" :label="item.name" :value="item.data" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="索引" prop="labelData">
            <el-select v-model="formData.labelData" remote :remote-method="searchRemoteLabelData" filterable
              placeholder="Select">
              <el-option v-for="item in labelTemplateDatas" :key="item.id" :label="item.name" :value="item.data" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="打印机" prop="printer">
            <el-select v-model="formData.printer" placeholder="请选择打印机">
              <el-option v-for="item in printers" :key="item" :label="item" :value="item" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="批次信息" prop="batch">
            <el-input v-model="formData.batch" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="位数" prop="runningNumberLength">
            <el-input-number v-model="formData.runningNumberLength" :min="1" :max="6" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="流水号" prop="runningNumber">
            <el-input-number v-model="formData.runningNumber" @input="handleRunningNumberInput" :min="1"
              :max="Math.pow(10, formData.runningNumberLength) - 1" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="批次号" prop="batchCode">
            <el-input v-model="batchCode" disabled />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="副本" prop="copies">
            <el-input-number v-model="formData.copies" :min="1" :max="100" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="数量" prop="num">
            <el-input-number v-model="formData.num" :min="1" :max="100" />
          </el-form-item>
        </el-col>
        <el-col :span="4">
          <el-form-item>
            <el-button type="primary" @click="resetLabelData">重置</el-button>
          </el-form-item>
        </el-col>
        <el-col :span="4">
          <el-form-item>
            <el-button type="primary" @click="print">打印</el-button>
          </el-form-item>
        </el-col>
        <el-col :span="8" v-for="item in formData.labelData">
          <el-form-item :label="item.name" :prop="item.key">
            <el-input v-model="item.value" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { printer } from "../api/printer";

const defaultFormData = {
  labelTemplate: "",
  labelData: [
    { name: "规格", key: "规格", value: "100*200mm" },
    {
      name: "数量",
      key: "数量",
      value: 100,
    },
    {
      name: "批次号",
      key: "批次号",
      value: "",
    },
    { name: "有效期至", key: "有效期至", value: "2023-06-01" },
    {
      name: "生产日期",
      key: "生产日期",
      value: "2023-05-01",
    },
    {
      name: "供应商",
      key: "供应商",
      value: "供应商名称",
    },
  ],
  printer: "",
  copies: 1,
  num: 1,
  runningNumber: 1,
  runningNumberLength: 1,
  batch: "",
};

const formData = reactive({ ...defaultFormData });
const labelTemplates = ref([]);
const labelTemplateDatas = ref([]);
const printers = ref([]);
const runningNumberCounter = ref(1);
const runningNumberLength = ref(1);

const batchCode = computed({
  get() {
    return (
      formData.batch +
      runningNumberCounter.value
        .toString()
        .slice(-formData.runningNumberLength)
        .padStart(formData.runningNumberLength, "0")
    );
  },
});
async function print() {
  console.log(formData);
}

async function searchRemoteLabelData(query) {
  console.log("searchRemoteLabelData");
}

async function resetLabelData() {
  Object.assign(formData, defaultFormData);
}

function handleRunningNumberInput(value) {
  console.log(value);
}

onMounted(async () => {
  printers.value = await printer.list();
});
</script>

<style scoped>
.common-layout {
  height: 100%;
  display: flex;
  width: 100%;
}

.printForm .el-form-item {
  width: 90%;
}
</style>
