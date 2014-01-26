The Great Gitsby
================

[deploy]: https://help.github.com/articles/managing-deploy-keys#deploy-keys

Handle git commit hooks like a champ.

The general idea behind Gitsby is similar to that of danneu's
[captain-githook](https://github.com/danneu/captain-githook) - you provide a
config file (either explicity through the `-config=file` flag, or
`~/gitsby/gitsby.json` by default) and Gitsby will handle commit hooks from
GitHub for you.

Gitsby will bind by default to `0.0.0.0:9999` - you can change this via the
`-host=0.0.0.0` and `-port=9999` flags.

# Usage

* Create `~/gitsby`

```
~/
└── gitsby/
    └── config.json
```

* List all of your repos in `~/gitsby/config.json`:

```json
{
	"landing": true,
	"repos": [
		{
			"url": "https://github.com/username/repo.git"
		},
		{
			"url": "git@github.com:username/other-repo.git",
			"directory": "~/somewhere-fun",
			"hidden": true
		}
	]
}
```

Repos MUST have at least a `url` (NB: if you want to deploy private repos,
check out [GitHub:Help 'Managing deploy keys'][deploy] and `ssh-agent`/
`ssh-add`), but can optionally contain `directory` (where it will be clone'd
to) and `hidden` (whether it's hidden on the - currently non-implemented -
landing page).

* Install The Great Gitsby:

```
$ go get github.com/plausibility/gitsby
...
$ go install github.com/plausibility/gitsby
...
$ gitsby -h
Usage of gitsby:
  -config="/home/<user>/gitsby/config.json": Gitsby config file
  -host="0.0.0.0": host to bind web.go to
  -port=9999: port to bind web.go to
```

* Throw a party:

```
$ gitsby
2014/01/26 17:31:10 The Great Gitsby is throwing a party!
2014/01/26 17:31:10 ➜ Preparing 2 repo(s) for sync
2014/01/26 17:31:14 ✔ [username/repo] synced
Already up-to-date.
2014/01/26 17:31:14 ✖ [username/repo] failed to 'make autodeploy'! Here's why:
make: *** No rule to make target 'autodeploy'.  Stop.
2014/01/26 17:31:14 ✔ [username/other-repo] synced
Already up-to-date.
2014/01/26 17:31:14 ✖ [username/other-repo] failed to 'make autodeploy'! Here's why:
make: *** No rule to make target 'autodeploy'.  Stop.
2014/01/26 17:31:14 The party is here: 0.0.0.0:9999
```

* Add a WebHook URL to your GitHub repository.

	* Make sure it's `http://example.com:9999/github` (for example); note the
	trailing `/github`!

* GitHub will notify Gitsby when repos need redeployment, and he will:

	* Run `git pull origin` in the appropriate repository.
	* Run `make autodeploy` in the repository root, which you can optionally
	use to automagically deploy your project.

# Why make this?

I've been meaning to learn Go for several months, but just haven't had a good
reason to. In these months, [Aki](https://github.com/aki--aki) has taken it
upon herself to berate me for not having an autodeploy setup for the DongCorp
website repo. Seeing as my friend danneu had made _captain-githook_, I figured
I'd try my hand at an autodeploy script, and that it was the perfect chance
to finally learn Go.

TL;DR: It was to scratch my own back.
