// scripts/build-lambdas.js
import { build } from "esbuild";
import { readdirSync, statSync } from "fs";
import path from "path";

const lambdaSrcDir = path.resolve("src/lambda");
const lambdaDistDir = path.resolve("dist/lambda");

// Find all Lambda handler entrypoints in src/lambda/
function findLambdaHandlers(dir: string): string[] {
  const handlers: string[] = [];
  for (const item of readdirSync(dir)) {
    const fullPath = path.join(dir, item);
    const stats = statSync(fullPath);

    if (stats.isDirectory()) {
      const indexTs = path.join(fullPath, "index.ts");
      const indexJs = path.join(fullPath, "index.js");

      if (statSync(indexTs, { throwIfNoEntry: false })) {
        handlers.push(indexTs);
      }
    }
  }
  return handlers;
}

const entryPoints = findLambdaHandlers(lambdaSrcDir);

if (entryPoints.length === 0) {
  console.log("⚠️  No Lambda handlers found in src/lambda/");
  process.exit(0);
}

console.log(`📦 Bundling ${entryPoints.length} Lambda(s) with esbuild...`);

await Promise.all(
  entryPoints.map((entry) => {
    const relPath = path.relative(lambdaSrcDir, entry);
    const outFile = path.join(
      lambdaDistDir,
      relPath.replace(/\.ts$/, ".js").replace(/\.js$/, ".js")
    );

    return build({
      entryPoints: [entry],
      bundle: true,
      platform: "node",
      format: "esm", // AWS Lambda ESM
      target: "node20",
      sourcemap: true,
      outfile: outFile,
    }).then(() => {
      console.log(`✅ Built: ${relPath} → ${path.relative(".", outFile)}`);
    });
  })
);

console.log("🎉 All Lambdas bundled successfully!");