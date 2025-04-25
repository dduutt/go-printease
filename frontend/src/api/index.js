export * from "./template";
export * from "./file";
export * from "./printer";
export * from "./printRecord";

export async function invokeWithMessage(f = async () => {}, ...args) {
  const r = await invoke(f, ...args);
  if (r.status) {
    ElMessage.success(r.message);
  } else {
    ElMessage.error(r.message);
  }
  return r;
}

export async function invokeWithLoading(f = async () => {}, ...args) {
  const loading = ElLoading.service({
    lock: true,
    text: "数据加载中...",
    background: "rgba(0, 0, 0, 0.7)",
  });
  const r = await invokeWithMessage(f, ...args);
  loading.close();
  return r;
}

export async function invoke(f = async () => {}, ...args) {
  try {
    const result = await f(...args);
    return {
      data: result,
      status: true,
      message: "操作成功",
    };
  } catch (error) {
    console.error(error);
    return {
      data: "",
      status: false,
      message: String(error),
    };
  }
}
