/**
 * Local Lambda runner for CDK TypeScript projects
 *
 * Usage:
 *   Run 'npm run build-test-lambdas' first to build the Lambda functions.
 *   Then, you can run a specific Lambda function with:
 *   npm run test-lambda -- dist/lambda/commit/index.js handler test/commit/event.json
 */

import path from "path";
import fs from "fs";

// Get CLI args
const [, , lambdaPath, handlerName, eventArg] = process.argv;

if (!lambdaPath || !handlerName) {
  console.error(
    "Usage: npm run local-lambda -- <lambdaPath> <handlerName> [eventFileOrJson]"
  );
  process.exit(1);
}

// Resolve the Lambda file path
const fullPath = path.resolve(lambdaPath);

// Load event (either from JSON file or inline JSON string)
let event: any = {};
if (eventArg) {
  if (fs.existsSync(eventArg)) {
    event = JSON.parse(fs.readFileSync(eventArg, "utf8"));
  } else {
    try {
      event = JSON.parse(eventArg);
    } catch {
      console.error("Invalid event JSON:", eventArg);
      process.exit(1);
    }
  }
}

(async () => {
  try {
    // Dynamically import the Lambda module
    const lambdaModule = await import(fullPath);

    if (!lambdaModule[handlerName]) {
      throw new Error(
        `Handler "${handlerName}" not found in module ${fullPath}. ` +
          `Available exports: ${Object.keys(lambdaModule).join(", ")}`
      );
    }

    console.log(`Invoking ${handlerName} from ${fullPath}...`);
    const result = await lambdaModule[handlerName](event, {}, () => {});
    console.log("Lambda result:", result);
  } catch (err) {
    console.error("Error running Lambda:", err);
    process.exit(1);
  }
})();