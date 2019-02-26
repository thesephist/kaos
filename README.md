# kaos

Kaos is an experimental to-do list manager using an open text file format.

## Synopsis

Kaos is an experimental testbed for various ideas I'm thinking about incorporating into my full-time task and project manager, [Sigil](https://github.com/thesephist/sigil).

## Usage

Initialize a Kaos database by creating a file called `kaosfile` in your project directory.

```
$ touch ./kaosfile
```

Now, the `kaos` CLI will save tasks to the file, and reference it to keep track of tasks. For example, to add a new task, run `kaos create`, and the CLI will guide you through a few questions to create a new task.

```
$ kaos create

> Project?
testproject

> Size?
3

> Description?
Something I really need to do

Created:
#kqmuvqlsnw [2019/02/23T22:04:05|-|-] @-
testproject     (3): Something I really need to do
```

See all of your tasks with `kaos list` (this is also the default command).

All tasks in Kaos are referenced by a string ID of lowercase alphabetical characters (why? Because it makes it quick to type in a terminal). You can reference unique substrings of these IDs (`Ref`s) to start, complete, remove, undo, and update tasks through the CLI.

- `kaos find <query>`: Show a list of tasks that match the query as a strict substring
- `kaos remove kqmuv`: Remove the task whose Ref contains the substring `kqmuv`
- `kaos due kqmuv`: Set a due date for this task
- `kaos project kqmuv`: Update the task project name (always a sing lexical word, no whitespace)
- `kaos size kqmuv`: Adjust the marked task size (how much effort it requires)
- `kaos describe kqmuv`: Update the task description
- `kaos start kqmuv`: Mark the task as started (will add a started timestamp to the task)
- `kaos finish kqmuv`: Mark the task as finished (will add a finished timestamp to the task)

## File format

Kaos uses an open file format that's human-readable to keep track of tasks. This has several key advantages

- The data is human-readable. It's just a file, and I can open it and read it and know exactly where I am on a project with a few `grep` and `sed` chained together.
- The data is not tied to the tool. Although `kaos` the CLI makes managing the file easier, it's not a necessity to interact with the data -- you can just use your text editor or write an alternative CLI. This openness is important to me as a user and a developer.
- The data is trivial to back up / share / replicate. It's just a text file.
- The data can be easily and esnsibly version controlled using e.g. `git`. Having a human-readable format also means version control tools for source code work well with the format. Although `kaos` has some history-tracking built in via timestamps on task updates, VCS work well with Kaos's file format.

The exact syntax for serializing a task into a `kaosfile` is still in flux, but should be easy to find in `pkg/kaos/task.go` as `Task.Parse() and Task.String()`.

