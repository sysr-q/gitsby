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
    └── gitsby.json
```

* List all of your repos in `~/gitsby/gitsby.json`:

```json
{
  "repos": [
    {
      "url": "https://github.com/username/repo.git",
      "command": {"cmd": "make", "args": ["bar", "baz"]}
    },
    {
      "url": "git@github.com:username/other-repo.git",
      "directory": "~/somewhere-fun",
      "hidden": true,
      "command": {"cmd": "./deploy", "args": ["all"]}
    }
  ]
}
```

Repos MUST have at least a `url` (NB: if you want to deploy private repos,
check out [GitHub:Help 'Managing deploy keys'][deploy] and `ssh-agent`/
`ssh-add`), and a `command` (e.g. `{"cmd": "make", "args": ["all", "deploy"]}`),
but can optionally contain `directory` (where it will be clone'd to) and
`hidden` (whether it's hidden on the - currently non-implemented - landing
page). They MAY also include `type`, but that's not implemented or checked yet
- it may in the future if I implement other repo types than git, however.

* Install The Great Gitsby:

```
$ go get github.com/sysr-q/gitsby
...
$ go install github.com/sysr-q/gitsby
...
$ gitsby -h
Usage of gitsby:
  -config="/home/<user>/gitsby/gitsby.json": Gitsby config file
  -host="0.0.0.0": host to bind web.go to
  -port=9999: port to bind web.go to
```

* Throw a party:

```
$ gitsby
2014/03/22 17:29:44 The Great Gitsby is sending invites to 2 repos.
2014/03/22 17:29:50 [username/repo] successfully synced: (stdout)
Already up-to-date.
2014/03/22 17:29:50 [username/repo] failed to deploy! Here's why: (stderr)
make: *** No rule to make target 'autodeploy'.  Stop.
2014/03/22 17:29:50 [username/other-repo] doesn't exist, syncing!
2014/03/22 17:29:59 [username/other-repo] successfully cloned to: /home/you/somewhere-fun (stdout)
2014/03/22 17:29:59 [username/other-repo] failed to deploy! Here's why: (stderr)
make: *** No rule to make target 'autodeploy'.  Stop.
2014/03/22 17:29:59 The party is here: 0.0.0.0:9999
```

* Add a WebHook URL to your GitHub repository.

    * Make sure it's `http://example.com:9999/github` (for example); note the
    trailing `/github`!

* GitHub will notify Gitsby when repos need redeployment, and he will:

    * Run `git pull origin` in the appropriate repository.
    * Run the specified command in the repository root, which you can use to
    automagically deploy your project.
