# Demo

Example files for trying out Atheon. Every secret in here is fake and intentionally structured to trigger a pattern so you can see Atheon in action.

## Files

| File | What it contains |
|------|-----------------|
| `config.yaml` | App config with a fake OpenAI key and AWS key buried inside |
| `server.go` | Clean file — no findings, intentional contrast |
| `logs/app.log` | Fake application log with a Stripe key and Twilio SID inside an error trace |

## Try it

Scan the whole demo folder:
```
atheon ./demo
```

Pipe the log file directly:
```
cat demo/logs/app.log | atheon -
```

## Pre-commit hook

To wire Atheon into git so it blocks commits automatically, copy the hook from the `hooks/` folder at the root of this repo:

```
cp hooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

Every `git commit` will now run Atheon first. If anything is found, the commit is blocked.

Try staging `demo/config.yaml` and committing — it will be caught before it ever lands.
