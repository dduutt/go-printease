export * from "./template";

export async function withMessage(f = async () => {}, ...args) {
  try {
    const result = await f(...args);
    console.log(result);
    ElMessage.success({
      message: "操作成功",
    });
    return {
      data: result,
      status: true,
    };
  } catch (error) {
    ElMessage.error({
      message: "操作失败",
    });
    return {
      data: "",
      status: false,
    };
  }
}
