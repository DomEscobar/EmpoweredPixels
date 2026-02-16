#!/usr/bin/env node
/**
 * OpenCode Agency — Advanced Director
 * Capability-based routing, workflow orchestration, integration with SQLite queue
 */

const { AgentBase } = require('./agent-base');
const { v4: uuidv4 } = require('uuid');

class DirectorAgent extends AgentBase {
  constructor(config, queue) {
    super(config, queue);
    this.workflows = this.loadWorkflows();
    this.capabilityRegistry = new Map(); // agent -> capabilities
    // Director itself doesn't poll tasks; it's invoked by webhook or manually
    // But we implement execute() to handle "director.plan" tasks
  }

  loadWorkflows() {
    // In production, load from agency-config.json
    // For now, inline same data
    return {
      feature: {
        steps: [
          { agent: 'coder', action: 'implement' },
          { agent: 'reviewer', action: 'review' },
          { agent: 'tester', action: 'validate' }
        ]
      },
      bugfix: {
        steps: [
          { agent: 'coder', action: 'fix' },
          { agent: 'reviewer', action: 'approve' },
          { agent: 'tester', action: 'verify' }
        ]
      }
    };
  }

  async execute(task) {
    this.logger.info({ taskId: task.id, type: task.type }, 'Director executing');

    if (task.type === 'director.plan') {
      return await this.planTask(task);
    } else if (task.type === 'workflow.start') {
      return await this.startWorkflow(task);
    } else if (task.type === 'director.status') {
      return this.reportStatus();
    } else if (task.type === 'director.rerun') {
      return await this.rerunTask(task);
    } else {
      throw new Error(`Unknown task type for director: ${task.type}`);
    }
  }

  async planTask(task) {
    const { description, priority = 0, tags = [] } = task.payload;
    const planId = uuidv4();

    this.logger.info({ planId, description }, 'Creating plan');

    // Simple heuristic: choose agents based on tags
    const subtasks = [
      {
        id: uuidv4(),
        type: 'implement',
        payload: { description, planId },
        required_capability: 'coding',
        priority
      },
      {
        id: uuidv4(),
        type: 'review',
        payload: { planId, description },
        required_capability: 'review',
        priority
      },
      {
        id: uuidv4(),
        type: 'test',
        payload: { planId, description },
        required_capability: 'testing',
        priority
      }
    ];

    // Enqueue subtasks
    for (const subtask of subtasks) {
      const assignedAgent = this.selectAgentForCapability(subtask.required_capability);
      if (!assignedAgent) {
        throw new Error(`No agent available for capability: ${subtask.required_capability}`);
      }
      // Add assignment
      const queueTask = {
        id: subtask.id,
        type: subtask.type,
        payload: subtask.payload,
        priority: subtask.priority,
        assigned_to: assignedAgent
      };
      this.queue.enqueue(queueTask);
      this.logger.info({ taskId: subtask.id, assigned_to: assignedAgent }, 'Enqueued subtask');
    }

    return {
      planId,
      status: 'planned',
      subtasks: subtasks.map(t => t.id)
    };
  }

  async startWorkflow(task) {
    const { workflow, description, priority = 0 } = task.payload;
    const definition = this.workflows[workflow];
    if (!definition) {
      throw new Error(`Unknown workflow: ${workflow}`);
    }

    const workflowId = uuidv4();
    this.logger.info({ workflowId, workflow, description }, 'Starting workflow');

    // Create a task for each step, with dependency chain via parent_id or metadata
    let previousTaskId = null;
    for (const step of definition.steps) {
      const stepTaskId = uuidv4();
      const stepTask = {
        id: stepTaskId,
        type: `${step.agent}.${step.action}`,
        payload: {
          workflowId,
          description,
          step: step.action,
          // Pass along context from previous step if needed
          ...(previousTaskId ? { previous_task_id: previousTaskId } : {})
        },
        priority
        // assigned_to will be step.agent (explicit)
      };
      this.queue.enqueue(stepTask);
      this.logger.info({ taskId: stepTaskId, agent: step.agent, action: step.action }, 'Enqueued workflow step');
      previousTaskId = stepTaskId;
    }

    return { workflowId, status: 'started' };
  }

  selectAgentForCapability(capability) {
    // In full implementation, query the registry for agents advertising this capability
    // For now, simple mapping
    const capabilityMap = {
      'coding': 'coder',
      'review': 'reviewer',
      'testing': 'tester',
      'deployment': 'deployer'
    };
    return capabilityMap[capability] || null;
  }

  async rerunTask(task) {
    const { task_id } = task.payload;
    const original = this.queue.getTask(task_id);
    if (!original) {
      throw new Error(`Task not found: ${task_id}`);
    }
    // Re-enqueue with cleared locks and retry count
    const newTask = { ...original, retry_count: 0, locked_by: null, locked_at: null, scheduled_at: new Date().toISOString() };
    this.queue.enqueue(newTask);
    return { status: 'requeued', task_id };
  }

  reportStatus() {
    const metrics = this.queue.getMetrics();
    return {
      status: 'ok',
      queue_metrics: metrics,
      timestamp: new Date().toISOString()
    };
  }
}

module.exports = { DirectorAgent };
