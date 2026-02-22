import globals from "globals";
import tseslint from "@typescript-eslint/eslint-plugin";
import vueEslint from "eslint-plugin-vue";
import vueParser from "vue-eslint-parser";

export default [
  // Global ignores
  { ignores: ["dist/**", "node_modules/**", ".git/**"] },

  // Vue files with TypeScript parser
  {
    files: ["**/*.vue"],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        ecmaVersion: 2022,
        sourceType: "module",
        parser: {
          ts: "@typescript-eslint/parser",
          js: "espree",
          '<template>': "espree"
        },
        extraFileExtensions: [".vue"],
        globals: {
          ...globals.browser,
          ...globals.node
        }
      }
    }
  },

  // TypeScript recommended (relaxed)
  ...tseslint.configs["flat/recommended"].map(config => ({
    ...config,
    rules: {
      ...config.rules,
      "@typescript-eslint/no-unused-vars": "off",
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/no-invalid-void-type": "off"
    }
  })),

  // Vue recommended (relaxed)
  ...vueEslint.configs["flat/recommended"].map(config => ({
    ...config,
    rules: {
      ...config.rules,
      "vue/multi-word-component-names": "off",
      "vue/max-attributes-per-line": "off",
      "vue/html-self-closing": "off",
      "vue/require-default-prop": "off",
      "vue/no-parsing-error": "warn"
    }
  })),

  // JavaScript config
  {
    files: ["**/*.js"],
    languageOptions: {
      ecmaVersion: 2022,
      sourceType: "module",
      globals: {
        ...globals.browser,
        ...globals.node
      }
    }
  }
];
