import tseslint from "@typescript-eslint/eslint-plugin";

console.log("Type of config:", typeof tseslint.configs["flat/recommended"]);
console.log("Is array:", Array.isArray(tseslint.configs["flat/recommended"]));
console.log("Keys:", Object.keys(tseslint.configs["flat/recommended"]));
console.log("Value:", JSON.stringify(tseslint.configs["flat/recommended"], null, 2));
