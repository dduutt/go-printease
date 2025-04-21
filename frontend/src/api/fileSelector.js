import { GetSelectedPath } from "../../wailsjs/go/main/App";

export async function fileSelector(displayName, ext) {
  return await GetSelectedPath(displayName, ext);
}
