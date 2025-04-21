import { withMessage } from ".";
import { GetPinters } from "../../wailsjs/go/main/App";
export const printerAPI = {
  list,
  print,
};

const url =
  "http://localhost:5159/api/actions?Wait=30s&MessageCount=200&MessageSeverity=Info";

async function list() {
  try {
    const r = await GetPinters();
    return {
      data: r,
      status: true,
    };
  } catch (e) {
    return {
      data: [],
      status: false,
    };
  }
}

async function print(formData) {
  const { batchCode, copies, printer, labelTemplatePath, num, printTemplate } =
    formData;
  const data = {};
  printTemplate?.fields?.forEach((item) => {
    data[item.key] = item.value;
  });

  data["batchCode"] = batchCode;

  // 构建 Bartender REST API 请求数据
  const printRequestData = {
    PrintBTWAction: {
      Document: labelTemplatePath,
      Printer: printer,
      SaveAfterPrint: false,
      NamedDataSources: data,
      Copies: copies,
      // PrintToFileFolder: labelTemplatePath.substring(
      //   0,
      //   labelTemplatePath.lastIndexOf("\\")
      // ),
      // PrintToFileFileName: "PrintByPrintBTWAction.pdf",
    },
  };

  // 发送打印请求到 Bartender REST API

  return await withMessage(fetch, url, {
    method: "POST",
    body: JSON.stringify(printRequestData),
    credentials: "include",
  });
}
