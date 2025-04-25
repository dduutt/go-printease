import { FileSelector } from "../../wailsjs/go/main/App";
import { invoke } from ".";

export const fileAPI = {
  selector,
};

async function selector(title, displayName, ext) {
  return await invoke(FileSelector, title, displayName, ext);
}
