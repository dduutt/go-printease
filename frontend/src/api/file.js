import { FileSelector } from "../../wailsjs/go/main/App";
import { withMessage } from ".";
export const fileAPI = {
  selector,
};

async function selector(title, displayName, ext) {
  return await withMessage(FileSelector, title, displayName, ext);
}
