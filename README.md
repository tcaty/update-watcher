## Update watcher

Always stay up to date!

## How it works?

1. Specify targets and webhooks.
2. Wait for new version being released.
3. Receive notifications.

![image](https://github.com/tcaty/update-watcher/assets/79706809/2fb24ae8-0016-4da6-b5c5-809c0e771622)

4. Update it!

## How to use it?

1. Run `make prepare` from repo root. This command will create `.env` files in `deploy/*` folders.
2. Go to `deploy/*` folders, write your secrets to `.env` files guided by `.env.example`.
3. Run `make run-dev` or `make run-prod` from repo root.
4. Always stay up to date!
5. Remove all created resources by command `make clean`.
