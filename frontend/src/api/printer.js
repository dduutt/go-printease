import { GetPinters } from "../../wailsjs/go/main/App";

export const printer = {
  list,
};

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
