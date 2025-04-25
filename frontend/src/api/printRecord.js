import { invokeWithMessage, invoke } from ".";
import { Create, FindByBatchCode } from "../../wailsjs/go/internal/PrintRecord";

export const printRecordAPI = {
  create,
  findByBatchCode,
};

async function create(formData) {
  const { batchCode, printTemplate, printer } = formData;
  const { fields, name } = printTemplate;
  const record = {
    batch_code: batchCode,
    template_name: name,
    fields,
    printer: printer,
  };
  return await invoke(Create, record);
}

async function findByBatchCode(batchCode) {
  return await invoke(FindByBatchCode, batchCode);
}
