export async function withMessage(wm = {}) {
  const {
    successMessage = "操作成功",
    errorMessage = "操作失败",
    successType = "success",
    errorType = "error",
    f = async () => {},
  } = wm;

  try {
    const result = await f();
    ElMessage({
      message: successMessage,
      type: successType,
    });
    return {
      data: result,
      status: true,
    };
  } catch (error) {
    ElMessage({
      message: errorMessage,
      type: errorType,
    });
    console.error(`${f.name} error:`, error);
    return {
      data: "",
      status: false,
    };
  }
}
