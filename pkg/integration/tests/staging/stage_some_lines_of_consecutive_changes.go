package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageSomeLinesOfConsecutiveChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage only some lines of a block of consecutive changes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "1\n2\n3\n4\n5\n6\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "1\n2b\n3b\n4b\n5b\n6\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().Staging().
			IsFocused().
			ContainsLines(
				Contains(" 1"),
				Contains("-2"),
				Contains("-3"),
				Contains("-4"),
				Contains("-5"),
				Contains("+2b"),
				Contains("+3b"),
				Contains("+4b"),
				Contains("+5b"),
				Contains(" 6"),
			).
			SelectedLines(Contains("-2")).
			Press(keys.Universal.RangeSelectDown).
			PressPrimaryAction().
			NavigateToLine(Contains("+2b")).
			Press(keys.Universal.RangeSelectDown).
			PressPrimaryAction()

		t.Views().StagingSecondary().
			ContainsLines(
				/* EXPECTED:
				Contains(" 1"),
				Contains("-2"),
				Contains("-3"),
				Contains("+2b"),
				Contains("+3b"),
				Contains(" 4"),
				Contains(" 5"),
				Contains(" 6"),
				ACTUAL: */
				Contains(" 1"),
				Contains("-2"),
				Contains("-3"),
				Contains(" 4"),
				Contains(" 5"),
				Contains("+2b"),
				Contains("+3b"),
				Contains(" 6"),
			)
	},
})
