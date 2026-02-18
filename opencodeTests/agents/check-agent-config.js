const path = require('path');
const baseDir = __dirname;
const configPath = path.resolve(baseDir, '../agency-config.json');
const config = require(configPath).agents.coder;
console.log('Coder agent env:', config.env);
console.log('OPENCODE_MODEL:', config.env.OPENCODE_MODEL);
