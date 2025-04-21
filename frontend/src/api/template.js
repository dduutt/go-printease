import { withMessage } from ".";
import {
  Create,
  ListByName,
  Update,
  Delete,
} from "../../wailsjs/go/internal/Template";

export const templateAPI = {
  listByName,
  create,
  update,
  deleteById,
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
  console.log("deleteById", id);
  return await withMessage(Delete, id);
}

