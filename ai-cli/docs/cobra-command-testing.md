# Testing cobra commands in one process — shared-tree state leakage

> Status: **Note** · 2026-09-13 · Scope: the table-driven CLI test in
> `ai-cli/cmds/root_test.go` (`TestCmd`) and the shared global command tree it drives.

## The symptom

`TestCmd` runs a list of cases, each calling `Execute()` on the shared global `cmds.RootCmd`.
Observed:

- With the **help** group enabled (`prompt execute --help`) *and* the **prompt** group enabled, a
  breakpoint in the prompt command's `Run` (or `validateArgs`/`doWork`) is **never hit** — yet the
  test stays green.
- Disable the help group → the breakpoint **is** hit.

## Why it happens

The whole suite runs in **one process**, so every case shares the same `cmds.RootCmd` tree. Cobra
and pflag **do not reset flag state between `Execute()` calls**: pflag only *sets* flags present in
the new args; it never clears flags a previous parse set.

The concrete leak is cobra's auto-added `--help` bool:

1. The help case runs `prompt execute --help` → cobra sets the `help` flag to `true`.
2. Nothing resets it. On the next `Execute()` (the prompt case), the args don't contain `--help`,
   so `help` stays `true`.
3. Early in execution cobra checks that flag and, when true, **prints help and returns without
   calling the command's `Run`**. So the breakpoint is never reached.
4. The help path returns `nil`, so `require.NoError` passes — which is why the test looks green and
   the problem is easy to miss.

With the help group disabled, `--help` is never set on the shared tree, so the prompt case runs
normally. That asymmetry is the tell: it's leaked flag state, not a problem with the prompt case.

`--help` is not the only sticky state. **Every** flag set in one case keeps its value *and* its
`Changed=true` on the shared tree. That matters for `MarkFlagRequired`: cobra treats a flag as
"provided" when `Changed` is true, so if an earlier case set `--llm`, a later case that omits it
would **wrongly pass** the required check. Today's two cases don't collide (the help case sets none
of the prompt flags), but two prompt-style cases would.

## Current fix (kept) — targeted help reset

Clear the sticky `--help` across the whole command tree before each `Execute()`:

```go
// resetHelpFlag clears cobra's sticky --help bool across the whole command tree.
// cobra adds it lazily to the executed command, so reset every command's copy.
func resetHelpFlag(c *cobra.Command) {
	if f := c.Flags().Lookup("help"); f != nil {
		_ = f.Value.Set("false")
		f.Changed = false
	}
	for _, sub := range c.Commands() {
		resetHelpFlag(sub)
	}
}

// ... inside the per-case subtest, before running:
resetHelpFlag(cmds.RootCmd)
cmds.RootCmd.SetArgs(tArgs)
err := cmds.RootCmd.Execute()
require.NoError(t, err)
```

This addresses the reported symptom. If the suite grows to multiple flag-setting cases and the
`Changed`-leak above starts to bite, generalize it to reset **all** flags (value → default,
`Changed = false`) across the tree before each run. Caveat: resetting a custom `pflag.Value` (e.g.
the `enumValue` used by `--llm`) to its empty default triggers that value's validation error on the
reset call — harmless, because the subsequent parse overwrites it, but worth a comment so it isn't
mistaken for a real failure.

## The structural alternative — a `NewCmd` factory

The reset works by *undoing* leaked state. The alternative removes the leak **by construction**:
don't reuse a global tree — build a fresh one per run.

### What we do today (global + `init()`)

- `cmds` has a package-level `var RootCmd = &cobra.Command{...}` — one instance for the process.
- Each command package has a package-level `var theCmd` and an `init()` that calls
  `cmds.RootCmd.AddCommand(theCmd)` and registers flags.
- A blank import (`_ ".../cmds/prompt"`) runs that `init()`, auto-attaching the command.

Result: exactly one tree, created once, shared by every case — state accumulates on it.

### Factory shape

A factory is a constructor that returns a **new** tree each call:

```go
// package cmds
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "ai-cli",
		Short: "...",
	}
	// explicit wiring instead of init() auto-registration:
	root.AddCommand(prompt.NewCmd())
	// root.AddCommand(otherpkg.NewCmd()) ...
	return root
}
```

```go
// package prompt — replace the global `var theCmd` + init() with a constructor.
// Flag-bound state moves into a struct so each built command owns its own copy,
// rather than writing into package globals.
type options struct {
	model                string
	outFolder            string
	llmProvider          string
	promptDefinitionFile string
	cfgFileName          string
	batch                bool
	batchPollInterval    time.Duration
	promptVars           map[string]string
}

func NewCmd() *cobra.Command {
	o := &options{}

	c := &cobra.Command{
		Use:   "prompt",
		Short: "invokes an llm to execute a prefilled prompt",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.validate(); err != nil {
				return err
			}
			return o.run()
		},
	}

	c.Flags().StringVarP(&o.model, "model", "M", "", "...")
	c.Flags().StringVarP(&o.outFolder, "out-folder", "O", "", "...")
	c.Flags().StringVarP(&o.promptDefinitionFile, "prompt", "P", "", "...")
	c.Flags().VarP(newEnumValue(&o.llmProvider, "", "anthropic", "ollama", "vllm"),
		"llm", "L", "provider to use (one of: anthropic, ollama, vllm)")
	c.Flags().StringVarP(&o.cfgFileName, "cfg-file", "C", "", "...")
	c.Flags().BoolVarP(&o.batch, "batch", "B", false, "...")
	c.Flags().DurationVar(&o.batchPollInterval, "poll-interval", 30*time.Second, "...")
	c.Flags().StringToStringVarP(&o.promptVars, "var", "V", nil, "...")

	for _, name := range []string{"prompt", "out-folder", "llm", "cfg-file", "model"} {
		_ = c.MarkFlagRequired(name)
	}
	return c
}
```

The test then builds its own tree each iteration — no reset needed:

```go
root := cmds.NewRootCmd()
root.SetArgs(tArgs)
err := root.Execute()
```

Every iteration gets a fresh `*cobra.Command` with flags freshly registered at their defaults, so
there is no sticky `--help`, no stale `Changed`, and no `MarkFlagRequired` cross-talk.

### Trade-offs

- **More work:** convert *every* command package from `init()`-registers-a-global to a `NewCmd()`
  constructor — a structural change across the CLI, not just the test.
- **Lose auto-registration:** today a blank import wires a command in; a factory lists subcommands
  explicitly in `NewRootCmd()`. Often considered cleaner/more testable, but it is a style change.
- **Flag-bound variables:** moving them from package globals into an `options` struct (as above)
  is what gives true per-build isolation; leaving them as globals mostly works for sequential tests
  (re-registration resets them) but keeps process-global state.

## Decision

Keep the **targeted help reset** for now — it fixes the observed symptom with a few lines and no
restructuring. The `NewCmd` factory is the right move *if/when* the CLI itself needs to be
re-instantiable (broader test isolation, embedding it in another program, running it twice in one
process), or when the `Changed`-leak starts affecting multi-case runs. It's an architecture choice,
not a bug fix.
