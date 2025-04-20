<template>
  <div class="template-manage">
    <div class="crud-header">
      <el-row :gutter="20">
        <el-col :span="18">
          <el-input
            v-model="searchText"
            placeholder="请输入搜索关键字"
            clearable
            @change="handleSearch"
          />
        </el-col>
        <el-col :span="6">
          <div class="flex-end">
            <el-button type="primary" @click="openDialog('create')"
              >新增模板</el-button
            >
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 数据表格 -->
    <el-table
      :data="tableData"
      stripe
      border
      table-layout="auto"
      style="width: 100%"
    >
      <el-table-column
        label="序号"
        align="center"
        prop="rowNum"
      ></el-table-column>
      <el-table-column
        prop="name"
        label="模板名称"
        align="center"
      ></el-table-column>
      <el-table-column
        prop="description"
        label="模板描述"
        align="center"
      ></el-table-column>
      <el-table-column
        prop="createdAt"
        label="创建时间"
        align="center"
      ></el-table-column>
      <el-table-column label="启用" align="center">
        <template #default="scope">
          <el-switch
            v-model="scope.row.inUse"
            :active-value="1"
            :inactive-value="0"
            @change="switchChange(scope.row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center">
        <template #default="scope">
          <el-button size="small" @click="openDialog('edit', scope.row)"
            >编辑</el-button
          >
          <el-button size="small" type="danger" @click="handleDelete(scope.row)"
            >删除</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      class="pagination"
      background
      :hide-on-single-page="true"
      :default-page-size="5"
      :total="total"
      :page-size="pageSize"
      :current-page="currentPage"
      @current-change="handlePageChange"
    />
    <!-- 编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="60%">
      <el-form :model="formData" :rules="formRules" ref="formRef">
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="模板路径" prop="path">
          <el-button
            text
            type="primary"
            @click="openFileSelector"
            :disabled="dialogType === 'edit'"
          >
            {{ formData.path || "选择模板文件" }}
          </el-button>
        </el-form-item>
        <el-form-item label="模板描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入模板描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from "vue";
import { fileAPI, templateAPI } from "../api/";

const defaultFormData = {
  id: null,
  name: "",
  path: "",
  description: "",
};

const searchText = ref("");
// 表格数据
const tableData = ref([]);

// 分页配置
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(5);

// 对话框控制
const dialogTitle = computed(() => {
  return dialogType.value === "create" ? "新增模板" : "编辑模板";
});
const dialogVisible = ref(false);
const dialogType = ref("create");
const formRef = ref("formRef");
const formData = reactive({ ...defaultFormData });

// 表单验证规则
const formRules = {
  name: [{ required: true, message: "请输入模板名称", trigger: "blur" }],
  path: [{ required: true, message: "请选择模板路径", trigger: "blur" }],
  description: [{ max: 100, message: "描述最多100个字符", trigger: "blur" }],
};

// 打开对话框
function openDialog(type, row) {
  dialogType.value = type;
  type === "edit"
    ? Object.assign(formData, row)
    : Object.assign(formData, { ...defaultFormData });
  dialogVisible.value = true;
}
// 切换启用状态
async function switchChange(row) {
  const r = await templateAPI.update(row);
  if (!r.status) {
    row.inUse = 1 - row.inUse; // 切换状态失败，恢复原状态
  }
}
// 提交表单
async function submitForm() {
  const ok = await formRef.value.validate();
  if (!ok) return;
  let r;
  if (dialogType.value === "edit") {
    // 编辑操作
    r = await templateAPI.update(formData);
  } else {
    // 新增操作
    r = await templateAPI.create(formData);
  }

  if (r.status) {
    // 这里添加实际提交逻辑
    dialogVisible.value = false;
    // 更新表格数据
    getTemplateList();
  }
}

// 删除操作
async function handleDelete(row) {
  // 这里添加实际删除逻辑
  const r = await templateAPI.deleteById(row.id);
  if (r.status) {
    getTemplateList();
  }
}

function handleSearch() {
  // 这里添加搜索逻辑
  getTemplateList();
}

// 打开文件选择器
async function openFileSelector() {
  const r = await fileAPI.selector("请选择模板文件", "Xlsx", "*.xlsx");
  if (r.status) {
    formData.path = r.data;
  }
}

// 获取模板列表
async function getTemplateList() {
  const r = await templateAPI.listByName(
    searchText.value,
    currentPage.value,
    pageSize.value
  );
  console.log(r);
  if (r.status) {
    tableData.value = r.data.list;
    total.value = r.data.total;
  }
}

// 分页变化
function handlePageChange(page) {
  currentPage.value = page;
  // 这里添加分页请求逻辑
  getTemplateList();
}

// 组件挂载后获取数据
onMounted(() => {
  getTemplateList();
});
</script>

<style scoped>
.crud-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  width: 100%;
}

.pagination {
  margin-top: 20px;
}

.el-button + .el-button {
  margin-left: 10px;
}
</style>

<style scoped>
.template-manage {
  width: 100%;
}
</style>
