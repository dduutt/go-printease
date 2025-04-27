import { invokeWithMessage, invoke } from ".";
import {
  Create,
  CreateMany,
  FindByBatchCode,
  FindByBatchCodes,
} from "../../wailsjs/go/internal/PrintRecord";

export const printRecordAPI = {
  create,
  findByBatchCode,
  findByBatchCodes,
  createMany,
  genBatchCodes,
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

async function findByBatchCodes(batchCodes) {
  return await invoke(FindByBatchCodes, batchCodes);
}

async function createMany(formData, codes = []) {
  const { batchCode, printTemplate, printer, num, runningNumberLength } =
    formData;
  const { fields, name } = printTemplate;
  const records = [];
  codes.forEach((code) => {
    records.push({
      batch_code: code,
      template_name: name,
      fields,
      printer: printer,
    });
  });
  return await invoke(CreateMany, records);
}

function genBatchCodes(batch, runningNumberLength, runningNumberCounter, num) {
  const batchCodes = [];
  for (let i = 0; i < num; i++) {
    const runningNumber = (runningNumberCounter + i)
      .toString()
      .slice(-runningNumberLength)
      .padStart(runningNumberLength, "0");
    const batchCode = `${batch}${runningNumber}`;
    batchCodes.push(batchCode);
  }
  return batchCodes;
}
