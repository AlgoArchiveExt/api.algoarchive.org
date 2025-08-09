# Make sure you run this where the body.json file is located
# - Anthony
node -e 'const obj=require("./body.json");console.log(JSON.stringify(obj).replace(/"/g,"\\\""))'