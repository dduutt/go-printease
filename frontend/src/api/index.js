export * from "./template";
export * from "./file";
export * from "./printer";

export async function withMessage(f = async () => {}, ...args) {
  const loading = ElLoading.service({
    lock: true,
    text: "数据加载中...",
    background: "rgba(0, 0, 0, 0.7)",
  });
  try {
    const result = await f(...args);
    console.log("withMessage result", result);
    ElMessage.success({
      message: "操作成功",
    });
    return {
      data: result,
      status: true,
    };
  } catch (error) {
    console.error("withMessage error", error);
    ElMessage.error({
      message: "操作失败",
    });
    return {
      data: "",
      status: false,
    };
  } finally {
    loading.close();
  }
}
