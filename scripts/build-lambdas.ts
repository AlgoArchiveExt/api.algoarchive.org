// scripts/build-lambdas.ts
import { build } from "esbuild";
import { readdirSync, statSync } from "fs";
import path from "path";

const lambdaSrcDir = path.resolve("src/lambda");
const lambdaDistDir = path.resolve("dist/lambda");

function findHandlers(dir: string): string[] {
  return readdirSync(dir)
    .map((name) => path.join(dir, name))
    .filter((p) => statSync(p).isDirectory())
    .map((dir) => path.join(dir, "index.ts"))
    .filter((p) => statSync(p, { throwIfNoEntry: false }));
}

const handlers = findHandlers(lambdaSrcDir);

Promise.all(
  handlers.map((entry: string) => {
    const rel = path.relative(lambdaSrcDir, entry);
    const outFile = path.join(
      lambdaDistDir,
      rel.replace(/\.ts$/, ".js")
    );

    return build({
      entryPoints: [entry],
      bundle: true,
      platform: "node",
      format: "cjs",
      target: "node16",
      outfile: outFile,
      sourcemap: true,
      external: [],
      absWorkingDir: path.dirname(entry)
    }).then(() => {
      console.log(`✅ Bundled ${rel}`);
    });
  })
);

console.log("🎉 All Lambdas bundled successfully!");