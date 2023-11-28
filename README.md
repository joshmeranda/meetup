# Meetup

A simple command line interface to managing your meeting notes.

## Installing

As of now the only way to instlal is to clone and build from source:

```
>>> git clone https://github.com/joshmeranda/meetup.git"
>>> cd meetup
>>> make install
```

## Ovedrview

todo: meeting path builder
todo: make date required for removing meetings

### Configuration

The main driver behind meeting up is the Manager, which handles all of the file management. The manager reads its configuration both from command line arguments as well from a configuration foudn in your user's default configuration directory under `meetup/config.yaml`. This configuration handles high-lvel concepts such as what editor to use and default domains.

| key                | type       | default       | description                                                          |
|--------------------|------------|---------------|----------------------------------------------------------------------|
| `root_dir`         | string     | $HOME/.meetup | The local directory where meetup meetings are stored.                |
| `editor`           | []string   | $EDITOR       | The command to use to open files.                                    |
| `default_metadata` | Metadata   |               | Override the default meetup metadata when creating a new meetup dir. |

Some values you can only configure at the metup directory level (eg GroupBy). These can be found at `<meetup_dir>/.metadata`:

| key        | type   | default | description                                                                                                                                            |
|------------|--------|---------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| `group_by` | string | domain  | Specify how to group each meeting. Must be one of `date`, or `domain`. NOTE: do not change this manually, change via `meetup meeting group-by` instead |

### Meetings

Opening a new or existing meeting can be done through the `open` subcommand. Note that the `--date` defaults to the current date, so when opening an old meeting, be sure to provide the right date.

Once meetings are created, you can view your meetings with the `list` subcommand. You can provide various filters on the date, domain, and name as simple wildcards.

If you decide you no longer need the notes you made for a meeting you can remove it with the `remove` subcommand.

See below for various examples on handling meetings:

 - Open a new meeting for today

```
meetup open work.product.team scheduling
```

 - Open a meeting from the early Januruary 23rd 2001

```
meetup open --date 2001-01-23 work.product.team scheduling
```

 - View all work meetings from 2010 with "frog" in the name

```
meetup list --date '2010-*' --domain 'work.*' --name '*frog*'
```

 - Remove that meeting from before

```
meetup remove --date 2001-01-23 work.product.team scheduling
```

### Templates

Meetup allows you to create templates that you can create meetings from. These templates should be in the form of go templates. You have access to all fields of `meetup.Meeting`. See [./examples/templates]() for examples of
templates. Managing templates can be done with the `add`, `list`, and `remove` subcommands. For exapmles, see below:

 - Add a file as a template

```
meetup template add ./examples/templates/simple.md
```

 - List available templates:

```
meetup template list
```

 - Create a meeting from a template you can simply provide the name of the template when opening the meeting

```
meetup open --template simple.md work.product.team scheduling
```

 - Remove uneeded templates

 ```
 meetup template remove simple.md
 ```

### Todos / Tasks

Meetup also provides some basic support for tracking tasks accross meetings. To do this, we use the markdown task list syntax:

```
# Meeting Tasks

 - [ ] make schedule
 - [ ] distribute schedule
 - [x] walk the office dog
```

Accessing the list above can be done with `meetup task` (or `meetup todo`). The command should render output like below:

```
[2023-11-27 tasks.test example] ❌ make schedule
[2023-11-27 tasks.test example] ❌ distribute schedule
[2023-11-27 tasks.test example] ✅ walk the office dog
```

See `meetpup task --help` for more details.