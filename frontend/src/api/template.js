import { InsertOne, List, QueryByName } from "../../wailsjs/go/internal/Template";

export const templateAPI = {
  add,
  update,
  getList,
  getByName,
};

async function add(template) {
  console.log(template);
  const res = await InsertOne(template);
  console.log(res);
  return [];
}

async function getByName(name) {
  return [];
}

async function update(template) {
  return [];
}

async function getList(searchText = "", currentPage = 1, pageSize = 10) {
  const res = await List(searchText, currentPage, pageSize);
  console.log(res);
  return res;
}
