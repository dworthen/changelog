import path from "node:path";

export default function run() {
  const dirname = import.meta.dirname;
  const bin = path.join(dirname, "bin", "changelog.exe");
  return bin;
}
