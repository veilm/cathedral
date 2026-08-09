Perfect. In the API documentation, ignore any large sections about 'OAuth 2.0 Flow' for now. We're looking for something much simpler, usually under a heading like 'Authentication', 'Personal Access Tokens', or 'API Keys'.

Your goal is to find the page that lets you generate a 'token' or 'key' for your own personal use. It's usually a long string of random characters. Many services have a dedicated section in your account settings for this.

Once you've generated a key, we'll need two pieces of information from the docs:
1. The **URL** for the tasks endpoint (e.g., `https://api.todoapp.com/v1/tasks`).
2. The name of the **HTTP Header** used for the key (e.g., `Authorization: Bearer YOUR_KEY_HERE` or `X-Api-Key: YOUR_KEY_HERE`).

Find that key and those two details, and we'll assemble the command.
