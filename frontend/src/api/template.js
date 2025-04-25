import { invokeWithLoading, invoke } from ".";
import {
  Create,
  ListByName,
  Update,
  Delete,
  FindDatasByKeys,
} from "../../wailsjs/go/internal/Template";

export const templateAPI = {
  listByName,
  create,
  update,
  deleteById,
  findDatasByKeys,
};

async function listByName(name = "", currentPage = 1, PageSize = 0) {
  const offset = (currentPage - 1) * PageSize;
  return await invokeWithLoading(ListByName, name, offset, PageSize);
}

async function create(template) {
  return await invokeWithLoading(Create, template);
}

async function update(template) {
  return await invokeWithLoading(Update, template);
}

async function deleteById(id = "") {
  return await invokeWithLoading(Delete, id);
}

async function findDatasByKeys(id, keys) {
  return await invoke(FindDatasByKeys, id, [...keys]);
}
