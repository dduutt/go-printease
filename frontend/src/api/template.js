import { withMessage } from ".";
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
  return await withMessage(ListByName, name, offset, PageSize);
}

async function create(template) {
  return await withMessage(Create, template);
}

async function update(template) {
  return await withMessage(Update, template);
}

async function deleteById(id = "") {
  return await withMessage(Delete, id);
}

async function findDatasByKeys(id, keys) {
  try {
    const r = await FindDatasByKeys(id, [...keys]);
    console.log(id, keys);
    console.log(r);
    return {
      data: r,
      status: true,
    };
  } catch (err) {
    console.log(err);
    return {
      status: false,
      data: [],
    };
  }
}
