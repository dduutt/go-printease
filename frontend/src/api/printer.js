import { invoke, invokeWithLoading, invokeWithMessage } from ".";
import { GetPinters } from "../../wailsjs/go/main/App";
export const printerAPI = {
  list,
  printMany,
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
// 批量打印多个标签
async function printMany(formData, codes = []) {
  const { batchCode, copies, printer, num, printTemplate, date } = formData;
  let labelTemplatePath = formData.labelTemplatePath;
  const baseData = {};
  printTemplate?.fields?.forEach((item) => {
    baseData[item.key] = item.value;
  });
  baseData["date"] = date;
  // 判断是否打印3c标志版本
  if (baseData["3C"] == "是") {
    // 修改标签模板路径
    labelTemplatePath = labelTemplatePath.replace(".btw", "_3c.btw");
  }

  // PDF打印设置
  const pdfSettings = {};
  if (printer.toUpperCase().includes("PDF")) {
    pdfSettings.PrintToFileFolder = labelTemplatePath.substring(
      0,
      labelTemplatePath.lastIndexOf("\\")
    );
    pdfSettings.PrintToFileFileName = "PrintByPrintBTWAction.pdf";
    pdfSettings.Printer = "PDF";
  }

  const printRequestDatas = [];
  // 为每个批次码创建独立的数据对象和打印请求
  codes.forEach((code) => {
    // 创建新的数据对象，避免引用同一个对象
    const data = { ...baseData, batchCode: code };

    // 构建独立的打印请求数据
    const printRequestData = {
      PrintBTWAction: {
        Document: labelTemplatePath,
        Printer: printer,
        SaveAfterPrint: false,
        NamedDataSources: data,
        Copies: copies,
        ...pdfSettings,
      },
    };

    printRequestDatas.push(printRequestData);
  });

  // 发送打印请求到 Bartender REST API
  const printResp = await fetch(url, {
    method: "POST",
    body: JSON.stringify(printRequestDatas),
    credentials: "include",
  }).then((response) => response.json());
  if (printResp.Status == "Faulted") {
    throw new Error("打印失败");
  }
  return printResp;
}
