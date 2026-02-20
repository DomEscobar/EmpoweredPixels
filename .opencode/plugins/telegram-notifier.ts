import { onTaskComplete, onTaskFail } from '@opencode/agent';

const TELEGRAM_TOKEN = process.env.TELEGRAM_BOT_TOKEN;
const TELEGRAM_CHAT_ID = process.env.TELEGRAM_CHAT_ID;

async function sendTelegram(message: string) {
  if (!TELEGRAM_TOKEN || !TELEGRAM_CHAT_ID) return;
  try {
    await fetch(`https://api.telegram.org/bot${TELEGRAM_TOKEN}/sendMessage`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        chat_id: TELEGRAM_CHAT_ID,
        text: message,
        parse_mode: 'Markdown'
      })
    });
  } catch (e) {
    console.error('Failed to send Telegram notification', e);
  }
}

onTaskComplete(async (task) => {
  await sendTelegram(`✅ *Task Completed*\n\n*ID:* ${task.id}\n*Summary:* ${task.summary || task.description}`);
});

onTaskFail(async (task, error) => {
  await sendTelegram(`❌ *Task Failed*\n\n*ID:* ${task.id}\n*Error:* ${error.message}`);
});
