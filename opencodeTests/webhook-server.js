#!/usr/bin/env node
/**
 * OpenCode Agency — Webhook Server
 * Receives external events (GitHub, Discord, cron) and transforms them into tasks
 */

const http = require('http');
const { verify: verifySignature } = require('crypto').webcrypto.subtle;
const { TaskQueue } = require('./task-queue');
const pino = require('pino');

class WebhookServer {
  constructor(config, queue) {
    this.config = config;
    this.queue = queue;
    this.logger = pino({ level: process.env.LOG_LEVEL || 'info' });
    this.port = config.agency.webhooks.port || 9091;
    this.secret = config.agency.webhooks.secret || '';
    this.routes = config.agency.webhooks.routes || {};
    this.server = null;
  }

  async start() {
    this.server = http.createServer(async (req, res) => {
      if (req.method !== 'POST') {
        res.writeHead(405);
        return res.end('Method Not Allowed');
      }

      const url = new URL(req.url, `http://localhost:${this.port}`);
      const routeKey = url.pathname;

      if (!this.routes[routeKey]) {
        res.writeHead(404);
        return res.end('Not Found');
      }

      let body = '';
      req.on('data', chunk => body += chunk);
      req.on('end', async () => {
        try {
          // Verify signature if secret configured (e.g., GitHub)
          if (this.secret && req.headers['x-hub-signature-256']) {
            const signature = req.headers['x-hub-signature-256'].split('=')[1];
            const digest = await this.calculateHmac(body, this.secret);
            const expected = Buffer.from(signature, 'hex').toString('hex');
            if (digest !== expected) {
              this.logger.warn('Invalid webhook signature');
              res.writeHead(401);
              return res.end('Invalid Signature');
            }
          }

          const payload = JSON.parse(body);
          await this.handleEvent(this.routes[routeKey], payload, req.headers);
          res.writeHead(200);
          res.end('OK');
        } catch (err) {
          this.logger.error({ error: err.message, body: body.slice(0, 500) }, 'Webhook handling error');
          res.writeHead(500);
          res.end('Internal Server Error');
        }
      });
    });

    await new Promise(resolve => this.server.listen(this.port, 'localhost', resolve));
    this.logger.info({ port: this.port }, 'Webhook server listening');
  }

  async calculateHmac(data, secret) {
    const encoder = new TextEncoder();
    const key = await this.webcrypto.subtle.importKey(
      'raw',
      encoder.encode(secret),
      { name: 'HMAC', hash: 'SHA-256' },
      false,
      ['sign']
    );
    const signature = await this.webcrypto.subtle.sign(
      'HMAC',
      key,
      encoder.encode(data)
    );
    return Buffer.from(signature).toString('hex');
  }

  async handleEvent(routeType, payload, headers) {
    this.logger.info({ routeType, event: payload.event || 'unknown' }, 'Processing webhook');

    if (routeType === 'github_issue') {
      // GitHub issue opened
      if (payload.action === 'opened' && payload.issue) {
        const task = {
          type: 'director.plan',
          payload: {
            description: `GitHub Issue #${payload.issue.number}: ${payload.issue.title}`,
            source: 'github',
            issue_url: payload.issue.html_url,
            body: payload.issue.body,
            priority: 5
          },
          priority: 5
        };
        this.queue.enqueue(task);
        this.logger.info({ issue: payload.issue.number }, 'Enqueued from GitHub issue');
      }
    } else if (routeType === 'manual_task') {
      // Direct POST /webhook/task with custom payload
      const task = {
        type: payload.type || 'director.plan',
        payload: payload.payload || payload,
        priority: payload.priority || 0
      };
      this.queue.enqueue(task);
      this.logger.info('Enqueued manual task');
    } else {
      this.logger.warn({ routeType }, 'Unhandled webhook route');
    }
  }

  stop() {
    if (this.server) this.server.close();
  }
}

module.exports = { WebhookServer };
