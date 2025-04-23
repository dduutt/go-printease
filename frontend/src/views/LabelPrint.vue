<template>
  <div class="common-layout">
    <el-form ref="ruleFormRef" :rules="formRules" :inline="true" :model="formData" label-position="top"
      class="printForm">
      <el-row>
        <el-col :span="8">
          <el-form-item label="打印模板" prop="printTemplate">
            <el-select v-model="formData.printTemplate" value-key="id" placeholder="请选择打印模板">
              <el-option v-for="item in printTemplates" :key="item.id" :label="item.name" :value="item" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="标签模板" prop="labelTemplatePath">
            <el-button text type="primary" @click="openFileSelector">
              {{ formData.labelTemplatePath || "选择模板文件" }}
            </el-button>
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
            <el-input-number v-model="formData.runningNumberCounter" :min="1"
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
        <el-col :span="8">
          <el-form-item label=" " prop="actions" class="button-group">
            <div>
              <el-button type="info" @click="resetLabelData(ruleFormRef)">重置</el-button>
              <el-button type="primary" @click="print">打印</el-button>
            </div>
          </el-form-item>
        </el-col>
        <el-divider />
        <el-col :span="8" v-for="item in formData?.printTemplate?.fields">
          <el-form-item :label="item.name" :prop="item.key">
            <el-autocomplete v-model="item.value" clearable
              :fetch-suggestions="(...args) => { querySearchAsync(item.key, ...args) }" placeholder="请输入"
              @select="handleSelect" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { templateAPI, fileAPI, printerAPI } from "../api";

const defaultFormData = {
  printer: "",
  copies: 1,
  num: 1,
  runningNumberCounter: 1,
  runningNumberLength: 3,
  batch: "",
  printTemplate: {},
  labelTemplatePath: "",
};

const formData = reactive({ ...defaultFormData });
const ruleFormRef = ref(null);
const printers = ref([]);
// 打印模板列表
const printTemplates = reactive([]);

const options = reactive({});
const loading = ref({})


// 表单验证规则
const formRules = {
  batchCode: [{ required: true, message: "批次号不能为空", trigger: "blur" }],
  labelTemplatePath: [
    { required: true, message: "请选择标签模板路径", trigger: "blur" },
  ],
  printTemplate: [
    { required: true, message: "请选择打印模板路径", trigger: "blur" },
  ],
  batch: [{ required: true, message: "请选择打印模板路径", trigger: "blur" }],
};

const batchCode = computed({
  get() {
    const b =
      formData.batch +
      formData.runningNumberCounter
        .toString()
        .slice(-formData.runningNumberLength)
        .padStart(formData.runningNumberLength, "0");
    formData.batchCode = b;
    return b;
  },
});
async function print() {
  const ok = await ruleFormRef.value.validate().catch((err) => false);
  if (!ok) {
    return;
  }
  const r = await printerAPI.print(formData);
  if (r.status) {
    formData.runningNumberCounter++;
  }
}
// 过滤keys中key字段的重复项
function filterUniqueValuesByKey(arr, key) {
  const v = arr.map((i) => i[key]);
  const fv = [...new Set(v)];
  return fv.map((i) => ({ value: i }));
}

async function querySearchAsync(key, query, cb) {
  const r = await querySearch();
  if (r.status) {
    const data = filterUniqueValuesByKey(r.data, key);
    cb(data);
  }

}

async function querySearch() {
  const keys = [...formData.printTemplate.fields];
  const nKeys = keys.filter(k => k.value != "");
  return await templateAPI.findDatasByKeys(formData.printTemplate.id, nKeys);
}

// 选择后自动填充
async function handleSelect(item, key) {
  const r = await querySearch();
  if (r.status) {
    const data = r.data
    if (data.length === 1) {
      formData.printTemplate.fields.forEach((i) => {
        i.value = data[0][i.key];
      });
    }
  }

}


async function resetLabelData() {
  formData.printTemplate.fields.forEach((i) => (i.value = ""));
}

async function getTemplates() {
  const r = await templateAPI.listByName();
  if (r.status) {
    Object.assign(printTemplates, r.data.list);
  }
}

async function openFileSelector() {
  const r = await fileAPI.selector("请选择模板文件", "BTW", "*.btw");
  if (r.status) {
    formData.labelTemplatePath = r.data;
  }
}

async function getPrinters() {
  const r = await printerAPI.list();
  if (r.status) {
    printers.value = r.data;
  }
}

onMounted(async () => {
  getPrinters();
  getTemplates();
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

.button-group {
  display: flex;
  justify-content: space-between;
}
</style>
