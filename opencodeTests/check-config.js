const configPath = process.env.AGENCY_CONFIG || './agency-config.json';
const config = require(configPath);
console.log('Coder agent OPENCODE_MODEL:', config.agents.coder.env.OPENCODE_MODEL);
console.log('AGENT_WORKDIR:', config.agents.coder.env.AGENT_WORKDIR);
