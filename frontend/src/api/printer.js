import { invoke, invokeWithLoading, invokeWithMessage } from ".";
import { GetPinters } from "../../wailsjs/go/main/App";
export const printerAPI = {
  list,
  print,
};

const url =
  "http://localhost:5159/api/actions?Wait=30s&MessageCount=200&MessageSeverity=Info";

async function list() {
  return await invoke(GetPinters);
}

async function print(formData) {
  const { batchCode, copies, printer, num, printTemplate, date } = formData;
  let labelTemplatePath = formData.labelTemplatePath;
  const data = {};
  printTemplate?.fields?.forEach((item) => {
    data[item.key] = item.value;
  });

  data["batchCode"] = batchCode;
  data["date"] = date;

  // 判断是否打印3c标志版本
  if (data["3C"] == "是") {
    // 修改标签模板路径
    labelTemplatePath = labelTemplatePath.replace(".btw", "_3c.btw");
  }
  // 构建 Bartender REST API 请求数据
  const printRequestData = {
    PrintBTWAction: {
      Document: labelTemplatePath,
      Printer: printer,
      SaveAfterPrint: false,
      NamedDataSources: data,
      Copies: copies,
    },
  };

  if (printer.toUpperCase().includes("PDF")) {
    printRequestData.PrintBTWAction.PrintToFileFolder =
      labelTemplatePath.substring(0, labelTemplatePath.lastIndexOf("\\"));
    printRequestData.PrintBTWAction.PrintToFileFileName =
      "PrintByPrintBTWAction.pdf";
    printRequestData.PrintBTWAction.Printer = "PDF";
  }

  // 发送打印请求到 Bartender REST API
  const printResp = await fetch(url, {
    method: "POST",
    body: JSON.stringify(printRequestData),
    credentials: "include",
  }).then((response) => response.json());
  if (printResp.Status == "Faulted") {
    throw new Error("打印失败");
  }
  return printResp;
}
