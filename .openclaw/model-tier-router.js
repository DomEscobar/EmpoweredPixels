/**
 * Model Tier Router - Routes agent tasks to appropriate model based on cost tier
 * Usage: const model = require('./model-tier-router').getModel(agentId, taskComplexity);
 */

const TIER_MODELS = {
    basic: 'openrouter/google/gemini-3-flash-preview',
    standard: 'vercel-ai-gateway/zai/glm-4.7',
    advanced: 'openrouter/anthropic/claude-opus-4.6'
};

const AGENT_TIERS = {
    supervisor: 'basic',
    player_casual: 'basic',
    player_hardcore: 'basic',
    releaser: 'basic',
    
    coder: 'standard',
    foundry: 'standard',
    forge: 'standard',
    guardian: 'standard',
    balancer: 'standard',
    creator: 'standard',
    
    main: 'advanced',
    analyst: 'advanced'
};

function getModel(agentId, taskComplexity = 'normal') {
    // Override with task complexity if specified
    if (taskComplexity === 'critical' || taskComplexity === 'strategic') {
        return TIER_MODELS.advanced;
    }
    if (taskComplexity === 'simple') {
        return TIER_MODELS.basic;
    }
    
    // Default to agent's tier
    const tier = AGENT_TIERS[agentId] || 'standard';
    return TIER_MODELS[tier];
}

function getTierForAgent(agentId) {
    return AGENT_TIERS[agentId] || 'standard';
}

module.exports = {
    getModel,
    getTierForAgent,
    TIER_MODELS,
    AGENT_TIERS
};